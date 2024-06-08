package psql

import (
	"fmt"
	"math"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/dobermanndotdev/dobermann/internal/adapters/models"
	"github.com/dobermanndotdev/dobermann/internal/app/query"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

func mapUserToModel(user *account.User) models.User {
	return models.User{
		ID:              user.ID().String(),
		FirstName:       null.String{},
		LastName:        null.String{},
		Email:           user.Email().Address(),
		LoginProviderID: user.LoginProviderID(),
		Role:            user.Role().String(),
		AccountID:       user.AccountID().String(),
		CreatedAt:       user.CreatedAt(),
	}
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

func mapIncidentToModel(incident *monitor.Incident) *models.Incident {
	return &models.Incident{
		ID:             incident.ID().String(),
		MonitorID:      incident.MonitorID().String(),
		ResolvedAt:     null.TimeFromPtr(incident.ResolvedAt()),
		Cause:          null.StringFrom(incident.Cause()),
		ResponseStatus: null.Int16FromPtr(incident.ResponseStatusCode()),
		CheckedURL:     incident.CheckedURL(),
		CreatedAt:      incident.CreatedAt(),
	}
}

func mapModelToIncident(model *models.Incident) (*monitor.Incident, error) {
	id, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	monitorID, err := domain.NewIdFromString(model.MonitorID)
	if err != nil {
		return nil, err
	}

	return monitor.NewIncident(id,
		monitorID,
		model.ResolvedAt.Ptr(),
		model.CreatedAt,
		model.CheckedURL,
		nil,
		model.Cause.String,
		model.ResponseStatus.Ptr(),
	)
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
		ID:        acc.ID().String(),
		CreatedAt: time.Now(),
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
		ID:               checkResult.ID().String(),
		MonitorID:        monitorID.String(),
		CheckedAt:        checkResult.CheckedAt(),
		Region:           checkResult.Region().String(),
		ResponseTimeInMS: checkResult.ResponseTimeInMs(),
		StatusCode:       null.Int16FromPtr(checkResult.StatusCode()),
	}
}

func mapModelToResponseTimeStats(modelList []*models.MonitorCheckResult) []query.ResponseTimeStat {
	result := make([]query.ResponseTimeStat, len(modelList))

	for i, model := range modelList {
		result[i] = query.ResponseTimeStat{
			Value:  int(model.ResponseTimeInMS),
			Region: model.Region,
			Date:   model.CheckedAt,
		}
	}

	return result
}
