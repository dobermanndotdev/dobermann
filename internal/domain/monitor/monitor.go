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
	subscribers   []*Subscriber
	incidents     []*Incident
	createdAt     time.Time
	checkInterval time.Duration
	lastCheckedAt *time.Time
}

func NewMonitor(
	id domain.ID,
	endpointUrl string,
	accountID domain.ID,
	isEndpointUp bool,
	incidents []*Incident,
	subscribers []*Subscriber,
	createdAt time.Time,
	lastCheckedAt *time.Time,
) (*Monitor, error) {
	return &Monitor{
		id:            id,
		endpointUrl:   endpointUrl,
		accountID:     accountID,
		isEndpointUp:  isEndpointUp,
		incidents:     incidents,
		subscribers:   subscribers,
		createdAt:     createdAt,
		checkInterval: time.Second * 30,
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

func (m *Monitor) Incidents() []*Incident {
	return m.incidents
}

func (m *Monitor) LastCheckedAt() *time.Time {
	return m.lastCheckedAt
}

func (m *Monitor) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Monitor) CheckInterval() time.Duration {
	return m.checkInterval
}

func (m *Monitor) SetEndpointCheckResult(isUp bool) {
	m.isEndpointUp = isUp
	lastChecked := time.Now().UTC()
	m.lastCheckedAt = &lastChecked
}

func (m *Monitor) Subscribers() []*Subscriber {
	return m.subscribers
}

func (m *Monitor) HasIncidentUnresolved() bool {
	for _, incident := range m.incidents {
		if !incident.IsResolved() {
			return true
		}
	}

	return false
}

func (m *Monitor) IncidentUnresolved() *Incident {
	for _, incident := range m.incidents {
		if !incident.IsResolved() {
			return incident
		}
	}

	return nil
}

func (m *Monitor) IsValid() bool {
	return !m.ID().IsEmpty()
}
