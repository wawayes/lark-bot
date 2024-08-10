package handlers

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type BotHelper struct {
	BotOpenIDs map[string]string // key: openID, value: botname
}

func NewBotHelper(botConfigs map[string]string) *BotHelper {
	return &BotHelper{
		BotOpenIDs: botConfigs,
	}
}

func (h *BotHelper) SendMsg(idType, id, msgType, templateJson string, client *lark.Client) (*larkim.CreateMessageResp, error) {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(idType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(id).
			MsgType(msgType).
			Content(string(templateJson)).
			Build()).
		Build()
	resp, err := client.Im.Message.Create(context.Background(), req)
	return resp, err
}

// 目前仅判断艾特一个人的情况
func (h *BotHelper) WhichBotMentioned(event *larkim.P2MessageReceiveV1) *string {
	if len(event.Event.Message.Mentions) == 0 {
		return nil
	}
	openID := event.Event.Message.Mentions[0].Id.OpenId
	if _, exists := h.BotOpenIDs[*openID]; exists {
		return openID
	}
	return nil
}

// 解析text消息内容
func (h *BotHelper) ParseTextMsg(text string) (content string, err error) {
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(text), &m)
	if err != nil {
		return
	}
	content = m["text"].(string)
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]+\s`)
	content = regex.ReplaceAllString(content, "")
	return
}

// 处理消息内容
func (h *BotHelper) ProcessMessage(msg interface{}) (string, error) {
	msg = strings.TrimSpace(msg.(string))
	msgB, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	msgStr := string(msgB)

	if len(msgStr) >= 2 {
		msgStr = msgStr[1 : len(msgStr)-1]
	}
	return msgStr, nil
}
