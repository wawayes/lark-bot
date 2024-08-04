package global

var weatherErrorMessages = map[string]string{
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

func GetWeatherErrorMessage(code string) string {
    if msg, ok := weatherErrorMessages[code]; ok {
        return msg
    }
    return "未知错误"
}