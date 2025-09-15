package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger 接口定义
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Close() error
}

// Config 日志配置
type Config struct {
	Type         string // 日志类型
	Level        string // 日志级别
	Format       string // 输出格式
	Output       string // 输出位置
	LogDir       string // 日志根目录
	MaxSize      int    // 日志文件最大大小(MB)
	MaxAge       int    // 日志文件最大保存天数
	MaxBackups   int    // 日志文件最大备份数量
	Compress     bool   // 是否压缩日志文件
	LocalTime    bool   // 是否使用本地时间
	EnableCaller bool   // 记录调用者信息
	EnableTrace  bool   // 记录堆栈跟踪
}

// logrusLogger 基于logrus的日志器实现
type logrusLogger struct {
	logger *logrus.Logger
	config *Config
}

// Factory 日志工厂接口
type Factory interface {
	Create(config *Config) (Logger, error)
}

// logrusFactory logrus日志工厂
type logrusFactory struct{}

// NewFactory 创建日志工厂
func NewFactory(factoryType string) (Factory, error) {
	switch factoryType {
	case "logrus", "":
		return &logrusFactory{}, nil
	default:
		return nil, fmt.Errorf("不支持的日志工厂类型: %s", factoryType)
	}
}

// GetFactory 获取日志工厂（兼容旧接口）
func GetFactory(factoryType string) (Factory, error) {
	return NewFactory(factoryType)
}

// Create 创建logrus日志器
func (f *logrusFactory) Create(config *Config) (Logger, error) {
	return New(config)
}

// Close 关闭日志器（添加Close方法以满足接口要求）
func (l *logrusLogger) Close() error {
	// logrus没有Close方法，这里返回nil
	return nil
}

// New 创建新的日志器
func New(config *Config) (Logger, error) {
	if config == nil {
		config = &Config{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			LogDir:     "./logs",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 3,
			Compress:   true,
			LocalTime:  true,
		}
	}

	l := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	l.SetLevel(level)

	// 设置格式
	if config.Format == "json" {
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		l.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// 设置输出
	switch config.Output {
	case "file":
		dateDir := time.Now().Format("2006-01-02")
		logDir := filepath.Join(config.LogDir, dateDir)
		
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}
		
		filename := filepath.Join(logDir, "app.log")
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("创建日志文件失败: %w", err)
		}
		l.SetOutput(file)
	case "both":
		// 同时输出到文件和标准输出
		dateDir := time.Now().Format("2006-01-02")
		logDir := filepath.Join(config.LogDir, dateDir)
		
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}
		
		filename := filepath.Join(logDir, "app.log")
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("创建日志文件失败: %w", err)
		}
		l.SetOutput(file)
	default:
		l.SetOutput(os.Stdout)
	}

	return &logrusLogger{
		logger: l,
		config: config,
	}, nil
}

// 实现Logger接口的方法
func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
