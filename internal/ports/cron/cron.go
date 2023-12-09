package cron

import (
	"time"

	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/common/kron"
)

type handlers struct {
	region      string
	application *app.App
}

func NewService(application *app.App, region string) *kron.Service {
	c := kron.NewService()
	allHandlers := handlers{
		region:      region,
		application: application,
	}

	c.AddJob(kron.NewJob(time.Second*5, allHandlers.BulkCheckEndpoints))

	return c
}
