package events

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/app/command"
)

type Publisher struct {
	eventPublisher message.Publisher
}

func NewPublisher(publisher message.Publisher) Publisher {
	return Publisher{
		eventPublisher: publisher,
	}
}

func (p Publisher) PublishMonitorCreated(ctx context.Context, event command.MonitorCreatedEvent) error {
	m, err := mapEventToMessage(MonitorCreatedEvent{
		Header:    NewHeader(MonitorCreatedEvent{}.EventName(), ""),
		ID:        event.ID,
		CreatedAt: event.CreatedAt,
	})
	if err != nil {
		return err
	}

	err = p.eventPublisher.Publish(MonitorCreatedEvent{}.EventName(), m)
	if err != nil {
		return err
	}

	return nil
}

func (p Publisher) PublishEndpointCheckFailed(ctx context.Context, event command.EndpointCheckFailedEvent) error {
	m, err := mapEventToMessage(EndpointCheckFailed{
		At:        event.At,
		MonitorID: event.MonitorID,
		Header:    NewHeader(EndpointCheckFailed{}.EventName(), ""),
	})
	if err != nil {
		return err
	}

	err = p.eventPublisher.Publish(EndpointCheckFailed{}.EventName(), m)
	if err != nil {
		return err
	}

	return nil
}

func (p Publisher) PublishIncidentCreated(ctx context.Context, event command.IncidentCreatedEvent) error {
	m, err := mapEventToMessage(IncidentCreatedEvent{
		At:         event.At,
		IncidentID: event.IncidentID,
		MonitorID:  event.MonitorID,
		Header:     NewHeader(IncidentCreatedEvent{}.EventName(), ""),
	})
	if err != nil {
		return err
	}

	err = p.eventPublisher.Publish(IncidentCreatedEvent{}.EventName(), m)
	if err != nil {
		return err
	}

	return nil
}
