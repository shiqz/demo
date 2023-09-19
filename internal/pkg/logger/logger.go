// Package logger 记录系统日志
package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var (
	// 日志格式
	formatter = &prefixed.TextFormatter{
		ForceColors:      true,
		ForceFormatting:  true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		DisableSorting:   false,
		TimestampFormat:  "2006/01/02 15:04:05",
	}
	// 日志输出
	writer = io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "runtime/logs/app.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
)

type Logger struct {
	*logrus.Logger
}

// init 默认日志
func init() {
	logrus.SetFormatter(formatter)
	logrus.SetOutput(writer)
}

// New 创建日志
func New(lv string) *Logger {
	lg := logrus.New()
	lg.SetOutput(writer)
	lg.SetFormatter(formatter)
	lg.SetLevel(GetLogLevel(lv))
	return &Logger{lg}
}

func GetLogLevel(lv string) logrus.Level {
	level, err := logrus.ParseLevel(lv)
	if err != nil {
		level = logrus.TraceLevel
	}
	return level
}
