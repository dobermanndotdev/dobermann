package monitors

import (
	"github.com/volatiletech/null/v8"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func mapMonitorToModel(m *monitor.Monitor) models.Monitor {
	return models.Monitor{
		ID:            m.ID().String(),
		AccountID:     m.AccountID().String(),
		EndpointURL:   m.EndpointUrl(),
		IsEndpointUp:  m.IsEndpointUp(),
		CreatedAt:     m.CreatedAt(),
		LastCheckedAt: null.TimeFromPtr(m.LastCheckedAt()),
	}
}
