package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/redis/go-redis/v9"
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

// 天气信息
type WeatherInfo struct {
	CityName      string
	DailyForecast qweather.GridDailyWeatherForecastResponse
	Indices       qweather.LifeIndexResponse
	AirQuality    qweather.AirQualityResponse
}

type NowWeatherInfo struct {
	CityName string
	Now      qweather.CurrentWeather
}

type WarningWeatherInfo struct {
	CityName string
	Warning  qweather.WarningWeatherResponse
}

type RainWeatherInfo struct {
	CityName string
	Rain     qweather.MinutelyPrecipitationResponse
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
		err = h.HandleTodayCommand(ctx, event, location)
	case "now":
		h.Logger.Info("handling now weather")
		err = h.HandleNowCommand(ctx, event, location)
	case "warning":
		h.Logger.Info("handling warning weather")
		err = h.HandleWarningCommand(ctx, event, location)
	case "rain":
		h.Logger.Info("handling rain weather")
		err = h.HandleRainCommand(ctx, event, location)
	default:
		h.BotHelper.SendMsg("chat_id", *event.Event.Message.ChatId, "text", `{"text": "不支持的命令"}`, h.LarkClient)
	}
	if err != nil {
		h.Logger.Errorf("handle command error: %s", err.Error())
		return err
	}
	return nil
}

func (h *QWeatherBotHandler) HandleTodayCommand(ctx context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	return h.SendTodayWeather(ctx, *event.Event.Message.ChatId, location)
}

func (h *QWeatherBotHandler) HandleNowCommand(ctx context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	return h.SendNowWeather(ctx, *event.Event.Message.ChatId, location)
}

func (h *QWeatherBotHandler) HandleWarningCommand(ctx context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	return h.SendWarningWeather(ctx, *event.Event.Message.ChatId, location)
}

func (h *QWeatherBotHandler) HandleRainCommand(ctx context.Context, event *larkim.P2MessageReceiveV1, location string) error {
	return h.SendRainWeather(ctx, *event.Event.Message.ChatId, location)
}

func (h *QWeatherBotHandler) getLonLaAndCityName(ctx context.Context, openID, location string) (string, string, error) {
	var (
		lonla, cityName string
	)
	if location == "" {
		locationObj, exists := h.LocationService.GetLocation(ctx, openID)
		if !exists {
			h.Logger.Infof("Location not found for openID: %s", openID)
			resp, err := h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "没有找到您的位置信息，请先发送您的位置信息"}`, h.LarkClient)
			if err != nil || !resp.Success() {
				h.Logger.Errorf("SendMsg error: %s", err.Error())
				return "", "", err
			}
			return "", "", nil
		}
		cityName = locationObj.Name
		lonla = fmt.Sprintf("%s,%s", locationObj.Longtitude, locationObj.Latitude)
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

func (h *QWeatherBotHandler) getDailyWeatherInfo(_ context.Context, lonla string) (*WeatherInfo, error) {
	daily, err := h.QWeatherClient.GetGridDailyWeather(lonla, 3)
	if err != nil {
		return nil, fmt.Errorf("GetDailyForecast error: %w", err)
	}

	indicesType := []string{"8", "3"}
	indices, err := h.QWeatherClient.GetIndicesWeather(indicesType, lonla, 1)
	if err != nil {
		return nil, fmt.Errorf("GetIndicesWeather error: %w", err)
	}

	air, err := h.QWeatherClient.GetAirQuality(lonla)
	if err != nil {
		return nil, fmt.Errorf("GetAirQuality error: %w", err)
	}

	return &WeatherInfo{
		DailyForecast: *daily,
		Indices:       *indices,
		AirQuality:    *air,
	}, nil
}

func (h *QWeatherBotHandler) getNowWeatherInfo(_ context.Context, lonla string) (*NowWeatherInfo, error) {
	now, err := h.QWeatherClient.GetCurrentWeather(lonla)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentWeather error: %w", err)
	}
	return &NowWeatherInfo{Now: *now}, nil
}

func (h *QWeatherBotHandler) getWarningWeatherInfo(_ context.Context, lonla string) (*WarningWeatherInfo, error) {
	warning, err := h.QWeatherClient.GetWarningWeather(lonla)
	if err != nil {
		return nil, fmt.Errorf("GetWarningWeather error: %w", err)
	}
	return &WarningWeatherInfo{Warning: *warning}, nil
}

func (h *QWeatherBotHandler) getRainWeatherInfo(_ context.Context, lonla string) (*RainWeatherInfo, error) {
	rain, err := h.QWeatherClient.GetMinutelyPrecipitation(lonla)
	if err != nil {
		return nil, fmt.Errorf("GetMinutelyPrecipitation error: %w", err)
	}
	return &RainWeatherInfo{Rain: *rain}, nil
}

// 构建天气卡片
func (h *QWeatherBotHandler) buildWeatherCard(info *WeatherInfo) *ContentTemplate {
	return &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:    info.CityName,
				MaxTemperature:  info.DailyForecast.Daily[0].TempMax,
				MinTemperature:  info.DailyForecast.Daily[0].TempMin,
				TomorrowWeather: info.DailyForecast.Daily[1].TextDay,
				WeatherText:     info.DailyForecast.Daily[0].TextDay,
				AirCondition:    info.AirQuality.AQI[0].Category,
				ComfText:        info.Indices.Daily[0].Text,
				ComfLevel:       info.Indices.Daily[0].Category,
				WearLevel:       info.Indices.Daily[1].Category,
				WearText:        info.Indices.Daily[1].Text,
				WeatherUrl:      ThirdUrl{PcUrl: info.DailyForecast.FxLink, AndroidUrl: info.DailyForecast.FxLink, IOSUrl: info.DailyForecast.FxLink},
			},
		},
	}
}

func (h *QWeatherBotHandler) buildNowWeatherCard(info *NowWeatherInfo) *ContentTemplate {
	obsTime := utils.ParseTime(info.Now.Now.ObsTime)
	return &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_NOW_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:   info.CityName,
				NowTemperature: info.Now.Now.Temp,
				WeatherText:    info.Now.Now.Text,
				ObsTime:        obsTime,
				Humidity:       info.Now.Now.Humidity,
				FeelTemp:       info.Now.Now.FeelsLike,
				Vis:            info.Now.Now.Vis,
				WeatherUrl:     ThirdUrl{PcUrl: info.Now.FxLink, AndroidUrl: info.Now.FxLink, IOSUrl: info.Now.FxLink},
			},
		},
	}
}

func (h *QWeatherBotHandler) buildWarningWeatherCard(info *WarningWeatherInfo) *ContentTemplate {
	return &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_WARINING_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:   info.CityName,
				WeatherWarning: info.Warning.Warning[0].Text,
				WeatherUrl:     ThirdUrl{PcUrl: info.Warning.FxLink, AndroidUrl: info.Warning.FxLink, IOSUrl: info.Warning.FxLink},
			},
		},
	}
}

func (h *QWeatherBotHandler) buildRainWeatherCard(info *RainWeatherInfo) *ContentTemplate {
	return &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_RAIN_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation: info.CityName,
				Content:      info.Rain.Summary,
				WeatherUrl:   ThirdUrl{PcUrl: info.Rain.FxLink, AndroidUrl: info.Rain.FxLink, IOSUrl: info.Rain.FxLink},
			},
		},
	}
}

func (h *QWeatherBotHandler) SendTodayWeather(ctx context.Context, openID, location string) error {
	lonla, cityName, err := h.getLonLaAndCityName(ctx, openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}

	weatherInfo, err := h.getDailyWeatherInfo(ctx, lonla)
	if err != nil {
		h.Logger.Errorf("getDailyWeatherInfo error: %s", err.Error())
		return err
	}
	weatherInfo.CityName = cityName

	card := h.buildWeatherCard(weatherInfo)

	err = h.sendWeatherCard(ctx, openID, card)
	if err != nil {
		h.Logger.Errorf("sendWeatherCard error: %s", err.Error())
		return err
	}

	return nil
}

func (h *QWeatherBotHandler) SendNowWeather(ctx context.Context, openID, location string) error {
	lonla, cityName, err := h.getLonLaAndCityName(ctx, openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}

	weatherInfo, err := h.getNowWeatherInfo(ctx, lonla)
	if err != nil {
		h.Logger.Errorf("getNowWeatherInfo error: %s", err.Error())
		return err
	}
	weatherInfo.CityName = cityName

	card := h.buildNowWeatherCard(weatherInfo)

	return h.sendWeatherCard(ctx, openID, card)
}

func (h *QWeatherBotHandler) SendWarningWeather(ctx context.Context, openID, location string) error {
	lonla, cityName, err := h.getLonLaAndCityName(ctx, openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}

	weatherInfo, err := h.getWarningWeatherInfo(ctx, lonla)
	if err != nil {
		h.Logger.Errorf("getWarningWeatherInfo error: %s", err.Error())
		return err
	}
	weatherInfo.CityName = cityName

	// 使用通用的哈希计算函数
	newHash, err := utils.CalculateHash(weatherInfo.Warning)
	if err != nil {
		h.Logger.Errorf("CalculateHash error: %s", err.Error())
		return err
	}

	// 从 Redis 获取旧的哈希值
	redisKey := "warning_hash:" + lonla
	oldHash, err := h.RedisClient.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		h.Logger.Errorf("Redis Get error: %s", err.Error())
		return err
	}

	// 如果有预警信息并且, 新旧哈希值不一致，发送预警消息, 如果没有预警信息或者新旧哈希值一致, 发送无预警信息的消息
	if len(weatherInfo.Warning.Warning) > 0 && newHash != oldHash {
		card := h.buildWarningWeatherCard(weatherInfo)
		err = h.sendWeatherCard(ctx, openID, card)
		if err != nil {
			h.Logger.Errorf("sendWeatherCard error: %s", err.Error())
			return err
		}
		// 更新哈希值
		err = h.RedisClient.Set(ctx, redisKey, newHash, 24*time.Hour).Err()
		if err != nil {
			h.Logger.Errorf("Redis Set error: %s", err.Error())
			return err
		}
	} else {
		resp, err := h.BotHelper.SendMsg("chat_id", openID, "text", `{"text": "暂无预警信息"}`, h.LarkClient)
		if err != nil || !resp.Success() {
			h.Logger.Errorf("SendMsg error: %s", err.Error())
			return err
		}
	}
	return nil
}

func (h *QWeatherBotHandler) SendRainWeather(ctx context.Context, openID, location string) error {
	lonla, cityName, err := h.getLonLaAndCityName(ctx, openID, location)
	if err != nil {
		h.Logger.Errorf("getLonLaAndCityName error: %s", err.Error())
		return err
	}

	weatherInfo, err := h.getRainWeatherInfo(ctx, lonla)
	if err != nil {
		h.Logger.Errorf("getRainWeatherInfo error: %s", err.Error())
		return err
	}
	weatherInfo.CityName = cityName

	card := h.buildRainWeatherCard(weatherInfo)

	return h.sendWeatherCard(ctx, openID, card)
}

func (h *QWeatherBotHandler) sendWeatherCard(_ context.Context, openID string, card *ContentTemplate) error {
	b, _ := json.Marshal(card)
	resp, err := h.BotHelper.SendMsg("chat_id", openID, "interactive", string(b), h.LarkClient)
	if err != nil || resp.Code != 0 {
		return fmt.Errorf("SendMsg error: %w, resp: %s", err, resp.Msg)
	}
	return nil
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
