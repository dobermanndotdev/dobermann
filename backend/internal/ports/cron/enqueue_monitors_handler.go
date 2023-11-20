package cron

import (
	"context"

	"github.com/flowck/doberman/internal/app"
	"github.com/flowck/doberman/internal/app/command"
)

type enqueueMonitorsTask struct {
	application *app.App
}

func NewEnqueueMonitorsTask(application *app.App) Task {
	return enqueueMonitorsTask{
		application: application,
	}
}

func (e enqueueMonitorsTask) Run(ctx context.Context) error {
	return e.application.EnqueueMonitors.Execute(ctx, command.EnqueueMonitors{})
}

func (e enqueueMonitorsTask) Name() string {
	return "enqueue_monitors_task"
}
