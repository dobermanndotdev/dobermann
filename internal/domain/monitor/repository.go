package monitor

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Repository interface {
	Insert(ctx context.Context, monitor *Monitor) error
	FindByID(ctx context.Context, ID domain.ID) (*Monitor, error)
	Update(ctx context.Context, ID domain.ID, fn func(monitor *Monitor) error) error
	UpdateForCheck(ctx context.Context, fn func(foundMonitors []*Monitor) error) error
	SaveCheckResult(ctx context.Context, monitorID domain.ID, checkResult *CheckResult) error
}
