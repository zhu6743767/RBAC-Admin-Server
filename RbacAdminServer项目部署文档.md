<<<<<<< HEAD
## 🎯 项目概述

RBAC管理员服务器是一个基于Go语言开发的权限管理系统，采用RBAC（基于角色的访问控制）模型，支持多环境配置、灵活的数据库支持和完善的安全机制。本系统提供了用户管理、角色管理、权限管理、部门管理、菜单管理等核心功能，适用于企业级应用的权限控制需求。

## 📋 环境要求

### 基础环境
- **Go版本**: 1.24.0 或更高版本
- **Git**: 用于版本控制
- **操作系统**: Windows/Linux/macOS

### 数据库支持
- **开发环境**: SQLite（文件或内存数据库，零配置）
- **测试环境**: MySQL 8.0+ 或 PostgreSQL 14+ 或 SQLite
- **生产环境**: MySQL 8.0+ 或 PostgreSQL 14+，推荐配合Redis使用

### 生产环境额外要求
- **Redis**: 6.0+（用于会话管理和缓存）
- **反向代理**: Nginx/Apache（推荐配置SSL证书）

## 🚀 快速部署

### 1. 获取项目代码

```bash
# 克隆项目代码
git clone <your-repo-url>
cd rbac_admin_server

# 安装依赖
go mod download
go mod tidy
```

### 2. 配置文件准备

项目使用`settings.yaml`作为主配置文件，支持通过环境变量覆盖敏感配置：

```bash
# 创建配置文件（如果不存在）
if not exist settings.yaml (copy settings.yaml.example settings.yaml)
```

### 3. 数据库初始化

```bash
# 创建数据库表结构
go run main.go -m db -t migrate

# 创建管理员用户
go run main.go -m user -t create
# 按照提示输入用户名和密码
```

### 4. 环境部署

#### 🔧 开发环境部署（推荐）

开发环境默认使用SQLite数据库，无需额外配置，适合快速开发和测试：

```bash
# 直接运行，零配置启动
go run main.go -env dev
```

默认配置：
- 端口: 8080
- 数据库: SQLite文件数据库 (rbac_admin.db)
- 日志级别: debug
- 调试模式: 开启

#### 🧪 测试环境部署

测试环境推荐使用MySQL或PostgreSQL数据库：

```bash
# 1. 创建测试数据库
# MySQL
sql -u root -p -e "CREATE DATABASE rbac_admin_test CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

# PostgreSQL
createdb -U postgres -h localhost -p 5432 rbac_admin_test

# 2. 修改测试环境配置文件 settings_test.yaml
# 可以从主配置文件复制并修改
copy settings.yaml settings_test.yaml
# 编辑 settings_test.yaml 文件，配置数据库连接

# 3. 运行测试环境
go run main.go -env test
```

#### 🏭 生产环境部署

生产环境需要完整的配置和安全设置：

```bash
# 1. 创建生产环境配置文件
copy settings.yaml settings_prod.yaml
# 编辑 settings_prod.yaml 文件，配置生产环境参数

# 2. 配置环境变量（推荐使用.env文件）
echo "DB_PASSWORD=your_db_password" > .env
echo "REDIS_PASSWORD=your_redis_password" >> .env
echo "JWT_SECRET=your_jwt_secret_key_minimum_32_characters" >> .env

# 3. 创建生产数据库
# MySQL
sql -u root -p -e "CREATE DATABASE rbac_admin_prod CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

echo "CREATE USER 'rbac_user'@'localhost' IDENTIFIED BY 'your_db_password';" | mysql -u root -p
echo "GRANT ALL PRIVILEGES ON rbac_admin_prod.* TO 'rbac_user'@'localhost';" | mysql -u root -p
echo "FLUSH PRIVILEGES;" | mysql -u root -p

# 4. 构建应用
go build -o rbac_admin_server main.go

# 5. 运行应用（使用生产环境配置）
./rbac_admin_server -env prod
```

### 5. 生产环境推荐部署方式

#### 使用Systemd（Linux系统）

创建systemd服务文件：

```bash
# 创建服务文件
vi /etc/systemd/system/rbac_admin_server.service
```

服务文件内容：

```ini
[Unit]
Description=RBAC Admin Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/path/to/rbac_admin_server
ExecStart=/path/to/rbac_admin_server/rbac_admin_server -env prod
EnvironmentFile=/path/to/rbac_admin_server/.env
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
systemctl daemon-reload
systemctl start rbac_admin_server
systemctl enable rbac_admin_server

# 查看服务状态
systemctl status rbac_admin_server

# 查看日志
journalctl -u rbac_admin_server -f
```

#### 使用Docker（容器化部署）

创建Dockerfile：

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o rbac_admin_server main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/rbac_admin_server .
COPY settings_prod.yaml .
EXPOSE 8080
CMD ["./rbac_admin_server", "-env", "prod"]
```

创建docker-compose.yml：

```yaml
version: '3'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
      - redis
    environment:
      - DB_HOST=db
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
    restart: unless-stopped

  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=${DB_ROOT_PASSWORD}
      - MYSQL_DATABASE=rbac_admin_prod
      - MYSQL_USER=rbac_user
      - MYSQL_PASSWORD=${DB_PASSWORD}
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:6-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    restart: unless-stopped

volumes:
  mysql-data:
  redis-data:
```

启动Docker容器：

```bash
docker-compose up -d
```

## ⚙️ 详细配置说明

### 配置文件结构

项目使用YAML配置文件，支持多环境配置（dev/test/prod）：

#### 完整配置示例 (settings.yaml)

```yaml
# 服务器配置
system:
  ip: "0.0.0.0"           # 绑定IP
  port: 8080              # 监听端口
  name: "RBAC管理员系统"   # 系统名称
  version: "1.0.0"        # 系统版本
  timezone: "Asia/Shanghai" # 时区

# 数据库配置
db:
  mode: "mysql"           # 数据库类型: mysql/postgres/sqlite
  host: "localhost"       # 数据库主机
  port: 3306              # 数据库端口
  user: "root"            # 数据库用户名
  password: "${DB_PASSWORD}" # 数据库密码（从环境变量读取）
  dbname: "rbac_admin"    # 数据库名称
  sslmode: "disable"      # SSL模式
  timeout: "30s"          # 连接超时
  charset: "utf8mb4"      # 字符集
  collation: "utf8mb4_general_ci" # 排序规则
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数
  conn_max_lifetime: "1h"  # 连接最大生命周期
  conn_max_idle_time: "30m" # 连接最大空闲时间

# Redis配置
redis:
  addr: "localhost:6379"  # Redis地址
  password: "${REDIS_PASSWORD}" # Redis密码（从环境变量读取）
  db: 0                   # Redis数据库索引
  pool_size: 10           # 连接池最大连接数
  min_idle_conns: 5       # 最小空闲连接数
  max_conn_age: "1h"      # 连接最大存活时间
  pool_timeout: "30s"     # 从连接池获取连接的超时时间
  idle_timeout: "5m"      # 空闲连接的超时时间
  idle_check_frequency: "1m" # 空闲连接检查频率
  read_timeout: "3s"      # 读取超时
  write_timeout: "3s"     # 写入超时
  dial_timeout: "5s"      # 连接超时
  max_retries: 3          # 最大重试次数
  min_retry_backoff: "1ms" # 最小重试间隔
  max_retry_backoff: "500ms" # 最大重试间隔

# JWT配置
jwt:
  secret: "${JWT_SECRET}" # JWT密钥（从环境变量读取）
  expire_hours: 24        # JWT过期时间（小时）
  issuer: "rbac-admin"    # 颁发者
  audience: "rbac-admin"  # 受众
  signing_method: "HS256" # 签名方法
  token_name: "Authorization" # Token名称

# 日志配置
log:
  level: "info"           # 日志级别: debug/info/warn/error/fatal
  dir: "./logs"           # 日志目录
  filename: "app.log"     # 日志文件名
  max_size: 100           # 单个日志文件最大大小（MB）
  max_backups: 10         # 保留的最大日志文件数
  max_age: 30             # 日志文件最大保留天数
  compress: true          # 是否压缩日志文件
  stdout: true            # 是否输出到标准输出

# 安全配置
security:
  cors_origins: ["*"]    # 允许的CORS源
  csrf_secret: "${CSRF_SECRET}" # CSRF密钥（从环境变量读取）
  xss_protection: true    # 启用XSS保护
  frame_options: "DENY"   # X-Frame-Options
  content_security_policy: "default-src 'self'" # 内容安全策略
  rate_limit: 100         # 每分钟请求限制
  brute_force_protection: true # 启用暴力破解保护
  password_complexity: 8  # 密码最小长度
  login_attempts_limit: 5 # 登录尝试次数限制
  login_lockout_time: 30  # 登录锁定时间（分钟）

# CORS配置
cors:
  allow_origins: ["*"]    # 允许的源
  allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"] # 允许的HTTP方法
  allow_headers: ["Origin", "Content-Type", "Authorization"] # 允许的HTTP头
  allow_credentials: true # 是否允许凭证
  expose_headers: []      # 暴露的HTTP头
  max_age: 600            # 预检请求缓存时间（秒）

# 性能配置
performance:
  max_request_size: 10    # 最大请求大小（MB）
  request_timeout: "30s"  # 请求超时时间
  response_compression: true # 启用响应压缩
  gzip_level: 6           # Gzip压缩级别
  cache_control: "no-cache" # 缓存控制
  etag: true              # 启用ETag

# 上传配置
upload:
  dir: "./uploads"        # 上传文件目录
  max_size: 50            # 最大上传文件大小（MB）
  allowed_types: ["image/jpeg", "image/png", "application/pdf", "application/zip"] # 允许的文件类型
  file_permissions: 0644  # 文件权限
  dir_permissions: 0755   # 目录权限
```

### 环境变量配置

生产环境必须使用环境变量管理敏感信息：

```bash
# 创建.env文件
touch .env
```

.env文件内容示例：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=rbac_user
DB_PASSWORD=your_secure_password_here
DB_NAME=rbac_admin_prod

# JWT配置
JWT_SECRET=your_jwt_secret_key_minimum_32_characters
JWT_EXPIRE_HOURS=24
JWT_ISSUER=rbac-admin
JWT_AUDIENCE=rbac-admin

# Redis配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password_here

# 安全配置
CSRF_SECRET=your_csrf_secret_key_here
CORS_ORIGINS=https://your-domain.com

# 系统配置
SYSTEM_PORT=8080
SYSTEM_IP=0.0.0.0
```

## 🔧 管理员用户创建

### 命令行创建管理员用户

```bash
go run main.go -m user -t create
```

按照提示输入用户名和密码：

```
请输入用户名
superadmin
请输入密码
请再次输入密码
创建用户成功
```

### 默认管理员账号

如果首次启动时没有创建管理员用户，系统会自动创建一个默认管理员账号（仅开发环境）：
- 用户名: admin
- 密码: admin123
- 首次登录后请立即修改密码

## 🌐 API接口文档

### API基础路径

所有API接口的基础路径为：`http://your-server:8080/api/v1/`

### 认证接口

#### 登录

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "your-username",
  "password": "your-password"
}
```

#### 退出登录

```
POST /api/v1/auth/logout
Authorization: Bearer your-jwt-token
```

#### 刷新令牌

```
POST /api/v1/auth/refresh
Authorization: Bearer your-jwt-token
```

### 用户管理接口

#### 获取用户列表

```
GET /api/v1/users
Authorization: Bearer your-jwt-token
```

#### 创建用户

```
POST /api/v1/users
Authorization: Bearer your-jwt-token
Content-Type: application/json

{
  "username": "new-user",
  "password": "new-password",
  "nickname": "新用户",
  "status": 1,
  "department_id": 1,
  "role_ids": [1, 2]
}
```

#### 获取用户详情

```
GET /api/v1/users/{id}
Authorization: Bearer your-jwt-token
```

#### 更新用户

```
PUT /api/v1/users/{id}
Authorization: Bearer your-jwt-token
Content-Type: application/json

{
  "nickname": "更新后的用户",
  "status": 1,
  "department_id": 2,
  "role_ids": [1]
}
```

#### 删除用户

```
DELETE /api/v1/users/{id}
Authorization: Bearer your-jwt-token
```

### 角色管理接口

#### 获取角色列表

```
GET /api/v1/roles
Authorization: Bearer your-jwt-token
```

#### 创建角色

```
POST /api/v1/roles
Authorization: Bearer your-jwt-token
Content-Type: application/json

{
  "name": "新角色",
  "description": "角色描述",
  "status": 1,
  "permission_ids": [1, 2, 3]
}
```

#### 获取角色详情

```
GET /api/v1/roles/{id}
Authorization: Bearer your-jwt-token
```

#### 更新角色

```
PUT /api/v1/roles/{id}
Authorization: Bearer your-jwt-token
Content-Type: application/json

{
  "name": "更新后的角色",
  "description": "更新后的描述",
  "status": 1,
  "permission_ids": [1, 2]
}
```

#### 删除角色

```
DELETE /api/v1/roles/{id}
Authorization: Bearer your-jwt-token
```

## 🛠️ 开发指南

### 项目结构

```
rbac_admin_server/
├── config/            # 配置相关
├── core/              # 核心功能
├── global/            # 全局变量
├── middleware/        # 中间件
├── models/            # 数据模型
├── pwd/               # 密码处理
├── routes/            # 路由定义
├── utils/             # 工具函数
├── main.go            # 入口文件
├── go.mod             # Go模块定义
├── go.sum             # 依赖版本锁定
├── settings.yaml      # 配置文件
└── RbacAdminServer项目部署文档.md # 部署文档
```

### 开发流程

1. **克隆项目代码**

```bash
git clone <your-repo-url>
cd rbac_admin_server
go mod download
go mod tidy
```

2. **配置开发环境**

创建开发环境配置文件：

```bash
copy settings.yaml settings_dev.yaml
```

3. **启动开发服务器**

```bash
go run main.go -env dev
```

4. **代码规范**

- 使用Go的标准代码格式化工具：`gofmt -s -w .`
- 运行静态代码分析：`go vet ./...`
- 运行单元测试：`go test ./... -v`

### 调试技巧

1. **启用调试日志**

在配置文件中将`log.level`设置为`debug`，可以查看更详细的日志信息。

2. **使用Delve调试器**

```bash
dlv debug main.go -- -env dev
```

3. **常见问题排查**

- 数据库连接失败：检查数据库服务是否启动，配置是否正确
- Redis连接失败：检查Redis服务是否启动，配置是否正确
- 端口占用：修改配置文件中的`system.port`值

## 🚨 故障排除

### 常见问题及解决方法

1. **数据库连接失败**

症状：启动时出现`数据库连接失败`错误

解决方法：
- 检查数据库服务是否启动
- 验证数据库配置是否正确（主机、端口、用户名、密码）
- 确认数据库用户是否有足够的权限

2. **Redis连接失败**

症状：启动时出现`Redis连接失败`错误

解决方法：
- 检查Redis服务是否启动
- 验证Redis配置是否正确（地址、密码）
- 确认Redis防火墙设置是否允许连接

3. **端口被占用**

症状：启动时出现`address already in use`错误

解决方法：
- 修改配置文件中的`system.port`值
- 停止占用该端口的其他进程

4. **JWT认证失败**

症状：API请求返回`401 Unauthorized`错误

解决方法：
- 检查JWT令牌是否过期
- 确认请求头中的`Authorization`字段格式正确（`Bearer token`）
- 验证配置文件中的JWT密钥是否正确

5. **权限不足**

症状：API请求返回`403 Forbidden`错误

解决方法：
- 检查当前用户是否有足够的权限执行该操作
- 确认角色和权限配置是否正确

## 📊 性能优化

### 数据库优化

1. **创建索引**

为常用查询字段创建索引，例如：

```sql
-- 用户表索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_status ON users(status);

-- 角色表索引
CREATE INDEX idx_roles_status ON roles(status);

-- 权限表索引
CREATE INDEX idx_permissions_name ON permissions(name);
```

2. **优化查询**

- 使用`SELECT`指定需要的字段，避免`SELECT *`
- 使用分页查询，限制返回结果数量
- 合理使用预加载（Preload）减少N+1查询问题

### Redis优化

1. **设置合理的过期时间**

根据数据特性设置合适的过期时间，避免Redis内存占用过大。

2. **使用管道操作**

批量执行Redis命令，减少网络往返次数。

3. **使用连接池**

配置合理的连接池大小，避免频繁创建和关闭连接。

## 📝 版本历史

### v1.0.0 (2024-xx-xx)

- 首次发布
- 支持用户、角色、权限、部门、菜单管理
- 支持MySQL、PostgreSQL、SQLite数据库
- 支持Redis缓存
- 支持多环境配置
- 完善的API接口

## 📚 附录

### 数据库表结构

主要表结构关系：

- **users**: 用户表
- **roles**: 角色表
- **permissions**: 权限表
- **user_roles**: 用户-角色关联表
- **role_permissions**: 角色-权限关联表
- **departments**: 部门表
- **menus**: 菜单表
- **role_menus**: 角色-菜单关联表
- **casbin_rules**: Casbin规则表（用于权限控制）

### 配置文件环境变量替换规则

配置文件中以`${环境变量名}`格式的字符串会被自动替换为对应的环境变量值。如果环境变量不存在，则保留原字符串。

例如：
```yaml
password: "${DB_PASSWORD}" # 如果存在环境变量DB_PASSWORD，则替换为其值
```

### 安全建议

1. **生产环境安全配置**

- 不要在配置文件中明文存储密码等敏感信息，使用环境变量
- 配置强密码策略
- 启用HTTPS加密通信
- 限制访问IP和端口

2. **定期维护**

- 定期更新依赖包，修复已知漏洞
- 定期备份数据库
- 定期检查系统日志，发现异常及时处理

3. **监控建议**

- 监控系统资源使用情况（CPU、内存、磁盘）
- 监控API请求量和响应时间
- 监控数据库连接和查询性能
- 监控Redis内存使用和连接情况
=======
# RBAC管理员服务器部署文档

## 1. 项目概况

RBAC管理员服务器是一款基于角色的访问控制系统（Role-Based Access Control），旨在为企业级应用提供灵活、安全、可扩展的权限管理解决方案。该系统通过精确的角色划分和权限控制，帮助企业构建多层级的安全访问体系，有效保障数据安全和业务流程合规性。

**核心功能特性：**
- 完善的用户管理体系，支持用户CRUD、状态控制和信息维护
- 多维度角色定义与管理，实现精细化权限控制
- 细粒度的权限点管理，支持API、菜单、按钮等多维度权限控制
- 可视化的菜单管理，支持多级菜单和动态路由
- 部门组织结构管理，适配企业级组织架构
- 完善的JWT认证机制，支持令牌刷新和过期控制
- 全面的系统监控与健康检查
- 支持多环境配置（开发、测试、生产）和环境变量覆盖

**典型应用场景：**
- 企业内部管理系统的权限控制
- SaaS平台的多租户权限隔离
- 政府及金融行业的分级访问控制
- 需要严格权限审计的业务系统

## 2. 技术栈

| 技术/组件 | 版本/说明 | 用途 | 溯源 |
|---------|----------|------|------|
| Go | 1.19+ | 主要开发语言 | <mcfile name="go.mod" path="e:\myblog\Go项目学习\rbac_admin_server\go.mod"></mcfile> |
| GORM | v2 | 数据库ORM框架 | <mcfile name="core\db.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\db.go"></mcfile> |
| MySQL | 5.7+/8.0+ | 关系型数据库（默认） | <mcfile name="config\config.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\config.go"></mcfile> |
| PostgreSQL | 可选 | 关系型数据库 | <mcfile name="core\db.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\db.go"></mcfile> |
| SQLite | 可选 | 轻量级数据库 | <mcfile name="core\db.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\db.go"></mcfile> |
| Redis | 6.0+ | 缓存服务 | <mcfile name="core\redis.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\redis.go"></mcfile> |
| JWT | 自定义实现 | 身份认证 | <mcfile name="config\config.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\config.go"></mcfile> |
| Swagger | 内置 | API文档自动生成 | <mcfile name="config\config.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\config.go"></mcfile> |
| logrus | v1 | 结构化日志系统 | <mcfile name="core\init.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\init.go"></mcfile> |
| godotenv | 最新版 | 环境变量加载 | <mcfile name="config\loader.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\loader.go"></mcfile> |
| yaml.v3 | 最新版 | 配置文件解析 | <mcfile name="config\loader.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\loader.go"></mcfile> |

## 3. 部署环境准备

### 3.1 硬件要求

| 环境类型 | CPU | 内存 | 存储空间 | 网络 |
|---------|-----|------|---------|------|
| 开发环境 | 2核 | 4GB | 50GB | 千兆网卡 |
| 测试环境 | 4核 | 8GB | 100GB | 千兆网卡 |
| 生产环境 | 8核+ | 16GB+ | 500GB+ | 千兆网卡/万兆网卡 |

### 3.2 软件要求

| 环境类型 | Go版本 | 数据库 | Redis | 操作系统 |
|---------|-------|--------|-------|---------|
| 开发环境 | 1.19+ | SQLite/MySQL | 可选 | Windows/Linux/macOS |
| 测试环境 | 1.19+ | MySQL 5.7+/PostgreSQL 12+ | 6.0+ | Linux |
| 生产环境 | 1.19+ | MySQL 8.0+/PostgreSQL 14+ | 6.0+ | Linux (CentOS/RHEL 7+/Ubuntu 20.04+) |

### 3.3 网络要求

- 生产环境需配置防火墙规则，仅开放必要端口（默认8080）
- 数据库服务器应配置内网访问限制，禁止外部直接访问
- Redis服务应配置密码认证和访问控制
- 建议为API服务配置HTTPS加密访问

## 4. 项目目录结构和组件

### 4.1 整体目录结构

```
rbac_admin_server/
├── api/            # API接口实现
├── config/         # 配置管理
│   ├── config.go   # 配置结构体定义
│   └── loader.go   # 配置加载逻辑
├── core/           # 核心功能模块
│   ├── db.go       # 数据库连接管理
│   ├── init.go     # 系统初始化
│   ├── logger.go   # 日志系统
│   ├── redis.go    # Redis连接管理
│   └── validator.go # 数据验证器
├── global/         # 全局变量和配置
│   └── global.go   # 全局变量定义
├── middleware/     # HTTP中间件
├── models/         # 数据模型定义
│   ├── api.go      # API模型
│   ├── base.go     # 基础模型
│   ├── menu.go     # 菜单模型
│   └── user.go     # 用户模型
├── routes/         # 路由定义
├── service/        # 业务逻辑层
├── utils/          # 工具函数
├── .env            # 环境变量配置（本地开发用，不提交）
├── .env.example    # 环境变量配置示例
├── config.go       # 配置入口文件
├── deploy.bat      # Windows部署脚本
├── deploy.sh       # Linux部署脚本
├── go.mod          # Go模块定义
├── go.sum          # 依赖版本锁定
├── main.go         # 程序入口
├── settings.yaml   # 主配置文件
├── settings_dev.yaml # 开发环境配置
├── settings_test.yaml # 测试环境配置
└── settings_prod.yaml # 生产环境配置
```

### 4.2 核心组件说明

#### 4.2.1 配置管理组件

配置管理组件负责加载、解析和管理系统配置，支持多环境配置文件、环境变量覆盖和默认值处理。配置系统基于YAML格式，结合环境变量实现配置的灵活调整和安全管理。

**主要文件：**
- <mcfile name="config\config.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\config.go"></mcfile>: 定义了完整的配置结构体
- <mcfile name="config\loader.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\loader.go"></mcfile>: 实现了配置加载和环境变量替换逻辑
- <mcfile name="global\global.go" path="e:\myblog\Go项目学习\rbac_admin_server\global\global.go"></mcfile>: 提供全局配置访问点

#### 4.2.2 数据库组件

数据库组件基于GORM框架，支持MySQL、PostgreSQL和SQLite三种数据库，实现了连接池配置、健康检查和自动迁移功能。

**主要文件：**
- <mcfile name="core\db.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\db.go"></mcfile>: 数据库初始化和连接管理
- <mcfile name="core\init.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\init.go"></mcfile>: 数据库表结构自动迁移
- <mcfile name="models\" path="e:\myblog\Go项目学习\rbac_admin_server\models\"></mcfile>: 数据模型定义

#### 4.2.3 Redis缓存组件

Redis组件提供了缓存服务的初始化、连接管理和常用操作封装，支持连接池配置和连接测试。

**主要文件：**
- <mcfile name="core\redis.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\redis.go"></mcfile>: Redis客户端初始化和操作封装

#### 4.2.4 系统初始化组件

系统初始化组件负责按顺序初始化日志、验证器、数据库和Redis等核心服务，是系统启动的关键环节。

**主要文件：**
- <mcfile name="core\init.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\init.go"></mcfile>: 系统初始化流程和资源管理

## 5. 项目配置详解

### 5.1 配置文件结构

系统采用分层配置结构，支持多环境配置文件和环境变量覆盖，配置项分为12个主要模块：

```yaml
# 系统基本配置
system:
  ip: 127.0.0.1    # 监听IP
  port: 8080       # 监听端口

# 数据库配置
db:
  mode: mysql       # 数据库类型: mysql, postgres, sqlite
  host: localhost   # 数据库主机
  port: 3306        # 数据库端口
  user: root        # 用户名
  password: ""      # 密码（建议使用环境变量）
  dbname: rbac_admin # 数据库名称
  # 连接池配置
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h
  conn_max_idle_time: 30m

# Redis配置
redis:
  addr: localhost:6379 # Redis地址
  password: ""       # Redis密码
  db: 0              # 数据库编号

# JWT配置
jwt:
  secret: ""         # JWT密钥（必须保密）
  expire_hours: 24   # 令牌有效期
  refresh_expire_hours: 168 # 刷新令牌有效期
  issuer: rbac-admin-server
  audience: rbac-client

# 其他配置模块...
log: {}
security: {}
performance: {}
upload: {}
monitoring: {}
cors: {}
swagger: {}
app: {}
```

### 5.2 多环境配置策略

系统支持三种运行环境配置：

| 环境类型 | 配置文件 | 主要特点 | 适用场景 |
|---------|---------|---------|---------|
| 开发环境 | settings_dev.yaml | 日志级别低、调试模式开启、SQLite或本地MySQL | 本地开发和调试 |
| 测试环境 | settings_test.yaml | 严格的日志记录、测试数据库、性能监控开启 | 集成测试和质量验证 |
| 生产环境 | settings_prod.yaml | 日志级别高、关闭调试模式、生产数据库集群、安全配置严格 | 正式线上环境 |

环境切换通过命令行参数 `-env` 控制，例如：`go run main.go -env prod`

### 5.3 环境变量支持

系统支持通过环境变量覆盖配置文件中的设置，主要环境变量如下：

| 环境变量名称 | 对应配置项 | 说明 |
|------------|-----------|------|
| SYSTEM_IP | system.ip | 服务器监听IP |
| SYSTEM_PORT | system.port | 服务器监听端口 |
| DB_MODE | db.mode | 数据库类型 |
| DB_HOST | db.host | 数据库主机 |
| DB_PORT | db.port | 数据库端口 |
| DB_USERNAME | db.user | 数据库用户名 |
| DB_PASSWORD | db.password | 数据库密码 |
| DB_NAME | db.dbname | 数据库名称 |
| REDIS_ADDR | redis.addr | Redis地址 |
| REDIS_PASSWORD | redis.password | Redis密码 |
| JWT_SECRET | jwt.secret | JWT签名密钥 |
| JWT_EXPIRE_HOURS | jwt.expire_hours | JWT过期时间（小时） |
| LOG_LEVEL | log.level | 日志级别 |
| LOG_DIR | log.log_dir | 日志目录 |
| APP_ENVIRONMENT | app.environment | 应用环境 |
| APP_DEBUG | app.debug | 调试模式 |

环境变量可通过 `.env` 文件加载，推荐在生产环境中使用环境变量配置敏感信息。

### 5.4 配置优先级

系统配置加载时遵循以下优先级（从高到低）：
1. 命令行参数指定的配置文件
2. 环境变量设置的值
3. 根据环境选择的配置文件（settings_dev.yaml/settings_test.yaml/settings_prod.yaml）
4. 代码中的默认配置值

## 6. 部署步骤

### 6.1 开发环境部署

**前置条件：**
- 已安装Go 1.19+环境
- 已安装Git

**部署步骤：**

1. **克隆代码仓库**
   ```bash
   git clone <repository-url> rbac_admin_server
   cd rbac_admin_server
   ```

2. **配置环境**
   ```bash
   # 复制环境变量示例文件
   cp .env.example .env
   # 根据本地环境修改.env文件
   # 对于Windows环境，使用复制命令
   copy .env.example .env
   ```

3. **安装依赖**
   ```bash
   go mod download
   go mod tidy
   ```

4. **启动服务**
   ```bash
   go run main.go -env dev
   # Windows环境下也可以使用部署脚本
   deploy.bat
   ```

5. **验证服务**
   - 访问：http://localhost:8080/health 验证服务健康状态
   - 访问：http://localhost:8080/swagger/index.html 查看API文档

### 6.2 测试环境部署

**前置条件：**
- 已安装Go 1.19+环境
- 已安装MySQL或PostgreSQL数据库
- 已安装Redis服务

**部署步骤：**

1. **准备数据库**
   ```sql
   CREATE DATABASE rbac_admin_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
   ```

2. **配置环境**
   ```bash
   # 设置环境变量或修改settings_test.yaml
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USERNAME=test_user
   export DB_PASSWORD=test_password
   export DB_NAME=rbac_admin_test
   export REDIS_ADDR=localhost:6379
   ```

3. **构建和启动**
   ```bash
   # Linux/Mac环境
   chmod +x deploy.sh
   ./deploy.sh
   
   # 或手动构建
   go build -o rbac_admin_server .
   ./rbac_admin_server -env test
   ```

4. **运行验证**
   - 检查服务日志是否有错误
   - 验证数据库连接和表结构是否正确创建
   - 执行自动化测试（如适用）

### 6.3 生产环境部署

**前置条件：**
- 已准备生产服务器（推荐Linux系统）
- 已配置MySQL/PostgreSQL数据库集群
- 已配置Redis缓存服务
- 已配置HTTPS证书（如需要）

**部署步骤：**

1. **系统准备**
   ```bash
   # 更新系统包
   sudo apt-get update && sudo apt-get upgrade -y  # Ubuntu/Debian
   sudo yum update -y  # CentOS/RHEL
   
   # 安装必要工具
   sudo apt-get install -y git curl wget  # Ubuntu/Debian
   sudo yum install -y git curl wget  # CentOS/RHEL
   ```

2. **安装Go环境**
   ```bash
   # 下载并安装Go 1.19+
   wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
   
   # 配置环境变量
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
   echo 'export GOPATH=$HOME/go' >> ~/.profile
   source ~/.profile
   
   # 验证安装
   go version
   ```

3. **创建部署目录**
   ```bash
   sudo mkdir -p /opt/rbac_admin_server
   sudo chown -R $USER:$USER /opt/rbac_admin_server
   cd /opt/rbac_admin_server
   ```

4. **获取代码**
   ```bash
   git clone <repository-url> .
   ```

5. **配置生产环境**
   ```bash
   # 创建生产环境配置文件
   cp settings_prod.yaml.example settings_prod.yaml
   
   # 编辑配置文件，设置生产环境参数
   vim settings_prod.yaml
   
   # 设置环境变量（敏感信息）
   cat > .env << 'EOF'
   # 数据库配置
   DB_PASSWORD=your-secure-password
   
   # Redis配置
   REDIS_PASSWORD=your-redis-password
   
   # JWT配置
   JWT_SECRET=your-32-character-secret-key
   
   # 应用环境
   APP_ENVIRONMENT=production
   APP_DEBUG=false
   EOF
   
   # 配置文件权限
   chmod 600 .env settings_prod.yaml
   ```

6. **构建应用**
   ```bash
   # 构建生产版本
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o rbac_admin_server
   
   # 赋予执行权限
   chmod +x rbac_admin_server
   ```

7. **配置Systemd服务**
   ```bash
   # 创建systemd服务文件
   sudo vim /etc/systemd/system/rbac_admin_server.service
   ```

   内容如下：
   ```ini
   [Unit]
   Description=RBAC Admin Server
   After=network.target mysql.service redis.service
   Requires=mysql.service redis.service
   
   [Service]
   Type=simple
   User=your-user
   WorkingDirectory=/opt/rbac_admin_server
   ExecStart=/opt/rbac_admin_server/rbac_admin_server -env prod
   Restart=on-failure
   RestartSec=5s
   EnvironmentFile=/opt/rbac_admin_server/.env
   
   [Install]
   WantedBy=multi-user.target
   ```

8. **启动服务并设置开机自启**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start rbac_admin_server
   sudo systemctl enable rbac_admin_server
   
   # 查看服务状态
   sudo systemctl status rbac_admin_server
   
   # 查看日志
   journalctl -u rbac_admin_server -f
   ```

9. **配置Nginx反向代理（可选）**
   ```bash
   # 安装Nginx
   sudo apt-get install nginx  # Ubuntu/Debian
   sudo yum install nginx  # CentOS/RHEL
   
   # 创建Nginx配置
   sudo vim /etc/nginx/sites-available/rbac_admin_server
   ```

   内容如下：
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       
       # 重定向到HTTPS（如果已配置证书）
       return 301 https://$host$request_uri;
   }
   
   server {
       listen 443 ssl;
       server_name your-domain.com;
       
       # SSL证书配置
       ssl_certificate /path/to/your/cert.pem;
       ssl_certificate_key /path/to/your/key.pem;
       
       # SSL优化配置
       ssl_protocols TLSv1.2 TLSv1.3;
       ssl_prefer_server_ciphers off;
       ssl_session_timeout 1d;
       ssl_session_cache shared:SSL:10m;
       ssl_session_tickets off;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
       
       # 健康检查端点不缓存
       location /health {
           proxy_pass http://localhost:8080;
           proxy_cache_bypass 1;
       }
   }
   ```

   启用配置：
   ```bash
   sudo ln -s /etc/nginx/sites-available/rbac_admin_server /etc/nginx/sites-enabled/
   sudo nginx -t  # 测试配置
   sudo systemctl reload nginx  # 重启Nginx
   ```

10. **验证生产环境部署**
    - 访问：https://your-domain.com/health 验证服务健康状态
    - 检查日志是否有异常
    - 执行基本功能验证

## 7. 项目代码运行加载流程

RBAC管理员服务器的启动和运行遵循以下完整流程：

### 7.1 启动流程

1. **命令行参数解析**
   - 解析 `-env` 参数确定运行环境（dev/test/prod）
   - 解析 `-config` 参数确定配置文件路径（可选）
   <mcfile name="main.go" path="e:\myblog\Go项目学习\rbac_admin_server\main.go"></mcfile>

2. **配置加载**
   - 根据环境或指定路径加载对应的YAML配置文件
   - 加载 `.env` 文件中的环境变量
   - 替换配置文件中的环境变量占位符
   - 应用默认配置值
   - 验证配置有效性
   <mcfile name="config\loader.go" path="e:\myblog\Go项目学习\rbac_admin_server\config\loader.go"></mcfile>

3. **系统初始化**
   - 初始化日志系统（基于logrus）
   - 初始化验证器
   - 初始化数据库连接（基于GORM）
   - 初始化Redis连接
   - 自动迁移数据库表结构
   <mcfile name="core\init.go" path="e:\myblog\Go项目学习\rbac_admin_server\core\init.go"></mcfile>

4. **服务启动**
   - 启动HTTP服务
   - 注册API路由
   - 注册中间件
   - 启动监控和健康检查端点
   - 启动Swagger文档服务（如启用）

5. **等待退出信号**
   - 监听系统信号（SIGINT、SIGTERM）
   - 收到信号后执行优雅关闭
   - 清理系统资源（数据库连接、Redis连接等）

### 7.2 配置加载机制详解

配置加载是系统启动的关键环节，具体流程如下：

1. **确定配置文件**：根据 `-env` 参数或 `-config` 参数确定要加载的配置文件
2. **读取配置文件**：使用 `os.ReadFile` 读取YAML配置文件内容
3. **环境变量替换**：替换配置文件中的 `${ENV_VAR}` 格式占位符
4. **YAML解析**：使用 `yaml.v3` 包将配置内容解析到 `Config` 结构体
5. **环境变量覆盖**：读取系统环境变量，覆盖对应配置项的值
6. **应用默认值**：为未设置的配置项应用默认值
7. **配置验证**：验证关键配置项的有效性（如端口范围、JWT密钥等）
8. **设置全局配置**：将加载的配置设置到全局变量中供系统使用

### 7.3 数据库初始化流程

数据库初始化是系统可用性的基础，具体流程如下：

1. **确定数据库类型**：根据配置中的 `db.mode` 确定使用的数据库类型
2. **构建连接字符串**：根据数据库类型和配置参数构建连接字符串
3. **配置GORM**：设置命名策略、外键约束、日志级别等GORM配置
4. **打开数据库连接**：调用 `gorm.Open` 打开数据库连接
5. **配置连接池**：设置最大连接数、最大空闲连接数、连接生命周期等参数
6. **测试连接**：执行 `Ping` 操作测试数据库连接是否成功
7. **自动迁移**：调用 `AutoMigrate` 自动创建或更新数据库表结构

## 8. 部署模式介绍

RBAC管理员服务器支持多种部署模式，可根据企业实际需求选择合适的部署方案。

### 8.1 单机部署

**特点**：所有服务组件部署在同一台服务器上，结构简单，适合小型应用或开发测试环境。

**架构**：
- 应用服务 + 数据库 + Redis 均部署在同一服务器
- 使用本地文件系统存储日志和上传文件

**适用场景**：小型项目、开发环境、测试环境、资源受限的场景

### 8.2 传统服务器集群部署

**特点**：应用服务部署在多台服务器上，通过负载均衡器分发请求，数据库和Redis独立部署。

**架构**：
- 多台应用服务器运行RBAC管理员服务实例
- 前端部署Nginx作为负载均衡器
- 独立的数据库服务器或数据库集群
- 独立的Redis服务器或Redis集群
- 可选的共享文件存储服务

**适用场景**：中大型应用、生产环境、需要高可用性的业务场景

### 8.3 容器化部署（Docker）

**特点**：使用Docker容器封装应用及其依赖，实现标准化部署和快速扩缩容。

**架构**：
- 应用服务打包为Docker镜像
- 使用Docker Compose或Kubernetes编排容器
- 数据库和Redis可使用容器化部署或独立部署
- 使用Docker卷或外部存储服务持久化数据

**示例Dockerfile**：
```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o rbac_admin_server

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/rbac_admin_server .
COPY settings_prod.yaml .
COPY .env .

RUN mkdir -p logs uploads

EXPOSE 8080

CMD ["./rbac_admin_server", "-env", "prod"]
```

**适用场景**：现代化部署环境、云原生应用、需要快速扩缩容的场景

### 8.4 云服务部署

**特点**：利用云服务提供商的托管服务，降低运维成本，提高系统可靠性。

**架构**：
- 应用服务部署在云服务器或容器服务中（如ECS、EKS、ECS）
- 使用云数据库服务（如RDS for MySQL/PostgreSQL）
- 使用云缓存服务（如ElastiCache for Redis）
- 使用云存储服务（如S3、OSS）存储静态资源和日志
- 使用云负载均衡服务（如ELB、SLB）分发流量
- 可选使用容器服务编排（如EKS、ACK）

**适用场景**：企业级应用、云原生应用、对可靠性和可扩展性有高要求的场景

## 9. 安全与性能优化

### 9.1 安全配置

1. **敏感信息管理**
   - 所有敏感信息（如数据库密码、JWT密钥）必须通过环境变量或密钥管理服务配置
   - 禁止在代码或配置文件中硬编码敏感信息
   - 定期轮换密码和密钥

2. **访问控制**
   - 生产环境应禁用Swagger UI和调试模式
   - 配置严格的防火墙规则，仅开放必要端口
   - 启用CSRF保护和XSS防护
   - 配置API访问频率限制

3. **数据安全**
   - 数据库连接使用SSL加密（如适用）
   - 敏感数据传输使用HTTPS
   - 用户密码使用bcrypt加密存储
   - 定期备份数据库

4. **日志安全**
   - 避免在日志中记录敏感信息
   - 配置适当的日志级别，生产环境建议使用INFO或WARN级别
   - 日志文件权限设置为仅管理员可读

### 9.2 性能优化

1. **数据库优化**
   - 合理配置数据库连接池参数
   - 为常用查询添加索引
   - 优化复杂查询，避免全表扫描
   - 考虑读写分离架构（如适用）

2. **缓存优化**
   - 合理使用Redis缓存热点数据
   - 设置适当的缓存过期时间
   - 实现缓存预热和缓存穿透防护
   - 考虑使用分布式缓存（如适用）

3. **应用性能优化**
   - 启用响应压缩
   - 配置合理的请求超时时间
   - 实现请求限流和熔断机制
   - 定期进行性能分析和优化

4. **部署优化**
   - 使用负载均衡提高并发处理能力
   - 根据业务需求合理配置服务器资源
   - 考虑使用CDN加速静态资源访问（如适用）
   - 配置自动扩缩容策略（如适用）

## 10. 监控与维护

### 10.1 系统监控

1. **健康检查**
   - 系统提供 `/health` 端点用于健康检查
   - 可集成到监控系统（如Prometheus、Zabbix）进行定期检查
   - 健康检查结果包含数据库和Redis连接状态

2. **性能监控**
   - 系统提供 `/metrics` 端点输出监控指标
   - 推荐集成Prometheus和Grafana构建监控仪表盘
   - 重点监控CPU使用率、内存占用、请求响应时间、数据库连接数等指标

3. **日志监控**
   - 定期检查系统日志，及时发现异常
   - 推荐使用ELK、Loki等日志收集和分析工具
   - 配置日志告警规则，对错误和异常日志及时通知

### 10.2 常规维护

1. **备份策略**
   - 数据库定期全量备份（至少每日一次）
   - 重要配置文件和数据定期备份
   - 备份数据存储在安全的异地位置
   - 定期验证备份的可用性

2. **更新与升级**
   - 定期更新系统依赖包，修复安全漏洞
   - 按照规范流程进行版本升级
   - 升级前进行充分测试，制定回滚方案
   - 记录更新内容和变更日志

3. **常见问题排查**
   - 服务无法启动：检查配置文件、数据库连接、端口占用情况
   - 数据库连接失败：检查数据库服务状态、连接参数、网络连通性
   - Redis连接失败：检查Redis服务状态、密码配置、网络连通性
   - API请求错误：检查日志中的错误信息、参数有效性、权限配置
   - 性能问题：检查系统资源使用情况、数据库查询效率、缓存命中率

## 11. 附录

### 11.1 配置项速查表

| 配置模块 | 核心配置项 | 默认值 | 说明 |
|---------|-----------|-------|------|
| system | port | 8080 | 服务器监听端口 |
| db | mode | sqlite | 数据库类型 |
| db | max_open_conns | 100 | 数据库最大连接数 |
| db | max_idle_conns | 10 | 数据库最大空闲连接数 |
| jwt | secret | - | JWT签名密钥（必填） |
| jwt | expire_hours | 24 | JWT令牌有效期（小时） |
| log | level | info | 日志级别 |
| log | log_dir | ./logs | 日志存储目录 |
| security | bcrypt_cost | 10 | 密码加密强度 |
| security | max_login_attempts | 5 | 最大登录尝试次数 |
| monitoring | health_check_path | /health | 健康检查路径 |
| monitoring | metrics_path | /metrics | 指标收集路径 |
| swagger | enable | true | 是否启用Swagger |
| swagger | enable_ui | true | 是否启用Swagger UI |

### 11.2 开发与部署命令速查表

| 命令 | 说明 | 适用环境 |
|-----|------|---------|
| go run main.go -env dev | 启动开发环境服务 | 开发环境 |
| go run main.go -env test | 启动测试环境服务 | 测试环境 |
| go run main.go -env prod | 启动生产环境服务 | 生产环境 |
| go build -o rbac_admin_server . | 构建应用程序 | 所有环境 |
| ./deploy.sh | 执行Linux部署脚本 | Linux环境 |
| deploy.bat | 执行Windows部署脚本 | Windows环境 |
| systemctl start rbac_admin_server | 启动Systemd服务 | Linux生产环境 |
| systemctl status rbac_admin_server | 查看服务状态 | Linux生产环境 |
| journalctl -u rbac_admin_server -f | 查看服务日志 | Linux生产环境 |

### 11.3 常见部署问题解决指南

1. **JWT密钥长度不足32位**
   - 问题：启动时提示JWT密钥不足32位
   - 解决：生成至少32位的随机密钥，通过环境变量JWT_SECRET设置
   - 示例：`JWT_SECRET=$(openssl rand -hex 16)`

2. **数据库连接超时**
   - 问题：无法连接到数据库服务器
   - 解决：检查数据库服务是否运行、连接参数是否正确、网络是否连通
   - 验证：使用命令行工具（如mysql、psql）测试数据库连接

3. **Redis连接失败**
   - 问题：无法连接到Redis服务器
   - 解决：检查Redis服务是否运行、密码是否正确、网络是否连通
   - 验证：使用redis-cli工具测试Redis连接

4. **端口被占用**
   - 问题：启动时提示端口已被占用
   - 解决：通过环境变量SYSTEM_PORT或配置文件修改监听端口
   - 查找：使用`netstat -tuln | grep <port>`或`lsof -i :<port>`查找占用端口的进程

5. **权限不足**
   - 问题：无法读取配置文件或写入日志
   - 解决：检查文件和目录权限，确保运行服务的用户有相应权限
   - 修复：使用`chmod`和`chown`命令调整权限

---

**文档版本**: v1.0.0
**发布日期**: 2023年12月
**适用项目**: RBAC管理员服务器
>>>>>>> 03404e4ec9e063e0d69d4af944091fb9ab46f525
