package command_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/adapters/transaction"
	"github.com/flowck/dobermann/backend/internal/app/command"
)

func TestNewBulkCheckEndpointsHandler(t *testing.T) {
	endpointsChecker, err := endpoint_checkers.NewHttpChecker("europe")
	require.NoError(t, err)
	txProvider := transaction.NewPsqlProviderMock()
	monitorRepository := psql.NewMonitorRepositoryMock()

	handler := command.NewBulkCheckEndpointsHandler(endpointsChecker, txProvider, monitorRepository)

	err = handler.Execute(ctx, command.BulkCheckEndpoints{})
	require.NoError(t, err)
}
