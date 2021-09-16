package domain

import (
	"context"
	"time"
)

// Role 代表角色資訊
type Role struct {
	ID          int64     `gorm:"column:id;primaryKey;not null"`
	Namespace   string    `gorm:"column:namespace;type:string;size:256;uniqueIndex:uniq_name;not null"`
	Name        string    `gorm:"column:name;type:string;size:32;uniqueIndex:uniq_name;not null"`
	Version     uint32    `gorm:"column:version;type:int;not null"`
	CreatorID   uint64    `gorm:"column:creator_id;type:bigint;not null"`
	CreatorName string    `gorm:"column:creator_name;type:string;size:32;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
	UpdaterID   uint64    `gorm:"column:updater_id;type:bigint;not null"`
	UpdaterName string    `gorm:"column:updater_name;type:string;size:32;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
}

// TableName gorm callback function, get table name
func (r *Role) TableName() string {
	return "roles"
}

// FindRolesOptions 用來當查詢 Roles 的條件
type FindRolesOptions struct {
	Namespace string
	AccountID int64
}

// RoleServicer 用來處理 Role 相關業務操作的 service layer
type RoleUsecase interface {
	Role(ctx context.Context, roleID int64) (Role, error)
	Roles(ctx context.Context, opts FindRolesOptions) ([]Role, int32, error)
	Count(ctx context.Context, opts FindRolesOptions) (int32, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
}

// RoleRepository 用來處理 Role 物件的存儲的行為 repository layer
type RoleRepository interface {
	CreateRole(ctx context.Context, role *Role) error
	Role(ctx context.Context, roleID int64) (*Role, error)
	Roles(ctx context.Context, opts FindRolesOptions) ([]Role, error)
	CountRoles(ctx context.Context, namespace string) (int32, error)
	UpdateRole(ctx context.Context, role *Role) error
	RolesByAccountID(ctx context.Context, accountID int64) ([]Role, error)
}
