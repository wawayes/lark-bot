// /*
//   - @Description:
//   - @Author: wangmingyao@duxiaoman.com
//   - @version:
//   - @Date: 2024-07-28 07:03:02
//   - @LastEditTime: 2024-07-31 00:40:42
//     */
package weather

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawayes/lark-bot/initialization"
)

var (
	key, location string
	latitude = "116.27,40.15"
	client *WeatherClient
)

func TestGetDailyWeather(t *testing.T) {
	resp, err := client.GetDailyWeather(3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(resp.Daily))
}

func TestGetHourlyWeather(t *testing.T) {
	resp, err := client.GetHourlyWeather()
	assert.Nil(t, err)
	assert.Equal(t, 24, len(resp.Hourly))
}

func TestGetNowWeather(t *testing.T) {
	resp, err := client.GetNowWeather()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, resp.Now.Temp)
}

func TestGetDailyIndices(t *testing.T) {
	resp, err := client.GetDailyIndices()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp.Daily))
}

func TestGetGridWeather(t *testing.T) {
	resp, err := client.GetGridWeather(latitude)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, resp.Now.Temp)
}

func TestGetMinutelyWeather(t *testing.T) {
	resp, err := client.GetMinutelyWeather(latitude)
	assert.Nil(t, err)
	assert.NotEqual(t, "", resp.Summary)
}

func TestGetWarningWeather(t *testing.T) {
	resp, err := client.GetWeatherWarning()
	assert.Nil(t, err)
	assert.Equal(t, "200", resp.Code)
}

func TestGetAirQuality(t *testing.T) {
	resp, err := client.GetAirQuality()
	assert.Nil(t, err)
	assert.Equal(t, "200", resp.Code)
}

func TestMain(m *testing.M) {
	path := "../../config.yaml"
	conf := initialization.GetTestConfig(path)
	key = conf.QWeather.Key
	location = conf.QWeather.Location
	client = NewWeatherClient(key, location)

	// This will run the TestGetDailyWeather and any other test functions you have.
    os.Exit(m.Run())
}