package examples

import (
	"fmt"
	"rbac.admin/config"
	"rbac.admin/database"
	"rbac.admin/database/models"
)

func MainExample() {
	fmt.Println("=== RBAC管理员服务器使用示例 ===")

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
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		return
	}

	fmt.Printf("用户总数: %d\n", len(users))
	for _, user := range users {
		fmt.Printf("用户: %s (%s)\n", user.Username, user.Email)
	}

	fmt.Println("示例运行完成！")
}