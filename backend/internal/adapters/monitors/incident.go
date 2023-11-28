package monitors

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type IncidentRepository struct {
	db boil.ContextExecutor
}

func NewIncidentRepository(db boil.ContextExecutor) IncidentRepository {
	return IncidentRepository{db: db}
}

func (i IncidentRepository) Create(ctx context.Context, monitorID domain.ID, incident *monitor.Incident) error {
	model := mapIncidentToModel(incident, monitorID)
	err := model.Insert(ctx, i.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
