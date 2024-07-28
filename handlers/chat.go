package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/config"
	"github.com/wawayes/lark-bot/services"
	"github.com/wawayes/lark-bot/utils"
)

const (
	TEMPLATE_ID           = "AAq0BtEvPm8DZ"
	TEMPLATE_VERSION_NAME = "1.0.6"
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
	MaxTemperature string        `json:"max_temperature"`
	MinTemperature string        `json:"min_temperature"`
	WeatherInfos   []WeatherInfo `json:"weather_infos"`
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

// 发送天气卡片
func SendChatMsg() error {
	client := config.GetLarkClient()
	contentTemplate := "{\"type\": \"\", \"data\": { \"template_id\": \"\",\n\"template_version_name\": \"\",\"template_variable\": {} } }"
	var template ContentTemplate
	err := json.Unmarshal([]byte(contentTemplate), &template)
	if err != nil {
		l.Errorf("json marshal err: %s", err.Error())
		return err
	}

	//获取当日天气
	respDaily, err := services.GetWeatherDay(3)
	if err != nil {
		l.Errorf("get weather day err: %s", err.Error())
		return err
	}
	if len(respDaily.Daily) <= 0 {
		l.Errorf("get weather failed")
		return err
	}
	// 获取逐时天气
	respHourly, err := services.GetWeatherHourly()
	if err != nil {
		l.Errorf("get weather hourly failed")
		return err
	}
	if len(respHourly.Hourly) <= 0 {
		l.Errorf("get weather failed")
		return err
	}
	weatherInfos := make([]WeatherInfo, 0)
	isNextDay := false
	for _, v := range respHourly.Hourly {
		var tempIcon string
		tempNum, _ := strconv.Atoi(v.Temp)
		if tempNum >= 30 {
			tempIcon = "ANGRY"
		} else if tempNum >= 25 {
			tempIcon = "OBSESSED"
		} else if tempNum >= 15 {
			tempIcon = "BETRAYED"
		} else {
			tempIcon = "SKULL"
		}

		// 第二天
		if utils.ParseTime(v.FxTime) == "00:00" {
			isNextDay = true
		}

		weatherInfos = append(weatherInfos,
			WeatherInfo{
				Time:        fmt.Sprintf("%s%s", getPreviewTime(isNextDay), utils.ParseTime(v.FxTime)),
				Emoji:       tempIcon,
				Temperature: v.Temp,
				Weather:     fmt.Sprintf("%s %s", weatherToEmoji[v.Text], v.Text),
			})
	}
	template.Type = "template"
	template.Data.TemplateVariable.MaxTemperature = respDaily.Daily[0].TempMax // 最高温度
	template.Data.TemplateVariable.MinTemperature = respDaily.Daily[0].TempMin // 最低温度
	template.Data.TemplateID = TEMPLATE_ID                                     // 模版ID
	template.Data.TemplateVersionName = TEMPLATE_VERSION_NAME                  // 模板版本号
	template.Data.TemplateVariable.WeatherInfos = weatherInfos                 // 天气

	templateJson, err := json.Marshal(template)
	if err != nil {
		l.Errorf("json marshal err: %s", err.Error())
		return err
	}

	// 创建请求对象
	chats, err := GetChatList()
	if err != nil {
		l.Errorf("获取群ids失败: %s", err.Error())
		return err
	}
	// 遍历发送群
	for i := 0; i < len(chats); i++ {
		chatID := chats[i]
		// 正式群
		// if chatID == "oc_ee9a94ca81e2fbce54e739144392c266" {
		// 	continue
		// }
		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(`chat_id`).
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(chatID).
				MsgType(`interactive`).
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
