package monitor

import (
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

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

type IncidentAction struct {
	takerUserID domain.ID
	takenAt     time.Time
	actionType  IncidentActionType
}

func (i IncidentAction) TakerUserID() domain.ID {
	return i.takerUserID
}

func (i IncidentAction) TakenAt() time.Time {
	return i.takenAt
}

func (i IncidentAction) ActionType() IncidentActionType {
	return i.actionType
}
