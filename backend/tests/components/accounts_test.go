package components_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/flowck/doberman/internal/adapters/models"
	"github.com/flowck/doberman/internal/app/command"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/tests/client"
)

func TestAccounts(t *testing.T) {
	router, err := message.NewRouter(message.RouterConfig{}, wmLogger)

	eventAccountID := ""
	router.AddNoPublisherHandler("account_created_handler", command.AccountCreatedEvent{}.EventName(), subscriber, func(msg *message.Message) error {
		event, err := unMarshallMessageToEvent[command.AccountCreatedEvent](msg)
		if err != nil {
			t.Fatal(err)
		}

		eventAccountID = event.ID

		return nil
	})

	go func() {
		if err = router.Run(ctx); err != nil {
			t.Logf("Router failed: %v", err)
		}
	}()

	account := fixtureAccountRequest()
	resp01, err := cli.CreateAnAccount(ctx, account)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp01.StatusCode)
	accountID := findLastAccountCreatedID(t)

	assert.Eventually(t, func() bool {
		return eventAccountID == accountID
	}, time.Second*3, time.Millisecond*250)

	accountOwner := findAccountOwner(t, accountID)

	resp02, err := cli.ConfirmAccount(ctx, accountOwner.ConfirmationCode.String)
	require.NoError(t, err)
	accountOwner = findAccountOwner(t, accountID)

	assert.Equal(t, http.StatusOK, resp02.StatusCode)
	assert.Empty(t, accountOwner.ConfirmationCode.String, "Confirmation code should be empty after confirmation")

	resp03, err := cli.ConfirmAccount(ctx, domain.NewID().String())
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp03.StatusCode, "Account not found")
}

func findLastAccountCreatedID(t *testing.T) string {
	model, err := models.Accounts(qm.OrderBy("created_at DESC")).One(ctx, db)
	require.NoError(t, err)

	return model.ID
}

func findAccountOwner(t *testing.T, accID string) *models.User {
	model, err := models.Users(models.UserWhere.AccountID.EQ(accID)).One(ctx, db)
	require.NoError(t, err)

	return model
}

func fixtureAccountRequest() client.CreateAnAccountRequest {
	return client.CreateAnAccountRequest{
		Email:       gofakeit.Email(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		Password:    gofakeit.Password(true, true, true, true, false, 12),
		AccountName: gofakeit.Company(),
	}
}
