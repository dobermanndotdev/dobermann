package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func FixtureMonitor(t *testing.T, accountID, teamID, memberID domain.ID) *monitor.Monitor {
	alertTriggers := []monitor.AlertTrigger{
		monitor.AlertTriggerUrlIsUnavailable,
	}

	endpoint, err := monitor.NewEndpoint(gofakeit.URL())
	require.NoError(t, err)

	onCallEscalation, err := monitor.NewOnCallEscalation([]monitor.NotificationMethod{
		monitor.NotificationMethodEmail,
	}, []domain.ID{
		memberID,
	})
	require.NoError(t, err)

	mo, err := monitor.New(domain.NewID(), accountID, teamID, alertTriggers, endpoint, onCallEscalation)
	require.NoError(t, err)

	return mo
}

func FixtureAndSaveTeamMember(ctx context.Context, t *testing.T, db *sql.DB, accountID domain.ID) (domain.ID, models.User) {
	memberID := domain.NewID()
	model := models.User{
		ID:                   memberID.String(),
		AccountID:            accountID.String(),
		FirstName:            gofakeit.FirstName(),
		LastName:             gofakeit.LastName(),
		PrimaryPhoneNumber:   gofakeit.Phone(),
		SecondaryPhoneNumber: gofakeit.Phone(),
		Email:                gofakeit.Email(),
		AvatarURL:            null.String{},
		Timezone:             null.String{},
		OnHolidaysUntil:      null.Time{},
	}
	require.NoError(t, model.Insert(ctx, db, boil.Infer()))

	return memberID, model
}

func FixtureAndSaveTeam(ctx context.Context, t *testing.T, db *sql.DB, accountID domain.ID) (domain.ID, models.Team) {
	teamID := domain.NewID()
	model := models.Team{
		ID:        teamID.String(),
		AccountID: accountID.String(),
		Name:      gofakeit.Company(),
	}
	require.NoError(t, model.Insert(ctx, db, boil.Infer()))

	return teamID, model
}

func FixtureAndSaveAccount(ctx context.Context, t *testing.T, db *sql.DB) domain.ID {
	accID := domain.NewID()
	model := models.Account{
		ID:   accID.String(),
		Name: gofakeit.Company(),
	}

	require.NoError(t, model.Insert(ctx, db, boil.Infer()))

	return accID
}
