package grpc

import (
	"context"
	"reflect"
	"testing"

	"identity/pkg/proto"
)

func TestNewIdentityServer(t *testing.T) {
	tests := []struct {
		name string
		want *IdentityServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdentityServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdentityServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentityServer_CreateAccount(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *proto.Account
	}
	tests := []struct {
		name    string
		s       *IdentityServer
		args    args
		want    *proto.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateAccount(tt.args.in0, tt.args.in1)
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
