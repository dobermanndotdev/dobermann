package kron

import (
	"context"
	"fmt"
	"time"
)

type JobHandler func(ctx context.Context) error
type MiddlewareFunc func(ctx context.Context) (context.Context, error)

type Job struct {
	MaxErrors   int
	Handler     JobHandler
	Interval    time.Duration
	middlewares []MiddlewareFunc
}

func NewJob(interval time.Duration, handler JobHandler) *Job {
	return &Job{
		MaxErrors: 3,
		Handler:   handler,
		Interval:  interval,
	}
}

func (j *Job) start(ctx context.Context) error {
	var err error
	errorCount := 0
	ticker := time.NewTicker(j.Interval)

	jobCtx := ctx

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if len(j.middlewares) > 0 {
				for _, middleware := range j.middlewares {
					jobCtx, err = middleware(jobCtx)
					if err != nil {
						errorCount++
						continue
					}
				}
			}

			err = j.Handler(jobCtx)
			if err != nil {
				errorCount++
			}

			if errorCount == j.MaxErrors {
				return fmt.Errorf("handler has failed too many times. last error: %v", err)
			}
		}
	}
}

func (j *Job) AddMiddleware(middleware MiddlewareFunc) {
	j.middlewares = append(j.middlewares, middleware)
}
