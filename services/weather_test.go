package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHourlyWeather(t *testing.T) {
	resp, err := GetWeatherHourly()
	assert.Nil(t, err)
	fmt.Printf("hourly response: %v", resp)
}
