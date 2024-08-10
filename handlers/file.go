package handlers

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 文件处理器
type FileMessageHandler struct {
	BaseHandler
}

func (f *HandlerFactory) CreateFileHandler() MessageHandler {
	return &FileMessageHandler{BaseHandler: BaseHandler{Logger: f.Logger}}
}

func (h *FileMessageHandler) Handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return nil
}
