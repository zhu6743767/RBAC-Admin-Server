package global

import (
	"errors"
	"rbac_admin_server/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// RefreshClaims 刷新令牌声明结构
type RefreshClaims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
// 使用全局配置中的JWT密钥和过期时间
// info: 用户信息，包含用户ID、用户名和角色列表
func GenerateToken(info ClaimsUserInfo) (string, error) {
	if Config == nil || Config.JWT.Secret == "" {
		return "", errors.New("JWT配置未初始化")
	}

	// 选择签名方法
	var signingMethod jwt.SigningMethod
	switch strings.ToUpper(Config.JWT.SigningMethod) {
	case "HS256":
		signingMethod = jwt.SigningMethodHS256
	case "HS384":
		signingMethod = jwt.SigningMethodHS384
	case "HS512":
		signingMethod = jwt.SigningMethodHS512
	default:
		signingMethod = jwt.SigningMethodHS256 // 默认使用HS256
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

	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString([]byte(Config.JWT.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint) (string, error) {
	if Config == nil || Config.JWT.Secret == "" {
		return "", errors.New("JWT配置未初始化")
	}

	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(Config.JWT.RefreshExpireHours) * time.Hour)),
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
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
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
		// 验证受众
		validAudience := false
		for _, aud := range claims.Audience {
			if aud == Config.JWT.Audience {
				validAudience = true
				break
			}
		}
		if !validAudience {
			return nil, errors.New("token受众无效")
		}
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	if Config == nil || Config.JWT.Secret == "" {
		return nil, errors.New("JWT配置未初始化")
	}

	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(Config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		if claims.Issuer != Config.JWT.Issuer {
			return nil, errors.New("token签发者无效")
		}
		// 验证受众
		validAudience := false
		for _, aud := range claims.Audience {
			if aud == Config.JWT.Audience {
				validAudience = true
				break
			}
		}
		if !validAudience {
			return nil, errors.New("token受众无效")
		}
		return claims, nil
	}

	return nil, errors.New("无效的刷新令牌")
}

// RefreshToken 刷新JWT token
// refreshToken: 有效的刷新令牌
func RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新令牌
	refreshClaims, err := ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 获取用户角色信息
	roleList, err := GetUserRoles(refreshClaims.UserID)
	if err != nil {
		return "", "", err
	}

	// 创建新的用户信息
	userInfo := ClaimsUserInfo{
		UserID:   refreshClaims.UserID,
		Username: "", // 在实际应用中应该从数据库获取用户名
		RoleList: roleList,
	}

	// 生成新的访问令牌和刷新令牌
	newAccessToken, err := GenerateToken(userInfo)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(refreshClaims.UserID)
	if err != nil {
		return newAccessToken, "", err
	}

	return newAccessToken, newRefreshToken, nil
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
