package adapters

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wawayes/lark-bot/infrastructure"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(config infrastructure.Config) *Logger {
	logger := logrus.New()

	// 设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
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

	return &Logger{Logger: logger}
}

func (l *Logger) LogInfo(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

func (l *Logger) LogError(msg string, err error, fields map[string]interface{}) {
	l.WithFields(fields).WithError(err).Error(msg)
}

func (l *Logger) LogDebug(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Debug(msg)
}
