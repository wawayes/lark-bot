package application

import (
	"context"

	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/global"
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

func (h *HandleTextMessage) Execute(ctx context.Context, message domain.Message) *global.BasicError {
	if len(message.Mentions) == 0 {
		return nil
	}
	// 获取消息中艾特的机器人
	serviceField, err := adapters.WhichMentioned(ctx, h.adapters.Redis(), message)
	if err != nil {
		return global.NewBasicError(global.CodeServerError, "get mentioned bot error", serviceField, err)
	}

	// 创建对应的卡片工厂
	var cardFactory domain.CardFactory
	switch serviceField {
	case domain.ServiceFiledWeather:
		cardFactory = &WeatherCardFactory{}
	case domain.ServiceFiledLLM:
		cardFactory = &LLMCardFactory{}
	case domain.ServiceFieldFlomo:
		cardFactory = &FlomoCardFactory{}
	default:
		cardFactory = &DefaultCardFactory{}
	}

	// 使用工厂创建卡片
	cardJson, err := cardFactory.CreateCard()
	if err != nil {
		return global.NewBasicError(global.CodeServerError, "create card error", cardJson, err)
	}

	// 发送卡片消息
	return h.adapters.Lark().SendCardMsg(ctx, message.ChatID, cardJson)
}

// LLM卡片工厂
type LLMCardFactory struct{}

func (f *LLMCardFactory) CreateCard() (string, *global.BasicError) {
	// TODO: 实现LLM卡片的创建逻辑
	return "llm card json", nil
}

// Flomo卡片工厂
type FlomoCardFactory struct{}

func (f *FlomoCardFactory) CreateCard() (string, *global.BasicError) {
	// TODO: 实现Flomo卡片的创建逻辑
	return "flomo card json", nil
}

// 默认卡片工厂
type DefaultCardFactory struct{}

func (f *DefaultCardFactory) CreateCard() (string, *global.BasicError) {
	return `{"text": "template card"}`, nil
}
