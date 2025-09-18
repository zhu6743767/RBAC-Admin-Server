package global

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims JWT声明结构
// 包含用户ID和过期时间
// 用于JWT token的payload部分
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
// 使用全局配置中的JWT密钥和过期时间
func GenerateToken(userID uint) (string, error) {
	if Config == nil || Config.JWT.Secret == "" {
		return "", errors.New("JWT配置未初始化")
	}

	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(Config.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    Config.JWT.Issuer,
			Audience:  jwt.ClaimStrings{Config.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Config.JWT.Secret))
}

// ParseToken 解析JWT token
// 验证token有效性并提取用户信息
func ParseToken(tokenString string) (*JWTClaims, error) {
	if Config == nil || Config.JWT.Secret == "" {
		return nil, errors.New("JWT配置未初始化")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}
