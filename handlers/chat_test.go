/*
 * @Description:
 * @Author: wangmingyao@duxiaoman.com
 * @version:
 * @Date: 2024-07-28 07:03:02
 * @LastEditTime: 2024-07-31 01:30:21
 */
package handlers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawayes/lark-bot/config"
)

func TestCronTaskRun(t *testing.T) {
	conf := config.InitTestConfig()
	config.LoadLarkClient(*conf)
	CronTaskRun()
}

// oc_8f43329098cfc919213c3c87007f5409 测试群
// oc_ee9a94ca81e2fbce54e739144392c266 全员群
func TestGetChatList(t *testing.T) {
	conf := config.InitTestConfig()
	config.LoadLarkClient(*conf)
	chats, err := GetChatList()
	assert.Nil(t, err)
	fmt.Println(chats)
}

func TestSendChatMsg(t *testing.T) {
	conf := config.InitTestConfig()
	config.LoadLarkClient(*conf)
	templateJson, err := ConstructWeatherContent()
	assert.Nil(t, err)
	var (
		idType  = "chat_id"
		id      = "oc_8f43329098cfc919213c3c87007f5409"
		msgType = "interactive"
	)
	err = SendChatMsg(idType, id, msgType, templateJson)
	assert.Nil(t, err)
}
