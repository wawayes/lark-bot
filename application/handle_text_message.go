package application

import (
	"context"

	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

// 文本消息处理器
type HandleTextMessage struct {
	adapters adapters.Adapter
}

func NewHandleTextMessage(adapter adapters.Adapter) *HandleTextMessage {
	return &HandleTextMessage{
		adapters: adapter,
	}
}

func (h *HandleTextMessage) Execute(ctx context.Context, message domain.Message) error {
	// TODO implement the business logic of HandleTextMessage
	// message := convertEventToMessage(event)
	return nil
}

func buildMessageCard() (string, error) {
	// TODO 构建消息卡片的逻辑（与原来的逻辑相同）
	return "", nil
}
