package identity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// Role 代表角色資訊
type Role struct {
	ID              int64     `json:"id,omitempty" gorm:"column:id"`
	UUID            string    `json:"uuid,omitempty" gorm:"column:uuid"`
	Namespace       string    `json:"namespace,omitempty" gorm:"column:namespace"`
	Name            string    `json:"name,omitempty" gorm:"column:name"`
	Rules           []Rule    `json:"-" gorm:"column:-"`
	RulesJSONString string    `json:"rule_json" gorm:"column:rules"`
	CreatorID       int64     `json:"creatorID" gorm:"column:creator_id"`
	CreatorName     string    `json:"creatorName" gorm:"column:creator_name"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdaterID       int64     `json:"updaterID" gorm:"column:updater_id"`
	UpdaterName     string    `json:"updaterName" gorm:"column:updater_name"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// Rule 代表資源存取資訊
type Rule struct {
	Namespace     string   `json:"namespace"`
	Resources     []string `json:"resources"`
	ResourceNames []string `json:"resource_names"`
	Verbs         []string `json:"verbs"`
}

// BridgeAccountsRoles 代表account對應role的資訊
type BridgeAccountsRoles struct {
	AccountID   int64     `json:"account_id,omitempty" gorm:"column:accounts_id"`
	RoleID      int64     `json:"role_id,omitempty" gorm:"column:roles_id"`
	UpdaterID   int64     `json:"updater_id" gorm:"column:updater_id" `
	UpdaterName string    `json:"updater_name" gorm:"column:updater_name" `
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at" `
}

// FindRolesOptions 用來當查詢 Roles 的條件
type FindRolesOptions struct {
	Namespace string
	AccountID int64
}

// RoleServicer 用來處理 Role 相關業務操作的 service layer
type RoleServicer interface {
	Role(ctx context.Context, roleID int64) (Role, error)
	Roles(ctx context.Context, opts FindRolesOptions) ([]Role, int32, error)
	Count(ctx context.Context, opts FindRolesOptions) (int32, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
}

// RoleRepository 用來處理 Role 物件的存儲的行為 repository layer
type RoleRepository interface {
	WriteDB() *gorm.DB
	CreateRole(ctx context.Context, role *Role) error
	Role(ctx context.Context, roleID int64) (*Role, error)
	Roles(ctx context.Context, opts FindRolesOptions) ([]Role, error)
	CountRoles(ctx context.Context, namespace string) (int32, error)
	UpdateRole(ctx context.Context, role *Role) error
	RolesByAccountID(ctx context.Context, accountID int64) ([]Role, error)
}

// TableName gorm callback function, get table name
func (r *Role) TableName() string {
	return "roles"
}
