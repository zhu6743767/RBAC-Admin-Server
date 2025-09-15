package errors

import (
	"fmt"
	"net/http"
)

// ErrorType 错误类型枚举
type ErrorType string

const (
	// 网络相关错误
	NetworkError    ErrorType = "NETWORK_ERROR"
	DatabaseError   ErrorType = "DATABASE_ERROR"
	RedisError      ErrorType = "REDIS_ERROR"
	ExternalAPIError ErrorType = "EXTERNAL_API_ERROR"
	
	// 业务逻辑错误
	ValidationError ErrorType = "VALIDATION_ERROR"
	AuthError       ErrorType = "AUTHENTICATION_ERROR"
	AuthorizationError ErrorType = "AUTHORIZATION_ERROR"
	BusinessError   ErrorType = "BUSINESS_ERROR"
	
	// 系统错误
	SystemError     ErrorType = "SYSTEM_ERROR"
	ConfigError     ErrorType = "CONFIG_ERROR"
	
	// 用户操作错误
	UserNotFound    ErrorType = "USER_NOT_FOUND"
	InvalidInput    ErrorType = "INVALID_INPUT"
	DuplicateEntry  ErrorType = "DUPLICATE_ENTRY"
)

// AppError 统一错误结构
type AppError struct {
	Type    ErrorType `json:"type"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	TraceID string    `json:"trace_id,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Message, e.Detail)
}

// NewError 创建新的应用错误
func NewError(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    getHTTPStatus(errorType),
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// NewErrorWithDetail 创建带详细信息的错误
func NewErrorWithDetail(errorType ErrorType, message, detail string) *AppError {
	err := NewError(errorType, message)
	err.Detail = detail
	return err
}

// WithContext 添加上下文信息
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// WithTraceID 添加追踪ID
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

// getHTTPStatus 根据错误类型返回HTTP状态码
func getHTTPStatus(errorType ErrorType) int {
	switch errorType {
	case NetworkError, DatabaseError, RedisError, ExternalAPIError:
		return http.StatusServiceUnavailable
	case ValidationError, InvalidInput:
		return http.StatusBadRequest
	case AuthError:
		return http.StatusUnauthorized
	case AuthorizationError:
		return http.StatusForbidden
	case UserNotFound:
		return http.StatusNotFound
	case DuplicateEntry:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// DatabaseError 数据库连接错误示例
func DatabaseConnectionError(host string, port int, err error) *AppError {
	return NewErrorWithDetail(DatabaseError, "数据库连接失败", err.Error()).
		WithContext("host", host).
		WithContext("port", port).
		WithContext("suggestion", "请检查数据库地址、端口、用户名和密码是否正确")
}

// NetworkError 网络错误示例
func NetworkConnectionError(service string, endpoint string, err error) *AppError {
	return NewErrorWithDetail(NetworkError, "网络连接失败", err.Error()).
		WithContext("service", service).
		WithContext("endpoint", endpoint)
}

// ValidationError 参数验证错误示例
func ValidationErrorWithField(field string, value interface{}) *AppError {
	return NewError(ValidationError, "参数验证失败").
		WithContext("field", field).
		WithContext("value", value)
}