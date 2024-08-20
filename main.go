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
	"github.com/wawayes/lark-bot/infrastructure"
	"github.com/wawayes/lark-bot/infrastructure/adapters"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 测试群chatID: oc_8f43329098cfc919213c3c87007f5409

func main() {
	r := gin.Default()

	// 初始化配置
	conf := infrastructure.GetConfig()
	adapter := adapters.NewAdapter(*conf)

	commandFactory := application.NewCommandFactory(adapter)

	// ... 设置事件处理器 ...
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

	// 监听并在 0.0.0.0:9000 上启动服务
	if err := infrastructure.StartServer(*conf, r); err != nil {
		fmt.Printf("failed to start server: %v", err)
	}
}
