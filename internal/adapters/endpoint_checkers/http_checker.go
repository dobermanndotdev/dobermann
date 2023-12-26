package endpoint_checkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type HttpChecker struct {
	client *http.Client
	region monitor.Region
}

func NewHttpChecker(region string, timeoutInSeconds int) (HttpChecker, error) {
	reg, err := monitor.NewRegion(region)
	if err != nil {
		return HttpChecker{}, err
	}

	if timeoutInSeconds < 0 || timeoutInSeconds > 30 {
		return HttpChecker{}, fmt.Errorf("the timeout should be within the range of 1 and 30")
	}

	return HttpChecker{
		client: &http.Client{
			Timeout: time.Second * time.Duration(timeoutInSeconds),
		},
		region: reg,
	}, nil
}

func (h HttpChecker) Check(ctx context.Context, endpointUrl string) (*monitor.CheckResult, error) {
	hCtx, cancel := context.WithTimeout(ctx, h.client.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(hCtx, http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	req.Header.Add("Accept-Encoding", "identity")

	startedAt := time.Now()
	resp, err := h.client.Do(req)
	if err != nil {
		if errors.Is(hCtx.Err(), context.Canceled) || errors.Is(hCtx.Err(), context.DeadlineExceeded) {
			return h.createCheckResults(startedAt, nil, true)
		}

		return nil, fmt.Errorf("unable to check endpoint %s: %v", endpointUrl, err)
	}

	defer func() { _ = resp.Body.Close() }()

	return h.createCheckResults(startedAt, &resp.StatusCode, false)
}

func (h HttpChecker) createCheckResults(startedAt time.Time, statusCode *int, isForcedTimeout bool) (*monitor.CheckResult, error) {
	checkDuration := time.Since(startedAt)

	var sCode int16

	if statusCode != nil {
		sCode = int16(*statusCode)
	}

	if isForcedTimeout {
		sCode = int16(http.StatusRequestTimeout)
	}

	checkResult, err := monitor.NewCheckResult(sCode, h.region, time.Now(), int16(checkDuration.Milliseconds()))
	if err != nil {
		return nil, fmt.Errorf("unable to create CheckResult: %v", err)
	}

	return checkResult, nil
}
