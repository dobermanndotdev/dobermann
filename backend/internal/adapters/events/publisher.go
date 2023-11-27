package events

import (
	"context"
	"fmt"

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
		return fmt.Errorf("unable to publish event: %v", err)
	}

	return nil
}
