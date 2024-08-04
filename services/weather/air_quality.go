package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type AirResponse struct {
    Code       string `json:"code"`
    UpdateTime string `json:"updateTime"`
    FxLink     string `json:"fxLink"`
    Now        struct {
        Category string `json:"category"`
    } `json:"now"`
}

func (c *WeatherClient) GetAirQuality() (*AirResponse, error) {
    url := fmt.Sprintf("%s/air/now?location=%s&key=%s", WeatherBaseURL, c.LocationID, c.APIKey)
    air := &AirResponse{}
    err := c.GetWeather(url, air)
    code := air.Code
    if code != "200" {
        l.Errorf("get air quality code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return air, err
}