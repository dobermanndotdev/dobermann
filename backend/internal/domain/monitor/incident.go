package monitor

import (
	"time"

	"github.com/friendsofgo/errors"
)

type Incident struct {
	id      string
	actions []IncidentAction
}

func (i *Incident) ID() string {
	return i.id
}

func (i *Incident) Actions() []IncidentAction {
	return i.actions
}

func NewIncident(id string, actions []IncidentAction) (*Incident, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty or invalid")
	}

	return &Incident{
		id:      id,
		actions: actions,
	}, nil
}

type IncidentAction struct {
	takerUserID string
	takenAt     time.Time
	actionType  IncidentActionType
}

func (i IncidentAction) TakerUserID() string {
	return i.takerUserID
}

func (i IncidentAction) TakenAt() time.Time {
	return i.takenAt
}

func (i IncidentAction) ActionType() IncidentActionType {
	return i.actionType
}

var (
	IncidentActionTypeResolved     = "resolved"
	IncidentActionTypeAcknowledged = "acknowledged"
)

type IncidentActionType struct {
	value string
}

func (t IncidentActionType) String() string {
	return t.value
}
