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
}

func NewHttpChecker() HttpChecker {
	return HttpChecker{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (h HttpChecker) Check(ctx context.Context, endpointUrl string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to check endpoint %s: %v", endpointUrl, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return monitor.ErrEndpointIsDown
	}

	return nil
}
