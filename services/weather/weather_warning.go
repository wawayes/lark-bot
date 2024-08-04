package weather

import "fmt"

type WarningResponse struct {
    Code       string           `json:"code"`
    UpdateTime string           `json:"updateTime"`
    FxLink     string           `json:"fxLink"`
    Warning    []struct{
		ID            string `json:"id"`
		Sender        string `json:"sender"`
		PubTime       string `json:"pubTime"`
		Title         string `json:"title"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
		Status        string `json:"status"`
		Level         string `json:"level"`
		Type          string `json:"type"`
		TypeName      string `json:"typeName"`
		Text          string `json:"text"`
		Related       string `json:"related"`
	} `json:"warning"`
}

func (c *WeatherClient) GetWeatherWarning() (*WarningResponse, error) {
    url := fmt.Sprintf("%s/warning/now?key=%s&location=%s", WeatherBaseURL, c.APIKey, c.LocationID)
    weather := &WarningResponse{}
    err := c.GetWeather(url, weather)
    return weather, err
}
