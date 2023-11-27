package monitor

import (
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Monitor struct {
	id            domain.ID
	accountID     domain.ID
	endpointUrl   string
	isEndpointUp  bool
	incidents     []Incident
	createdAt     time.Time
	lastCheckedAt *time.Time
}

func NewMonitor(
	id domain.ID,
	endpointUrl string,
	accountID domain.ID,
	isEndpointUp bool,
	incidents []Incident,
	createdAt time.Time,
	lastCheckedAt *time.Time,
) (*Monitor, error) {
	return &Monitor{
		id:            id,
		endpointUrl:   endpointUrl,
		accountID:     accountID,
		isEndpointUp:  isEndpointUp,
		incidents:     incidents,
		createdAt:     createdAt,
		lastCheckedAt: lastCheckedAt,
	}, nil
}

func (m *Monitor) ID() domain.ID {
	return m.id
}

func (m *Monitor) EndpointUrl() string {
	return m.endpointUrl
}

func (m *Monitor) AccountID() domain.ID {
	return m.accountID
}

func (m *Monitor) IsEndpointUp() bool {
	return m.isEndpointUp
}

func (m *Monitor) Incidents() []Incident {
	return m.incidents
}

func (m *Monitor) LastCheckedAt() *time.Time {
	return m.lastCheckedAt
}

func (m *Monitor) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Monitor) SetEndpointCheckResult(isUp bool) {
	m.isEndpointUp = isUp
	lastChecked := time.Now().UTC()
	m.lastCheckedAt = &lastChecked
}
