package domain

import "context"

type Bot interface {
	HandleMessage(ctx context.Context, msg Message) error
	HandleAction(ctx context.Context, action Action) error
}

type WeatherBot interface {
	Bot
	GetDailyForecast(ctx context.Context, location string, days int) (*DailyForecast, error)
	GetCurrentWeather(ctx context.Context, location string, hours int) (*CurrentWeather, error)
	GetRainSnow(ctx context.Context, location string) (*RainSnow, error)
	GetWarningWeather(ctx context.Context, location string) (*WarningWeather, error)
}

type LLMBot interface {
	Bot
	GenerateResponse(ctx context.Context, prompt Message) (string, error)
}
