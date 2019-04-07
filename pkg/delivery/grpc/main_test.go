package grpc

import (
	"fmt"
	"identity/internal/config"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/cenk/backoff"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
	"github.com/jinzhu/gorm"

	identityDatabase "identity/pkg/repository/database"
	identitySvc "identity/pkg/service"

	_ "github.com/go-sql-driver/mysql"
)

var (
	_identityServer *IdentityServer
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println("Current test filename: " + filename)

	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	dir = filepath.Join(dir, "../../")
	fmt.Println(dir)

	// setup test env
	config := config.Configuration{
		Logs: []config.LogSetting{
			{
				Name:     "console",
				Type:     "console",
				MinLevel: "info",
			},
		},
	}

	config.Database.Address = "localhost:3306"
	config.Database.Username = "test"
	config.Database.Password = "test"
	config.Database.DBName = "identity_db"

	initLogger("identity_test", &config)

	db, err := setupDatabase(&config)
	if err != nil {
		panic(err)
	}

	accountRepo := identityDatabase.NewAccountRepo(db)
	accountSvc := identitySvc.NewAccountService(accountRepo)

	_identityServer = NewIdentityServer(accountSvc)

	// start testing
	m.Run()
}

func initLogger(appID string, config *config.Configuration) {
	// set up log target
	log.SetAppID(appID)
	for _, target := range config.Logs {
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

func setupDatabase(config *config.Configuration) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", config.Database.Username, config.Database.Password, config.Database.Address, config.Database.DBName)

	log.Debugf("main: database connection string: %s", connectionString)
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
