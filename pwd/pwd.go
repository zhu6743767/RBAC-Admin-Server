package pwd

// 密码加密

import (
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) string {

	// 使用bcrypt库的GenerateFromPassword函数进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func ComparePassword(hashedPassword, inputPassword string) bool {
	// 使用bcrypt库的CompareHashAndPassword函数进行密码校验
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return false
	}
	return true
}