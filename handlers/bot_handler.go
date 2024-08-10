package handlers

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type BotHandler interface {
	Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}
