package monitor

import (
	"context"

	"github.com/flowck/doberman/internal/domain"
)

type Repository interface {
	FindByID(ctx context.Context, ID domain.ID) (*Monitor, error)
	Insert(ctx context.Context, mo *Monitor) error
	Update(ctx context.Context, ID domain.ID, updater func(mo *Monitor) error) error
	UpdateAllForCheck(ctx context.Context, updater func(mos []*Monitor) error) error
}
