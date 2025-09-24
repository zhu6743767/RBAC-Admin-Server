package core

import (
	"os"
	"rbac_admin_server/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ReadConfig 读取配置文件
func ReadConfig(filePath string) *config.Config {
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Fatalf("❌ 配置文件读取失败: %v", err.Error())
		return nil
	}
	var c *config.Config
	err = yaml.Unmarshal(byteData, &c)
	if err != nil {
		logrus.Fatalf("❌ 配置文件格式解析失败: %v", err.Error())
		return nil
	}
	logrus.Infof("✅ 配置文件加载成功: %s", filePath)
	return c
}

// InitLogger 初始化日志系统
func InitLogger(logDir string) {
	// 设置日志输出目录
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			logrus.Fatalf("创建日志目录失败: %v", err)
		}
	}

	// 设置日志格式为JSON
	global.Logger = logrus.New()
	global.Logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	// 设置日志级别
	level, err := logrus.ParseLevel(global.Config.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	global.Logger.SetLevel(level)

	// 同时输出到控制台和文件
	// global.Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	logrus.Info("✅ 日志系统初始化成功")
}