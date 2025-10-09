package flags

import (
	"fmt"
	"rbac_admin_server/config"
	"rbac_admin_server/core"
	"rbac_admin_server/core/init_gorm"
	"rbac_admin_server/global"
	"rbac_admin_server/models"
	"rbac_admin_server/utils"

	"gorm.io/gorm"
)

// HandleCommandLineArgs 处理命令行参数
func HandleCommandLineArgs(args CommandLineArgs) error {
	// 加载配置
	var err error
	global.Config, err = config.Load(args.Config)
	if err != nil {
		return fmt.Errorf("配置文件加载失败: %v", err)
	}

	// 初始化日志
	if err := core.InitLogger(&global.Config.Log); err != nil {
		return fmt.Errorf("日志系统初始化失败: %v", err)
	}

	// 根据模式处理
	switch args.Mode {
	case ModeServer:
		// 启动服务器模式 - 在main.go中处理
		return nil
	case ModeDatabase:
		// 数据库操作模式
		return handleDatabaseCommand(args.Type)
	case ModeUser:
		// 用户管理模式
		return handleUserCommand(args.Type, args.Username, args.Password)
	default:
		return fmt.Errorf("不支持的操作模式: %s", args.Mode)
	}
}

// handleDatabaseCommand 处理数据库相关命令
func handleDatabaseCommand(typeArg string) error {
	// 初始化数据库
	db, err := init_gorm.InitGorm()
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %v", err)
	}
	global.DB = db

	switch typeArg {
	case DatabaseMigrate:
		// 数据库迁移
		if err := core.AutoMigrateModels(); err != nil {
			return fmt.Errorf("数据库迁移失败: %v", err)
		}
		global.Logger.Info("✅ 数据库迁移成功")

		// 初始化基础数据
		if err := initBaseData(db); err != nil {
			return fmt.Errorf("初始化基础数据失败: %v", err)
		}
		global.Logger.Info("✅ 基础数据初始化成功")

	case DatabaseSeed:
		// 数据库种子数据
		if err := initBaseData(db); err != nil {
			return fmt.Errorf("初始化种子数据失败: %v", err)
		}
		global.Logger.Info("✅ 种子数据初始化成功")

	case DatabaseReset:
		// 重置数据库
		if err := resetDatabase(db); err != nil {
			return fmt.Errorf("重置数据库失败: %v", err)
		}
		global.Logger.Info("✅ 数据库重置成功")

	default:
		return fmt.Errorf("不支持的数据库操作类型: %s", typeArg)
	}

	return nil
}

// handleUserCommand 处理用户相关命令
func handleUserCommand(typeArg, username, password string) error {
	// 初始化数据库
	db, err := init_gorm.InitGorm()
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %v", err)
	}
	global.DB = db

	switch typeArg {
	case UserCreateAdmin:
		// 创建管理员用户
		if password == "" {
			password = "admin123" // 默认密码
			global.Logger.Warn("使用默认密码创建管理员用户，请尽快修改")
		}

		if err := createAdminUser(db, username, password); err != nil {
			return fmt.Errorf("创建管理员用户失败: %v", err)
		}
		global.Logger.Infof("✅ 管理员用户 %s 创建成功", username)

	case UserList:
		// 列出用户
		if err := listUsers(db); err != nil {
			return fmt.Errorf("列出用户失败: %v", err)
		}

	case UserReset:
		// 重置用户密码
		if password == "" {
			password = "123456" // 默认密码
			global.Logger.Warn("使用默认密码重置用户密码，请尽快修改")
		}

		if err := resetUserPassword(db, username, password); err != nil {
			return fmt.Errorf("重置用户密码失败: %v", err)
		}
		global.Logger.Infof("✅ 用户 %s 密码重置成功", username)

	default:
		return fmt.Errorf("不支持的用户操作类型: %s", typeArg)
	}

	return nil
}

// initBaseData 初始化基础数据
func initBaseData(db *gorm.DB) error {
	// 这里可以初始化一些基础数据，如默认角色、权限等
	// 例如创建超级管理员角色
	var superAdminRole models.Role
	if err := db.Where("name = ?", "超级管理员").First(&superAdminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			superAdminRole = models.Role{
				Name:        "超级管理员",
				Description: "系统超级管理员，拥有所有权限",
				Status:      1,
			}
			if err := db.Create(&superAdminRole).Error; err != nil {
				return fmt.Errorf("创建超级管理员角色失败: %v", err)
			}
			global.Logger.Info("✅ 创建超级管理员角色成功")
		} else {
			return fmt.Errorf("查询超级管理员角色失败: %v", err)
		}
	} else {
		global.Logger.Info("超级管理员角色已存在，跳过创建")
	}

	return nil
}

// resetDatabase 重置数据库
func resetDatabase(db *gorm.DB) error {
	// 获取所有表名
	var tables []string
	db.Raw("SHOW TABLES").Scan(&tables)

	// 删除所有表
	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS `" + table + "`").Error; err != nil {
			return fmt.Errorf("删除表 %s 失败: %v", table, err)
		}
	}
	global.Logger.Info("✅ 所有表已删除")

	// 重新创建表
	if err := core.AutoMigrateModels(); err != nil {
		return fmt.Errorf("重新创建表失败: %v", err)
	}
	global.Logger.Info("✅ 所有表已重新创建")

	// 重新初始化基础数据
	if err := initBaseData(db); err != nil {
		return fmt.Errorf("重新初始化基础数据失败: %v", err)
	}

	return nil
}

// createAdminUser 创建管理员用户
func createAdminUser(db *gorm.DB, username, password string) error {
	// 检查用户是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("用户 %s 已存在", username)
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询用户失败: %v", err)
	}

	// 创建管理员用户
	hashedPassword := utils.MakePassword(password)

	adminUser := models.User{
		Username:  username,
		Password:  hashedPassword,
		Nickname:  "系统管理员",
		Avatar:    "/uploads/avatar/default.png",
		Email:     "admin@example.com",
		Phone:     "13800138000",
		Status:    1,
		IsAdmin:   true,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("创建管理员用户失败: %v", err)
	}

	// 分配超级管理员角色
	var superAdminRole models.Role
	if err := db.Where("name = ?", "超级管理员").First(&superAdminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			global.Logger.Warn("超级管理员角色不存在，跳过角色分配")
		} else {
			return fmt.Errorf("查询超级管理员角色失败: %v", err)
		}
	} else {
		userRole := models.UserRole{
			UserID: adminUser.ID,
			RoleID: superAdminRole.ID,
		}
		if err := db.Create(&userRole).Error; err != nil {
			return fmt.Errorf("分配超级管理员角色失败: %v", err)
		}
		global.Logger.Info("✅ 已为管理员用户分配超级管理员角色")
	}

	return nil
}

// listUsers 列出用户
func listUsers(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("查询用户列表失败: %v", err)
	}

	global.Logger.Info("\n用户列表:")
	global.Logger.Info("ID\t用户名\t\t昵称\t\t邮箱\t\t\t状态\t是否管理员")
	global.Logger.Info("--------------------------------------------------------------------------------------------------------")
	for _, user := range users {
		global.Logger.Infof("%d\t%s\t\t%s\t\t%s\t\t%d\t\t%t",
			user.ID, user.Username, user.Nickname, user.Email, user.Status, user.IsAdmin)
	}
	global.Logger.Info("--------------------------------------------------------------------------------------------------------")
	global.Logger.Infof("总计 %d 个用户", len(users))

	return nil
}

// resetUserPassword 重置用户密码
func resetUserPassword(db *gorm.DB, username, password string) error {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户 %s 不存在", username)
		}
		return fmt.Errorf("查询用户失败: %v", err)
	}

	// 更新密码
	hashedPassword := utils.MakePassword(password)

	user.Password = hashedPassword
	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("更新用户密码失败: %v", err)
	}

	return nil
}