package command_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/adapters/psql"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
	"github.com/dobermanndotdev/dobermann/tests"
)

var (
	ctx context.Context
)

func TestNewCreateMonitorHandler(t *testing.T) {
	account00 := tests.FixtureAccount(t)
	eventPublisher := events.NewPublisherMock()
	monitorRepository := psql.NewMonitorRepositoryMock()
	handler := command.NewCreateMonitorHandler(monitorRepository, eventPublisher)

	testCases := []struct {
		name        string
		expectedErr string
		monitor     *monitor.Monitor
	}{
		{
			name:        "error_invalid_monitor",
			expectedErr: "monitor cannot be invalid",
			monitor:     &monitor.Monitor{},
		},
		{
			name:        "account_created_successfully",
			expectedErr: "",
			monitor:     tests.FixtureMonitor(t, account00),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := handler.Execute(ctx, command.CreateMonitor{
				Monitor: tc.monitor,
			})

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)

			createdMonitor, err := monitorRepository.FindByID(ctx, tc.monitor.ID())
			require.NoError(t, err)
			assert.Equal(t, tc.monitor.ID(), createdMonitor.ID())

			event := eventPublisher.GetEventByID(tc.monitor.ID().String())
			require.NotNil(t, event)

			monitorCreatedEvent, ok := event.(events.MonitorCreatedEvent)
			require.True(t, ok)
			assert.Equal(t, monitorCreatedEvent.ID, tc.monitor.ID().String())
		})
	}
}

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	os.Exit(m.Run())
}
