package monitor

import (
	"errors"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type CheckResult struct {
	id               domain.ID
	statusCode       *int16
	responseTimeInMs int16
	region           Region
	checkedAt        time.Time
}

func NewCheckResult(
	id domain.ID,
	statusCode *int16,
	region Region,
	checkedAt time.Time,
	responseTimeInMs int16,
) (*CheckResult, error) {
	if statusCode != nil {
		if *statusCode < 100 || *statusCode > 599 {
			return nil, errors.New("invalid status code")
		}
	}

	if time.Now().Before(checkedAt) {
		return nil, errors.New("checkedAt cannot be set in the past")
	}

	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	return &CheckResult{
		id:               id,
		statusCode:       statusCode,
		region:           region,
		checkedAt:        checkedAt,
		responseTimeInMs: responseTimeInMs,
	}, nil
}

func (c *CheckResult) ID() domain.ID {
	return c.id
}

func (c *CheckResult) StatusCode() *int16 {
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
	if c.statusCode != nil {
		return *c.statusCode >= 400
	}

	return true
}
