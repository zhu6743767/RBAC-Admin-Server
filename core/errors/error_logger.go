package errors

import (
	"encoding/json"
	"runtime"
	"time"

	"rbac.admin/global"
)

// ErrorLogEntry 错误日志条目
type ErrorLogEntry struct {
	Timestamp   time.Time              `json:"@timestamp"`
	Level       string                 `json:"level"`
	Type        ErrorType              `json:"error_type"`
	Code        int                    `json:"error_code"`
	Message     string                 `json:"message"`
	Detail      string                 `json:"detail,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	TraceID     string                 `json:"trace_id,omitempty"`
	StackTrace  string                 `json:"stack_trace,omitempty"`
	Function    string                 `json:"function"`
	File        string                 `json:"file"`
	Line        int                    `json:"line"`
	UserID      string                 `json:"user_id,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	Endpoint    string                 `json:"endpoint,omitempty"`
	HTTPMethod  string                 `json:"http_method,omitempty"`
}

// LogError 记录错误到日志系统
func LogError(appError *AppError, requestContext ...map[string]interface{}) {
	if appError == nil {
		return
	}

	entry := ErrorLogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Type:      appError.Type,
		Code:      appError.Code,
		Message:   appError.Message,
		Detail:    appError.Detail,
		Context:   appError.Context,
		TraceID:   appError.TraceID,
	}

	// 添加上下文信息
	if len(requestContext) > 0 {
		ctx := requestContext[0]
		if userID, ok := ctx["user_id"].(string); ok {
			entry.UserID = userID
		}
		if requestID, ok := ctx["request_id"].(string); ok {
			entry.RequestID = requestID
		}
		if userAgent, ok := ctx["user_agent"].(string); ok {
			entry.UserAgent = userAgent
		}
		if ip, ok := ctx["ip"].(string); ok {
			entry.IP = ip
		}
		if endpoint, ok := ctx["endpoint"].(string); ok {
			entry.Endpoint = endpoint
		}
		if method, ok := ctx["http_method"].(string); ok {
			entry.HTTPMethod = method
		}
	}

	// 获取调用栈信息
	entry.Function, entry.File, entry.Line = getCallerInfo(3)
	
	// 获取堆栈跟踪
	entry.StackTrace = getStackTrace()

	// 序列化为JSON并记录
	logData, _ := json.Marshal(entry)
	global.Logger.Error(string(logData))
}

// LogWarning 记录警告
func LogWarning(message string, context map[string]interface{}) {
	entry := map[string]interface{}{
		"@timestamp": time.Now(),
		"level":      "WARNING",
		"message":    message,
		"context":    context,
	}
	
	logData, _ := json.Marshal(entry)
	global.Logger.Warn(string(logData))
}

// getCallerInfo 获取调用者信息
func getCallerInfo(skip int) (function, file string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", 0
	}
	
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", file, line
	}
	
	return fn.Name(), file, line
}

// getStackTrace 获取堆栈跟踪
func getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// DatabaseConnectionError 数据库连接错误记录
func LogDatabaseConnectionError(host string, port int, username string, err error, context map[string]interface{}) {
	appErr := DatabaseConnectionError(host, port, err)
	if username != "" {
		appErr.WithContext("username", username)
	}
	LogError(appErr, context)
}

// NetworkConnectionError 网络连接错误记录
func LogNetworkConnectionError(service, endpoint string, err error, context map[string]interface{}) {
	appErr := NetworkConnectionError(service, endpoint, err)
	LogError(appErr, context)
}

// ValidationError 参数验证错误记录
func LogValidationError(field string, value interface{}, context map[string]interface{}) {
	appErr := ValidationErrorWithField(field, value)
	LogError(appErr, context)
}