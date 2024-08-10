package handlers

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

// Middleware is a function that wraps a MessageHandler.
type Middleware func(MessageHandler) MessageHandler

func LoggingMiddleware(l *logrus.Logger) Middleware {
	return func(next MessageHandler) MessageHandler {
		return MessageHandlerFunc(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			l.Infof("Handling message: %s", *event.Event.Message.Content)
			err := next.Handle(ctx, event)
			l.Infof("Finishing handling message: %s", *event.Event.Message.Content)
			if err != nil {
				l.Errorf("handling message: %s failed: %s", *event.Event.Message.Content, err.Error())
			}
			return err
		})
	}
}

type MessageHandlerFunc func(ctx context.Context, event *larkim.P2MessageReceiveV1) error

func (f MessageHandlerFunc) Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return f(ctx, event)
}
