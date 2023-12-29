package monitor

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Incident struct {
	id         domain.ID
	monitorID  domain.ID
	createdAt  time.Time
	resolvedAt *time.Time
	checkedURL string
	details    IncidentDetails
	actions    []IncidentAction
}

func (i *Incident) MonitorID() domain.ID {
	return i.monitorID
}

func (i *Incident) CheckedURL() string {
	return i.checkedURL
}

func (i *Incident) Details() IncidentDetails {
	return i.details
}

func NewIncident(
	id,
	monitorID domain.ID,
	resolvedAt *time.Time,
	createdAt time.Time,
	checkedURL string,
	actions []IncidentAction,
	details IncidentDetails,
) (*Incident, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be empty or invalid")
	}

	var resolvedAtUTC *time.Time
	if resolvedAt != nil {
		if resolvedAt.UTC().After(time.Now().UTC()) {
			return nil, errors.New("resolvedAt cannot be set in the future")
		}

		r := resolvedAt.UTC()
		resolvedAtUTC = &r
	}

	if details.Status < 100 || details.Status > 599 {
		return nil, errors.New("the status provided is invalid")
	}

	checkedURL = strings.TrimSpace(checkedURL)
	if checkedURL == "" {
		return nil, errors.New("checkedURL cannot be empty")
	}

	if _, err := url.Parse(checkedURL); err != nil {
		return nil, fmt.Errorf("the url provided is invalid")
	}

	return &Incident{
		id:         id,
		monitorID:  monitorID,
		actions:    actions,
		details:    details,
		checkedURL: checkedURL,
		resolvedAt: resolvedAtUTC,
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

type IncidentDetails struct {
	Cause           string
	Status          int16
	ResponseBody    string
	ResponseHeaders string
	RequestHeaders  string
}

//
// Repo
//

type IncidentRepository interface {
	FindByID(ctx context.Context, id domain.ID) (*Incident, error)
	Create(ctx context.Context, incident *Incident) error
	Update(ctx context.Context, id domain.ID, fn func(incident *Incident) error) error
	AppendIncidentAction(ctx context.Context, incidentID domain.ID, incidentAction *IncidentAction) error
}
