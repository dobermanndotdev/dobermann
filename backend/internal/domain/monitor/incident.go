package monitor

import (
	"time"

	"github.com/friendsofgo/errors"
)

type Incident struct {
	id      string
	actions []IncidentAction
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

var (
	IncidentActionTypeResolved     = "resolved"
	IncidentActionTypeAcknowledged = "acknowledged"
)

type IncidentActionType struct {
	value string
}
