package config

import (
	"io/ioutil"
	"sync"

	l "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseUrl           string `yaml:"BASE_URL"`
	AppID             string `yaml:"APP_ID"`
	AppSecret         string `yaml:"APP_SECRET"`
	VerificationToken string `yaml:"VERIFICATION_TOKEN"`
	EncryptToken      string `yaml:"ENCRYPT_TOKEN"`
	TestChatID        string `yaml:"TEST_CHAT_ID"`
	AllChatID         string `yaml:"ALL_CHAT_ID"`
}

var (
	conf = &Config{}
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		path := "config.yaml"
		data, err := ioutil.ReadFile(path)
		if err != nil {
			l.Errorf("读取配置文件路径失败: [path]: %s, err: %s", path, err.Error())
		}
		err = yaml.Unmarshal(data, conf)
		if err != nil {
			l.Errorf("yaml 解析配置文件内容失败: [yaml内容]: %s, err: %s", data, err.Error())
		}
	})
	return conf
}

func InitTestConfig() *Config {
	once.Do(func() {
		path := "../config.yaml"
		data, err := ioutil.ReadFile(path)
		if err != nil {
			l.Errorf("读取配置文件路径失败: [path]: %s, err: %s", path, err.Error())
		}
		err = yaml.Unmarshal(data, conf)
		if err != nil {
			l.Errorf("yaml 解析配置文件内容失败: [yaml内容]: %s, err: %s", data, err.Error())
		}
	})
	return conf
}
