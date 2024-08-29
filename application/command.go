package application

import (
	"context"

	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/global"
)

type Command interface {
	Execute(ctx context.Context, message domain.Message) *global.BasicError
}
