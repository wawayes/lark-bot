package handlers

import (
	"context"
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 接收消息
func ReceiveHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	HandleChatMsg(ctx, event)
	return nil
}

// 消息已读
func ReadHandler(ctx context.Context, event *larkim.P2MessageReadV1) error {
	readId := event.Event.Reader.ReaderId.OpenId
	fmt.Printf("msg is read by : %v \n", *readId)
	return nil
}
