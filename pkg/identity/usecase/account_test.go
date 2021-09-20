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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AccountTestSuite struct {
	suite.Suite
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
			Colorful:                  false,         // Disable color
		},
	)

	gormConfig := gorm.Config{
		Logger: dbLogger,
	}

	dsn := "root:root@tcp(mysql:3306)/identity_db?charset=utf8&parseTime=true&timeout=60s"
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
	if err != nil {
		panic(err)
	}

	global.DB = db

	eventLogRepo := identityMysql.NewEventLogRepo(global.DB)
	accountRepo := identityMysql.NewAccountRepo()

	usecase := NewAccountUsecase(accountRepo, eventLogRepo)

	accountTestSuite := AccountTestSuite{
		db:          db,
		accountRepo: accountRepo,
		usecase:     usecase,
		namespace:   "test.identity",
	}

	suite.Run(t, &accountTestSuite)
}

func (suite *AccountTestSuite) SetupTest() {
	err := suite.db.Migrator().DropTable(domain.EventLog{}, domain.Account{}, domain.Role{}, domain.Permission{})
	suite.Require().NoError(err)

	err = suite.db.AutoMigrate(domain.EventLog{}, domain.Account{}, domain.Role{})
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

	newAccount1, err := suite.usecase.Account(ctx, newAccount.ID)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), newAccount1.UUID, newAccount.UUID)

	newAccount1, err = suite.usecase.AccountByUUID(ctx, newAccount.UUID)
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
			assert.ErrorIs(suite.T(), err, domain.ErrUsernameOrPasswordIncorrect)
		}

		assert.Equal(suite.T(), int32(5), account.FailedPasswordAttempt)
		err = suite.usecase.LockAccount(ctx, account.ID)
		suite.Require().NoError(err)
	})
}
