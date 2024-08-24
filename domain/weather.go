package domain

type Geo struct {
	City   string
	FxLink string
}

type DailyForecast struct {
	City      string
	FxDate    string
	TempMax   string
	TempMin   string
	TextDay   string
	TextNight string
	Humidity  string
	WindSpeed string
	FxLink    string
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
	Summary string
	FxLink  string
}
