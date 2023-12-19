package http

import (
	"time"

	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func mapMonitorToResponseItem(m *monitor.Monitor) Monitor {
	return Monitor{
		CreatedAt:              m.CreatedAt(),
		EndpointUrl:            m.EndpointUrl(),
		Id:                     m.ID().String(),
		Incidents:              mapIncidentsToResponse(m.Incidents()),
		IsEndpointUp:           m.IsEndpointUp(),
		IsPaused:               m.IsPaused(),
		LastCheckedAt:          m.LastCheckedAt(),
		CheckIntervalInSeconds: int(m.CheckInterval().Seconds()),
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

func mapRequestToMonitor(body CreateMonitorRequest, user *authenticatedUser) (*monitor.Monitor, error) {
	subscriber, err := monitor.NewSubscriber(user.ID)
	if err != nil {
		return nil, err
	}

	return monitor.NewMonitor(
		domain.NewID(),
		body.EndpointUrl,
		user.AccountID,
		false,
		false,
		nil,
		[]*monitor.Subscriber{subscriber},
		time.Now().UTC(),
		time.Second*time.Duration(body.CheckIntervalInSeconds),
		nil,
	)
}

func mapMonitorResponseTimeStatsToResponse(stats query.ResponseTimeStats) GetMonitorResponseTimeStatsPayload {
	result := ResponseTimeStats{
		ResponseTimePerRegion: make([]ResponseTimePerRegion, len(stats.ResponseTimePerRegion)),
	}

	for i, region := range stats.ResponseTimePerRegion {
		result.ResponseTimePerRegion[i] = ResponseTimePerRegion{
			Region: region.Region.String(),
			Data:   make([]ResponseTimePerDate, len(region.Data)),
		}

		for j, dataPoint := range region.Data {

			result.ResponseTimePerRegion[i].Data[j] = ResponseTimePerDate{
				Date:  dataPoint.Date,
				Value: int(dataPoint.Value),
			}
		}
	}
	
	return GetMonitorResponseTimeStatsPayload{
		Data: result,
	}
}
