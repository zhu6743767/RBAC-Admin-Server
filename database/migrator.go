package database

import (
	"rbac.admin/database/models"

	"gorm.io/gorm"
)

// Migrator 数据库迁移器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrate 自动迁移所有模型
func (m *Migrator) AutoMigrate() error {
	modelList := []interface{}{
		&models.UserModel{},
		&models.RoleModel{},
		&models.MenuModel{},
		&models.ApiModel{},
		&models.UserRole{},
		&models.RoleMenu{},
		&models.RoleApi{},
		&models.LoginLog{},
		&models.OperationLog{},
	}

	return m.db.AutoMigrate(modelList...)
}

// MigrateWithOptions 带选项的迁移
func (m *Migrator) MigrateWithOptions(dropIndex bool, modelList ...interface{}) error {
	if len(modelList) == 0 {
		modelList = []interface{}{
			&models.UserModel{},
			&models.RoleModel{},
			&models.MenuModel{},
			&models.ApiModel{},
			&models.UserRole{},
			&models.RoleMenu{},
			&models.RoleApi{},
			&models.LoginLog{},
			&models.OperationLog{},
		}
	}

	if dropIndex {
		// 删除索引后重新创建
		for _, model := range modelList {
			if err := m.db.Migrator().DropTable(model); err != nil {
				// 注意：这里暂时移除日志，避免循环依赖
				// 在实际应用中，可以通过构造函数注入日志实例
			}
		}
	}

	return m.db.AutoMigrate(modelList...)
}

// SeedData 初始化基础数据
func (m *Migrator) SeedData() error {
	// 检查是否已经存在数据
	var count int64
	m.db.Model(&models.UserModel{}).Count(&count)
	if count > 0 {
		// 暂时移除日志，避免循环依赖
		return nil
	}

	// 创建默认API权限
	apis := []models.ApiModel{
		{Name: "创建用户", Path: "/api/users", Method: "POST", Group: "用户管理", Description: "创建用户"},
		{Name: "查看用户", Path: "/api/users", Method: "GET", Group: "用户管理", Description: "查看用户"},
		{Name: "更新用户", Path: "/api/users/:id", Method: "PUT", Group: "用户管理", Description: "更新用户"},
		{Name: "删除用户", Path: "/api/users/:id", Method: "DELETE", Group: "用户管理", Description: "删除用户"},
		{Name: "创建角色", Path: "/api/roles", Method: "POST", Group: "角色管理", Description: "创建角色"},
		{Name: "查看角色", Path: "/api/roles", Method: "GET", Group: "角色管理", Description: "查看角色"},
		{Name: "更新角色", Path: "/api/roles/:id", Method: "PUT", Group: "角色管理", Description: "更新角色"},
		{Name: "删除角色", Path: "/api/roles/:id", Method: "DELETE", Group: "角色管理", Description: "删除角色"},
		{Name: "查看菜单", Path: "/api/menus", Method: "GET", Group: "菜单管理", Description: "查看菜单"},
		{Name: "创建菜单", Path: "/api/menus", Method: "POST", Group: "菜单管理", Description: "创建菜单"},
	}

	if err := m.db.Create(&apis).Error; err != nil {
		return err
	}

	// 创建默认菜单
	menus := []models.MenuModel{
		{Name: "系统管理", Path: "/system", Component: "Layout", Sort: 1, Status: 1, Type: 0, Title: "系统管理", Icon: "system"},
		{Name: "用户管理", Path: "/system/users", Component: "system/user/index", Sort: 1, Status: 1, Type: 1, Title: "用户管理", Icon: "user"},
		{Name: "角色管理", Path: "/system/roles", Component: "system/role/index", Sort: 2, Status: 1, Type: 1, Title: "角色管理", Icon: "role"},
		{Name: "菜单管理", Path: "/system/menus", Component: "system/menu/index", Sort: 3, Status: 1, Type: 1, Title: "菜单管理", Icon: "menu"},
	}

	if err := m.db.Create(&menus).Error; err != nil {
		return err
	}

	// 创建管理员角色
	adminRole := models.RoleModel{
		Name:        "admin",
		Description: "系统管理员",
		Status:      1,
	}
	if err := m.db.Create(&adminRole).Error; err != nil {
		return err
	}

	// 为管理员角色分配所有API权限
	var allApis []models.ApiModel
	m.db.Find(&allApis)
	for _, api := range allApis {
		roleApi := models.RoleApi{
			RoleID: adminRole.ID,
			ApiID:  api.ID,
		}
		m.db.Create(&roleApi)
	}

	// 为管理员角色分配所有菜单权限
	var allMenus []models.MenuModel
	m.db.Find(&allMenus)
	for _, menu := range allMenus {
		roleMenu := models.RoleMenu{
			RoleID: adminRole.ID,
			MenuID: menu.ID,
		}
		m.db.Create(&roleMenu)
	}

	// 创建默认管理员用户
	adminUser := models.UserModel{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Nickname: "系统管理员",
		Status:   1,
	}
	if err := m.db.Create(&adminUser).Error; err != nil {
		return err
	}

	// 为用户分配管理员角色
	userRole := models.UserRole{
		UserID: adminUser.ID,
		RoleID: adminRole.ID,
	}
	if err := m.db.Create(&userRole).Error; err != nil {
		return err
	}
	if err := m.db.Create(&userRole).Error; err != nil {
		return err
	}

	return nil
}

// CheckConnection 检查数据库连接
func (m *Migrator) CheckConnection() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetDatabaseInfo 获取数据库信息
func (m *Migrator) GetDatabaseInfo() map[string]interface{} {
	info := make(map[string]interface{})
	
	// 获取表信息
	var tables []string
	m.db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tables)
	info["tables"] = tables

	// 获取用户数量
	var userCount int64
	m.db.Model(&models.UserModel{}).Count(&userCount)
	info["user_count"] = userCount

	// 获取角色数量
	var roleCount int64
	m.db.Model(&models.RoleModel{}).Count(&roleCount)
	info["role_count"] = roleCount

	// 获取菜单数量
	var menuCount int64
	m.db.Model(&models.MenuModel{}).Count(&menuCount)
	info["menu_count"] = menuCount

	// 获取API数量
	var apiCount int64
	m.db.Model(&models.ApiModel{}).Count(&apiCount)
	info["api_count"] = apiCount

	return info
}

// BackupDatabase 备份数据库（仅支持MySQL）
func (m *Migrator) BackupDatabase(backupPath string) error {
	// 这里可以实现数据库备份逻辑
	// 实际项目中可以使用mysqldump或其他工具
	return nil
}

// OptimizeDatabase 优化数据库
func (m *Migrator) OptimizeDatabase() error {
	// 分析表并优化
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	// 这里可以添加数据库优化逻辑
	// 如：分析表、重建索引、清理碎片等
	_ = sqlDB // 避免变量未使用的编译错误
	return nil
}