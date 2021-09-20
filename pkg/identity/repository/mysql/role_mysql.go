package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"identity/internal/pkg/database"
	"identity/pkg/domain"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nite-coder/blackbear/pkg/log"
	"gorm.io/gorm"
)

type RoleRepo struct {
}

func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}

func (repo *RoleRepo) CreateRole(ctx context.Context, role *domain.Role) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	role.CreatedAt = time.Now().UTC()

	if err := db.Create(role).Error; err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if mysqlErr.Number == 1062 {
				return fmt.Errorf("mysql: the role has already exists.  %w", domain.ErrAlreadyExists)
			}
		}
		logger.Err(err).Interface("params", role).Error("mysql: create role fail")
		return err
	}

	return nil
}

func (repo *RoleRepo) UpdateRole(ctx context.Context, role *domain.Role) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	args := make(map[string]interface{})
	args["name"] = role.Name
	args["desc"] = role.Desc
	args["state"] = role.State
	args["updater_id"] = role.UpdaterID
	args["updater_name"] = role.UpdaterName
	args["updated_at"] = time.Now().UTC()

	result := db.Model(role).
		Where("id = ?", role.ID).
		Where("version = @version", sql.Named("version", role.Version)).
		UpdateColumns(args).
		UpdateColumn("version", gorm.Expr("version + 1"))

	err := result.Error
	if err != nil {
		logger.Err(err).Error("mysql: update role failed")
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrStale
	}

	return nil
}

func (repo *RoleRepo) Role(ctx context.Context, namespace string, id uint64) (*domain.Role, error) {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	role := domain.Role{}
	err := db.Model(domain.Role{}).Where("id = ?", id).Where("namespace = ?", namespace).Find(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("mysql: role id %d was not found. %w", id, domain.ErrNotFound)
		}
		logger.Err(err).Interface("params", role).Error("mysql: get role fail")
		return nil, err
	}

	return &role, nil
}

func (repo *RoleRepo) RolesByAccountID(ctx context.Context, namespace string, accountID uint64) ([]domain.Role, error) {
	//logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	roles := []domain.Role{}
	err := db.Model(domain.Role{}).
		Joins("left join accounts_roles on accounts_roles.role_id = roles.id").
		Where("account_id = ?", accountID).
		Where("namespace = ?", namespace).
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("mysql: get role by account id failed. %w", err)
	}

	return roles, nil
}

func (repo *RoleRepo) Roles(ctx context.Context, opts domain.FindRoleOptions) ([]domain.Role, error) {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	var roles []domain.Role

	db = repo.buildWhereClause(db, opts)

	if err := db.Find(&roles).Error; err != nil {
		logger.Err(err).Interface("params", opts).Error("mysql: get roles failed")
		return nil, err
	}

	return roles, nil
}

func (repo *RoleRepo) AddAccountsToRole(ctx context.Context, accountIDs []uint64, roleID uint64) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	err := db.Where("role_id = ?", roleID).Delete(&domain.AccountRole{}).Error
	if err != nil {
		logger.Err(err).Error("mysql: delete role failed")
		return err
	}

	accountRoles := []domain.AccountRole{}

	for _, accountID := range accountIDs {
		accountRole := domain.AccountRole{
			AccountID: accountID,
			RoleID:    roleID,
		}
		accountRoles = append(accountRoles, accountRole)
	}

	err = db.CreateInBatches(accountRoles, 100).Error
	if err != nil {
		logger.Err(err).Error("mysql: add accounts to role failed.")
		return err
	}

	return nil
}

func (repo *RoleRepo) buildWhereClause(db *gorm.DB, opts domain.FindRoleOptions) *gorm.DB {
	if opts.Namespace != "" {
		db = db.Where("namespace = ?", opts.Namespace)
	}

	if opts.Name != "" {
		db = db.Where("name = ?", opts.Name)
	}

	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}

	if opts.Offset > 0 {
		db = db.Offset(opts.Offset)
	}

	if opts.Sort != "" {
		db = db.Order(opts.Sort)
	}

	return db
}
