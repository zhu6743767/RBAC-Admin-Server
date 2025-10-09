# RBAC Admin Server 部署指南

## 1. 项目简介

RBAC Admin Server 是一个基于角色的访问控制系统后端服务，提供用户管理、权限控制、角色分配等核心功能。

## 2. 技术栈

- Go 1.24.0
- Gin Web框架
- GORM 数据库ORM
- Redis 缓存
- JWT 身份认证
- Casbin 权限管理
- MySQL/SQLite/PostgreSQL 数据库支持

## 3. 环境准备

### 3.1 必要依赖

- Go 1.24.0 或更高版本
- Redis 服务器
- MySQL/PostgreSQL 数据库（可选，默认支持SQLite）
- Git

### 3.2 安装Go环境

请参考[Go官方文档](https://golang.org/doc/install)安装适合您操作系统的Go版本。

## 4. 项目获取

```bash
# 克隆项目
git clone https://github.com/your-username/rbac_admin_server.git
cd rbac_admin_server
```

## 5. 配置文件设置

### 5.1 主要配置文件

项目包含以下配置文件：

- `settings.yaml` - 主配置文件
- `.env` - 环境变量配置文件（包含敏感信息）
- `settings.yaml.example` - 配置文件示例
- `.env.example` - 环境变量配置示例

### 5.2 配置文件说明

#### 5.2.1 settings.yaml 配置

```yaml
# 系统基础配置
system:
  mode: "debug"          # 运行模式：debug/release
  ip: "127.0.0.1"        # 服务器IP
  port: 8080             # 服务器端口
  read_timeout: 30       # 读取超时时间（秒）
  write_timeout: 30      # 写入超时时间（秒）
  max_header_bytes: 1048576 # 最大请求头字节数

# 数据库配置
db:
  mode: "mysql"          # 数据库类型：mysql/sqlite/postgres
  host: "localhost"      # 数据库主机
  port: 3306             # 数据库端口
  user: "root"           # 数据库用户名
  password: "password"   # 数据库密码
  db_name: "rbac_admin"  # 数据库名
  path: "./rbac_admin.db" # SQLite数据库路径
  max_idle_conns: 10     # 最大空闲连接数
  max_open_conns: 100    # 最大打开连接数
  conn_max_lifetime: 3600 # 连接最大生命周期（秒）

# Redis配置
redis:
  addr: "localhost:6379" # Redis服务器地址
  password: ""           # Redis密码
  db: 0                  # Redis数据库索引
  pool_size: 10          # 连接池大小
  min_idle_conns: 5      # 最小空闲连接数
  dial_timeout: 3        # 拨号超时时间（秒）
  read_timeout: 3        # 读取超时时间（秒）
  write_timeout: 3       # 写入超时时间（秒）
  idle_timeout: 1800     # 空闲超时时间（秒）
  max_conn_age: 3600     # 连接最大年龄（秒）
  pool_timeout: 30       # 连接池超时时间（秒）
  idle_check_frequency: 60 # 空闲检查频率（秒）
  max_retries: 3         # 最大重试次数
  min_retry_backoff: 100 # 最小重试退避时间（毫秒）
  max_retry_backoff: 500 # 最大重试退避时间（毫秒）

# JWT配置
jwt:
  secret: "your-secret-key" # JWT密钥（请确保安全）
  expire_hours: 24        # 令牌有效期（小时）
  issuer: "rbac-admin-server" # 令牌签发者
  renew_window: 6         # 续期窗口（小时）

# 日志配置
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

# 安全配置
security:
  xss_protection: "1"    # XSS保护
  content_type_nosniff: "nosniff" # 内容类型嗅探
  x_frame_options: "DENY" # X-Frame-Options
  csrf_protection: true  # CSRF保护
  rate_limit: 100        # 速率限制（请求/秒）
  bcrypt_cost: 12        # BCrypt加密成本

# CORS配置
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

# 监控配置
monitoring:
  enabled: true          # 是否启用监控
  prometheus_port: 9090  # Prometheus端口
  health_check_path: "/health" # 健康检查路径
  metrics_path: "/metrics" # 指标路径
  trace_sampling_rate: 0.1 # 跟踪采样率

# Swagger配置
swagger:
  enabled: true          # 是否启用Swagger
  path: "/swagger"       # Swagger路径
  title: "RBAC Admin Server API" # API标题
  description: "RBAC Admin Server API Documentation" # API描述
  version: "1.0.0"       # API版本
  terms_of_service: ""   # 服务条款
  contact_name: "Admin"  # 联系人
  contact_url: ""        # 联系URL
  contact_email: "admin@example.com" # 联系邮箱
  license_name: "MIT"    # 许可证名称
  license_url: "https://opensource.org/licenses/MIT" # 许可证URL

# 应用配置
app:
  name: "RBAC Admin Server" # 应用名称
  version: "1.0.0"       # 应用版本
  description: "A RBAC Admin Server implemented in Go" # 应用描述
  copyright: "© 2023 RBAC Admin Server" # 版权信息
  timezone: "Asia/Shanghai" # 时区
  language: "zh-CN"      # 语言
  debug: true            # 调试模式

# 上传配置
upload:
  path: "./uploads"      # 上传文件保存路径
  max_size: 10           # 最大文件大小（MB）
  allowed_extensions:    # 允许的文件扩展名
    - ".jpg"
    - ".jpeg"
    - ".png"
    - ".gif"
    - ".pdf"
    - ".doc"
    - ".docx"
    - ".xls"
    - ".xlsx"
```

#### 5.2.2 .env 环境变量配置

```env
# 系统环境变量
SYSTEM_MODE=debug
SYSTEM_IP=127.0.0.1
SYSTEM_PORT=8080

# 数据库环境变量
DB_MODE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=rbac_admin

# Redis环境变量
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT环境变量
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
JWT_ISSUER=rbac-admin-server

# 日志环境变量
LOG_LEVEL=info
LOG_DIR=./logs
LOG_STDOUT=true

# 应用环境变量
APP_NAME=RBAC Admin Server
APP_VERSION=1.0.0
APP_DEBUG=true
```

### 5.3 配置文件使用说明

1. 复制配置文件示例：
   ```bash
   cp settings.yaml.example settings.yaml
   cp .env.example .env
   ```

2. 根据您的环境修改配置文件中的相应值。

3. 敏感信息（如数据库密码、JWT密钥）建议在`.env`文件中设置，而不是直接在`settings.yaml`中硬编码。

## 6. 编译和运行

### 6.1 直接运行

```bash
# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

### 6.2 编译后运行

```bash
# 编译项目
go build

# 运行编译后的二进制文件
./rbac_admin_server
```

### 6.3 指定配置文件

```bash
# 使用自定义配置文件
./rbac_admin_server -settings custom_settings.yaml
```

### 6.4 脚本运行

项目提供了便捷的运行脚本：

#### Windows系统
```powershell
# 运行简单测试
.	est_simple.bat

# 部署服务
.\deploy.bat
```

#### Linux/Mac系统
```bash
# 运行简单测试
chmod +x ./simple_test.ps1
./simple_test.ps1

# 部署服务
chmod +x ./deploy.sh
./deploy.sh
```

## 7. 验证服务

服务启动后，可以通过以下方式验证：

1. 访问健康检查接口：`http://localhost:8080/health`
2. 访问Swagger文档：`http://localhost:8080/swagger/index.html`（如果启用了Swagger）

## 8. API接口测试

项目提供了多个PowerShell脚本用于测试API接口：

- `test_login_detailed.ps1` - 详细测试登录功能
- `test_admin_user_list.ps1` - 测试获取用户列表
- `test_admin_crud_operations.ps1` - 测试管理员CRUD操作
- `test_all_admin_apis.ps1` - 测试所有管理员API

运行测试脚本示例：
```powershell
# 运行登录测试
.	est_login_detailed.ps1
```

## 9. 项目目录结构

```
rbac_admin_server/
├── api/               # API接口定义
├── config/            # 配置文件和配置加载逻辑
├── core/              # 核心初始化和功能实现
├── global/            # 全局变量和函数
├── middleware/        # 中间件
├── models/            # 数据模型
├── routes/            # 路由定义
├── utils/             # 工具函数
├── .env               # 环境变量配置
├── .env.example       # 环境变量配置示例
├── settings.yaml      # 主配置文件
├── settings.yaml.example # 配置文件示例
├── main.go            # 程序入口
└── DEPLOYMENT_GUIDE.md # 部署指南
```

## 10. 常见问题与解决方案

### 10.1 数据库连接失败

- 检查数据库服务是否正常运行
- 验证数据库配置是否正确（主机、端口、用户名、密码、数据库名）
- 确保数据库用户有足够的权限

### 10.2 Redis连接失败

- 检查Redis服务是否正常运行
- 验证Redis配置是否正确（地址、密码、数据库索引）
- 确保防火墙没有阻止Redis连接

### 10.3 JWT认证失败

- 检查JWT密钥是否配置正确
- 确保客户端请求头中包含正确的Authorization头
- 验证token是否过期

### 10.4 端口占用

- 如果8080端口已被占用，可以修改配置文件中的`system.port`值

## 11. 安全建议

1. 在生产环境中使用强JWT密钥
2. 不要在代码库中提交包含敏感信息的`.env`文件
3. 生产环境中设置`system.mode`为`release`
4. 生产环境中设置`log.level`为`info`或更高
5. 定期备份数据库
6. 配置适当的防火墙规则

## 12. 开发指南

### 12.1 添加新API

1. 在`api`目录下创建新的API包
2. 实现API逻辑
3. 在`routes/routes.go`中注册新路由

### 12.2 添加新模型

1. 在`models`目录下创建新的模型文件
2. 在`core/init_gorm/enter.go`中的`MigrateTables`函数中添加新模型

### 12.3 添加新中间件

1. 在`middleware`目录下创建新的中间件文件
2. 在`routes/routes.go`中的`SetupRouter`函数中应用新中间件

## 13. 部署到生产环境

### 13.1 Docker部署（推荐）

项目支持Docker部署，可以创建以下Dockerfile：

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rbac_admin_server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/rbac_admin_server .
COPY settings.yaml .
COPY .env .
COPY config/casbin/ ./config/casbin/

EXPOSE 8080
CMD ["./rbac_admin_server"]
```

然后构建和运行Docker容器：

```bash
docker build -t rbac_admin_server .
docker run -p 8080:8080 --env-file .env rbac_admin_server
```

### 13.2 传统部署

1. 在目标服务器上安装Go环境、Redis和数据库
2. 复制编译后的二进制文件和配置文件到目标服务器
3. 创建系统服务或使用进程管理工具（如Supervisor、Systemd）管理服务
4. 配置反向代理（如Nginx）以支持HTTPS

## 14. 维护与更新

### 14.1 更新流程

1. 备份数据库和配置文件
2. 拉取最新代码
3. 更新依赖：`go mod tidy`
4. 重新编译：`go build`
5. 重启服务

### 14.2 日志管理

- 日志文件默认保存在`./logs`目录
- 根据配置自动滚动和压缩
- 定期清理旧日志以节省磁盘空间

## 15. 联系与支持

如有问题或建议，请联系：admin@example.com