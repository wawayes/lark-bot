package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/wawayes/lark-bot/services"
	"github.com/wawayes/lark-bot/utils"
	qweather "github.com/wawayes/qweather-sdk-go"
)

const (
	WEATHER_TEMPLATE_ID          = "AAq0BtEvPm8DZ" // 每日天气模板
	WEATHER_WARINING_TEMPLATE_ID = "AAq0KrVFKi8kd" // 天气预警模板
	WEATHER_NOW_TEMPLATE_ID      = "AAq0elKwjDqMg" // 实时天气模板
	WEATHER_RAIN_TEMPLATE_ID     = "AAq0aPQGCCqzZ" // 降雨预报模板
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
	Content         string   `json:"content"`          // 文本内容
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
	cmd, location, err := h.ParseCmd(text)
	if err != nil {
		h.Logger.Errorf("ParseTextMsg error: %s", err.Error())
		return err
	}
	switch cmd {
	case "today":
		h.Logger.Info("handling today weather")
		err = h.handleTodayCommand(ctx, event, location)
	case "now":
		h.Logger.Info("handling now weather")
		err = h.handleNowCommand(ctx, event, location)
	case "warning":
		h.Logger.Info("handling warning weather")
		err = h.handleWarningCommand(ctx, event, location)
	case "rain":
		h.Logger.Info("handling rain weather")
		err = h.handleRainCommand(ctx, event, location)
	default:
		h.BotHelper.SendMsg("chat_id", *event.Event.Message.ChatId, "text", `{"text": "不支持的命令"}`, h.LarkClient)
	}
	if err != nil {
		h.Logger.Errorf("handle command error: %s", err.Error())
		return err
	}
	return nil
}

func (h *QWeatherBotHandler) handleTodayCommand(_ context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	openID := *event.Event.Message.ChatId
	lonla, cityName, err := h.getLonLaAndCityName(openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}
	// 获取未来三日天气预报
	daily, err := h.QWeatherClient.GetGridDailyWeather(lonla, 3)
	if err != nil {
		h.Logger.Errorf("GetDailyForecast error: %s", err.Error())
		return err
	}
	// 获取生活指数
	indicesType := []string{"3", "8"}
	indices, err := h.QWeatherClient.GetIndicesWeather(indicesType, lonla, 1)
	if err != nil {
		h.Logger.Errorf("GetIndicesWeather error: %s", err.Error())
		return err
	}
	// 获取空气质量
	air, err := h.QWeatherClient.GetAirQuality(lonla)
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
				CityLocation:    cityName,
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

func (h *QWeatherBotHandler) handleNowCommand(_ context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	openID := *event.Event.Message.ChatId
	lonla, cityName, err := h.getLonLaAndCityName(openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}
	// 获取实时天气
	now, err := h.QWeatherClient.GetCurrentWeather(lonla)
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
				CityLocation:   cityName,
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

func (h *QWeatherBotHandler) handleWarningCommand(_ context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	openID := *event.Event.Message.ChatId
	lonla, cityName, err := h.getLonLaAndCityName(openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}
	// 获取天气预警信息
	warning, err := h.QWeatherClient.GetWarningWeather(lonla)
	if err != nil {
		h.Logger.Errorf("GetWarningWeather error: %s", err.Error())
		return err
	}
	if len(warning.Warning) == 0 {
		h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "此地区暂无天气预警信息"}`, h.LarkClient)
		return nil
	}
	// 构建天气预警信息
	card := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_WARINING_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:   cityName,
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

func (h *QWeatherBotHandler) handleRainCommand(_ context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	openID := *event.Event.Message.ChatId
	lonla, cityName, err := h.getLonLaAndCityName(openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "查询的数据或地区不存在"}`, h.LarkClient)
		return err
	}
	// 获取降雨预报
	rain, err := h.QWeatherClient.GetMinutelyPrecipitation(lonla)
	if err != nil {
		h.Logger.Errorf("GetMinutelyPrecipitation error: %s", err.Error())
		return err
	}
	// 构建降雨预报信息
	card := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_RAIN_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation: cityName,
				Content:      rain.Summary,
				WeatherUrl:   ThirdUrl{PcUrl: rain.FxLink, AndroidUrl: rain.FxLink, IOSUrl: rain.FxLink},
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

func (h *QWeatherBotHandler) getLonLaAndCityName(openID, location string) (string, string, error) {
	var (
		lonla, cityName string
	)
	if location == "" {
		locationObj, exists := h.LocationService.GetLocation(openID)
		cityName = locationObj.Name
		lonla = fmt.Sprintf("%s,%s", locationObj.Longtitude, locationObj.Latitude)
		if !exists {
			h.Logger.Infof("Location not found for openID: %s", openID)
			resp, err := h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "没有找到您的位置信息，请先发送您的位置信息"}`, h.LarkClient)
			if err != nil || !resp.Success() {
				h.Logger.Errorf("SendMsg error: %s", err.Error())
				return "", "", err
			}
			return "", "", nil
		}
	} else {
		resp, err := h.QWeatherClient.CityLookup(location, "1")
		if err != nil {
			h.Logger.Errorf("CityLookup error: %s", err.Error())
			return "", "", err
		}
		cityName = resp.Location[0].Name
		lonla = fmt.Sprintf("%s,%s", resp.Location[0].Lon, resp.Location[0].Lat)
	}
	return lonla, cityName, nil
}

func (h *QWeatherBotHandler) ParseCmd(text string) (cmd, location string, err error) {
	content, err := h.BotHelper.ParseTextMsg(text)
	if err != nil {
		h.Logger.Errorf("ParseTextMsg error: %s", err.Error())
		return "", "", err
	}
	if !strings.Contains(content, "/") {
		h.Logger.Errorf("invalid command: %s", content)
		return content, "", nil
	}
	// 分割命令和位置
	parts := strings.SplitN(content, "/", 3)
	if len(parts) < 2 {
		h.Logger.Errorf("invalid command: %s", content)
		err = fmt.Errorf("invalid command: %s", content)
	}
	cmd = parts[1]
	if len(parts) == 3 {
		location = parts[2]
	}
	return cmd, location, nil
}
