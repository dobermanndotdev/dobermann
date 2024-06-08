package amqp

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/app"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type EndpointCheckFailedHandler struct {
	application *app.App
}

func (e EndpointCheckFailedHandler) HandlerName() string {
	return "EndpointCheckFailed_CreateIncident_Handler"
}

func (e EndpointCheckFailedHandler) EventName() string {
	return events.EndpointCheckFailed{}.EventName()
}

func (e EndpointCheckFailedHandler) Handle(m *message.Message) error {
	event, err := events.NewEventFromMessage[events.EndpointCheckFailed](m)
	if err != nil {
		return fmt.Errorf("unable to map the event '%v': %v", event, err)
	}

	monitorID, err := domain.NewIdFromString(event.MonitorID)
	if err != nil {
		return fmt.Errorf("unable to map the monitor id '%s': %v", event.MonitorID, err)
	}

	checkResultID, err := domain.NewIdFromString(event.CheckResultID)
	if err != nil {
		return fmt.Errorf("unable to map the check result id '%s': %v", event.CheckResultID, err)
	}

	err = e.application.Commands.CreateIncident.Execute(m.Context(), command.CreateIncident{
		MonitorID:     monitorID,
		CheckResultID: checkResultID,
		CheckedUrl:    event.CheckedURL,
	})
	if err != nil {
		return err
	}

	return nil
}
