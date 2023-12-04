package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func NewMonitorRepository(db boil.ContextExecutor) MonitorRepository {
	return MonitorRepository{
		db: db,
	}
}

type MonitorRepository struct {
	db boil.ContextExecutor
}

func (p MonitorRepository) FindByID(ctx context.Context, id domain.ID) (*monitor.Monitor, error) {
	model, err := models.Monitors(
		models.MonitorWhere.ID.EQ(id.String()),
		qm.Load(models.MonitorRels.Users),
		qm.Load(models.MonitorRels.Incidents),
	).One(ctx, p.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, monitor.ErrMonitorNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("unable to query to find monitor with id %s: %v", id, err)
	}

	return mapModelToMonitor(model)
}

func (p MonitorRepository) Insert(ctx context.Context, m *monitor.Monitor) error {
	model := mapMonitorToModel(m)
	err := model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return err
	}

	if len(m.Subscribers()) == 0 {
		return nil
	}

	var userModel *models.User
	for _, subscriber := range m.Subscribers() {
		userModel, err = models.FindUser(ctx, p.db, subscriber.UserID().String())
		if err != nil {
			return err
		}

		err = model.AddUsers(ctx, p.db, false, userModel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p MonitorRepository) Update(
	ctx context.Context,
	id domain.ID,
	fn func(monitor *monitor.Monitor) error,
) error {
	model, err := models.FindMonitor(ctx, p.db, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return monitor.ErrMonitorNotFound
	}

	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (p MonitorRepository) FindAll(
	ctx context.Context,
	accID domain.ID,
	params query.PaginationParams,
) (query.PaginatedResult[*monitor.Monitor], error) {
	mods := []qm.QueryMod{
		models.MonitorWhere.AccountID.EQ(accID.String()),
		qm.Load(models.MonitorRels.Incidents),
		qm.Offset(mapPaginationParamsToOffset(params.Page, params.Limit)),
		qm.Limit(params.Limit),
		qm.OrderBy("created_at DESC"),
	}

	modelList, err := models.Monitors(mods...).All(ctx, p.db)
	if err != nil {
		return query.PaginatedResult[*monitor.Monitor]{}, fmt.Errorf("error while querying monitors: %v", err)
	}

	count, err := models.Monitors(models.MonitorWhere.AccountID.EQ(accID.String())).Count(ctx, p.db)
	if err != nil {
		return query.PaginatedResult[*monitor.Monitor]{}, fmt.Errorf("error while counting monitors: %v", err)
	}

	monitors, err := mapModelsToMonitors(modelList)
	if err != nil {
		return query.PaginatedResult[*monitor.Monitor]{}, err
	}

	return query.PaginatedResult[*monitor.Monitor]{
		TotalCount: count,
		Data:       monitors,
		Page:       params.Page,
		PerPage:    params.Limit,
		PageCount:  mapPaginationPerPageCount(count, params.Limit),
	}, nil
}
