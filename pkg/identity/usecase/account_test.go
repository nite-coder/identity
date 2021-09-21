package usecase

import (
	"context"
	"identity/internal/pkg/global"
	"identity/pkg/domain"
	identityMysql "identity/pkg/identity/repository/mysql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AccountTestSuite struct {
	suite.Suite
	id          string
	db          *gorm.DB
	accountRepo domain.AccountRepository
	usecase     domain.AccountUsecase
	namespace   string
}

func TestAccountTestSuite(t *testing.T) {
	var err error

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	gormConfig := gorm.Config{
		Logger: dbLogger,
	}

	dsn := "root:root@tcp(mysql:3306)/identity_db?charset=utf8mb4&parseTime=true&timeout=60s"
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
	if err != nil {
		panic(err)
	}

	global.DB = db

	eventLogRepo := identityMysql.NewEventLogRepo(global.DB)
	accountRepo := identityMysql.NewAccountRepo()

	usecase := NewAccountUsecase(accountRepo, eventLogRepo)

	accountTestSuite := AccountTestSuite{
		id:          uuid.NewString(),
		db:          db,
		accountRepo: accountRepo,
		usecase:     usecase,
		namespace:   "test.identity",
	}

	suite.Run(t, &accountTestSuite)

}

func (suite *AccountTestSuite) SetupTest() {
	domain.TableNameEventLog = "event_logs" + "_" + uuid.NewString()
	domain.TableNameAccount = "accounts" + "_" + uuid.NewString()
	domain.TableNameAccountRole = "accounts_roles" + "_" + uuid.NewString()
	domain.TableNameRoles = "roles" + "_" + uuid.NewString()
	domain.TableNamePermission = "permission" + "_" + uuid.NewString()

	err := suite.db.AutoMigrate(domain.EventLog{}, domain.Account{}, domain.Role{}, domain.Permission{})
	suite.Require().NoError(err)

}

func (suite *AccountTestSuite) TearDownTest() {
	domain.TableNameEventLog = "event_logs" + "_" + uuid.NewString()
	domain.TableNameAccount = "accounts" + "_" + uuid.NewString()
	domain.TableNameAccountRole = "accounts_roles" + "_" + uuid.NewString()
	domain.TableNameRoles = "roles" + "_" + uuid.NewString()
	domain.TableNamePermission = "permission" + "_" + uuid.NewString()

	err := suite.db.Migrator().DropTable(domain.EventLog{}, domain.Account{}, domain.Role{}, domain.Permission{})
	suite.Require().NoError(err)
}

func (suite *AccountTestSuite) TestCreateAccount() {
	ctx := context.Background()

	account := domain.Account{
		Namespace:       suite.namespace,
		Username:        "halo",
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	newAccount, err := suite.usecase.CreateAccount(ctx, &account)
	suite.Require().NoError(err)

	//assert.Equal(t, 1, newAccount)
	assert.Equal(suite.T(), account.Username, newAccount.Username)

	newAccount1, err := suite.usecase.Account(ctx, suite.namespace, newAccount.ID)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), newAccount1.UUID, newAccount.UUID)

	newAccount1, err = suite.usecase.AccountByUUID(ctx, suite.namespace, newAccount.UUID)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), newAccount1.Username, newAccount.Username)

	opts := domain.FindAccountOptions{
		Namespace: suite.namespace,
	}

	count, err := suite.usecase.CountAccounts(ctx, opts)
	suite.Require().NoError(err)
	suite.Assert().Equal(int64(1), count)
}

func (suite *AccountTestSuite) TestLogin() {
	ctx := context.Background()

	account := domain.Account{
		Namespace:       suite.namespace,
		Username:        "halo",
		PasswordEncrypt: "123456",
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	_, err := suite.usecase.CreateAccount(ctx, &account)
	suite.Require().NoError(err)

	suite.Run("username was not found", func() {
		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "no_this_user",
			Password:  "123456",
		}
		account, err := suite.usecase.Login(ctx, login)
		suite.Assert().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)
		suite.Assert().Nil(account)
	})

	suite.Run("login successfully", func() {
		now := time.Now()
		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "halo",
			Password:  "1111",
		}
		newAccount, err := suite.usecase.Login(ctx, login)
		suite.Assert().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)

		login.Password = "123456"
		newAccount, err = suite.usecase.Login(ctx, login)
		suite.Require().NoError(err)
		suite.Assert().Equal(account.Username, newAccount.Username)
		suite.Assert().True(now.Before(newAccount.LastLoginAt))
		suite.Assert().Equal(int32(0), newAccount.FailedPasswordAttempt)
	})

	suite.Run("login failed and account is locked", func() {
		var account *domain.Account
		var err error

		for i := 0; i < 5; i++ {
			login := domain.LoginInfo{
				Namespace: suite.namespace,
				Username:  "halo",
				Password:  "111",
			}
			account, err = suite.usecase.Login(ctx, login)
			suite.Assert().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)
		}

		suite.Assert().Equal(int32(5), account.FailedPasswordAttempt)

		changeStateRequest := domain.ChangeStateRequest{
			Namespace:   suite.namespace,
			AccountID:   account.ID,
			State:       domain.AccountStatusLocked,
			UpdaterID:   domain.SystemID,
			UpdaterName: domain.SystemName,
		}
		err = suite.usecase.ChangeState(ctx, changeStateRequest)
		suite.Require().NoError(err)

		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "halo",
			Password:  "123456",
		}
		_, err = suite.usecase.Login(ctx, login)
		suite.Assert().ErrorIs(err, domain.ErrAccountLocked)
	})
}

func (suite *AccountTestSuite) TestChangePassword() {
	ctx := context.Background()

	oldPassword := "123456"
	newPassword1 := "111111"
	newPassword2 := "222222"

	account := domain.Account{
		Namespace:       suite.namespace,
		Username:        "halo",
		PasswordEncrypt: oldPassword,
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	suite.Run("update account password", func() {
		_, err := suite.usecase.CreateAccount(ctx, &account)
		suite.Require().NoError(err)

		request := domain.UpdateAccountPasswordRequest{
			Namespace:   suite.namespace,
			AccountID:   account.ID,
			OldPassword: oldPassword,
			NewPassword: newPassword1,
			UpdaterID:   2,
			UpdaterName: "halo",
		}
		err = suite.usecase.UpdateAccountPassword(ctx, request)
		suite.Require().NoError(err)

		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "halo",
			Password:  newPassword1,
		}
		_, err = suite.usecase.Login(ctx, login)
		suite.Require().NoError(err)
	})

	suite.Run("force update account password", func() {
		request := domain.ForceUpdateAccountPasswordRequest{
			Namespace:   suite.namespace,
			AccountID:   account.ID,
			NewPassword: newPassword2,
			UpdaterID:   1,
			UpdaterName: "admin",
		}
		err := suite.usecase.ForceUpdateAccountPassword(ctx, request)
		suite.Require().NoError(err)

		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "halo",
			Password:  newPassword2,
		}
		_, err = suite.usecase.Login(ctx, login)
		suite.Require().NoError(err)
	})
}

func (suite *AccountTestSuite) TestChangeState() {
	ctx := context.Background()

	account := domain.Account{
		Namespace:       suite.namespace,
		Username:        "halo",
		PasswordEncrypt: "123456",
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	_, err := suite.usecase.CreateAccount(ctx, &account)
	suite.Require().NoError(err)

	changeStateRequest := domain.ChangeStateRequest{
		Namespace:   suite.namespace,
		AccountID:   account.ID,
		State:       domain.AccountStatusLocked,
		UpdaterID:   domain.SystemID,
		UpdaterName: domain.SystemName,
	}
	err = suite.usecase.ChangeState(ctx, changeStateRequest)
	suite.Require().NoError(err)

	newAccount, err := suite.usecase.Account(ctx, suite.namespace, account.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal(domain.AccountStatusLocked, newAccount.State)
}
