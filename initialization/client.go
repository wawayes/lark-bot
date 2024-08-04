package initialization

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

var larkClient *lark.Client

func LoadLarkClient(config Config) {
	options := []lark.ClientOptionFunc{
		lark.WithLogLevel(larkcore.LogLevelDebug),
	}
	if config.Lark.BaseUrl != "" {
		options = append(options, lark.WithOpenBaseUrl(config.Lark.BaseUrl))
	}

	larkClient = lark.NewClient(config.Lark.AppID, config.Lark.AppSecret, options...)

}

func GetLarkClient() *lark.Client {
	return larkClient
}
