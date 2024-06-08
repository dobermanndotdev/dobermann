package resend

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/common/logs"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type serviceMock struct {
	logger *logs.Logger
}

func NewServiceMock(logger *logs.Logger) Service {
	return &serviceMock{logger: logger}
}

func (s *serviceMock) SendEmailAboutIncident(
	ctx context.Context,
	user *account.User,
	m *monitor.Monitor,
	incident *monitor.Incident,
) error {
	s.logger.Warnf("[ResendMock] sending email about incident to %s", user.Email().Address())
	s.logger.Warnf("[ResendMock] monitor %s and incident id %s", m.EndpointUrl(), incident.ID())

	//IDEA: Simulate errors
	return nil
}

func (s *serviceMock) SendEmailIncidentResolution(ctx context.Context, user *account.User, m *monitor.Monitor, incidentID domain.ID) error {
	s.logger.Warnf("[ResendMock] sending email about incident being resolved to %s", user.Email().Address())
	s.logger.Warnf("[ResendMock] monitor %s", m.EndpointUrl())

	return nil
}
