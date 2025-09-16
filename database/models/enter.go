package models

import (
	"time"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;comment:删除时间(软删除)" json:"deleted_at,omitempty"`
}

// UserModel 用户模型 - 存储用户信息
type UserModel struct {
	BaseModel
	Username string      `gorm:"column:username;type:varchar(32);uniqueIndex;not null;comment:用户名" json:"username"`
	Nickname string      `gorm:"column:nickname;type:varchar(64);comment:昵称" json:"nickname"`
	Avatar   string      `gorm:"column:avatar;type:varchar(512);comment:头像URL" json:"avatar"`
	Email    string      `gorm:"column:email;type:varchar(128);uniqueIndex;comment:邮箱地址" json:"email"`
	Password string      `gorm:"column:password;type:varchar(256);not null;comment:密码哈希" json:"-"` // 禁止JSON序列化
	Status   int8        `gorm:"column:status;type:tinyint;default:1;comment:状态(1正常 2禁用)" json:"status"`
	Roles    []RoleModel `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID;comment:用户角色关联" json:"roles,omitempty"`
}

// RoleModel 角色模型 - 存储角色信息
type RoleModel struct {
	BaseModel
	Name        string      `gorm:"column:name;type:varchar(32);uniqueIndex;not null;comment:角色名称" json:"name"`
	Description string      `gorm:"column:description;type:varchar(255);comment:角色描述" json:"description"`
	Sort        int         `gorm:"column:sort;type:int;default:0;comment:排序权重" json:"sort"`
	Status      int8        `gorm:"column:status;type:tinyint;default:1;comment:状态(1启用 2禁用)" json:"status"`
	Users       []UserModel `gorm:"many2many:user_roles;joinForeignKey:RoleID;joinReferences:UserID;comment:角色用户关联" json:"-"` // 禁止JSON序列化
	Menus       []MenuModel `gorm:"many2many:role_menus;joinForeignKey:RoleID;joinReferences:MenuID;comment:角色菜单关联" json:"menus,omitempty"`
	Apis        []ApiModel  `gorm:"many2many:role_apis;joinForeignKey:RoleID;joinReferences:ApiID;comment:角色接口关联" json:"apis,omitempty"`
}

// MenuModel 菜单模型 - 存储菜单信息
type MenuModel struct {
	BaseModel
	Name       string      `gorm:"column:name;type:varchar(32);uniqueIndex;not null;comment:菜单名称" json:"name"`
	Path       string      `gorm:"column:path;type:varchar(128);comment:菜单路径" json:"path"`
	Component  string      `gorm:"column:component;type:varchar(128);comment:组件路径" json:"component"`
	Icon       string      `gorm:"column:icon;type:varchar(64);comment:图标" json:"icon"`
	Title      string      `gorm:"column:title;type:varchar(32);comment:菜单标题" json:"title"`
	Type       int8        `gorm:"column:type;type:tinyint;default:1;comment:类型(1目录 2菜单 3按钮)" json:"type"`
	Sort       int         `gorm:"column:sort;type:int;default:0;comment:排序权重" json:"sort"`
	Status     int8        `gorm:"column:status;type:tinyint;default:1;comment:状态(1启用 2禁用)" json:"status"`
	ParentID   *uint       `gorm:"column:parent_id;index;comment:父菜单ID" json:"parent_id,omitempty"`
	Children   []MenuModel `gorm:"foreignKey:ParentID;references:ID;comment:子菜单" json:"children,omitempty"`
	Roles      []RoleModel `gorm:"many2many:role_menus;joinForeignKey:MenuID;joinReferences:RoleID;comment:菜单角色关联" json:"-"`
}

// ApiModel API接口模型 - 存储接口信息
type ApiModel struct {
	BaseModel
	Name        string      `gorm:"column:name;type:varchar(32);uniqueIndex;not null;comment:接口名称" json:"name"`
	Path        string      `gorm:"column:path;type:varchar(128);uniqueIndex:idx_path_method;not null;comment:接口路径" json:"path"`
	Method      string      `gorm:"column:method;type:varchar(10);uniqueIndex:idx_path_method;not null;comment:请求方法(GET POST PUT DELETE)" json:"method"`
	Group       string      `gorm:"column:group;type:varchar(32);comment:接口分组" json:"group"`
	Description string      `gorm:"column:description;type:varchar(255);comment:接口描述" json:"description"`
	Status      int8        `gorm:"column:status;type:tinyint;default:1;comment:状态(1启用 2禁用)" json:"status"`
	Roles       []RoleModel `gorm:"many2many:role_apis;joinForeignKey:ApiID;joinReferences:RoleID;comment:接口角色关联" json:"-"`
}

// UserRole 用户角色关联表
type UserRole struct {
	BaseModel
	UserID uint `gorm:"column:user_id;index;not null;comment:用户ID" json:"user_id"`
	RoleID uint `gorm:"column:role_id;index;not null;comment:角色ID" json:"role_id"`
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	BaseModel
	RoleID uint `gorm:"column:role_id;index;not null;comment:角色ID" json:"role_id"`
	MenuID uint `gorm:"column:menu_id;index;not null;comment:菜单ID" json:"menu_id"`
}

// RoleApi 角色接口关联表
type RoleApi struct {
	BaseModel
	RoleID uint `gorm:"column:role_id;index;not null;comment:角色ID" json:"role_id"`
	ApiID  uint `gorm:"column:api_id;index;not null;comment:接口ID" json:"api_id"`
}

// 表名方法
func (UserModel) TableName() string  { return "sys_users" }
func (RoleModel) TableName() string  { return "sys_roles" }
func (MenuModel) TableName() string  { return "sys_menus" }
func (ApiModel) TableName() string   { return "sys_apis" }
func (UserRole) TableName() string   { return "sys_user_roles" }
func (RoleMenu) TableName() string   { return "sys_role_menus" }
func (RoleApi) TableName() string    { return "sys_role_apis" }
