package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/initialization"
	"github.com/wawayes/lark-bot/services/weather"
)

func HandleChatMsg(ctx context.Context, event *larkim.P2MessageReceiveV1) (err error) {
	conf := initialization.GetConfig()
	client := weather.NewWeatherClient(conf.QWeather.Key, conf.QWeather.Location)
	msgType := event.Event.Message.ChatType
	msgID := event.Event.Message.MessageId
	if *msgType != "group" {
		l.Infof("暂时先处理群聊消息, 暂未考虑清晰合适的私聊场景.")
		return err
	}
	chatID := event.Event.Message.ChatId
	content, err := parseTextMsg(*event.Event.Message.Content)
	if err != nil {
		l.Errorf("解析 接收消息 错误, err: %s", err.Error())
		return err
	}
	switch content {
	case "/today":
		templateJson, err := ConstructTodayWeatherContent(client)
		if err != nil {
			l.Errorf("construct weather content err: %s", err.Error())
			return err
		}
		err = SendChatMsg("chat_id", *chatID, "interactive", templateJson)
		if err != nil {
			l.Errorf("Send chat msg err: %s", err.Error())
			return err
		}
	case "/now":
		templateJson, err := ConstructNowWeatherContent()
		if err != nil {
			l.Errorf("construct now weather err %s", err.Error())
			return err
		}
		err = SendChatMsg("chat_id", *chatID, "interactive", templateJson)
		if err != nil {
			l.Errorf("Send chat msg err: %s", err.Error())
			return err
		}
	case "/rain", "/snow":
		msg, err := ReplyMinutelyWeather()
		if err != nil {
			l.Errorf("reply minutely weather err: %s", err.Error())
			return err
		}
		err = replyMsg(ctx, msg, msgID)
		if err != nil {
			l.Errorf("Send chat msg err: %s", err.Error())
			return err
		}
	default:
		err := replyMsg(ctx, "我暂时理解不了. 你可以使用\n/today: 查看今日天气预报\n/now查询当前北京天气\n/rain或者/snow查询降水情况, 目前会根据您的家&苏州街&西二旗的坐标反馈降水信息.\n祝生活愉快,期待再次见面.", msgID)
		if err != nil {
			l.Errorf("Send chat msg err: %s", err.Error())
			return err
		}
	}
	return nil
}

// 解析text消息内容
func parseTextMsg(text string) (content string, err error) {
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
func processMessage(msg interface{}) (string, error) {
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

// 回复文本消息
func replyMsg(ctx context.Context, msg string, msgId *string) error {
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

// 获取用户或机器人所在群组ID
func GetChatList() ([]string, error) {
	client := initialization.GetLarkClient()
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
