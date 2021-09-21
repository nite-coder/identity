package initialize

import (
	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/config/provider/file"
)

func InitConfig() error {
	fileProvder := file.New()

	err := fileProvder.Load()
	if err != nil {
		return err
	}

	config.AddProvider(fileProvder)

	return nil
}
