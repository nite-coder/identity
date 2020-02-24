package grpc

import (
	"context"
	"fmt"
	"identity/internal/pkg/config"
	"net"
	"os"
	"testing"
	"time"

	identityProto "identity/pkg/identity/proto"
	identityDatabase "identity/pkg/identity/repository/database"
	identitySvc "identity/pkg/identity/service"

	"github.com/cenk/backoff"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	_identityServer *IdentityServer
	_identityClient identityProto.IdentityServiceClient
	_lis            *bufconn.Listener
)

const _bufSize = 1024 * 1024

func initialize(cfg config.Configuration) error {
	initLogger("identity", cfg)

	db, err := setupDatabase(cfg)
	if err != nil {
		return err
	}

	accountRepo := identityDatabase.NewAccountRepo(db)
	accountSvc := identitySvc.NewAccountService(accountRepo)

	_identityServer = NewIdentityServer(accountSvc)

	ctx := context.Background()
	clientConn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}

	_identityClient = identityProto.NewIdentityServiceClient(clientConn)

	return nil

}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return _lis.Dial()
}

func initLogger(appID string, cfg config.Configuration) {
	// set up log target
	defaultFeilds := log.Fields{
		"app_id": appID,
		"env":    cfg.Env,
	}
	log.WithDefaultFields(defaultFeilds)

	for _, target := range cfg.Logs {
		switch target.Type {
		case "console":
			clog := console.New()
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.RegisterHandler(clog, levels...)
		case "gelf":
			graylog := gelf.New(target.ConnectionString)
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.RegisterHandler(graylog, levels...)
		}
	}
}

func setupDatabase(config config.Configuration) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", config.Database.Username, config.Database.Password, config.Database.Address, config.Database.DBName)

	var db *gorm.DB
	var err error
	err = backoff.Retry(func() error {
		db, err = gorm.Open("mysql", connectionString)
		if err != nil {
			log.Errorf("main: mysql open failed: %v", err)
			return err
		}
		err = db.DB().Ping()
		if err != nil {
			log.Errorf("main: mysql ping error: %v", err)
			return err
		}
		return nil
	}, bo)

	if err != nil {
		log.Panicf("main: mysql connect err: %s", err.Error())
	}

	log.Infof("database ping success")
	db.DB().SetMaxIdleConns(150)
	db.DB().SetMaxOpenConns(300)
	db.DB().SetConnMaxLifetime(14400 * time.Second)

	return db, nil
}

func TestMain(m *testing.M) {
	fmt.Println("===begin===")

	config.EnvPrefix = "IDENTITY"
	cfg := config.New("app.yml")
	err := initialize(cfg)
	if err != nil {
		log.Panicf("main: initialize failed: %v", err)
		return
	}

	// start grpc server
	_lis = bufconn.Listen(_bufSize)
	if err != nil {
		log.Fatalf("main: bind identity grpc failed: %v", err)
	}
	grpcServer := grpc.NewServer()

	identityProto.RegisterIdentityServiceServer(grpcServer, _identityServer)
	log.Info("main: grpc service started")
	go func() {
		if err = grpcServer.Serve(_lis); err != nil {
			log.Fatalf("main: failed to start grpc server: %v", err)
		}
	}()

	code := m.Run() // run all tests
	grpcServer.GracefulStop()

	fmt.Println("===end===")
	os.Exit(code)
}
