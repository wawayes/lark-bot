package tasks

import (
	"context"
	"fmt"

	"github.com/wawayes/lark-bot/handlers"
	"github.com/wawayes/lark-bot/services"
)

type WeatherWarningtTask struct {
	schedule        string
	qweatherHandler *handlers.QWeatherBotHandler
	locationService *services.LocationService
}

func NewWeatherWarningTask(schedule string, qweatherHandler *handlers.QWeatherBotHandler, locationService *services.LocationService) *WeatherWarningtTask {
	return &WeatherWarningtTask{
		schedule:        schedule,
		qweatherHandler: qweatherHandler,
		locationService: locationService,
	}
}

func (t *WeatherWarningtTask) Run() error {
	// 实现检查和发送天气预警的逻辑
	// 1. 获取所有用户的位置信息
	ctx := context.Background()
	locations := t.locationService.GetAllLocations(ctx)
	for openID, location := range locations {
		// 2. 使用 qweatherHandler 获取天气预警信息并发送
		t.qweatherHandler.SendWarningWeather(ctx, openID, fmt.Sprintf("%s,%s", location.Longtitude, location.Latitude))
	}
	return nil
}

func (t *WeatherWarningtTask) GetSchedule() string {
	return t.schedule
}
