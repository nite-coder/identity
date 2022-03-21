package domain

import (
	"context"
	"time"
)

type RoleState uint32

const (
	RoleStatusDefault  RoleState = 0
	RoleStatusNormal   RoleState = 1
	RoleStatusDisabled RoleState = 2
)

// Role 代表角色資訊
type Role struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Namespace   string    `gorm:"column:namespace;type:string;size:256;uniqueIndex:uniq_name;not null"`
	Name        string    `gorm:"column:name;type:string;size:32;uniqueIndex:uniq_name;not null"`
	Desc        string    `gorm:"column:desc;type:string;size:512;not null"`
	State       RoleState `gorm:"column:state;type:int;not null"`
	Version     uint32    `gorm:"column:version;type:int;not null"`
	CreatorID   uint64    `gorm:"column:creator_id;type:bigint;not null"`
	CreatorName string    `gorm:"column:creator_name;type:string;size:32;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:1970-01-01 00:00:00;not null"`
	UpdaterID   uint64    `gorm:"column:updater_id;type:bigint;not null"`
	UpdaterName string    `gorm:"column:updater_name;type:string;size:32;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:1970-01-01 00:00:00;not null"`
}

// FindRolesOptions 用來當查詢 Roles 的條件
type FindRoleOptions struct {
	Namespace string
	Name      string
	Limit     int
	Offset    int
	Sort      string
}

// RoleUsecase 用來處理 Role 相關業務操作的場景
type RoleUsecase interface {
	Role(ctx context.Context, namespace string, id uint64) (*Role, error)
	Roles(ctx context.Context, opts FindRoleOptions) ([]Role, error)
	CreateRole(ctx context.Context, role *Role) error
	RolesByAccountID(ctx context.Context, namespace string, accountID uint64) ([]Role, error)
	// Count(ctx context.Context, opts FindRoleOptions) (uint64, error)
	UpdateRole(ctx context.Context, role *Role) error
	AddAccountsToRole(ctx context.Context, accountIDs []uint64, roleID uint64) error
}

// RoleRepository 用來處理 Role 物件的存儲的行為 repository layer
type RoleRepository interface {
	Role(ctx context.Context, namespace string, id uint64) (*Role, error)
	Roles(ctx context.Context, opts FindRoleOptions) ([]Role, error)
	RolesByAccountID(ctx context.Context, namespace string, accountID uint64) ([]Role, error)
	CreateRole(ctx context.Context, role *Role) error
	// CountRoles(ctx context.Context, namespace string) (int32, error)
	UpdateRole(ctx context.Context, role *Role) error
	AddAccountsToRole(ctx context.Context, accountIDs []uint64, roleID uint64) error
}
