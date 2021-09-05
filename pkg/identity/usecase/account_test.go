package usecase

import (
	"context"
	"identity/pkg/domain"
	identityMysql "identity/pkg/identity/repository/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	gormConfig := gorm.Config{
		//PrepareStmt: true,
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsn := "root:root@tcp(mysql:3306)/identity_db?charset=utf8&parseTime=true&timeout=60s"
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
	if err != nil {
		panic(err)
	}

	accountRepo := identityMysql.NewAccountRepo(db)
	usecase := NewAccountUsecase(accountRepo)

	accountTestSuite := AccountTestSuite{
		db:          db,
		accountRepo: accountRepo,
		usecase:     usecase,
		namespace:   "test.identity",
	}

	suite.Run(t, &accountTestSuite)
}

func (suite *AccountTestSuite) SetupTest() {
	err := suite.db.Migrator().DropTable(domain.Account{})
	require.NoError(suite.T(), err)

	err = suite.db.AutoMigrate(domain.Account{})
	require.NoError(suite.T(), err)
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
	require.NoError(suite.T(), err)

	//assert.Equal(t, 1, newAccount)
	assert.Equal(suite.T(), account.Username, newAccount.Username)

	newAccount1, err := suite.usecase.Account(ctx, newAccount.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), newAccount1.UUID, newAccount.UUID)

	newAccount1, err = suite.usecase.AccountByUUID(ctx, newAccount.UUID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), newAccount1.Username, newAccount.Username)

	opts := domain.FindAccountOptions{
		Namespace: suite.namespace,
	}

	count, err := suite.usecase.CountAccounts(ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count)
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
	require.NoError(suite.T(), err)

	suite.Run("username was not found", func() {
		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "no_this_user",
			Password:  "123456",
		}
		account, err := suite.usecase.Login(ctx, login)
		assert.ErrorIs(suite.T(), err, domain.ErrUsernameOrPasswordIncorrect)
		assert.Nil(suite.T(), account)
	})

	suite.Run("login successfully", func() {
		now := time.Now()
		login := domain.LoginInfo{
			Namespace: suite.namespace,
			Username:  "halo",
			Password:  "1111",
		}
		newAccount, err := suite.usecase.Login(ctx, login)
		assert.ErrorIs(suite.T(), err, domain.ErrUsernameOrPasswordIncorrect)

		login.Password = "123456"
		newAccount, err = suite.usecase.Login(ctx, login)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), account.Username, newAccount.Username)
		assert.True(suite.T(), now.Before(newAccount.LastLoginAt))
		assert.Equal(suite.T(), int32(0), newAccount.FailedPasswordAttempt)
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
		require.NoError(suite.T(), err)
	})
}
