package query

import (
	"context"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type MonitorResponseTimeStats struct {
	ID          domain.ID
	RangeInDays int
}

type ResponseTimeStats struct {
	ResponseTimePerRegion []ResponseTimePerRegion
}

type ResponseTimePerRegion struct {
	Region monitor.Region
	Data   []ResponseTimePerDate
}

type ResponseTimePerDate struct {
	Value int16
	Date  time.Time
}

type ResponseTimeStatsOptions struct {
	RangeInDays int
	MonitorID   domain.ID
}

type responseTimeFinder interface {
	ResponseTimeStats(ctx context.Context, options ResponseTimeStatsOptions) (ResponseTimeStats, error)
}

type MonitorResponseTimeStatsHandler struct {
	responseTimeFinder responseTimeFinder
}

func NewMonitorResponseTimeStatsHandler(responseTimeFinder responseTimeFinder) MonitorResponseTimeStatsHandler {
	return MonitorResponseTimeStatsHandler{
		responseTimeFinder: responseTimeFinder,
	}
}

func (h MonitorResponseTimeStatsHandler) Execute(ctx context.Context, q MonitorResponseTimeStats) (ResponseTimeStats, error) {
	return h.responseTimeFinder.ResponseTimeStats(ctx, ResponseTimeStatsOptions{
		RangeInDays: q.RangeInDays,
		MonitorID:   q.ID,
	})
}
