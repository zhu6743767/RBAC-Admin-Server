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
# 使用开发环境配置（使用SQLite内存数据库）
go run main.go -env dev
```

#### 测试环境
```bash
# 使用测试环境配置（需要MySQL数据库）
go run main.go -env test
```

#### 生产环境
```bash
# 1. 复制环境变量模板
cp .env.example .env

# 2. 编辑 .env 文件，设置实际的环境变量值

# 3. 使用生产环境配置（需要MySQL数据库和Redis）
go run main.go -env prod
```

### 4. 自定义配置

可以通过以下方式指定配置文件：
```bash
go run main.go -config /path/to/your/config.yaml
```

## 📁 项目结构

```
rbac_admin_server/
├── api/                    # API接口
├── config/                 # 配置管理
│   ├── config.go          # 配置结构体定义
│   └── loader.go          # 配置加载器
├── core/                   # 核心模块
│   ├── db.go              # 数据库初始化
│   ├── init.go              # 系统初始化
│   ├── logger.go          # 日志初始化
│   ├── redis.go              # Redis初始化
│   └── validator.go          # 验证器初始化
├── middleware/              # 中间件
├── models/                  # 数据模型
├── routes/                  # 路由
├── settings_dev.yaml      # 开发环境配置
├── settings_test.yaml     # 测试环境配置
├── settings_prod.yaml     # 生产环境配置
├── .env.example           # 环境变量模板
├── .env                   # 生产环境变量（需要创建）
├── CONFIG_GUIDE.md        # 配置指南
├── main.go                # 主程序入口
└── README.md              # 项目文档
```

## ⚙️ 配置系统特性

### 🎯 多环境支持
- **开发环境** (`dev`): 使用SQLite内存数据库，调试模式开启，无需外部依赖
- **测试环境** (`test`): 使用MySQL，调试模式关闭，模拟生产环境
- **生产环境** (`prod`): 使用MySQL+Redis，调试模式关闭，环境变量配置

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
- 数据库连接加密
- API访问限流

## 🧪 测试配置

项目提供了配置测试功能：

```bash
# 测试开发环境配置（无需外部依赖）
go run main.go -env dev

# 测试测试环境配置（需要MySQL数据库）
go run main.go -env test

# 测试生产环境配置（需要MySQL数据库和Redis）
go run main.go -env prod
```

## 🔐 安全配置

**⚠️ 重要安全提醒：**

详细的安全配置指南请参考 [SECURITY.md](SECURITY.md)。

### 🛡️ 敏感信息保护
- **生产环境**：所有敏感信息必须通过环境变量配置
- **禁止提交**：切勿将 `.env` 文件或包含真实密码的配置文件提交到版本库
- **强密码策略**：使用强密码和长JWT密钥（32位以上）

### 📋 环境变量配置

生产环境需要设置以下环境变量（参考 `.env.example`）：

```bash
# 🖥️ 系统配置
SYSTEM_PORT=8080

# 🗄️ 数据库配置（使用实际的数据库信息）
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=rbac_user
DB_PASSWORD=your_secure_password_here  # ⚠️ 生产环境必须修改
DB_NAME=rbac_admin

# 🔐 JWT配置（使用强密钥）
JWT_SECRET=your_jwt_secret_key_minimum_32_characters  # ⚠️ 必须修改
JWT_EXPIRE_HOURS=24
JWT_ISSUER=rbac-admin
JWT_AUDIENCE=rbac-admin

# 🔄 Redis配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password_here  # 如有密码

# 📝 日志配置
LOG_LEVEL=info
LOG_DIR=./logs

# 🎯 应用配置
APP_NAME=RBAC管理员
APP_VERSION=1.0.0
APP_ENVIRONMENT=production
APP_DEBUG=false

# 🔒 安全配置
CSRF_SECRET=your_csrf_secret_key_here  # ⚠️ 必须修改

# 🌐 CORS配置
CORS_ORIGINS=https://your-domain.com
```

## 🚀 开发指南

### 开发环境搭建

1. **克隆项目**
   ```bash
   git clone <your-repo-url>
   cd rbac_admin_server
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置环境**
   - 开发环境：直接使用默认配置，无需额外设置
   - 测试环境：需要MySQL数据库
   - 生产环境：需要MySQL数据库和Redis，复制 `.env.example` 为 `.env` 并配置

4. **运行项目**
   ```bash
   # 开发环境（推荐，零配置启动）
   go run main.go -env dev
   
   # 测试环境（需要MySQL）
   go run main.go -env test
   
   # 生产环境（需要MySQL+Redis）
   go run main.go -env prod
   ```

### 添加新配置
1. 在 `config/config.go` 中添加新的配置结构体
2. 在 `config/loader.go` 中添加默认值和环境变量处理
3. 在所有环境的YAML配置文件中添加相应配置
4. 更新 `.env.example` 文件
5. 更新 `CONFIG_GUIDE.md` 文档

### 配置验证
使用配置测试脚本快速验证配置：
```bash
go run main.go -env dev
```

## 🤝 贡献

欢迎提交Issue和Pull Request来改进配置系统！

## 📄 许可证

MIT License