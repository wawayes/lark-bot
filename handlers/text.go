package handlers

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/wawayes/lark-bot/services"
	qweather "github.com/wawayes/qweather-sdk-go"
)

// 消息处理器
type TextMessageHandler struct {
	BaseHandler
	QweatherClient  *qweather.Client          // 天气客户端
	LocationService *services.LocationService // 位置服务
	BotHandlers     map[string]BotHandler     // 机器人处理器
}

func (f *HandlerFactory) CreateTextHandler(base *BaseHandler, qweather *qweather.Client, locationService *services.LocationService) MessageHandler {
	botHandlers := map[string]BotHandler{
		"ou_16f6a982f4a0415201701fc2dd85ef8c": CreateQWeatherBotHandler(base, f.QweatherClient, f.LocationService), // 及时雨大人
	}
	return &TextMessageHandler{BaseHandler: *base, QweatherClient: qweather, LocationService: locationService, BotHandlers: botHandlers}
}

// TODO 不艾特时, 会空指针
func (h *TextMessageHandler) Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	h.Logger.Infof("进入到了消息处理器, 消息内容: %s", *event.Event.Message.Content)
	// 判断机器人是否被艾特
	botOpenID := h.BotHelper.WhichBotMentioned(event)
	if botOpenID == nil {
		h.Logger.Infof("未艾特机器人, 不处理")
		return nil
	}
	// 根据不同的机器人, 走不同的逻辑
	if handler, exists := h.BotHandlers[*botOpenID]; exists {
		return handler.Handle(ctx, event)
	} else {
		h.Logger.Infof("未找到对应的机器人处理器, 不处理")
		return nil
	}
}
