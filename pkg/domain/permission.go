package domain

import "time"

type Permission struct {
	ID          int64     `gorm:"column:id;primaryKey;not null"`
	Namespace   string    `gorm:"column:namespace;type:string;size:256;uniqueIndex:uniq_name;not null"`
	Name        string    `gorm:"column:name;type:string;size:32;uniqueIndex:uniq_name;not null"`
	GroupID     uint32    `gorm:"column:group_id;type:int;not null"`
	Version     uint32    `gorm:"column:version;type:int;not null"`
	CreatorID   uint64    `gorm:"column:creator_id;type:bigint;not null"`
	CreatorName string    `gorm:"column:creator_name;type:string;size:32;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
	UpdaterID   uint64    `gorm:"column:updater_id;type:bigint;not null"`
	UpdaterName string    `gorm:"column:updater_name;type:string;size:32;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
}

func (p *Permission) TableName() string {
	return "permissions"
}
