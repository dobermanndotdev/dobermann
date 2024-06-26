package resend

import (
	"context"
	"fmt"

	resendsdk "github.com/resendlabs/resend-go/v2"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

// Service This interface is needed to allow the initialisation of either the concrete struct service or a mock
// service to run the project in CI environments
type Service interface {
	SendEmailIncidentResolution(context.Context, *account.User, *monitor.Monitor, domain.ID) error
	SendEmailAboutIncident(context.Context, *account.User, *monitor.Monitor, *monitor.Incident) error
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
%s,<br>
An incident has been created for the monitor %s. <br>
For more details please follow the link %s. <br>

Dobermann - Endpoint monitoring
`, getGreetings(user), m.EndpointUrl(), getIncidentLink(s.hostname, incident.ID()))

	_, err := s.client.Emails.SendWithContext(ctx, &resendsdk.SendEmailRequest{
		From:    fmt.Sprintf("Dobermann <%s>", s.from),
		To:      []string{user.Email().Address()},
		Subject: fmt.Sprintf("[Dobermann] - New Incident for %s", m.EndpointUrl()),
		Html:    body,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendEmailIncidentResolution(
	ctx context.Context,
	user *account.User,
	m *monitor.Monitor,
	incidentID domain.ID,
) error {
	body := fmt.Sprintf(`
%s,<br>
The last incident reported on the monitor %s has been resolved. <br>
For more details please follow the link %s. <br>

Dobermann - Endpoint monitoring
`, getGreetings(user), m.EndpointUrl(), getIncidentLink(s.hostname, incidentID))

	_, err := s.client.Emails.SendWithContext(ctx, &resendsdk.SendEmailRequest{
		From:    fmt.Sprintf("Dobermann <%s>", s.from),
		To:      []string{user.Email().Address()},
		Subject: fmt.Sprintf("[Dobermann] - Incident resolved %s", m.EndpointUrl()),
		Html:    body,
	})
	if err != nil {
		return err
	}

	return nil
}

func getIncidentLink(host string, incidentID domain.ID) string {
	link := fmt.Sprintf("%s/dashboard/incidents/%s", host, incidentID)
	return fmt.Sprintf(`<a href="%s">%s</a>`, link, link)
}

func getGreetings(user *account.User) string {
	if user.FirstName() == "" {
		return "Hi"
	}

	return fmt.Sprintf("Hi %s", user.FirstName())
}
