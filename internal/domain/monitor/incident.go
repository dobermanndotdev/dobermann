package monitor

import (
	"context"
	"errors"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Incident struct {
	id         domain.ID
	isResolved bool
	createdAt  time.Time
	actions    []IncidentAction
}

func NewIncident(id domain.ID, isResolved bool, createdAt time.Time, actions []IncidentAction) (*Incident, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be empty or invalid")
	}

	return &Incident{
		id:         id,
		actions:    actions,
		isResolved: isResolved,
		createdAt:  createdAt.UTC(),
	}, nil
}

func (i *Incident) ID() domain.ID {
	return i.id
}

func (i *Incident) Actions() []IncidentAction {
	return i.actions
}

func (i *Incident) IsResolved() bool {
	return i.isResolved
}

func (i *Incident) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Incident) Resolve() {
	i.isResolved = true
}

//
// Repo
//

type IncidentRepository interface {
	FindByID(ctx context.Context, id domain.ID) (*Incident, error)
	Create(ctx context.Context, monitorID domain.ID, incident *Incident) error
	Update(ctx context.Context, id, monitorID domain.ID, fn func(incident *Incident) error) error
}
