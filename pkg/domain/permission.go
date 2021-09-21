package domain

import (
	"context"
	"time"
)

var (
	TableNamePermission = "permissions"
)

type Permission struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Namespace   string    `gorm:"column:namespace;type:string;size:256;uniqueIndex:uniq_code;not null"`
	Code        string    `gorm:"column:name;type:string;size:32;uniqueIndex:uniq_code;not null"`
	AccountID   uint64    `gorm:"column:account_id;type:bigint;uniqueIndex:uniq_code;not null"`
	CreatorID   uint64    `gorm:"column:creator_id;type:bigint;not null"`
	CreatorName string    `gorm:"column:creator_name;type:string;size:32;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
}

func (p *Permission) TableName() string {
	return TableNamePermission
}

type PermissionUsecase interface {
	PermissionsByAccountID(ctx context.Context, namespace string, accountID uint64) ([]Permission, error)
	UpdatePermissions(ctx context.Context, namespace string, permissions []*Permission) error
}

type PermissionRepository interface {
	PermissionsByAccountID(ctx context.Context, namespace string, accountID uint64) ([]Permission, error)
	DeletePermissionsByAccountID(ctx context.Context, namespace string, accountID uint64) error
	CreatePermissions(ctx context.Context, permissions []*Permission) error
}
