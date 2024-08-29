package application

import (
	"context"
	"fmt"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/global"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
	"github.com/wawayes/lark-bot/services"
	"github.com/wawayes/lark-bot/utils"
)

type WeatherCommand interface {
	Execute(ctx context.Context) (*adapters.LarkMessageCard, *global.BasicError)
	IsBack() bool
}

type DailyWeatherCommand struct {
	weatherService services.WeatherService
	adapter        adapters.Adapter
	cardAction     *larkcard.CardAction
}

func NewDailyWeatherCommand(weatherService services.WeatherService, adapter *adapters.Adapter, cardAction *larkcard.CardAction) *DailyWeatherCommand {
	return &DailyWeatherCommand{
		weatherService: weatherService,
		adapter:        *adapter,
		cardAction:     cardAction,
	}
}

func (c *DailyWeatherCommand) Execute(ctx context.Context) (*adapters.LarkMessageCard, *global.BasicError) {
	responseCard := adapters.NewLarkMessageCard()
	responseCard.AddHeader("每日天气", larkcard.TemplateBlue)

	todo := c.cardAction.Action.Value["todo"]
	openID := c.cardAction.OpenID
	location, err := GetLocationByOpenID(ctx, c.adapter.Redis().Client, openID)
	if err != nil {
		global.Log.Errorf("failed to get location by openID: %s, err: %+v", openID, err)
		return nil, global.NewBasicError(global.CodeServerError, "failed to get location by openID", nil, err)
	}
	lonla := fmt.Sprintf("%s,%s", location.Longitude, location.Latitude)
	switch todo {
	case "today":
		weather, err := c.weatherService.GetDailyForecast(ctx, lonla, 3)
		if !err.Ok() {
			global.Log.Errorf("failed to get today's weather: %+v, err: %+v", weather, err)
			return nil, err
		}
		responseCard = c.buildTodayWeatherCard(responseCard, weather)
	case "three_days":
		weather, err := c.weatherService.GetDailyForecast(ctx, lonla, 3)
		if !err.Ok() {
			global.Log.Errorf("failed to get three days' weather: %+v, err: %+v", weather, err)
			return nil, err
		}
		responseCard = c.buildThreeDaysWeatherCard(responseCard, weather)
	case "seven_days":
		weather, err := c.weatherService.GetDailyForecast(ctx, lonla, 7)
		if !err.Ok() {
			global.Log.Errorf("failed to get seven days' weather: %+v, err: %+v", weather, err)
			return nil, err
		}
		responseCard = c.buildSevenDaysWeatherCard(responseCard, weather)
	case "back":
		return nil, nil
	default:
		responseCard.AddTextElement("请选择你想查看的天气信息:")
	}

	responseCard.AddCardAction(larkcard.MessageCardActionLayoutFlow.Ptr(), c.buildDailyWeatherBtns())

	return responseCard, nil
}

func (c *DailyWeatherCommand) IsBack() bool {
	todo := c.cardAction.Action.Value["todo"]
	return todo == "back"
}

func (c *DailyWeatherCommand) buildDailyWeatherBtns() []larkcard.MessageCardActionElement {
	card := adapters.NewLarkMessageCard()
	btns := []larkcard.MessageCardActionElement{}

	btnToday := card.AddButton("今日天气", map[string]interface{}{
		"todo": "today",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())

	btnThreeDays := card.AddButton("未来三日天气", map[string]interface{}{
		"todo": "three_days",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())

	btnSevenDays := card.AddButton("未来七日天气", map[string]interface{}{
		"todo": "seven_days",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())

	btnBack := card.AddButton("返回", map[string]interface{}{
		"todo": "back",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())

	btns = append(btns, btnToday, btnThreeDays, btnSevenDays, btnBack)
	return btns
}

func (c *DailyWeatherCommand) buildDailyWeatherCard(card *adapters.LarkMessageCard, dailyForecast *domain.DailyForecast, days int) *adapters.LarkMessageCard {
	title := fmt.Sprintf("未来%d日天气", days)
	if days == 1 {
		title = "今日天气"
	}
	card.AddHeader(title, larkcard.TemplateBlue)

	cityInfo := fmt.Sprintf("**地区**: %s 🏙️\n", dailyForecast.City)
	card.AddLarkMd(cityInfo)
	for i := 0; i < days; i++ {
		dayInfo := fmt.Sprintf("**日期**: %s 📅\n**最高温度**: %s°C 🌡️\n**最低温度**: %s°C 🌡️\n**白天天气**: %s %s\n**夜间天气**: %s %s\n**湿度**: %s%% 💧\n**风速**: %sm/s 🍃\n",
			dailyForecast.Daily[i].FxDate,
			dailyForecast.Daily[i].TempMax,
			dailyForecast.Daily[i].TempMin,
			dailyForecast.Daily[i].TextDay, getWeatherEmoji(dailyForecast.Daily[i].TextDay),
			dailyForecast.Daily[i].TextNight, getWeatherEmoji(dailyForecast.Daily[i].TextNight),
			dailyForecast.Daily[i].Humidity,
			dailyForecast.Daily[i].WindSpeedDay,
		)
		card.AddLarkMd(dayInfo)
		card.AddHr()
	}

	return card
}

func (c *DailyWeatherCommand) buildTodayWeatherCard(card *adapters.LarkMessageCard, dailyForecast *domain.DailyForecast) *adapters.LarkMessageCard {
	return c.buildDailyWeatherCard(card, dailyForecast, 1)
}

func (c *DailyWeatherCommand) buildThreeDaysWeatherCard(card *adapters.LarkMessageCard, dailyForecast *domain.DailyForecast) *adapters.LarkMessageCard {
	return c.buildDailyWeatherCard(card, dailyForecast, 3)
}

func (c *DailyWeatherCommand) buildSevenDaysWeatherCard(card *adapters.LarkMessageCard, dailyForecast *domain.DailyForecast) *adapters.LarkMessageCard {
	return c.buildDailyWeatherCard(card, dailyForecast, 7)
}

type CurrentWeatherCommand struct {
	weatherService services.WeatherService
	adapter        adapters.Adapter
	cardAction     *larkcard.CardAction
}

func NewCurrentWeatherCommand(weatherService services.WeatherService, adapter *adapters.Adapter, cardAction *larkcard.CardAction) *CurrentWeatherCommand {
	return &CurrentWeatherCommand{
		weatherService: weatherService,
		adapter:        *adapter,
		cardAction:     cardAction,
	}
}

func (c *CurrentWeatherCommand) Execute(ctx context.Context) (*adapters.LarkMessageCard, *global.BasicError) {
	responseCard := adapters.NewLarkMessageCard()
	responseCard.AddHeader("实时天气", larkcard.TemplateBlue)
	location, err := GetLocationByOpenID(ctx, c.adapter.Redis().Client, c.cardAction.OpenID)
	if err != nil {
		global.Log.Errorf("failed to get location by openID: %s, err: %+v", c.cardAction.OpenID, err)
		return nil, global.NewBasicError(global.CodeServerError, "failed to get location by openID", nil, err)
	}
	lonla := fmt.Sprintf("%s,%s", location.Longitude, location.Latitude)
	forecast, err := c.weatherService.GetCurrentWeather(ctx, lonla)
	if !err.Ok() {
		global.Log.Errorf("failed to get weather forecast: %+v, err: %+v", forecast, err)
		return nil, err
	}
	responseCard = c.buildCurrentWeatherCard(responseCard, forecast)
	// TODO FxLink
	return responseCard, nil
}

func (c *CurrentWeatherCommand) IsBack() bool {
	return false
}

func (c *CurrentWeatherCommand) buildCurrentWeatherCard(card *adapters.LarkMessageCard, current *domain.CurrentWeather) *adapters.LarkMessageCard {
	cityInfo := fmt.Sprintf("**地区**: %s 🏙️\n", current.City)
	card.AddLarkMd(cityInfo)

	weatherInfo := fmt.Sprintf("**天气**: %s %s\n**温度**: %s°C 🌡️\n**湿度**: %s%% 💧\n**观测时间**: %s 🕙\n",
		current.Text, getWeatherEmoji(current.Text),
		current.Temp,
		current.Humidity,
		utils.ParseTime(current.ObsTime),
	)
	card.AddLarkMd(weatherInfo)
	card.AddCardAction(larkcard.MessageCardActionLayoutFlow.Ptr(), buildWeatherHomeBtns())

	return card
}

type RainSnowCommand struct {
	weatherService services.WeatherService
	adapter        adapters.Adapter
	cardAction     *larkcard.CardAction
}

func NewRainSnowCommand(weatherService services.WeatherService, adapter *adapters.Adapter, cardAction *larkcard.CardAction) *RainSnowCommand {
	return &RainSnowCommand{
		weatherService: weatherService,
		adapter:        *adapter,
		cardAction:     cardAction,
	}
}

func (c *RainSnowCommand) Execute(ctx context.Context) (*adapters.LarkMessageCard, *global.BasicError) {
	responseCard := adapters.NewLarkMessageCard()
	responseCard.AddHeader("雨雪查询", larkcard.TemplateBlue)
	location, err := GetLocationByOpenID(ctx, c.adapter.Redis().Client, c.cardAction.OpenID)
	if err != nil {
		global.Log.Errorf("failed to get location by openID: %s, err: %+v", c.cardAction.OpenID, err)
		return nil, global.NewBasicError(global.CodeServerError, "failed to get location by openID", nil, err)
	}
	lonla := fmt.Sprintf("%s,%s", location.Longitude, location.Latitude)
	rainSnow, err := c.weatherService.GetRainSnow(ctx, lonla)
	if !err.Ok() {
		global.Log.Errorf("failed to get rain snow: %+v, err: %+v", rainSnow, err)
		return nil, err
	}
	responseCard = buildRainSnowCard(responseCard, rainSnow)
	return responseCard, nil
}

func (c *RainSnowCommand) IsBack() bool {
	return false
}

func buildRainSnowCard(card *adapters.LarkMessageCard, rainSnow *domain.RainSnow) *adapters.LarkMessageCard {
	city := fmt.Sprintf("**地区**: %s 🏙️\n", rainSnow.City)
	card.AddLarkMd(city)
	card.AddHr()
	card.AddLarkMd(fmt.Sprintf("☔️❄️: %s", rainSnow.Summary))
	card.AddCardAction(larkcard.MessageCardActionLayoutFlow.Ptr(), buildWeatherHomeBtns())
	return card
}

type WarningWeatherCommand struct {
	weatherService services.WeatherService
	adapter        adapters.Adapter
	cardAction     *larkcard.CardAction
}

func NewWarningWeatherCommand(weatherService services.WeatherService, adapter *adapters.Adapter, cardAction *larkcard.CardAction) *WarningWeatherCommand {
	return &WarningWeatherCommand{
		weatherService: weatherService,
		adapter:        *adapter,
		cardAction:     cardAction,
	}
}

func (c *WarningWeatherCommand) Execute(ctx context.Context) (*adapters.LarkMessageCard, *global.BasicError) {
	responseCard := adapters.NewLarkMessageCard()
	responseCard.AddHeader("预警信息", larkcard.TemplateBlue)
	location, err := GetLocationByOpenID(ctx, c.adapter.Redis().Client, c.cardAction.OpenID)
	if err != nil {
		global.Log.Errorf("failed to get location by openID: %s, err: %+v", c.cardAction.OpenID, err)
		return nil, global.NewBasicError(global.CodeServerError, "failed to get location by openID", nil, err)
	}
	lonla := fmt.Sprintf("%s,%s", location.Longitude, location.Latitude)
	warning, err := c.weatherService.GetWarningWeather(ctx, lonla)
	if !err.Ok() {
		global.Log.Errorf("failed to get warning weather: %+v, err: %+v", warning, err)
		return nil, err
	}
	responseCard = buildWarningWeatherCard(responseCard, warning)
	return responseCard, nil
}

func (c *WarningWeatherCommand) IsBack() bool {
	return false
}

func buildWarningWeatherCard(card *adapters.LarkMessageCard, warning *domain.WarningWeather) *adapters.LarkMessageCard {
	city := fmt.Sprintf("**地区**: %s 🏙️\n", warning.City)
	card.AddLarkMd(city)
	for _, v := range warning.Warning {
		warningInfo := fmt.Sprintf("**发布机构**: %s\n**发布时间**: %s\n**标题**: %s\n**状态**: %s\n**严重程度**: %s\n**详细信息**: %s\n",
			v.Sender,
			utils.ParseTime(v.PubTime),
			v.Title,
			v.Status,
			v.Severity,
			v.Text,
		)
		card.AddLarkMd(warningInfo)
		card.AddHr()
	}
	if len(warning.Warning) == 0 {
		card.AddLarkMd("暂无预警信息")
	}
	card.AddCardAction(larkcard.MessageCardActionLayoutFlow.Ptr(), buildWeatherHomeBtns())
	return card
}

var weatherEmoji = map[string]string{
	"晴":   "🌞",
	"多云":  "⛅",
	"阴":   "☁️",
	"阵雨":  "🌦️",
	"雷阵雨": "⛈️",
	"小雨":  "🌧️",
	"中雨":  "🌧️",
	"大雨":  "🌧️",
	"暴雨":  "🌧️",
	"小雪":  "🌨️",
	"中雪":  "🌨️",
	"大雪":  "🌨️",
	"暴雪":  "🌨️",
	"雾":   "🌫️",
	"霾":   "🌫️",
}

func getWeatherEmoji(weather string) string {
	if emoji, ok := weatherEmoji[weather]; ok {
		return emoji
	}
	return ""
}
