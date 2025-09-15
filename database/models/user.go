package models

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Role 角色模型
type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Resource    string    `gorm:"size:100" json:"resource"` // 资源名称
	Action      string    `gorm:"size:50" json:"action"`    // 操作类型: create, read, update, delete
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

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
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
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
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (LoginLog) TableName() string {
	return "login_logs"
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}