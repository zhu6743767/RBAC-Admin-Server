package flags

import (
	"flag"
)

// Mode 操作模式枚举
const (
	ModeServer   = "server"  // 启动服务器模式
	ModeDatabase = "db"      // 数据库操作模式
	ModeUser     = "user"    // 用户管理模式
)

// DatabaseType 数据库操作类型枚举
const (
	DatabaseMigrate   = "migrate"    // 数据库迁移
	DatabaseSeed      = "seed"       // 数据库种子数据
	DatabaseReset     = "reset"      // 重置数据库
)

// UserType 用户操作类型枚举
const (
	UserCreateAdmin = "create"      // 创建管理员用户
	UserList        = "list"        // 列出用户
	UserReset       = "reset"       // 重置用户密码
)

// CommandLineArgs 命令行参数结构体
type CommandLineArgs struct {
	Mode      string // 操作模式
	Type      string // 操作类型
	Config    string // 配置文件路径
	Username  string // 用户名
	Password  string // 密码
}

// ParseCommandLineArgs 解析命令行参数
func ParseCommandLineArgs() CommandLineArgs {
	// 定义命令行参数
	mode := flag.String("m", ModeServer, "操作模式: server(启动服务器), db(数据库操作), user(用户管理)")
	typeArg := flag.String("t", "", "操作类型: 对于db模式可以是migrate/seed/reset, 对于user模式可以是create/list/reset")
	config := flag.String("settings", "settings.yaml", "配置文件路径")
	username := flag.String("username", "admin", "用户名")
	password := flag.String("password", "", "密码")

	// 解析命令行参数
	flag.Parse()

	return CommandLineArgs{
		Mode:      *mode,
		Type:      *typeArg,
		Config:    *config,
		Username:  *username,
		Password:  *password,
	}
}