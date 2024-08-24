package application

import (
	"context"
	"encoding/json"
	"fmt"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/wawayes/lark-bot/global"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

type CardHandler interface {
	Handle(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, *global.BasicError)
}

type CardHandlerImpl struct {
	weatherService WeatherService
	// llmService LLMService
	larkClient *adapters.LarkClient
}

func NewCardHandler(weatherService WeatherService, larkClient *adapters.LarkClient) *CardHandlerImpl {
	return &CardHandlerImpl{
		weatherService: weatherService,
		larkClient:     larkClient,
	}
}

func (h *CardHandlerImpl) Handle(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
	action := cardAction.Action.Value

	var responseCard *adapters.LarkMessageCard

	todo := action["todo"]
	switch todo {
	case "today_weather":
		weather, err := h.weatherService.GetDailyForecast(ctx, "101010700", 3)
		if !err.Ok() {
			global.Log.Errorf("failed to get today's weather: %+v, err: %+v", weather, err)
			return createErrorResponse("Failed to get today's weather"), nil
		}
		responseCard = createWeatherCard("Today's Weather", fmt.Sprintf("City: %s\nDate: %s\nMax Temp: %s\nMin Temp: %s\nDay: %s\nNight: %s\nHumidity: %s\nWind Speed: %s", weather.City, weather.FxDate, weather.TempMax, weather.TempMin, weather.TextDay, weather.TextNight, weather.Humidity, weather.WindSpeed))
	case "current_weather":
		forecast, err := h.weatherService.GetCurrentWeather(ctx, "101010700")
		if !err.Ok() {
			global.Log.Errorf("failed to get weather forecast: %+v, err: %+v", forecast, err)
			return createErrorResponse("Failed to get weather forecast"), nil
		}
		responseCard = createWeatherCard("Weather Forecast", fmt.Sprintf("City: %s\nObsTime: %s\nTemp: %s\nText: %s\nHumidity: %s\nFxLink: %s", forecast.City, forecast.ObsTime, forecast.Temp, forecast.Text, forecast.Humidity, forecast.FxLink))
	default:
		return createErrorResponse("Unknown action"), nil
	}

	cardJSON, err := responseCard.ToJson()
	if err != nil {
		return createErrorResponse("Failed to create response card"), nil
	}
	h.larkClient.SendCardMsg(ctx, cardAction.OpenChatId, cardJSON)
	return createSuccessResponse(cardJSON), nil
}

func createSuccessResponse(cardJSON string) map[string]interface{} {
	return map[string]interface{}{
		"toast": map[string]interface{}{
			"type":    "success",
			"content": "操作成功",
		},
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

func createWeatherCard(title string, weather string) *adapters.LarkMessageCard {
	card := adapters.NewLarkMessageCard()
	card.AddHeader(title)
	card.AddTextElement(weather)
	return card
}

func (h *CardHandlerImpl) HandleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	// This method is not used in this scenario, but kept for compatibility
	return nil
}
