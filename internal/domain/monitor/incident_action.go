package monitor

import (
	"errors"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

var (
	IncidentActionTypeCreated      = IncidentActionType{"created"}
	IncidentActionTypeResolved     = IncidentActionType{"resolved"}
	IncidentActionTypeAcknowledged = IncidentActionType{"acknowledged"}
)

type IncidentActionType struct {
	value string
}

func (t IncidentActionType) String() string {
	return t.value
}

type IncidentAction struct {
	id                domain.ID
	takenByUserWithID *domain.ID
	at                time.Time
	description       string
	actionType        IncidentActionType
}

func NewIncidentAction(
	id domain.ID,
	takenByUserWithID *domain.ID,
	at time.Time,
	description string,
	actionType IncidentActionType,
) (*IncidentAction, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	if at.IsZero() {
		return nil, errors.New("at cannot be invalid")
	}

	at = at.UTC()
	if at.After(time.Now().UTC()) {
		return nil, errors.New("at cannot be set in the future")
	}

	return &IncidentAction{
		id:                id,
		at:                at,
		actionType:        actionType,
		description:       description,
		takenByUserWithID: takenByUserWithID,
	}, nil
}

func (i *IncidentAction) Id() domain.ID {
	return i.id
}

func (i *IncidentAction) Description() string {
	return i.description
}

func (i *IncidentAction) TakerUserID() *domain.ID {
	return i.takenByUserWithID
}

func (i *IncidentAction) TakenAt() time.Time {
	return i.at
}

func (i *IncidentAction) ActionType() IncidentActionType {
	return i.actionType
}
