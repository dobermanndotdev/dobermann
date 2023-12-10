package monitor

import "time"

type CheckResult struct {
	statusCode       int
	region           string
	checkedAt        time.Time
	responseTimeInMs time.Duration
}
