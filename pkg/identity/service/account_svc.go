package service

import (
	"context"
	identity "identity/pkg/identity"
)

type AccountService struct {
	accountRepo identity.AccountRepository
}

func NewAccountService(accountRepo identity.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}
func (svc *AccountService) Account(ctx context.Context, accountID int64) (identity.Account, error) {
	panic("not implemented")
}

func (svc *AccountService) AccountByUUID(ctx context.Context, accountUUID string) (identity.Account, error) {
	panic("not implemented")
}

func (svc *AccountService) Accounts(ctx context.Context, opts identity.FindAccountOptions) ([]identity.Account, error) {
	panic("not implemented")
}

func (svc *AccountService) CountAccounts(ctx context.Context, opts identity.FindAccountOptions) (int32, error) {
	panic("not implemented")
}

func (svc *AccountService) CreateAccount(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (svc *AccountService) UpdateAccount(ctx context.Context, account *identity.Account) error {
	panic("not implemented")
}

func (svc *AccountService) UpdateAccountPassword(ctx context.Context, accountID int64, oldPassword string, newPassword string, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (svc *AccountService) DeleteAccount(ctx context.Context, accountID int64, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (svc *AccountService) ForceUpdateAccountPassword(ctx context.Context, accountID int64, newPassword string, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (svc *AccountService) LockAccount(ctx context.Context, accountID int64, lockedType int32, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (svc *AccountService) LockAccounts(ctx context.Context, accountIDs []int64, lockedTypes []int32, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (svc *AccountService) UnlockAccount(ctx context.Context, accountID int64, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (svc *AccountService) Login(ctx context.Context, loginInfo *identity.LoginInfo) (*identity.Account, error) {
	panic("not implemented")
}

func (svc *AccountService) ClearOTP(ctx context.Context, accountUUID string) error {
	panic("not implemented")
}

func (svc *AccountService) GenerateOTPAuth(ctx context.Context, accountID int64) (string, error) {
	panic("not implemented")
}

func (svc *AccountService) SetOTPExpireTime(ctx context.Context, accountUUID string, duration int64) error {
	panic("not implemented")
}

func (svc *AccountService) VerifyOTP(ctx context.Context, accountUUID string, otpCode string) (*identity.Account, error) {
	panic("not implemented")
}

//AccountIDsByRoleName(ctx context.Context, namespace, roleName string) ([]int64, error)
func (svc *AccountService) UpdateAccountRole(ctx context.Context, accountID int64, roles []int64, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}
