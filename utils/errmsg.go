package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 错误码常量定义
const (
	SUCCESS        = 200
	ERROR          = 500
	// 用户模块错误
	ERROR_USERNAME_USED   = 1001
	ERROR_PASSWORD_WRONG  = 1002
	ERROR_USER_NOT_EXIST  = 1003
	ERROR_TOKEN_EXIST     = 1004
	ERROR_TOKEN_RUNTIME   = 1005
	ERROR_TOKEN_WRONG     = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT   = 1008
	ERROR_TOKEN_INVALID   = 1009
	ERROR_UNAUTHORIZED    = 401
	ERROR_GET_USER        = 1010
	ERROR_INVALID_PARAM   = 400
	ERROR_UPDATE_USER     = 1011
	ERROR_ENCRYPT_PASSWORD = 1012
	// 文章模块错误
	ERROR_ART_NOT_EXIST   = 2001
	// 分类模块错误
	ERROR_CATENAME_USED   = 3001
	ERROR_CATE_NOT_EXIST  = 3002
	// 权限错误
	ERROR_PERMISSION_DENIED = 4001
	// 验证码和邮件错误
	ERROR_CAPTCHA_WRONG   = 5001
	ERROR_CAPTCHA_EXPIRE  = 5002
	ERROR_EMAIL_SEND      = 5003
	ERROR_EMAIL_CODE_WRONG = 5004
	ERROR_EMAIL_CODE_EXPIRE = 5005
	ERROR_EMAIL_CONFIG    = 5006
)

// GetValidationError 将validator错误转换为字符串
func GetValidationError(errs validator.ValidationErrors) string {
	for _, err := range errs {
		// 简单实现，根据需要可以扩展为更详细的错误信息
		return err.Field() + "字段验证失败"
	}
	return "参数验证失败"
}

// 错误信息映射表
var codeMsg = map[int]string{
	SUCCESS:              "OK",
	ERROR:                "FAIL",
	ERROR_USERNAME_USED:  "用户名已存在",
	ERROR_PASSWORD_WRONG: "密码错误",
	ERROR_USER_NOT_EXIST: "用户不存在",
	ERROR_TOKEN_EXIST:    "TOKEN不存在",
	ERROR_TOKEN_RUNTIME:  "TOKEN已过期",
	ERROR_TOKEN_WRONG:    "TOKEN不正确",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误",
	ERROR_USER_NO_RIGHT:  "用户无权限",
	ERROR_TOKEN_INVALID:  "无效的TOKEN",
	ERROR_ART_NOT_EXIST:  "文章不存在",
	ERROR_CATENAME_USED:  "该分类已存在",
	ERROR_CATE_NOT_EXIST: "分类不存在",
	ERROR_PERMISSION_DENIED: "权限不足",
	ERROR_UNAUTHORIZED:    "未授权",
	ERROR_GET_USER:        "获取用户信息失败",
	ERROR_INVALID_PARAM:   "参数无效",
	ERROR_UPDATE_USER:     "更新用户信息失败",
	ERROR_ENCRYPT_PASSWORD: "密码加密失败",
	ERROR_CAPTCHA_WRONG:   "验证码错误",
	ERROR_CAPTCHA_EXPIRE:  "验证码已过期",
	ERROR_EMAIL_SEND:      "邮件发送失败",
	ERROR_EMAIL_CODE_WRONG: "邮箱验证码错误",
	ERROR_EMAIL_CODE_EXPIRE: "邮箱验证码已过期",
	ERROR_EMAIL_CONFIG:    "邮箱配置错误",
}

// GetErrMsg 根据错误码获取错误信息
func GetErrMsg(code int) string {
	msg, ok := codeMsg[code]
	if !ok {
		return codeMsg[ERROR]
	}
	return msg
}

// Error 响应错误信息
func Error(c *gin.Context, status int, code int, data interface{}) {
	c.JSON(status, gin.H{
		"code": code,
		"msg":  GetErrMsg(code),
		"data": data,
	})
}

// Success 响应成功信息
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": SUCCESS,
		"msg":  GetErrMsg(SUCCESS),
		"data": data,
	})
}