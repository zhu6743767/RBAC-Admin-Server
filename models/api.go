package models

// API API模型
type API struct {
	BaseModel
	Name        string `gorm:"size:64;not null;comment:API名称" json:"name" validate:"required"`
	Path        string `gorm:"size:128;not null;comment:API路径" json:"path" validate:"required"`
	Method      string `gorm:"size:16;not null;comment:请求方法" json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	Group       string `gorm:"size:64;comment:API分组" json:"group"`
	Description string `gorm:"size:255;comment:API描述" json:"description"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Sort        int    `gorm:"type:int;default:0;comment:排序" json:"sort"`
}

// TableName 设置表名
func (API) TableName() string {
	return "apis"
}

// Dict 字典模型
type Dict struct {
	BaseModel
	Name        string     `gorm:"size:64;not null;comment:字典名称" json:"name" validate:"required"`
	Key         string     `gorm:"size:64;uniqueIndex;not null;comment:字典标识" json:"key" validate:"required"`
	Description string     `gorm:"size:255;comment:字典描述" json:"description"`
	Status      int        `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Sort        int        `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Items       []DictItem `gorm:"foreignKey:DictID" json:"items,omitempty"`
}

// TableName 设置表名
func (Dict) TableName() string {
	return "dicts"
}

// DictItem 字典项模型
type DictItem struct {
	BaseModel
	DictID      uint   `gorm:"comment:字典ID" json:"dict_id"`
	Label       string `gorm:"size:64;not null;comment:标签" json:"label" validate:"required"`
	Value       string `gorm:"size:64;not null;comment:值" json:"value" validate:"required"`
	Description string `gorm:"size:255;comment:描述" json:"description"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Sort        int    `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Dict        Dict   `gorm:"foreignKey:DictID" json:"dict,omitempty"`
}

// TableName 设置表名
func (DictItem) TableName() string {
	return "dict_items"
}

// Config 配置模型
type Config struct {
	BaseModelNoDelete
	Key         string `gorm:"size:64;uniqueIndex;not null;comment:配置键" json:"key" validate:"required"`
	Value       string `gorm:"type:text;comment:配置值" json:"value"`
	Type        string `gorm:"size:32;default:'string';comment:配置类型(string,int,bool,json)" json:"type"`
	Description string `gorm:"size:255;comment:配置描述" json:"description"`
	IsSystem    int    `gorm:"type:tinyint;default:2;comment:是否系统配置(1:是,2:否)" json:"is_system"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Group       string `gorm:"size:64;comment:配置分组" json:"group"`
}

// TableName 设置表名
func (Config) TableName() string {
	return "configs"
}

// File 文件模型
type File struct {
	BaseModel
	Name       string `gorm:"size:128;not null;comment:文件名称" json:"name" validate:"required"`
	Path       string `gorm:"size:255;not null;comment:文件路径" json:"path" validate:"required"`
	Size       int64  `gorm:"not null;comment:文件大小" json:"size"`
	Type       string `gorm:"size:64;comment:文件类型" json:"type"`
	MimeType   string `gorm:"size:128;comment:MIME类型" json:"mime_type"`
	Extension  string `gorm:"size:32;comment:文件扩展名" json:"extension"`
	Hash       string `gorm:"size:64;comment:文件哈希" json:"hash"`
	Status     int    `gorm:"type:tinyint;default:1;comment:状态(1:正常,2:禁用)" json:"status"`
	Category   string `gorm:"size:64;comment:文件分类" json:"category"`
	UploadedBy uint   `gorm:"comment:上传用户ID" json:"uploaded_by"`
	User       User   `gorm:"foreignKey:UploadedBy" json:"user,omitempty"`
}

// TableName 设置表名
func (File) TableName() string {
	return "files"
}

// Log 日志模型
type Log struct {
	BaseModel
	UserID      uint   `gorm:"comment:用户ID" json:"user_id"`
	Username    string `gorm:"size:64;comment:用户名" json:"username"`
	IP          string `gorm:"size:64;comment:IP地址" json:"ip"`
	UserAgent   string `gorm:"size:255;comment:用户代理" json:"user_agent"`
	Method      string `gorm:"size:16;comment:请求方法" json:"method"`
	Path        string `gorm:"size:128;comment:请求路径" json:"path"`
	StatusCode  int    `gorm:"comment:状态码" json:"status_code"`
	RequestBody string `gorm:"type:text;comment:请求体" json:"request_body"`
	Response    string `gorm:"type:text;comment:响应内容" json:"response"`
	Latency     int64  `gorm:"comment:响应时间(毫秒)" json:"latency"`
	Error       string `gorm:"type:text;comment:错误信息" json:"error"`
	Module      string `gorm:"size:64;comment:模块" json:"module"`
	Action      string `gorm:"size:64;comment:操作" json:"action"`
	Description string `gorm:"size:255;comment:描述" json:"description"`
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 设置表名
func (Log) TableName() string {
	return "logs"
}
