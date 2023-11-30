package monitor

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Incident struct {
	id        domain.ID
	createdAt time.Time
	actions   []IncidentAction
}

func (i *Incident) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Incident) ID() domain.ID {
	return i.id
}

func (i *Incident) Actions() []IncidentAction {
	return i.actions
}

func NewIncident(id domain.ID, createdAt time.Time, actions []IncidentAction) (*Incident, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be empty or invalid")
	}

	return &Incident{
		id:        id,
		actions:   actions,
		createdAt: createdAt.UTC(),
	}, nil
}

//
// Repo
//

type IncidentRepository interface {
	FindByID(ctx context.Context, id domain.ID) (*Incident, error)
	Create(ctx context.Context, monitorID domain.ID, incident *Incident) error
}
