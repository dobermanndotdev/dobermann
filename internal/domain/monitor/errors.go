package monitor

import "errors"

var (
	ErrEndpointIsDown   = errors.New("endpoint is down")
	ErrMonitorNotFound  = errors.New("monitor not found")
	ErrIncidentNotFound = errors.New("incident not found")
)
