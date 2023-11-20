package command

import (
	"context"
	"time"

	"github.com/flowck/doberman/internal/domain"
)

type EventPublisher interface {
	PublishMonitorsEnqueued(ctx context.Context, event MonitorsEnqueuedEvent) error
	PublishMonitorEndpointCheckFailed(ctx context.Context, monitorID domain.ID) error
	PublishMonitorEndpointCheckSucceeded(ctx context.Context, monitorID domain.ID) error
	PublishIncidentResolved(ctx context.Context, incidentID domain.ID) error
	PublishIncidentCreated(ctx context.Context, incidentID domain.ID) error
	PublishIncidentEscalated(ctx context.Context, incidentID domain.ID) error
	PublishCommentOnIncidentPosted(ctx context.Context, commentID domain.ID) error
	PublishAccountCreated(ctx context.Context, event AccountCreatedEvent) error
}

type Event interface {
	EventName() string
}

type MonitorsEnqueuedEvent struct {
	IDs        []string  `json:"id"`
	EnqueuedAt time.Time `json:"enqueued_at"`
}

func (m MonitorsEnqueuedEvent) EventName() string {
	return "MonitorEnqueued_v1"
}

type AccountCreatedEvent struct {
	ID string `json:"id"`
}

func (a AccountCreatedEvent) EventName() string {
	return "AccountCreatedEvent_v1"
}
