package command

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CheckEndpoint struct {
	MonitorID domain.ID
}

type CheckEndpointHandler struct {
	httpChecker       httpChecker
	eventPublisher    EventPublisher
	monitorRepository monitor.Repository
}

type httpChecker interface {
	Check(ctx context.Context, endpointUrl string) (CheckResult, error)
}

type CheckResult struct {
	Result          *monitor.CheckResult
	ResponseHeaders string
	RequestHeaders  string
	ResponseStatus  int16
	ResponseBody    string
}

func NewCheckEndpointHandler(
	httpChecker httpChecker,
	monitorRepository monitor.Repository,
	eventPublisher EventPublisher,
) CheckEndpointHandler {
	return CheckEndpointHandler{
		httpChecker:       httpChecker,
		eventPublisher:    eventPublisher,
		monitorRepository: monitorRepository,
	}
}

func (c CheckEndpointHandler) Execute(ctx context.Context, cmd CheckEndpoint) error {
	checkSucceeded := false
	var checkResult CheckResult
	var foundMonitor *monitor.Monitor

	err := c.monitorRepository.Update(ctx, cmd.MonitorID, func(m *monitor.Monitor) error {
		foundMonitor = m

		var err error
		checkResult, err = c.httpChecker.Check(ctx, m.EndpointUrl())
		if err != nil {
			return fmt.Errorf("error checking endpoint %s due to: %v", m.EndpointUrl(), err)
		}

		if checkResult.Result.IsEndpointDown() {
			m.MarkEndpointAsDown()
		} else {
			m.MarkEndpointAsUp()
			checkSucceeded = true
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating monitor during check: %v", err)
	}

	err = c.monitorRepository.SaveCheckResult(ctx, cmd.MonitorID, checkResult.Result)
	if err != nil {
		return fmt.Errorf("unable to save the check result of monitor with id %s: %v", cmd.MonitorID, err)
	}

	if checkSucceeded {
		err = c.eventPublisher.PublishEndpointCheckSucceeded(ctx, EndpointCheckSucceededEvent{
			At:        time.Now(),
			MonitorID: cmd.MonitorID.String(),
		})

		if err != nil {
			return fmt.Errorf("error publishing event EndpointCheckSucceededEvent: %v", err)
		}
	} else {
		err = c.eventPublisher.PublishEndpointCheckFailed(ctx, EndpointCheckFailedEvent{
			At:              time.Now(),
			MonitorID:       cmd.MonitorID.String(),
			CheckedURL:      foundMonitor.EndpointUrl(),
			ResponseHeaders: checkResult.ResponseHeaders,
			ResponseBody:    checkResult.ResponseBody,
			ResponseStatus:  checkResult.ResponseStatus,
			RequestHeaders:  checkResult.RequestHeaders,
			Cause:           fmt.Sprintf("Request failed with status code %d", checkResult.ResponseStatus),
		})
		if err != nil {
			return fmt.Errorf("error publishing event EndpointCheckFailedEvent: %v", err)
		}
	}

	return nil
}
