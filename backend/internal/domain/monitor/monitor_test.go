package monitor_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/monitor"
)

func TestNewMonitor(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string

		id                  domain.ID
		teamID              domain.ID
		accountID           domain.ID
		endpoint            string
		peopleToBeNotified  []domain.ID
		alertTriggers       []monitor.AlertTrigger
		notificationMethods []monitor.NotificationMethod
	}{
		{
			name:        "new_monitor",
			expectedErr: "",

			id:                  domain.NewID(),
			teamID:              domain.NewID(),
			accountID:           domain.NewID(),
			endpoint:            gofakeit.URL(),
			alertTriggers:       []monitor.AlertTrigger{monitor.AlertTriggerUrlIsUnavailable},
			peopleToBeNotified:  []domain.ID{domain.NewID()},
			notificationMethods: []monitor.NotificationMethod{monitor.NotificationMethodSMS, monitor.NotificationMethodEmail},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			endpoint, err := monitor.NewEndpoint(tc.endpoint)
			require.NoError(t, err)

			onCallEscalation, err := monitor.NewOnCallEscalation(
				tc.notificationMethods,
				tc.peopleToBeNotified,
			)
			require.NoError(t, err)

			moni, err := monitor.New(
				tc.id,
				tc.accountID,
				tc.teamID,
				tc.alertTriggers,
				endpoint,
				onCallEscalation,
			)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.id, moni.ID())
			assert.Equal(t, tc.accountID, moni.AccountID())
			assert.Equal(t, tc.teamID, moni.TeamID())
		})
	}
}
