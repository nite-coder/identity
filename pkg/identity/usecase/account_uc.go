package usecase

import (
	"context"
	"errors"
	"fmt"
	"identity/internal/pkg/database"
	"identity/pkg/domain"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountUsecase struct {
	accountRepo  domain.AccountRepository
	eventLogRepo domain.EventLogRepository
	loginRepo    domain.LoginLogRepository
	ipDB         geoip2.Reader
}

func NewAccountUsecase(accountRepo domain.AccountRepository, eventLogRepo domain.EventLogRepository, loginRepo domain.LoginLogRepository, ipDB *geoip2.Reader) *AccountUsecase {
	return &AccountUsecase{
		accountRepo:  accountRepo,
		eventLogRepo: eventLogRepo,
		loginRepo:    loginRepo,
		ipDB:         *ipDB,
	}
}
func (uc *AccountUsecase) Account(ctx context.Context, namespace string, accountID uint64) (*domain.Account, error) {
	return uc.accountRepo.Account(ctx, namespace, accountID)
}

func (uc *AccountUsecase) AccountByUUID(ctx context.Context, namespace string, accountUUID string) (*domain.Account, error) {
	return uc.accountRepo.AccountByUUID(ctx, namespace, accountUUID)
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

	db := database.FromContext(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		ctx = database.ToContext(ctx, tx)

		err = uc.accountRepo.CreateAccount(ctx, account)
		if err != nil {
			return err
		}

		return uc.eventLogRepo.CreateEventLog(ctx, &domain.EventLog{
			Namespace: account.Namespace,
			Action:    "create",
			TargetID:  strconv.FormatUint(account.ID, 10),
			Actor:     account.CreatorName,
			Message:   "account is created",
			State:     domain.EventLogSuccess,
		})
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *AccountUsecase) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return uc.accountRepo.UpdateAccount(ctx, account)
}

func (uc *AccountUsecase) UpdateAccountPassword(ctx context.Context, request domain.UpdateAccountPasswordRequest) error {
	account, err := uc.accountRepo.Account(ctx, request.Namespace, request.AccountID)
	if err != nil {
		return err
	}

	//check old password
	err = isPasswordValid(account.PasswordEncrypt, request.OldPassword)
	if err != nil {
		return domain.ErrUsernameOrPasswordIncorrect
	}

	//update
	newPassword, err := encryptPassword(request.NewPassword)
	if err != nil {
		return err
	}
	account.PasswordEncrypt = newPassword
	account.UpdaterID = request.UpdaterID
	account.UpdaterName = request.UpdaterName
	account.UpdatedAt = time.Now().UTC()

	return uc.accountRepo.UpdateAccountPassword(ctx, account)
}

func (uc *AccountUsecase) ForceUpdateAccountPassword(ctx context.Context, request domain.ForceUpdateAccountPasswordRequest) error {
	account, err := uc.Account(ctx, request.Namespace, request.AccountID)
	if err != nil {
		return err
	}

	//update
	newPassword, err := encryptPassword(request.NewPassword)
	if err != nil {
		return err
	}

	updateAccount := domain.Account{
		ID:              account.ID,
		PasswordEncrypt: newPassword,
	}
	updateAccount.UpdaterID = request.UpdaterID
	updateAccount.UpdaterName = request.UpdaterName
	updateAccount.UpdatedAt = time.Now().UTC()

	return uc.accountRepo.UpdateAccountPassword(ctx, &updateAccount)
}

func (uc *AccountUsecase) ChangeState(ctx context.Context, request domain.ChangeStateRequest) error {
	account, err := uc.Account(ctx, request.Namespace, request.AccountID)
	if err != nil {
		return err
	}

	if account.State == request.State {
		return nil
	}

	db := database.FromContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		ctx = database.ToContext(ctx, tx)

		account.State = request.State
		account.UpdaterID = request.AccountID
		account.UpdaterName = request.UpdaterName

		err = uc.accountRepo.UpdateState(ctx, account)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("change state from %s to %s", account.State.String(), request.State.String())
		return uc.eventLogRepo.CreateEventLog(ctx, &domain.EventLog{
			Namespace: request.Namespace,
			Action:    request.State.String(),
			TargetID:  strconv.FormatUint(account.ID, 10),
			Actor:     account.UpdaterName,
			Message:   msg,
			State:     domain.EventLogSuccess,
		})
	})
}

func (uc *AccountUsecase) Login(ctx context.Context, request domain.LoginInfo) (*domain.Account, error) {
	if len(request.Namespace) == 0 || len(request.Username) == 0 {
		return nil, domain.ErrInvalidInput
	}

	options := domain.FindAccountOptions{
		Namespace: request.Namespace,
		Username:  request.Username,
	}
	accounts, err := uc.accountRepo.Accounts(ctx, options)

	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, domain.ErrUsernameOrPasswordIncorrect
	}

	account := accounts[0]

	//check status
	switch account.State {
	case domain.AccountStatusLocked:
		return nil, domain.ErrAccountLocked
	case domain.AccountStatusDisabled:
		return nil, domain.ErrAccountDisabled
	}

	db := database.FromContext(ctx)
	//compare password
	err = isPasswordValid(account.PasswordEncrypt, request.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {

			//帳密錯誤更新錯誤次數
			err = db.Transaction(func(tx *gorm.DB) error {
				account.FailedPasswordAttempt = account.FailedPasswordAttempt + 1
				err := uc.accountRepo.UpdateAccount(ctx, &account)
				if err != nil {
					return err
				}

				var contryCode, cityName string
				if len(request.ClientIP) > 0 {
					ip := net.ParseIP(request.ClientIP)
					record, err := uc.ipDB.City(ip)
					if err != nil {
						return err
					}

					if len(record.Subdivisions) > 0 {
						cityName = record.Subdivisions[0].Names["zh-CN"]
					}
					contryCode = record.Country.IsoCode
				}

				return uc.loginRepo.CreateLoginLog(ctx, &domain.LoginLog{
					Namespace:   request.Namespace,
					TargetID:    strconv.FormatUint(account.ID, 10),
					Actor:       account.Username,
					CountryCode: contryCode,
					CityName:    cityName,
					DeviceType:  request.DeviceType,
					State:       domain.LoginLogFail,
				})
			})

			if err != nil {
				return nil, err
			}

			return &account, domain.ErrUsernameOrPasswordIncorrect
		}

		return &account, err
	}

	//登入成功，清除登入失敗次數
	err = db.Transaction(func(tx *gorm.DB) error {
		account.FailedPasswordAttempt = 0
		account.LastLoginAt = time.Now().UTC()
		err = uc.accountRepo.UpdateAccount(ctx, &account)
		if err != nil {
			return err
		}

		var contryCode, cityName string
		if len(request.ClientIP) > 0 {
			ip := net.ParseIP(request.ClientIP)
			record, err := uc.ipDB.City(ip)
			if err != nil {
				return err
			}

			if len(record.Subdivisions) > 0 {
				cityName = record.Subdivisions[0].Names["zh-CN"]
			}
			contryCode = record.Country.IsoCode
		}

		return uc.loginRepo.CreateLoginLog(ctx, &domain.LoginLog{
			Namespace:   request.Namespace,
			TargetID:    strconv.FormatUint(account.ID, 10),
			Actor:       account.Username,
			CountryCode: contryCode,
			CityName:    cityName,
			DeviceType:  request.DeviceType,
			State:       domain.LoginLogSuccess,
		})
	})

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (uc *AccountUsecase) ClearOTP(ctx context.Context, accountUUID string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) GenerateOTPAuth(ctx context.Context, accountID uint64) (string, error) {
	panic("not implemented")
}

func (uc *AccountUsecase) SetOTPExpireTime(ctx context.Context, accountUUID string, duration int64) error {
	panic("not implemented")
}

func (uc *AccountUsecase) VerifyOTP(ctx context.Context, accountUUID string, otpCode string) (*domain.Account, error) {
	panic("not implemented")
}

//AccountIDsByRoleName(ctx context.Context, namespace, roleName string) ([]int64, error)
func (uc *AccountUsecase) UpdateAccountRole(ctx context.Context, accountID uint64, roles []int64, updaterAccountID uint64, updaterUsername string) error {
	panic("not implemented")
}

func (uc *AccountUsecase) AddRolesToAccount(ctx context.Context, request domain.AddRolesToAccountRequest) error {
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
