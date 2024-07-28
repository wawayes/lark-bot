package handlers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawayes/lark-bot/config"
)

// oc_8f43329098cfc919213c3c87007f5409 测试群
// oc_ee9a94ca81e2fbce54e739144392c266 全员群
func TestGetChatList(t *testing.T) {
	conf := config.GetConfig()
	config.LoadLarkClient(*conf)
	chats, err := GetChatList()
	assert.Nil(t, err)
	fmt.Println(chats)
}

func TestSendChatMsg(t *testing.T) {
	conf := config.GetConfig()
	config.LoadLarkClient(*conf)
	err := SendChatMsg()
	assert.Nil(t, err)
}
