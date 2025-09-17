package models

// Menu 菜单模型
type Menu struct {
	BaseModel
	Name       string   `gorm:"size:64;not null;comment:菜单名称" json:"name" validate:"required"`
	Path       string   `gorm:"size:128;comment:菜单路径" json:"path"`
	Component  string   `gorm:"size:128;comment:组件路径" json:"component"`
	Redirect   string   `gorm:"size:128;comment:重定向路径" json:"redirect"`
	Icon       string   `gorm:"size:64;comment:菜单图标" json:"icon"`
	Type       string   `gorm:"size:32;default:'menu';comment:菜单类型(menu,button)" json:"type"`
	Permission string   `gorm:"size:64;comment:权限标识" json:"permission"`
	Sort       int      `gorm:"type:int;default:0;comment:排序" json:"sort"`
	ParentID   uint     `gorm:"default:0;comment:上级菜单ID" json:"parent_id"`
	Status     int      `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Hidden     int      `gorm:"type:tinyint;default:2;comment:是否隐藏(1:是,2:否)" json:"hidden"`
	KeepAlive  int      `gorm:"type:tinyint;default:1;comment:是否缓存(1:是,2:否)" json:"keep_alive"`
	AlwaysShow int      `gorm:"type:tinyint;default:2;comment:是否总是显示(1:是,2:否)" json:"always_show"`
	Breadcrumb int      `gorm:"type:tinyint;default:1;comment:是否显示面包屑(1:是,2:否)" json:"breadcrumb"`
	Affix      int      `gorm:"type:tinyint;default:2;comment:是否固定(1:是,2:否)" json:"affix"`
	NoCache    int      `gorm:"type:tinyint;default:2;comment:是否不缓存(1:是,2:否)" json:"no_cache"`
	Roles      []Role   `gorm:"many2many:role_menus;" json:"roles,omitempty"`
	Children   []Menu   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Meta       MenuMeta `gorm:"embedded;embeddedPrefix:meta_" json:"meta"`
}

// TableName 设置表名
func (Menu) TableName() string {
	return "menus"
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title      string `gorm:"size:128;comment:菜单标题" json:"title"`
	Icon       string `gorm:"size:64;comment:菜单图标" json:"icon"`
	NoCache    bool   `gorm:"comment:是否不缓存" json:"no_cache"`
	Breadcrumb bool   `gorm:"comment:是否显示面包屑" json:"breadcrumb"`
	Affix      bool   `gorm:"comment:是否固定" json:"affix"`
	ActiveMenu string `gorm:"size:128;comment:高亮菜单" json:"active_menu"`
	Roles      string `gorm:"size:255;comment:角色权限" json:"roles"`
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	RoleID uint `gorm:"primaryKey;comment:角色ID" json:"role_id"`
	MenuID uint `gorm:"primaryKey;comment:菜单ID" json:"menu_id"`
}

// TableName 设置表名
func (RoleMenu) TableName() string {
	return "role_menus"
}
