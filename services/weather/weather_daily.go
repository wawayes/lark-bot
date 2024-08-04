package weather

import (
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/global"
)

type DailyWeatherResponse struct {
    Code       string          `json:"code"`
    UpdateTime string          `json:"updateTime"`
    FxLink     string          `json:"fxLink"`
    Daily      []struct{
		FxDate        string `json:"fxDate"`
		Sunrise       string `json:"sunrise"`
		Sunset        string `json:"sunset"`
		Moonrise      string `json:"moonrise"`
		Moonset       string `json:"moonset"`
		MoonPhase     string `json:"moonPhase"`
		MoonPhaseIcon string `json:"moonPhaseIcon"`
		TempMax       string `json:"tempMax"`
		TempMin       string `json:"tempMin"`
		IconDay       string `json:"iconDay"`
		TextDay       string `json:"textDay"`
		IconNight     string `json:"iconNight"`
		TextNight     string `json:"textNight"`
		Wind360Day    string `json:"wind360Day"`
		WindDirDay    string `json:"windDirDay"`
		WindScaleDay  string `json:"windScaleDay"`
		WindSpeedDay  string `json:"windSpeedDay"`
		Wind360Night  string `json:"wind360Night"`
		WindDirNight  string `json:"windDirNight"`
		WindScaleNight string `json:"windScaleNight"`
		WindSpeedNight string `json:"windSpeedNight"`
		Humidity       string `json:"humidity"`
		Precip         string `json:"precip"`
		Pressure       string `json:"pressure"`
		Vis            string `json:"vis"`
		Cloud          string `json:"cloud"`
		UvIndex        string `json:"uvIndex"`
	} `json:"daily"`
}

func (c *WeatherClient) GetDailyWeather(days int) (*DailyWeatherResponse, error) {
    url := fmt.Sprintf("%s/weather/%dd?key=%s&location=%s", WeatherBaseURL, days, c.APIKey, c.LocationID)
    weather := &DailyWeatherResponse{}
    err := c.GetWeather(url, weather)
	code := weather.Code
    if code != "200" {
        l.Errorf("get daily weather code is not 200")
        return nil, global.NewBasicError(code, global.GetWeatherErrorMessage(code))
    }
    return weather, err
}
