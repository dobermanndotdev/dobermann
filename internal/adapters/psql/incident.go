package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/dobermanndotdev/dobermann/internal/adapters/models"
	"github.com/dobermanndotdev/dobermann/internal/app/query"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

func NewIncidentRepository(db boil.ContextExecutor) IncidentRepository {
	return IncidentRepository{db: db}
}

type IncidentRepository struct {
	db boil.ContextExecutor
}

func (i IncidentRepository) FindAll(
	ctx context.Context,
	accountID domain.ID,
	params query.PaginationParams,
) (query.PaginatedResult[*monitor.Incident], error) {
	monitorModels, err := models.Monitors(models.MonitorWhere.AccountID.EQ(accountID.String()), qm.Select("id")).All(ctx, i.db)
	if err != nil {
		return query.PaginatedResult[*monitor.Incident]{}, fmt.Errorf("unable to query monitors of account %s due to: %v", accountID, err)
	}

	monitorIDs := make([]string, len(monitorModels))
	for idx, m := range monitorModels {
		monitorIDs[idx] = m.ID
	}

	mods := []qm.QueryMod{
		models.IncidentWhere.MonitorID.IN(monitorIDs),
		qm.Offset(mapPaginationParamsToOffset(params.Page, params.Limit)),
		qm.Limit(params.Limit),
		qm.OrderBy("created_at DESC"),
	}

	modelList, err := models.Incidents(mods...).All(ctx, i.db)
	if err != nil {
		return query.PaginatedResult[*monitor.Incident]{}, fmt.Errorf("error while querying monitors: %v", err)
	}

	count, err := models.Incidents(models.IncidentWhere.MonitorID.IN(monitorIDs)).Count(ctx, i.db)
	if err != nil {
		return query.PaginatedResult[*monitor.Incident]{}, fmt.Errorf("error while counting monitors: %v", err)
	}

	incidents, err := mapModelsToIncidents(modelList)
	if err != nil {
		return query.PaginatedResult[*monitor.Incident]{}, err
	}

	return query.PaginatedResult[*monitor.Incident]{
		TotalCount: count,
		Data:       incidents,
		Page:       params.Page,
		PerPage:    params.Limit,
		PageCount:  mapPaginationPerPageCount(count, params.Limit),
	}, nil
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

func (i IncidentRepository) Create(ctx context.Context, incident *monitor.Incident) error {
	model := mapIncidentToModel(incident)
	err := model.Insert(ctx, i.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (i IncidentRepository) Update(
	ctx context.Context,
	id domain.ID,
	fn func(incident *monitor.Incident) error,
) error {
	incident, err := i.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = fn(incident)
	if err != nil {
		return err
	}

	model := mapIncidentToModel(incident)

	_, err = model.Update(ctx, i.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (i IncidentRepository) AppendIncidentAction(
	ctx context.Context,
	incidentID domain.ID,
	action *monitor.IncidentAction,
) error {
	model := mapIncidentActionToModel(action, incidentID)
	err := model.Insert(ctx, i.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
