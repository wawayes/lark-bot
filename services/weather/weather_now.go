package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type NowWeatherResponse struct {
    Code       string `json:"code"`
    UpdateTime string `json:"updateTime"`
    FxLink     string `json:"fxLink"`
    Now        struct {
        ObsTime   string `json:"obsTime"`
        FeelsLike string `json:"feelsLike"`
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
        Vis       string `json:"vis"`
        Cloud     string `json:"cloud"`
        Dew       string `json:"dew"`
    } `json:"now"`
}

func (c *WeatherClient) GetNowWeather() (*NowWeatherResponse, error) {
    url := fmt.Sprintf("%s/weather/now?key=%s&location=%s", WeatherBaseURL, c.APIKey, c.LocationID)
    weather := &NowWeatherResponse{}
    err := c.GetWeather(url, weather)
    code := weather.Code
    if code != "200" {
        l.Errorf("get now weather code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return weather, err
}