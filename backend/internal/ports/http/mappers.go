package http

import (
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/monitor"
)

func mapReqBodyToMonitor(body CreateMonitorRequest) (*monitor.Monitor, error) {
	//TODO: review this
	accID, _ := domain.NewIdFromString("01H8YBGZF8QY158H4NP4NDXFHD")
	teamID, _ := domain.NewIdFromString("01HAF76XEDP97014GMFPKXACY0")

	alertTriggers := make([]monitor.AlertTrigger, len(body.AlertTriggers))
	for i, value := range body.AlertTriggers {
		alertTrigger, err := monitor.NewAlertTrigger(value)
		if err != nil {
			return nil, err
		}

		alertTriggers[i] = alertTrigger
	}

	endpoint, err := monitor.NewEndpoint(body.Endpoint)
	if err != nil {
		return nil, err
	}

	notificationMethods := make([]monitor.NotificationMethod, len(body.OnCallEscalation.NotificationMethods))
	for i, value := range body.OnCallEscalation.NotificationMethods {
		notificationMethod, err := monitor.NewNotificationMethod(value)
		if err != nil {
			return nil, err
		}

		notificationMethods[i] = notificationMethod
	}

	teamMembersToBeNotified := make([]domain.ID, len(body.OnCallEscalation.TeamMembersToBeNotified))
	for i, value := range body.OnCallEscalation.TeamMembersToBeNotified {
		teamMemberID, err := domain.NewIdFromString(value)
		if err != nil {
			return nil, err
		}

		teamMembersToBeNotified[i] = teamMemberID
	}

	onCallEscalation, err := monitor.NewOnCallEscalation(notificationMethods, teamMembersToBeNotified)
	if err != nil {
		return nil, err
	}

	return monitor.New(domain.NewID(), accID, teamID, alertTriggers, endpoint, onCallEscalation)
}
