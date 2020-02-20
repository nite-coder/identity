package identity

import (
	"context"
)

type Claims map[string]interface{}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Claims      Claims `json:"claims,omitempty"`
}

type LoginInfo struct {
	Namespace string `json:"namespace,omitempty" db:"namespace"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
}

func NewContext(ctx context.Context, claim Claims) context.Context {
	return context.WithValue(ctx, "identity_claims", claim)
}

func FromContext(ctx context.Context) (Claims, bool) {
	val, ok := ctx.Value("identity_claims").(Claims)
	if !ok {
		return nil, false
	}
	return val, true
}

type TokenServicer interface {
	Token(ctx context.Context, app, accessToken string) (*Token, error)
	DeleteToken(ctx context.Context, app, accessToken string) error
	CreateToken(ctx context.Context, login *LoginInfo) (*Token, error, int)
}
