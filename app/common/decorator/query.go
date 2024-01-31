package decorator

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], logger logger.Logger, tracer *telemetry.OtelSdk) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: tracer,
		},
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
