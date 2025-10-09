package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword 密码加密
// 输入明文密码，返回加密后的密码
func HashedPassword(password string) string {
	// 生成密码哈希值
	// bcrypt.DefaultCost是bcrypt的默认工作因子，通常为10
	// 较高的工作因子会增加计算时间，但也会提高安全性
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// 在实际应用中，应该记录错误并处理
		// 这里为了简化，返回空字符串
		return ""
	}
	// 返回哈希后的密码字符串
	return string(hash)
}

// ComparePassword 密码验证
// 输入存储的哈希密码和待验证的明文密码，返回是否匹配
func ComparePassword(hashedPassword, password string) bool {
	// 使用bcrypt.CompareHashAndPassword验证密码是否匹配
	// 第一个参数是存储的哈希密码，第二个参数是待验证的明文密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// 如果err为nil，表示密码匹配成功
	return err == nil
}

// MakePassword 密码加密（与HashedPassword功能相同，用于兼容API）
// 输入明文密码，返回加密后的密码
func MakePassword(password string) string {
	return HashedPassword(password)
}