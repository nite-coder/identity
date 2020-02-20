package identity

import (
	"context"
	"time"
)

type Role struct {
	ID        int64      `json:"id,omitempty" db:"id"`
	UUID      string     `json:"uuid,omitempty" db:"uuid"`
	Namespace string     `json:"namespace,omitempty" db:"namespace"`
	Name      string     `json:"name,omitempty" db:"name"`
	Rules     []Rule     `json:"rules" db:"-"`
	RulesJSON string     `json:"-" db:"rulesJSON"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Rule struct {
	Namespace     string   `json:"namespace"`
	Resources     []string `json:"resources"`
	ResourceNames []string `json:"resource_names"`
	Verbs         []string `json:"verbs"`
}

type RoleServicer interface {
	Roles(ctx context.Context) ([]*Role, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateUserRole(ctx context.Context, app, accountID string, roles []string) error
}
