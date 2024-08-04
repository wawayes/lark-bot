package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type MinutelyResponse struct {
    Code       string `json:"code"`
    UpdateTime string `json:"updateTime"`
    FxLink     string `json:"fxLink"`
    Summary    string `json:"summary"`
}

// 获取分钟级降水
func (c *WeatherClient) GetMinutelyWeather(location string) (*MinutelyResponse, error) {
    url := fmt.Sprintf("%s/minutely/5m?location=%s&key=%s", WeatherBaseURL, location, c.APIKey)
    weather := &MinutelyResponse{}
    err := c.GetWeather(url, weather)
    code := weather.Code
    if code != "200" {
        l.Errorf("get minutely rain/snow code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return weather, err
}
