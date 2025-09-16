package models

import (
	"time"
)

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50;not null" json:"username"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 成功, 0: 失败
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	User UserModel `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Username    string    `gorm:"size:50;not null" json:"username"`
	Operation   string    `gorm:"size:100;not null" json:"operation"`
	Resource    string    `gorm:"size:100" json:"resource"`
	ResourceID  string    `gorm:"size:50" json:"resource_id"`
	OldValue    string    `gorm:"type:text" json:"old_value"`
	NewValue    string    `gorm:"type:text" json:"new_value"`
	IP          string    `gorm:"size:45" json:"ip"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	Status      int       `gorm:"default:1" json:"status"`
	Duration    int64     `gorm:"default:0" json:"duration"` // 操作耗时(毫秒)
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	User UserModel `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 设置表名
func (LoginLog) TableName() string {
	return "login_logs"
}

func (OperationLog) TableName() string {
	return "operation_logs"
}