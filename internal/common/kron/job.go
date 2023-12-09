package kron

import (
	"context"
	"fmt"
	"time"
)

type jobHandler func(ctx context.Context) error

type Job struct {
	MaxErrors int
	Handler   jobHandler
	Interval  time.Duration
}

func NewJob(interval time.Duration, handler jobHandler) *Job {
	return &Job{
		Interval:  interval,
		Handler:   handler,
		MaxErrors: 3,
	}
}

func (j *Job) start(ctx context.Context) error {
	var err error
	errorCount := 0
	ticker := time.NewTicker(j.Interval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err = j.Handler(ctx)
			if err != nil {
				errorCount++
			}

			if errorCount == j.MaxErrors {
				return fmt.Errorf("handler has failed too many times. last error: %v", err)
			}
		}
	}
}
