package handlers

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/initialization"
	"github.com/wawayes/lark-bot/services/weather"
)

// 定时向群聊发送逐时天气信息
func CronTaskRun() {
	conf := initialization.GetConfig()
	client := weather.NewWeatherClient(conf.QWeather.Key, conf.QWeather.Location)
	chatID := conf.Lark.AllChatID	
	idType := "chat_id"
	msg_type := "interactive"
	warningFlag := false
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() {
		fmt.Println("Weather warning check executed at:", time.Now().Format("2006-01-02 15:04:05"))
		templateJson, err := ConstructWeatherWarningContent()
		if err != nil {
			l.Errorf("构建天气发送内容失败 err: %s", err.Error())
		}
		if templateJson == "" || warningFlag {
			return
		}
		err = SendChatMsg(idType, chatID, msg_type, templateJson)
		if err != nil {
			l.Errorf("发送天气预警失败失败")
		}
		warningFlag = true
	})
	c.AddFunc("0 55 6 * * *", func() { 
		fmt.Println("Task executed at:", time.Now().Format("2006-01-02 15:04:05"))
		templateJson, err := ConstructTodayWeatherContent(client)
		if err != nil {
			l.Errorf("构建天气发送内容失败 err: %s", err.Error())
		}
		SendChatMsg(idType, chatID, msg_type, templateJson)
	})
	c.Start()
}