package email

import (
	"crypto/tls"
	"rbac_admin_server/global"

	"gopkg.in/gomail.v2"
)

// SendEmail 发送邮件
// user: 收件人邮箱地址
// subject: 邮件主题
// content: 邮件内容（支持HTML）
// 返回错误信息，成功时为nil
func SendEmail(user string, subject string, content string) error {
	e := global.Config.Email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", e.User)     // 发送者
	mailer.SetHeader("To", user)         // 接收者
	mailer.SetHeader("Subject", subject) // 邮件主题
	mailer.SetBody("text/html", content) // 邮件内容
	// 构建SMTP客户端
	dialer := gomail.NewDialer(e.Host, e.Port, e.User, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 忽略证书校验
	if err := dialer.DialAndSend(mailer); err != nil {
		global.Logger.Errorf("发送邮件失败: %v", err)
		return err
	}
	global.Logger.Infof("邮件发送成功: 收件人=%s, 主题=%s", user, subject)
	return nil
}