# RBAC Admin Server 全面部署指南

## 🚀 快速开始

如果你想快速部署和运行RBAC Admin Server，可以按照以下步骤操作：

### 使用编译后运行（推荐）

```bash
# 克隆项目
git clone https://github.com/rbacadmin/rbac_admin_server.git
cd rbac_admin_server

# 创建配置文件
cp settings.yaml.example settings.yaml
cp .env.example .env

# 根据实际情况修改配置文件
# 修改settings.yaml中的服务器、数据库、Redis等配置
# 修改.env文件中的敏感信息

# 安装依赖
go mod tidy

# 编译项目
go build

# 执行数据库迁移
go run main.go -m db -t migrate

# 创建管理员用户
go run main.go -m user -t create -username admin -password admin123

# 启动服务
./rbac_admin_server
```

## 📋 目录

1. [项目简介](#1-项目简介)
2. [技术栈](#2-技术栈)
3. [环境准备](#3-环境准备)
4. [项目获取](#4-项目获取)
5. [配置文件设置](#5-配置文件设置)
6. [编译和运行](#6-编译和运行)
7. [数据库初始化](#7-数据库初始化)
8. [服务验证](#8-服务验证)
9. [API接口测试](#9-api接口测试)
10. [常见问题与解决方案](#10-常见问题与解决方案)
11. [安全建议](#11-安全建议)

## 1. 项目简介

RBAC Admin Server 是一个基于 Go 语言开发的 RBAC（基于角色的访问控制）系统后端服务，提供用户管理、角色管理、权限管理以及文件管理等功能。它采用前后端分离架构，为企业级应用提供完整的用户权限管理解决方案。

## 2. 技术栈

- Go 1.24.0
- Gin Web框架
- GORM 数据库ORM
- Redis 缓存
- JWT 身份认证
- Casbin 权限管理
- MySQL/SQLite/PostgreSQL 数据库支持
- Swagger API文档
- 验证码(Captcha)安全验证

## 3. 环境准备

在部署 RBAC Admin Server 之前，需要准备以下环境：

- **Go 1.24.0 或更高版本**：用于编译和运行项目
- **Redis 7.0 或更高版本**：用于缓存和会话管理
- **MySQL 8.0 或更高版本**：用于数据存储
- **Git**：用于获取项目代码

### 3.1 Go 环境安装

#### Windows 系统

1. 访问 [Go 官方下载页面](https://golang.org/dl/)，下载 Go 1.24.0 安装包
2. 运行安装包并按照提示完成安装
3. 打开命令提示符，运行 `go version` 验证安装是否成功

```cmd
# 验证安装
> go version
go version go1.24.0 windows/amd64

# 配置 Go 代理（可选，加速依赖下载）
> go env -w GOPROXY=https://goproxy.cn,direct
> go env -w GO111MODULE=on
```

### 3.2 Redis 安装

#### Windows 系统

1. 访问 [Redis 官方下载页面](https://redis.io/download/)，下载最新的 Redis 安装包
2. 解压安装包并按照说明完成安装
3. 打开命令提示符，运行 `redis-server` 启动 Redis 服务

```cmd
# 启动 Redis 服务
> redis-server

# 验证安装（另开一个命令提示符）
> redis-cli ping
PONG
```

### 3.3 数据库安装

#### MySQL 8.0 安装

**Windows 系统**

1. 访问 [MySQL 官方下载页面](https://dev.mysql.com/downloads/installer/)，下载 MySQL 8.0 安装包
2. 运行安装包并按照提示完成安装
3. 设置 root 密码和其他必要配置

### 3.4 数据库初始化

安装完成后，需要创建一个数据库供 RBAC Admin Server 使用。请按照以下步骤执行数据库初始化：

1. **登录 MySQL 数据库**：
   ```bash
   # 打开命令行工具，输入以下命令登录MySQL
   mysql -u root -p
   # 然后输入MySQL root用户的密码
   ```

2. **创建数据库并设置字符集**：
   ```sql
   -- 创建数据库（如果不存在）
   CREATE DATABASE IF NOT EXISTS rbacadmin 
   CHARACTER SET utf8mb4 
   COLLATE utf8mb4_unicode_ci;
   
   -- 查看数据库是否创建成功
   SHOW DATABASES;
   ```

3. **创建数据库用户并授权**：
   ```sql
   -- 创建一个专门的数据库用户（推荐生产环境使用）
   CREATE USER IF NOT EXISTS 'rbacadmin'@'%' IDENTIFIED BY 'YourSecurePassword';
   
   -- 授予该用户对rbacadmin数据库的所有权限
   GRANT ALL PRIVILEGES ON rbacadmin.* TO 'rbacadmin'@'%';
   
   -- 刷新权限
   FLUSH PRIVILEGES;
   
   -- 查看用户权限
   SHOW GRANTS FOR 'rbacadmin'@'%';
   ```

4. **配置文件中的数据库连接信息**：
   请确保在`.env`文件中正确配置数据库连接信息：
   ```env
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=rbacadmin  # 或root，取决于你使用哪个用户
   DB_PASSWORD=YourSecurePassword
   DB_DBNAME=rbacadmin
   ```

**注意事项**：
- 在生产环境中，强烈建议使用专门的数据库用户（如'rbacadmin'）而不是root用户
- 请使用强密码并妥善保管
- utf8mb4字符集支持完整的Unicode字符集，包括emoji表情符号
- '%'表示允许从任何IP地址连接，生产环境中可根据需要限制连接IP

## 4. 项目获取

### 4.1 从 GitHub 克隆

```bash
git clone https://github.com/rbacadmin/rbac_admin_server.git
cd rbac_admin_server
```

### 4.2 下载源码包

如果没有安装 Git，也可以直接下载源码包：

1. 访问 [项目 GitHub 页面](https://github.com/rbacadmin/rbac_admin_server)
2. 点击 "Code" 按钮，然后选择 "Download ZIP"
3. 解压下载的 ZIP 文件
4. 进入解压后的目录

## 5. 配置文件设置

### 5.1 主要配置文件

项目包含以下配置文件：

- `settings.yaml` - 主配置文件（包含系统核心配置）
- `.env` - 环境变量配置文件（包含敏感信息和环境特定配置）
- `settings.yaml.example` - 配置文件示例
- `.env.example` - 环境变量配置示例

### 5.2 配置文件说明

#### 5.2.1 settings.yaml 配置

settings.yaml 是项目的主要配置文件，包含系统运行所需的各种参数。以下是完整的配置文件结构说明：

#### 5.2.1.1 系统基础配置

```yaml
# 🖥️ 服务器配置
system:
  ip: "0.0.0.0"         # 服务器绑定IP，生产环境建议设置为具体IP
  port: 8080            # 服务器监听端口
  mode: "debug"         # 运行模式: debug, release（生产环境建议使用release）
  name: "RBAC Admin Server"  # 应用名称
  version: "1.0.0"      # 应用版本号
  timezone: "Asia/Shanghai"  # 时区设置
```

#### 5.2.1.2 数据库配置

```yaml
# 🗄️ 数据库配置
db:
  mode: "mysql"            # 数据库类型: mysql, postgres, sqlite
  host: "127.0.0.1"       # 数据库主机
  port: 3306              # 数据库端口
  user: "root"             # 数据库用户名
  password: "${DB_PASSWORD}"     # 数据库密码（通过环境变量注入，更安全）
  dbname: "rbacadmin"     # 数据库名称
  sslmode: "disable"      # SSL模式: disable, require, verify-ca, verify-full
  timeout: "30s"          # 数据库连接超时时间
  charset: "utf8mb4"      # 数据库字符集
  collation: "utf8mb4_unicode_ci"  # 数据库排序规则
  max_open_conns: 100     # 最大连接数
  max_idle_conns: 10      # 空闲连接数
  conn_max_lifetime: 3600 # 连接生命周期(秒)
  conn_max_idle_time: 1800 # 连接最大空闲时间
  path: "rbac_admin.db"   # SQLite数据库路径（仅SQLite模式有效）
```

#### 5.2.1.3 Redis配置

```yaml
# 🔄 Redis配置
redis:
  addr: "${REDIS_ADDR:127.0.0.1:6379}"  # Redis地址（支持环境变量替换）
  password: "${REDIS_PASSWORD}"         # Redis密码
  db: ${REDIS_DB:3}                      # Redis数据库编号
  pool_size: 20                          # 连接池大小
  min_idle_conns: 5                      # 最小空闲连接数
  max_conn_age: 3600                     # 连接最大存活时间(秒)
  pool_timeout: 30                       # 连接池获取连接超时时间(秒)
  idle_timeout: 1800                     # 空闲连接超时时间(秒)
  read_timeout: 3                        # 读取超时时间(秒)
  write_timeout: 3                       # 写入超时时间(秒)
  dial_timeout: 3                        # 拨号超时时间(秒)
```

#### 5.2.1.4 JWT认证配置

```yaml
# 🔐 JWT认证配置
jwt:
  secret: "${JWT_SECRET:aB3kL9mN7xY2qR8sT1uV4wE6zC0pF5gH}"  # JWT密钥（至少32个字符）
  expire_hours: 72                            # Token过期时间(小时)
  refresh_expire_hours: 168                   # 刷新Token过期时间(小时)
  issuer: "rbacAdmin"                         # Token签发者
  audience: "rbac-client"                     # Token受众
  signing_method: "HS256"                     # 签名方法
  token_name: "Authorization"                 # 令牌在请求头中的名称
```

#### 5.2.1.5 验证码配置

```yaml
# 🧩 验证码配置（重要）
# 验证码功能默认启用，登录接口必须使用验证码
captcha:
  enable: true          # 是否启用验证码：true(启用)/false(禁用)
  width: 120            # 验证码图片宽度
  height: 40            # 验证码图片高度
  length: 4             # 验证码长度
  expire_seconds: 300   # 验证码有效期(秒)
  noise_level: 1        # 噪点级别
  fonts: []             # 自定义字体列表

# 🔧 如何修改验证码配置
# 在settings.yaml文件中找到captcha部分，根据需要修改以下参数：
# - enable: 设置为false可以禁用验证码功能
# - expire_seconds: 修改验证码有效期（单位：秒）
# - length: 修改验证码字符长度
# 示例：将验证码有效期延长至10分钟（600秒）
# captcha:
#   enable: true
#   expire_seconds: 600
#   length: 6
```

#### 5.2.1.6 日志配置

```yaml
# 📝 日志配置
log:
  level: "info"          # 日志级别：debug/info/warn/error/fatal/panic
  dir: "./logs"          # 日志目录
  filename: "rbac_admin.log" # 日志文件名
  format: "text"         # 日志格式：text/json
  max_size: 100          # 单文件最大大小（MB）
  max_age: 7             # 最大保留天数
  max_backups: 3         # 最大备份数量
  compress: true         # 是否压缩
  stdout: true           # 是否输出到标准输出
  enable_caller: true    # 是否显示调用者信息
  output: "both"         # 日志输出方式: stdout, file, both
```

#### 5.2.1.7 安全配置

```yaml
# 🔒 安全配置
security:
  xss_protection: "1"    # XSS保护
  content_type_nosniff: "nosniff" # 内容类型嗅探
  x_frame_options: "DENY" # X-Frame-Options
  csrf_protection: true  # CSRF保护
  rate_limit: 100        # 速率限制（请求/秒）
  bcrypt_cost: 12        # BCrypt加密成本
  brute_force_protection: true  # 是否启用暴力破解保护
  password_complexity: 3  # 密码复杂度要求（1-5）
  login_attempts_limit: 5  # 登录尝试次数限制
  login_lockout_time: 30  # 登录锁定时间（分钟）
  session_timeout: 24h   # 会话超时时间
```

#### 5.2.1.8 CORS配置

```yaml
# 🌐 CORS配置
cors:
  allow_origins:         # 允许的源
    - "http://localhost:3000"
    - "http://localhost:8080"
  allow_methods:         # 允许的HTTP方法
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:         # 允许的HTTP头
    - "Origin"
    - "Content-Type"
    - "Authorization"
  allow_credentials: true # 是否允许凭证
  max_age: 3600          # 预检请求缓存时间（秒）
  expose_headers: []     # 暴露的HTTP头列表
```

#### 5.2.1.9 性能优化配置

```yaml
# ⚡ 性能优化配置
performance:
  max_request_size: 10   # 最大请求大小（MB）
  max_upload_size: "100M" # 最大上传文件大小
  request_timeout: 30s   # 请求处理超时时间
  response_compression: true # 是否启用响应压缩
  gzip_level: 6          # GZIP压缩级别（1-9）
  cache_control: "no-cache" # Cache-Control头设置
  etag: true             # 是否启用ETag
```

#### 5.2.1.10 文件上传配置

```yaml
# 📤 文件上传配置
upload:
  dir: "uploads"         # 上传文件保存目录
  max_size: 100          # 最大文件大小（MB）
  allowed_types: ["image/jpeg", "image/png", "application/pdf"] # 允许的文件类型
  file_permissions: 0644 # 文件权限
  dir_permissions: 0755  # 目录权限
  max_files_per_request: 10 # 每次请求最多上传文件数
```

#### 5.2.1.11 监控配置

```yaml
# 📊 监控配置
monitoring:
  enabled: true          # 是否启用监控
  endpoint: "/metrics"   # 监控指标端点
  prometheus_enabled: true # 是否启用Prometheus
```

#### 5.2.1.12 API文档配置

```yaml
# 📚 API文档配置
swagger:
  enabled: true          # 是否启用Swagger文档
  path: "/swagger"       # Swagger文档访问路径
  title: "RBAC Admin Server API" # API文档标题
  description: "RBAC权限管理系统API文档" # API文档描述
  version: "1.0.0"       # API版本
```

#### 5.2.1.13 应用特定配置

```yaml
# 🎯 应用配置
app:
  debug: false           # 是否启用调试模式（生产环境必须为false）
  graceful_shutdown_timeout: 30s # 优雅关闭超时时间
  request_id: true       # 是否为每个请求生成唯一ID
  trace_id: true         # 是否为每个请求生成跟踪ID
```

#### 5.2.1.14 环境变量使用说明

配置文件中使用 `${ENV_VAR:default_value}` 格式可以引用环境变量，系统会自动从 .env 文件或系统环境变量中读取对应值。推荐将敏感信息（如密码、密钥等）通过环境变量注入。

示例用法：
- `${REDIS_ADDR:127.0.0.1:6379}` 表示：优先使用环境变量 REDIS_ADDR 的值，如果不存在则使用默认值 127.0.0.1:6379
- `.env` 文件中的配置会自动加载为环境变量

**注意：** 所有敏感信息（如密码、密钥等）必须通过环境变量配置，不要在settings.yaml中直接硬编码敏感信息。

#### 5.2.2 .env 环境变量配置

.env文件用于存储敏感信息和环境特定配置，以下是完整的环境变量说明，按功能分组：

#### 5.2.2.1 核心敏感配置 - 生产环境必须修改

```env
DB_PASSWORD=your_secure_database_password  # 数据库密码（必填）
REDIS_PASSWORD=your_secure_redis_password  # Redis密码（如无密码留空）
JWT_SECRET=your_secure_jwt_secret_minimum_32_characters  # JWT密钥（必填，建议32位以上）
CSRF_SECRET=your_secure_csrf_secret  # CSRF密钥（必填）
```

#### 5.2.2.2 数据库配置

```env
DB_MODE=mysql                    # 数据库类型: mysql, postgres, sqlite
DB_HOST=localhost                # 数据库主机地址
DB_PORT=3306                     # 数据库端口
DB_USER=root                     # 数据库用户名
DB_DBNAME=rbacadmin              # 数据库名称
DB_PATH=./data/rbac_admin.db     # SQLite数据库路径（当DB_MODE=sqlite时使用）
DB_SSLMODE=disable               # SSL模式: disable, require, verify-ca, verify-full
DB_TIMEOUT=30s                   # 数据库连接超时时间
```

#### 5.2.2.3 Redis配置

```env
REDIS_ADDR=localhost:6379        # Redis服务器地址和端口
REDIS_DB=0                       # Redis数据库索引
```

#### 5.2.2.4 JWT认证配置

```env
JWT_EXPIRE_HOURS=24              # JWT令牌有效期（小时）
JWT_REFRESH_EXPIRE_HOURS=72      # 刷新令牌有效期（小时）
JWT_ISSUER=RBAC Admin Server     # 令牌颁发者
JWT_AUDIENCE=rbac-admin-users    # 令牌受众
JWT_SIGNING_METHOD=HS256         # 签名方法
```

#### 5.2.2.5 系统配置

```env
SYSTEM_IP=0.0.0.0                # 服务器绑定IP
SYSTEM_PORT=8080                 # 服务器监听端口
SYSTEM_NAME=RBAC Admin Server    # 应用名称
SYSTEM_VERSION=1.0.0             # 应用版本号
SYSTEM_TIMEZONE=Asia/Shanghai    # 时区设置
```

#### 5.2.2.6 日志配置

```env
LOG_LEVEL=info                   # 日志级别: debug, info, warn, error, fatal, panic
LOG_DIR=./logs                   # 日志目录
LOG_FILENAME=rbac_admin.log      # 日志文件名
LOG_STDOUT=true                  # 是否输出到标准输出
LOG_FORMAT=text                  # 日志格式: text, json
```

#### 5.2.2.7 安全配置

```env
SECURITY_RATE_LIMIT=100          # 请求频率限制（每分钟）
SECURITY_LOGIN_ATTEMPTS_LIMIT=5  # 登录尝试次数限制
SECURITY_LOCK_DURATION_MINUTES=30  # 登录锁定时间（分钟）
SECURITY_BCRYPT_COST=12          # bcrypt加密成本
```

#### 5.2.2.8 CORS配置

```env
CORS_ALLOW_ORIGINS=https://your-domain.com  # 允许的源，生产环境建议设置为具体域名
```

#### 5.2.2.9 应用配置

```env
APP_DEBUG=false                  # 是否启用调试模式（生产环境必须为false）
APP_ENVIRONMENT=production       # 应用环境: development, testing, production
APP_GRACEFUL_SHUTDOWN_TIMEOUT=30s  # 优雅关闭超时时间
```

#### 5.2.2.10 监控配置

```env
MONITORING_ENABLED=true          # 是否启用监控
MONITORING_ENDPOINT=/metrics     # 监控指标端点
```

#### 5.2.2.11 API文档配置

```env
SWAGGER_ENABLED=true             # 是否启用Swagger文档
SWAGGER_PATH=/swagger            # Swagger文档访问路径
```

#### 5.2.2.12 性能优化配置

```env
PERFORMANCE_MAX_REQUEST_SIZE=10  # 最大请求大小（MB）
PERFORMANCE_REQUEST_TIMEOUT=30s  # 请求处理超时时间
PERFORMANCE_RESPONSE_COMPRESSION=true  # 是否启用响应压缩
```

#### 5.2.2.13 文件上传配置

```env
UPLOAD_DIR=./uploads             # 上传文件保存目录
UPLOAD_MAX_SIZE=100              # 最大文件大小（MB）
```

#### 5.2.2.14 生产环境配置示例

```env
# =================================================================================
# 🏗️ RBAC管理员服务器 - 生产环境配置
# =================================================================================

# 🖥️ 系统配置
SYSTEM_PORT=8080

# 🗄️ 数据库配置
DB_MODE=mysql
DB_HOST=192.168.10.199
DB_PORT=3306
DB_USER=root
DB_PASSWORD=Zdj_7819!
DB_DBNAME=rbacadmin

# 🔐 JWT配置
JWT_SECRET=AbCdEfGhIjKlMnOpQrStUvWxYz1234567890
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_HOURS=168

# 🔄 Redis配置
REDIS_ADDR=192.168.10.199:6379
REDIS_PASSWORD=
REDIS_DB=4

# 📝 日志配置
LOG_LEVEL=info
LOG_DIR=./logs

# 🎯 应用配置
APP_NAME=RBAC管理员
APP_VERSION=1.0.0
APP_ENVIRONMENT=production
APP_DEBUG=false
```

#### 5.2.2.15 安全提示

- 不要将包含敏感信息的.env文件提交到代码仓库
- 敏感信息（密码、密钥等）应通过环境变量注入
- 确保文件权限设置正确，防止未授权访问
- 生产环境中请使用强密码和密钥，并妥善保管此文件

#### 5.2.2.16 提示

- 如需覆盖更多配置，请参考settings.yaml.example中的配置项
- 环境变量命名格式：配置项路径中的点(.)替换为下划线(_)，全部大写
- 例如：`system.port` -> `SYSTEM_PORT`
- 环境变量会覆盖settings.yaml中的同名配置

### 5.3 配置文件使用说明

#### 5.3.1 创建配置文件

1. 复制配置文件示例：
   ```bash
   # Linux/Mac系统
   cp settings.yaml.example settings.yaml
   cp .env.example .env
   
   # Windows系统(PowerShell)
   Copy-Item -Path settings.yaml.example -Destination settings.yaml
   Copy-Item -Path .env.example -Destination .env
   ```

2. 根据您的环境修改配置文件中的相应值。

3. 敏感信息（如数据库密码、JWT密钥）**必须**在`.env`文件中设置，而不是直接在`settings.yaml`中硬编码。

#### 5.3.2 环境变量替换机制

settings.yaml 支持通过 `${ENV_VAR:default_value}` 语法从环境变量中读取配置，优先级高于直接在YAML中定义的值：

- `${REDIS_ADDR:127.0.0.1:6379}` 表示：优先使用环境变量 REDIS_ADDR 的值，如果不存在则使用默认值 127.0.0.1:6379
- `.env` 文件中的配置会自动加载为环境变量

#### 5.3.3 配置优先级

配置项的优先级从高到低依次为：
1. 命令行参数（如 `-settings custom_settings.yaml`）
2. `.env` 文件中定义的环境变量
3. settings.yaml 中使用 `${}` 语法定义的默认值
4. settings.yaml 中直接定义的值

## 6. 编译和运行

### 6.1 配置系统说明

RBAC Admin Server 采用灵活的配置系统，支持多环境配置和环境变量替换。以下是配置系统的关键特性：

#### 6.1.1 多环境支持

系统默认支持三种环境配置：

- **开发环境** (`settings_dev.yaml`) - 调试模式，SQLite数据库，便于开发和调试
- **测试环境** (`settings_test.yaml`) - 测试模式，独立数据库，用于自动化测试
- **生产环境** (`settings_prod.yaml`) - 生产模式，环境变量配置，更高的安全性

#### 6.1.2 配置优先级

配置项的优先级从高到低依次为：
1. **命令行参数** `-config` 指定的配置文件
2. **命令行参数** `-env` 选择的环境配置
3. **环境变量** 覆盖配置文件中的值
4. `.env` 文件中定义的环境变量
5. settings.yaml 中使用 `${}` 语法定义的默认值
6. settings.yaml 中直接定义的值

#### 6.1.3 环境变量替换

settings.yaml 支持通过 `${ENV_VAR:default_value}` 语法从环境变量中读取配置：

- `${REDIS_ADDR:127.0.0.1:6379}` 表示：优先使用环境变量 REDIS_ADDR 的值，如果不存在则使用默认值 127.0.0.1:6379
- `.env` 文件中的配置会自动加载为环境变量

#### 6.1.4 配置验证

系统启动时会自动验证配置的有效性：
- ✅ 配置文件格式正确性
- ✅ 必填字段完整性
- ✅ 环境变量解析
- ✅ 默认值应用
- ✅ 路径存在性检查

如果配置验证失败，系统会输出详细的错误信息并退出。

### 6.2 安装依赖

```bash
go mod tidy
```

### 6.3 直接运行

```bash
# 使用默认配置运行（默认使用dev环境）
go run main.go

# 指定环境运行
go run main.go -env prod

# 指定配置文件运行
go run main.go -config custom_settings.yaml

# 查看帮助信息
go run main.go -h
```

### 6.4 编译后运行

```bash
# 编译项目
go build

# 运行编译后的二进制文件
./rbac_admin_server

# 指定环境运行
./rbac_admin_server -env prod

# 指定配置文件运行
./rbac_admin_server -config custom_settings.yaml
```

### 6.5 命令行参数说明

系统支持以下命令行参数：

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| `-env` | 指定运行环境 | `dev` | `./rbac_admin_server -env prod` |
| `-config` | 指定配置文件路径 | `settings.yaml` | `./rbac_admin_server -config custom.yaml` |
| `-m` | 指定运行模式 | `server` | `./rbac_admin_server -m db -t migrate` |
| `-t` | 指定运行子模式（与-m配合使用） | - | `./rbac_admin_server -m user -t create` |
| `-h` | 显示帮助信息 | - | `./rbac_admin_server -h` |

### 6.6 部署脚本说明

项目根目录提供了多个部署相关的脚本，可以简化部署过程：

#### 6.6.1 Windows 部署脚本

1. **deploy.bat** - Windows系统一键部署脚本
   ```bash
   # 在Windows命令提示符中运行
   deploy.bat
   ```
   该脚本会自动执行以下操作：
   - 安装依赖
   - 编译项目
   - 执行数据库迁移
   - 启动服务

2. **create_admin_user.bat** - Windows系统创建管理员用户脚本
   ```bash
   # 在Windows命令提示符中运行
   create_admin_user.bat
   ```
   该脚本会引导你创建一个新的管理员用户。

3. **start_and_monitor.bat** - Windows系统启动并监控服务脚本
   ```bash
   # 在Windows命令提示符中运行
   start_and_monitor.bat
   ```
   该脚本会启动服务并监控服务运行状态。

#### 6.6.2 Linux/Mac 部署脚本

1. **deploy.sh** - Linux/Mac系统一键部署脚本
   ```bash
   # 在Linux/Mac终端中运行
   chmod +x deploy.sh
   ./deploy.sh
   ```
   该脚本会自动执行以下操作：
   - 安装依赖
   - 编译项目
   - 执行数据库迁移
   - 启动服务

**注意事项**：
- 运行脚本前，请确保已正确配置`.env`文件
- 在Linux/Mac系统上，可能需要先赋予脚本执行权限
- 脚本执行过程中，可能需要根据提示输入必要的信息

### 6.7 Makefile使用说明

项目提供了完整的Makefile，用于简化构建和运行流程。以下是常用的Makefile命令：

```bash
# 下载依赖
make deps

# 更新依赖
make update-deps

# 构建完整版服务器
make build-full

# 运行完整版服务器
make run-full

# 构建简化版服务器
make build-simple

# 运行简化版服务器
make run-simple

# 构建无CGO版本（用于交叉编译）
make build-nocgo

# 运行测试
make test

# 清理构建文件
make clean

# 初始化数据目录（创建data和logs目录）
make init

# 显示帮助信息
make help
```

使用示例：

```bash
# 构建并运行完整版服务器
make build-full && make run-full

# 清理并重新构建
make clean && make build-full

# 初始化环境并运行测试
make init && make deps && make test```

## 7. 数据库初始化

### 7.1 数据库迁移

在首次运行或更新项目后，必须执行数据库迁移以创建或更新表结构。这是确保系统正常运行的重要步骤：

1. **执行数据库迁移**：
   ```bash
   # 在项目根目录下执行以下命令
   go run main.go -m db -t migrate
   
   # 预期输出（成功时）
   # 2023/10/13 15:30:25 数据库迁移成功
   # 2023/10/13 15:30:25 创建默认角色和权限
   ```

2. **检查迁移结果**：
   迁移成功后，可以登录MySQL数据库查看创建的表：
   ```sql
   USE rbacadmin;
   SHOW TABLES;
   ```
   应该能看到系统所需的所有表，如users、roles、permissions等。

### 7.2 创建管理员用户

数据库迁移完成后，必须创建一个管理员用户才能登录系统：

1. **创建管理员用户**：
   ```bash
   # 格式：go run main.go -m user -t create -username [用户名] -password [密码]
   go run main.go -m user -t create -username admin -password admin123
   
   # 预期输出（成功时）
   # 2023/10/13 15:35:10 管理员用户admin创建成功
   # 2023/10/13 15:35:10 已为用户分配超级管理员角色
   ```

2. **重置管理员密码**：
   如果管理员用户已存在或需要更改密码，可以使用以下命令：
   ```bash
   # 格式：go run main.go -m user -t reset -username [用户名] -password [新密码]
   go run main.go -m user -t reset -username admin -password new_secure_password
   
   # 预期输出（成功时）
   # 2023/10/13 15:40:05 用户admin密码重置成功
   ```

**重要提示**：
- 首次启动系统前，必须完成数据库迁移和管理员用户创建这两个步骤
- 生产环境中，请使用强密码并妥善保管
- 系统默认情况下，只有管理员用户具有所有权限

## 8. 服务验证

### 8.1 默认端口说明

服务默认运行在`8080`端口上。在服务成功启动后，你会在控制台看到类似以下的输出：
```
后端服务运行在 http://127.0.0.1:8080
```

### 8.2 如何更改端口号

要更改服务运行的端口号，请修改`settings.yaml`文件中的`system.port`配置项：
```yaml
# 🖥️ 服务器配置
system:
  ip: 127.0.0.1           # 服务IP地址
  port: 8081              # 修改为你想要的端口号
  mode: "debug"           # 运行模式: debug, release
```

### 8.3 验证方法

服务启动后，可以通过以下方式验证：

1. 访问健康检查接口：`http://localhost:8080/health`（如果更改了端口号，请使用新的端口）
2. 访问Swagger文档：`http://localhost:8080/swagger/index.html`（如果启用了Swagger）

## 9. API接口测试

### 9.1 测试说明

系统API接口需要进行身份认证，大部分接口需要在请求头中提供有效的JWT令牌。以下是测试API接口的基本流程：

1. **获取验证码**：调用`/public/captcha/get`接口获取验证码ID和验证码内容
2. **登录**：使用获取的验证码和用户名密码调用`/public/login`接口获取JWT令牌
3. **调用受保护接口**：在请求头中添加`Authorization: Bearer {token}`来调用需要认证的接口

### 9.2 验证码功能说明（重要）

系统默认启用了验证码功能，在调用登录接口前必须先获取验证码。验证码有效期为300秒（5分钟）。

### 9.3 登录流程

1. 获取验证码：
   ```bash
   curl -X GET http://localhost:8080/public/captcha/get
   ```

2. 使用获取的验证码和用户名密码登录：
   ```bash
   curl -X POST http://localhost:8080/public/login \
   -H "Content-Type: application/json" \
   -d '{"username":"admin","password":"admin123","captcha":"1234","captchaId":"captcha_id_from_step_1"}'
   ```

3. 使用返回的token访问受保护接口：
   ```bash
   curl -X GET http://localhost:8080/admin/user/list \
   -H "Authorization: Bearer your_jwt_token_here"
   ```

### 9.4 项目自带测试工具

项目根目录提供了API测试工具，可以方便地进行接口测试：

1. **api_client_test.go** - API客户端测试文件
   ```bash
   # 运行API客户端测试
   go run api_client_test.go
   ```

2. **api_test目录** - 包含详细的API测试用例
   ```bash
   # 运行所有API测试
   go test ./api_test/...
   
   # 运行特定的API测试
   go test ./api_test/test_login_flow.go -v
   ```

3. **test_login.json** - 登录测试配置文件
   可以修改此文件中的配置来进行登录测试。

### 9.5 测试结果示例

API测试成功后的输出示例：

```
=== RUN   TestLoginFlow
=== PAUSE TestLoginFlow
=== CONT  TestLoginFlow
✅  获取验证码成功: captcha_id=6d7a8b9c
✅  登录成功，token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
✅  访问用户列表成功，获取到10个用户
✅  测试流程完成
--- PASS: TestLoginFlow (2.53s)
PASS
ok      command-line-arguments  2.542s
```

## 10. 常见问题与解决方案

### 10.1 编译错误

**问题**: 编译时出现重复声明错误

**解决方案**: 检查是否有多个包含main函数的文件在根目录。将测试文件移动到专门的测试目录中，例如：
```bash
mkdir -p api_test
mv test_*.go api_test/
```

### 10.2 数据库连接错误

**问题**: 无法连接到数据库

**解决方案**: 检查以下几点：
1. MySQL服务是否正常运行
2. `.env`文件中的数据库配置是否正确
3. 数据库用户是否有足够的权限
4. 防火墙是否阻止了连接

### 10.3 Redis连接错误

**问题**: 无法连接到Redis

**解决方案**: 检查以下几点：
1. Redis服务是否正常运行
2. `.env`文件中的Redis配置是否正确
3. 防火墙是否阻止了连接

### 10.4 端口被占用

**问题**: 启动服务器时提示端口被占用

**解决方案**: 修改settings.yaml中的system.port配置，使用一个未被占用的端口。

## 11. 安全建议

### 11.1 敏感信息保护策略

#### 11.1.1 环境变量优先原则
- **生产环境**：所有敏感信息必须通过环境变量配置
- **开发环境**：可以使用配置文件，但建议使用环境变量
- **禁止**：将真实密码、密钥提交到版本库

#### 11.1.2 敏感信息分类

##### 🔑 数据库凭据
- `DB_PASSWORD` - 数据库密码
- `DB_HOST` - 数据库主机地址
- `DB_USER` - 数据库用户名

##### 🔐 JWT安全密钥
- `JWT_SECRET` - JWT签名密钥（建议32位以上）
- `JWT_ISSUER` - JWT发行者
- `JWT_AUDIENCE` - JWT受众

##### 🔄 Redis配置
- `REDIS_PASSWORD` - Redis密码
- `REDIS_ADDR` - Redis地址

##### 🔒 其他安全配置
- `CSRF_SECRET` - CSRF保护密钥
- 各类API密钥和令牌

### 11.2 配置方法

#### 11.2.1 开发环境配置
1. 复制 `.env.example` 为 `.env`
2. 填写真实的配置值
3. 确保 `.env` 文件在 `.gitignore` 中

#### 11.2.2 生产环境配置
1. **不要创建 `.env` 文件**
2. 通过操作系统的环境变量设置
3. 使用密钥管理服务（如Docker Secrets、Kubernetes Secrets等）

#### 11.2.3 Docker环境示例
```bash
# 使用环境变量运行
docker run -e DB_PASSWORD=your_password \
           -e JWT_SECRET=your_jwt_secret \
           -p 8080:8080 \
           rbac-admin-server
```

#### 11.2.4 Linux系统环境变量
```bash
# 临时设置
export DB_PASSWORD="your_password"
export JWT_SECRET="your_jwt_secret_key_minimum_32_characters"

# 永久设置（添加到 /etc/environment 或 ~/.bashrc）
echo 'export DB_PASSWORD="your_password"' >> ~/.bashrc
echo 'export JWT_SECRET="your_jwt_secret"' >> ~/.bashrc
source ~/.bashrc
```

#### 11.2.5 Windows系统环境变量
```powershell
# 临时设置
$env:DB_PASSWORD = "your_password"
$env:JWT_SECRET = "your_jwt_secret_key_minimum_32_characters"

# 永久设置（系统环境变量）
[Environment]::SetEnvironmentVariable("DB_PASSWORD", "your_password", "Machine")
[Environment]::SetEnvironmentVariable("JWT_SECRET", "your_jwt_secret", "Machine")
```

### 11.3 安全警告

#### 11.3.1 ❌ 禁止行为
- 不要将真实密码、密钥提交到Git仓库
- 不要在代码中硬编码敏感信息
- 不要将 `.env` 文件上传到生产服务器
- 不要使用弱密码或短密钥

#### 11.3.2 ✅ 推荐做法
- 使用强密码（12位以上，包含大小写字母、数字、特殊字符）
- JWT密钥至少32位，使用随机生成的字符串
- 定期更换密钥和密码
- 使用密钥管理服务
- 启用SSL/TLS加密传输
- 限制数据库和Redis的网络访问

### 11.4 配置验证

#### 11.4.1 开发环境验证
```bash
# 检查配置是否加载成功
go run validate_config.go -env dev

# 启动服务测试
go run main.go -env dev
```

#### 11.4.2 生产环境验证
```bash
# 验证生产环境配置（干运行模式）
go run validate_config.go -env prod

# 检查环境变量是否设置正确
echo $DB_PASSWORD  # Linux/Mac
$env:DB_PASSWORD   # PowerShell
```

### 11.5 部署检查清单

部署到生产环境前，请确认：

- [ ] 所有敏感信息通过环境变量配置
- [ ] `.env` 文件不存在于生产环境
- [ ] 数据库密码强度足够
- [ ] JWT密钥长度至少32位且为随机字符串
- [ ] Redis已设置密码保护（如需要）
- [ ] 数据库和Redis只允许可信IP访问
- [ ] 启用SSL/TLS证书
- [ ] 定期备份和密钥轮换策略

### 11.6 问题反馈

如发现安全配置问题，请立即：
1. 停止服务
2. 更换所有受影响的密钥和密码
3. 检查访问日志
4. 联系安全团队

---

**记住：安全是第一优先级！永远不要在代码中暴露敏感信息。**

---

**最后更新时间**: 2025-10-13
**版本**: 1.0.0