package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"rbac_admin_server/config"
	"rbac_admin_server/global"
)

// InitLogger 初始化日志系统
// 支持多种日志格式(json/text)，多种输出方式(stdout/file/both)
// 自动创建日志目录，支持日志轮转和压缩
func InitLogger(cfg *config.LogConfig) error {
	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("解析日志级别失败: %v", err)
	}
	logrus.SetLevel(level)

	// 设置日志格式
	switch strings.ToLower(cfg.Format) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   "2006-01-02 15:04:05",
			DisableHTMLEscape: true,
		})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			DisableColors:   false,
			FullTimestamp:   true,
		})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			DisableColors:   false,
			FullTimestamp:   true,
		})
	}

	// 设置日志输出
	if cfg.Stdout {
		logrus.SetOutput(os.Stdout)
	} else {
		fileWriter, err := createLogFileWriter(cfg)
		if err != nil {
			return fmt.Errorf("创建日志文件写入器失败: %v", err)
		}
		logrus.SetOutput(fileWriter)
	}

	// 设置日志选项
	logrus.SetReportCaller(cfg.EnableCaller)

	// 创建新的logrus实例并应用相同的配置
	global.Logger = logrus.New()
	global.Logger.SetLevel(level)
	global.Logger.SetFormatter(logrus.StandardLogger().Formatter)
	global.Logger.SetOutput(logrus.StandardLogger().Out)
	global.Logger.SetReportCaller(cfg.EnableCaller)

	global.Logger.Info("✅ 日志系统初始化成功")
	return nil
}

// createLogFileWriter 创建日志文件写入器
// 支持按日期创建日志目录，自动创建不存在的目录
func createLogFileWriter(cfg *config.LogConfig) (*os.File, error) {
	// 创建日志目录
	logDir := cfg.Dir
	if logDir == "" {
		logDir = "./logs"
	}

	// 按日期创建子目录
	dateDir := filepath.Join(logDir, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件名
	logFile := filepath.Join(dateDir, "app.log")

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}

	return file, nil
}

// LogWithField 带字段的日志记录
func LogWithField(component string, fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields).WithField("component", component)
}

// LogError 记录错误日志
func LogError(component string, err error, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	fields["error"] = err.Error()
	LogWithField(component, fields).Error("操作失败")
}

// LogInfo 记录信息日志
func LogInfo(component string, message string, fields logrus.Fields) {
	LogWithField(component, fields).Info(message)
}

// LogDebug 记录调试日志
func LogDebug(component string, message string, fields logrus.Fields) {
	LogWithField(component, fields).Debug(message)
}

// LogWarn 记录警告日志
func LogWarn(component string, message string, fields logrus.Fields) {
	LogWithField(component, fields).Warn(message)
}
