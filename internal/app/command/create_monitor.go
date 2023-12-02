package command

import (
	"context"
	"fmt"

	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CreateMonitor struct {
	Monitor *monitor.Monitor
}

type CreateMonitorHandler struct {
	eventPublisher    EventPublisher
	monitorRepository monitor.Repository
}

func NewCreateMonitorHandler(monitorRepository monitor.Repository, eventPublisher EventPublisher) CreateMonitorHandler {
	return CreateMonitorHandler{
		eventPublisher:    eventPublisher,
		monitorRepository: monitorRepository,
	}
}

func (h CreateMonitorHandler) Execute(ctx context.Context, cmd CreateMonitor) error {
	err := h.monitorRepository.Insert(ctx, cmd.Monitor)
	if err != nil {
		return err
	}

	err = h.eventPublisher.PublishMonitorCreated(ctx, MonitorCreatedEvent{
		ID:        cmd.Monitor.ID().String(),
		CreatedAt: cmd.Monitor.CreatedAt(),
	})
	if err != nil {
		return fmt.Errorf("unable to publish event PublishMonitorCreated: %v", err)
	}

	return nil
}
