package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/redis/go-redis/v9"
	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/infrastructure"
)

type LarkClient struct {
	Client *lark.Client
}

func NewLarkClient(config infrastructure.Config) *LarkClient {
	options := []lark.ClientOptionFunc{
		lark.WithLogLevel(larkcore.LogLevelDebug),
	}
	if config.Lark.BaseUrl != "" {
		options = append(options, lark.WithOpenBaseUrl(config.Lark.BaseUrl))
	}

	client := lark.NewClient(config.Lark.AppID, config.Lark.AppSecret, options...)
	return &LarkClient{Client: client}
}

// 回复消息
func (lc *LarkClient) ReplyMsg(ctx context.Context, receiveID, msgType, contentJson string) error {
	resp, err := lc.Client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(receiveID).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(msgType).
			Uuid(uuid.New().String()).
			Content(contentJson).
			Build()).
		Build())
	if err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("failed to reply message: %s", resp.Msg)
	}
	return nil
}

func (lc *LarkClient) SendCardMsg(ctx context.Context, chatID, cardJson string) error {
	// TODO 发送消息卡片
	resp, err := lc.Client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(chatID).
			Content(cardJson).
			Build()).
		Build())

	if err != nil {

		return err
	}
	if !resp.Success() {
		return fmt.Errorf("failed to send card message: %s", resp.Msg)
	}

	return nil
}

func ConvertEventToMessage(event *larkim.P2MessageReceiveV1) domain.Message {
	// 因为可能空指针，创建一个辅助函数来安全地获取字符串值
	getString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	// 创建一个辅助函数来转换 Mentions
	convertMentions := func(mentions []*larkim.MentionEvent) []domain.Mention {
		var result []domain.Mention
		for _, m := range mentions {
			if m != nil {
				mention := domain.Mention{
					Key:  getString(m.Key),
					Name: getString(m.Name),
					ID: domain.IDObject{
						UserID:  getString(m.Id.UserId),
						UnionID: getString(m.Id.UnionId),
						OpenID:  getString(m.Id.OpenId),
					},
					TenantKey: getString(m.TenantKey),
				}
				result = append(result, mention)
			}
		}
		return result
	}

	return domain.Message{
		MessageID: getString(event.Event.Message.MessageId),
		RootID:    getString(event.Event.Message.RootId),
		ParentID:  getString(event.Event.Message.ParentId),
		MsgType:   domain.MessageType(getString(event.Event.Message.MessageType)),
		ChatID:    getString(event.Event.Message.ChatId),
		ChatType:  domain.ChatType(getString(event.Event.Message.ChatType)),
		Content:   getString(event.Event.Message.Content),
		Sender: domain.Sender{
			SenderID: domain.IDObject{
				UserID:  getString(event.Event.Sender.SenderId.UserId),
				UnionID: getString(event.Event.Sender.SenderId.UnionId),
				OpenID:  getString(event.Event.Sender.SenderId.OpenId),
			},
			SenderType: getString(event.Event.Sender.SenderType),
			TenantKey:  getString(event.Event.Sender.TenantKey),
		},
		CreateTime: getString(event.Event.Message.CreateTime),
		UpdateTime: getString(event.Event.Message.UpdateTime),
		ThreadID:   getString(event.Event.Message.ThreadId),
		Mentions:   convertMentions(event.Event.Message.Mentions),
		UserAgent:  getString(event.Event.Message.UserAgent),
	}
}

func WhichMentioned(ctx context.Context, redisClient *RedisClient, message domain.Message) (domain.ServiceFiled, error) {
	// 从 Redis 中获取机器人的服务字段
	serviceField, err := redisClient.Client.Get(ctx, fmt.Sprintf("bot:%s", message.Mentions[0].Name)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("bot not found")
	} else if err != nil {
		return "", fmt.Errorf("redis get error: %w", err)
	}
	return domain.ServiceFiled(serviceField), nil
}
