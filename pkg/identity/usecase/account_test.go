package usecase

import (
	"context"
	"identity/pkg/domain"
	identityMysql "identity/pkg/identity/repository/mysql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateAccount(t *testing.T) {
	gormConfig := gorm.Config{
		//PrepareStmt: true,
		Logger: logger.Default.LogMode(logger.Silent),
	}
	dsn := "root:root@tcp(mysql:3306)/identity_db?charset=utf8&parseTime=true&timeout=60s"
	db, err := gorm.Open(gormMySQL.Open(dsn), &gormConfig)
	require.NoError(t, err)

	err = db.Migrator().DropTable(domain.Account{})
	require.NoError(t, err)

	err = db.AutoMigrate(domain.Account{})
	require.NoError(t, err)

	accountRepo := identityMysql.NewAccountRepo(db)
	usecase := NewAccountUsecase(accountRepo)

	ctx := context.Background()

	account := domain.Account{
		Namespace:       "test.abc",
		Username:        "halo",
		PasswordEncrypt: "123456",
		CreatorID:       1,
		CreatorName:     "admin",
		State:           domain.AccountStatusNormal,
	}

	newAccount, err := usecase.CreateAccount(ctx, &account)
	require.NoError(t, err)

	//assert.Equal(t, 1, newAccount)
	assert.Equal(t, account.Username, newAccount.Username)

	newAccount1, err := usecase.Account(ctx, newAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, newAccount1.UUID, newAccount.UUID)

	newAccount1, err = usecase.AccountByUUID(ctx, newAccount.UUID)
	require.NoError(t, err)
	assert.Equal(t, newAccount1.Username, newAccount.Username)

	opts := domain.FindAccountOptions{
		Namespace: "test.abc",
	}

	count, err := usecase.CountAccounts(ctx, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
