package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type HourlyWeatherResponse struct {
    Code       string           `json:"code"`
    UpdateTime string           `json:"updateTime"`
    FxLink     string           `json:"fxLink"`
    Hourly     []struct{
        FxTime    string `json:"fxTime"`
        Temp      string `json:"temp"`
        Icon      string `json:"icon"`
        Text      string `json:"text"`
        Wind360   string `json:"wind360"`
        WindDir   string `json:"windDir"`
        WindScale string `json:"windScale"`
        WindSpeed string `json:"windSpeed"`
        Humidity  string `json:"humidity"`
        Pop       string `json:"pop"`
        Precip    string `json:"precip"`
        Pressure  string `json:"pressure"`
        Cloud     string `json:"cloud"`
        Dew       string `json:"dew"`
    } `json:"hourly"`
}

func (c *WeatherClient) GetHourlyWeather() (*HourlyWeatherResponse, error) {
    url := fmt.Sprintf("%s/weather/24h?key=%s&location=%s", WeatherBaseURL, c.APIKey, c.LocationID)
    weather := &HourlyWeatherResponse{}
    err := c.GetWeather(url, weather)
    code := weather.Code
    if code != "200" {
        l.Errorf("get hourly weather code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return weather, err
}