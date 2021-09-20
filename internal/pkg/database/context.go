package database

import (
	"context"
	"identity/internal/pkg/global"

	"gorm.io/gorm"
)

var (
	ctxKey = &struct {
		name string
	}{
		name: "db",
	}
)

func FromContext(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(ctxKey).(*gorm.DB)
	if !ok {
		return global.DB
	}
	return db
}

func ToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey, db)
}
