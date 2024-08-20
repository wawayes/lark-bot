package application

import (
	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

type CommandFactory struct {
	adapter adapters.Adapter
}

func NewCommandFactory(adapter adapters.Adapter) *CommandFactory {
	return &CommandFactory{
		adapter: adapter,
	}
}

func (f *CommandFactory) CreateCommand(messageType domain.MessageType) Command {
	switch messageType {
	case domain.MsgTypeText:
		return NewHandleTextMessage(f.adapter)
	case domain.MsgTypeLocation:
		return NewHandleLocationMessage(f.adapter)
	default:
		return nil
	}
}
