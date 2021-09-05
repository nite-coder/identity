package database

import (
	"context"
	"errors"
	"fmt"
	"identity/pkg/domain"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nite-coder/blackbear/pkg/log"
	"gorm.io/gorm"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (repo *AccountRepo) Account(ctx context.Context, opts domain.FindAccountOptions) (domain.Account, error) {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	var account domain.Account

	if err := db.Where("id = ?", opts.ID).Find(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return account, fmt.Errorf("mysql: account not found. %w", domain.ErrNotFound)
		}
		logger.Err(err).Interface("opts", opts).Error("mysql: get account failed.")
		return account, err
	}

	return account, nil
}

func (repo *AccountRepo) AccountByUUID(ctx context.Context, opts domain.FindAccountOptions) (domain.Account, error) {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	var account domain.Account

	if err := db.Where("uuid = ?", opts.UUID).Find(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return account, fmt.Errorf("mysql: account not found. %w", domain.ErrNotFound)
		}
		logger.Err(err).Interface("opts", opts).Error("mysql: get account by uuid failed.")
		return account, err
	}

	return account, nil
}

func (repo *AccountRepo) Accounts(ctx context.Context, opts domain.FindAccountOptions) ([]domain.Account, error) {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	var accounts []domain.Account

	db = repo.buildWhereClause(db, opts)

	if err := db.Find(&accounts).Error; err != nil {
		logger.Err(err).Interface("opts", opts).Error("mysql: get accounts failed")
		return nil, err
	}

	return accounts, nil
}

func (repo *AccountRepo) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	if err := db.Create(account).Error; err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if mysqlErr.Number == 1062 {
				return nil, fmt.Errorf("mysql: the account has already exists.  %w", domain.ErrAlreadyExists)
			}
		}
		logger.Err(err).Interface("account", account).Error("mysql: create account fail")
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) UpdateAccount(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	args := make(map[string]interface{})
	//use map to implement update
	args["id"] = account.ID

	args["uuid"] = account.UUID

	args["namespace"] = account.Namespace

	args["type"] = account.Type

	args["username"] = account.Username

	args["otp_enable"] = account.OTPEnable

	args["otp_secret"] = account.OTPSecret

	args["otp_effective_at"] = account.OTPEffectiveAt

	args["first_name"] = account.FirstName

	args["last_name"] = account.LastName

	args["avatar"] = account.Avatar

	args["email"] = account.Email

	args["mobile"] = account.Mobile

	args["external_id"] = account.ExternalID

	args["state"] = account.State

	args["failed_password_attempt"] = account.FailedPasswordAttempt

	args["client_ip"] = account.ClientIP

	args["last_login_at"] = account.LastLoginAt

	args["updater_id"] = account.UpdaterID

	args["updater_name"] = account.UpdaterName

	args["updated_at"] = time.Now().UTC()

	if err := db.Model(account).Where("id = ?", account.ID).UpdateColumns(args).Error; err != nil {
		logger.Err(err).Error("mysql: update account failed")
		return err
	}
	return nil
}

func (repo *AccountRepo) UpdateAccountPassword(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	args := make(map[string]interface{})
	args["password_encrypt"] = account.PasswordEncrypt
	args["updater_id"] = account.UpdaterID
	args["updater_name"] = account.UpdaterName
	args["updated_at"] = time.Now().UTC()

	if err := db.Model(account).Where("id = ?", account.ID).UpdateColumns(args).Error; err != nil {
		logger.Err(err).Interface("account", account).Error("mysql:uUpdate account password failed")
		return err
	}
	return nil
}

func (repo *AccountRepo) UpdateAccountLockedOut(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	args := make(map[string]interface{})
	args["state"] = account.State
	args["updater_id"] = account.UpdaterID
	args["updater_name"] = account.UpdaterName
	args["updated_at"] = account.UpdatedAt

	if err := db.Model(account).Where("id = ?", account.ID).UpdateColumns(args).Error; err != nil {
		logger.Errorf("database: UpdateAccountLockedOut : %#v fail: %v", account, err)
		return err
	}
	return nil
}

func (repo *AccountRepo) DeleteAccount(ctx context.Context, account *domain.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) ClearOTP(ctx context.Context, accountUUID string) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountOTPExpireTime(ctx context.Context, account *domain.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountOTPSecret(ctx context.Context, account *domain.Account) (string, error) {
	panic("not implemented")
}

func (repo *AccountRepo) CountAccounts(ctx context.Context, options domain.FindAccountOptions) (int64, error) {
	logger := log.FromContext(ctx)
	db := repo.db.WithContext(ctx)

	var total int64

	db = db.Model(domain.Account{})

	if err := repo.buildWhereClause(db, options).Count(&total).Error; err != nil {
		logger.Err(err).Interface("opts", options).Error("mysql: count account failed.")
		return 0, err
	}

	return total, nil
}

func (repo *AccountRepo) buildWhereClause(db *gorm.DB, options domain.FindAccountOptions) *gorm.DB {
	if options.LoginTimeEnd.Unix() > 0 {
		db = db.Where("last_login_at BETWEEN ? AND ?", options.LoginTimeStart, options.LoginTimeEnd)
	}

	if options.CreatedTimeStart.Unix() > 0 {
		db = db.Where("created_at BETWEEN ? AND ?", options.CreatedTimeStart, options.CreatedTimeEnd)
	}

	if options.ID != 0 {
		db = db.Where(" id = ?", options.ID)
	}

	if options.UUID != "" {
		db = db.Where(" uuid = ?", options.UUID)
	}

	if options.ExternalID != "" {
		db = db.Where(" external_id = ?", options.ExternalID)
	}

	if options.Namespace != "" {
		db = db.Where(" namespace = ?", options.Namespace)
	}

	if options.Username != "" {
		db = db.Where(" username = ?", options.Username)
	}

	if options.Email != "" {
		db = db.Where(" email = ?", options.Email)
	}

	if options.Mobile != "" {
		db = db.Where(" mobile = ?", options.Mobile)
	}

	if options.FirstName != "" {
		db = db.Where(" first_name = ?", options.FirstName)
	}

	if len(options.Role) != 0 {
		db = db.Where(" roles_id in (?)", options.Role)
	}

	if options.State > 0 {
		db = db.Where(" state = ?", options.State)
	}

	if options.Keyword != "" {
		db = db.Where(" CONCAT(username,email,first_name) like ? ", "%"+options.Keyword+"%")
	}

	if options.Type > 0 {
		db = db.Where("type = ?", options.Type)
	}

	if options.Limit > 0 {
		db = db.Limit(options.Limit)
	}

	if options.Offset > 0 {
		db = db.Offset(options.Offset)
	}

	if options.Sort != "" {
		db = db.Order(options.Sort)
	}

	return db
}
