package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/wawayes/lark-bot/services"
	"github.com/wawayes/lark-bot/utils"
	qweather "github.com/wawayes/qweather-sdk-go"
)

const (
	WEATHER_TEMPLATE_ID          = "AAq0BtEvPm8DZ" // 每日天气模板
	WEATHER_WARINING_TEMPLATE_ID = "AAq0KrVFKi8kd" // 天气预警模板
	WEATHER_NOW_TEMPLATE_ID      = "AAq0elKwjDqMg" // 实时天气模板
)

type QWeatherBotHandler struct {
	BaseHandler
	QWeatherClient  *qweather.Client
	LocationService *services.LocationService
}

func CreateQWeatherBotHandler(base *BaseHandler, qweather *qweather.Client, locationService *services.LocationService) *QWeatherBotHandler {
	return &QWeatherBotHandler{BaseHandler: *base, QWeatherClient: qweather, LocationService: locationService}
}

type ContentTemplate struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	TemplateID          string           `json:"template_id"`
	TemplateVersionName string           `json:"template_version_name"`
	TemplateVariable    TemplateVariable `json:"template_variable"`
}

type TemplateVariable struct {
	CityLocation    string   `json:"city_location"`    // 城市-地区
	WeatherText     string   `json:"weather_text"`     // 天气
	AirCondition    string   `json:"air_condition"`    // 空气质量
	MaxTemperature  string   `json:"max_temperature"`  // 最高温度
	MinTemperature  string   `json:"min_temperature"`  // 最低温度
	TomorrowWeather string   `json:"tomorrow_weather"` // 明日天气
	ComfText        string   `json:"comf_text"`        // 舒适度描述
	ComfLevel       string   `json:"comf_level"`       // 舒适度指数
	WearText        string   `json:"wear_text"`        // 穿衣描述
	WearLevel       string   `json:"wear_level"`       // 穿衣指数
	HomeTemp        string   `json:"home_temp"`        // 家的温度
	HomeText        string   `json:"home_text"`        // 家的天气
	SuZhouJieTemp   string   `json:"suzhoujie_temp"`   // 苏州街温度
	SuZhouJieText   string   `json:"suzhoujie_text"`   // 苏州街天气
	XiErQiTemp      string   `json:"xierqi_temp"`      // 西二旗温度
	XiErQiText      string   `json:"xierqi_text"`      // 西二旗天气
	WeatherWarning  string   `json:"weather_warning"`  // 天气预警信息
	NowTemperature  string   `json:"now_temperature"`  // 此时温度
	ObsTime         string   `json:"obs_time"`         // 观测时间
	Humidity        string   `json:"humidity"`         // 相对湿度
	FeelTemp        string   `json:"feel_temp"`        // 体感温度
	Vis             string   `json:"vis"`              // 能见度
	WeatherUrl      ThirdUrl `json:"weather_url"`      // 天气URL
}

type ThirdUrl struct {
	Url        string `json:"url"`
	PcUrl      string `json:"pc_url"`
	IOSUrl     string `json:"ios_url"`
	AndroidUrl string `json:"android_url"`
}

func (h *QWeatherBotHandler) Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	h.Logger.Info("QWeatherHandler is handling the message")
	text := *event.Event.Message.Content
	cmd, err := h.BotHelper.ParseTextMsg(text)
	if err != nil {
		h.Logger.Errorf("ParseTextMsg error: %s", err.Error())
		return err
	}
	switch cmd {
	case "/today":
		h.Logger.Info("handling today weather")
		err = h.handleTodayCommand(ctx, event)
	case "/now":
		h.Logger.Info("handling now weather")
		err = h.handleNowCommand(ctx, event)
	case "/warning":
		h.Logger.Info("handling warning weather")
		err = h.handleWarningCommand(ctx, event)
	default:
		h.BotHelper.SendMsg("chat_id", *event.Event.Message.ChatId, "text", `{"text": "不支持的命令"}`, h.LarkClient)
	}
	if err != nil {
		h.Logger.Errorf("handle command error: %s", err.Error())
		return err
	}
	return nil
}

func (h *QWeatherBotHandler) handleTodayCommand(_ context.Context, event *larkim.P2MessageReceiveV1) error {
	// 获取当前位置
	openID := *event.Event.Message.ChatId
	location, exists := h.LocationService.GetLocation(openID)
	lalon := fmt.Sprintf("%s,%s", location.Longtitude, location.Latitude)
	if !exists {
		h.Logger.Infof("Location not found for openID: %s", openID)
		resp, err := h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "没有找到您的位置信息，请先发送您的位置信息"}`, h.LarkClient)
		if err != nil || !resp.Success() {
			h.Logger.Errorf("SendMsg error: %s", err.Error())
			return err
		}
		return nil
	}
	// 查询城市信息
	city, err := h.QWeatherClient.CityLookup(lalon, "3")
	if err != nil {
		h.Logger.Errorf("CityLookup error: %s", err.Error())
		return err
	}
	locationID := city.Location[0].ID
	// 获取未来三日天气预报
	daily, err := h.QWeatherClient.GetDailyForecast(locationID, 3)
	if err != nil {
		h.Logger.Errorf("GetDailyForecast error: %s", err.Error())
		return err
	}
	// 获取生活指数
	indicesType := []string{"3", "8"}
	indices, err := h.QWeatherClient.GetIndicesWeather(indicesType, locationID, 1)
	if err != nil {
		h.Logger.Errorf("GetIndicesWeather error: %s", err.Error())
		return err
	}
	// 获取空气质量
	air, err := h.QWeatherClient.GetAirQuality(locationID)
	if err != nil {
		h.Logger.Errorf("GetAirQuality error: %s", err.Error())
		return err
	}
	// 构建今日天气预报信息
	card := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:    city.Location[0].Name,
				MaxTemperature:  daily.Daily[0].TempMax,
				MinTemperature:  daily.Daily[0].TempMin,
				TomorrowWeather: daily.Daily[1].TextDay,
				WeatherText:     daily.Daily[0].TextDay,
				AirCondition:    air.AQI[0].Category,
				ComfText:        indices.Daily[0].Text,
				ComfLevel:       indices.Daily[0].Category,
				WearLevel:       indices.Daily[1].Category,
				WearText:        indices.Daily[1].Text,
				WeatherUrl:      ThirdUrl{PcUrl: daily.FxLink, AndroidUrl: daily.FxLink, IOSUrl: daily.FxLink},
			},
		},
	}
	b, _ := json.Marshal(card)
	resp, err := h.BotHelper.SendMsg("chat_id", openID, "interactive", string(b), h.LarkClient)
	if err != nil || resp.Code != 0 {
		h.Logger.Errorf("SendMsg error: %s, resp: %s", err.Error(), resp.Msg)
		return err
	}
	return nil
}

func (h *QWeatherBotHandler) handleNowCommand(_ context.Context, event *larkim.P2MessageReceiveV1) error {
	// 获取当前位置
	openID := *event.Event.Message.ChatId
	location, exists := h.LocationService.GetLocation(openID)
	lalon := fmt.Sprintf("%s,%s", location.Longtitude, location.Latitude)
	if !exists {
		h.Logger.Infof("Location not found for openID: %s", openID)
		resp, err := h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "没有找到您的位置信息，请先发送您的位置信息"}`, h.LarkClient)
		if err != nil || !resp.Success() {
			h.Logger.Errorf("SendMsg error: %s", err.Error())
			return err
		}
		return nil
	}
	// 查询城市信息
	city, err := h.QWeatherClient.CityLookup(lalon, "3")
	if err != nil {
		h.Logger.Errorf("CityLookup error: %s", err.Error())
		return err
	}
	// 获取实时天气
	now, err := h.QWeatherClient.GetCurrentWeather(lalon)
	if err != nil {
		h.Logger.Errorf("GetNowWeather error: %s", err.Error())
		return err
	}
	// 格式化时间
	obsTime := utils.ParseTime(now.Now.ObsTime)
	// 构建实时天气信息
	card := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_NOW_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:   city.Location[0].Name,
				NowTemperature: now.Now.Temp,
				WeatherText:    now.Now.Text,
				ObsTime:        obsTime,
				Humidity:       now.Now.Humidity,
				FeelTemp:       now.Now.FeelsLike,
				Vis:            now.Now.Vis,
				WeatherUrl:     ThirdUrl{PcUrl: now.FxLink, AndroidUrl: now.FxLink, IOSUrl: now.FxLink},
			},
		},
	}
	b, _ := json.Marshal(card)
	resp, err := h.BotHelper.SendMsg("chat_id", openID, "interactive", string(b), h.LarkClient)
	if err != nil || !resp.Success() {
		h.Logger.Errorf("SendMsg error: %s", err.Error())
		return err
	}
	return nil
}

func (h *QWeatherBotHandler) handleWarningCommand(_ context.Context, larkim *larkim.P2MessageReceiveV1) error {
	// 获取当前位置
	openID := *larkim.Event.Message.ChatId
	location, exists := h.LocationService.GetLocation(openID)
	lalon := fmt.Sprintf("%s,%s", location.Longtitude, location.Latitude)
	if !exists {
		h.Logger.Infof("Location not found for openID: %s", openID)
		resp, err := h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "没有找到您的位置信息，请先发送您的位置信息"}`, h.LarkClient)
		if err != nil || !resp.Success() {
			h.Logger.Errorf("SendMsg error: %s", err.Error())
			return err
		}
		return nil
	}
	city, err := h.QWeatherClient.CityLookup(lalon, "3")
	if err != nil {
		h.Logger.Errorf("CityLookup error: %s", err.Error())
		return err
	}
	// 获取天气预警信息
	warning, err := h.QWeatherClient.GetWarningWeather(lalon)
	if err != nil {
		h.Logger.Errorf("GetWarningWeather error: %s", err.Error())
		return err
	}
	// 构建天气预警信息
	card := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_WARINING_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:   city.Location[0].Adm2,
				WeatherWarning: warning.Warning[0].Text,
			},
		},
	}
	b, _ := json.Marshal(card)
	resp, err := h.BotHelper.SendMsg("chat_id", openID, "interactive", string(b), h.LarkClient)
	if err != nil || !resp.Success() {
		h.Logger.Errorf("SendMsg error: %s", err.Error())
		return err
	}
	return nil
}
