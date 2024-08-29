package global

import (
	"fmt"
)

// ErrorCode 定义错误码类型
type ErrorCode string

// 通用错误码
const (
	OK = "0"

	CodeSuccess        = "200"
	CodeNoData         = "204"
	CodeInvalidRequest = "400"
	CodeAuthFailed     = "401"
	CodeExceededLimit  = "402"
	CodeNoPermission   = "403"
	CodeNotFound       = "404"
	CodeExceededQPM    = "429"
	CodeServerError    = "500"
)

// ErrorMessage 定义错误消息映射
type ErrorMessage map[string]string

var WeatherErrorMessages = map[string]string{
	CodeSuccess:        "请求成功",
	CodeNoData:         "请求成功,但你查询的地区暂时没有你需要的数据",
	CodeInvalidRequest: "请求错误,可能包含错误的请求参数或缺少必选的请求参数",
	CodeAuthFailed:     "认证失败,可能使用了错误的KEY、数字签名错误、KEY的类型错误(如使用SDK的KEY去访问Web API)",
	CodeExceededLimit:  "超过访问次数或余额不足以支持继续访问服务,你可以充值、升级访问量或等待访问量重置",
	CodeNoPermission:   "无访问权限,可能是绑定的PackageName、BundleID、域名IP地址不一致,或者是需要额外付费的数据",
	CodeNotFound:       "查询的数据或地区不存在",
	CodeExceededQPM:    "超过限定的QPM(每分钟访问次数),请参考QPM说明",
	CodeServerError:    "无响应或超时,接口服务异常请联系我们",
}

// BasicError 定义基本错误结构
type BasicError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Err        error       `json:"-"`
	StatusCode int         `json:"-"`
}

func (e *BasicError) Error() string {
	if e != nil {
		return fmt.Sprintf("code: %s, msg: %s, data: %+v, err: %s", e.Code, e.Message, e.Data, e.Err)
	}
	return "err is nil"
}

// NewBasicError 创建新的BasicError
func NewBasicError(code string, msg string, data interface{}, err error) *BasicError {
	return &BasicError{
		Code:    code,
		Message: msg,
		Data:    data,
		Err:     err,
	}
}

// 辅助函数
func (e *BasicError) Ok() bool {
	if e.Err != nil {
		return false
	}
	return true
}
