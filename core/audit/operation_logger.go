package audit

import (
	"encoding/json"
	"time"

	"rbac.admin/core/logger"
)

// OperationType 操作类型
type OperationType string

const (
	UserRegister   OperationType = "USER_REGISTER"
	UserLogin      OperationType = "USER_LOGIN"
	UserLogout     OperationType = "USER_LOGOUT"
	PasswordChange OperationType = "PASSWORD_CHANGE"
	PasswordReset  OperationType = "PASSWORD_RESET"
	ProfileUpdate  OperationType = "PROFILE_UPDATE"
	RoleAssignment OperationType = "ROLE_ASSIGNMENT"
	PermissionGrant OperationType = "PERMISSION_GRANT"
	DataCreate     OperationType = "DATA_CREATE"
	DataUpdate     OperationType = "DATA_UPDATE"
	DataDelete     OperationType = "DATA_DELETE"
	DataQuery      OperationType = "DATA_QUERY"
	SystemConfig   OperationType = "SYSTEM_CONFIG"
	APIAccess      OperationType = "API_ACCESS"
)

// OperationStatus 操作状态
type OperationStatus string

const (
	Success OperationStatus = "SUCCESS"
	Failed  OperationStatus = "FAILED"
	Partial OperationStatus = "PARTIAL"
)

// OperationLog 操作日志结构
type OperationLog struct {
	Timestamp    time.Time              `json:"@timestamp"`
	Operation    OperationType          `json:"operation"`
	Status       OperationStatus        `json:"status"`
	UserID       string                 `json:"user_id"`
	Username     string                 `json:"username,omitempty"`
	IP           string                 `json:"ip"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	RequestID    string                 `json:"request_id"`
	Endpoint     string                 `json:"endpoint"`
	HTTPMethod   string                 `json:"http_method"`
	Duration     int64                  `json:"duration_ms"`
	RequestBody  map[string]interface{} `json:"request_body,omitempty"`
	ResponseBody map[string]interface{} `json:"response_body,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	BeforeData   map[string]interface{} `json:"before_data,omitempty"`
	AfterData    map[string]interface{} `json:"after_data,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// OperationLogger 操作日志记录器
type OperationLogger struct {
	logger logger.Logger
}

// NewOperationLogger 创建操作日志记录器
func NewOperationLogger(log logger.Logger) *OperationLogger {
	return &OperationLogger{logger: log}
}

// Log 记录操作日志
func (ol *OperationLogger) Log(log *OperationLog) {
	if log == nil {
		return
	}
	
	// 设置时间戳
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	
	// 序列化为JSON
	logData, err := json.Marshal(log)
	if err != nil {
		ol.logger.Error("Failed to marshal operation log: %v", err)
		return
	}
	
	// 根据操作状态选择日志级别
	switch log.Status {
	case Success:
		ol.logger.Info(string(logData))
	case Failed:
		ol.logger.Error(string(logData))
	case Partial:
		ol.logger.Warn(string(logData))
	}
}

// LogUserOperation 快捷记录用户操作
func (ol *OperationLogger) LogUserOperation(
	operation OperationType,
	userID, username string,
	ip, endpoint, method string,
	status OperationStatus,
	errorMsg string,
	metadata map[string]interface{},
) {
	log := &OperationLog{
		Operation:    operation,
		Status:       status,
		UserID:       userID,
		Username:     username,
		IP:           ip,
		Endpoint:     endpoint,
		HTTPMethod:   method,
		ErrorMessage: errorMsg,
		Metadata:     metadata,
	}
	
	ol.Log(log)
}

// LogAuthOperation 记录认证相关操作
func (ol *OperationLogger) LogAuthOperation(
	operation OperationType,
	userID, username string,
	ip string,
	status OperationStatus,
	metadata map[string]interface{},
) {
	log := &OperationLog{
		Operation:  operation,
		Status:     status,
		UserID:     userID,
		Username:   username,
		IP:         ip,
		Metadata:   metadata,
	}
	
	ol.Log(log)
}

// LogDataOperation 记录数据操作
func (ol *OperationLogger) LogDataOperation(
	operation OperationType,
	userID string,
	endpoint string,
	before, after map[string]interface{},
	status OperationStatus,
	metadata map[string]interface{},
) {
	log := &OperationLog{
		Operation:  operation,
		Status:     status,
		UserID:     userID,
		Endpoint:   endpoint,
		BeforeData: before,
		AfterData:  after,
		Metadata:   metadata,
	}
	
	ol.Log(log)
}

// LogUserRegister 记录用户注册
func (ol *OperationLogger) LogUserRegister(userID, username, ip string, status OperationStatus, metadata map[string]interface{}) {
	ol.LogAuthOperation(UserRegister, userID, username, ip, status, metadata)
}

// LogUserLogin 记录用户登录
func (ol *OperationLogger) LogUserLogin(userID, username, ip string, status OperationStatus, metadata map[string]interface{}) {
	ol.LogAuthOperation(UserLogin, userID, username, ip, status, metadata)
}

// LogPasswordChange 记录密码修改
func (ol *OperationLogger) LogPasswordChange(userID, username, ip string, status OperationStatus, metadata map[string]interface{}) {
	ol.LogAuthOperation(PasswordChange, userID, username, ip, status, metadata)
}

// LogPasswordReset 记录密码重置
func (ol *OperationLogger) LogPasswordReset(userID, username, ip string, status OperationStatus, metadata map[string]interface{}) {
	ol.LogAuthOperation(PasswordReset, userID, username, ip, status, metadata)
}