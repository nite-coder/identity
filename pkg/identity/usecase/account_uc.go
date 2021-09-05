package usecase

import (
	"context"
	"fmt"
	"identity/pkg/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountUsecase struct {
	accountRepo domain.AccountRepository
}

func NewAccountUsecase(accountRepo domain.AccountRepository) *AccountUsecase {
	return &AccountUsecase{
		accountRepo: accountRepo,
	}
}
func (uc *AccountUsecase) Account(ctx context.Context, accountID uint64) (domain.Account, error) {
	options := domain.FindAccountOptions{
		ID: accountID,
	}

	return uc.accountRepo.Account(ctx, options)
}

func (uc *AccountUsecase) AccountByUUID(ctx context.Context, accountUUID string) (domain.Account, error) {
	options := domain.FindAccountOptions{
		UUID: accountUUID,
	}

	return uc.accountRepo.AccountByUUID(ctx, options)
}

func (uc *AccountUsecase) Accounts(ctx context.Context, opts domain.FindAccountOptions) ([]domain.Account, error) {
	account, err := uc.accountRepo.Accounts(ctx, opts)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (uc *AccountUsecase) CountAccounts(ctx context.Context, opts domain.FindAccountOptions) (int64, error) {
	total, err := uc.accountRepo.CountAccounts(ctx, opts)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (uc *AccountUsecase) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	var err error

	if account.Namespace == "" {
		return nil, fmt.Errorf("namespace is empty. %w", domain.ErrInvalidInput)
	}

	if account.Username == "" {
		return nil, fmt.Errorf("username is empty. %w", domain.ErrInvalidInput)
	}

	if account.PasswordEncrypt == "" {
		return nil, fmt.Errorf("username is empty. %w", domain.ErrInvalidInput)
	}

	if account.CreatorName == "" {
		return nil, fmt.Errorf("creator name is empty. %w", domain.ErrInvalidInput)
	}

	account.UUID = uuid.NewString()
	account.PasswordEncrypt, err = encryptPassword(account.PasswordEncrypt)
	account.OTPEffectiveAt = time.Unix(0, 0)
	account.OTPLastResetAt = time.Unix(0, 0)
	account.LastLoginAt = time.Unix(0, 0)
	if err != nil {
		return nil, err
	}

	newAccount, err := uc.accountRepo.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (uc *AccountUsecase) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return uc.accountRepo.UpdateAccount(ctx, account)
}

func (uc *AccountUsecase) UpdateAccountPassword(ctx context.Context, accountID int64, oldPassword string, newPassword string, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) DeleteAccount(ctx context.Context, accountID int64, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) ForceUpdateAccountPassword(ctx context.Context, accountID uint64, newPassword string, updaterAccountID uint64, updaterUsername string) error {
	//find account
	account, err := uc.Account(ctx, accountID)
	if err != nil {
		return err
	}

	//update
	newPassword, err = encryptPassword(newPassword)
	if err != nil {
		return err
	}

	updateAccount := domain.Account{
		ID:              account.ID,
		PasswordEncrypt: newPassword,
	}
	updateAccount.UpdaterID = updaterAccountID
	updateAccount.UpdaterName = updaterUsername
	updateAccount.UpdatedAt = time.Now().UTC()

	return uc.accountRepo.UpdateAccountPassword(ctx, &updateAccount)
}

func (uc *AccountUsecase) LockAccount(ctx context.Context, accountID int64, lockedType int32, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) LockAccounts(ctx context.Context, accountIDs []int64, lockedTypes []int32, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) UnlockAccount(ctx context.Context, accountID int64, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) Login(ctx context.Context, loginInfo *domain.LoginInfo) (*domain.Account, error) {
	//find account
	options := domain.FindAccountOptions{
		Namespace: loginInfo.Namespace,
		Username:  loginInfo.Username,
	}
	accounts, err := uc.Accounts(ctx, options)

	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, domain.ErrUsernameOrPasswordIncorrect
	}

	account := accounts[0]

	//check status
	if account.State == domain.AccountStatusLocked || account.State == domain.AccountStatusDisabled {
		return nil, domain.ErrAccountDisable
	}

	//compare password
	err = isPasswordValid(account.PasswordEncrypt, loginInfo.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			//帳密錯誤更新錯誤次數
			account.FailedPasswordAttempt = account.FailedPasswordAttempt + 1
			err := uc.accountRepo.UpdateAccount(ctx, &account)
			if err != nil {
				return nil, err
			}
		}

		return nil, domain.ErrUsernameOrPasswordIncorrect
	}

	//清除登入失敗次數
	account.FailedPasswordAttempt = 0
	err = uc.accountRepo.UpdateAccount(ctx, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (uc *AccountUsecase) ClearOTP(ctx context.Context, accountUUID string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) GenerateOTPAuth(ctx context.Context, accountID int64) (string, error) {
	panic("not implemented")
}

func (uc *AccountUsecase) SetOTPExpireTime(ctx context.Context, accountUUID string, duration int64) error {
	panic("not implemented")
}

func (uc *AccountUsecase) VerifyOTP(ctx context.Context, accountUUID string, otpCode string) (*domain.Account, error) {
	panic("not implemented")
}

//AccountIDsByRoleName(ctx context.Context, namespace, roleName string) ([]int64, error)
func (uc *AccountUsecase) UpdateAccountRole(ctx context.Context, accountID int64, roles []int64, updaterAccountID int64, updaterUsername string) error {
	panic("not implemented")
}

func encryptPassword(password string) (string, error) {
	newEncrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(newEncrypt), nil
}

// isPasswordValid 比對password是否正確
func isPasswordValid(encryptPassword, oldPassword string) error {
	// compare password
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(oldPassword))
	if err != nil {
		return err
	}
	return nil
}
