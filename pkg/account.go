package identity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID                    int64     `json:"id,omitempty" gorm:"column:id;primary_key"`
	UUID                  string    `json:"uuid,omitempty" gorm:"column:uuid"`
	Namespace             string    `json:"namespace,omitempty" gorm:"column:namespace"`
	Username              string    `json:"username,omitempty" gorm:"column:username"`
	PasswordHash          string    `json:"-" gorm:"column:password_hash"`
	FirstName             string    `json:"first_name" gorm:"column:first_name"`
	LastName              string    `json:"last_name" gorm:"column:last_name"`
	Avatar                string    `json:"avatar,omitempty" gorm:"column:avatar"`
	Email                 string    `json:"email,omitempty" gorm:"column:email"`
	Mobile                string    `json:"mobile,omitempty" gorm:"column:mobile"`
	ExternalID            string    `json:"external_id,omitempty" gorm:"column:external_id"`
	IsLockedOut           bool      `json:"is_locked_out,omitempty" gorm:"column:is_locked_out"`
	FailedPasswordAttempt int       `json:"failed_password_attempt_count,omitempty" gorm:"column:failed_password_attempt_count"`
	Roles                 []*Role   `json:"roles,omitempty"`
	ClientIP              string    `json:"client_ip,omitempty" gorm:"column:client_ip"`
	UserAgent             string    `json:"user_agent,omitempty" gorm:"column:user_agent"`
	Notes                 string    `json:"notes,omitempty" gorm:"column:notes"`
	LastLoginAt           time.Time `json:"last_login_at,omitempty" gorm:"column:last_login_at"`
	CreatorID             int64     `json:"creator_id,omitempty" gorm:"column:creator_id"`
	CreatorName           string    `json:"creator_name,omitempty" gorm:"column:creator_name"`
	CreatedAt             time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdaterID             int64     `json:"updater_id,omitempty" gorm:"column:updater_id"`
	UpdaterName           string    `json:"updater_name,omitempty" gorm:"column:updater_name"`
	UpdatedAt             time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type FindAccountOptions struct {
	ID               string `json:"id" gorm:"id"`
	ExternalID       string `json:"external_id" gorm:"external_id"`
	App              string `json:"app" gorm:"app"`
	Username         string `json:"username" gorm:"username"`
	Email            string `json:"email,omitempty" gorm:"email"`
	Mobile           string `json:"mobile,omitempty" gorm:"mobile"`
	Role             string `json:"role"`
	IsLockedOut      int    `json:"is_locked_out" gorm:"is_locked_out"`
	Skip             int    `gorm:"skip"`
	Take             int    `gorm:"take"`
	SortBy           string `gorm:"sortby"`
	Sort             string
	CreatedTimeStart *time.Time `gorm:"created_start_time"`
	CreatedTimeEnd   *time.Time `gorm:"created_end_time"`
	LoginTimeStart   *time.Time `gorm:"login_start_time"`
	LoginTimeEnd     *time.Time `gorm:"login_end_time"`
}

type AccountServicer interface {
	Account(ctx context.Context, accountID string) (*Account, error)
	Accounts(ctx context.Context, opts *FindAccountOptions) ([]*Account, error)
	CountAccounts(ctx context.Context, opts *FindAccountOptions) (int, error)
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccountPassword(ctx context.Context, accountID string, newPassword string) error
	LockAccount(ctx context.Context, namespace, accountID string) error
	UnlockAccount(ctx context.Context, namespace, accountID string) error
}

type AccountRepository interface {
	DB() *gorm.DB
	CreateAccount(ctx context.Context, account *Account) error
}
