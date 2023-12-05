package messaging

import (
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/common/watermill_logger"
)

type AmqpPublisher struct {
	publisher message.Publisher
}

func NewAmqpPublisher(amqpUrl string, logger *logs.Logger) (AmqpPublisher, error) {
	amqpPublisher, err := amqp.NewPublisher(amqp.NewDurableQueueConfig(amqpUrl), watermill_logger.NewWatermillLogger(logger))
	if err != nil {
		return AmqpPublisher{}, err
	}

	return AmqpPublisher{
		publisher: NewCorrelationIdPublisher(amqpPublisher),
	}, nil
}

func (p AmqpPublisher) Publish(topic string, messages ...*message.Message) error {
	return p.publisher.Publish(topic, messages...)
}

func (p AmqpPublisher) Close() error {
	return p.publisher.Close()
}
