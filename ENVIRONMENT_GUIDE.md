# 🌍 环境切换使用指南

## 📋 概述

本项目支持通过命令行参数切换不同运行环境，每个环境使用独立的配置文件，便于开发、测试和生产部署。

## 🎯 支持的环境

| 环境 | 参数 | 配置文件 | 用途 |
|------|------|----------|------|
| 开发环境 | `dev` 或 `development` | `settings_dev.yaml` | 本地开发，调试模式 |
| 测试环境 | `test` 或 `testing` | `settings_test.yaml` | 功能测试，模拟生产 |
| 生产环境 | `prod` 或 `production` | `settings_prod.yaml` | 线上部署，安全优化 |

## 🚀 使用方法

### 1. 直接运行（默认开发环境）
```bash
go run main.go
# 或使用编译后的二进制文件
./rbac-admin-server
```

### 2. 指定运行环境
```bash
# 开发环境（默认）
go run main.go -env=dev

# 测试环境
go run main.go -env=test

# 生产环境
go run main.go -env=prod

# 简写形式
go run main.go -env=development
go run main.go -env=testing
go run main.go -env=production
```

### 3. 指定自定义配置文件
```bash
# 完全自定义配置文件（优先级最高）
go run main.go -config=/path/to/your/config.yaml

# 同时使用环境和自定义配置
go run main.go -env=prod -config=/etc/rbac-admin/settings.yaml
```

### 4. 查看帮助信息
```bash
go run main.go -h
```

## 📁 配置文件说明

### 🔧 开发环境 (settings_dev.yaml)
- **数据库**: SQLite（无需安装MySQL）
- **日志级别**: debug（最详细）
- **调试模式**: 启用
- **CORS**: 允许所有来源（便于前后端分离开发）
- **Swagger**: 启用
- **性能分析**: 启用（pprof）

### 🧪 测试环境 (settings_test.yaml)
- **数据库**: MySQL测试数据库
- **日志级别**: info（适中）
- **调试模式**: 关闭
- **CORS**: 限制指定来源
- **安全配置**: 启用所有安全功能
- **性能**: 模拟生产环境限制

### 🏭 生产环境 (settings_prod.yaml)
- **数据库**: MySQL生产数据库
- **日志级别**: info（生产级）
- **调试模式**: 关闭
- **安全配置**: 最高级别
- **SSL**: 强制启用
- **性能**: 优化配置
- **敏感信息**: 全部使用环境变量

## 🔐 环境变量配置

### 生产环境必需的环境变量
```bash
# 数据库配置
export DB_HOST=your-db-host
export DB_USERNAME=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=rbac_admin

# JWT配置
export JWT_SECRET=your-256-bit-secret-key

# Redis配置（可选）
export REDIS_HOST=localhost
export REDIS_PASSWORD=your-redis-password

# CSRF保护
export CSRF_SECRET=your-csrf-secret

# 日志目录
export LOG_DIR=/var/log/rbac-admin
```

## 🛠️ 快速开始

### 开发环境快速启动
```bash
# 1. 克隆项目
git clone <repository-url>
cd rbac-admin-server

# 2. 安装依赖
go mod tidy

# 3. 启动开发环境
go run main.go -env=dev

# 4. 访问应用
# 打开浏览器访问: http://localhost:8080
# API文档: http://localhost:8080/swagger/index.html
```

### 生产环境部署
```bash
# 1. 编译应用
go build -o rbac-admin-server main.go

# 2. 设置环境变量
export DB_HOST=your-prod-db
export DB_PASSWORD=your-prod-password
export JWT_SECRET=your-prod-jwt-secret

# 3. 启动生产环境
./rbac-admin-server -env=prod
```

## 📊 环境验证

启动后，控制台会显示当前环境信息：

```
🚀 RBAC管理员服务器启动中...
╔═══════════════════════════════════════╗
║          RBAC Admin Server            ║
║    Role-Based Access Control System   ║
╚═══════════════════════════════════════╝

🌍 运行环境: DEV
📁 配置文件: settings_dev.yaml
🗄️ 数据库: SQLite(./data/rbac_admin_dev.db)
──────────────────────────────────────────────────
✅ RBAC管理员服务器启动成功!
══════════════════════════════════════════════════
🌐 访问地址: http://localhost:8080
📊 健康检查: http://localhost:8080/health
📈 监控指标: http://localhost:8080/metrics
📚 API文档: http://localhost:8080/swagger/index.html
🗄️ 数据库: MySQL(root@localhost:3306/rbac_admin)
📊 日志级别: debug
══════════════════════════════════════════════════
```

## 🔄 环境切换脚本

创建便捷的启动脚本：

### run-dev.sh（Linux/Mac）
```bash
#!/bin/bash
echo "🚀 启动开发环境..."
go run main.go -env=dev
```

### run-dev.bat（Windows）
```batch
@echo off
echo 🚀 启动开发环境...
go run main.go -env=dev
pause
```

## 📝 最佳实践

### 1. 开发阶段
- 使用 `-env=dev` 启动
- 利用SQLite快速迭代
- 开启详细日志和调试信息

### 2. 测试阶段
- 使用 `-env=test` 启动
- 连接独立测试数据库
- 验证所有功能正常

### 3. 生产部署
- 使用 `-env=prod` 启动
- 确保所有环境变量已设置
- 使用负载均衡和健康检查

### 4. 配置管理
- 开发环境：直接修改配置文件
- 测试环境：配置文件 + 部分环境变量
- 生产环境：全部使用环境变量

## 🚨 注意事项

1. **敏感信息**：生产环境所有敏感信息必须使用环境变量
2. **数据库**：不同环境使用不同的数据库实例
3. **日志**：生产环境日志目录需要有写入权限
4. **端口**：确保各环境端口不冲突
5. **备份**：生产环境配置定期备份

## 🔍 故障排查

### 配置文件未找到
```bash
# 检查文件是否存在
ls -la settings_*.yaml

# 或指定完整路径
go run main.go -config=/full/path/to/settings.yaml
```

### 端口冲突
```bash
# 检查端口占用
netstat -an | grep 8080

# 修改配置文件中的端口
# 或在配置中使用环境变量 PORT
```

### 环境变量未生效
```bash
# 检查环境变量
env | grep DB_

# 在Linux/Mac上临时设置
export DB_PASSWORD=yourpassword

# 在Windows上设置
set DB_PASSWORD=yourpassword
```