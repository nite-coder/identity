package grpc

import (
	"context"
	"fmt"
	identityProto "identity/pkg/identity/proto"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	fmt.Println("test account")

	ctx := context.Background()

	account := identityProto.Account{
		Id: "123",
	}

	createAccountReq := identityProto.CreateAccountRequest{
		Account: &account,
	}

	createAccountResp, err := _identityClient.CreateAccount(ctx, &createAccountReq)
	require.Nil(t, err)
	require.Equal(t, "123", createAccountResp.Account.Id)

}

func TestRole(t *testing.T) {
	fmt.Println("test role")
}
