package decorator

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger logger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	d.logger.Debug(fmt.Sprintf("Executing command: %s", handlerType))
	defer func() {
		if err == nil {
			d.logger.Info("Command executed successfully")
		} else {
			d.logger.Error(fmt.Errorf("failed execute the command: %w", err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger logger.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	handlerType := generateActionName(cmd)

	d.logger.Debug(fmt.Sprintf("Executing query: %s", handlerType))
	defer func() {
		if err == nil {
			d.logger.Info("Query executed successfully")
		} else {
			d.logger.Error(fmt.Errorf("failed execute the query: %w", err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}
