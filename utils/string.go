package utils

import (
	"time"
)

func ParseTime(isoTime string) string {
	// 使用自定义的格式字符串进行解析，确保与输入的时间字符串格式完全匹配
	layout := "2006-01-02T15:04-07:00"
	parsedTime, err := time.Parse(layout, isoTime)
	if err != nil {
		return ""
	}

	// Format the time in "HH:mm" format
	formattedTime := parsedTime.Format("1月2日15:04")
	return formattedTime
}
