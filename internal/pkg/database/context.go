package database

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

var (
	ctxKey = &struct {
		name string
	}{
		name: "db",
	}
)

var _db *gorm.DB
var _isMock bool

func SetDB(db *gorm.DB) {
	_db = db
}

func SetMockMode(flag bool) {
	_isMock = flag
}

func FromContext(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(ctxKey).(*gorm.DB)
	if !ok {
		return _db
	}
	return db
}

func ToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey, db)
}

func Transaction(ctx context.Context, unitOfWork func(ctx context.Context) error, opts ...sql.TxOptions) error {
	if _isMock {
		return unitOfWork(ctx)
	}

	db := FromContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		ctx = ToContext(ctx, tx)
		return unitOfWork(ctx)
	})
}
