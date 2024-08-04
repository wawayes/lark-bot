package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawayes/lark-bot/initialization"
	"github.com/wawayes/lark-bot/services/weather"
)

func TestConstructTodayWeatherContent(t *testing.T) {
	conf := initialization.GetTestConfig("../config.yaml")
	client := weather.NewWeatherClient(conf.QWeather.Key, conf.QWeather.Location)
	templateJson, err := ConstructTodayWeatherContent(client)
	assert.Nil(t, err)
	assert.NotEmpty(t, templateJson)
}