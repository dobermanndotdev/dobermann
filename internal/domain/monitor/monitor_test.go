package monitor_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	"github.com/flowck/dobermann/backend/tests"
)

func TestMonitor_NewMonitor(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string

		id            domain.ID
		accountID     domain.ID
		isEndpointUp  bool
		isPaused      bool
		endpointUrl   string
		createdAt     time.Time
		lastCheckedAt *time.Time
		checkInterval time.Duration
		incidents     []*monitor.Incident
		subscribers   []*monitor.Subscriber
	}{
		{
			name:          "new_monitor",
			expectedErr:   "",
			id:            domain.NewID(),
			accountID:     domain.NewID(),
			isEndpointUp:  true,
			isPaused:      false,
			endpointUrl:   gofakeit.URL(),
			createdAt:     time.Now(),
			lastCheckedAt: nil,
			checkInterval: time.Second * 30,
			incidents:     nil,
			subscribers:   nil,
		},
		{
			name:          "error_invalid_id",
			expectedErr:   "id cannot be invalid",
			id:            domain.ID{},
			accountID:     domain.NewID(),
			isEndpointUp:  true,
			isPaused:      false,
			endpointUrl:   gofakeit.URL(),
			createdAt:     time.Now(),
			lastCheckedAt: nil,
			checkInterval: time.Second * 30,
			incidents:     nil,
			subscribers:   nil,
		},
		{
			name:          "error_empty_or_invalid_url",
			expectedErr:   "the url must start with http/https/tcp instead of ''",
			id:            domain.NewID(),
			accountID:     domain.NewID(),
			isEndpointUp:  true,
			isPaused:      false,
			endpointUrl:   "",
			createdAt:     time.Now(),
			lastCheckedAt: nil,
			checkInterval: time.Second * 30,
			incidents:     nil,
			subscribers:   nil,
		},
		{
			name:          "error_created_at_set_in_the_future",
			expectedErr:   "createdAt cannot be set in the future",
			id:            domain.NewID(),
			accountID:     domain.NewID(),
			isEndpointUp:  true,
			isPaused:      false,
			endpointUrl:   gofakeit.URL(),
			createdAt:     time.Date(time.Now().Year()+1, 1, 1, 1, 1, 1, 1, time.UTC),
			lastCheckedAt: nil,
			checkInterval: time.Second * 30,
			incidents:     nil,
			subscribers:   nil,
		},
		{
			name:          "error_too_low_of_a_check_interval",
			expectedErr:   "checkInterval cannot be less than 30 seconds",
			id:            domain.NewID(),
			accountID:     domain.NewID(),
			isEndpointUp:  true,
			isPaused:      false,
			endpointUrl:   gofakeit.URL(),
			createdAt:     time.Now(),
			lastCheckedAt: nil,
			checkInterval: time.Second * 29,
			incidents:     nil,
			subscribers:   nil,
		},
		{
			name:          "error_last_checked_at_set_in_the_future",
			expectedErr:   "lastCheckedAt cannot be set in the future",
			id:            domain.NewID(),
			accountID:     domain.NewID(),
			isEndpointUp:  true,
			isPaused:      false,
			endpointUrl:   gofakeit.URL(),
			createdAt:     time.Now(),
			lastCheckedAt: tests.ToPtr(time.Date(time.Now().Year()+1, 1, 1, 1, 1, 1, 1, time.UTC)),
			checkInterval: time.Second * 30,
			incidents:     nil,
			subscribers:   nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m, err := monitor.NewMonitor(
				tc.id,
				tc.endpointUrl,
				tc.accountID,
				tc.isEndpointUp,
				tc.isPaused,
				tc.incidents,
				tc.subscribers,
				tc.createdAt,
				tc.checkInterval,
				tc.lastCheckedAt,
			)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, m)

			endpointURL, err := monitor.NewURL(tc.endpointUrl)
			require.NoError(t, err)

			assert.Equal(t, tc.id, m.ID())
			assert.Equal(t, tc.accountID, m.AccountID())
			assert.Equal(t, endpointURL.String(), m.EndpointUrl())
			assert.Equal(t, tc.isEndpointUp, m.IsEndpointUp())
			assert.Equal(t, tc.isPaused, m.IsPaused())
			assert.ElementsMatch(t, tc.incidents, m.Incidents())
			assert.ElementsMatch(t, tc.subscribers, m.Subscribers())
			assert.Equal(t, tc.createdAt, m.CreatedAt())
			assert.Equal(t, tc.checkInterval, m.CheckInterval())
			assert.Equal(t, tc.lastCheckedAt, m.LastCheckedAt())
		})
	}
}

func TestMonitor_Edit(t *testing.T) {
	// Given
	m := mustNewMonitor(t)
	newEndpointURL := gofakeit.URL()
	newCheckInterval := time.Second * 180

	// When
	err := m.Edit(newEndpointURL, newCheckInterval)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, newEndpointURL, m.EndpointUrl())
	assert.Equal(t, newCheckInterval, m.CheckInterval())
}

func TestMonitor_Pause(t *testing.T) {
	// Given
	m := mustNewMonitor(t)
	require.False(t, m.IsPaused())

	// When
	m.Pause()

	// Then
	assert.True(t, m.IsPaused())
}

func TestMonitor_UnPause(t *testing.T) {
	// Given
	m := mustNewMonitor(t)
	m.Pause()
	require.True(t, m.IsPaused())

	// When
	m.UnPause()

	// Then
	assert.False(t, m.IsPaused())
}

func TestMonitor_HasIncidentUnresolved(t *testing.T) {
	// Given/When
	m := mustNewMonitor(t)

	// Then
	assert.True(t, m.HasIncidentUnresolved())
}

func TestMonitor_MarkEndpointAsDown(t *testing.T) {
	// Given
	m := mustNewMonitor(t)
	require.True(t, m.IsEndpointUp())

	// When
	m.MarkEndpointAsDown()

	// Then
	assert.False(t, m.IsEndpointUp())
}

func TestMonitor_MarkEndpointAsUp(t *testing.T) {
	// Given
	m := mustNewMonitor(t)
	require.True(t, m.IsEndpointUp())

	// When
	m.MarkEndpointAsDown()

	// Then
	assert.False(t, m.IsEndpointUp())
}

func mustNewMonitor(t *testing.T) *monitor.Monitor {
	monitorID := domain.NewID()
	endpointURL := gofakeit.URL()
	createdAt := time.Date(time.Now().Year()-1, 1, 1, 1, 1, 1, 1, time.UTC)

	var err error
	var incident *monitor.Incident
	var incidents []*monitor.Incident

	for i := 0; i < 10; i++ {
		incident, err = monitor.NewIncident(
			domain.NewID(),
			monitorID,
			nil,
			time.Now(),
			endpointURL,
			nil,
			"unresponsive url",
			tests.ToPtr(int16(http.StatusInternalServerError)),
		)

		incidents = append(incidents, incident)
	}

	m, err := monitor.NewMonitor(
		monitorID,
		endpointURL,
		domain.NewID(),
		true,
		false,
		incidents,
		nil,
		createdAt,
		time.Second*30,
		nil,
	)
	require.NoError(t, err)

	return m
}

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	os.Exit(m.Run())
}
