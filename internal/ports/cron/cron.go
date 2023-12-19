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

func NewService(application *app.App, region string, isProduction bool) *kron.Service {
	c := kron.NewService()
	allHandlers := handlers{
		region:      region,
		application: application,
	}

	interval := time.Second * 5
	if isProduction {
		interval = time.Second * 25
	}

	bulkCheckEndpointsJob := kron.NewJob(interval, allHandlers.BulkCheckEndpoints)
	bulkCheckEndpointsJob.AddMiddleware(withCorrelationIdMiddleware)

	c.AddJob(bulkCheckEndpointsJob)

	return c
}
