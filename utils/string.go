package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wawayes/lark-bot/domain"
)

func ParseTime(input string) string {
	// 定义输入时间的格式
	const inputFormat = "2006-01-02T15:04-07:00"
	// 定义输出时间的格式
	const outputFormat = "2006年1月2日15:04"

	// 解析输入的时间字符串
	parsedTime, err := time.Parse(inputFormat, input)
	if err != nil {
		return ""
	}

	// 格式化解析后的时间
	formattedTime := parsedTime.Format(outputFormat)
	return formattedTime
}

// BuildReplyContent 将 Reply 结构体的内容转换为飞书 API 所需的 JSON 格式
func BuildMessageContent(reply domain.Reply) (string, error) {
	var content map[string]interface{}

	switch reply.MsgType {
	case domain.MsgTypeText:
		content = map[string]interface{}{
			"text": reply.Content,
		}
	case domain.MsgTypePost:
		content = map[string]interface{}{
			"post": reply.Content, // 假设 Content 字段已经是正确格式的富文本内容
		}
	// 可以添加更多类型的处理
	default:
		return "", fmt.Errorf("unsupported message type: %s", reply.MsgType)
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}

	return string(contentJSON), nil
}
