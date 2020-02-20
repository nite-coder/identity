package main

import (
	"fmt"
	"identity/internal/pkg/config"
	identityGRPC "identity/pkg/identity/delivery/grpc"
	identityDatabase "identity/pkg/identity/repository/database"
	identitySvc "identity/pkg/identity/service"
	"time"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"

	"github.com/cenk/backoff"
	"github.com/jinzhu/gorm"
)

var (
	_identityServer *identityGRPC.IdentityServer
)

func initialize(cfg config.Configuration) error {
	initLogger("identity", cfg)

	db, err := setupDatabase(cfg)
	if err != nil {
		return err
	}

	accountRepo := identityDatabase.NewAccountRepo(db)
	accountSvc := identitySvc.NewAccountService(accountRepo)

	_identityServer = identityGRPC.NewIdentityServer(accountSvc)

	return nil

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
