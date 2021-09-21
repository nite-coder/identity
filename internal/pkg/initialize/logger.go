package initialize

import (
	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/nite-coder/blackbear/pkg/log/handler/console"
	"github.com/nite-coder/blackbear/pkg/log/handler/gelf"
)

type LogSetting struct {
	Name             string
	Type             string
	MinLevel         string `mapstructure:"min_level"`
	ConnectionString string `mapstructure:"connection_string"`
}

func InitLogger() error {
	logSettings := []LogSetting{}
	err := config.Scan("log", &logSettings)
	if err != nil {
		return err
	}

	logger := log.New()

	for _, logSetting := range logSettings {
		switch logSetting.Type {
		case "console":
			clog := console.New()
			levels := log.GetLevelsFromMinLevel(logSetting.MinLevel)
			logger.AddHandler(clog, levels...)
		case "gelf":
			graylog := gelf.New(logSetting.ConnectionString)
			levels := log.GetLevelsFromMinLevel(logSetting.MinLevel)
			logger.AddHandler(graylog, levels...)
		}
	}

	log.SetLogger(logger)
	return nil
}
