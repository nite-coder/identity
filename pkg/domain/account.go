package domain

import (
	"context"
	"time"
)

const (
	AccountStatusNone     AccountState = 0
	AccountStatusNormal   AccountState = 1 //狀態正常
	AccountStatusDisabled AccountState = 2 //人工鎖定
	AccountStatusLocked   AccountState = 3 //密碼錯誤次數過多

	SystemName = "system"
	SystemID   = 0
)

var (
	TableNameAccount     = "accounts"
	TableNameAccountRole = "accounts_roles"
)

type AccountState int32

func (state AccountState) String() string {
	switch state {
	case AccountStatusNormal:
		return "normal"
	case AccountStatusDisabled:
		return "disabled"
	case AccountStatusLocked:
		return "locked"
	default:
		return "unknown"
	}
}

// Account represent account information
type Account struct {
	ID                    uint64       `gorm:"primaryKey;autoIncrement;not null"`
	UUID                  string       `gorm:"column:uuid;type:char(36);size:36;uniqueIndex:uniq_uuid;not null"`
	Namespace             string       `gorm:"column:namespace;type:string;size:256;uniqueIndex:uniq_username;not null"`
	Type                  int32        `gorm:"column:type;type:int;not null"`
	Username              string       `gorm:"column:username;type:string;size:32;uniqueIndex:uniq_username;not null"`
	PasswordEncrypt       string       `gorm:"column:password_encrypt;type:string;size:128;not null"`
	FirstName             string       `gorm:"column:first_name;type:string;size:24;not null"`
	LastName              string       `gorm:"column:last_name;type:string;size:24;not null"`
	Avatar                string       `gorm:"column:avatar;type:string;size:24;not null"`
	Email                 string       `gorm:"column:email;type:string;size:128;not null"`
	MobileCountryCode     string       `gorm:"column:mobile_country_code;type:string;size:5;not null"`
	Mobile                string       `gorm:"column:mobile;type:string;size:20;not null"`
	ExternalID            string       `gorm:"column:external_id;type:string;size:128;not null"`
	FailedPasswordAttempt int32        `gorm:"column:failed_password_attempt;type:int;not null"`
	OTPEnable             int32        `gorm:"column:otp_enable;type:tinyint;not null"`
	OTPSecret             string       `gorm:"column:otp_secret;type:string;size:64;not null"`
	OTPEffectiveAt        time.Time    `gorm:"column:otp_effective_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
	OTPLastResetAt        time.Time    `gorm:"column:otp_last_reset_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
	ClientIP              string       `gorm:"column:client_ip;type:string;size:64;not null"`
	Note                  string       `gorm:"column:notes;type:string;size:512;not null"`
	LastLoginAt           time.Time    `gorm:"column:last_login_at;type:datetime;not null"`
	IsAdmin               int32        `gorm:"column:is_admin;type:tinyint;not null"`
	State                 AccountState `gorm:"column:state;type:int;not null"`
	StateChangedAt        time.Time    `gorm:"column:state_changed_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
	Version               uint32       `gorm:"column:version;type:int;not null"`
	CreatorID             uint64       `gorm:"column:creator_id;type:bigint;not null"`
	CreatorName           string       `gorm:"column:creator_name;type:string;size:32;not null"`
	CreatedAt             time.Time    `gorm:"column:created_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
	UpdaterID             uint64       `gorm:"column:updater_id;type:bigint;not null"`
	UpdaterName           string       `gorm:"column:updater_name;type:string;size:32;not null"`
	UpdatedAt             time.Time    `gorm:"column:updated_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
}

// TableName 用來取 Account 的資料表名稱
func (a *Account) TableName() string {
	return TableNameAccount
}

type AccountRole struct {
	AccountID uint64 `gorm:"primaryKey;not null"`
	RoleID    uint64 `gorm:"primaryKey;not null"`
}

func (a *AccountRole) TableName() string {
	return TableNameAccountRole
}

// FindAccountOptions 用來查詢 Account 的選項
type FindAccountOptions struct {
	ID               uint64
	UUID             string
	ExternalID       string
	Namespace        string
	Username         string
	Email            string
	Mobile           string
	Role             []string
	State            AccountState
	FirstName        string
	Offset           int
	Limit            int
	SortBy           string
	Sort             string
	CreatedTimeStart time.Time
	CreatedTimeEnd   time.Time
	LoginTimeStart   time.Time
	LoginTimeEnd     time.Time
	Keyword          string
	Type             int32
}

// LoginInfo 用來傳遞登入資訊
type LoginInfo struct {
	Namespace string `json:"namespace,omitempty"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type UpdateAccountPasswordRequest struct {
	Namespace   string `json:"namespace,omitempty"`
	AccountID   uint64
	OldPassword string
	NewPassword string
	UpdaterID   uint64
	UpdaterName string
}

type ForceUpdateAccountPasswordRequest struct {
	Namespace   string `json:"namespace,omitempty"`
	AccountID   uint64
	NewPassword string
	UpdaterID   uint64
	UpdaterName string
}

type ChangeStateRequest struct {
	Namespace   string `json:"namespace,omitempty"`
	AccountID   uint64
	State       AccountState
	UpdaterID   uint64
	UpdaterName string
}

// AccountUsecase 用來處理 Account 相關業務操作的 service layer
type AccountUsecase interface {
	Account(ctx context.Context, namespace string, accountID uint64) (*Account, error)
	AccountByUUID(ctx context.Context, namespace string, uuid string) (*Account, error)
	Accounts(ctx context.Context, opts FindAccountOptions) ([]Account, error)
	CountAccounts(ctx context.Context, opts FindAccountOptions) (int64, error)
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	UpdateAccountPassword(ctx context.Context, request UpdateAccountPasswordRequest) error
	ForceUpdateAccountPassword(ctx context.Context, request ForceUpdateAccountPasswordRequest) error
	ChangeState(ctx context.Context, request ChangeStateRequest) error
	Login(ctx context.Context, loginInfo LoginInfo) (*Account, error)
	ClearOTP(ctx context.Context, accountUUID string) error
	GenerateOTPAuth(ctx context.Context, accountID uint64) (string, error)
	SetOTPExpireTime(ctx context.Context, accountUUID string, duration int64) error
	VerifyOTP(ctx context.Context, accountUUID string, otpCode string) (*Account, error)
	//AccountIDsByRoleName(ctx context.Context, namespace, roleName string) ([]int64, error)
	UpdateAccountRole(ctx context.Context, accountID uint64, roles []int64, updaterAccountID uint64, updaterUsername string) error
}

// AccountRepository 用來處理 Account 物件的存儲的行為 repository layer
type AccountRepository interface {
	Account(ctx context.Context, namespace string, accountID uint64) (*Account, error)
	AccountByUUID(ctx context.Context, namespace string, uuid string) (*Account, error)
	Accounts(ctx context.Context, opts FindAccountOptions) ([]Account, error)
	AccountsByRoleID(ctx context.Context, namespace string, roleID uint64) ([]Account, error)
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *Account) error
	UpdateAccountPassword(ctx context.Context, account *Account) error
	UpdateState(ctx context.Context, account *Account) error
	ClearOTP(ctx context.Context, accountUUID string) error
	UpdateAccountOTPExpireTime(ctx context.Context, account *Account) error
	UpdateAccountOTPSecret(ctx context.Context, account *Account) (string, error)
	CountAccounts(ctx context.Context, options FindAccountOptions) (int64, error)
	//AccountIDsByRoleName(ctx context.Context, roleID int64) ([]int64, error)
	//UpdateAccountRole(ctx context.Context, accountID int64, roles []int64, updaterAccountID int64, updaterUsername string) error
}
