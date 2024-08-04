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

	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/wawayes/lark-bot/handlers"
	"github.com/wawayes/lark-bot/initialization"

	"github.com/gin-gonic/gin"
	l "github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	conf := initialization.GetConfig()
	initialization.LoadLarkClient(*conf)

	// 当访问 "/ping" 路径时，返回 "pong" 字符串
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// 事件配置
	eventHandler := dispatcher.NewEventDispatcher(
		conf.Lark.VerificationToken, conf.Lark.EncryptToken).
		OnP2MessageReceiveV1(handlers.ReceiveHandler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			return handlers.ReadHandler(ctx, event)
		})
	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(eventHandler))

	handlers.CronTaskRun()

	// 监听并在 0.0.0.0:9000 上启动服务
	if err := initialization.StartServer(*conf, r); err != nil {
		l.Fatalf("failed to start server: %v", err)
	}
}
