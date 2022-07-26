package mysql

import (
	"context"
	"identity/internal/pkg/database"
	"identity/pkg/domain"
	"time"

	"github.com/nite-coder/blackbear/pkg/log"
)

type LoginLogRepo struct {
}

func NewLoginLogRepo() *EventLogRepo {
	return &EventLogRepo{}
}

func (repo *EventLogRepo) CreateLoginLog(ctx context.Context, loginLog *domain.LoginLog) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	loginLog.CreatedAt = time.Now().UTC()

	err := db.Create(loginLog).Error
	if err != nil {
		logger.Err(err).Any("params", loginLog).Error("mysql: create login log fail")
		return err
	}

	return nil
}
