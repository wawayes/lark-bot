package handlers

import (
	"os"
	"testing"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/wawayes/lark-bot/initialization"
)

var (
	client *lark.Client
)

func TestGetBotInfo(t *testing.T) {
	botInfo, err := GetBotInfo(client)
	if err != nil {
		t.Errorf("TestGetBotInfo failed: %v", err)
	}
	t.Logf("TestGetBotInfo success: %s", botInfo)
}

func TestMain(m *testing.M) {
	path := "../config.yaml"
	conf := initialization.GetTestConfig(path)
	initialization.LoadLarkClient(*conf)
	client = initialization.GetLarkClient()
	// This will run the TestGetDailyWeather and any other test functions you have.
	os.Exit(m.Run())
}
