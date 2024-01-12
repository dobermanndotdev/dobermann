package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/common/ptr"
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

	password, err := account.NewPassword(email.Address())
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
		time.Second*30,
		nil,
	)
	require.NoError(t, err)

	return newMonitor
}

func FixtureIncident(t *testing.T, monitorID string) *monitor.Incident {
	mID, err := domain.NewIdFromString(monitorID)
	require.NoError(t, err)

	incident, err := monitor.NewIncident(
		domain.NewID(),
		mID,
		nil,
		time.Now().UTC(),
		gofakeit.URL(),
		nil,
		gofakeit.Sentence(10),
		ptr.ToPtr(int16(500)),
	)
	require.NoError(t, err)

	return incident
}

type FixtureClient struct {
	Db  *sql.DB
	Ctx context.Context
}

func (f *FixtureClient) FixtureAndInsertIncidents(t *testing.T, monitorID domain.ID, count int) []models.Incident {
	var model models.Incident
	incidents := make([]models.Incident, count)

	for i := 0; i < count; i++ {
		model = models.Incident{
			ID:             domain.NewID().String(),
			MonitorID:      monitorID.String(),
			ResolvedAt:     null.TimeFromPtr(nil),
			Cause:          null.String{},
			ResponseStatus: null.Int16From(500),
			CheckedURL:     gofakeit.URL(),
			CreatedAt:      time.Now(),
		}
		require.NoError(t, model.Insert(f.Ctx, f.Db, boil.Infer()))

		incidentAction := models.IncidentAction{
			ID:                domain.NewID().String(),
			Description:       null.StringFrom(gofakeit.Sentence(20)),
			ActionType:        monitor.IncidentActionTypeCreated.String(),
			IncidentID:        model.ID,
			TakenByUserWithID: null.String{},
			At:                time.Now(),
		}
		require.NoError(t, incidentAction.Insert(f.Ctx, f.Db, boil.Infer()))

		incidents[i] = model
	}

	return incidents
}

func (f *FixtureClient) FixtureCheckResults(
	t *testing.T,
	monitorID domain.ID,
	startCheckedAt time.Time,
	rangeInDays, checksPerDay, checkIntervalInSeconds int,
) {
	for i := 0; i < rangeInDays; i++ {
		checkedAt := startCheckedAt.Add(time.Hour * 24 * time.Duration(i+1))

		for j := 0; j < checksPerDay; j++ {
			checkedAt = checkedAt.Add(time.Second * time.Duration(checkIntervalInSeconds))
			model := models.MonitorCheckResult{
				ID:               domain.NewID().String(),
				StatusCode:       null.Int16From(200),
				CheckedAt:        checkedAt,
				MonitorID:        monitorID.String(),
				ResponseTimeInMS: int16(gofakeit.Number(150, 350)),
				Region:           monitor.Regions[gofakeit.Number(0, len(monitor.Regions)-1)].String(),
			}

			require.NoError(t, model.Insert(f.Ctx, f.Db, boil.Infer()))
		}
	}
}

func EndpointUrlGenerator(isUp bool) string {
	isUpParam := "false"

	if isUp {
		isUpParam = "true"
	}

	return fmt.Sprintf("%s?id=%s&is_up=%s", SimulatorEndpointUrl, domain.NewID().String(), isUpParam)
}
