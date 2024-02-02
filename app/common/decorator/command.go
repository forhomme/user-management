package decorator

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"strings"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], logger *baselogger.Logger, tracer *telemetry.OtelSdk) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandMetricsDecorator[H]{
			base:   handler,
			client: tracer,
		},
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
