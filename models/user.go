package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	BaseModel
	Username     string     `gorm:"size:64;uniqueIndex;not null;comment:用户名" json:"username" validate:"required,username"`
	Password     string     `gorm:"size:128;not null;comment:密码" json:"-" validate:"required,password"`
	Nickname     string     `gorm:"size:64;comment:昵称" json:"nickname"`
	Email        string     `gorm:"size:128;uniqueIndex;comment:邮箱" json:"email" validate:"omitempty,email"`
	Phone        string     `gorm:"size:16;uniqueIndex;comment:手机号" json:"phone" validate:"omitempty,phone"`
	Avatar       string     `gorm:"size:255;comment:头像" json:"avatar"`
	Status       int        `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	LastLoginAt  *time.Time `gorm:"type:datetime;comment:最后登录时间" json:"last_login_at"`
	LastLoginIP  string     `gorm:"size:64;comment:最后登录IP" json:"last_login_ip"`
	LoginCount   int        `gorm:"type:int;default:0;comment:登录次数" json:"login_count"`
	DepartmentID uint       `gorm:"comment:部门ID" json:"department_id"`
	DeptID       uint       `gorm:"comment:部门ID(别名)" json:"dept_id"`
	Gender       int        `gorm:"type:tinyint;default:0;comment:性别(0:未知,1:男,2:女)" json:"gender"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Roles        []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	IsAdmin      bool       `gorm:"type:tinyint;default:0;comment:是否管理员" json:"is_admin"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == 0 {
		u.Status = 1
	}
	return nil
}

// Department 部门模型
type Department struct {
	BaseModel
	Name     string       `gorm:"size:64;not null;comment:部门名称" json:"name" validate:"required"`
	ParentID uint         `gorm:"default:0;comment:上级部门ID" json:"parent_id"`
	Sort     int          `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status   int          `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 设置表名
func (Department) TableName() string {
	return "departments"
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string       `gorm:"size:64;uniqueIndex;not null;comment:角色名称" json:"name" validate:"required"`
	Key         string       `gorm:"size:64;uniqueIndex;not null;comment:角色标识" json:"key" validate:"required"`
	Description string       `gorm:"size:255;comment:角色描述" json:"description"`
	Status      int          `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Sort        int          `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Menus       []Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	BaseModel
	Name        string `gorm:"size:64;not null;comment:权限名称" json:"name" validate:"required"`
	Key         string `gorm:"size:64;uniqueIndex;not null;comment:权限标识" json:"key" validate:"required"`
	Description string `gorm:"size:255;comment:权限描述" json:"description"`
	Type        string `gorm:"size:32;default:'api';comment:权限类型(api,menu,button)" json:"type"`
	Method      string `gorm:"size:16;comment:请求方法" json:"method"`
	Path        string `gorm:"size:128;comment:请求路径" json:"path"`
	Component   string `gorm:"size:128;comment:组件路径" json:"component"`
	Icon        string `gorm:"size:64;comment:图标" json:"icon"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	ParentID    uint   `gorm:"default:0;comment:上级权限ID" json:"parent_id"`
	Sort        int    `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 设置表名
func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey;comment:用户ID" json:"user_id"`
	RoleID uint `gorm:"primaryKey;comment:角色ID" json:"role_id"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;comment:角色ID" json:"role_id"`
	PermissionID uint `gorm:"primaryKey;comment:权限ID" json:"permission_id"`
}

// TableName 设置表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
