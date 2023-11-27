package monitors

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func NewPsqlRepository(db boil.ContextExecutor) PsqlRepository {
	return PsqlRepository{
		db: db,
	}
}

type PsqlRepository struct {
	db boil.ContextExecutor
}

func (p PsqlRepository) Insert(ctx context.Context, m *monitor.Monitor) error {
	model := mapMonitorToModel(m)
	err := model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (p PsqlRepository) Update(ctx context.Context, id domain.ID, fn func(monitor *monitor.Monitor) error) error {
	return nil
}
