package grpc

import (
	"context"
	identity "identity/pkg"

	"identity/pkg/proto"
)

// IdentityServer is server
type IdentityServer struct {
	accountSvc identity.AccountServicer
}

// NewIdentityServer generate a new identity server instance
func NewIdentityServer(accountSvc identity.AccountServicer) *IdentityServer {
	return &IdentityServer{
		accountSvc: accountSvc,
	}
}

// CreateAccount function create an account
func (s *IdentityServer) CreateAccount(ctx context.Context, account *proto.Account) (*proto.Account, error) {
	account.ID = 9
	return account, nil
}
