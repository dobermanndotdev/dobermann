package http

import (
	"time"

	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
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
		UpSince:                m.UpSince(),
		DownSince:              m.DownSince(),
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
			Cause:      incident.Details().Cause,
			CheckedUrl: incident.CheckedURL(),
			CreatedAt:  incident.CreatedAt(),
			Id:         incident.ID().String(),
			ResolvedAt: incident.ResolvedAt(),
		}
	}

	return result
}

func mapIncidentToFullIncidentResponse(incident *monitor.Incident) FullIncident {
	return FullIncident{
		Id:              incident.ID().String(),
		MonitorId:       incident.MonitorID().String(),
		CreatedAt:       incident.CreatedAt(),
		CheckedUrl:      incident.CheckedURL(),
		ResolvedAt:      incident.ResolvedAt(),
		Cause:           incident.Details().Cause,
		ResponseStatus:  int(incident.Details().Status),
		ResponseBody:    incident.Details().ResponseBody,
		RequestHeaders:  incident.Details().RequestHeaders,
		ResponseHeaders: incident.Details().RequestHeaders,
	}
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

func mapMonitorResponseTimeStatsToResponse(stats []query.ResponseTimeStat) GetMonitorResponseTimeStatsPayload {
	result := make([]ResponseTimeStat, len(stats))

	for i, stat := range stats {
		result[i] = ResponseTimeStat{
			Region: stat.Region,
			Date:   stat.Date,
			Value:  stat.Value,
		}
	}

	return GetMonitorResponseTimeStatsPayload{
		Data: result,
	}
}

func mapUserToResponse(user *account.User) GetProfileDetailsPayload {
	return GetProfileDetailsPayload{Data: User{
		Id:        user.ID().String(),
		Email:     user.Email().Address(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Role:      user.Role().String(),
		CreatedAt: user.CreatedAt(),
	}}
}
