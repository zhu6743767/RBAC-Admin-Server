package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// Validate 全局验证器实例
	Validate *validator.Validate
)

// InitValidator 初始化验证器
// 注册自定义验证规则和错误消息
func InitValidator() error {
	Validate = validator.New()

	// 注册自定义验证规则
	if err := registerCustomValidations(); err != nil {
		return fmt.Errorf("注册自定义验证规则失败: %v", err)
	}

	// 注册自定义错误消息
	registerCustomMessages()

	return nil
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations() error {
	// 验证手机号
	if err := Validate.RegisterValidation("phone", validatePhone); err != nil {
		return err
	}

	// 验证用户名（字母、数字、下划线，3-20位）
	if err := Validate.RegisterValidation("username", validateUsername); err != nil {
		return err
	}

	// 验证密码强度
	if err := Validate.RegisterValidation("password", validatePassword); err != nil {
		return err
	}

	// 验证中文姓名
	if err := Validate.RegisterValidation("chinese_name", validateChineseName); err != nil {
		return err
	}

	// 验证身份证号
	if err := Validate.RegisterValidation("id_card", validateIDCard); err != nil {
		return err
	}

	return nil
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 空值由required规则处理
	}

	// 中国手机号正则表达式
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true // 空值由required规则处理
	}

	// 用户名正则表达式（字母、数字、下划线，3-20位）
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return usernameRegex.MatchString(username)
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true // 空值由required规则处理
	}

	// 密码强度要求：8-20位，包含大小写字母、数字、特殊字符中的至少3种
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	var count int
	// 检查是否包含大写字母
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		count++
	}
	// 检查是否包含小写字母
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		count++
	}
	// 检查是否包含数字
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		count++
	}
	// 检查是否包含特殊字符
	if regexp.MustCompile(`[!@#$%^&*()_+=\-[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		count++
	}

	return count >= 3
}

// validateChineseName 验证中文姓名
func validateChineseName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" {
		return true // 空值由required规则处理
	}

	// 中文姓名正则表达式（2-10个汉字）
	chineseNameRegex := regexp.MustCompile(`^[\u4e00-\u9fa5]{2,10}$`)
	return chineseNameRegex.MatchString(name)
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if idCard == "" {
		return true // 空值由required规则处理
	}

	// 18位身份证号正则表达式
	idCardRegex := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	return idCardRegex.MatchString(idCard)
}

// registerCustomMessages 注册自定义错误消息
func registerCustomMessages() {
	// 这里可以注册自定义的错误消息
	// 实际项目中可以根据需要国际化
}

// ValidateStruct 验证结构体
func ValidateStruct(data interface{}) error {
	if Validate == nil {
		return fmt.Errorf("验证器未初始化")
	}
	return Validate.Struct(data)
}

// ValidateVar 验证单个变量
func ValidateVar(value interface{}, tag string) error {
	if Validate == nil {
		return fmt.Errorf("验证器未初始化")
	}
	return Validate.Var(value, tag)
}

// FormatValidationError 格式化验证错误
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s不能为空", e.Field())
			case "email":
				errors[field] = "请输入有效的邮箱地址"
			case "phone":
				errors[field] = "请输入有效的手机号"
			case "username":
				errors[field] = "用户名必须是3-20位的字母、数字或下划线"
			case "password":
				errors[field] = "密码必须是8-20位，包含大小写字母、数字、特殊字符中的至少3种"
			case "chinese_name":
				errors[field] = "请输入2-10个汉字的中文姓名"
			case "id_card":
				errors[field] = "请输入有效的身份证号"
			case "min":
				errors[field] = fmt.Sprintf("%s长度不能少于%s个字符", e.Field(), e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s长度不能超过%s个字符", e.Field(), e.Param())
			case "len":
				errors[field] = fmt.Sprintf("%s长度必须是%s个字符", e.Field(), e.Param())
			default:
				errors[field] = fmt.Sprintf("%s格式不正确", e.Field())
			}
		}
	}

	return errors
}
