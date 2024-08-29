package global

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/infrastructure"
)

var (
	Log  *Logger
	once sync.Once
)

type Logger struct {
	*logrus.Logger
}

// InitLogger 初始化全局日志实例
func InitLogger(config infrastructure.Config) {
	once.Do(func() {
		Log = newLogger(config)
	})
}

func newLogger(config infrastructure.Config) *Logger {
	logger := logrus.New()

	// 设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})

	// 设置输出
	if config.Log.Output == "stdout" {
		logger.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(config.Log.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			fmt.Printf("Failed to log to file, using default stderr: %v", err)
		}
	}

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	logger.SetReportCaller(true)

	return &Logger{Logger: logger}
}

// 以下是一些辅助函数，可以直接使用全局 Log 实例

func LogInfo(msg string, fields map[string]interface{}) {
	if Log != nil {
		Log.WithFields(fields).Info(msg)
	}
}

func LogError(msg string, err error, fields map[string]interface{}) {
	if Log != nil {
		Log.WithFields(fields).WithError(err).Error(msg)
	}
}

func LogDebug(msg string, fields map[string]interface{}) {
	if Log != nil {
		Log.WithFields(fields).Debug(msg)
	}
}
