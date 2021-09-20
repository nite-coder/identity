package domain

import (
	"context"
	"time"
)

type EventLogState uint32

const (
	EventLogNone    EventLogState = 0
	EventLogSuccess EventLogState = 1
	EventLogFail    EventLogState = 2
)

var (
	TableNameEventLog = "event_logs"
)

type EventLog struct {
	ID        uint64        `gorm:"primaryKey;autoIncrement;not null"`
	Namespace string        `gorm:"column:namespace;type:string;size:256;not null"`
	Action    string        `gorm:"column:action;type:string;size:64;not null"`
	TargetID  string        `gorm:"column:target_id;type:string;size:256;not null"`
	Actor     string        `gorm:"column:actor;type:string;size:32;not null"`
	Message   string        `gorm:"column:namespace;type:string;size:256;not null"`
	State     EventLogState `gorm:"column:state;type:int;not null"`
	ClientIP  string        `gorm:"column:client_ip;type:string;size:64;not null"`
	CreatedAt time.Time     `gorm:"column:created_at;type:datetime;default:'1970-01-01 00:00:00';not null"`
}

func (e *EventLog) TableName() string {
	return TableNameEventLog
}

type EventLogRepository interface {
	CreateEventLog(ctx context.Context, eventLog *EventLog) error
}
