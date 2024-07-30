package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/robfig/cron/v3"
	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/config"
	"github.com/wawayes/lark-bot/services"
)

const (
	TEMPLATE_ID           = "AAq0BtEvPm8DZ"
	TEMPLATE_VERSION_NAME = "1.0.8"
)

var (
	weatherToEmoji = map[string]string{
		"晴":      "☀️", // 晴
		"多云":     "⛅",  // 部分多云
		"阴天":     "☁️", // 阴天
		"小雨":     "🌧️", // 雨
		"中雨":     "🌧️", // 雨
		"大雨":     "🌧️", // 雨
		"暴雨":     "🌧️", // 雨
		"雷阵雨":    "⛈️", // 雷雨
		"雷阵雨转暴雨": "⛈️", // 用雷雨代表
		"雪":      "❄️", // 雪
		"阴":      "☁️", // 阴天
		"小雪":     "🌨️", // 轻雪
		"中雪":     "🌨️", // 中雪
		"大雪":     "🌨️", // 大雪
		"暴雪":     "🌨️", // 暴雪
		"雾":      "🌫️", // 雾
		"霾":      "🌫️", // 雾霾
		"风":      "🌬️", // 风
		"台风":     "🌀",  // 台风
		"沙尘暴":    "🌪️", // 沙尘暴
	}
)

type WeatherInfo struct {
	Time        string `json:"time"`
	Temperature string `json:"temperature"`
	Weather     string `json:"weather"`
	Emoji       string `json:"emoji"`
}

type TemplateVariable struct {
	CityLocation    string     `json:"city_location"`    // 城市-地区
	TodayWeather    string     `json:"today_weather"`    // 今日天气
	AirCondition    string     `json:"air_condition"`    // 空气质量
	MaxTemperature  string     `json:"max_temperature"`  // 最高温度
	MinTemperature  string     `json:"min_temperature"`  // 最低温度
	TomorrowWeather string     `json:"tomorrow_weather"` // 明日天气
	ComfText        string     `json:"comf_text"`        // 舒适度描述
	ComfLevel       string     `json:"comf_level"`       // 舒适度指数
	WearText        string     `json:"wear_text"`        // 穿衣描述
	WearLevel       string     `json:"wear_level"`       // 穿衣指数
	HomeTemp        string     `json:"home_temp"`        // 家的温度
	HomeText        string     `json:"home_text"`        // 家的天气
	SuZhouJieTemp   string     `json:"suzhoujie_temp"`   // 苏州街温度
	SuZhouJieText   string     `json:"suzhoujie_text"`   // 苏州街天气
	XiErQiTemp      string     `json:"xierqi_temp"`      // 西二旗温度
	XiErQiText      string     `json:"xierqi_text"`      // 西二旗天气
	WeatherUrl      WeatherUrl `json:"weather_url"`
}

type WeatherUrl struct {
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

// 定时向群聊发送逐时天气信息
func CronTaskRun() {
	conf := config.GetConfig()
	testChatID := conf.TestChatID
	idType := "chat_id"
	msg_type := "interactive"
	templateJson, err := ConstructWeatherContent()
	if err != nil {
		l.Errorf("构建天气发送内容失败 err: %s", err.Error())
	}
	c := cron.New()
	c.AddFunc("45 6 * * *", func() { // 每天的6:45执行
		fmt.Println("Task executed at:", time.Now().Format("2006-01-02 15:04:05"))
		SendChatMsg(idType, testChatID, msg_type, templateJson)
	})
	c.Start()
}

// 根据温度匹配一个飞书表情
func GenTempIcon(temperature string) (tempIcon string) {
	tempNum, _ := strconv.Atoi(temperature)
	if tempNum >= 30 {
		tempIcon = "ANGRY"
	} else if tempNum >= 25 {
		tempIcon = "OBSESSED"
	} else if tempNum >= 15 {
		tempIcon = "BETRAYED"
	} else {
		tempIcon = "SKULL"
	}
	return
}

// 构建content发送消息的消息体
func ConstructWeatherContent() (string, error) {
	//获取今明两天天气
	respDaily, err := services.GetWeatherDay(3)
	if err != nil {
		l.Errorf("get weather day err: %s", err.Error())
		return "", err
	}
	// 获取空气质量
	airCondition, err := services.GetAirCondition()
	if err != nil {
		l.Errorf("get air condition err: %s", err.Error())
		return "", err
	}
	// 获取生活指数
	respLife, err := services.GetDaily()
	if err != nil {
		l.Errorf("get life level err: %s", err.Error())
		return "", err
	}
	// 获取格点天气
	gridResp1, _ := services.GetGridWeather(services.HOME_LOCATION)
	gridResp2, _ := services.GetGridWeather(services.SUZHOUJIE_LOCATION)
	gridResp3, _ := services.GetGridWeather(services.XIERQI_LOCATION)
	// 构建content template
	var template = &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID:          TEMPLATE_ID,
			TemplateVersionName: TEMPLATE_VERSION_NAME,
			TemplateVariable: TemplateVariable{
				CityLocation:    "北京",
				TodayWeather:    respDaily.Daily[0].TextDay,
				AirCondition:    airCondition,
				MaxTemperature:  respDaily.Daily[0].TempMax,
				MinTemperature:  respDaily.Daily[0].TempMin,
				TomorrowWeather: respDaily.Daily[1].TextDay,
				ComfText:        respLife.Daily[0].Text,
				ComfLevel:       respLife.Daily[0].Category,
				WearText:        respLife.Daily[1].Text,
				WearLevel:       respLife.Daily[1].Category,
				HomeTemp:        gridResp1.Now.Temp,
				HomeText:        gridResp1.Now.Text,
				SuZhouJieTemp:   gridResp2.Now.Temp,
				SuZhouJieText:   gridResp2.Now.Text,
				XiErQiTemp:      gridResp3.Now.Temp,
				XiErQiText:      gridResp3.Now.Text,
				WeatherUrl:      WeatherUrl{Url: services.WEATHER_URL, PcUrl: services.WEATHER_URL, IOSUrl: services.WEATHER_URL, AndroidUrl: services.WEATHER_URL},
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

// 发送群聊天气卡片
func SendChatMsg(idType, id, msgType, templateJson string) error {
	client := config.GetLarkClient()
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

// 获取用户或机器人所在群组ID
func GetChatList() ([]string, error) {
	client := config.GetLarkClient()
	// 创建请求对象
	req := larkim.NewListChatReqBuilder().
		UserIdType("user_id").
		SortType("ByCreateTimeAsc").
		PageSize(20).
		Build()
	l.Infof("client: %v", client)
	resp, err := client.Im.Chat.List(context.Background(), req)
	if err != nil {
		l.Errorf("获取机器人所在群聊列表失败, err: %s", err.Error())
		return nil, err
	}
	if resp.Code != 0 {
		l.Errorf("获取机器人所在群聊列表 响应体错误码非0, resp: %v", resp.CodeError.Msg)
		l.Debugf("get chat list resp: %v", resp.Data.Items)
	}
	chats := make([]string, 0)
	for _, v := range resp.Data.Items {
		chatID := v.ChatId
		external := v.External
		// 如果不是外部的群聊, 就能发消息
		if !*external {
			chats = append(chats, *chatID)
		}
	}
	return chats, nil
}

func getPreviewTime(isNextDay bool) string {
	if isNextDay {
		return "n."
	}
	return ""
}
