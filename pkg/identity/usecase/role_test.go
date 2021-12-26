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
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RoleTestSuite struct {
	suite.Suite
	db          *gorm.DB
	roleRepo    domain.RoleRepository
	accountRepo domain.AccountRepository
	usecase     domain.RoleUsecase
	namespace   string
}

func TestRoleTestSuite(t *testing.T) {
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

	dsn := "root:root@tcp(localhost:3306)/identity_db?charset=utf8mb4&parseTime=true&timeout=60s"
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
	if err != nil {
		panic(err)
	}

	global.DB = db

	roleRepo := identityMysql.NewRoleRepo()
	accountRepo := identityMysql.NewAccountRepo()
	usecase := NewRoleUsecase(roleRepo)

	roleTestSuite := RoleTestSuite{
		db:          db,
		roleRepo:    roleRepo,
		accountRepo: accountRepo,
		usecase:     usecase,
		namespace:   "test.identity",
	}

	suite.Run(t, &roleTestSuite)
}

func (suite *RoleTestSuite) SetupTest() {


	err := suite.db.Set("gorm:table_options", "AUTO_INCREMENT=100000").AutoMigrate(domain.Account{})
	suite.Require().NoError(err)
	err = suite.db.AutoMigrate(domain.EventLog{}, domain.AccountRole{}, domain.Role{}, domain.Permission{}, domain.LoginLog{})
	suite.Require().NoError(err)

}

func (suite *RoleTestSuite) TearDownTest() {


	err := suite.db.Migrator().DropTable(domain.EventLog{}, domain.Account{}, domain.AccountRole{}, domain.Role{}, domain.Permission{}, domain.LoginLog{})
	suite.Require().NoError(err)
}

func (suite *RoleTestSuite) TestCRUDRole() {
	ctx := context.Background()

	role := domain.Role{
		Namespace:   suite.namespace,
		Name:        "admin",
		State:       domain.RoleStatusNormal,
		CreatorID:   1,
		CreatorName: "admin",
	}

	err := suite.usecase.CreateRole(ctx, &role)
	suite.Require().NoError(err)

	suite.Equal(uint64(1), role.ID)

	role.Name = "users"
	role.Desc = "user groups"
	role.State = domain.RoleStatusDisabled
	role.UpdaterID = 1
	role.UpdaterName = "admin"

	err = suite.usecase.UpdateRole(ctx, &role)
	suite.Require().NoError(err)

	newRole, err := suite.roleRepo.Role(ctx, role.Namespace, role.ID)
	suite.Require().NoError(err)

	suite.Equal(role.Name, newRole.Name)
	suite.Equal(role.Desc, newRole.Desc)
	suite.Equal(role.State, newRole.State)
	suite.Equal(role.UpdaterID, newRole.UpdaterID)
	suite.Equal(role.UpdaterName, newRole.UpdaterName)
}

func (suite *RoleTestSuite) TestAddAccountsToRole() {
	ctx := context.Background()

	role := domain.Role{
		Namespace:   suite.namespace,
		Name:        "finance",
		State:       domain.RoleStatusNormal,
		CreatorID:   1,
		CreatorName: "admin",
	}

	err := suite.usecase.CreateRole(ctx, &role)
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

	account2 := domain.Account{
		Namespace: suite.namespace,
		UUID:      uuid.NewString(),
		Username: sql.NullString{
			String: "user002",
			Valid:  true,
		},
		FirstName:       "jordan",
		PasswordEncrypt: "123456",
		State:           domain.AccountStatusNormal,
		CreatorID:       1,
		CreatorName:     "admin",
	}

	err = suite.accountRepo.CreateAccount(ctx, &account2)
	suite.Require().NoError(err)

	accountIds := []uint64{account1.ID, account2.ID}

	err = suite.usecase.AddAccountsToRole(ctx, accountIds, role.ID)
	suite.Require().NoError(err)

	accounts, err := suite.accountRepo.AccountsByRoleID(ctx, suite.namespace, role.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal(2, len(accounts))
	suite.Assert().Equal("user001", accounts[0].Username.String)

	roles, err := suite.usecase.RolesByAccountID(ctx, suite.namespace, account2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal(1, len(roles))
	suite.Assert().Equal("finance", roles[0].Name)
}
