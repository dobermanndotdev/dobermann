package monitors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
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
	model, err := models.FindMonitor(ctx, p.db, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return monitor.ErrMonitorNotFound
	}

	if err != nil {
		return fmt.Errorf("unable to query to find monitor with id %s: %v", id, err)
	}

	m, err := mapModelToMonitor(model)
	if err != nil {
		return err
	}

	err = fn(m)
	if err != nil {
		return err
	}

	model = mapMonitorToModel(m)
	_, err = model.Update(ctx, p.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("unable to update monitor with id %s: %v", id, err)
	}

	return nil
}
