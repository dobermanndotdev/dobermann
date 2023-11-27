package command

import "time"

type MonitorCreatedEvent struct {
	ID        string
	CreatedAt time.Time
}
