package grpc

import (
	"context"
	"identity/pkg/domain"
	identityProto "identity/pkg/identity/proto"
)

// IdentityServer is server
type IdentityServer struct {
	accountSvc domain.AccountUsecase
}

// NewIdentityServer generate a new identity server instance
func NewIdentityServer(accountSvc domain.AccountUsecase) *IdentityServer {
	return &IdentityServer{
		accountSvc: accountSvc,
	}
}
func (s *IdentityServer) Account(ctx context.Context, _ *identityProto.AccountRequest) (*identityProto.AccountResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) Accounts(ctx context.Context, _ *identityProto.AccountsRequest) (*identityProto.AccountsResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) CountAccounts(ctx context.Context, _ *identityProto.CountAccountsRequest) (*identityProto.CountAccountsResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) CreateAccount(ctx context.Context, _ *identityProto.CreateAccountRequest) (*identityProto.CreateAccountResponse, error) {

	// account := identityProto.Account{
	// 	Id: "123",
	// }

	resp := identityProto.CreateAccountResponse{
		//		Account: &account,
	}

	return &resp, nil
}

func (s *IdentityServer) UpdateAccount(ctx context.Context, _ *identityProto.UpdateAccountRequest) (*identityProto.UpdateAccountResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) UpdateAccountPassword(ctx context.Context, _ *identityProto.UpdateAccountPasswordRequest) (*identityProto.UpdateAccountPasswordResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) ForcedUpdatePassword(ctx context.Context, _ *identityProto.ForcedUpdatePasswordRequest) (*identityProto.ForcedUpdatePasswordResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) LockAccount(ctx context.Context, _ *identityProto.LockAccountRequest) (*identityProto.LockAccountResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) LockAccounts(ctx context.Context, _ *identityProto.LockAccountsRequest) (*identityProto.LockAccountsResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) UnlockAccount(ctx context.Context, _ *identityProto.UnlockAccountRequest) (*identityProto.UnlockAccountResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) DeleteAccount(ctx context.Context, _ *identityProto.DeleteAccountRequest) (*identityProto.DeleteAccountResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) Login(ctx context.Context, _ *identityProto.LoginRequest) (*identityProto.LoginResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) ClearOTP(ctx context.Context, _ *identityProto.ClearOTPRequest) (*identityProto.ClearOTPResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) GenerateOTPAuth(ctx context.Context, _ *identityProto.GenerateOTPAuthRequest) (*identityProto.GenerateOTPAuthResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) SetOTPExpireTime(ctx context.Context, _ *identityProto.SetOTPExpireTimeRequest) (*identityProto.SetOTPExpireTimeResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) VerifyOTP(ctx context.Context, _ *identityProto.VerifyOTPRequest) (*identityProto.VerifyOTPResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) Role(ctx context.Context, _ *identityProto.RoleRequest) (*identityProto.RoleResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) Roles(ctx context.Context, _ *identityProto.RolesRequest) (*identityProto.RolesResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) CreateRole(ctx context.Context, _ *identityProto.CreateRoleRequest) (*identityProto.CreateRoleResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) UpdateRole(ctx context.Context, _ *identityProto.UpdateRoleRequest) (*identityProto.UpdateRoleResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) UpdateAccountRole(ctx context.Context, _ *identityProto.UpdateAccountRoleRequest) (*identityProto.UpdateAccountRoleResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) AccountRoles(ctx context.Context, _ *identityProto.AccountRolesRequest) (*identityProto.AccountRolesResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) CreateToken(ctx context.Context, _ *identityProto.CreateTokenRequest) (*identityProto.CreateTokenResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) CreateRefreshToken(ctx context.Context, _ *identityProto.CreateRefreshTokenRequest) (*identityProto.CreateRefreshTokenResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) Token(ctx context.Context, _ *identityProto.TokenRequest) (*identityProto.TokenResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) DeleteTokenByRoleName(ctx context.Context, _ *identityProto.DeleteTokenByRoleNameRequest) (*identityProto.DeleteTokenByRoleNameResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) DeleteTokenByAccountID(ctx context.Context, _ *identityProto.DeleteTokenByAccountIDRequest) (*identityProto.DeleteTokenByAccountIDResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) RenewToken(ctx context.Context, _ *identityProto.RenewTokenRequest) (*identityProto.RenewTokenResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) RefreshToken(ctx context.Context, _ *identityProto.RefreshTokenRequest) (*identityProto.RefreshTokenResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) BindHashToken(ctx context.Context, _ *identityProto.BindHashTokenRequest) (*identityProto.BindHashTokenResponse, error) {
	panic("not implemented")
}

func (s *IdentityServer) DeleteHash(ctx context.Context, _ *identityProto.DeleteHashRequest) (*identityProto.DeleteHashResponse, error) {
	panic("not implemented")
}
