package events

import (
	"context"
	"sync"

	"github.com/flowck/dobermann/backend/internal/app/command"
)

type PublisherMock struct {
	mutex  *sync.RWMutex
	events map[string]Event
}

func NewPublisherMock() *PublisherMock {
	return &PublisherMock{
		mutex:  &sync.RWMutex{},
		events: make(map[string]Event),
	}
}

func (p PublisherMock) PublishMonitorCreated(ctx context.Context, event command.MonitorCreatedEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.events[event.ID] = MonitorCreatedEvent{
		Header:    Header{},
		ID:        event.ID,
		CreatedAt: event.CreatedAt,
	}

	return nil
}

func (p PublisherMock) PublishEndpointCheckFailed(ctx context.Context, event command.EndpointCheckFailedEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return nil
}

func (p PublisherMock) PublishIncidentCreated(ctx context.Context, event command.IncidentCreatedEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return nil
}

func (p PublisherMock) PublishEndpointCheckSucceeded(ctx context.Context, event command.EndpointCheckSucceededEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return nil
}

func (p PublisherMock) PublishIncidentResolved(ctx context.Context, event command.IncidentResolvedEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return nil
}

func (p PublisherMock) GetEventByID(id string) Event {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	event, exists := p.events[id]
	if !exists {
		return nil
	}

	return event
}
