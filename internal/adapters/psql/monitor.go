package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
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

func (p MonitorRepository) UpdateForCheck(ctx context.Context, fn func(foundMonitors []*monitor.Monitor) error) error {
	mods := []qm.QueryMod{
		models.MonitorWhere.IsPaused.EQ(false),
		qm.Where("DATE_PART('seconds', now()::timestamp - last_checked_at::timestamp) >= check_interval_in_seconds"),
		qm.For("UPDATE SKIP LOCKED"),
		qm.Limit(100),
	}

	modelsFound, err := models.Monitors(mods...).All(ctx, p.db)
	if err != nil {
		return err
	}

	monitorsFound, err := mapModelsToMonitors(modelsFound)
	if err != nil {
		return err
	}

	err = fn(monitorsFound)
	if err != nil {
		return err
	}

	modelsFound = mapMonitorsToModels(monitorsFound)

	for _, model := range modelsFound {
		_, err = model.Update(ctx, p.db, boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
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

func (p MonitorRepository) SaveCheckResult(ctx context.Context, monitorID domain.ID, checkResult *monitor.CheckResult) error {
	exists, err := models.Monitors(models.MonitorWhere.ID.EQ(monitorID.String())).Exists(ctx, p.db)
	if err != nil {
		return fmt.Errorf("unable to check if monitor with id %s exists: %v", monitorID, err)
	}

	// monitor could have been deleted, therefore don't do anything
	if !exists {
		return nil
	}

	model := mapCheckResultToModel(monitorID, checkResult)
	err = model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("unable to save check result of monitor with id %s: %v", monitorID, err)
	}

	return nil
}

func (p MonitorRepository) Delete(ctx context.Context, ID domain.ID) error {
	model, err := models.Monitors(
		models.MonitorWhere.ID.EQ(ID.String()),
	).One(ctx, p.db)
	if errors.Is(err, sql.ErrNoRows) {
		return monitor.ErrMonitorNotFound
	}

	if err != nil {
		return fmt.Errorf("error querying the monitor by id %s: %v", ID, err)
	}

	incidents, err := models.Incidents(models.IncidentWhere.MonitorID.EQ(ID.String())).All(ctx, p.db)
	if err != nil {
		return fmt.Errorf("error querying incidents of monitor with id '%s': %v", ID, err)
	}

	for _, incident := range incidents {
		_, err = models.IncidentActions(models.IncidentActionWhere.IncidentID.EQ(incident.ID)).DeleteAll(ctx, p.db)
		if err != nil {
			return fmt.Errorf("error deleting monitor incident actions of monitor with id '%s': %v", ID, err)
		}
	}

	_, err = models.Incidents(models.IncidentWhere.MonitorID.EQ(ID.String())).All(ctx, p.db)
	if err != nil {
		return fmt.Errorf("error deleting monitor incident of monitor with id '%s': %v", ID, err)
	}

	_, err = queries.Raw("DELETE FROM subscribers WHERE monitor_id = $1", ID.String()).ExecContext(ctx, p.db)
	if err != nil {
		return fmt.Errorf("error deleting subscribers of monitor with id '%s': %v", ID, err)
	}

	_, err = model.Delete(ctx, p.db)
	if err != nil {
		return fmt.Errorf("error deleting monitor with id '%s': %v", ID, err)
	}

	return nil
}
