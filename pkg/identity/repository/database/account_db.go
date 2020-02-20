package database

import (
	"context"
	identity "identity/pkg/identity"

	"github.com/jinzhu/gorm"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (repo *AccountRepo) WriteDB() *gorm.DB {
	panic("not implemented")
}

func (repo *AccountRepo) Account(ctx context.Context, opts identity.FindAccountOptions) (identity.Account, error) {
	panic("not implemented")
}

func (repo *AccountRepo) AccountByUUID(ctx context.Context, opts identity.FindAccountOptions) (identity.Account, error) {
	panic("not implemented")
}

func (repo *AccountRepo) Accounts(ctx context.Context, opts identity.FindAccountOptions) ([]identity.Account, error) {
	panic("not implemented")
}

func (repo *AccountRepo) CreateAccount(ctx context.Context, account *identity.Account) (*identity.Account, error) {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccount(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountPassword(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountLockedOut(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountLockedOutTX(ctx context.Context, DB *gorm.DB, account *identity.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) DeleteAccount(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) ClearOTP(ctx context.Context, accountUUID string) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountOTPExpireTime(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (repo *AccountRepo) UpdateAccountOTPSecret(ctx context.Context, account *identity.Account) (string, error) {
	panic("not implemented")
}

func (repo *AccountRepo) CountAccounts(ctx context.Context, options *identity.FindAccountOptions) (int32, error) {
	panic("not implemented")
}
