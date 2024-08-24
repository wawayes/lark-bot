package application

import (
	"context"
	"fmt"

	"github.com/wawayes/lark-bot/domain"
	"github.com/wawayes/lark-bot/global"
	qweather "github.com/wawayes/qweather-sdk-go"
)

type WeatherService interface {
	GetDailyForecast(ctx context.Context, location string, days int) (*domain.DailyForecast, *global.BasicError)
	GetCurrentWeather(ctx context.Context, location string) (*domain.CurrentWeather, *global.BasicError)
	GetRainSnow(ctx context.Context, location string) (*domain.RainSnow, *global.BasicError)
	GetWarningWeather(ctx context.Context, location string) (*domain.WarningWeather, *global.BasicError)
}

type WeatherServiceImpl struct {
	sdk *qweather.Client
}

func NewWeatherService(apiKey string) *WeatherServiceImpl {
	client := qweather.NewClient(apiKey)
	return &WeatherServiceImpl{
		sdk: client,
	}
}

func (s *WeatherServiceImpl) GetDailyForecast(ctx context.Context, location string, days int) (*domain.DailyForecast, *global.BasicError) {
	geo, err := s.sdk.CityLookup(location)
	if err != nil {
		global.Log.Errorf("failed to get geo: %+v, err: %+v", geo, err)
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	dailyForecast, err := s.sdk.GetDailyForecast(location, days)
	if err != nil {
		global.Log.Errorf("failed to get daily forecast: %+v, err: %+v", dailyForecast, err)
		return nil, global.NewBasicError(dailyForecast.Code, fmt.Sprintf("failed to get daily forecast: %+v", dailyForecast), nil, err)
	}
	return &domain.DailyForecast{
		City:      geo.Location[0].Name,
		FxDate:    dailyForecast.Daily[0].FxDate,
		TempMax:   dailyForecast.Daily[0].TempMax,
		TempMin:   dailyForecast.Daily[0].TempMin,
		TextDay:   dailyForecast.Daily[0].TextDay,
		TextNight: dailyForecast.Daily[0].TextNight,
		Humidity:  dailyForecast.Daily[0].Humidity,
		WindSpeed: dailyForecast.Daily[0].WindSpeedDay,
	}, &global.BasicError{}
}

func (s *WeatherServiceImpl) GetCurrentWeather(ctx context.Context, location string) (*domain.CurrentWeather, *global.BasicError) {
	geo, err := s.sdk.CityLookup(location)
	if err != nil {
		global.Log.Errorf("failed to get geo: %+v, err: %+v", geo, err)
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	currentWeather, err := s.sdk.GetCurrentWeather(location)
	if err != nil {
		global.Log.Errorf("failed to get current weather: %+v, err: %+v", currentWeather, err)
		return nil, global.NewBasicError(currentWeather.Code, fmt.Sprintf("failed to get current weather: %+v", currentWeather), nil, err)
	}
	return &domain.CurrentWeather{
		City:     geo.Location[0].Name,
		Temp:     currentWeather.Now.Temp,
		Text:     currentWeather.Now.Text,
		Humidity: currentWeather.Now.Humidity,
		ObsTime:  currentWeather.Now.ObsTime,
	}, &global.BasicError{}
}

func (s *WeatherServiceImpl) GetRainSnow(ctx context.Context, location string) (*domain.RainSnow, *global.BasicError) {
	geo, err := s.sdk.CityLookup(location)
	if err != nil {
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	rainSnow, err := s.sdk.GetMinutelyPrecipitation(location)
	if err != nil {
		return nil, global.NewBasicError(rainSnow.Code, fmt.Sprintf("failed to get rain snow: %+v", rainSnow), nil, err)
	}
	return &domain.RainSnow{
		City:    geo.Location[0].Name,
		Summary: rainSnow.Summary,
	}, &global.BasicError{}
}

func (s *WeatherServiceImpl) GetWarningWeather(ctx context.Context, location string) (*domain.WarningWeather, *global.BasicError) {
	geo, err := s.sdk.CityLookup(location)
	if err != nil {
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	warningWeather, err := s.sdk.GetWarningWeather(location)
	if err != nil {
		return nil, global.NewBasicError(warningWeather.Code, fmt.Sprintf("failed to get warning weather: %+v", warningWeather), nil, err)
	}
	return &domain.WarningWeather{
		City:    geo.Location[0].Name,
		Summary: warningWeather.Warning[0].Text,
	}, &global.BasicError{}
}
