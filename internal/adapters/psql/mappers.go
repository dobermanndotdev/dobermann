package psql

import (
	"fmt"
	"math"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func mapUserToModel(user *account.User) models.User {
	return models.User{
		ID:        user.ID().String(),
		FirstName: null.StringFrom(user.FirstName()),
		LastName:  null.StringFrom(user.LastName()),
		Email:     user.Email().Address(),
		Password:  user.Password().String(),
		Role:      user.Role().String(),
		AccountID: user.AccountID().String(),
		CreatedAt: user.CreatedAt(),
	}
}

func mapModelToUser(model *models.User) (*account.User, error) {
	id, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	accountId, err := domain.NewIdFromString(model.AccountID)
	if err != nil {
		return nil, err
	}

	email, err := account.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	role, err := account.NewRole(model.Role)
	if err != nil {
		return nil, err
	}

	password, err := account.NewPasswordFromHash(model.Password)
	if err != nil {
		return nil, err
	}

	return account.NewUser(id, model.FirstName.String, model.LastName.String, email, role, password, accountId)
}

func mapMonitorToModel(m *monitor.Monitor) *models.Monitor {
	return &models.Monitor{
		ID:                     m.ID().String(),
		AccountID:              m.AccountID().String(),
		EndpointURL:            m.EndpointUrl(),
		IsPaused:               m.IsPaused(),
		IsEndpointUp:           m.IsEndpointUp(),
		CreatedAt:              m.CreatedAt(),
		CheckIntervalInSeconds: int(math.Floor(m.CheckInterval().Seconds())),
		LastCheckedAt:          null.TimeFromPtr(m.LastCheckedAt()),
	}
}

func mapMonitorsToModels(monitorList []*monitor.Monitor) []*models.Monitor {
	result := make([]*models.Monitor, len(monitorList))

	for i, m := range monitorList {
		result[i] = mapMonitorToModel(m)
	}

	return result
}

func mapModelsToMonitors(modelList []*models.Monitor) ([]*monitor.Monitor, error) {
	result := make([]*monitor.Monitor, len(modelList))

	var err error
	var m *monitor.Monitor

	for i, model := range modelList {
		m, err = mapModelToMonitor(model)
		if err != nil {
			return nil, err
		}

		result[i] = m
	}

	return result, nil
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

	var incidents []*monitor.Incident
	if model.R != nil && model.R.Incidents != nil {
		incidents, err = mapModelsToIncidents(model.R.Incidents)
		if err != nil {
			return nil, fmt.Errorf("error mapping incidents to monitor: %v", err)
		}
	}

	var subscribers []*monitor.Subscriber
	if model.R != nil && model.R.Users != nil {
		subscribers, err = mapModelsToSubscribers(model.R.Users)
		if err != nil {
			return nil, fmt.Errorf("error mapping subscribers to monitor: %v", err)
		}
	}

	return monitor.NewMonitor(
		id,
		model.EndpointURL,
		accountID,
		model.IsEndpointUp,
		model.IsPaused,
		incidents,
		subscribers,
		model.CreatedAt,
		time.Second*time.Duration(model.CheckIntervalInSeconds),
		model.LastCheckedAt.Ptr(),
	)
}

func mapModelsToSubscribers(modelList []*models.User) ([]*monitor.Subscriber, error) {
	result := make([]*monitor.Subscriber, len(modelList))

	var err error
	var subscriber *monitor.Subscriber
	for i, m := range modelList {
		subscriber, err = mapModelToSubscriber(m)
		if err != nil {
			return nil, err
		}

		result[i] = subscriber
	}

	return result, nil
}

func mapModelToSubscriber(model *models.User) (*monitor.Subscriber, error) {
	userID, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	return monitor.NewSubscriber(userID)
}

func mapIncidentToModel(incident *monitor.Incident, monitorID domain.ID) *models.Incident {
	return &models.Incident{
		ID:         incident.ID().String(),
		IsResolved: incident.IsResolved(),
		MonitorID:  monitorID.String(),
		CreatedAt:  incident.CreatedAt(),
	}
}

func mapModelToIncident(model *models.Incident) (*monitor.Incident, error) {
	id, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	return monitor.NewIncident(id, model.IsResolved, model.CreatedAt, nil)
}

func mapModelsToIncidents(modelList []*models.Incident) ([]*monitor.Incident, error) {
	result := make([]*monitor.Incident, len(modelList))

	var err error
	var incident *monitor.Incident
	for i, m := range modelList {
		incident, err = mapModelToIncident(m)
		if err != nil {
			return nil, err
		}

		result[i] = incident
	}

	return result, nil
}

func mapAccountToModel(acc *account.Account) models.Account {
	return models.Account{
		ID:         acc.ID().String(),
		Name:       acc.Name(),
		VerifiedAt: null.TimeFromPtr(acc.VerifiedAt()),
		CreatedAt:  time.Now(),
	}
}

func mapPaginationParamsToOffset(page, limit int) int {
	// page from query is 1-based whereas postgres offset is 0-based
	p := page - 1

	return p * limit
}

func mapPaginationPerPageCount(total int64, limit int) int {
	// round up
	return int(math.Ceil(float64(total) / float64(limit)))
}

func mapIncidentActionToModel(action *monitor.IncidentAction, incidentID domain.ID) models.IncidentAction {
	model := models.IncidentAction{
		ID:          action.Id().String(),
		Description: null.StringFrom(action.Description()),
		ActionType:  action.ActionType().String(),
		IncidentID:  incidentID.String(),
		At:          action.TakenAt(),
	}

	if action.TakerUserID() != nil {
		value := action.TakerUserID().String()
		model.TakenByUserWithID = null.StringFrom(value)
	}

	return model
}

func mapCheckResultToModel(monitorID domain.ID, checkResult *monitor.CheckResult) models.MonitorCheckResult {
	return models.MonitorCheckResult{
		MonitorID:        monitorID.String(),
		StatusCode:       checkResult.StatusCode(),
		Region:           checkResult.Region().String(),
		CheckedAt:        checkResult.CheckedAt(),
		ResponseTimeInMS: checkResult.ResponseTimeInMs(),
	}
}
