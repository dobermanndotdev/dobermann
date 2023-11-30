package watermill_logger

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"

	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type WatermillLogger struct {
	logger *logs.Logger
}

func (w WatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	w.logger.WithFields(logrus.Fields(fields)).Error(err)
}

func (w WatermillLogger) Info(msg string, fields watermill.LogFields) {
	w.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (w WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	w.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (w WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	w.logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (w WatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return WatermillLogger{logger: w.logger.WithFields(logrus.Fields(fields)).Logger}
}

func NewWatermillLogger(logger *logs.Logger) watermill.LoggerAdapter {
	return WatermillLogger{logger: logger}
}
