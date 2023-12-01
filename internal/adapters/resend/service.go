package resend

import (
	"context"
	"fmt"

	resendsdk "github.com/resendlabs/resend-go/v2"

	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

// Service This interface is needed to allow the initialisation of either the concrete struct service or a mock
// service to run the project in CI environments
type Service interface {
	SendEmailAboutIncident(
		ctx context.Context,
		user *account.User,
		m *monitor.Monitor,
		incident *monitor.Incident,
	) error
}

type service struct {
	from     string
	hostname string
	client   *resendsdk.Client
}

func NewService(apiKey, from, hostname string) Service {
	return &service{
		from:     from,
		hostname: hostname,
		client:   resendsdk.NewClient(apiKey),
	}
}

func (s *service) SendEmailAboutIncident(
	ctx context.Context,
	user *account.User,
	m *monitor.Monitor,
	incident *monitor.Incident,
) error {
	body := fmt.Sprintf(`
Hi %s,</br>
An incident has been created for the monitor %s. </br>
For more details please follow the link %s/monitors/%s/incidents/%s. </br>

Dobermann - Endpoint monitoring
`, user.FirstName(), m.EndpointUrl(), s.hostname, m.ID().String(), incident.ID().String())

	_, err := s.client.Emails.SendWithContext(ctx, &resendsdk.SendEmailRequest{
		From:    s.from,
		To:      []string{user.Email().Address()},
		Subject: fmt.Sprintf("[Dobermann] - New Incident for %s", m.EndpointUrl()),
		Html:    body,
	})
	if err != nil {
		return err
	}

	return nil
}
