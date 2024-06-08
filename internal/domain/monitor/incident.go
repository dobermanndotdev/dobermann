package monitor

import (
	"errors"
	"time"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type Incident struct {
	id                 domain.ID
	monitorID          domain.ID
	createdAt          time.Time
	resolvedAt         *time.Time
	checkedURL         URL
	cause              string
	responseStatusCode *int16
	actions            []IncidentAction
}

func NewIncident(
	id,
	monitorID domain.ID,
	resolvedAt *time.Time,
	createdAt time.Time,
	checkedURL string,
	actions []IncidentAction,
	cause string,
	responseStatusCode *int16,
) (*Incident, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be empty or invalid")
	}

	parsedCheckedURL, err := NewURL(checkedURL)
	if err != nil {
		return nil, err
	}

	var resolvedAtUTC *time.Time
	if resolvedAt != nil {
		if resolvedAt.UTC().After(time.Now().UTC()) {
			return nil, errors.New("resolvedAt cannot be set in the future")
		}

		r := resolvedAt.UTC()
		resolvedAtUTC = &r
	}

	if responseStatusCode != nil {
		if *responseStatusCode < 100 || *responseStatusCode > 599 {
			return nil, errors.New("the status provided is invalid")
		}
	}

	return &Incident{
		id:                 id,
		cause:              cause,
		actions:            actions,
		monitorID:          monitorID,
		checkedURL:         parsedCheckedURL,
		resolvedAt:         resolvedAtUTC,
		createdAt:          createdAt.UTC(),
		responseStatusCode: responseStatusCode,
	}, nil
}

func (i *Incident) Cause() string {
	return i.cause
}

func (i *Incident) ResponseStatusCode() *int16 {
	return i.responseStatusCode
}

func (i *Incident) MonitorID() domain.ID {
	return i.monitorID
}

func (i *Incident) CheckedURL() string {
	return i.checkedURL.String()
}

func (i *Incident) ID() domain.ID {
	return i.id
}

func (i *Incident) Actions() []IncidentAction {
	return i.actions
}

func (i *Incident) IsResolved() bool {
	return i.resolvedAt != nil
}

func (i *Incident) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Incident) Resolve() {
	now := time.Now()
	i.resolvedAt = &now
}

func (i *Incident) ResolvedAt() *time.Time {
	return i.resolvedAt
}
