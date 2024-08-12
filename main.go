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
	"github.com/wawayes/lark-bot/services"
	"github.com/wawayes/lark-bot/tasks"
	qweather "github.com/wawayes/qweather-sdk-go"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 测试群chatID: oc_8f43329098cfc919213c3c87007f5409

func main() {
	r := gin.Default()

	// 初始化配置
	conf := initialization.GetConfig()
	// 飞书客户端初始化
	initialization.LoadLarkClient(*conf)
	larkClient := initialization.GetLarkClient()
	// 日志初始化
	l := logrus.New()
	// 和风天气客户端初始化
	qweatherClient := qweather.NewClient(conf.QWeather.Key)
	// 机器人助手初始化
	botHelper := handlers.NewBotHelper(map[string]string{"ou_16f6a982f4a0415201701fc2dd85ef8c": "及时雨大人"})
	// redis初始化
	initialization.InitRedis(*conf)
	redisClient := initialization.GetRedisClient()
	cache := services.NewUserCache(redisClient)
	locationService := services.NewLocationService(cache)

	// 初始化定时任务
	taskManager := tasks.NewTaskManager(l)
	dailyWeatherTask := tasks.NewDailyWeatherTask(
		conf.Tasks.DailyWeather,
		handlers.CreateQWeatherBotHandler(&handlers.BaseHandler{Logger: l, BotHelper: botHelper, LarkClient: larkClient}, qweatherClient, locationService),
		locationService,
	)
	taskManager.AddTask(dailyWeatherTask)
	weatherWarningTask := tasks.NewWeatherWarningTask(
		conf.Tasks.WeatherWarning,
		handlers.CreateQWeatherBotHandler(&handlers.BaseHandler{Logger: l, BotHelper: botHelper, LarkClient: larkClient, RedisClient: redisClient}, qweatherClient, locationService),
		locationService,
	)
	taskManager.AddTask(weatherWarningTask)
	taskManager.Start()
	defer taskManager.Stop()

	// 初始化处理器工厂
	factory := &handlers.HandlerFactory{
		Logger:          l,
		BotHelper:       botHelper,
		LarkClient:      larkClient,
		LocationService: locationService,
		QweatherClient:  qweatherClient,
		RedisClient:     redisClient,
	}
	bot := handlers.NewChatBot(
		factory,
		handlers.LoggingMiddleware(l),
	)

	// 当访问 "/ping" 路径时，返回 "pong" 字符串
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// 事件配置
	eventHandler := dispatcher.NewEventDispatcher(
		conf.Lark.VerificationToken, conf.Lark.EncryptToken).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			return bot.HandleMessage(ctx, event)
		}).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			return bot.ReadHandler(ctx, event)
		})
	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(eventHandler))

	// 监听并在 0.0.0.0:9000 上启动服务
	if err := initialization.StartServer(*conf, r); err != nil {
		l.Fatalf("failed to start server: %v", err)
	}
}
