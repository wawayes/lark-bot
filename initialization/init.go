package initialization

import (
	"io/ioutil"
	"sync"

	l "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config 表示顶层配置结构
type Config struct {
    Env      EnvConfig      `yaml:"env"`
    Lark     LarkConfig     `yaml:"lark"`
    QWeather QWeatherConfig `yaml:"qweather"`
}

// EnvConfig 包含环境相关配置
type EnvConfig struct {
    HttpPort string `yaml:"http_port"`
}

// LarkConfig 包含Lark应用的配置
type LarkConfig struct {
    BaseUrl string `yaml:"base_url"`
    AppID             string `yaml:"app_id"`
    AppSecret         string `yaml:"app_secret"`
    EncryptToken      string `yaml:"encrypt_token"`
    VerificationToken string `yaml:"verification_token"`
    TestChatID        string `yaml:"test_chat_id"`
    AllChatID         string `yaml:"all_chat_id"`
}

// QWeatherConfig 包含天气应用的配置
type QWeatherConfig struct {
    Key      string `yaml:"key"`
    Location string `yaml:"location"`
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
            l.Errorf("Failed to read config file: [path]: %s, err: %s", path, err.Error())
            return
        }
        err = yaml.Unmarshal(data, conf)
        if err != nil {
            l.Errorf("Failed to parse YAML config: [yaml content]: %s, err: %s", data, err.Error())
        }
    })
    return conf
}

func GetTestConfig(path string) *Config {
    once.Do(func() {
        data, err := ioutil.ReadFile(path)
        if err != nil {
            l.Errorf("Failed to read config file: [path]: %s, err: %s", path, err.Error())
            return
        }
        err = yaml.Unmarshal(data, conf)
        if err != nil {
            l.Errorf("Failed to parse YAML config: [yaml content]: %s, err: %s", data, err.Error())
        }
    })
    return conf
}