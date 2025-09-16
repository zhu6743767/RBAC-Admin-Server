# 🏗️ RBAC管理员服务器配置指南

## 📋 配置概述

本项目支持三种环境配置：
- **开发环境** (`settings_dev.yaml`) - 调试模式，SQLite数据库
- **测试环境** (`settings_test.yaml`) - 测试模式，独立数据库
- **生产环境** (`settings_prod.yaml`) - 生产模式，环境变量配置

## 🚀 快速开始

### 1. 开发环境启动
```bash
# 使用开发配置
go run main.go -env dev

# 或指定配置文件
go run main.go -config settings_dev.yaml
```

### 2. 测试环境启动
```bash
# 使用测试配置
go run main.go -env test

# 或指定配置文件
go run main.go -config settings_test.yaml
```

### 3. 生产环境启动
```bash
# 使用生产配置（需要环境变量）
go run main.go -env prod

# 或指定配置文件
go run main.go -config settings_prod.yaml
```

## ⚙️ 配置优先级

1. **命令行参数** `-config` 指定的配置文件
2. **命令行参数** `-env` 选择的环境配置
3. **环境变量** 覆盖配置文件中的值

## 🔧 环境变量配置

复制 `.env.example` 为 `.env` 并修改：

```bash
cp .env.example .env
# 编辑 .env 文件，设置你的配置值
```

### 关键环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `SERVER_PORT` | 服务器端口 | `8080` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PASSWORD` | 数据库密码 | `your_secure_password` |
| `JWT_SECRET` | JWT密钥 | `your_jwt_secret` |
| `REDIS_PASSWORD` | Redis密码 | `your_redis_password` |

## 🧪 测试配置

运行配置测试：

```bash
# 测试开发环境配置
go run main_simple.go -env dev

# 测试测试环境配置  
go run main_simple.go -env test

# 测试生产环境配置
go run main_simple.go -env prod
```

测试结果示例：

### 开发环境
```
🚀 RBAC管理员服务器配置测试
✅ 配置加载成功!
📋 应用信息: RBAC管理员 - 开发环境 v1.0.0-dev (development)
🖥️  服务器端口: 8080
🗄️  数据库类型: sqlite
🔐 JWT颁发者: rbac-admin-dev
📝 日志级别: debug
```

### 测试环境
```
🚀 RBAC管理员服务器配置测试
✅ 配置加载成功!
📋 应用信息: RBAC管理员 - 测试环境 v1.0.0-test (testing)
🖥️  服务器端口: 8081
🗄️  数据库类型: mysql
🔐 JWT颁发者: rbac-admin-test
📝 日志级别: info
```

### 生产环境
```
🚀 RBAC管理员服务器配置测试
✅ 配置加载成功!
📋 应用信息: RBAC管理员 v1.0.0 (production)
🖥️  服务器端口: 8080
🗄️  数据库类型: mysql
🔐 JWT颁发者: rbac-admin
📝 日志级别: info
```

## 📁 配置文件结构

```
rbac_admin_server/
├── settings_dev.yaml      # 开发环境配置
├── settings_test.yaml     # 测试环境配置
├── settings_prod.yaml     # 生产环境配置
├── .env.example          # 环境变量示例
├── config/
│   ├── config.go         # 配置结构体定义
│   └── loader.go         # 配置加载器
└── test_config.go        # 配置测试脚本
```

## 🔍 配置验证

配置加载器会自动验证：
- ✅ 配置文件格式正确性
- ✅ 必填字段完整性
- ✅ 环境变量解析
- ✅ 默认值应用
- ✅ 路径存在性检查

## 🚨 注意事项

### 生产环境
- 必须使用强JWT密钥
- 数据库密码必须复杂
- 启用SSL连接
- 关闭调试模式
- 配置合理的日志级别

### 开发环境
- 使用SQLite便于开发
- 启用详细日志
- 允许CORS跨域
- 启用Swagger文档

### 测试环境
- 使用独立测试数据库
- 模拟生产环境配置
- 启用所有安全功能

## 📞 故障排除

### 常见问题

1. **配置文件找不到**
   ```bash
   # 检查文件路径
   ls -la settings_*.yaml
   ```

2. **环境变量不生效**
   ```bash
   # 检查环境变量
   env | grep -E "(SERVER|DB|JWT|REDIS)"
   ```

3. **配置验证失败**
   ```bash
   # 运行测试脚本查看详细错误
   go run test_config.go
   ```

## 🎯 下一步

1. 根据你的环境修改对应的配置文件
2. 设置必要的环境变量
3. 运行测试确保配置正确
4. 启动服务并验证功能