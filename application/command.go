package application

import (
	"context"

	"github.com/wawayes/lark-bot/domain"
)

type Command interface {
	Execute(ctx context.Context, message domain.Message) error
}
