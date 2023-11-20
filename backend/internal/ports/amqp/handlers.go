package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/doberman/internal/app"
)

type Handler interface {
	Name() string
	EventName() string
	Handle(msg *message.Message) error
}

func NewHandlers(application *app.App) []Handler {
	return []Handler{
		checkMonitorEndpointHandler{application: application},
	}
}
