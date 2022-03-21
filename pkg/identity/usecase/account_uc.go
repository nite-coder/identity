package usecase

import (
	"context"
	"encoding/json"
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
	"gorm.io/datatypes"
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

func (uc *AccountUsecase) CreateAccount(ctx context.Context, account *domain.Account) error {
	var err error

	if account.Namespace == "" {
		return fmt.Errorf("namespace can't empty. %w", domain.ErrInvalidInput)
	}

	if !account.Username.Valid && (!account.MobileCountryCode.Valid && !account.Mobile.Valid) && !account.Email.Valid {
		return fmt.Errorf("login name can't be empty. %w", domain.ErrInvalidInput)
	}

	if account.PasswordEncrypt == "" {
		return fmt.Errorf("password can't be empty. %w", domain.ErrInvalidInput)
	}

	if account.CreatorName == "" {
		return fmt.Errorf("creator name can't empty. %w", domain.ErrInvalidInput)
	}

	account.UUID = uuid.NewString()
	account.PasswordEncrypt, err = encryptPassword(account.PasswordEncrypt)
	account.OTPLastResetAt = time.Unix(0, 0)
	account.LastLoginAt = time.Unix(0, 0)
	if err != nil {
		return err
	}

	db := database.FromContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		ctx = database.ToContext(ctx, tx)

		err = uc.accountRepo.CreateAccount(ctx, account)
		if err != nil {
			return err
		}

		newStatus, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		return uc.eventLogRepo.CreateEventLog(ctx, &domain.EventLog{
			Namespace: "identity.account",
			Action:    "create",
			TargetID:  strconv.FormatUint(account.ID, 10),
			Message:   "account is created",
			OldStatus: datatypes.JSON([]byte("{}")),
			NewStatus: newStatus,
			State:     domain.EventLogSuccess,
			Actor:     account.CreatorName,
		})
	})
}

func (uc *AccountUsecase) UpdateAccount(ctx context.Context, request *domain.Account) error {
	account, err := uc.accountRepo.Account(ctx, request.Namespace, request.ID)
	if err != nil {
		return err
	}

	if account.Version != request.Version {
		return domain.ErrStale
	}

	db := database.FromContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		ctx = database.ToContext(ctx, tx)

		oldStatus, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		err = uc.accountRepo.UpdateAccount(ctx, request)
		if err != nil {
			return err
		}

		newStatus, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		return uc.eventLogRepo.CreateEventLog(ctx, &domain.EventLog{
			Namespace: "identity.account",
			Action:    "update",
			TargetID:  strconv.FormatUint(account.ID, 10),
			Message:   "update account",
			OldStatus: oldStatus,
			NewStatus: newStatus,
			State:     domain.EventLogSuccess,
			Actor:     request.UpdaterName,
		})
	})
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
	account, err := uc.accountRepo.Account(ctx, request.Namespace, request.AccountID)
	if err != nil {
		return err
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

		oldStatus, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		account.State = request.State
		account.UpdaterID = request.AccountID
		account.UpdaterName = request.UpdaterName

		err = uc.accountRepo.UpdateState(ctx, account)
		if err != nil {
			return err
		}

		newStatus, err := json.Marshal(&account)
		if err != nil {
			return err
		}

		return uc.eventLogRepo.CreateEventLog(ctx, &domain.EventLog{
			Namespace: "identity.account",
			Action:    "change_state",
			TargetID:  strconv.FormatUint(account.ID, 10),
			Message:   fmt.Sprintf("change state from %s to %s", account.State.String(), request.State.String()),
			OldStatus: oldStatus,
			NewStatus: newStatus,
			State:     domain.EventLogSuccess,
			Actor:     account.UpdaterName,
		})
	})
}

func (uc *AccountUsecase) Login(ctx context.Context, request domain.LoginInfo) (*domain.Account, error) {
	if len(request.Namespace) == 0 || request.LoginType == domain.LoginTypeDefault {
		return nil, domain.ErrInvalidInput
	}

	opts := domain.FindAccountOptions{
		Namespace: request.Namespace,
	}

	switch request.LoginType {
	case domain.LoginTypeUsername:
		if len(request.Username) == 0 {
			return nil, domain.ErrInvalidInput
		}

		opts.Username = request.Username
	case domain.LoginTypeEmail:
		if len(request.Email) == 0 {
			return nil, domain.ErrInvalidInput
		}

		opts.Email = request.Email
	case domain.LoginTypeMobile:
		if len(request.MobileCountryCode) == 0 && len(request.Mobile) == 0 {
			return nil, domain.ErrInvalidInput
		}

		opts.MobileCountryCode = request.MobileCountryCode
		opts.Mobile = request.Mobile
	}
	accounts, err := uc.accountRepo.Accounts(ctx, opts)

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

func (uc *AccountUsecase) ResetOTPSecret(ctx context.Context, request domain.ResetOTPSecretRequest) (string, error) {
	panic("not implemented")
}

func (uc *AccountUsecase) VerifyOTP(ctx context.Context, accountUUID string, request domain.VerifyOTPRequest) error {
	panic("not implemented")
}

func (uc *AccountUsecase) AddRolesToAccount(ctx context.Context, request domain.AddRolesToAccountRequest) error {
	return uc.accountRepo.AddRolesToAccount(ctx, request)
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
