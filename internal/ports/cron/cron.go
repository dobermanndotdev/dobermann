package cron

import (
	"time"

	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/common/kron"
)

type Handlers struct {
	application *app.App
}

func NewService(application *app.App) *kron.Service {
	c := kron.NewService()
	handlers := Handlers{application: application}

	c.AddJob(kron.NewJob(time.Second*5, handlers.BulkCheckEndpoints))

	return c
}
