package command

import (
	"context"

	"github.com/flowck/doberman/internal/domain/monitor"
)

type CreateMonitor struct {
	Monitor *monitor.Monitor
}

type createMonitorHandler struct {
	repo monitor.Repository
}

func NewCreateMonitorHandler(repo monitor.Repository) createMonitorHandler {
	return createMonitorHandler{
		repo: repo,
	}
}

func (h createMonitorHandler) Execute(ctx context.Context, cmd CreateMonitor) error {
	err := h.repo.Insert(ctx, cmd.Monitor)
	if err != nil {
		return err
	}

	return nil
}
