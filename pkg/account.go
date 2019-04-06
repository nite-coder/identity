package identity

import (
	"context"
	"time"
)

type Account struct {
	ID                    int64      `json:"id,omitempty" db:"id"`
	UUID                  string     `json:"uuid,omitempty" db:"uuid"`
	Namespace             string     `json:"namespace,omitempty" db:"namespace"`
	Username              string     `json:"username,omitempty" db:"username"`
	PasswordHash          string     `json:"-" db:"password_hash"`
	FirstName             string     `json:"first_name" db:"first_name"`
	LastName              string     `json:"last_name" db:"last_name"`
	Avatar                string     `json:"avatar,omitempty" db:"avatar"`
	Email                 string     `json:"email,omitempty" db:"email"`
	Mobile                string     `json:"mobile,omitempty" db:"mobile"`
	ExternalID            string     `json:"external_id,omitempty" db:"external_id"`
	IsLockedOut           bool       `json:"is_locked_out,omitempty" db:"is_locked_out"`
	FailedPasswordAttempt int        `json:"failed_password_attempt_count,omitempty" db:"failed_password_attempt_count"`
	Roles                 []*Role    `json:"roles,omitempty"`
	ClientIP              string     `json:"client_ip,omitempty" db:"client_ip"`
	UserAgent             string     `json:"user_agent,omitempty" db:"user_agent"`
	Note                  string     `json:"note,omitempty" db:"note"`
	LastLoginAt           *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatorID             int64
	CreatorName           string
	CreatedAt             *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdaterID             int64
	UpdatedAt             *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	UpdaterName           string
}

type FindAccountOptions struct {
	ID               string `json:"id" db:"id"`
	ExternalID       string `json:"external_id" db:"external_id"`
	App              string `json:"app" db:"app"`
	Username         string `json:"username" db:"username"`
	Email            string `json:"email,omitempty" db:"email"`
	Mobile           string `json:"mobile,omitempty" db:"mobile"`
	Role             string `json:"role"`
	IsLockedOut      int    `json:"is_locked_out" db:"is_locked_out"`
	Skip             int    `db:"skip"`
	Take             int    `db:"take"`
	SortBy           string `db:"sortby"`
	Sort             string
	CreatedTimeStart *time.Time `db:"created_start_time"`
	CreatedTimeEnd   *time.Time `db:"created_end_time"`
	LoginTimeStart   *time.Time `db:"login_start_time"`
	LoginTimeEnd     *time.Time `db:"login_end_time"`
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
