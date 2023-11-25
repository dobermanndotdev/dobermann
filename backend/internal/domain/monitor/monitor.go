package monitor

type Monitor struct {
	id           string
	endpointUrl  string
	accountID    string
	isEndpointUp bool
	incidents    []Incident
}

func NewMonitor(
	id string,
	endpointUrl string,
	accountID string,
	isEndpointUp bool,
	incidents []Incident,
) (*Monitor, error) {
	return &Monitor{
		id:           id,
		endpointUrl:  endpointUrl,
		accountID:    accountID,
		isEndpointUp: isEndpointUp,
		incidents:    incidents,
	}, nil
}
