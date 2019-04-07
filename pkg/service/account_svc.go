package service

import (
	"context"

	identity "identity/pkg"
)

type AccountService struct {
	accountRepo identity.AccountRepository
}

func NewAccountService(accountRepo identity.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}
func (svc *AccountService) Account(ctx context.Context, accountID string) (*identity.Account, error) {
	panic("not implemented")
}

func (svc *AccountService) Accounts(ctx context.Context, opt *identity.FindAccountOptions) ([]*identity.Account, error) {
	panic("not implemented")
}

func (svc *AccountService) CountAccounts(ctx context.Context, opt *identity.FindAccountOptions) (int, error) {
	panic("not implemented")
}

func (svc *AccountService) CreateAccount(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (svc *AccountService) UpdateAccountPassword(ctx context.Context, accountID string, newPassword string) error {
	panic("not implemented")
}

func (svc *AccountService) LockAccount(ctx context.Context, app string, accountID string) error {
	panic("not implemented")
}

func (svc *AccountService) UnlockAccount(ctx context.Context, app string, accountID string) error {
	panic("not implemented")
}
