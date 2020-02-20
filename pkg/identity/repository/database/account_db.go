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

func (repo *AccountRepo) DB() *gorm.DB {
	return repo.db
}

func (repo *AccountRepo) CreateAccount(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}
