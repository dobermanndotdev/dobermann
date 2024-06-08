package monitor

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type Repository interface {
	Delete(ctx context.Context, ID domain.ID) error
	Insert(ctx context.Context, monitor *Monitor) error
	FindByID(ctx context.Context, ID domain.ID) (*Monitor, error)
	Update(ctx context.Context, ID domain.ID, fn func(monitor *Monitor) error) error
	UpdateForCheck(ctx context.Context, fn func(foundMonitors []*Monitor) error) error
	SaveCheckResult(ctx context.Context, ID domain.ID, checkResult *CheckResult) error
}

type IncidentRepository interface {
	FindByID(ctx context.Context, id domain.ID) (*Incident, error)
	Create(ctx context.Context, incident *Incident) error
	Update(ctx context.Context, id domain.ID, fn func(incident *Incident) error) error
	AppendIncidentAction(ctx context.Context, incidentID domain.ID, incidentAction *IncidentAction) error
}
