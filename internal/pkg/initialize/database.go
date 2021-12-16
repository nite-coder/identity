package initialize

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Name             string
	ConnectionString string `mapstructure:"connection_string"`
	Type             string
	Migration        bool `mapstructure:"migration"`
	EnableLogger     bool `mapstructure:"enable_logger"`
}

func InitDatabase(name string) (*gorm.DB, error) {
	databases := []Database{}
	err := config.Scan("database", &databases)
	if err != nil {
		return nil, err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(30) * time.Second

	for _, database := range databases {

		log.Str("connection_string", database.ConnectionString).Debug("database is initialing.")

		if strings.EqualFold(database.Name, name) {

			// migrate database if needed
			if database.Migration {
				path := fmt.Sprintf("./deployment/%s/%s", database.Type, database.Name)
				path = filepath.ToSlash(path) // due to migrate package path issue on window os, therefore, we need to run this
				source := fmt.Sprintf("file://%s", path)
				migrateDBURL := fmt.Sprintf("%s://%s", database.Type, database.ConnectionString)

				m, err := migrate.New(
					source,
					migrateDBURL,
				)
				if err != nil {
					// make sure migration package is using v4 above
					return nil, fmt.Errorf("db migration config was wrong. db_name: %s, source: %s, migrateDBURL: %s, error: %w", database.Name, source, migrateDBURL, err)
				}

				err = m.Up()
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return nil, fmt.Errorf("db migration failed. db: %s, source: %s, migrateDBURL: %s, error: %w", database.Name, source, migrateDBURL, err)
				}

				log.Infof("%s database was migrated", database.Name)
			}

			var db *gorm.DB
			var err error
			err = backoff.Retry(func() error {

				gormConfig := gorm.Config{
					//PrepareStmt: true,
					Logger: logger.Default.LogMode(logger.Info),
				}

				switch strings.ToLower(database.Type) {
				case "mysql":
					db, err = gorm.Open(mysql.Open(database.ConnectionString), &gormConfig)
				case "postgres":
					db, err = gorm.Open(postgres.New(postgres.Config{
						DSN:                  database.ConnectionString,
						PreferSimpleProtocol: false, // disables implicit prepared statement usage
					}), &gormConfig)
				}

				if err != nil {
					return fmt.Errorf("startup: database open failed: %w", err)
				}

				sqlDB, err := db.DB()
				if err != nil {
					return err
				}

				sqlDB.SetMaxIdleConns(150)
				sqlDB.SetMaxOpenConns(300)
				sqlDB.SetConnMaxLifetime(14400 * time.Second)

				err = sqlDB.Ping()
				if err != nil {
					return fmt.Errorf("startup: database ping failed. name: %s, error: %w", database.Name, err)
				}

				return nil
			}, bo)

			if err != nil {
				return nil, fmt.Errorf("startup: database connect failed.  connection_string: %s, error: %w", database.ConnectionString, err)
			}

			return db, nil
		}
	}

	return nil, fmt.Errorf("startup: database name was not found. name: %s", name)

}
