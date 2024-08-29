package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/global"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

// 位置消息处理器
type HandleLocationMessage struct {
	adapter adapters.Adapter
}

func NewHandleLocationMessage(adapter adapters.Adapter) *HandleLocationMessage {
	return &HandleLocationMessage{
		adapter: adapter,
	}
}

func (h *HandleLocationMessage) Execute(ctx context.Context, message domain.Message) *global.BasicError {
	// TODO implement the business logic of HandleLocationMessage
	location, err := parseLocation(message.Content)
	if err != nil {
		global.Log.Errorf("failed to parse location: %+v", err)
		return global.NewBasicError(global.CodeServerError, "failed to parse location", nil, err)
	}
	locationBytes, err := json.Marshal(location)
	if err != nil {
		global.Log.Errorf("failed to marshal location: %+v", err)
		return global.NewBasicError(global.CodeServerError, "failed to marshal location", nil, err)
	}
	err = h.adapter.Redis().Client.Set(ctx, "location:"+message.Sender.SenderID.OpenID, locationBytes, 24*time.Hour).Err()
	if err != nil {
		global.Log.Errorf("failed to set location in Redis: %+v", err)
		return global.NewBasicError(global.CodeServerError, "failed to set location in Redis", nil, err)
	}

	// 构建 JSON 格式的文本内容
	textContent := map[string]string{
		"text": fmt.Sprintf("已保存位置信息🎉: %s\n经度: %s\n纬度: %s", location.Name, location.Longitude, location.Latitude),
	}
	jsonContent, err := json.Marshal(textContent)
	if err != nil {
		global.Log.Errorf("failed to marshal JSON content: %+v", err)
		return global.NewBasicError(global.CodeServerError, "failed to marshal JSON content", nil, err)
	}
	h.adapter.Lark().SendTextMsg(ctx, message.ChatID, string(jsonContent))
	return nil
}

func GetLocationByOpenID(ctx context.Context, redis *redis.Client, openID string) (*domain.Location, *global.BasicError) {
	locationJson := redis.Get(ctx, fmt.Sprintf("location:%s", openID)).Val()
	// 解析 JSON 格式的位置信息
	var location domain.Location
	err := json.Unmarshal([]byte(locationJson), &location)
	if err != nil {
		global.Log.Errorf("failed to unmarshal location JSON: %+v", err)
		return nil, global.NewBasicError(global.CodeServerError, "failed to unmarshal location JSON", nil, err)
	}
	return &location, nil
}

func parseLocation(content string) (*domain.Location, error) {
	// 移除字符串开头的 "Content:" 前缀
	jsonStr := strings.TrimPrefix(content, "Content:")

	// 解析 JSON 字符串为 map[string]interface{}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, global.NewBasicError(global.CodeServerError, "failed to unmarshal JSON", nil, err)
	}

	// 从 map 中提取所需的字段值
	name, ok := data["name"].(string)
	if !ok {
		name = ""
	}

	longitude, ok := data["longitude"].(string)
	if !ok {
		longitude = ""
	}

	latitude, ok := data["latitude"].(string)
	if !ok {
		latitude = ""
	}

	// 创建 Location 结构体并返回
	location := &domain.Location{
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}

	return location, nil
}
