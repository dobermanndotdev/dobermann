package monitors

import (
	"github.com/volatiletech/null/v8"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func mapMonitorToModel(m *monitor.Monitor) *models.Monitor {
	return &models.Monitor{
		ID:            m.ID().String(),
		AccountID:     m.AccountID().String(),
		EndpointURL:   m.EndpointUrl(),
		IsEndpointUp:  m.IsEndpointUp(),
		CreatedAt:     m.CreatedAt(),
		LastCheckedAt: null.TimeFromPtr(m.LastCheckedAt()),
	}
}

func mapModelToMonitor(model *models.Monitor) (*monitor.Monitor, error) {
	id, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	accountID, err := domain.NewIdFromString(model.AccountID)
	if err != nil {
		return nil, err
	}

	return monitor.NewMonitor(
		id,
		model.EndpointURL,
		accountID,
		model.IsEndpointUp,
		nil,
		model.CreatedAt,
		model.LastCheckedAt.Ptr(),
	)
}

func mapIncidentToModel(incident *monitor.Incident, monitorID domain.ID) *models.Incident {
	return &models.Incident{
		ID:        incident.ID().String(),
		MonitorID: monitorID.String(),
		CreatedAt: incident.CreatedAt(),
	}
}
