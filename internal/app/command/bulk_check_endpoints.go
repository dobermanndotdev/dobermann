package command

import (
	"context"
)

type BulkCheckEndpoints struct {
	FromRegion string
}

type BulkCheckEndpointsHandler struct {
	httpChecker httpChecker
	txProvider  TransactionProvider
}

func NewBulkCheckEndpointsHandler(httpChecker httpChecker, txProvider TransactionProvider) BulkCheckEndpointsHandler {
	return BulkCheckEndpointsHandler{
		txProvider:  txProvider,
		httpChecker: httpChecker,
	}
}

func (c BulkCheckEndpointsHandler) Execute(ctx context.Context, cmd BulkCheckEndpoints) error {
	return nil
}
