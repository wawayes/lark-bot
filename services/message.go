package services

import (
	"context"

	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

type MessageService struct {
	larkClient *adapters.LarkClient
}

func NewMessageService(larkClient *adapters.LarkClient) *MessageService {
	return &MessageService{
		larkClient: larkClient,
	}
}

func (s *MessageService) SendReply(ctx context.Context, message domain.Message, content string) error {
	reply := domain.Reply{
		ReceiveID:     message.MessageID,
		ReceiveIDType: "message_id",
		Content:       content,
		MsgType:       domain.MsgTypeText,
	}
	return s.larkClient.ReplyMsg(ctx, reply)
}

func (s *MessageService) SendCard(ctx context.Context, message domain.Message, card string) error {
	reply := domain.Reply{
		ReceiveID:     message.MessageID,
		ReceiveIDType: "message_id",
		Content:       card,
		MsgType:       domain.MsgTypeInteractive,
	}
	return s.larkClient.ReplyMsg(ctx, reply)
}
