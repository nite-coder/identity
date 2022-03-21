package domain

import (
	"context"
	"time"

	"gorm.io/datatypes"
)

type EventLogState uint32

const (
	EventLogDefault    EventLogState = 0
	EventLogSuccess EventLogState = 1
	EventLogFail    EventLogState = 2
)



type EventLog struct {
	ID        uint64         `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Namespace string         `gorm:"column:namespace;type:string;size:256;not null"`
	Action    string         `gorm:"column:action;type:string;size:64;not null"`
	TargetID  string         `gorm:"column:target_id;type:string;size:256;not null"`
	Message   string         `gorm:"column:message;type:string;size:512;not null"`
	OldStatus datatypes.JSON `gorm:"column:old_status;type:json;not null"`
	NewStatus datatypes.JSON `gorm:"column:new_status;type:json;not null"`
	State     EventLogState  `gorm:"column:state;type:int;not null"`
	ClientIP  string         `gorm:"column:client_ip;type:string;size:64;not null"`
	Actor     string         `gorm:"column:actor;type:string;size:32;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;default:1970-01-01 00:00:00;not null"`
}



type EventLogRepository interface {
	CreateEventLog(ctx context.Context, eventLog *EventLog) error
}
