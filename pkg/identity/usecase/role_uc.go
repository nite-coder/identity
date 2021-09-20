package usecase

import (
	"context"
	"identity/pkg/domain"
)

type RoleUsecase struct {
	roleRepo domain.RoleRepository
}

func NewRoleUsecase(repo domain.RoleRepository) *RoleUsecase {
	return &RoleUsecase{
		roleRepo: repo,
	}
}

func (uc *RoleUsecase) CreateRole(ctx context.Context, role *domain.Role) error {
	return uc.roleRepo.CreateRole(ctx, role)
}

func (uc *RoleUsecase) UpdateRole(ctx context.Context, role *domain.Role) error {
	return uc.roleRepo.UpdateRole(ctx, role)
}

func (uc *RoleUsecase) Role(ctx context.Context, namespace string, id uint64) (*domain.Role, error) {
	return uc.roleRepo.Role(ctx, namespace, id)
}

func (uc *RoleUsecase) Roles(ctx context.Context, opts domain.FindRoleOptions) ([]domain.Role, error) {
	return uc.roleRepo.Roles(ctx, opts)
}

func (uc *RoleUsecase) AddAccountsToRole(ctx context.Context, accountIDs []uint64, roleID uint64) error {
	return uc.roleRepo.AddAccountsToRole(ctx, accountIDs, roleID)
}

func (uc *RoleUsecase) RolesByAccountID(ctx context.Context, namespace string, accountID uint64) ([]domain.Role, error) {
	return uc.roleRepo.RolesByAccountID(ctx, namespace, accountID)
}
