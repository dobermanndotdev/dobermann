package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func NewIncidentRepository(db boil.ContextExecutor) IncidentRepository {
	return IncidentRepository{db: db}
}

type IncidentRepository struct {
	db boil.ContextExecutor
}

func (i IncidentRepository) FindByID(ctx context.Context, id domain.ID) (*monitor.Incident, error) {
	model, err := models.FindIncident(ctx, i.db, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, monitor.ErrIncidentNotFound
	}

	if err != nil {
		return nil, err
	}

	return mapModelToIncident(model)
}

func (i IncidentRepository) Create(ctx context.Context, monitorID domain.ID, incident *monitor.Incident) error {
	model := mapIncidentToModel(incident, monitorID)
	err := model.Insert(ctx, i.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (i IncidentRepository) Update(ctx context.Context, id, monitorID domain.ID, fn func(incident *monitor.Incident) error) error {
	incident, err := i.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = fn(incident)
	if err != nil {
		return err
	}

	model := mapIncidentToModel(incident, monitorID)

	_, err = model.Update(ctx, i.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
