package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type GridWeather struct {
    Code       string `json:"code"`
    UpdateTime string `json:"updateTime"`
    FxLink     string `json:"fxLink"`
    Now        struct {
        Temp      string `json:"temp"`
        Icon      string `json:"icon"`
        Text      string `json:"text"`
        Wind360   string `json:"wind360"`
        WindDir   string `json:"windDir"`
        WindScale string `json:"windScale"`
        WindSpeed string `json:"windSpeed"`
        Humidity  string `json:"humidity"`
        Precip    string `json:"precip"`
        Pressure  string `json:"pressure"`
        Cloud     string `json:"cloud"`
    } `json:"now"`
}

func (c *WeatherClient) GetGridWeather(location string) (*GridWeather, error) {
    url := fmt.Sprintf("%s/grid-weather/now?key=%s&location=%s", WeatherBaseURL, c.APIKey, location)
    weather := &GridWeather{}
    err := c.GetWeather(url, weather)
    code := weather.Code
    if code != "200" {
        l.Errorf("get grid weather code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return weather, err
}
