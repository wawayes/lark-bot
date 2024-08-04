package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type IndicesResponse struct {
    Code       string       `json:"code"`
    UpdateTime string       `json:"updateTime"`
    FxLink     string       `json:"fxLink"`
    Daily      []struct{
		Date     string `json:"date"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Level    string `json:"level"`
		Category string `json:"category"`
		Text     string `json:"text"`
	} `json:"daily"`
}

func (c *WeatherClient) GetDailyIndices() (*IndicesResponse, error) {
    url := fmt.Sprintf("%s/indices/1d?type=3,8&key=%s&location=%s", WeatherBaseURL, c.APIKey, c.LocationID)
    indices := &IndicesResponse{}
    err := c.GetWeather(url, indices)
    code := indices.Code
    if code != "200" {
        l.Errorf("get indices code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return indices, err
}