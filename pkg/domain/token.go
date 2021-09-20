package domain

import (
	"context"
	"time"

	errPKG "github.com/rotisserie/eris"
)

var (
	// ErrKeyNotFound returns errors.ResourceNotFound
	ErrKeyNotFound = errPKG.New("key not found")
)

// Claims 用來代表登入後的資料
// 儲存於redis中的account資訊, 由其他service帶入
type Claims map[string]interface{}

// Token 登入用的令牌
type Token struct {
	AccountID        int64
	Namespace        string
	ExpiresIn        int64
	TokenString      string
	Claims           map[string]string
	Username         string
	AccountType      int32
	RefreshExpiresIn int64
}

type key int

const (
	// IdentityClaims 用來取 context 裡面的 claims
	IdentityClaims key = iota

	// PairTokenKey 配對的TokenKey(refreshToken <-> accessToken)
	PairTokenKey = "PairTokenKey"
	// BindHashKey 額外綁定的hashKey(目前只支援 1對1)
	BindHashKey = "BindHashKey"
)

// NewContext 產生一個新的包含 claims 的新 context
func NewContext(ctx context.Context, claim Claims) context.Context {
	return context.WithValue(ctx, IdentityClaims, claim)
}

var (
	ctxKey = &struct {
		name string
	}{
		name: "log",
	}
)

// FromContext 從 context 裡面取得 claims
func FromContext(ctx context.Context) (Claims, bool) {
	val, ok := ctx.Value(ctxKey).(Claims)
	if !ok {
		return nil, false
	}
	return val, true
}

// TokenUsecase 用來處理 Token 相關業務操作的場景
type TokenUsecase interface {
	CreateToken(ctx context.Context, accessToken Token, prefixTokens ...string) (string, string, error)
	Token(ctx context.Context, tokenKey string) (*Token, error)
	RefreshToken(ctx context.Context, tokenKey string) (string, string, error)
	BindHashToken(ctx context.Context, hashKey, accessTokenKey string) error
	DeleteHash(ctx context.Context, hashKey string) error
	DeleteTokenByAccountID(ctx context.Context, accountID int64, prefixTokens ...string) error
	RenewToken(ctx context.Context, tokenKey string, duration int64) error
	CreateRefreshToken(ctx context.Context, token *Token) (string, error)
}

// TokenRepository 用來處理 token 物件存儲的行為 repository layer
type TokenRepository interface {
	CreateAccountHash(ctx context.Context, accountID, refreshKey, accessKey string, d time.Duration) error
	BindHashToken(ctx context.Context, hashKey, accessToken string) error
	DeleteHash(ctx context.Context, hashID string) error
	SetToken(ctx context.Context, prefix string, token Token, d time.Duration) (string, error)
	DeleteToken(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string) (string, string, error)
	GetToken(ctx context.Context, tokenString string) (Token, error)
	RenewToken(ctx context.Context, tokenString string, d time.Duration) error
	GetRcc() interface{}

	//need to suspend
	GetAuthToken(ctx context.Context, tokenString string) (*Token, error)
	SetAuthToken(ctx context.Context, token *Token, d time.Duration) (string, error)
	DeleteAuthTokenByAccountID(ctx context.Context, accountID int64) error
	CreateRefreshToken(ctx context.Context, token *Token, d time.Duration) (string, error)
}
