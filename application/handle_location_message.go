package application

import (
	"context"

	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

// 位置消息处理器
type HandleLocationMessage struct {
	adapter adapters.Adapter
}

func NewHandleLocationMessage(adapter adapters.Adapter) *HandleLocationMessage {
	return &HandleLocationMessage{
		adapter: adapter,
	}
}

func (h *HandleLocationMessage) Execute(ctx context.Context, message domain.Message) error {
	// TODO implement the business logic of HandleLocationMessage
	// message := convertEventToMessage(event)
	return nil
}

func extractLocation(message domain.Message) (domain.Location, error) {
	// TODO 从消息中提取位置信息的逻辑
	return domain.Location{}, nil
}

func (h *HandleLocationMessage) saveLocation(ctx context.Context, message domain.Message, location domain.Location) error {
	// TODO 保存位置信息的逻辑
	return nil
}

func (h *HandleLocationMessage) sendReply(ctx context.Context, message domain.Message, location domain.Location) error {
	// TODO 发送回复的逻辑
	return nil
}
