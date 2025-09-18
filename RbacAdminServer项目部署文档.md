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