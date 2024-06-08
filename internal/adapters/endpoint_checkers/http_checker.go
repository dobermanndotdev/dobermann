package endpoint_checkers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dobermanndotdev/dobermann/internal/common/logs"
	"github.com/dobermanndotdev/dobermann/internal/common/ptr"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
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

const (
	maxRetries            = 3
	delayPerReqFailedInMs = 500
)

type HttpChecker struct {
	timeout time.Duration
	region  monitor.Region
	logger  *logs.Logger
}

func NewHttpChecker(region string, timeoutInSeconds int, logger *logs.Logger) (HttpChecker, error) {
	reg, err := monitor.NewRegion(region)
	if err != nil {
		return HttpChecker{}, err
	}

	if timeoutInSeconds < 0 || timeoutInSeconds > 30 {
		return HttpChecker{}, fmt.Errorf("the timeout should be within the range of 1 and 30")
	}

	return HttpChecker{
		region:  reg,
		logger:  logger,
		timeout: time.Second * time.Duration(timeoutInSeconds),
	}, nil
}

func (h HttpChecker) Check(ctx context.Context, endpointUrl string) (*monitor.CheckResult, error) {
	client := http.Client{
		Timeout: h.timeout / 2,
	}

	counter := 0
	var err error
	var req *http.Request
	var resp *http.Response
	var startedAt time.Time
	var responseTimeInMs int16

	for counter < maxRetries {
		counter++

		req, err = http.NewRequestWithContext(ctx, http.MethodGet, endpointUrl, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("unable to create request: %v", err)
		}
		req.Close = true

		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
		req.Header.Add("User-Agent", randomUserAgent(&userAgents, len(userAgents)-1))
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")

		startedAt = time.Now()
		resp, err = client.Do(req)
		if err != nil {
			h.logger.Errorf("GET request to %s failed due to: %v", endpointUrl, err)
			time.Sleep(time.Millisecond * delayPerReqFailedInMs)
			continue
		}

		if resp != nil {
			_ = resp.Body.Close()
		}

		responseTimeInMs = int16(time.Since(startedAt).Milliseconds())
	}

	var responseStatusCode *int16
	if resp != nil {
		responseStatusCode = ptr.ToPtr(int16(resp.StatusCode))
	}

	return monitor.NewCheckResult(
		domain.NewID(),
		responseStatusCode,
		h.region,
		startedAt,
		responseTimeInMs,
	)
}

func randomUserAgent(userAgents *[]string, max int) string {
	//nolint:gosec
	return (*userAgents)[rand.Intn(max-0)+0]
}
