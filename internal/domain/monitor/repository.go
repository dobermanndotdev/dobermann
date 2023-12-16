package monitor

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Repository interface {
	Delete(ctx context.Context, ID domain.ID) error
	Insert(ctx context.Context, monitor *Monitor) error
	FindByID(ctx context.Context, ID domain.ID) (*Monitor, error)
	Update(ctx context.Context, ID domain.ID, fn func(monitor *Monitor) error) error
	UpdateForCheck(ctx context.Context, fn func(foundMonitors []*Monitor) error) error
	SaveCheckResult(ctx context.Context, ID domain.ID, checkResult *CheckResult) error
}
