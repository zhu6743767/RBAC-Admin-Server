package examples

import (
	"fmt"
	"rbac.admin/config"
	"rbac.admin/database"
	"rbac.admin/database/models"
	"gorm.io/gorm"
)

// 数据库工厂模式使用示例
func DatabaseUsageExample() {
	fmt.Println("=== 数据库工厂模式使用示例 ===")

	// 1. 加载配置
	cfg := config.DefaultConfig()
	cfg.Database.Type = "mysql"
	cfg.Database.Host = "192.168.10.199"
	cfg.Database.Port = 3306
	cfg.Database.Username = "root"
	cfg.Database.Password = "Zdj_7819!"
	cfg.Database.Database = "rbac_admin"

	// 2. 创建数据库管理器
	dbManager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	defer dbManager.Close()

	// 3. 获取数据库连接
	db := dbManager.GetDB()

	// 4. 执行迁移
	migrator := database.NewMigrator(db)
	if err := migrator.AutoMigrate(); err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		return
	}

	// 5. 初始化数据
	if err := migrator.SeedData(); err != nil {
		fmt.Printf("数据初始化失败: %v\n", err)
		return
	}

	// 6. 查询示例
	var users []models.UserModel
	if err := db.Find(&users).Error; err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		return
	}

	fmt.Printf("用户总数: %d\n", len(users))
	for _, user := range users {
		fmt.Printf("用户: %s (%s)\n", user.Username, user.Email)
	}

	// 7. 创建新用户示例
	newUser := models.UserModel{
		Username: "test_user",
		Email:    "test@example.com",
		Password: "hashed_password",
		Nickname: "测试用户",
		Status:   1,
	}
	if err := db.Create(&newUser).Error; err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}

	fmt.Printf("创建新用户成功: %s\n", newUser.Username)
}

// 不同数据库配置示例
func databaseConfigExamples() {
	examples := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "MySQL配置",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Type:     "mysql",
					Host:     "192.168.10.199",
					Port:     3306,
					Username: "root",
					Password: "Zdj_7819!",
					Database: "rbac_admin",
					Charset:  "utf8mb4",
				},
			},
		},
		{
			name: "PostgreSQL配置",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Type:     "postgres",
					Host:     "localhost",
					Port:     5432,
					Username: "postgres",
					Password: "",
					Database: "rbac_admin",
					SSLMode:  "disable",
				},
			},
		},
		{
			name: "SQLite配置",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Type: "sqlite",
					Path: "./rbac_admin.db",
				},
			},
		},
	}

	for _, example := range examples {
		fmt.Printf("\n=== %s ===\n", example.name)
		fmt.Printf("类型: %s\n", example.config.Database.Type)
		fmt.Printf("主机: %s\n", example.config.Database.Host)
		fmt.Printf("端口: %d\n", example.config.Database.Port)
		fmt.Printf("数据库: %s\n", example.config.Database.Database)
	}
}

// 事务处理示例
func transactionExample(db *gorm.DB) {
	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		fmt.Printf("事务开始失败: %v\n", tx.Error)
		return
	}

	// 在事务中执行操作
	user := models.UserModel{
		Username: "transaction_user",
		Email:    "transaction@example.com",
		Password: "hashed_password",
		Status:   1,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		fmt.Printf("事务创建用户失败: %v\n", err)
		return
	}

	role := models.RoleModel{
		Name:        "transaction_role",
		Description: "事务测试角色",
		Status:      1,
	}

	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		fmt.Printf("事务创建角色失败: %v\n", err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("事务提交失败: %v\n", err)
		return
	}

	fmt.Println("事务操作成功完成")
}

// 关联查询示例
func associationExample(db *gorm.DB) {
	// 预加载关联数据
	var user models.UserModel
	if err := db.Preload("Roles").Preload("Roles.Menus").First(&user, 1).Error; err != nil {
		fmt.Printf("查询用户关联数据失败: %v\n", err)
		return
	}

	fmt.Printf("用户: %s\n", user.Username)
	fmt.Printf("角色数量: %d\n", len(user.Roles))
	
	for _, role := range user.Roles {
		fmt.Printf("  角色: %s - %s\n", role.Name, role.Description)
		fmt.Printf("  菜单数量: %d\n", len(role.Menus))
		for _, menu := range role.Menus {
			fmt.Printf("    菜单: %s - %s\n", menu.Name, menu.Title)
		}
	}
}

// 分页查询示例
func paginationExample(db *gorm.DB) {
	var users []models.UserModel
	var total int64

	// 查询总数
	db.Model(&models.UserModel{}).Count(&total)

	// 分页查询
	page := 1
	pageSize := 10
	offset := (page - 1) * pageSize

	if err := db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		fmt.Printf("分页查询失败: %v\n", err)
		return
	}

	fmt.Printf("总记录数: %d\n", total)
	fmt.Printf("当前页: %d\n", page)
	fmt.Printf("每页数量: %d\n", pageSize)
	fmt.Printf("返回记录数: %d\n", len(users))
}

// 复杂查询示例
func complexQueryExample(db *gorm.DB) {
	// 查询用户的角色和菜单
	var results []struct {
		Username string
		RoleName string
		MenuName string
	}

	db.Table("sys_users u").
		Select("u.username, r.name as role_name, m.name as menu_name").
		Joins("JOIN sys_user_roles ur ON u.id = ur.user_id").
		Joins("JOIN sys_roles r ON ur.role_id = r.id").
		Joins("JOIN sys_role_menus rm ON r.id = rm.role_id").
		Joins("JOIN sys_menus m ON rm.menu_id = m.id").
		Where("u.status = ?", 1).
		Where("r.status = ?", 1).
		Scan(&results)

	fmt.Println("用户-角色-菜单关联查询结果:")
	for _, result := range results {
		fmt.Printf("用户: %s, 角色: %s, 菜单: %s\n", result.Username, result.RoleName, result.MenuName)
	}
}