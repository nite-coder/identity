package grpc

import (
	"context"
	"identity/pkg/proto"
	"reflect"
	"testing"
)

func TestIdentityServer_CreateAccount(t *testing.T) {
	type args struct {
		ctx     context.Context
		account *proto.Account
	}

	createAccount := proto.Account{
		UUID: "123456",
	}

	tests := []struct {
		name    string
		s       *IdentityServer
		args    args
		want    *proto.Account
		wantErr bool
	}{
		{
			name: "simple",
			s:    _identityServer,
			args: args{
				ctx:     context.Background(),
				account: &createAccount,
			},
			want:    &createAccount,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateAccount(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("IdentityServer.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IdentityServer.CreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
