package monitor

import (
	"errors"
	"time"
)

type CheckResult struct {
	statusCode       int16
	responseTimeInMs int16
	region           Region
	checkedAt        time.Time
}

func (c *CheckResult) StatusCode() int16 {
	return c.statusCode
}

func (c *CheckResult) ResponseTimeInMs() int16 {
	return c.responseTimeInMs
}

func (c *CheckResult) Region() Region {
	return c.region
}

func (c *CheckResult) CheckedAt() time.Time {
	return c.checkedAt
}

func (c *CheckResult) IsEndpointDown() bool {
	return c.statusCode >= 400
}

func NewCheckResult(
	statusCode int16,
	region Region,
	checkedAt time.Time,
	responseTimeInMs int16,
) (*CheckResult, error) {
	if statusCode < 100 || statusCode > 599 {
		return nil, errors.New("invalid status code")
	}

	if time.Now().Before(checkedAt) {
		return nil, errors.New("checkedAt cannot be set in the past")
	}

	return &CheckResult{
		statusCode:       statusCode,
		region:           region,
		checkedAt:        checkedAt,
		responseTimeInMs: responseTimeInMs,
	}, nil
}
