package core

import (
	"os"
	"regexp"
	"rbac_admin_server/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// 注意：此文件中的ReadConfig函数已被config.Load函数替代
// 为保持兼容性暂时保留，但推荐使用config.Load

// ReadConfig 读取配置文件（已被config.Load替代，建议使用新函数）
func ReadConfig(filePath string) *config.Config {
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Fatalf("❌ 配置文件读取失败: %v", err.Error())
		return nil
	}

	// 替换配置文件中的环境变量
	content := replaceEnvVars(string(byteData))

	var c *config.Config
	err = yaml.Unmarshal([]byte(content), &c)
	if err != nil {
		logrus.Fatalf("❌ 配置文件格式解析失败: %v", err.Error())
		return nil
	}

	// 如果配置为空，使用默认配置
	if c == nil {
		c = config.DefaultConfig()
		logrus.Infof("✅ 使用默认配置")
	} else {
		logrus.Infof("✅ 配置文件加载成功: %s", filePath)
	}

	return c
}

// replaceEnvVars 替换字符串中的环境变量，格式为 ${ENV_VAR:default_value}
func replaceEnvVars(content string) string {
	// 定义环境变量正则表达式
	envRegex := regexp.MustCompile(`\$\{([^:}]+)(:([^}]*))?\}`)

	// 替换所有匹配的环境变量
	result := envRegex.ReplaceAllStringFunc(content, func(match string) string {
		// 提取环境变量名和默认值
		matches := envRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		envName := matches[1]
		defaultValue := ""
		if len(matches) > 3 {
			defaultValue = matches[3]
		}

		// 获取环境变量值，如果不存在则使用默认值
		value := os.Getenv(envName)
		if value == "" {
			value = defaultValue
		}

		return value
	})

	return result
}