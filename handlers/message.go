package handlers

import (
	"context"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/services"
	qweather "github.com/wawayes/qweather-sdk-go"
)

// MessageHandler is the interface that wraps the Handle method.
type MessageHandler interface {
	Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}

// 通用处理器
type BaseHandler struct {
	Logger      *logrus.Logger        // 日志
	BotHelper   *BotHelper            // 机器人相关
	LarkClient  *lark.Client          // 飞书客户端
	BotHandlers map[string]BotHandler // 机器人处理器
	RedisClient *redis.Client         // redis客户端
}

// 工厂模式 创建消息处理器
type HandlerFactory struct {
	Logger          *logrus.Logger
	BotHelper       *BotHelper
	LarkClient      *lark.Client
	LocationService *services.LocationService
	QweatherClient  *qweather.Client
	BotHandlers     map[string]BotHandler
	RedisClient     *redis.Client
}

func (f *HandlerFactory) CreateHandler(messageType string) MessageHandler {
	switch messageType {
	case "text":
		return f.CreateTextHandler(&BaseHandler{Logger: f.Logger, BotHelper: f.BotHelper, LarkClient: f.LarkClient, BotHandlers: f.BotHandlers, RedisClient: f.RedisClient}, f.QweatherClient, f.LocationService)
	case "image":
		return &ImageMessageHandler{BaseHandler: BaseHandler{Logger: f.Logger}}
	case "file":
		return &FileMessageHandler{BaseHandler: BaseHandler{Logger: f.Logger}}
	case "location":
		return f.CreateLocationHandler(&BaseHandler{Logger: f.Logger, BotHelper: f.BotHelper, LarkClient: f.LarkClient}, f.LocationService)
	default:
		return nil
	}
}
