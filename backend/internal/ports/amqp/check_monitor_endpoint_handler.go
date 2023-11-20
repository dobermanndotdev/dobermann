package amqp

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/doberman/internal/app"
	"github.com/flowck/doberman/internal/app/command"
	"github.com/flowck/doberman/internal/domain"
)

type checkMonitorEndpointHandler struct {
	application *app.App
}

func (c checkMonitorEndpointHandler) EventName() string {
	return ""
}

func (c checkMonitorEndpointHandler) Name() string {
	return ""
}

func (c checkMonitorEndpointHandler) Handle(msg *message.Message) error {
	event, err := mapPayloadToMonitorsEnqueuedEvent(msg.Payload)
	if err != nil {
		return err
	}

	var mID domain.ID
	for _, id := range event.IDs {
		if mID, err = domain.NewIdFromString(id); err != nil {
			return err
		}

		err = c.application.CheckMonitorEndpoint.Execute(msg.Context(), command.CheckMonitorEndpoint{
			MonitorID: mID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func mapPayloadToMonitorsEnqueuedEvent(p message.Payload) (command.MonitorsEnqueuedEvent, error) {
	var event command.MonitorsEnqueuedEvent
	if err := json.Unmarshal(p, &event); err != nil {
		return command.MonitorsEnqueuedEvent{}, err
	}

	return event, nil
}
