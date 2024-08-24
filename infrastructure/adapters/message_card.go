package adapters

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/wawayes/lark-bot/domain"
)

type LarkMessageCard struct {
	card *larkcard.MessageCard
}

func NewLarkMessageCard() *LarkMessageCard {
	return &LarkMessageCard{card: larkcard.NewMessageCard()}
}

func (l *LarkMessageCard) ToJson() (string, error) {
	return l.card.String()
}

func (l *LarkMessageCard) AddHeader(title string) *LarkMessageCard {
	l.card.Header(larkcard.NewMessageCardHeader().
		Title(larkcard.NewMessageCardPlainText().Content(title)))
	return l
}

func (l *LarkMessageCard) AddTextElement(content string) *LarkMessageCard {
	l.card.Elements(append(l.card.Elements_,
		larkcard.NewMessageCardDiv().
			Text(larkcard.NewMessageCardPlainText().Content(content))))
	return l
}

func (l *LarkMessageCard) AddLinkElement(text, url string) *LarkMessageCard {
	l.card.Elements(append(l.card.Elements_,
		larkcard.NewMessageCardDiv().
			Text(larkcard.NewMessageCardPlainText().Content(text)).
			Extra(larkcard.NewMessageCardEmbedButton().
				Text(larkcard.NewMessageCardPlainText().Content("点击访问")).
				Type(larkcard.MessageCardButtonTypeDefault).
				Url(url))))
	return l
}

func (l *LarkMessageCard) AddCardAction(layout *larkcard.MessageCardActionLayout, actions []larkcard.MessageCardActionElement) *LarkMessageCard {
	l.card.Elements(append(l.card.Elements_,
		larkcard.NewMessageCardAction().
			Actions(actions).
			Layout(layout)))
	return l
}

func (l *LarkMessageCard) AddButton(content string, value map[string]interface{}, typename larkcard.MessageCardButtonType) *larkcard.MessageCardEmbedButton {
	btn := larkcard.NewMessageCardEmbedButton().
		Type(typename).
		Value(value).
		Text(larkcard.NewMessageCardPlainText().Content(content))
	return btn
}

var _ domain.MessageCard = (*LarkMessageCard)(nil)
