# RBAC管理员服务器

一个基于Go的RBAC（基于角色的访问控制）管理系统，支持多环境配置和灵活的配置管理。

## 🚀 快速开始

### 1. 环境准备

确保已安装：
- Go 1.24.0 或更高版本
- Git

### 2. 获取项目

```bash
git clone <your-repo-url>
cd rbac_admin_server
```

### 3. 配置环境

#### 开发环境
```bash
# 使用开发环境配置
go run main_final.go -env dev
```

#### 测试环境
```bash
# 使用测试环境配置
go run main_final.go -env test
```

#### 生产环境
```bash
# 1. 复制环境变量模板
cp .env.example .env

# 2. 编辑 .env 文件，设置实际的环境变量值

# 3. 使用生产环境配置
go run main_final.go -env prod
```

### 4. 自定义配置

可以通过以下方式指定配置文件：
```bash
go run main_final.go -config /path/to/your/config.yaml
```

## 📁 项目结构

```
rbac_admin_server/
├── config/                 # 配置管理
│   ├── config.go          # 配置结构体定义
│   └── loader.go          # 配置加载器
├── settings_dev.yaml      # 开发环境配置
├── settings_test.yaml     # 测试环境配置
├── settings_prod.yaml     # 生产环境配置
├── .env.example           # 环境变量模板
├── .env                   # 生产环境变量（需要创建）
├── CONFIG_GUIDE.md        # 配置指南
├── main_simple.go         # 配置测试脚本
└── main_final.go          # 最终主程序
```

## ⚙️ 配置系统特性

### 🎯 多环境支持
- **开发环境** (`dev`): 使用SQLite，调试模式开启
- **测试环境** (`test`): 使用MySQL，调试模式关闭
- **生产环境** (`prod`): 使用MySQL，调试模式关闭，环境变量配置

### 🔧 配置优先级
1. 命令行参数 (`-config`)
2. 环境变量 (`.env`文件或系统环境变量)
3. YAML配置文件
4. 默认值

### 📊 配置验证
- 自动验证配置完整性
- 友好的错误提示
- 支持环境变量替换

### 🔒 安全特性
- 敏感信息通过环境变量管理
- JWT密钥配置
- CORS跨域配置
- 安全头配置

## 🧪 测试配置

项目提供了配置测试功能：

```bash
# 测试所有环境配置
go run main_simple.go -env dev   # 开发环境
go run main_simple.go -env test  # 测试环境
go run main_simple.go -env prod  # 生产环境
```

## 📋 环境变量

生产环境需要设置以下环境变量（参考 `.env.example`）：

```bash
# 服务器配置
SERVER_PORT=8080

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=rbac_prod

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT配置
JWT_SECRET=your-super-secret-jwt-key
JWT_ISSUER=rbac-admin-prod
JWT_AUDIENCE=rbac-admin-users

# 应用配置
APP_NAME=RBAC管理员
APP_VERSION=1.0.0
APP_ENVIRONMENT=production

# 安全与监控
SECURITY_ENABLED=true
CORS_ALLOWED_ORIGINS=https://yourdomain.com
METRICS_ENABLED=true
```

## 📝 开发指南

### 添加新配置
1. 在 `config/config.go` 中添加新的配置结构体
2. 在 `config/loader.go` 中添加默认值和环境变量处理
3. 在所有环境的YAML配置文件中添加相应配置
4. 更新 `.env.example` 文件
5. 更新 `CONFIG_GUIDE.md` 文档

### 配置调试
使用配置测试脚本快速验证配置：
```bash
go run main_simple.go -env dev
```

## 🤝 贡献

欢迎提交Issue和Pull Request来改进配置系统！

## 📄 许可证

MIT License