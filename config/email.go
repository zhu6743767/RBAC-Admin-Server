package config

// Email 邮箱配置
// 定义系统中邮件发送相关的配置项
type Email struct {
    User     string `yaml:"user"`     // 发送者邮箱地址
    Password string `yaml:"password"` // 邮箱授权码
    Host     string `yaml:"host"`     // SMTP服务器地址
    Port     int    `yaml:"port"`     // SMTP服务器端口
}

// Verify 验证邮箱配置是否完整
// 返回true表示配置完整，可以用于发送邮件
func (e Email) Verify() bool {
    if e.User == "" || e.Password == "" || e.Host == "" || e.Port == 0 {
        return false
    }
    return true
}