package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 日志级别
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	// 日志级别名称
	levelNames = map[int]string{
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
	}

	// 日志级别映射
	levelMap = map[string]int{
		"debug": LevelDebug,
		"info":  LevelInfo,
		"warn":  LevelWarn,
		"error": LevelError,
		"fatal": LevelFatal,
	}

	// 日志输出
	logWriter io.Writer
	// 当前日志级别
	currentLevel int
	// 标准日志
	stdLog *log.Logger
)

// Init 初始化日志系统
func Init(level string, path string) error {
	// 设置日志级别
	lvl, ok := levelMap[strings.ToLower(level)]
	if !ok {
		lvl = LevelInfo
	}
	currentLevel = lvl

	// 创建日志目录
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// 创建日志文件
	logFile := filepath.Join(path, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 同时输出到控制台和文件
	logWriter = io.MultiWriter(os.Stdout, file)
	stdLog = log.New(logWriter, "", log.LstdFlags)

	return nil
}

// Debug 调试日志
func Debug(msg string, keysAndValues ...interface{}) {
	if currentLevel <= LevelDebug {
		logWithLevel(LevelDebug, msg, keysAndValues...)
	}
}

// Info 信息日志
func Info(msg string, keysAndValues ...interface{}) {
	if currentLevel <= LevelInfo {
		logWithLevel(LevelInfo, msg, keysAndValues...)
	}
}

// Warn 警告日志
func Warn(msg string, keysAndValues ...interface{}) {
	if currentLevel <= LevelWarn {
		logWithLevel(LevelWarn, msg, keysAndValues...)
	}
}

// Error 错误日志
func Error(msg string, keysAndValues ...interface{}) {
	if currentLevel <= LevelError {
		logWithLevel(LevelError, msg, keysAndValues...)
	}
}

// Fatal 致命错误日志
func Fatal(msg string, err error) {
	if currentLevel <= LevelFatal {
		if err != nil {
			logWithLevel(LevelFatal, msg, "error", err)
		} else {
			logWithLevel(LevelFatal, msg)
		}
		os.Exit(1)
	}
}

// 带级别的日志输出
func logWithLevel(level int, msg string, keysAndValues ...interface{}) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	// 提取文件名
	file = filepath.Base(file)

	// 格式化键值对
	var kvStr string
	if len(keysAndValues) > 0 {
		kvPairs := make([]string, 0, len(keysAndValues)/2)
		for i := 0; i < len(keysAndValues); i += 2 {
			key := fmt.Sprintf("%v", keysAndValues[i])
			var value string
			if i+1 < len(keysAndValues) {
				value = fmt.Sprintf("%v", keysAndValues[i+1])
			} else {
				value = "<no value>"
			}
			kvPairs = append(kvPairs, fmt.Sprintf("%s=%s", key, value))
		}
		kvStr = " " + strings.Join(kvPairs, " ")
	}

	// 输出日志
	stdLog.Printf("[%s] %s:%d %s%s", levelNames[level], file, line, msg, kvStr)
}

// LogEntry 日志条目
type LogEntry struct {
	Time    time.Time
	Level   string
	Message string
	File    string
	Line    int
	Data    map[string]string
}

// GetLogs 获取日志列表
func GetLogs(level string, page, pageSize int) ([]*LogEntry, int, error) {
	// 简化实现，实际应该从日志文件中读取
	logs := make([]*LogEntry, 0)

	// 生成一些示例日志
	for i := 0; i < pageSize; i++ {
		logs = append(logs, &LogEntry{
			Time:    time.Now().Add(-time.Duration(i) * time.Hour),
			Level:   "INFO",
			Message: fmt.Sprintf("示例日志 %d", i+1),
			File:    "main.go",
			Line:    100 + i,
			Data:    map[string]string{"key": "value"},
		})
	}

	return logs, 100, nil // 假设总共有100条日志
}

// ClearLogs 清除日志
func ClearLogs(level string) error {
	// 简化实现，实际应该删除或清空日志文件
	Info("清除日志", "level", level)
	return nil
}
