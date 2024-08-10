package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/wawayes/lark-bot/services"
)

type LocationMessageHandler struct {
	*BaseHandler
	LocationService *services.LocationService // 位置服务
}

func (f *HandlerFactory) CreateLocationHandler(base *BaseHandler, locationService *services.LocationService) *LocationMessageHandler {
	return &LocationMessageHandler{
		BaseHandler:     base,
		LocationService: locationService,
	}
}

func (h *LocationMessageHandler) Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	h.Logger.Infof("收到位置消息: %s", *event.Event.Message.Content)
	h.Logger.Infof("消息体: chatID: %s, chatType: %s", *event.Event.Message.ChatId, *event.Event.Message.ChatType)

	var content struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}
	if err := json.Unmarshal([]byte(*event.Event.Message.Content), &content); err != nil {
		h.Logger.Errorf("解析位置消息失败: %v", err)
		return err
	}

	chatID := *event.Event.Message.ChatId

	latitudeFloat, _ := strconv.ParseFloat(content.Latitude, 64)
	latitude := fmt.Sprintf("%.2f", latitudeFloat)
	longitudeFloat, _ := strconv.ParseFloat(content.Longitude, 64)
	longitude := fmt.Sprintf("%.2f", longitudeFloat)

	h.LocationService.SaveLocation(chatID, latitude, longitude)

	h.Logger.Infof("已保存群聊 %s 的位置: 纬度 %s, 经度 %s", chatID, content.Latitude, content.Longitude)

	// 可以在这里发送一条确认消息给用户
	resp, err := h.BotHelper.SendMsg("chat_id", chatID, "text", `{"text": "已记录您的位置，您现在可以查询天气了。"}`, h.LarkClient)
	if err != nil || !resp.Success() {
		h.Logger.Errorf("发送确认消息失败: %s", resp.CodeError.Msg)
		return err
	}
	return nil
}
