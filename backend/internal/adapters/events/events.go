package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Event interface {
	EventName() string
}

type Header struct {
	Name          string `json:"name"`
	CorrelationID string `json:"correlation_id"`
}

func NewHeader(name, correlationID string) Header {
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
	Header    Header    `json:"header"`
	MonitorID string    `json:"monitor_id"`
	At        time.Time `json:"at"`
}

func (e EndpointCheckFailed) EventName() string {
	return "EndpointCheckFailed_v1"
}

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
// Utils
//

func mapEventToMessage(event Event) (*message.Message, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall event: %s due to: %v", event.EventName(), err)
	}

	return &message.Message{
		Payload: data,
	}, nil
}

func NewEventFromMessage[T any](m *message.Message) (T, error) {
	var event T
	err := json.Unmarshal(m.Payload, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}
