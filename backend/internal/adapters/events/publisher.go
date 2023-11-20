package events

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/flowck/doberman/internal/app/command"
	"github.com/flowck/doberman/internal/domain"
)

type Publisher struct {
	publisher message.Publisher
}

func (p Publisher) PublishAccountCreated(ctx context.Context, event command.AccountCreatedEvent) error {
	msg, err := mapToMessage(event)
	if err != nil {
		return err
	}

	return p.publisher.Publish(event.EventName(), msg)
}

func NewPublisher(publisher message.Publisher) Publisher {
	return Publisher{
		publisher: publisher,
	}
}

func (p Publisher) PublishMonitorsEnqueued(ctx context.Context, event command.MonitorsEnqueuedEvent) error {
	msg, err := mapToMessage(event)
	if err != nil {
		return err
	}

	return p.publisher.Publish(event.EventName(), msg)
}

func (p Publisher) PublishMonitorEndpointCheckFailed(ctx context.Context, monitorID domain.ID) error {
	//TODO implement me
	panic("implement me")
}

func (p Publisher) PublishMonitorEndpointCheckSucceeded(ctx context.Context, monitorID domain.ID) error {
	//TODO implement me
	panic("implement me")
}

func (p Publisher) PublishIncidentResolved(ctx context.Context, incidentID domain.ID) error {
	//TODO implement me
	panic("implement me")
}

func (p Publisher) PublishIncidentCreated(ctx context.Context, incidentID domain.ID) error {
	//TODO implement me
	panic("implement me")
}

func (p Publisher) PublishIncidentEscalated(ctx context.Context, incidentID domain.ID) error {
	//TODO implement me
	panic("implement me")
}

func (p Publisher) PublishCommentOnIncidentPosted(ctx context.Context, commentID domain.ID) error {
	//TODO implement me
	panic("implement me")
}

func mapToMessage(event command.Event) (*message.Message, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return &message.Message{
		UUID:     uuid.NewString(),
		Metadata: nil,
		Payload:  payload,
	}, nil
}
