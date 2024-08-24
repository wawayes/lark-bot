/*
 * @Description:
 * @Author: wangmingyao@duxiaoman.com
 * @version:
 * @Date: 2024-07-28 07:03:02
 * @LastEditTime: 2024-07-31 01:31:45
 */
package main

import (
	"context"
	"fmt"

	"github.com/wawayes/lark-bot/application"
	"github.com/wawayes/lark-bot/global"
	"github.com/wawayes/lark-bot/infrastructure"
	"github.com/wawayes/lark-bot/infrastructure/adapters"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
	r := gin.Default()

	// 初始化配置
	conf := infrastructure.GetConfig()
	global.InitLogger(*conf)
	adapter := adapters.NewAdapter(*conf)

	application.InitBots(context.Background(), conf, adapter.Redis())

	commandFactory := application.NewCommandFactory(adapter)

	// 初始化 CardHandler
	weatherService := application.NewWeatherService("27df0ab1a3014458b59906f1c8bfa6f7")
	cardHandler := application.NewCardHandler(weatherService, adapter.Lark())

	// 设置事件处理器
	eventHandler := dispatcher.NewEventDispatcher(conf.Lark.VerificationToken, conf.Lark.EncryptToken).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			message := adapters.ConvertEventToMessage(event)
			command := commandFactory.CreateCommand(message.MsgType)
			if command == nil {
				return fmt.Errorf("unsupported message type: %s", message.MsgType)
			}
			return command.Execute(ctx, message)
		})

	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(eventHandler))
	r.POST("/webhook/card", sdkginext.NewCardActionHandlerFunc(
		larkcard.NewCardActionHandler(conf.Lark.VerificationToken, conf.Lark.EncryptToken,
			func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
				return cardHandler.Handle(ctx, cardAction)
			})))

	// 监听并在 0.0.0.0:9000 上启动服务
	if err := infrastructure.StartServer(*conf, r); err != nil {
		fmt.Printf("failed to start server: %v", err)
	}
}
