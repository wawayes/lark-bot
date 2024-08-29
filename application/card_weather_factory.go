package application

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/wawayes/lark-bot/global"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

type WeatherCardFactory struct{}

func (f *WeatherCardFactory) CreateCard() (string, *global.BasicError) {
	card := adapters.NewLarkMessageCard()
	card.AddHeader("天气信息", larkcard.TemplateBlue)
	card.AddTextElement("请选择您需要的天气信息类型：")
	card.AddCardAction(larkcard.MessageCardActionLayoutFlow.Ptr(), buildWeatherHomeBtns())
	cardJson, err := card.ToJson()
	if err != nil {
		return "", global.NewBasicError(global.CodeServerError, "create weather card error", cardJson, err)
	}
	return cardJson, nil
}


