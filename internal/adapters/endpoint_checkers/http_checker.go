package endpoint_checkers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v3"

	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2.1 Safari/605.1.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.1",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 OPR/95.0.0.",
	"Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.3",
}

type HttpChecker struct {
	timeout time.Duration
	region  monitor.Region
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
		timeout: time.Second * time.Duration(timeoutInSeconds),
		region:  reg,
	}, nil
}

func (h HttpChecker) Check(_ context.Context, endpointUrl string) (command.CheckResult, error) {
	requestCtx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = h.timeout
	exponentialBackoff.InitialInterval = time.Millisecond * 250
	backoffStrategy := backoff.WithContext(exponentialBackoff, requestCtx)

	var result command.CheckResult

	err := backoff.Retry(func() error {
		result = command.CheckResult{}

		client := http.Client{
			Timeout: time.Second * 5,
		}

		req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, endpointUrl, nil)
		if err != nil {
			return fmt.Errorf("unable to create request: %v", err)
		}

		req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Add("User-Agent", randomUserAgent(&userAgents, len(userAgents)-1))

		startedAt := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			// fmt.Printf("context canceled? --> %v // context deadline exceeded? --> %v\n", errors.Is(requestCtx.Err(), context.Canceled), errors.Is(requestCtx.Err(), context.DeadlineExceeded))

			if errors.Is(requestCtx.Err(), context.DeadlineExceeded) {
				result, err = h.createCheckResults(startedAt, req, nil, true)
				if err != nil {
					return err
				}

				return nil
			}

			// fmt.Println("error --> ", err)
			return fmt.Errorf("unable to check endpoint %s: %v", endpointUrl, err)
		}
		defer func() { _ = resp.Body.Close() }()

		result, err = h.createCheckResults(startedAt, req, resp, false)
		if err != nil {
			return err
		}

		return nil
	}, backoffStrategy)

	return result, err
}

func (h HttpChecker) createCheckResults(
	startedAt time.Time,
	req *http.Request,
	resp *http.Response,
	isForcedTimeout bool,
) (command.CheckResult, error) {
	checkDuration := time.Since(startedAt)

	var statusCode int16
	var reqHeaders http.Header
	var resHeaders http.Header

	if resp != nil {
		resHeaders = resp.Header
		reqHeaders = req.Header
		statusCode = int16(resp.StatusCode)
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
		ResponseBody:    "",
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

func randomUserAgent(userAgents *[]string, max int) string {
	//nolint:gosec
	return (*userAgents)[rand.Intn(max-0)+0]
}
