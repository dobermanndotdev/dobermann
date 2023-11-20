package monitors

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/flowck/doberman/internal/adapters/models"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/monitor"
)

type PsqlRepository struct {
	db boil.ContextExecutor
}

func NewPsqlRepository(db boil.ContextExecutor) PsqlRepository {
	return PsqlRepository{
		db: db,
	}
}

func (p PsqlRepository) FindByID(ctx context.Context, ID domain.ID) (*monitor.Monitor, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlRepository) UpdateAllForCheck(ctx context.Context, updater func(mos []*monitor.Monitor) error) error {
	mods := []qm.QueryMod{
		models.MonitorWhere.CheckStatus.NEQ(monitor.CheckStatusEnqueued.String()),
		models.MonitorWhere.IsPaused.EQ(false),
		qm.Where("DATE_PART('minute', now()::timestamp - last_checked_at::timestamp) >= check_interval"),
		qm.Limit(100),
		qm.For("UPDATE"),
	}

	mdls, err := models.Monitors(mods...).All(ctx, p.db)
	if err != nil {
		return err
	}

	monitorList, err := mapFromModelsToMonitors(mdls)
	if err != nil {
		return err
	}

	if err = updater(monitorList); err != nil {
		return err
	}

	mdls = mapFromMonitorsToModels(monitorList)
	for _, model := range mdls {
		_, err = model.Update(ctx, p.db, boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p PsqlRepository) Insert(ctx context.Context, mo *monitor.Monitor) error {
	model := mapFromMonitorToModel(mo)
	if err := model.Insert(ctx, p.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func (p PsqlRepository) Update(ctx context.Context, id domain.ID, updater func(m *monitor.Monitor) error) error {
	//TODO implement me
	panic("implement me")
}
