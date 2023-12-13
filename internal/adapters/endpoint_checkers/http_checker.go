package endpoint_checkers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type HttpChecker struct {
	client *http.Client
	region monitor.Region
}

func NewHttpChecker(region string) (HttpChecker, error) {
	reg, err := monitor.NewRegion(region)
	if err != nil {
		return HttpChecker{}, err
	}

	return HttpChecker{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		region: reg,
	}, nil
}

func (h HttpChecker) Check(ctx context.Context, endpointUrl string) (*monitor.CheckResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	startedAt := time.Now()
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to check endpoint %s: %v", endpointUrl, err)
	}
	defer func() { _ = resp.Body.Close() }()

	checkDuration := time.Since(startedAt)
	checkResult, err := monitor.NewCheckResult(int16(resp.StatusCode), h.region, time.Now(), int16(checkDuration.Milliseconds()))
	if err != nil {
		return nil, fmt.Errorf("unable to create CheckResult: %v", err)
	}

	return checkResult, nil
}
