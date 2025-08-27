// utils/logger.go
package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"task4/config"
	"time"
)

// 日志级别
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

// Logger 日志记录器
type Logger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	logLevel string
	logFile  *os.File
}

var logger *Logger

// InitLogger 初始化日志记录器
func InitLogger(cfg *config.Config) error {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 创建多输出写入器（同时输出到文件和控制台）
	multiWriter := io.MultiWriter(os.Stdout, file)

	logger = &Logger{
		debugLog: log.New(multiWriter, "[DEBUG] ", log.Ldate|log.Ltime),
		infoLog:  log.New(multiWriter, "[INFO]  ", log.Ldate|log.Ltime),
		warnLog:  log.New(multiWriter, "[WARN]  ", log.Ldate|log.Ltime),
		errorLog: log.New(multiWriter, "[ERROR] ", log.Ldate|log.Ltime),
		logLevel: strings.ToUpper(cfg.LogLevel),
		logFile:  file,
	}

	// 记录日志系统启动信息
	logger.Info("日志系统初始化完成", map[string]interface{}{
		"log_level": logger.logLevel,
		"log_file":  cfg.LogFile,
	})

	return nil
}

// GetLogger 获取日志记录器实例
func GetLogger() *Logger {
	return logger
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// 获取调用者信息
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Debug 记录调试信息
func (l *Logger) Debug(message string, data map[string]interface{}) {
	if l.logLevel == DEBUG {
		caller := getCallerInfo()
		l.debugLog.Printf("%s %s %s", caller, message, formatData(data))
	}
}

// Info 记录普通信息
func (l *Logger) Info(message string, data map[string]interface{}) {
	if l.logLevel == DEBUG || l.logLevel == INFO || l.logLevel == WARN || l.logLevel == ERROR {
		caller := getCallerInfo()
		l.infoLog.Printf("%s %s %s", caller, message, formatData(data))
	}
}

// Warn 记录警告信息
func (l *Logger) Warn(message string, data map[string]interface{}) {
	if l.logLevel == DEBUG || l.logLevel == INFO || l.logLevel == WARN || l.logLevel == ERROR {
		caller := getCallerInfo()
		l.warnLog.Printf("%s %s %s", caller, message, formatData(data))
	}
}

// Error 记录错误信息
func (l *Logger) Error(message string, data map[string]interface{}) {
	if l.logLevel == DEBUG || l.logLevel == INFO || l.logLevel == WARN || l.logLevel == ERROR {
		caller := getCallerInfo()
		l.errorLog.Printf("%s %s %s", caller, message, formatData(data))
	}
}

// formatData 格式化日志数据
func formatData(data map[string]interface{}) string {
	if data == nil || len(data) == 0 {
		return ""
	}

	var pairs []string
	for k, v := range data {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(pairs, " ")
}

// LogRequest 记录HTTP请求信息
func (l *Logger) LogRequest(method, path, clientIP string, statusCode int, latency time.Duration) {
	l.Info("HTTP请求", map[string]interface{}{
		"method":    method,
		"path":      path,
		"client_ip": clientIP,
		"status":    statusCode,
		"latency":   latency.String(),
	})
}

// LogAuth 记录认证相关日志
func (l *Logger) LogAuth(action, username string, success bool, details map[string]interface{}) {
	status := "失败"
	if success {
		status = "成功"
	}

	if details == nil {
		details = make(map[string]interface{})
	}

	details["action"] = action
	details["username"] = username
	details["status"] = status

	l.Info("认证操作", details)
}
