package mysql

import (
	"context"
	"identity/internal/pkg/database"
	"identity/pkg/domain"
	"time"

	"github.com/nite-coder/blackbear/pkg/log"
	"gorm.io/gorm"
)

type EventLogRepo struct {
}

func NewEventLogRepo(db *gorm.DB) *EventLogRepo {
	return &EventLogRepo{}
}

func (repo *EventLogRepo) CreateEventLog(ctx context.Context, eventLog *domain.EventLog) error {
	logger := log.FromContext(ctx)
	db := database.FromContext(ctx)

	eventLog.CreatedAt = time.Now().UTC()

	err := db.Create(eventLog).Error
	if err != nil {
		logger.Err(err).Interface("params", eventLog).Error("mysql: create eventLog fail")
		return err
	}

	return nil
}
