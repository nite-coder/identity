package usecase

import (
	"context"
	"database/sql"
	"identity/internal/pkg/global"
	"identity/pkg/domain"
	identityMysql "identity/pkg/identity/repository/mysql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type AccountTestSuite struct {
	suite.Suite
	id          string
	db          *gorm.DB
	accountRepo domain.AccountRepository
	roleRepo    domain.RoleRepository
	usecase     domain.AccountUsecase
	namespace   string
}

func TestAccountTestSuite(t *testing.T) {
	var err error

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	gormConfig := gorm.Config{
		Logger: dbLogger,
	}

	dsn := "root:root@tcp(localhost:3306)/identity_db?charset=utf8mb4&parseTime=true&timeout=60s"
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
	if err != nil {
		panic(err)
	}

	global.DB = db

	eventLogRepo := identityMysql.NewEventLogRepo()
	accountRepo := identityMysql.NewAccountRepo()
	roleRepo := identityMysql.NewRoleRepo()
	loginRepo := identityMysql.NewLoginLogRepo()

	filePath := "/workspace/data/geoip/geolite.mmdb" // TODO: remove hard code
	ipDB, err := geoip2.Open(filePath)
	if err != nil {
		panic(err)
	}

	usecase := NewAccountUsecase(accountRepo, eventLogRepo, loginRepo, ipDB)

	accountTestSuite := AccountTestSuite{
		id:          uuid.NewString(),
		db:          db,
		accountRepo: accountRepo,
		roleRepo:    roleRepo,
		usecase:     usecase,
		namespace:   "test.identity",
	}

	suite.Run(t, &accountTestSuite)

}

func (suite *AccountTestSuite) SetupTest() {
	prefix := uuid.NewString()
	suite.db.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   prefix + "_", // table name prefix, table for `User` would be `t_users`
		SingularTable: true,         // use singular table name, table for `User` would be `user` with this option enabled
	}

	err := suite.db.Set("gorm:table_options", "AUTO_INCREMENT=100000").AutoMigrate(domain.Account{})
	suite.Require().NoError(err)
	err = suite.db.AutoMigrate(domain.EventLog{}, domain.AccountRole{}, domain.Role{}, domain.Permission{}, domain.LoginLog{})
	suite.Require().NoError(err)

}

func (suite *AccountTestSuite) TearDownTest() {
	err := suite.db.Migrator().DropTable(domain.EventLog{}, domain.Account{}, domain.AccountRole{}, domain.Role{}, domain.Permission{}, domain.LoginLog{})
	suite.Require().NoError(err)
}

func (suite *AccountTestSuite) TestCRUDAccount() {
	ctx := context.Background()

	account1 := domain.Account{
		Namespace: suite.namespace,
		Username: sql.NullString{
			String: "halo",
			Valid:  true,
		},
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	err := suite.usecase.CreateAccount(ctx, &account1)
	suite.Require().NoError(err)

	account2 := domain.Account{
		Namespace: suite.namespace,
		MobileCountryCode: sql.NullString{
			String: "886",
			Valid:  true,
		},
		Mobile: sql.NullString{
			String: "966123456",
		},
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	err = suite.usecase.CreateAccount(ctx, &account2)
	suite.Require().NoError(err)

	account3 := domain.Account{
		Namespace: suite.namespace,
		Email: sql.NullString{
			String: "halo@no.com",
			Valid:  true,
		},
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	err = suite.usecase.CreateAccount(ctx, &account3)
	suite.Require().NoError(err)

	opts := domain.FindAccountOptions{
		Namespace: suite.namespace,
	}
	count, err := suite.usecase.CountAccounts(ctx, opts)
	suite.Require().NoError(err)
	suite.Assert().Equal(int64(3), count)

	newAccount1, err := suite.usecase.Account(ctx, suite.namespace, account1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal(newAccount1.Username, account1.Username)
	suite.Assert().Equal(newAccount1.State, account1.State)

	newAccount1, err = suite.usecase.AccountByUUID(ctx, suite.namespace, account1.UUID)
	suite.Require().NoError(err)
	suite.Assert().Equal(newAccount1.Username, newAccount1.Username)

	suite.Run("update account", func() {
		account3.NickName = "nickname_halo"
		account3.FirstName = "halo"
		account3.LastName = "hello"

		err = suite.usecase.UpdateAccount(ctx, &account3)
		suite.Require().NoError(err)

		newAccount3, err := suite.usecase.Account(ctx, suite.namespace, account3.ID)
		suite.Require().NoError(err)

		suite.Assert().Equal(account3.NickName, newAccount3.NickName)
		suite.Assert().Equal(account3.FirstName, newAccount3.FirstName)
		suite.Assert().Equal(account3.LastName, newAccount3.LastName)
	})

}

func (suite *AccountTestSuite) TestLogin() {
	ctx := context.Background()

	account1 := domain.Account{
		Namespace: suite.namespace,
		Username: sql.NullString{
			String: "halo",
			Valid:  true,
		},
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	err := suite.usecase.CreateAccount(ctx, &account1)
	suite.Require().NoError(err)

	account2 := domain.Account{
		Namespace: suite.namespace,
		MobileCountryCode: sql.NullString{
			String: "886",
			Valid:  true,
		},
		Mobile: sql.NullString{
			String: "966123456",
			Valid:  true,
		},
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	err = suite.usecase.CreateAccount(ctx, &account2)
	suite.Require().NoError(err)

	account3 := domain.Account{
		Namespace: suite.namespace,
		Email: sql.NullString{
			String: "halo@no.com",
			Valid:  true,
		},
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	err = suite.usecase.CreateAccount(ctx, &account3)
	suite.Require().NoError(err)

	opts := domain.FindAccountOptions{
		Namespace: suite.namespace,
	}
	count, err := suite.usecase.CountAccounts(ctx, opts)
	suite.Require().NoError(err)
	suite.Assert().Equal(int64(3), count)

	suite.Run("username was not found", func() {
		login := domain.LoginInfo{
			Namespace:  suite.namespace,
			LoginType:  domain.LoginTypeUsername,
			Username:   "no_this_user",
			Password:   "123456",
			ClientIP:   "182.48.113.104",
			DeviceType: domain.DeviceTypeWeb,
		}
		account, err := suite.usecase.Login(ctx, login)
		suite.Require().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)
		suite.Assert().Nil(account)
	})

	suite.Run("username login", func() {
		now := time.Now()
		login := domain.LoginInfo{
			Namespace: suite.namespace,
			LoginType: domain.LoginTypeUsername,
			Username:  "halo",
			Password:  "1111",
			ClientIP:  "182.48.113.104",
		}
		newAccount, err := suite.usecase.Login(ctx, login)
		suite.Require().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)

		login.Password = "123456"
		newAccount, err = suite.usecase.Login(ctx, login)
		suite.Require().NoError(err)
		suite.Assert().Equal(account1.Username, newAccount.Username)
		suite.Assert().True(now.Before(newAccount.LastLoginAt))
		suite.Assert().Equal(int32(0), newAccount.FailedPasswordAttempt)
	})

	suite.Run("mobile login", func() {
		now := time.Now()
		login := domain.LoginInfo{
			Namespace:         suite.namespace,
			LoginType:         domain.LoginTypeMobile,
			MobileCountryCode: "886",
			Mobile:            "966123456",
			Password:          "1111",
			ClientIP:          "182.48.113.104",
		}
		newAccount, err := suite.usecase.Login(ctx, login)
		suite.Require().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)

		login.Password = "123456"
		newAccount, err = suite.usecase.Login(ctx, login)
		suite.Require().NoError(err)
		suite.Assert().Equal(account2.Mobile, newAccount.Mobile)
		suite.Assert().True(now.Before(newAccount.LastLoginAt))
		suite.Assert().Equal(int32(0), newAccount.FailedPasswordAttempt)
	})

	suite.Run("email login", func() {
		now := time.Now()
		login := domain.LoginInfo{
			Namespace: suite.namespace,
			LoginType: domain.LoginTypeEmail,
			Email:     "halo@no.com",
			Password:  "1111",
			ClientIP:  "182.48.113.104",
		}
		newAccount, err := suite.usecase.Login(ctx, login)
		suite.Require().ErrorIs(err, domain.ErrUsernameOrPasswordIncorrect)

		login.Password = "123456"
		newAccount, err = suite.usecase.Login(ctx, login)
		suite.Require().NoError(err)
		suite.Assert().Equal(account3.Email, newAccount.Email)
		suite.Assert().True(now.Before(newAccount.LastLoginAt))
		suite.Assert().Equal(int32(0), newAccount.FailedPasswordAttempt)
	})

	suite.Run("login failed and account is locked", func() {
		var account *domain.Account
		var err error

		for i := 0; i < 5; i++ {
			login := domain.LoginInfo{
				Namespace: suite.namespace,
				LoginType: domain.LoginTypeUsername,
				Username:  "halo",
				Password:  "111",
				ClientIP:  "",
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
			LoginType: domain.LoginTypeUsername,
			Username:  "halo",
			Password:  "123456",
			ClientIP:  "182.48.113.104",
		}
		_, err = suite.usecase.Login(ctx, login)
		suite.Require().ErrorIs(err, domain.ErrAccountLocked)
	})
}

func (suite *AccountTestSuite) TestChangePassword() {
	ctx := context.Background()

	oldPassword := "123456"
	newPassword1 := "111111"
	newPassword2 := "222222"

	account := domain.Account{
		Namespace: suite.namespace,
		Username: sql.NullString{
			String: "halo",
			Valid:  true,
		},
		PasswordEncrypt: oldPassword,
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	suite.Run("update account password", func() {
		err := suite.usecase.CreateAccount(ctx, &account)
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
			LoginType: domain.LoginTypeUsername,
			Username:  "halo",
			Password:  newPassword1,
			ClientIP:  "",
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
			LoginType: domain.LoginTypeUsername,
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
		Namespace: suite.namespace,
		Username: sql.NullString{
			String: "halo",
			Valid:  true,
		},
		PasswordEncrypt: "123456",
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	err := suite.usecase.CreateAccount(ctx, &account)
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

func (suite *AccountTestSuite) TestAddRoleToAccount() {
	ctx := context.Background()

	role1 := domain.Role{
		Namespace:   suite.namespace,
		Name:        "finance1",
		State:       domain.RoleStatusNormal,
		CreatorID:   1,
		CreatorName: "admin",
	}

	err := suite.roleRepo.CreateRole(ctx, &role1)
	suite.Require().NoError(err)

	role2 := domain.Role{
		Namespace:   suite.namespace,
		Name:        "finance2",
		State:       domain.RoleStatusNormal,
		CreatorID:   1,
		CreatorName: "admin",
	}

	err = suite.roleRepo.CreateRole(ctx, &role2)
	suite.Require().NoError(err)

	account1 := domain.Account{
		Namespace: suite.namespace,
		UUID:      uuid.NewString(),
		Username: sql.NullString{
			String: "user001",
			Valid:  true,
		},
		FirstName:       "angela",
		LastName:        "wang",
		PasswordEncrypt: "123456",
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	err = suite.accountRepo.CreateAccount(ctx, &account1)
	suite.Require().NoError(err)

	roleIDs := []uint64{role1.ID, role2.ID}

	AddRolesToAccountRequest := domain.AddRolesToAccountRequest{
		Namespace:   suite.namespace,
		RoleIDs:     roleIDs,
		AccountID:   account1.ID,
		UpdaterID:   1,
		UpdaterName: "admin",
	}

	err = suite.usecase.AddRolesToAccount(ctx, AddRolesToAccountRequest)
	suite.Require().NoError(err)

	roles, err := suite.roleRepo.RolesByAccountID(ctx, suite.namespace, account1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal(2, len(roles))
	suite.Assert().Equal("finance1", roles[0].Name)
	suite.Assert().Equal("finance2", roles[1].Name)
}
