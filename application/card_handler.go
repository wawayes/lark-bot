package application

import (
	"context"
	"encoding/json"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/wawayes/lark-bot/global"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
	"github.com/wawayes/lark-bot/services"
)

type CardHandler interface {
	Handle(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
}

type CardHandlerImpl struct {
	weatherService services.WeatherService
	adapter        adapters.Adapter
}

func NewCardHandler(weatherService services.WeatherService, adapter adapters.Adapter) *CardHandlerImpl {
	return &CardHandlerImpl{
		weatherService: weatherService,
		adapter:        adapter,
	}
}

func (h *CardHandlerImpl) Handle(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
	global.Log.Infof("receive card action: %+v", cardAction)
	cmd := NewWeatherCommand(h.weatherService, &h.adapter, cardAction)
	if cmd == nil {
		return createErrorResponse("Unknown action"), nil
	}

	// 如果是返回操作，发送首页卡片
	if cmd.IsBack() {
		cardFactory := &WeatherCardFactory{}
		cardJSON, err := cardFactory.CreateCard()
		if err != nil {
			global.Log.Errorf("failed to create card: %+v", err)
			return createErrorResponse("Failed to create response card"), nil
		}
		h.adapter.Lark().SendCardMsg(ctx, cardAction.OpenChatId, cardJSON)
		global.Log.Infof("send card message: %s", cardJSON)
		return createSuccessResponse(cardJSON), nil
	}

	responseCard, err := cmd.Execute(ctx)
	if err != nil {
		global.Log.Errorf("failed to execute command: %+v", err)
		return createErrorResponse("Failed to process request"), nil
	}

	cardJSON, errJson := responseCard.ToJson()
	if errJson != nil {
		global.Log.Errorf("failed to create card: %+v", errJson)
		return createErrorResponse("Failed to create response card"), nil
	}
	h.adapter.Lark().SendCardMsg(ctx, cardAction.OpenChatId, cardJSON)
	global.Log.Infof("send card message: %s", cardJSON)
	return createSuccessResponse(cardJSON), nil
}

func NewWeatherCommand(weatherService services.WeatherService, adapter *adapters.Adapter, cardAction *larkcard.CardAction) WeatherCommand {
	action := cardAction.Action.Value
	todo := action["todo"]
	switch todo {
	case "daily_weather", "today", "three_days", "seven_days", "back":
		return NewDailyWeatherCommand(weatherService, adapter, cardAction)
	case "current_weather":
		return NewCurrentWeatherCommand(weatherService, adapter, cardAction)
	case "rain_snow":
		return NewRainSnowCommand(weatherService, adapter, cardAction)
	case "warning_weather":
		return NewWarningWeatherCommand(weatherService, adapter, cardAction)
	default:
		return nil
	}
}

func buildWeatherHomeBtns() []larkcard.MessageCardActionElement {
	card := adapters.NewLarkMessageCard()
	btns := []larkcard.MessageCardActionElement{}
	btnDaily := card.AddButton("每日天气", map[string]interface{}{
		"todo":     "daily_weather",
		"nextCard": "daily_weather_card",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())
	btnCurrent := card.AddButton("实时天气", map[string]interface{}{
		"todo":     "current_weather",
		"nextCard": "current_weather_card",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())
	btnRainSnow := card.AddButton("雨雪提醒", map[string]interface{}{
		"todo":     "rain_snow",
		"nextCard": "rain_snow_card",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())
	btnWarningWeather := card.AddButton("预警天气", map[string]interface{}{
		"todo":     "warning_weather",
		"nextCard": "warning_weather_card",
	}, *larkcard.MessageCardButtonTypePrimary.Ptr())
	btns = append(btns, btnDaily, btnCurrent, btnRainSnow, btnWarningWeather)
	return btns
}

func createSuccessResponse(cardJSON string) map[string]interface{} {
	return map[string]interface{}{
		// "toast": map[string]interface{}{
		// 	"type":    "success",
		// 	"content": "操作成功",
		// },
		"card": map[string]interface{}{
			"type": "raw",
			"data": json.RawMessage(cardJSON),
		},
	}
}

func createErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"toast": map[string]interface{}{
			"type":    "error",
			"content": message,
		},
	}
}
