package query

import (
	"context"
	"time"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type MonitorResponseTimeStats struct {
	ID          domain.ID
	RangeInDays *int
}

type ResponseTimeStat struct {
	Value  int
	Region string
	Date   time.Time
}

type ResponseTimeStatsOptions struct {
	RangeInDays int
	MonitorID   domain.ID
}

type responseTimeFinder interface {
	ResponseTimeStats(ctx context.Context, monitorID domain.ID, rangeInDays *int) ([]ResponseTimeStat, error)
}

type MonitorResponseTimeStatsHandler struct {
	responseTimeFinder responseTimeFinder
}

func NewMonitorResponseTimeStatsHandler(responseTimeFinder responseTimeFinder) MonitorResponseTimeStatsHandler {
	return MonitorResponseTimeStatsHandler{
		responseTimeFinder: responseTimeFinder,
	}
}

func (h MonitorResponseTimeStatsHandler) Execute(ctx context.Context, q MonitorResponseTimeStats) ([]ResponseTimeStat, error) {
	return h.responseTimeFinder.ResponseTimeStats(ctx, q.ID, q.RangeInDays)
}
