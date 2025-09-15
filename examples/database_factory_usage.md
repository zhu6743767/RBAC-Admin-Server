# 数据库工厂模式使用指南

## 概述

本项目采用工厂模式实现多数据库支持，能够根据配置文件自动创建和管理不同类型的数据库连接。

## 支持的数据库类型

- **MySQL**: 生产环境首选
- **PostgreSQL**: 企业级应用
- **SQLite**: 开发测试环境
- **SQL Server**: Windows环境

## 配置文件示例

### MySQL配置
```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  name: "rbac_admin"
  ssl_mode: "disable"
  timeout: 30
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
```

### PostgreSQL配置
```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "postgres"
  name: "rbac_admin"
  ssl_mode: "disable"
  timeout: 30
```

### SQLite配置（推荐开发使用）
```yaml
database:
  type: "sqlite"
  path: "./rbac_admin.db"
```

## 环境变量配置

### Windows环境变量
```bash
# MySQL
set DB_TYPE=mysql
set DB_HOST=localhost
set DB_PORT=3306
set DB_USERNAME=root
set DB_PASSWORD=password
set DB_NAME=rbac_admin

# PostgreSQL
set DB_TYPE=postgres
set DB_HOST=localhost
set DB_PORT=5432
set DB_USERNAME=postgres
set DB_PASSWORD=postgres
set DB_NAME=rbac_admin

# SQLite
set DB_TYPE=sqlite
set DB_PATH=./rbac_admin.db
```

### Linux/Mac环境变量
```bash
# MySQL
export DB_TYPE=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=password
export DB_NAME=rbac_admin

# PostgreSQL
export DB_TYPE=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_NAME=rbac_admin

# SQLite
export DB_TYPE=sqlite
export DB_PATH=./rbac_admin.db
```

## 代码使用示例

### 1. 基本使用

```go
package main

import (
    "log"
    "rbac.admin/config"
    "rbac.admin/database"
    "rbac.admin/global"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("settings.yaml")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // 初始化数据库
    dbManager, err := database.NewDatabaseManager(cfg.Database)
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    
    // 设置全局数据库管理器
    global.DBManager = dbManager
    
    // 执行自动迁移
    migrator := database.NewMigrator(dbManager.GetDB())
    if err := migrator.AutoMigrate(); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
    
    log.Println("Database initialized successfully")
}
```

### 2. 事务处理示例

```go
func createUserWithRole(db *gorm.DB, userData *models.User, roleID uint) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        if err := tx.Create(userData).Error; err != nil {
            return err
        }
        
        // 关联角色
        userRole := &models.UserRole{
            UserID: userData.ID,
            RoleID: roleID,
        }
        
        return tx.Create(userRole).Error
    })
}
```

### 3. 复杂查询示例

```go
// 获取用户及其角色权限
func getUserWithPermissions(db *gorm.DB, userID uint) (*models.User, error) {
    var user models.User
    
    err := db.Preload("Roles.Permissions").
        First(&user, userID).Error
    
    return &user, err
}

// 分页查询用户
func getUsersPaginated(db *gorm.DB, page, pageSize int) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    offset := (page - 1) * pageSize
    
    db.Model(&models.User{}).Count(&total)
    
    err := db.Limit(pageSize).
        Offset(offset).
        Order("created_at desc").
        Find(&users).Error
    
    return users, total, err
}
```

## 数据库迁移

### 自动迁移
项目启动时会自动执行数据库迁移，创建所需的表结构：

- **users**: 用户表
- **roles**: 角色表
- **permissions**: 权限表
- **user_roles**: 用户角色关联表
- **role_permissions**: 角色权限关联表

### 手动迁移
```go
// 创建迁移器
migrator := database.NewMigrator(dbManager.GetDB())

// 执行迁移
if err := migrator.AutoMigrate(); err != nil {
    log.Fatal("Migration failed:", err)
}

// 初始化基础数据
if err := migrator.SeedData(); err != nil {
    log.Fatal("Seeding failed:", err)
}
```

## 数据初始化

### 默认管理员账户
首次启动时会自动创建以下基础数据：

1. **超级管理员角色**
   - 名称：超级管理员
   - 描述：系统最高权限角色

2. **管理员账户**
   - 用户名：admin
   - 密码：admin123
   - 邮箱：admin@example.com
   - 角色：超级管理员

3. **基础权限**
   - 用户管理权限
   - 角色管理权限
   - 权限管理权限
   - 系统管理权限

## 连接池配置

数据库连接池参数自动配置：

```yaml
# MySQL连接池参数
max_open_conns: 100    # 最大连接数
max_idle_conns: 10     # 空闲连接数
conn_max_lifetime: 3600 # 连接最大生命周期(秒)

# PostgreSQL连接池参数
max_open_conns: 50
max_idle_conns: 5
conn_max_lifetime: 1800

# SQLite参数
max_open_conns: 1      # SQLite单连接
max_idle_conns: 1
conn_max_lifetime: 0   # 永久连接
```

## 故障排查

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务是否启动
   - 验证连接参数是否正确
   - 检查防火墙设置

2. **权限不足**
   - 确保数据库用户有足够权限
   - 检查数据库是否存在
   - 验证用户名和密码

3. **SQLite文件权限**
   - 确保程序有写入权限
   - 检查磁盘空间是否充足

### 调试模式
```go
// 启用GORM调试日志
dbManager.GetDB().Debug()
```

## 性能优化

### 索引优化
自动创建的索引：
- users.username (唯一索引)
- users.email (唯一索引)
- roles.name (唯一索引)
- permissions.name (唯一索引)

### 查询优化
- 使用预加载减少N+1查询
- 合理使用分页避免大数据查询
- 使用事务保证数据一致性

## 扩展支持

### 添加新的数据库类型
在`database_factory.go`中添加新的数据库类型：

```go
// 1. 添加数据库类型常量
const (
    DatabaseTypeOracle DatabaseType = "oracle"
)

// 2. 实现Oracle配置结构体
type OracleConfig struct {
    BaseConfig
    // Oracle特有配置
}

// 3. 实现创建方法
func (f *DatabaseFactory) createOracleConnection(config DatabaseConfig) (*gorm.DB, error) {
    // Oracle连接实现
}
```

这套数据库工厂模式提供了完整的多数据库支持，可以根据项目需求灵活切换数据库类型，同时保持代码的一致性和可维护性。