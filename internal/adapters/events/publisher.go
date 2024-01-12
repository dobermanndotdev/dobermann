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
	m, err := mapEventToMessage(ctx, MonitorCreatedEvent{
		Header:    NewHeader(ctx, MonitorCreatedEvent{}.EventName()),
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
	m, err := mapEventToMessage(ctx, EndpointCheckFailed{
		Header:        NewHeader(ctx, EndpointCheckFailed{}.EventName()),
		At:            event.At,
		MonitorID:     event.MonitorID,
		CheckedURL:    event.CheckedURL,
		CheckResultID: event.CheckResultID,
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
	m, err := mapEventToMessage(ctx, IncidentCreatedEvent{
		At:         event.At,
		IncidentID: event.IncidentID,
		MonitorID:  event.MonitorID,
		Header:     NewHeader(ctx, IncidentCreatedEvent{}.EventName()),
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

func (p Publisher) PublishEndpointCheckSucceeded(ctx context.Context, event command.EndpointCheckSucceededEvent) error {
	m, err := mapEventToMessage(ctx, EndpointCheckSucceededEvent{
		At:        event.At,
		MonitorID: event.MonitorID,
		Header:    NewHeader(ctx, EndpointCheckSucceededEvent{}.EventName()),
	})

	if err != nil {
		return err
	}

	err = p.eventPublisher.Publish(EndpointCheckSucceededEvent{}.EventName(), m)
	if err != nil {
		return err
	}

	return nil
}

func (p Publisher) PublishIncidentResolved(ctx context.Context, event command.IncidentResolvedEvent) error {
	m, err := mapEventToMessage(ctx, IncidentResolvedEvent{
		At:         event.At,
		MonitorID:  event.MonitorID,
		IncidentID: event.IncidentID,
		Header:     NewHeader(ctx, IncidentResolvedEvent{}.EventName()),
	})
	if err != nil {
		return err
	}

	err = p.eventPublisher.Publish(IncidentResolvedEvent{}.EventName(), m)
	if err != nil {
		return err
	}

	return nil
}
