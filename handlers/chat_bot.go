package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type ChatBot struct {
	factory    *HandlerFactory
	middleware []Middleware
}

func NewChatBot(factory *HandlerFactory, middleware ...Middleware) *ChatBot {
	return &ChatBot{
		factory:    factory,
		middleware: middleware,
	}
}

func (bot *ChatBot) ReadHandler(ctx context.Context, event *larkim.P2MessageReadV1) error {
	userID := event.Event.Reader.ReaderId.UserId
	bot.factory.Logger.Infof("msg is read by : %v \n", *userID)
	return nil
}

func (bot *ChatBot) HandleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	msgType := *event.Event.Message.MessageType
	handler := bot.factory.CreateHandler(msgType)
	if handler == nil {
		return errors.New("unsupported message type")
	}

	for _, m := range bot.middleware {
		handler = m(handler)
	}

	return handler.Handle(ctx, event)
}

func GetBotInfo(client *lark.Client) (string, error) {
	resp, err := client.Do(context.Background(),
		&larkcore.ApiReq{
			HttpMethod:                http.MethodGet,
			ApiPath:                   "https://open.feishu.cn/open-apis/bot/v3/info",
			Body:                      nil,
			SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
		},
	)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(bytes.NewReader(resp.RawBody))
	if err != nil {
		return "", err
	}
	fmt.Println(string(b))
	return string(b), nil
}
