package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/dobermanndotdev/dobermann/internal/common/observability"
	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type Event interface {
	EventName() string
}

type Header struct {
	Name          string `json:"name"`
	CorrelationID string `json:"correlation_id"`
}

func NewHeader(ctx context.Context, name string) Header {
	correlationID, _ := observability.CorrelationIdFromContext(ctx)

	return Header{
		Name:          name,
		CorrelationID: correlationID,
	}
}

type MonitorCreatedEvent struct {
	Header    Header    `json:"header"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMonitorCreatedEventFromMessage(m *message.Message) (MonitorCreatedEvent, error) {
	var event MonitorCreatedEvent
	err := json.Unmarshal(m.Payload, &event)
	if err != nil {
		return MonitorCreatedEvent{}, err
	}

	return event, nil
}

func (e MonitorCreatedEvent) EventName() string {
	return "MonitorCreatedEvent_v1"
}

type EndpointCheckFailed struct {
	Header        Header    `json:"header"`
	MonitorID     string    `json:"monitor_id"`
	CheckResultID string    `json:"check_result_id"`
	CheckedURL    string    `json:"checked_url"`
	At            time.Time `json:"at"`
}

func (e EndpointCheckFailed) EventName() string {
	return "EndpointCheckFailed_v1"
}

//
// IncidentCreatedEvent
//

type IncidentCreatedEvent struct {
	Header     Header    `json:"header"`
	IncidentID string    `json:"incident_id"`
	MonitorID  string    `json:"monitor_id"`
	At         time.Time `json:"at"`
}

func (e IncidentCreatedEvent) EventName() string {
	return "IncidentCreated_v1"
}

//
// EndpointCheckSucceededEvent
//

type EndpointCheckSucceededEvent struct {
	Header    Header    `json:"header"`
	MonitorID string    `json:"monitor_id"`
	At        time.Time `json:"at"`
}

func (e EndpointCheckSucceededEvent) EventName() string {
	return "EndpointCheckSucceeded_v1"
}

//
// IncidentResolvedEvent
//

type IncidentResolvedEvent struct {
	Header     Header    `json:"header"`
	MonitorID  string    `json:"monitor_id"`
	IncidentID string    `json:"incident_id"`
	At         time.Time `json:"at"`
}

func (e IncidentResolvedEvent) EventName() string {
	return "IncidentResolved_v1"
}

//
// Utils
//

func mapEventToMessage(ctx context.Context, event Event) (*message.Message, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall event: %s due to: %v", event.EventName(), err)
	}

	msg := message.NewMessage(domain.NewID().String(), payload)
	msg.SetContext(ctx)

	return msg, nil
}

func NewEventFromMessage[T any](m *message.Message) (T, error) {
	var event T
	err := json.Unmarshal(m.Payload, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}
