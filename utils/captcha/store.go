package captcha

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/mojocn/base64Captcha"
	"rbac_admin_server/global"
)

// CaptchaStore 验证码存储实例
// 使用base64Captcha默认的内存存储
var CaptchaStore = base64Captcha.DefaultMemStore

// EmailCodeStore 邮件验证码存储
// 用于管理邮件验证码的存储、验证和过期清理

type EmailCodeStore struct {
	codes map[string]EmailCodeInfo
	mu    sync.RWMutex
}

// EmailCodeInfo 邮件验证码信息
// 包含验证码内容、关联邮箱和过期时间
type EmailCodeInfo struct {
	Code      string    // 验证码内容
	Email     string    // 关联的邮箱地址
	ExpiredAt time.Time // 过期时间
}

// EmailStore 全局邮件验证码存储实例
var EmailStore = &EmailCodeStore{
	codes: make(map[string]EmailCodeInfo),
}

// SetEmailCode 设置邮件验证码
// email: 接收验证码的邮箱
// code: 验证码内容
// duration: 验证码有效期
func (s *EmailCodeStore) SetEmailCode(email, code string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[email] = EmailCodeInfo{
		Code:      code,
		Email:     email,
		ExpiredAt: time.Now().Add(duration),
	}
}

// GetEmailCode 获取邮件验证码
// email: 邮箱地址
// 返回验证码内容和是否存在且未过期
func (s *EmailCodeStore) GetEmailCode(email string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, exists := s.codes[email]
	if !exists || time.Now().After(info.ExpiredAt) {
		return "", false
	}
	return info.Code, true
}

// VerifyEmailCode 验证邮件验证码
// email: 邮箱地址
// code: 用户输入的验证码
// 返回验证结果
func (s *EmailCodeStore) VerifyEmailCode(email, code string) bool {
	storedCode, exists := s.GetEmailCode(email)
	if !exists {
		return false
	}
	return storedCode == code
}

// CleanExpired 清理过期验证码
// 定期调用以释放内存
func (s *EmailCodeStore) CleanExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	count := 0
	for email, info := range s.codes {
		if now.After(info.ExpiredAt) {
			delete(s.codes, email)
			count++
		}
	}
	if count > 0 {
		global.Logger.Infof("清理了 %d 个过期的邮件验证码", count)
	}
}

// StartCleanupTimer 启动定期清理定时器
// 每5分钟执行一次清理
func (s *EmailCodeStore) StartCleanupTimer() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.CleanExpired()
		}
	}()
}

// GenerateAndStoreEmailCode 生成并存储邮件验证码
// email: 邮箱地址
// duration: 有效期
// 返回生成的验证码
func (s *EmailCodeStore) GenerateAndStoreEmailCode(email string, duration time.Duration) string {
	code := generateEmailCode()
	s.SetEmailCode(email, code, duration)
	return code
}

// generateEmailCode 生成6位数字验证码
func generateEmailCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}