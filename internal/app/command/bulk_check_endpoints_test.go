package command_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/adapters/endpoint_checkers"
	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/adapters/psql"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/common/logs"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
	"github.com/dobermanndotdev/dobermann/tests"
)

type mockTxProvider struct {
	EventPublisher    *events.PublisherMock
	MonitorRepository psql.MonitorRepositoryMock
}

func (p mockTxProvider) Transact(ctx context.Context, fn command.TransactFunc) error {
	return fn(command.TransactableAdapters{
		AccountRepository:  nil,
		UserRepository:     nil,
		IncidentRepository: nil,
		MonitorRepository:  p.MonitorRepository,
		EventPublisher:     p.EventPublisher,
	})
}

func TestNewBulkCheckEndpointsHandler(t *testing.T) {
	endpointsChecker, err := endpoint_checkers.NewHttpChecker("europe", 5, logs.New(false))
	require.NoError(t, err)

	eventPublisher := events.NewPublisherMock()
	txProvider := mockTxProvider{
		EventPublisher:    eventPublisher,
		MonitorRepository: psql.NewMonitorRepositoryMock(),
	}
	monitorRepository := psql.NewMonitorRepositoryMock()

	account00 := tests.FixtureAccount(t)
	handler := command.NewBulkCheckEndpointsHandler(endpointsChecker, txProvider, eventPublisher, monitorRepository)

	testCases := []struct {
		name         string
		withIncident bool
		monitor      *monitor.Monitor
	}{
		{
			name:    "monitor_is_down_with_incident",
			monitor: fixtureMonitor(t, account00, false),
		},
		{
			name:    "monitor_is_up_with_incident_to_resolve",
			monitor: fixtureMonitor(t, account00, true),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err = txProvider.MonitorRepository.Insert(ctx, tc.monitor)
			require.NoError(t, err)

			err = handler.Execute(ctx, command.BulkCheckEndpoints{})
			require.NoError(t, err)

			foundMonitor, err := txProvider.MonitorRepository.FindByID(ctx, tc.monitor.ID())
			require.NoError(t, err)

			assert.NotNil(t, foundMonitor.LastCheckedAt())

			require.NotEmpty(t, txProvider.EventPublisher.Events())

			event := txProvider.EventPublisher.GetEventByID(tc.monitor.ID().String())
			require.NotNil(t, event)

			if tc.monitor.IsEndpointUp() {
				assert.Equal(t, events.EndpointCheckSucceededEvent{}.EventName(), event.EventName())
			} else {
				assert.Equal(t, events.EndpointCheckFailed{}.EventName(), event.EventName())
			}
		})
	}
}

func fixtureMonitor(t *testing.T, acc *account.Account, isUp bool) *monitor.Monitor {
	subscribers := make([]*monitor.Subscriber, len(acc.Users()))
	endpointUrl := tests.EndpointUrlGenerator(isUp)
	endpointUrl = strings.Replace(endpointUrl, "endpoint_simulator", "localhost", -1)

	var err error
	var subscriber *monitor.Subscriber
	for i, user := range acc.Users() {
		subscriber, err = monitor.NewSubscriber(user.ID())
		require.NoError(t, err)

		subscribers[i] = subscriber
	}

	monitorID := domain.NewID()

	newMonitor, err := monitor.NewMonitor(
		monitorID,
		endpointUrl,
		acc.ID(),
		isUp,
		false,
		[]*monitor.Incident{tests.FixtureIncident(t, monitorID.String())},
		subscribers,
		time.Now().UTC(),
		time.Second*30,
		nil,
	)
	require.NoError(t, err)

	return newMonitor
}
