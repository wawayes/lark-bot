/*
 * @Description:
 * @Author: wangmingyao@duxiaoman.com
 * @version:
 * @Date: 2024-07-28 07:03:02
 * @LastEditTime: 2024-07-31 00:40:42
 */
package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTodayWeather(t *testing.T) {
	resp, err := GetWeatherDay(3)
	assert.Nil(t, err)
	fmt.Printf("today weather: %v", resp.Daily[0].TextDay)
}

func TestHourlyWeather(t *testing.T) {
	resp, err := GetWeatherHourly()
	assert.Nil(t, err)
	fmt.Printf("hourly response: %v", resp)
}

func TestGetWarningWeather(t *testing.T) {
	resp, err := GetWeatherWarning()
	assert.Nil(t, err)
	fmt.Printf("warning weather: %v", resp)
}

func TestGetGridWeather(t *testing.T) {
	resp, err := GetGridWeather(SUZHOUJIE_LOCATION)
	assert.Nil(t, err)
	fmt.Printf("grid weather: %s", resp.Now.Text)
}

func TestGetDaily(t *testing.T) {
	resp, err := GetDaily()
	assert.Nil(t, err)
	fmt.Printf("daily: %v", resp)
}
