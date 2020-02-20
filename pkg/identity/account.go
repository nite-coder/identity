package identity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// Account represent account information
type Account struct {
	ID                    int64     `json:"id,omitempty" gorm:"column:id;primary_key"`
	UUID                  string    `json:"uuid,omitempty" gorm:"column:uuid"`
	Namespace             string    `json:"namespace,omitempty" gorm:"column:namespace"`
	Type                  int32     `json:"type,omitempty" gorm:"column:type"`
	Username              string    `json:"username,omitempty" gorm:"column:username"`
	PasswordEncrypt       string    `json:"-" gorm:"column:password_encrypt"` //b crytp
	OTPEnable             bool      `json:"otp_enable" gorm:"column:otp_enable"`
	OTPSecret             string    `json:"-" gorm:"column:otp_secret"`
	OTPEffectiveAt        time.Time `json:"otp_effective_time,omitempty" gorm:"column:otp_effective_at"`
	FirstName             string    `json:"firstName" gorm:"column:first_name"`
	LastName              string    `json:"lastName" gorm:"column:last_name"`
	Avatar                string    `json:"avatar,omitempty" gorm:"column:avatar"`
	Email                 string    `json:"email,omitempty" gorm:"column:email"`
	Mobile                string    `json:"mobile,omitempty" gorm:"column:mobile"`
	ExternalID            string    `json:"externalID,omitempty" gorm:"column:external_id"`
	IsLockedOut           int32     `json:"isLockedOut,omitempty" gorm:"column:is_locked_out"`
	FailedPasswordAttempt int32     `json:"failedPassword_attempt,omitempty" gorm:"column:failed_password_attempt"`
	ClientIP              string    `json:"clientIP,omitempty" gorm:"column:client_ip"`
	UserAgent             string    `json:"userAgent,omitempty" gorm:"column:user_agent"`
	Notes                 string    `json:"notes,omitempty" gorm:"column:notes"`
	IsAdmin               bool      `json:"isAdmin,omitempty" gorm:"column:is_admin"`
	LastLoginAt           time.Time `json:"lastLoginTime,omitempty" gorm:"column:last_login_at"`
	Roles                 []Role    `json:"roles,omitempty"` //bridge
	CreatorID             int64     `json:"creatorID" gorm:"column:creator_id"`
	CreatorName           string    `json:"creatorName" gorm:"column:creator_name"`
	CreatedAt             time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdaterID             int64     `json:"updaterID" gorm:"column:updater_id"`
	UpdaterName           string    `json:"updaterName" gorm:"column:updater_name"`
	UpdatedAt             time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// FindAccountOptions 用來查詢 Account 的選項
type FindAccountOptions struct {
	ID               int64
	UUID             string
	ExternalID       string
	Namespace        string
	Username         string
	Email            string
	Mobile           string
	Role             []string
	IsLockedOut      int32
	FirstName        string
	Offset           int32
	Limit            int32
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

// AccountServicer 用來處理 Account 相關業務操作的 service layer
type AccountServicer interface {
	Account(ctx context.Context, accountID int64) (Account, error)
	AccountByUUID(ctx context.Context, accountUUID string) (Account, error)
	Accounts(ctx context.Context, opts FindAccountOptions) ([]Account, error)
	CountAccounts(ctx context.Context, opts FindAccountOptions) (int32, error)
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *Account) error
	UpdateAccountPassword(ctx context.Context, accountID int64, oldPassword string, newPassword string, updaterAccountID int64, updaterUsername string) error
	DeleteAccount(ctx context.Context, accountID int64, updaterAccountID int64, updaterUsername string) error
	ForceUpdateAccountPassword(ctx context.Context, accountID int64, newPassword string, updaterAccountID int64, updaterUsername string) error
	LockAccount(ctx context.Context, accountID int64, lockedType int32, updaterAccountID int64, updaterUsername string) error
	LockAccounts(ctx context.Context, accountIDs []int64, lockedTypes []int32, updaterAccountID int64, updaterUsername string) error
	UnlockAccount(ctx context.Context, accountID int64, updaterAccountID int64, updaterUsername string) error
	Login(ctx context.Context, loginInfo *LoginInfo) (*Account, error)
	ClearOTP(ctx context.Context, accountUUID string) error
	GenerateOTPAuth(ctx context.Context, accountID int64) (string, error)
	SetOTPExpireTime(ctx context.Context, accountUUID string, duration int64) error
	VerifyOTP(ctx context.Context, accountUUID string, otpCode string) (*Account, error)
	//AccountIDsByRoleName(ctx context.Context, namespace, roleName string) ([]int64, error)
	UpdateAccountRole(ctx context.Context, accountID int64, roles []int64, updaterAccountID int64, updaterUsername string) error
}

// AccountRepository 用來處理 Account 物件的存儲的行為 repository layer
type AccountRepository interface {
	WriteDB() *gorm.DB
	Account(ctx context.Context, opts FindAccountOptions) (Account, error)
	AccountByUUID(ctx context.Context, opts FindAccountOptions) (Account, error)
	Accounts(ctx context.Context, opts FindAccountOptions) ([]Account, error)
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	UpdateAccountPassword(ctx context.Context, account *Account) error
	UpdateAccountLockedOut(ctx context.Context, account *Account) error
	UpdateAccountLockedOutTX(ctx context.Context, DB *gorm.DB, account *Account) error
	DeleteAccount(ctx context.Context, account *Account) error
	ClearOTP(ctx context.Context, accountUUID string) error
	UpdateAccountOTPExpireTime(ctx context.Context, account *Account) error
	UpdateAccountOTPSecret(ctx context.Context, account *Account) (string, error)
	CountAccounts(ctx context.Context, options *FindAccountOptions) (int32, error)
	//AccountIDsByRoleName(ctx context.Context, roleID int64) ([]int64, error)
	//UpdateAccountRole(ctx context.Context, accountID int64, roles []int64, updaterAccountID int64, updaterUsername string) error
}

// TableName 用來取 Account 的資料表名稱
func (a *Account) TableName() string {
	return "accounts"
}
