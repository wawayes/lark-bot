package services

import (
	"encoding/json"
	"errors"
	"fmt"

	l "github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/utils"
)

// 五福家园3号楼 116.271227,40.151309

// 苏州街地铁站 116.312784,39.981726

// 西二旗地铁站 116.312564,40.059029

const (
	WEATHER_BASE_URL   = "https://devapi.qweather.com/v7"                          // 天气地址前缀
	WEATHER_URL        = "https://www.qweather.com/weather/beijing-101010100.html" // 天气地址
	KEY                = "27df0ab1a3014458b59906f1c8bfa6f7"                        // api key
	LOCATION           = "101010100"                                               // 北京的LocationID
	HOME_LOCATION      = "116.27,40.15"                                            // 家的坐标
	SUZHOUJIE_LOCATION = "116.31,39.98"                                            // 苏州街坐标
	XIERQI_LOCATION    = "116.31,40.05"                                            // 西二旗坐标
)

// CommonWeatherDetails - 存储多种天气数据中的共通字段
type CommonWeatherDetails struct {
	Temp      string `json:"temp"`      // 温度
	Icon      string `json:"icon"`      // 天气状况图标代码
	Text      string `json:"text"`      // 天气状况文字描述
	Wind360   string `json:"wind360"`   // 风向360角度
	WindDir   string `json:"windDir"`   // 风向
	WindScale string `json:"windScale"` // 风力等级
	WindSpeed string `json:"windSpeed"` // 风速
	Humidity  string `json:"humidity"`  // 相对湿度
	Precip    string `json:"precip"`    // 降水量
	Pressure  string `json:"pressure"`  // 大气压强
	Cloud     string `json:"cloud"`     // 云量 百分比数值
}

// 每日天气
type DailyWeather struct {
	FxDate        string `json:"fxDate"`        // 预报日期
	Sunrise       string `json:"sunrise"`       // 日出时间
	Sunset        string `json:"sunset"`        // 日落时间
	Moonrise      string `json:"moonrise"`      // 月升时间
	Moonset       string `json:"moonset"`       // 月落时间
	MoonPhase     string `json:"moonPhase"`     // 月相名称
	MoonPhaseIcon string `json:"moonPhaseIcon"` // 月相图标代码
	TempMax       string `json:"tempMax"`       // 当天最高温度
	TempMin       string `json:"tempMin"`       // 当天最低温度
	TextDay       string `json:"textDay"`       // 白天天气描述
	UvIndex       string `json:"uvIndex"`       // 紫外线强度指数
	Vis           string `json:"vis"`           // 能见度
	CommonWeatherDetails
}

// 实时天气
type NowWeather struct {
	ObsTime   string `json:"obsTime"`   // 数据观测时间
	FeelsLike string `json:"feelsLike"` // 体感温度
	Dew       string `json:"dew"`       // 露点温度
	CommonWeatherDetails
}

// 逐小时天气
type HourlyWeather struct {
	FxTime string `json:"fxTime"` // 预报时间
	Pop    string `json:"pop"`    // 逐小时预报降水概率 百分比数值，可能为空
	Dew    string `json:"dew"`    // 露点温度
	CommonWeatherDetails
}

// WeatherResponse - 通用天气响应体
type WeatherResponse struct {
	Code       string          `json:"code"`             // 状态码
	UpdateTime string          `json:"updateTime"`       // 更新时间
	FxLink     string          `json:"fxLink"`           // 数据链接
	Daily      []DailyWeather  `json:"daily,omitempty"`  // 每日天气数据
	Now        NowWeather      `json:"now,omitempty"`    // 实时天气数据
	Hourly     []HourlyWeather `json:"hourly,omitempty"` // 每小时天气数据
	Refer      struct {
		Sources []string `json:"sources"` // 数据来源
		License []string `json:"license"` // 许可信息
	} `json:"refer"`
}

// 天气预警
type WarningWeather struct {
	ID            string `json:"id"`            // 本条预警的唯一标识
	Sender        string `json:"sender"`        // 预警发布单位
	PubTime       string `json:"pubTime"`       // 预警发布时间
	Title         string `json:"title"`         // 预警信息标题
	StartTime     string `json:"startTime"`     // 预警开始时间
	EndTime       string `json:"endTime"`       // 预警结束时间
	Status        string `json:"status"`        // 预警信息的发布状态
	Severity      string `json:"severity"`      // 预警严重等级
	SeverityColor string `json:"severityColor"` // 预警严重等级颜色
	Type          string `json:"type"`          // 预警类型ID
	TypeName      string `json:"typeName"`      // 预警类型名称
	Urgency       string `json:"urgency"`       // 预警信息的紧迫程度
	Certainty     string `json:"certainty"`     // 预警信息的确定性
	Text          string `json:"text"`          // 预警详细文字描述
	Related       string `json:"related"`       // 相关联的预警ID
}

// 天气预警结构体
type WarningResponse struct {
	Code       string           `json:"code"`       // 状态码
	UpdateTime string           `json:"updateTime"` // 当前API的最近更新时间
	FxLink     string           `json:"fxLink"`     // 当前数据的响应式页面
	Warning    []WarningWeather `json:"warning"`    // 预警列表
	Refer      struct {
		Sources []string `json:"sources"` // 原始数据来源
		License []string `json:"license"` // 数据许可或版权声明
	} `json:"refer"`
}

// 格点天气响应体
type GridResponse struct {
	Code       string               `json:"code"`       // 状态码
	UpdateTime string               `json:"updateTime"` // 当前API的最近更新时间
	FxLink     string               `json:"fxLink"`     // 当前数据的响应式页面
	Now        CommonWeatherDetails `json:"now"`        // 格点实时天气
	Refer      struct {
		Sources []string `json:"sources"` // 原始数据来源
		License []string `json:"license"` // 数据许可或版权声明
	} `json:"refer"`
}

// 生活指数
type DailyLevel struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Category string `json:"category"`
	Text     string `json:"text"`
}

// 生活指数响应体
type DailyResponse struct {
	Code       string       `json:"code"`
	UpdateTime string       `json:"updateTime"`
	FxLink     string       `json:"fxLink"`
	Daily      []DailyLevel `json:"daily"`
	Refer      struct {
		Sources []string `json:"sources"` // 原始数据来源
		License []string `json:"license"` // 数据许可或版权声明
	} `json:"refer"`
}

// 实时天气接口
func GetWeatherNow() (*WeatherResponse, error) {
	url := fmt.Sprintf("%s/weather/now?key=%s&location=%s", WEATHER_BASE_URL, KEY, LOCATION)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get weather err: %s", err.Error())
		return nil, err
	}
	var resp WeatherResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("weather now json Unmarshal err: %s", err.Error())
		return nil, err
	}
	if resp.Code != "200" {
		l.Errorf("get weather warning code is not 200, resp: %v", resp)
		return nil, errors.New("code is not 200")
	}
	return &resp, nil
}

// 每日天气接口
// param: num 查询接下来num天的天气
func GetWeatherDay(num int) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s/weather/%dd?key=%s&location=%s", WEATHER_BASE_URL, num, KEY, LOCATION)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get daily weather err: %s", err.Error())
		return nil, err
	}
	var resp WeatherResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("daily weather now json Unmarshal err: %s", err.Error())
		return nil, err
	}
	if resp.Code != "200" {
		l.Errorf("get weather warning code is not 200, resp: %v", resp)
		return nil, errors.New("code is not 200")
	}
	return &resp, nil
}

// 逐小时天气接口
func GetWeatherHourly() (*WeatherResponse, error) {
	url := fmt.Sprintf("%s/weather/24h?key=%s&location=%s", WEATHER_BASE_URL, KEY, LOCATION)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get hourly daily weather err: %s", err.Error())
		return nil, err
	}
	var resp WeatherResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("hourly weather now json Unmarshal err: %s", err.Error())
		return nil, err
	}
	if resp.Code != "200" {
		l.Errorf("get weather warning code is not 200, resp: %v", resp)
		return nil, errors.New("code is not 200")
	}
	return &resp, nil
}

// 获取天气预警
func GetWeatherWarning() (*WarningResponse, error) {
	url := fmt.Sprintf("%s/warning/now?key=%s&location=%s", WEATHER_BASE_URL, KEY, LOCATION)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get weather warning err: %s", err.Error())
		return nil, err
	}
	var resp WarningResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("get weather warning json unmarshal err: %s", err.Error())
		return nil, err
	}
	if resp.Code != "200" {
		l.Errorf("get weather warning code is not 200, resp: %v", resp)
		return nil, errors.New("code is not 200")
	}
	return &resp, nil
}

// 获取格点实时天气
func GetGridWeather(location string) (*GridResponse, error) {
	url := fmt.Sprintf("%s/grid-weather/now?key=%s&location=%s", WEATHER_BASE_URL, KEY, location)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get grid weather err: %s", err.Error())
		return nil, err
	}
	var resp GridResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("json unmarshal err: %s", err.Error())
		return nil, err
	}
	if resp.Code != "200" {
		l.Errorf("get grid weather code is not 200")
		return nil, errors.New("code is not 200")
	}
	return &resp, nil
}

// 获取生活指数
func GetDaily() (*DailyResponse, error) {
	url := fmt.Sprintf("%s/indices/1d?type=3,8&key=%s&location=%s", WEATHER_BASE_URL, KEY, LOCATION)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get daily err: %s", err.Error())
		return nil, err
	}
	var resp DailyResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("json unmarshal err: %s", err.Error())
		return nil, err
	}
	if resp.Code != "200" {
		l.Errorf("get daily code is not 200")
		return nil, errors.New("code is not 200")
	}
	return &resp, nil
}

// 获取空气质量
func GetAirCondition() (string, error) {
	type airResponse struct {
		Code string `json:"code"`
		Now  struct {
			Category string `json:"category"` // 空气质量描述
		} `json:"now"`
	}
	url := fmt.Sprintf("%s/air/now?location=%s&key=%s", WEATHER_BASE_URL, LOCATION, KEY)
	data, err := utils.HttpGet(url)
	if err != nil {
		l.Errorf("get daily err: %s", err.Error())
		return "", err
	}
	var resp airResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		l.Errorf("json unmarshal err: %s", err.Error())
		return "", err
	}
	if resp.Code != "200" {
		l.Errorf("get daily code is not 200")
		return "", errors.New("code is not 200")
	}
	return resp.Now.Category, nil
}
