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
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.LoginLog{},
		&models.OperationLog{},
	}

	return m.db.AutoMigrate(modelList...)
}

// MigrateWithOptions 带选项的迁移
func (m *Migrator) MigrateWithOptions(dropIndex bool, modelList ...interface{}) error {
	if len(modelList) == 0 {
		modelList = []interface{}{
			&models.User{},
			&models.Role{},
			&models.Permission{},
			&models.UserRole{},
			&models.RolePermission{},
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
	m.db.Model(&models.User{}).Count(&count)
	if count > 0 {
		// 暂时移除日志，避免循环依赖
		return nil
	}

	// 创建默认权限
	permissions := []models.Permission{
		{Name: "user:create", Description: "创建用户", Resource: "user", Action: "create"},
		{Name: "user:read", Description: "查看用户", Resource: "user", Action: "read"},
		{Name: "user:update", Description: "更新用户", Resource: "user", Action: "update"},
		{Name: "user:delete", Description: "删除用户", Resource: "user", Action: "delete"},
		{Name: "role:create", Description: "创建角色", Resource: "role", Action: "create"},
		{Name: "role:read", Description: "查看角色", Resource: "role", Action: "read"},
		{Name: "role:update", Description: "更新角色", Resource: "role", Action: "update"},
		{Name: "role:delete", Description: "删除角色", Resource: "role", Action: "delete"},
		{Name: "permission:create", Description: "创建权限", Resource: "permission", Action: "create"},
		{Name: "permission:read", Description: "查看权限", Resource: "permission", Action: "read"},
		{Name: "permission:update", Description: "更新权限", Resource: "permission", Action: "update"},
		{Name: "permission:delete", Description: "删除权限", Resource: "permission", Action: "delete"},
	}

	if err := m.db.Create(&permissions).Error; err != nil {
		return err
	}

	// 创建管理员角色
	adminRole := models.Role{
		Name:        "admin",
		Description: "系统管理员",
		Status:      1,
	}
	if err := m.db.Create(&adminRole).Error; err != nil {
		return err
	}

	// 为管理员角色分配所有权限
	var allPermissions []models.Permission
	m.db.Find(&allPermissions)
	for _, permission := range allPermissions {
		rolePermission := models.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}
		m.db.Create(&rolePermission)
	}

	// 创建默认管理员用户
	adminUser := models.User{
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
	m.db.Model(&models.User{}).Count(&userCount)
	info["user_count"] = userCount

	// 获取角色数量
	var roleCount int64
	m.db.Model(&models.Role{}).Count(&roleCount)
	info["role_count"] = roleCount

	// 获取权限数量
	var permissionCount int64
	m.db.Model(&models.Permission{}).Count(&permissionCount)
	info["permission_count"] = permissionCount

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