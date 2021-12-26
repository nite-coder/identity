package global

import "gorm.io/gorm"

const (
	VersionAddOne = "version + 1"
)

var (
	DB *gorm.DB
)
