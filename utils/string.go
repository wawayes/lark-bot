package utils

import (
	"time"
)

func ParseTime(input string) string {
	// 定义输入时间的格式
	const inputFormat = "2006-01-02T15:04-07:00"
	// 定义输出时间的格式
	const outputFormat = "2006年1月2日15:04"

	// 解析输入的时间字符串
	parsedTime, err := time.Parse(inputFormat, input)
	if err != nil {
		return ""
	}

	// 格式化解析后的时间
	formattedTime := parsedTime.Format(outputFormat)
	return formattedTime
}
