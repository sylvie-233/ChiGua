package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Logger 全局日志实例
var Logger *logrus.Logger

// logFileHandle 全局日志文件句柄
var logFileHandle *os.File

// InitLogger 初始化日志
func InitLogger(logLevel, logFile string) {
	// 创建日志实例
	Logger = logrus.New()

	// 设置日志级别
	switch logLevel {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 输出到控制台
	Logger.SetOutput(os.Stdout)

	// 如果指定了日志文件，添加文件输出
	if logFile != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			return
		}

		// 创建日志文件
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("打开日志文件失败: %v\n", err)
			return
		}

		// 保存日志文件句柄
		logFileHandle = file

		// 同时输出到控制台和文件
		Logger.SetOutput(io.MultiWriter(os.Stdout, file))
	}
}

// CloseLogger 关闭日志文件
func CloseLogger() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}

// Debug 输出调试日志
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf 输出格式化调试日志
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info 输出信息日志
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof 输出格式化信息日志
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn 输出警告日志
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf 输出格式化警告日志
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error 输出错误日志
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf 输出格式化错误日志
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal 输出致命错误日志并退出
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf 输出格式化致命错误日志并退出
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}
