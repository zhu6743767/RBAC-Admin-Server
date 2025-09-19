package global

import (
	"errors"
	"rbac.admin/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ClaimsUserInfo 自定义JWT声明结构，包含用户基本信息
// 用于在JWT token中存储用户相关数据
type ClaimsUserInfo struct {
	UserID   uint   `json:"userID"`
	Username string `json:"username"`
	RoleList []uint `json:"roleList"`
}

// JWTClaims JWT声明结构
// 包含用户信息和标准声明
// 用于JWT token的payload部分
type JWTClaims struct {
	ClaimsUserInfo
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
// 使用全局配置中的JWT密钥和过期时间
// info: 用户信息，包含用户ID、用户名和角色列表
func GenerateToken(info ClaimsUserInfo) (string, error) {
	if Config == nil || Config.JWT.Secret == "" {
		return "", errors.New("JWT配置未初始化")
	}

	claims := JWTClaims{
		ClaimsUserInfo: info,
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
// tokenString: 要解析的JWT token字符串
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
		// 验证签发者
		if claims.Issuer != Config.JWT.Issuer {
			return nil, errors.New("token签发者无效")
		}
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// RefreshToken 刷新JWT token
// tokenString: 过期但未失效的JWT token
func RefreshToken(tokenString string) (string, error) {
	// 解析token但不验证过期时间
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, _, err := parser.ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", errors.New("无效的token格式")
	}

	// 检查token是否接近过期（5分钟内）或已过期
	now := time.Now()
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(now.Add(-5*time.Minute)) {
		return "", errors.New("token已过期太久，无法刷新")
	}

	// 生成新的token，保持原有的用户信息
	return GenerateToken(claims.ClaimsUserInfo)
}

// GetUserRoles 获取用户角色列表
// userID: 用户ID
func GetUserRoles(userID uint) ([]uint, error) {
	var userRoles []models.UserRole
	result := DB.Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, result.Error
	}

	roleList := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleList[i] = ur.RoleID
	}

	return roleList, nil
}
