package cron

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"

	"github.com/flowck/doberman/internal/common/logs"
	"github.com/flowck/doberman/internal/common/observability"
)

type Daemon struct {
	task      Task
	maxErrors int16
	logger    *logs.Logger
	interval  time.Duration
	done      chan interface{}
}

type Task interface {
	Run(ctx context.Context) error
	Name() string
}

func NewDaemon(task Task, interval time.Duration, logger *logs.Logger) *Daemon {
	return &Daemon{
		task:      task,
		logger:    logger,
		interval:  interval,
		maxErrors: 3,
		done:      make(chan interface{}),
	}
}

func (c *Daemon) Run(ctx context.Context) error {
	c.logger.WithFields(logs.Fields{"task": c.task.Name()}).Info("Task deamon has started")
	var err error

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			c.logger.Infof("Stopping task %s due to termination", c.task.Name())
			return nil
		case <-ctx.Done():
			c.logger.Infof("Stopping task %s due to context cancelation", c.task.Name())
			return nil
		case <-ticker.C:
			ctxWithCorrId := observability.NewContextWithCorrelationID(ctx)

			err = backoff.Retry(func() error {
				return c.task.Run(ctxWithCorrId)
			}, backoff.WithContext(backoff.NewConstantBackOff(time.Minute*1), ctxWithCorrId))

			if err != nil {
				c.logger.WithError(err).WithFields(logs.Fields{"task": c.task.Name()}).Warnf("Stopping execution due to the number of errors")
				return err
			}
		}
	}
}

func (c *Daemon) Stop() {
	close(c.done)
}
