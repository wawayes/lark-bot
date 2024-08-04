/*
 * @Description:
 * @Author: wangmingyao@duxiaoman.com
 * @version:
 * @Date: 2024-07-28 07:03:02
 * @LastEditTime: 2024-07-31 01:30:21
 */
package handlers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawayes/lark-bot/initialization"
)

func TestCronTaskRun(t *testing.T) {
	path := "../config.yaml"
	conf := initialization.GetTestConfig(path)
	initialization.LoadLarkClient(*conf)
	CronTaskRun()
}


func TestSendWarningMsg(t *testing.T) {
	conf := initialization.GetTestConfig("../config.yaml")
	initialization.LoadLarkClient(*conf)
	var (
		idType  = "chat_id"
		id      = "oc_8f43329098cfc919213c3c87007f5409"
		msgType = "interactive"
	)
	// 构建content template
	template := &ContentTemplate{
		Type: "template",
		Data: Data{
			TemplateID: WEATHER_WARINING_TEMPLATE_ID,
			// TemplateVersionName: TEMPLATE_VERSION_NAME,
			TemplateVariable: TemplateVariable{
				WeatherWarning:    "海中心气象台2023年04月03日10时30分发布大风蓝色预警[Ⅳ级/一般]：受江淮气旋影响，预计明天傍晚以前本市大部地区将出现6级阵风7-8级的东南大风，沿江沿海地区7级阵风8-9级，请注意防范大风对高空作业、交通出行、设施农业等的不利影响。",
				WeatherUrl: ThirdUrl{
					Url: "https://www.qweather.com/weather/haidian-101010100.html", 
					PcUrl: "https://www.qweather.com/weather/haidian-101010100.html",
					IOSUrl: "https://www.qweather.com/weather/haidian-101010100.html", 
					AndroidUrl: "https://www.qweather.com/weather/haidian-101010100.html"},
			},
		},
	}
	b, err := json.Marshal(template)
	assert.Nil(t, err)
	templateJson := string(b)
	err = SendChatMsg(idType, id, msgType, templateJson)
	assert.Nil(t, err)
}

// // oc_8f43329098cfc919213c3c87007f5409 测试群
// // oc_ee9a94ca81e2fbce54e739144392c266 全员群
// func TestGetChatList(t *testing.T) {
// 	conf := initialization.InitTestConfig()
// 	initialization.LoadLarkClient(*conf)
// 	chats, err := GetChatList()
// 	assert.Nil(t, err)
// 	fmt.Println(chats)
// }

// func TestSendChatMsg(t *testing.T) {
// 	conf := initialization.InitTestConfig()
// 	initialization.LoadLarkClient(*conf)
// 	templateJson, err := ConstructTodayWeatherContent()
// 	assert.Nil(t, err)
// 	var (
// 		idType  = "chat_id"
// 		id      = "oc_8f43329098cfc919213c3c87007f5409"
// 		msgType = "interactive"
// 	)
// 	err = SendChatMsg(idType, id, msgType, templateJson)
// 	assert.Nil(t, err)
// }
