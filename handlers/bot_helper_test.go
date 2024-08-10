package handlers

import (
	"testing"

	"github.com/wawayes/lark-bot/initialization"
)

func TestSendMsg(t *testing.T) {
	var h BotHelper
	conf := initialization.GetConfig()
	initialization.LoadLarkClient(*conf)
	larkClient := initialization.GetLarkClient()
	resp, err := h.SendMsg("user_id", "276f826e", "text", `{"text": "hello"}`, larkClient)
	if err != nil || !resp.Success() {
		t.Errorf("SendMsg() failed, err: %v", err)
	}
	t.Logf("resp: %v", resp)
}
