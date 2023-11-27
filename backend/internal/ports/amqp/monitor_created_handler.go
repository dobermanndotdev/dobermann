package amqp

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type MonitorCreatedHandler struct {
	application *app.App
}

func (h MonitorCreatedHandler) HandlerName() string {
	return fmt.Sprintf("%s_Handler", events.MonitorCreatedEvent{}.EventName())
}

func (h MonitorCreatedHandler) EventName() string {
	return events.MonitorCreatedEvent{}.EventName()
}

func (h MonitorCreatedHandler) Handle(m *message.Message) error {
	logs.Infof("event consumed: %s", string(m.Payload))
	return nil
}
