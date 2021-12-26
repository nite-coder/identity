package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"identity/internal/pkg/database"
	"identity/pkg/domain"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nite-coder/blackbear/pkg/log"
	"gorm.io/gorm"
)

type AccountRepo struct {
}

func NewAccountRepo() *AccountRepo {
	return &AccountRepo{}
}

func (repo *AccountRepo) Account(ctx context.Context, namespace string, accountID uint64) (*domain.Account, error) {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	var account domain.Account

	if err := db.Where("id = ?", accountID).Where("namespace = ?", namespace).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &account, fmt.Errorf("mysql: account not found. %w", domain.ErrNotFound)
		}
		logger.Err(err).Interface("params", accountID).Error("mysql: get account failed.")
		return &account, err
	}

	return &account, nil
}

func (repo *AccountRepo) AccountByUUID(ctx context.Context, namespace string, uuid string) (*domain.Account, error) {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	var account domain.Account

	if err := db.Where("uuid = ?", uuid).Where("namespace = ?", namespace).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &account, fmt.Errorf("mysql: account not found. %w", domain.ErrNotFound)
		}
		logger.Err(err).Interface("params", uuid).Error("mysql: get account by uuid failed.")
		return &account, err
	}

	return &account, nil
}

func (repo *AccountRepo) Accounts(ctx context.Context, opts domain.FindAccountOptions) ([]domain.Account, error) {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	var accounts []domain.Account

	db = repo.buildWhereClause(db, opts)

	if err := db.Find(&accounts).Error; err != nil {
		logger.Err(err).Interface("params", opts).Error("mysql: get accounts failed")
		return nil, err
	}

	return accounts, nil
}

func (repo *AccountRepo) CreateAccount(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	account.CreatedAt = time.Now().UTC()
	account.LastLoginAt = time.Unix(0, 0)

	if err := db.Create(account).Error; err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if mysqlErr.Number == 1062 {
				return fmt.Errorf("mysql: the account has already exists.  %w", domain.ErrAlreadyExists)
			}
		}
		logger.Err(err).Interface("params", account).Error("mysql: create account fail")
		return err
	}

	return nil
}

func (repo *AccountRepo) UpdateAccount(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	args := make(map[string]interface{})
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
	args["nick_name"] = account.NickName
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
	args["version"] = gorm.Expr("version + 1")

	result := db.Model(account).
		Where("id = ?", account.ID).
		Where("version = @version", sql.Named("version", account.Version)).
		Updates(args)

	err := result.Error
	if err != nil {
		logger.Err(err).Error("mysql: update account failed")
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrStale
	}
	return nil
}

func (repo *AccountRepo) UpdateAccountPassword(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	args := make(map[string]interface{})
	args["password_encrypt"] = account.PasswordEncrypt
	args["updater_id"] = account.UpdaterID
	args["updater_name"] = account.UpdaterName
	args["updated_at"] = time.Now().UTC()
	args["version"] = gorm.Expr("version + 1")

	result := db.Model(account).
		Where("id = ?", account.ID).
		Where("version = @version", sql.Named("version", account.Version)).
		Updates(args)

	err := result.Error
	if err != nil {
		logger.Err(err).Interface("params", account).Error("mysql:update account password failed")
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrStale
	}

	return nil
}

func (repo *AccountRepo) UpdateState(ctx context.Context, account *domain.Account) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	args := make(map[string]interface{})
	args["state"] = account.State
	args["state_changed_at"] = time.Now().UTC()
	args["updater_id"] = account.UpdaterID
	args["updater_name"] = account.UpdaterName
	args["updated_at"] = time.Now().UTC()
	args["version"] = gorm.Expr("version + 1")

	result := db.Model(account).
		Where("id = ?", account.ID).
		Where("version = @version", sql.Named("version", account.Version)).
		Updates(args)

	err := result.Error
	if err != nil {
		logger.Err(err).Interface("param", account).Error("mysql: update state failed")
	}

	if result.RowsAffected == 0 {
		return domain.ErrStale
	}

	return nil
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
	db := database.FromContext(ctx)

	var total int64

	db = db.Model(domain.Account{})

	if err := repo.buildWhereClause(db, options).Count(&total).Error; err != nil {
		logger.Err(err).Interface("params", options).Error("mysql: count account failed.")
		return 0, err
	}

	return total, nil
}

func (repo *AccountRepo) AccountsByRoleID(ctx context.Context, namespace string, roleID uint64) ([]domain.Account, error) {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	stmt := db.Statement

	err := stmt.Parse(&domain.Account{})
	if err != nil {
		return nil, err
	}
	accountTable := stmt.Schema.Table

	err = stmt.Parse(&domain.AccountRole{})
	if err != nil {
		return nil, err
	}
	accountRolesTable := stmt.Schema.Table

	joinStr := fmt.Sprintf("left join `%s` on `%s`.account_id = `%s`.id", accountRolesTable, accountRolesTable, accountTable)

	var accounts []domain.Account
	err = db.Model(domain.Account{}).
		Joins(joinStr).
		Where("role_id = ?", roleID).
		Where("namespace = ?", namespace).
		Find(&accounts).Error

	if err != nil {
		logger.Err(err).Error("mysql: get accounts by role id failed.")
		return accounts, err
	}

	return accounts, nil
}

func (repo *AccountRepo) AddRolesToAccount(ctx context.Context, request domain.AddRolesToAccountRequest) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	err := db.Where("account_id = ?", request.AccountID).Delete(&domain.AccountRole{}).Error
	if err != nil {
		logger.Err(err).Error("mysql: delete account role failed")
		return err
	}

	accountRoles := []domain.AccountRole{}

	for _, roleID := range request.RoleIDs {
		accountRole := domain.AccountRole{
			AccountID: request.AccountID,
			RoleID:    roleID,
		}
		accountRoles = append(accountRoles, accountRole)
	}

	err = db.CreateInBatches(accountRoles, 100).Error
	if err != nil {
		logger.Err(err).Error("mysql: add accounts to role failed")
		return err
	}

	return nil
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

	if options.MobileCountryCode != "" {
		db = db.Where(" mobile_country_code = ?", options.MobileCountryCode)
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
