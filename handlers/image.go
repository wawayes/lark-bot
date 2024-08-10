package handlers

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 图片处理器
type ImageMessageHandler struct {
	BaseHandler
}

func (f *HandlerFactory) CreateImageHandler() MessageHandler {
	return &ImageMessageHandler{BaseHandler: BaseHandler{Logger: f.Logger}}
}

func (h *ImageMessageHandler) Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	h.Logger.Infof("进入到了图片处理器, 消息内容: %s", *event.Event.Message.Content)
	return nil
}
