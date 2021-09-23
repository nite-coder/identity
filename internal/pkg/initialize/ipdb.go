package initialize

import (
	"os"
	"path/filepath"

	"github.com/oschwald/geoip2-golang"
)

func InitGeoIPDB() (*geoip2.Reader, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(path, "geolite.mmdb")

	db, err := geoip2.Open(filePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
