package tasks

import (
	"context"
	"fmt"

	"github.com/wawayes/lark-bot/handlers"
	"github.com/wawayes/lark-bot/services"
)

type DailyWeatherTask struct {
	schedule        string
	qweatherHandler *handlers.QWeatherBotHandler
	locationService *services.LocationService
}

func NewDailyWeatherTask(schedule string, qweatherHandler *handlers.QWeatherBotHandler, locationService *services.LocationService) *DailyWeatherTask {
	return &DailyWeatherTask{
		schedule:        schedule,
		qweatherHandler: qweatherHandler,
		locationService: locationService,
	}
}

func (t *DailyWeatherTask) Run() error {
	ctx := context.Background()
	locations := t.locationService.GetAllLocations(ctx)
	for openID, location := range locations {
		// 使用 qweatherHandler 获取天气信息并发送
		t.qweatherHandler.SendTodayWeather(ctx, openID, fmt.Sprintf("%s,%s", location.Longtitude, location.Latitude))
	}
	return nil
}

func (t *DailyWeatherTask) GetSchedule() string {
	return t.schedule
}
