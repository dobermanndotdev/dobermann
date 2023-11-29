package monitor

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Repository interface {
	Insert(ctx context.Context, monitor *Monitor) error
	FindByID(ctx context.Context, accountID, ID domain.ID) (*Monitor, error)
	Update(ctx context.Context, id domain.ID, fn func(monitor *Monitor) error) error
}
