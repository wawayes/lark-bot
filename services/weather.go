package services

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
	return &WeatherServiceImpl{
		sdk: qweather.NewClient(apiKey),
	}
}

func (s *WeatherServiceImpl) GetDailyForecast(ctx context.Context, location string, days int) (*domain.DailyForecast, *global.BasicError) {
	geo, err := s.sdk.CityLookup(location)
	if err != nil {
		global.Log.Errorf("failed to get geo: %+v, err: %+v", geo, err)
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	dailyForecastResp, err := s.sdk.GetDailyForecast(geo.Location[0].ID, days)
	if err != nil {
		global.Log.Errorf("failed to get daily forecast: %+v, err: %+v", dailyForecastResp, err)
		return nil, global.NewBasicError(dailyForecastResp.Code, fmt.Sprintf("failed to get daily forecast: %+v", dailyForecastResp), nil, err)
	}
	dailyForecast := &domain.DailyForecast{}
	dailyForecast.City = geo.Location[0].Name
	for _, v := range dailyForecastResp.Daily {
		dailyForecast.Daily = append(dailyForecast.Daily, domain.Daily{
			FxDate:       v.FxDate,
			TempMax:      v.TempMax,
			TempMin:      v.TempMin,
			TextDay:      v.TextDay,
			TextNight:    v.TextNight,
			Humidity:     v.Humidity,
			WindSpeedDay: v.WindSpeedDay,
			Precip:       v.Precip,
		})
	}
	return dailyForecast, &global.BasicError{}
}

func (s *WeatherServiceImpl) GetCurrentWeather(ctx context.Context, location string) (*domain.CurrentWeather, *global.BasicError) {
	geo, err := s.sdk.CityLookup(location)
	if err != nil {
		global.Log.Errorf("failed to get geo: %+v, err: %+v", geo, err)
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	currentWeather, err := s.sdk.GetCurrentWeather(geo.Location[0].ID)
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
		global.Log.Errorf("failed to get geo: %+v", geo)
		return nil, global.NewBasicError(geo.Code, fmt.Sprintf("failed to get geo: %+v", geo), nil, err)
	}
	warningWeather, err := s.sdk.GetWarningWeather(geo.Location[0].ID)
	if err != nil {
		global.Log.Errorf("failed to get warning weather: %+v", warningWeather)
		return nil, global.NewBasicError(warningWeather.Code, fmt.Sprintf("failed to get warning weather: %+v", warningWeather), nil, err)
	}
	warningList := make([]domain.Warning, 0)
	for _, v := range warningWeather.Warning {
		warningList = append(warningList, domain.Warning{
			Sender:   v.Sender,
			PubTime:  v.PubTime,
			Title:    v.Title,
			Status:   v.Status,
			Severity: v.Severity,
			Text:     v.Text,
		})
	}
	return &domain.WarningWeather{
		City:    geo.Location[0].Name,
		Warning: warningList,
	}, &global.BasicError{}
}
