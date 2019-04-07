package grpc

import (
	"context"
	identity "identity/pkg"

	"identity/pkg/proto"
)

type IdentityServer struct {
	accountSvc identity.AccountServicer
}

func NewIdentityServer(accountSvc identity.AccountServicer) *IdentityServer {
	return &IdentityServer{
		accountSvc: accountSvc,
	}
}

func (s *IdentityServer) CreateAccount(ctx context.Context, account *proto.Account) (*proto.Account, error) {
	account.ID = 9
	return account, nil
}
