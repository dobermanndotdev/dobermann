package observability

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
)

type correlationIdKey string

var (
	corrIdKey correlationIdKey = "appCorrelationId"
)

func NewContextWithCorrelationID(ctx context.Context) context.Context {
	return context.WithValue(ctx, corrIdKey, ulid.Make().String())
}

func NewContextWithCorrelationIdFromString(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, corrIdKey, id)
}

func GetCorrelationIdFromContext(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(corrIdKey))
}
