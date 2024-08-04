package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/initialization"
	"github.com/wawayes/lark-bot/services/weather"
	"github.com/wawayes/lark-bot/utils"
)

const (
	WEATHER_TEMPLATE_ID          = "AAq0BtEvPm8DZ" // 每日天气模板
	WEATHER_WARINING_TEMPLATE_ID = "AAq0KrVFKi8kd" // 天气预警模板
	WEATHER_NOW_TEMPLATE_ID      = "AAq0elKwjDqMg" // 实时天气模板
)

type TemplateVariable struct {
	CityLocation      string   `json:"city_location"`       // 城市-地区
	WeatherText       string   `json:"weather_text"`        // 天气
	AirCondition      string   `json:"air_condition"`       // 空气质量
	MaxTemperature    string   `json:"max_temperature"`     // 最高温度
	MinTemperature    string   `json:"min_temperature"`     // 最低温度
	TomorrowWeather   string   `json:"tomorrow_weather"`    // 明日天气
	ComfText          string   `json:"comf_text"`           // 舒适度描述
	ComfLevel         string   `json:"comf_level"`          // 舒适度指数
	WearText          string   `json:"wear_text"`           // 穿衣描述
	WearLevel         string   `json:"wear_level"`          // 穿衣指数
	HomeTemp          string   `json:"home_temp"`           // 家的温度
	HomeText          string   `json:"home_text"`           // 家的天气
	SuZhouJieTemp     string   `json:"suzhoujie_temp"`      // 苏州街温度
	SuZhouJieText     string   `json:"suzhoujie_text"`      // 苏州街天气
	XiErQiTemp        string   `json:"xierqi_temp"`         // 西二旗温度
	XiErQiText        string   `json:"xierqi_text"`         // 西二旗天气
	WeatherWarning    string   `json:"weather_warning"`     // 天气预警信息
	NowTemperature    string   `json:"now_temperature"`     // 此时温度
	ObsTime           string   `json:"obs_time"`            // 观测时间
	Humidity          string   `json:"humidity"`            // 相对湿度
	FeelTemp          string   `json:"feel_temp"`           // 体感温度
	Vis               string   `json:"vis"`                 // 能见度
	WeatherUrl        ThirdUrl `json:"weather_url"`         // 天气URL
}

type ThirdUrl struct {
	Url        string `json:"url"`
	PcUrl      string `json:"pc_url"`
	IOSUrl     string `json:"ios_url"`
	AndroidUrl string `json:"android_url"`
}

type Data struct {
	TemplateID          string           `json:"template_id"`
	TemplateVersionName string           `json:"template_version_name"`
	TemplateVariable    TemplateVariable `json:"template_variable"`
}

type ContentTemplate struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

// 构建实时天气发送信息
func ConstructNowWeatherContent() (templateJson string, err error) {
	conf := initialization.GetConfig()
	client := weather.NewWeatherClient(conf.QWeather.Key, conf.QWeather.Location)
	now, err := client.GetNowWeather()
	if err != nil {
		l.Errorf("get weather now err: %s", err.Error())
		return
	}
	obsTime := utils.ParseTime(now.Now.ObsTime)
	template := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_NOW_TEMPLATE_ID,
			TemplateVariable: TemplateVariable{
				CityLocation:   "北京",
				WeatherText:    now.Now.Text,
				NowTemperature: now.Now.Temp,
				FeelTemp:       now.Now.FeelsLike,
				Humidity:       now.Now.Humidity,
				Vis:            now.Now.Vis,
				ObsTime:        obsTime,
				WeatherUrl:     ThirdUrl{Url: now.FxLink, PcUrl: now.FxLink, IOSUrl: now.FxLink, AndroidUrl: now.FxLink},
			},
		},
	}
	b, err := json.Marshal(template)
	if err != nil {
		l.Errorf("json marshal err: %s", err.Error())
		return "", err
	}
	return string(b), nil
}

// 构建天气预警信息
func ConstructWeatherWarningContent() (string, error) {
	conf := initialization.GetConfig()
	client := weather.NewWeatherClient(conf.QWeather.Key, conf.QWeather.Location)
	warning, err := client.GetWeatherWarning()
	if err != nil {
		l.Errorf("获取天气预警失败,err:%s", err.Error())
		return "", err
	}
	if len(warning.Warning) == 0 {
		l.Infof("暂无天气预警")
		return "", err
	}
	content := warning.Warning[0].Text
	// 构建content template
	template := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_WARINING_TEMPLATE_ID,
			// TemplateVersionName: TEMPLATE_VERSION_NAME,
			TemplateVariable: TemplateVariable{
				WeatherWarning:    content,
				WeatherUrl: ThirdUrl{Url: warning.FxLink, PcUrl: warning.FxLink, IOSUrl: warning.FxLink, AndroidUrl: warning.FxLink},
			},
		},
	}
	templateJson, err := json.Marshal(template)
	if err != nil {
		l.Errorf("json marshal err: %s", err.Error())
		return "", err
	}
	return string(templateJson), nil
}

// 回复分钟降雨量信息
func ReplyMinutelyWeather() (string, error) {
	conf := initialization.GetConfig()
	client := weather.NewWeatherClient(conf.QWeather.Key, conf.QWeather.Location)
	home, err := client.GetMinutelyWeather(weather.HomeLatitude)
	if err != nil {
		l.Errorf("get minutely weather err: %s", err.Error())
		return "", err
	}
	xierqi, _ := client.GetMinutelyWeather(weather.XierqiLatitude)
	suzhoujie, _ := client.GetMinutelyWeather(weather.SuzhoujieLatitude)
	res := fmt.Sprintf("\n五福家园, %s, %s \n\n苏州街, %s, %s\n\n西二旗, %s, %s",
		weather.HomeLatitude, home.Summary,
		weather.SuzhoujieLatitude, suzhoujie.Summary, 
		weather.XierqiLatitude, xierqi.Summary)
	return res, nil
}

// 构建今日天气卡片
func ConstructTodayWeatherContent(client *weather.WeatherClient) (string, error) {
	//获取今明两天天气
	respDaily, err := client.GetDailyWeather(3)
	if err != nil {
		l.Errorf("get weather day err: %s", err.Error())
		return "", err
	}
	// 获取空气质量
	airCondition, err := client.GetAirQuality()
	if err != nil {
		l.Errorf("get air condition err: %s", err.Error())
		return "", err
	}
	// 获取生活指数
	respLife, err := client.GetDailyIndices()
	if err != nil {
		l.Errorf("get life level err: %s", err.Error())
		return "", err
	}
	// 构建content template
	var template = &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_TEMPLATE_ID,
			//TemplateVersionName: TEMPLATE_VERSION_NAME,
			TemplateVariable: TemplateVariable{
				CityLocation:    "北京",
				WeatherText:     respDaily.Daily[0].TextDay,
				AirCondition:    airCondition.Now.Category,
				MaxTemperature:  respDaily.Daily[0].TempMax,
				MinTemperature:  respDaily.Daily[0].TempMin,
				TomorrowWeather: respDaily.Daily[1].TextDay,
				WearText:        respLife.Daily[0].Text,
				WearLevel:       respLife.Daily[0].Category,
				ComfText:        respLife.Daily[1].Text,
				ComfLevel:       respLife.Daily[1].Category,
				WeatherUrl:      ThirdUrl{Url: respDaily.FxLink, PcUrl: respDaily.FxLink, IOSUrl: respDaily.FxLink, AndroidUrl: respDaily.FxLink},
			},
		},
	}
	templateJson, err := json.Marshal(template)
	if err != nil {
		l.Errorf("json marshal err: %s", err.Error())
		return "", err
	}
	return string(templateJson), nil
}

func SendChatMsg(idType, id, msgType, templateJson string) error {
	// conf := initialization.GetConfig()
	// initialization.LoadLarkClient(*conf)
	client := initialization.GetLarkClient()
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(idType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(id).
			MsgType(msgType).
			Content(string(templateJson)).
			Build()).
		Build()
	resp, err := client.Im.Message.Create(context.Background(), req)
	if err != nil {
		l.Errorf("send msg err: %s", err.Error())
		return err
	}
	if !resp.Success() {
		l.Errorf("send msg error success is false")
		return err
	}

	return nil
}
