package application

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/domain"
)

type LoggingDecorator struct {
	command Command
	logger  *logrus.Logger
}

func NewLoggingDecorator(command Command, logger *logrus.Logger) *LoggingDecorator {
	return &LoggingDecorator{
		command: command,
		logger:  logger,
	}
}

func (d *LoggingDecorator) Execute(ctx context.Context, message domain.Message) error {
	d.logger.Infof("Executing command: %T", d.command)
	err := d.command.Execute(ctx, message)
	if err != nil {
		d.logger.Errorf("Command failed: %v", err)
	}
	return err
}
