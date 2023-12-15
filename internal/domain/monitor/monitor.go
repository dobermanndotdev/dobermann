package monitor

import (
	"errors"
	"net/url"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
)

const minIntervalInSeconds = 30

type Monitor struct {
	id            domain.ID
	accountID     domain.ID
	endpointUrl   string
	isEndpointUp  bool
	isPaused      bool
	subscribers   []*Subscriber
	incidents     []*Incident
	createdAt     time.Time
	checkInterval time.Duration
	lastCheckedAt *time.Time
	checkResults  []CheckResult
}

func NewMonitor(
	id domain.ID,
	endpointUrl string,
	accountID domain.ID,
	isEndpointUp bool,
	isPaused bool,
	incidents []*Incident,
	subscribers []*Subscriber,
	createdAt time.Time,
	checkInterval time.Duration,
	lastCheckedAt *time.Time,
) (*Monitor, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	if _, err := url.Parse(endpointUrl); err != nil {
		return nil, errors.New("endpointUrl cannot be invalid")
	}

	if accountID.IsEmpty() {
		return nil, errors.New("accountID cannot be invalid")
	}

	if checkInterval < time.Second*minIntervalInSeconds {
		return nil, errors.New("checkInterval cannot be less than 30 seconds")
	}

	return &Monitor{
		id:            id,
		endpointUrl:   endpointUrl,
		accountID:     accountID,
		isEndpointUp:  isEndpointUp,
		isPaused:      isPaused,
		incidents:     incidents,
		subscribers:   subscribers,
		createdAt:     createdAt,
		checkInterval: checkInterval,
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

func (m *Monitor) IsPaused() bool {
	return m.isPaused
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

func (m *Monitor) CheckResults() []CheckResult {
	return m.checkResults
}

func (m *Monitor) SetEndpointCheckResult(isUp bool) {
	m.isEndpointUp = isUp
	lastChecked := time.Now().UTC()
	m.lastCheckedAt = &lastChecked
}

func (m *Monitor) MarkEndpointAsUp() {
	m.isEndpointUp = true
	lastChecked := time.Now().UTC()
	m.lastCheckedAt = &lastChecked
}

func (m *Monitor) MarkEndpointAsDown() {
	m.isEndpointUp = false
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

func (m *Monitor) Pause() {
	m.isPaused = true
}

func (m *Monitor) UnPause() {
	m.isPaused = false
}
