package http

import (
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func mapMonitorToResponseItem(m *monitor.Monitor) Monitor {
	return Monitor{
		CreatedAt:     m.CreatedAt(),
		EndpointUrl:   m.EndpointUrl(),
		Id:            m.ID().String(),
		Incidents:     mapIncidentsToResponse(m.Incidents()),
		IsEndpointUp:  m.IsEndpointUp(),
		IsPaused:      m.IsPaused(),
		LastCheckedAt: m.LastCheckedAt(),
	}
}

func mapMonitorsToResponseItems(monitors []*monitor.Monitor) []Monitor {
	result := make([]Monitor, len(monitors))

	for i, m := range monitors {
		result[i] = mapMonitorToResponseItem(m)
	}

	return result
}

func mapIncidentsToResponse(incidents []*monitor.Incident) []Incident {
	if incidents == nil {
		return make([]Incident, 0)
	}

	result := make([]Incident, len(incidents))

	for i, incident := range incidents {
		result[i] = Incident{
			Id:        incident.ID().String(),
			CreatedAt: incident.CreatedAt(),
		}
	}

	return result
}
