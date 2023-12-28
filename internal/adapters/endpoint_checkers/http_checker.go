package endpoint_checkers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/flowck/dobermann/backend/internal/app/command"
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

func (h HttpChecker) Check(ctx context.Context, endpointUrl string) (command.CheckResult, error) {
	hCtx, cancel := context.WithTimeout(ctx, h.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(hCtx, http.MethodGet, endpointUrl, nil)
	if err != nil {
		return command.CheckResult{}, fmt.Errorf("unable to create request: %v", err)
	}

	req.Header.Add("Accept-Encoding", "identity")

	startedAt := time.Now()
	resp, err := h.client.Do(req)
	if err != nil {
		if errors.Is(hCtx.Err(), context.Canceled) || errors.Is(hCtx.Err(), context.DeadlineExceeded) {
			return h.createCheckResults(startedAt, req, nil, true)
		}

		return command.CheckResult{}, fmt.Errorf("unable to check endpoint %s: %v", endpointUrl, err)
	}
	defer func() { _ = resp.Body.Close() }()

	return h.createCheckResults(startedAt, req, resp, false)
}

func (h HttpChecker) createCheckResults(
	startedAt time.Time,
	req *http.Request,
	resp *http.Response,
	isForcedTimeout bool,
) (command.CheckResult, error) {
	checkDuration := time.Since(startedAt)

	var statusCode int16
	var responseBody string
	var reqHeaders http.Header
	var resHeaders http.Header

	if resp != nil {
		resHeaders = resp.Header
		reqHeaders = req.Header
		statusCode = int16(resp.StatusCode)

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = ""
		} else {
			responseBody = string(data)
		}
	}

	if isForcedTimeout {
		statusCode = int16(http.StatusRequestTimeout)
	}

	result, err := monitor.NewCheckResult(statusCode, h.region, time.Now(), int16(checkDuration.Milliseconds()))
	if err != nil {
		return command.CheckResult{}, fmt.Errorf("unable to create CheckResult: %v", err)
	}

	return command.CheckResult{
		Result:          result,
		ResponseStatus:  statusCode,
		ResponseBody:    responseBody,
		ResponseHeaders: mapHeadersToString(resHeaders),
		RequestHeaders:  mapHeadersToString(reqHeaders),
	}, nil
}

func mapHeadersToString(headers http.Header) string {
	data, err := json.Marshal(headers)
	if err != nil {
		return ""
	}

	return string(data)
}
