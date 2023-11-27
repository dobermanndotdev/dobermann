package monitor

import "time"

type Monitor struct {
	id           string
	endpointUrl  string
	accountID    string
	isEndpointUp bool
	createdAt    time.Time
	incidents    []Incident
}

func NewMonitor(
	id string,
	endpointUrl string,
	accountID string,
	isEndpointUp bool,
	incidents []Incident,
	createdAt time.Time,
) (*Monitor, error) {
	return &Monitor{
		id:           id,
		endpointUrl:  endpointUrl,
		accountID:    accountID,
		isEndpointUp: isEndpointUp,
		incidents:    incidents,
		createdAt:    createdAt,
	}, nil
}

func (m *Monitor) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Monitor) ID() string {
	return m.id
}

func (m *Monitor) EndpointUrl() string {
	return m.endpointUrl
}

func (m *Monitor) AccountID() string {
	return m.accountID
}

func (m *Monitor) IsEndpointUp() bool {
	return m.isEndpointUp
}

func (m *Monitor) Incidents() []Incident {
	return m.incidents
}
