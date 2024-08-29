package domain

type Geo struct {
	City   string
	FxLink string
}

type DailyForecast struct {
	City       string
	UpdateTime string  `json:"updateTime"`
	FxLink     string  `json:"fxLink"`
	Daily      []Daily `json:"daily"`
}

type Daily struct {
	FxDate       string `json:"fxDate"`
	Sunrise      string `json:"sunrise"`
	Sunset       string `json:"sunset"`
	TempMax      string `json:"tempMax"`
	TempMin      string `json:"tempMin"`
	TextDay      string `json:"textDay"`
	TextNight    string `json:"textNight"`
	WindSpeedDay string `json:"windSpeedDay"`
	Humidity     string `json:"humidity"`
	Precip       string `json:"precip"`
	Pressure     string `json:"pressure"`
}

type CurrentWeather struct {
	City     string
	Temp     string
	Text     string
	Humidity string
	ObsTime  string
	FxLink   string
}

type RainSnow struct {
	City    string
	Summary string
	FxLink  string
}

type WarningWeather struct {
	City    string
	Warning []Warning
	FxLink  string
}

type Warning struct {
	Sender   string
	PubTime  string
	Title    string
	Status   string
	Severity string
	Text     string
}
