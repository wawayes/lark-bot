package application

import (
	"context"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
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
	// 获取消息中艾特的机器人
	serviceField, err := adapters.WhichMentioned(ctx, h.adapters.Redis(), message)
	if err != nil {
		return err
	}

	// 根据不同的机器人执行不同的逻辑
	switch serviceField {
	case domain.ServiceFiledWeather:
		cardJson, err := h.buildWeatherCard()
		if err != nil {
			return err
		}
		return h.adapters.Lark().SendCardMsg(ctx, message.ChatID, cardJson)
	case domain.ServiceFiledLLM:
		return h.adapters.Lark().SendCardMsg(ctx, message.ChatID, "llm card json")
	case domain.ServiceFieldFlomo:
		return h.adapters.Lark().SendCardMsg(ctx, message.ChatID, "flomo card json")
	default:
		return nil
	}
}

func (h *HandleTextMessage) buildWeatherCard() (string, error) {
	card := adapters.NewLarkMessageCard()
	card.AddHeader("天气信息")
	card.AddTextElement("请选择您需要的天气信息类型：")

	btnDaily := card.AddButton("今日天气", map[string]interface{}{
		"todo": "today_weather",
	}, *larkcard.MessageCardButtonTypeDefault.Ptr())
	btnCurrent := card.AddButton("实时天气", map[string]interface{}{
		"todo": "current_weather",
	}, *larkcard.MessageCardButtonTypeDefault.Ptr())
	card.AddCardAction(larkcard.MessageCardActionLayoutTrisection.Ptr(), []larkcard.MessageCardActionElement{btnDaily, btnCurrent})
	card.AddTextElement("更多天气信息：")
	card.AddLinkElement("访问和风天气官网", "https://www.qweather.com/")
	cardJson, err := card.ToJson()
	if err != nil {
		return "", err
	}
	return cardJson, nil

}
