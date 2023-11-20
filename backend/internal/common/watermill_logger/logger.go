package watermill_logger

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"
)

type WatermillLogrusLogger struct {
	l *logrus.Logger
}

func NewWatermillLogrusLogger(logger *logrus.Logger) WatermillLogrusLogger {
	return WatermillLogrusLogger{
		l: logger,
	}
}

func (w WatermillLogrusLogger) Error(msg string, err error, fields watermill.LogFields) {
	w.l.WithError(err).WithFields(logrus.Fields(fields)).Error(msg)
}

func (w WatermillLogrusLogger) Info(msg string, fields watermill.LogFields) {
	w.l.WithFields(logrus.Fields(fields)).Info(msg)
}

func (w WatermillLogrusLogger) Debug(msg string, fields watermill.LogFields) {
	w.l.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (w WatermillLogrusLogger) Trace(msg string, fields watermill.LogFields) {
	w.l.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (w WatermillLogrusLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	w.l = w.l.WithFields(logrus.Fields(fields)).Logger
	return w
}
