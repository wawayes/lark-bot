package weather

import (
	"encoding/json"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/utils"
)

const (
    WeatherBaseURL = "https://devapi.qweather.com/v7"
	HomeLatitude = "116.27,40.15" // 家坐标
	SuzhoujieLatitude = "116.31,39.98"
	XierqiLatitude = "116.31,40.05"
)

type WeatherClient struct {
    APIKey     string
    LocationID string
}

func NewWeatherClient(apiKey, locationID string) *WeatherClient {
    return &WeatherClient{
        APIKey:     apiKey,
        LocationID: locationID,
    }
}

type WeatherResponse interface {

}

func (c *WeatherClient) GetWeather(url string, data WeatherResponse) error {
    body, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("http get weather request err: %s", err.Error())
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		l.Errorf("get weather json unmarshal err: %s", err.Error())
		return err
	}
	return nil
}

