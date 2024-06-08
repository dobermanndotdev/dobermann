package observability

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/common/logs"
)

type commandLoggingDecorator[C any] struct {
	logger *logs.Logger
	base   CommandHandler[C]
}

func (c commandLoggingDecorator[C]) Execute(ctx context.Context, cmd C) (err error) {
	logger := c.logger.WithFields(logs.Fields{
		"name":    handlerName(c.base),
		"command": prettyPrint(cmd),
	})

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.WithError(err).Error("Failed to execute query")
		}
	}()

	return c.base.Execute(ctx, cmd)
}

type commandWithResultLoggingDecorator[Q any, R any] struct {
	logger *logs.Logger
	base   QueryHandler[Q, R]
}

func (c commandWithResultLoggingDecorator[C, R]) Execute(ctx context.Context, cmd C) (result R, err error) {
	logger := c.logger.WithFields(logs.Fields{
		"name":                handlerName(c.base),
		"command with result": prettyPrint(cmd),
	})

	logger.Debug("Executing command with result")
	defer func() {
		if err == nil {
			logger.Info("Command with result executed successfully")
		} else {
			logger.WithError(err).Error("Failed to execute command with result")
		}
	}()

	return c.base.Execute(ctx, cmd)
}

type queryLoggingDecorator[Q any, R any] struct {
	logger *logs.Logger
	base   QueryHandler[Q, R]
}

func (c queryLoggingDecorator[Q, R]) Execute(ctx context.Context, q Q) (result R, err error) {
	logger := c.logger.WithFields(logs.Fields{
		"name":  handlerName(c.base),
		"query": prettyPrint(q),
	})

	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query executed successfully")
		} else {
			logger.WithError(err).Error("Failed to execute query")
		}
	}()

	return c.base.Execute(ctx, q)
}
