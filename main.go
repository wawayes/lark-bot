package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/robfig/cron/v3"
	"github.com/wawayes/lark-bot/config"
	"github.com/wawayes/lark-bot/handlers"
)

func main() {
	r := gin.Default()

	conf := config.GetConfig()

	// 当访问 "/ping" 路径时，返回 "pong" 字符串
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// 事件配置
	eventHandler := dispatcher.NewEventDispatcher(
		conf.VerificationToken, conf.EncryptToken).
		OnP2MessageReceiveV1(handlers.Handler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			return handlers.ReadHandler(ctx, event)
		})
	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(eventHandler))

	c := cron.New()
	c.AddFunc("45 6,10,14,18,22 * * *", func() { // 每天的6:45, 10:45, 14:45, 18:45 和 22:45 执行
		fmt.Println("Task executed at:", time.Now().Format("2006-01-02 15:04:05"))
		handlers.SendChatMsg()
	})
	c.Start()

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run("0.0.0.0:9000")
}
