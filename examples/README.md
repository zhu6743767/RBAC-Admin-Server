# 示例文件使用指南

本目录包含了RBAC管理员服务器的各种使用示例。

## 📁 文件说明

### 示例文件
- `database_usage_example.go` - 数据库使用示例，展示如何连接数据库、执行迁移、查询数据
- `factory_usage_example.go` - 日志工厂模式使用示例
- `logging_usage_example.go` - 日志系统使用示例（注释形式）
- `main_example.go` - 完整的使用示例

### 运行文件
- `run_examples.go` - 主运行程序，提供交互式菜单运行各个示例

## 🚀 如何运行示例

### 方法1：使用运行程序（推荐）
```bash
# 编译运行程序
go build -o run_examples.exe run_examples.go

# 运行程序
./run_examples.exe

# 或者直接运行
go run run_examples.go
```

### 方法2：单独运行示例
```bash
# 运行特定示例
go run examples/database_usage_example.go

# 注意：由于包名限制，某些示例可能需要手动修改
```

### 方法3：编译所有示例
```bash
# 编译所有示例文件
go build -v ./examples/...
```

## 📊 示例内容

### 数据库使用示例
- 连接MySQL数据库
- 执行数据库迁移
- 初始化种子数据
- 查询用户信息
- 创建新用户
- 事务处理示例
- 关联查询示例
- 分页查询示例
- 复杂查询示例

### 日志工厂示例
- 使用logrus日志
- 动态切换日志配置
- 创建不同配置的日志实例

### 完整示例
- 完整的项目初始化流程
- 数据库连接和迁移
- 数据查询展示

## ⚙️ 配置说明

所有示例默认使用以下数据库配置：
- **类型**: MySQL
- **主机**: 192.168.10.199
- **端口**: 3306
- **用户名**: root
- **密码**: Zdj_7819!
- **数据库**: rbac_admin

## 🔧 自定义配置

要修改数据库配置，请编辑相应示例文件中的配置部分：

```go
cfg := config.DefaultConfig()
cfg.Database.Type = "mysql"
cfg.Database.Host = "your-host"
cfg.Database.Port = 3306
cfg.Database.Username = "your-user"
cfg.Database.Password = "your-password"
cfg.Database.Database = "your-database"
```

## 📝 注意事项

1. **包名限制**：示例文件使用`package examples`，不能直接运行
2. **依赖关系**：确保项目依赖已正确安装（`go mod tidy`）
3. **数据库连接**：确保数据库服务已启动且配置正确
4. **权限问题**：确保数据库用户有足够的权限

## 🐛 故障排查

### 编译错误
- 检查Go版本（需要1.21+）
- 运行`go mod tidy`更新依赖
- 确保所有文件在同一模块下

### 运行错误
- 检查数据库连接配置
- 确认数据库服务运行正常
- 检查防火墙和网络连接
- 验证数据库用户权限

### 数据库连接失败
```bash
# 测试MySQL连接
mysql -h 192.168.10.199 -u root -p rbac_admin
```