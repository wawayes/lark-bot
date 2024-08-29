package domain

import "github.com/wawayes/lark-bot/global"

type CardFactory interface {
	CreateCard() (string, *global.BasicError)
}

type MessageCard interface {
	ToJson() (string, error)
}
