package infrastructure

import (
	"io/ioutil"
	"sync"

	l "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config 表示顶层配置结构
type Config struct {
	Env      EnvConfig      `yaml:"env"`
	Redis    RedisConfig    `yaml:"redis"`
	Lark     LarkConfig     `yaml:"lark"`
	QWeather QWeatherConfig `yaml:"qweather"`
	Tasks    Tasks          `yaml:"tasks"`
	Log      Log            `yaml:"log"`
}

// EnvConfig 包含环境相关配置
type EnvConfig struct {
	HttpPort string `yaml:"http_port"`
}

// RedisConfig 包含Redis的配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LarkConfig 包含Lark应用的配置
type LarkConfig struct {
	BaseUrl           string `yaml:"base_url"`
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

// Tasks 包含定时任务的配置
type Tasks struct {
	DailyWeather   string `yaml:"daily_weather"`
	WeatherWarning string `yaml:"weather_warning"`
}

type Log struct {
	Level  string
	Output string
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
