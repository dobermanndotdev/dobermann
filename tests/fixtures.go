package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

const (
	SimulatorEndpointUrl = "http://endpoint_simulator:8090" // Hostname within docker's network
)

func FixturePassword() string {
	return gofakeit.Password(true, true, true, true, false, 12)
}

func FixtureAccount(t *testing.T) *account.Account {
	email, err := account.NewEmail(gofakeit.Email())
	require.NoError(t, err)

	password, err := account.NewPassword(FixturePassword())
	require.NoError(t, err)

	acc, err := account.NewFirstTimeAccount(gofakeit.Company(), email, password)
	require.NoError(t, err)

	return acc
}

func FixtureAndInsertAccount(t *testing.T, db *sql.DB, insertUser bool) *account.Account {
	acc := FixtureAccount(t)

	model := models.Account{
		Name: acc.Name(),
		ID:   acc.ID().String(),
	}

	require.NoError(t, model.Insert(context.Background(), db, boil.Infer()))

	if !insertUser {
		return acc
	}

	owner, err := acc.FirstAccountOwner()
	require.NoError(t, err)

	userModel := models.User{
		ID:        owner.ID().String(),
		Email:     owner.Email().Address(),
		Password:  owner.Password().String(),
		Role:      owner.Role().String(),
		AccountID: acc.ID().String(),
	}

	require.NoError(t, userModel.Insert(context.Background(), db, boil.Infer()))

	return acc
}

func FixtureMonitor(t *testing.T, acc *account.Account) *monitor.Monitor {
	subscribers := make([]*monitor.Subscriber, len(acc.Users()))

	var err error
	var subscriber *monitor.Subscriber
	for i, user := range acc.Users() {
		subscriber, err = monitor.NewSubscriber(user.ID())
		require.NoError(t, err)

		subscribers[i] = subscriber
	}

	newMonitor, err := monitor.NewMonitor(
		domain.NewID(),
		SimulatorEndpointUrl,
		acc.ID(),
		false,
		false,
		nil,
		subscribers,
		time.Now().UTC(),
		nil,
	)
	require.NoError(t, err)

	return newMonitor
}

func FixtureIncident(t *testing.T) *monitor.Incident {
	incident, err := monitor.NewIncident(domain.NewID(), false, time.Now().UTC(), nil)
	require.NoError(t, err)

	return incident
}
