package command

import "time"

type MonitorCreatedEvent struct {
	ID        string
	CreatedAt time.Time
}

type EndpointCheckFailed struct {
	MonitorID string
	At        time.Time
}
