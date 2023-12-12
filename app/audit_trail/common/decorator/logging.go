package decorator

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger logger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	d.logger.Debugf("Executing command %s", handlerType)
	defer func() {
		if err == nil {
			d.logger.Info("Command executed successfully")
		} else {
			d.logger.Errorf(err, "failed to execute command: %w", err)
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
	d.logger.Debugf("Executing query %s", handlerType)
	defer func() {
		defer func() {
			if err == nil {
				d.logger.Info("Command executed successfully")
			} else {
				d.logger.Errorf(err, "failed to execute query: %w", err)
			}
		}()
	}()

	return d.base.Handle(ctx, cmd)
}
