package config

// Captcha 验证码配置
// 定义系统中验证码相关的配置项
type Captcha struct {
    Enable bool `yaml:"enable"` // 是否启用验证码功能
}