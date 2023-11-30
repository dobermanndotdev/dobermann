package resend

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
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
