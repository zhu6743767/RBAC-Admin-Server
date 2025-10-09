# RBAC Admin Server 项目部署文档

## 1. 文档版本信息

| 版本 | 发布日期 | 更新内容 | 责任人 |
|------|---------|---------|--------|
| v1.3.0 | 2024-04-15 | 修复配置文档与实际代码不一致的问题，更新项目结构说明，修正配置项速查表 | 开发团队 |
| v1.2.0 | 2024-04-10 | 更新项目结构为模块化设计，调整核心组件实现，优化配置加载机制，更新Docker部署方案，完善配置文档和命令说明 | 开发团队 |
| v1.1.0 | 2024-03-15 | 更新配置系统为模块化设计，调整环境变量命名规则，更新配置文件结构 | 开发团队 |
| v1.0.0 | 2023-12-01 | 初始版本 | 开发团队 |

## 2. 项目概述

RBAC Admin Server是一个基于Go语言开发的权限管理系统后端服务，实现了完整的基于角色的访问控制（RBAC）功能。

### 2.1 核心功能
- 用户管理：用户的增删改查、密码修改、状态管理
- 角色管理：角色的增删改查、权限分配
- 菜单管理：动态菜单配置、权限控制
- 部门管理：组织架构维护
- 文件上传：支持文件上传和管理
- 系统日志：操作日志、登录日志记录
- JWT认证：基于JWT的身份验证机制

### 2.2 技术栈
- 后端框架：Go语言原生实现
- 数据库：支持MySQL、PostgreSQL、SQLite
- 缓存：Redis
- 权限控制：Casbin
- 配置管理：YAML配置文件 + 环境变量
- 日志系统：logrus

## 3. 项目结构

### 3.1 核心目录结构

```
rbac_admin_server/
├── api/                # API接口定义
│   ├── dept_api/       # 部门相关API
│   ├── file_api/       # 文件相关API
│   ├── log_api/        # 日志相关API
│   ├── menu_api/       # 菜单相关API
│   ├── permission_api/ # 权限相关API
│   ├── profile_api/    # 用户资料相关API
│   ├── role_api/       # 角色相关API
│   └── user_api/       # 用户相关API
├── config/             # 配置结构体定义
├── core/               # 核心组件初始化
│   ├── init_casbin/    # Casbin权限初始化
│   ├── init_gorm/      # 数据库初始化
│   └── init_redis/     # Redis初始化
├── global/             # 全局变量定义
├── middleware/         # 中间件
├── models/             # 数据模型
├── routes/             # 路由定义
├── utils/              # 工具函数
├── main.go             # 程序入口
├── go.mod              # Go模块定义
├── go.sum              # 依赖版本锁定
├── settings.yaml       # 主配置文件
├── settings.yaml.example # 配置文件示例
├── settings_dev.yaml   # 开发环境配置
└── settings_prod.yaml  # 生产环境配置
```

### 3.2 核心文件说明

| 文件名 | 功能描述 |
|-------|---------|
| main.go | 程序入口，负责加载配置、初始化组件、启动服务 |
| config/config.go | 定义全局配置结构体和默认配置 |
| core/init.go | 系统核心组件初始化逻辑 |
| core/read_config.go | 配置文件读取功能 |
| routes/routes.go | HTTP路由定义和注册 |
| settings.yaml | 主配置文件，包含系统所有配置项 |

## 4. 核心组件说明

### 4.1 配置系统

配置系统负责加载和解析YAML格式的配置文件，支持环境变量覆盖配置值。配置文件通过命令行参数指定，默认加载`settings.yaml`。

```go
// 配置加载流程
var configFile = flag.String("settings", "settings.yaml", "配置文件路径")
// 在main函数中加载配置
global.Config = core.ReadConfig(*configFile)
```

### 4.2 数据库组件

支持多种数据库类型，通过GORM框架实现数据库操作。数据库连接在系统初始化阶段建立，并在应用退出时自动关闭。

```go
// 数据库初始化
func InitGorm() (*gorm.DB, error) {
    // 根据配置创建数据库连接
}
```

### 4.3 Redis缓存组件

提供高性能的缓存服务，用于存储会话数据、热点数据等。Redis连接在系统初始化阶段建立，支持连接池管理。

```go
// Redis初始化
func InitRedis() (*redis.Client, error) {
    // 根据配置创建Redis连接
}
```

### 4.4 Casbin权限控制

基于Casbin实现的RBAC权限控制，支持角色和权限的灵活配置。

```go
// Casbin初始化
func InitCasbin() (*casbin.Enforcer, error) {
    // 初始化Casbin权限管理
}
```

### 4.5 日志系统

基于logrus的日志系统，支持多种日志级别、文件滚动、JSON格式输出等功能。

```go
// 日志系统初始化
func InitLogger(logConfig *config.LogConfig) error {
    // 配置和初始化日志系统
}
```

## 5. 快速开始

### 5.1 环境准备

- Go 1.16+ 环境
- MySQL/PostgreSQL/SQLite 数据库
- Redis 服务（可选）

### 5.2 克隆代码

```bash
# 克隆项目代码
git clone https://github.com/your-username/rbac_admin_server.git
cd rbac_admin_server
```

### 5.3 配置环境

1. 复制配置文件示例

```bash
# Linux/Mac
cp settings.yaml.example settings.yaml

# Windows
copy settings.yaml.example settings.yaml
```

2. 根据实际环境修改配置文件 `settings.yaml`

```yaml
# 配置数据库连接
db:
  mode: "mysql"  # 可选: mysql, postgres, sqlite
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_database_password"
  dbname: "rbac_admin"

# 配置Redis连接（可选）
redis:
  addr: "localhost:6379"
  password: "your_redis_password"
  db: 0

# 配置JWT认证
jwt:
  secret: "your_jwt_secret_key"  # 至少32个字符
  expire_hours: 24
```

3. 设置环境变量（可选，用于覆盖配置文件中的敏感信息）

```bash
# Linux/Mac
export DB_PASSWORD=your_database_password
export REDIS_PASSWORD=your_redis_password
export JWT_SECRET=your_jwt_secret_key

# Windows PowerShell
$env:DB_PASSWORD="your_database_password"
$env:REDIS_PASSWORD="your_redis_password"
$env:JWT_SECRET="your_jwt_secret_key"
```

### 5.4 启动服务器

```bash
# 直接运行
# Linux/Mac
go run main.go

# Windows PowerShell
go run main.go

# 指定配置文件运行
go run main.go -settings settings_dev.yaml
```

### 5.5 验证服务

服务启动后，可以通过以下方式验证：

```bash
# 检查服务是否启动成功
curl http://localhost:8080/health
```

## 6. 开发指南

### 6.1 代码规范
- 遵循Go语言标准编码规范
- 函数和变量命名清晰明了
- 关键代码添加适当注释
- 保持代码简洁，避免冗余

### 6.2 开发流程
1. 确保本地环境配置正确
2. 运行 `go mod tidy` 安装依赖
3. 启动开发服务器 `go run main.go -settings settings_dev.yaml`
4. 开发和测试新功能
5. 提交代码前运行 `go fmt ./...` 和 `go vet ./...`

### 6.3 调试技巧
- 使用 `go run -race main.go` 检测竞态条件
- 使用日志输出调试信息
- 对于复杂问题，可使用Delve调试器

## 7. 配置详解

### 7.1 配置文件结构

配置文件采用YAML格式，主要包含以下几个部分：

- system: 系统基础配置
- db: 数据库配置
- redis: Redis缓存配置
- jwt: JWT认证配置
- log: 日志系统配置
- security: 安全配置
- cors: 跨域资源共享配置
- swagger: API文档配置
- monitoring: 监控配置
- performance: 性能配置
- captcha: 验证码配置

### 7.2 核心配置项说明

#### 7.2.1 系统配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| mode | "dev" | 运行模式，可选值：dev, test, prod |
| ip | "0.0.0.0" | 服务器绑定IP地址 |
| port | 8080 | 服务器监听端口 |
| read_timeout | 30 | 读取超时时间（秒） |
| write_timeout | 30 | 写入超时时间（秒） |
| max_header_bytes | 1048576 | 最大请求头大小（字节） |

#### 7.2.2 数据库配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| mode | "sqlite" | 数据库类型，可选值：mysql, postgres, sqlite |
| host | "localhost" | 数据库主机地址 |
| port | 3306 | 数据库端口 |
| user | "root" | 数据库用户名 |
| password | "123456" | 数据库密码 |
| dbname | "rbac_admin.db" | 数据库名称 |
| path | "./data" | SQLite数据库文件路径 |
| max_idle_conns | 10 | 最大空闲连接数 |
| max_open_conns | 100 | 最大打开连接数 |
| conn_max_lifetime | 60 | 连接最大生命周期（分钟） |

#### 7.2.3 Redis配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| addr | "localhost:6379" | Redis服务器地址和端口 |
| password | "" | Redis密码 |
| db | 0 | Redis数据库索引 |
| pool_size | 100 | Redis连接池大小 |
| min_idle_conns | 10 | 最小空闲连接数 |
| dial_timeout | 5 | 拨号超时时间（秒） |
| read_timeout | 3 | 读取超时时间（秒） |
| write_timeout | 3 | 写入超时时间（秒） |
| idle_timeout | 60 | 空闲连接超时时间（分钟） |

#### 7.2.4 JWT配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| secret | "your-secret-key" | JWT签名密钥，至少32个字符 |
| expire_hours | 24 | JWT令牌有效期（小时） |
| issuer | "rbac-admin" | 令牌颁发者 |
| renew_window | 6 | 令牌自动续期窗口（小时） |

#### 7.2.5 日志配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| level | "info" | 日志级别，可选值：debug, info, warn, error, fatal, panic |
| dir | "./logs" | 日志文件存储目录 |
| filename | "app.log" | 日志文件名 |
| format | "text" | 日志格式，可选值：text, json |
| max_size | 50 | 单个日志文件最大大小（MB） |
| max_age | 7 | 日志文件保留天数 |
| max_backups | 10 | 保留的最大日志文件数 |
| compress | false | 是否压缩旧日志文件 |
| stdout | true | 是否输出到标准输出 |
| enable_caller | true | 是否记录调用者信息 |

#### 7.2.6 监控配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| enabled | true | 是否启用监控功能 |
| prometheus_port | 9090 | Prometheus监控端口 |
| health_check_path | "/health" | 健康检查API路径 |
| metrics_path | "/metrics" | Prometheus指标收集路径 |
| trace_sampling_rate | 0.1 | 分布式追踪采样率（0-1之间） |

#### 7.2.7 验证码配置

| 配置项 | 默认值 | 说明 |
|-------|-------|------|
| width | 120 | 验证码图片宽度 |
| height | 40 | 验证码图片高度 |
| length | 4 | 验证码字符长度 |
| expire_seconds | 300 | 验证码有效期（秒） |

### 7.3 多环境配置策略

项目支持多环境配置，可以通过不同的配置文件区分开发、测试和生产环境：

1. `settings_dev.yaml`: 开发环境配置
2. `settings_test.yaml`: 测试环境配置
3. `settings_prod.yaml`: 生产环境配置

启动时通过 `-settings` 参数指定要使用的配置文件：

```bash
# 开发环境
go run main.go -settings settings_dev.yaml

# 生产环境
go run main.go -settings settings_prod.yaml
```

### 7.4 环境变量支持

配置文件中的敏感信息（如密码、密钥等）可以通过环境变量注入，提高安全性：

```yaml
# 在配置文件中使用环境变量
password: "${DB_PASSWORD}"
```

设置环境变量：

```bash
# Linux/Mac
export DB_PASSWORD=your_database_password

# Windows PowerShell
$env:DB_PASSWORD="your_database_password"
```

### 7.5 配置优先级规则

配置项的优先级从高到低依次为：
1. 命令行参数
2. 环境变量
3. 指定的配置文件
4. 默认配置

### 7.6 生产环境配置最佳实践

1. 使用独立的配置文件（如`settings_prod.yaml`）
2. 敏感信息通过环境变量注入
3. 关闭调试日志，设置日志级别为`info`或`warn`
4. 配置合理的数据库连接池参数
5. 启用Redis缓存以提高性能
6. 使用强JWT密钥（至少32个字符）
7. 配置合理的令牌过期时间

## 8. 部署指南

### 8.1 验证服务

在部署前，确保本地开发环境中的服务能够正常运行和访问：

```bash
# 启动服务
go run main.go

# 验证服务
curl http://localhost:8080/health
```

### 8.2 测试环境部署

测试环境部署可以采用与开发环境类似的方式，但建议使用独立的配置文件：

```bash
# 上传代码到测试服务器
scp -r rbac_admin_server user@test-server:/path/to/deploy

# 登录测试服务器
ssh user@test-server

# 切换到部署目录
cd /path/to/deploy/rbac_admin_server

# 安装依赖
 go mod tidy

# 编译应用
go build -o rbac_admin_server main.go

# 启动服务（使用测试环境配置）
./rbac_admin_server -settings settings_test.yaml
```

### 8.3 生产环境部署

#### 8.3.1 系统准备

1. 安装必要的依赖（Go环境、数据库、Redis等）
2. 创建专用的系统用户运行服务
3. 配置防火墙规则，只开放必要的端口
4. 配置系统时间同步

#### 8.3.2 数据库准备

1. 创建数据库和用户

```sql
-- MySQL
CREATE DATABASE rbac_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'rbac_user'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON rbac_admin.* TO 'rbac_user'@'localhost';
FLUSH PRIVILEGES;
```

2. 导入初始数据（如果有）

#### 8.3.3 环境配置

1. 创建配置文件 `settings_prod.yaml`，根据生产环境进行配置
2. 设置环境变量存储敏感信息
3. 配置系统服务管理（如Systemd）

#### 8.3.4 部署目录创建

```bash
# 创建部署目录
mkdir -p /opt/rbac_admin_server
mkdir -p /etc/rbac_admin_server
mkdir -p /var/log/rbac_admin_server

# 创建系统用户
useradd -r -s /sbin/nologin rbac_admin

# 设置目录权限
chown -R rbac_admin:rbac_admin /opt/rbac_admin_server
chown -R rbac_admin:rbac_admin /var/log/rbac_admin_server
chown -R root:rbac_admin /etc/rbac_admin_server
chmod 750 /etc/rbac_admin_server
```

#### 8.3.5 代码获取与部署

```bash
# 克隆代码
cd /opt
 git clone https://github.com/your-username/rbac_admin_server.git

# 编译应用
cd /opt/rbac_admin_server
go mod tidy
go build -o rbac_admin_server main.go

# 复制配置文件
cp settings.yaml.example /etc/rbac_admin_server/settings_prod.yaml

# 修改生产环境配置
vi /etc/rbac_admin_server/settings_prod.yaml

# 设置文件权限
chown root:rbac_admin /etc/rbac_admin_server/settings_prod.yaml
chmod 640 /etc/rbac_admin_server/settings_prod.yaml
```

#### 8.3.6 Systemd服务配置

创建Systemd服务文件 `/etc/systemd/system/rbac_admin_server.service`：

```ini
[Unit]
Description=RBAC Admin Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=rbac_admin
Group=rbac_admin
WorkingDirectory=/opt/rbac_admin_server
ExecStart=/opt/rbac_admin_server/rbac_admin_server -settings /etc/rbac_admin_server/settings_prod.yaml
Restart=on-failure
RestartSec=5s
Environment="DB_PASSWORD=your_database_password"
Environment="REDIS_PASSWORD=your_redis_password"
Environment="JWT_SECRET=your_jwt_secret_key"

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
# 重新加载Systemd配置
systemctl daemon-reload

# 启动服务
systemctl start rbac_admin_server

# 设置开机自启
systemctl enable rbac_admin_server

# 查看服务状态
systemctl status rbac_admin_server
```

#### 8.3.7 Nginx反向代理配置

配置Nginx作为反向代理，提供负载均衡和HTTPS支持：

```nginx
server {
    listen 80;
    server_name rbac.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name rbac.example.com;

    # SSL配置
    ssl_certificate /etc/letsencrypt/live/rbac.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/rbac.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_session_tickets off;

    # 代理设置
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 请求限流
        limit_req zone=rbac_admin burst=10 nodelay;

        # 缓冲区设置
        proxy_buffering on;
        proxy_buffer_size 8k;
        proxy_buffers 8 8k;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        proxy_pass http://localhost:8080;
        expires 30d;
        add_header Cache-Control "public, max-age=2592000";
    }

    # 文件上传目录
    location /uploads/ {
        proxy_pass http://localhost:8080/uploads/;
        expires 1d;
        add_header Cache-Control "public, max-age=86400";
    }
}

# 限流配置
limit_req_zone $binary_remote_addr zone=rbac_admin:10m rate=20r/s;
```

启用Nginx配置并重启：

```bash
# 检查配置语法
nginx -t

# 重启Nginx
systemctl restart nginx

# 设置Nginx开机自启
systemctl enable nginx
```

### 8.4 验证生产环境部署

#### 8.4.1 服务状态检查

```bash
# 查看Systemd服务状态
systemctl status rbac_admin_server

# 检查监听端口
netstat -tuln | grep 8080
```

#### 8.4.2 日志查看

```bash
# 查看Systemd服务日志
journalctl -u rbac_admin_server -f

# 查看应用日志
cat /var/log/rbac_admin_server/app.log
```

#### 8.4.3 API测试

使用curl或其他工具测试API端点：

```bash
# 测试公共登录接口
curl -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"your_password"}' https://rbac.example.com/api/v1/system/auth/login
```

#### 8.4.4 HTTPS配置验证

使用浏览器访问 https://rbac.example.com/api/v1/system/auth/login 确认HTTPS配置正确

### 8.5 配置自动更新

创建配置更新脚本`/opt/rbac_admin_server/scripts/update_config.sh`：

```bash
#!/bin/bash
# 配置更新脚本

# 备份当前配置
cp /etc/rbac_admin_server/settings_prod.yaml /etc/rbac_admin_server/settings_prod.yaml.bak

# 从配置中心获取最新配置
# 这里以示例方式展示，实际应根据你的配置中心实现
wget -O /etc/rbac_admin_server/settings_prod.yaml.new https://config.example.com/rbac_admin_server/prod/settings_prod.yaml

# 验证新配置格式
if grep -q "db:" /etc/rbac_admin_server/settings_prod.yaml.new && grep -q "jwt:" /etc/rbac_admin_server/settings_prod.yaml.new;
then
    # 应用新配置
    mv /etc/rbac_admin_server/settings_prod.yaml.new /etc/rbac_admin_server/settings_prod.yaml
    chown root:rbac_admin /etc/rbac_admin_server/settings_prod.yaml
    chmod 640 /etc/rbac_admin_server/settings_prod.yaml
    
    # 重启服务应用新配置
systemctl restart rbac_admin_server
    echo "Configuration updated successfully at $(date)"
else
    echo "Invalid configuration format, update aborted at $(date)"
    rm -f /etc/rbac_admin_server/settings_prod.yaml.new
fi
```

设置定时任务自动检查配置更新：

```bash
# 设置执行权限
chmod +x /opt/rbac_admin_server/scripts/update_config.sh

# 创建定时任务（每天凌晨2点执行）
crontab -u rbac_admin -e
# 添加以下内容
0 2 * * * /opt/rbac_admin_server/scripts/update_config.sh >> /var/log/rbac_admin_server/config_update.log 2>&1
```

## 9. 部署模式

### 9.1 传统服务器集群部署

#### 9.1.1 架构

采用多服务器部署，通过负载均衡器分发请求，共享数据库和Redis服务。

#### 9.1.2 适用场景
- 对可用性要求较高的中大型应用
- 需要灵活扩展计算资源的场景
- 已有机房或物理服务器资源

#### 9.1.3 最佳实践
- 使用Nginx或HAProxy作为负载均衡器
- 配置会话共享，避免会话粘性问题
- 实现数据库读写分离，提高并发性能
- 配置Redis哨兵或集群模式，提高缓存服务可用性
- 实施滚动更新策略，减少服务中断时间

### 9.2 容器化部署（Docker）

#### 9.2.1 特点
- 环境一致性，消除"开发环境正常，生产环境异常"问题
- 快速部署和扩缩容
- 隔离性好，各服务互不影响
- 便于自动化部署和CI/CD集成

#### 9.2.2 架构

使用Docker Compose编排多个容器服务，包括应用容器、数据库容器和Redis容器。

#### 9.2.3 示例Dockerfile

```dockerfile
# 使用官方Go镜像作为构建环境
FROM golang:1.18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rbac_admin_server main.go

# 使用Alpine作为运行环境
FROM alpine:3.16

# 设置工作目录
WORKDIR /app

# 复制编译好的二进制文件
COPY --from=builder /app/rbac_admin_server .

# 复制配置文件模板
COPY settings.yaml.example .

# 设置时区
RUN apk --no-cache add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 创建日志目录
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./rbac_admin_server", "-settings", "settings_prod.yaml"]
```

#### 9.2.4 Docker Compose示例

```yaml
version: '3'
services:
  app:
    build: .
    container_name: rbac_admin_server
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=3306
      - DB_NAME=rbac_admin
      - DB_USER=root
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=redis
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - db
      - redis
    volumes:
      - ./logs:/app/logs
      - ./settings_prod.yaml:/app/settings_prod.yaml
    restart: always

  db:
    image: mysql:8.0
    container_name: rbac_admin_db
    environment:
      - MYSQL_ROOT_PASSWORD=${DB_PASSWORD}
      - MYSQL_DATABASE=rbac_admin
    volumes:
      - db_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: always

  redis:
    image: redis:6.0
    container_name: rbac_admin_redis
    volumes:
      - redis_data:/data
    restart: always

volumes:
  db_data:
  redis_data:
```

#### 9.2.5 部署命令示例

```bash
# 设置环境变量
export DB_PASSWORD=your_secure_db_password
export JWT_SECRET=your_secure_jwt_secret

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看应用日志
docker logs -f rbac_admin_server

# 停止服务
docker-compose down
```

#### 9.2.6 适用场景
- 快速开发和测试环境搭建
- 微服务架构应用
- 需要自动化部署流程的场景
- 云原生应用部署

#### 9.2.7 最佳实践
- 使用多阶段构建减小镜像体积
- 避免在容器内存储敏感信息，使用环境变量或配置中心
- 合理设置容器资源限制（CPU、内存）
- 使用Docker卷或外部存储持久化数据
- 实现健康检查和自动重启机制

## 10. 监控与维护

### 10.1 健康检查机制

系统提供了健康检查功能，可以通过以下方式验证服务状态：

```bash
# 检查服务状态
systemctl status rbac_admin_server

# 查看应用日志
journalctl -u rbac_admin_server -f
```

### 10.2 日志监控

应用日志默认存储在 `./logs` 目录（生产环境为 `/var/log/rbac_admin_server`），可以使用日志分析工具进行监控和分析。

### 10.3 常见问题排查指南

#### 10.3.1 服务启动问题

1. **JWT密钥长度不足32位**
   - **问题**：启动时提示"JWT secret must be at least 32 characters"错误
   - **原因**：JWT签名密钥长度不符合安全要求
   - **解决步骤**：
     1. 生成至少32位的随机密钥：
        - Linux/Mac: `JWT_SECRET=$(openssl rand -hex 16)`
        - Windows PowerShell: `$JWT_SECRET = [Convert]::ToBase64String([System.Security.Cryptography.RandomNumberGenerator]::GetBytes(32))`
     2. 通过环境变量设置密钥：`export JWT_SECRET=your_generated_secret`
     3. 或在配置文件中设置：`jwt.secret: your_generated_secret`

2. **端口被占用**
   - **问题**：启动时提示"listen tcp :8080: bind: address already in use"错误
   - **原因**：指定端口已被其他进程占用
   - **解决步骤**：
     1. 查找占用端口的进程：
        - Linux: `netstat -tuln | grep 8080` 或 `lsof -i :8080`
        - Windows: `netstat -ano | findstr :8080`
     2. 根据进程ID（PID）停止占用端口的进程：
        - Linux: `kill -9 <PID>`
        - Windows: `taskkill /F /PID <PID>`
     3. 或修改服务监听端口：通过环境变量`PORT=8081`或配置文件修改

3. **配置文件格式错误**
   - **问题**：启动时提示"yaml: line X: did not find expected key"等解析错误
   - **原因**：YAML配置文件格式不正确，通常是缩进、语法或字符编码问题
   - **解决步骤**：
     1. 检查配置文件的缩进是否正确（YAML使用空格，不支持Tab）
     2. 验证配置项名称和值的语法是否正确
     3. 检查是否有特殊字符或非ASCII字符导致编码问题
     4. 使用YAML验证工具（如yamllint.com）检查文件格式

#### 10.3.2 数据库连接问题

1. **数据库连接超时**
   - **问题**：启动时或运行中提示"dial tcp <host>:<port>: connect: connection timed out"错误
   - **原因**：无法连接到数据库服务器，可能是网络问题或数据库服务未启动
   - **解决步骤**：
     1. 检查数据库服务是否正常运行：
        - MySQL: `systemctl status mysql` 或 `service mysql status`
        - PostgreSQL: `systemctl status postgresql`
     2. 验证数据库连接参数（主机、端口、用户名、密码、数据库名）是否正确
     3. 检查网络连通性：`ping <db_host>` 和 `telnet <db_host> <db_port>`
     4. 确认数据库服务器防火墙是否允许连接

2. **数据库访问权限不足**
   - **问题**：连接数据库后提示"ERROR 1044 (42000): Access denied for user '<user>'@'<host>' to database '<db>'"错误
   - **原因**：数据库用户没有足够的权限访问指定数据库
   - **解决步骤**：
     1. 以管理员身份登录数据库：`mysql -u root -p`
     2. 授予用户访问权限：
        ```sql
        GRANT ALL PRIVILEGES ON <db_name>.* TO '<user>'@'%' IDENTIFIED BY '<password>';
        FLUSH PRIVILEGES;
        ```

## 11. 附录

### 11.1 配置项速查表

下表列出了RBAC管理员服务器的核心配置项及其说明，您可以根据实际需求进行配置调整。

| 配置模块 | 配置项 | 默认值 | 说明 | 配置文件路径 |
|---------|-------|-------|------|------------|
| **system** | mode | "dev" | 运行模式，可选值：dev, test, prod | config/system.go |
|  | ip | "0.0.0.0" | 服务器绑定IP地址 | config/system.go |
|  | port | 8080 | 服务器监听端口 | config/system.go |
|  | read_timeout | 30 | HTTP读取超时（秒） | config/system.go |
|  | write_timeout | 30 | HTTP写入超时（秒） | config/system.go |
|  | max_header_bytes | 1048576 | 最大请求头大小（字节） | config/system.go |
| **db** | mode | "sqlite" | 数据库类型，可选值：mysql, postgres, sqlite | config/db.go |
|  | host | "localhost" | 数据库主机地址 | config/db.go |
|  | port | 3306 | 数据库端口 | config/db.go |
|  | user | "root" | 数据库用户名 | config/db.go |
|  | password | "123456" | 数据库密码 | config/db.go |
|  | dbname | "rbac_admin.db" | 数据库名称 | config/db.go |
|  | path | "./data" | SQLite数据库文件路径 | config/db.go |
|  | max_idle_conns | 10 | 最大空闲连接数 | config/db.go |
|  | max_open_conns | 100 | 最大打开连接数 | config/db.go |
|  | conn_max_lifetime | 60 | 连接最大生命周期（分钟） | config/db.go |
| **redis** | addr | "localhost:6379" | Redis服务器地址和端口 | config/redis.go |
|  | password | "" | Redis密码 | config/redis.go |
|  | db | 0 | Redis数据库索引 | config/redis.go |
|  | pool_size | 100 | Redis连接池大小 | config/redis.go |
|  | min_idle_conns | 10 | 最小空闲连接数 | config/redis.go |
|  | dial_timeout | 5 | 拨号超时时间（秒） | config/redis.go |
|  | read_timeout | 3 | 读取超时时间（秒） | config/redis.go |
|  | write_timeout | 3 | 写入超时时间（秒） | config/redis.go |
|  | idle_timeout | 60 | 空闲连接超时时间（分钟） | config/redis.go |
| **jwt** | secret | "your-secret-key" | JWT签名密钥，至少32个字符 | config/jwt.go |
|  | expire_hours | 24 | JWT令牌有效期（小时） | config/jwt.go |
|  | issuer | "rbac-admin" | 令牌颁发者 | config/jwt.go |
|  | renew_window | 6 | 令牌自动续期窗口（小时） | config/jwt.go |
| **log** | level | "info" | 日志级别 | config/log.go |
|  | dir | "./logs" | 日志文件存储目录 | config/log.go |
|  | filename | "app.log" | 日志文件名 | config/log.go |
|  | format | "text" | 日志格式 | config/log.go |
|  | max_size | 50 | 单个日志文件最大大小（MB） | config/log.go |
|  | max_age | 7 | 日志文件保留天数 | config/log.go |
|  | max_backups | 10 | 保留的最大日志文件数 | config/log.go |
|  | compress | false | 是否压缩旧日志文件 | config/log.go |
|  | stdout | true | 是否输出到标准输出 | config/log.go |
|  | enable_caller | true | 是否记录调用者信息 | config/log.go |

### 11.2 开发与部署命令速查表

下表汇总了RBAC管理员服务器的常用开发和部署命令，方便快速查找和使用。

#### 11.2.1 开发环境命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `go mod tidy` | 安装依赖并清理未使用的依赖 | `go mod tidy` |
| `go run main.go` | 运行应用（开发模式） | `go run main.go` |
| `go run main.go -settings settings_dev.yaml` | 指定配置文件运行应用 | `go run main.go -settings settings_dev.yaml` |
| `go build -o rbac_admin_server main.go` | 构建应用 | `go build -o rbac_admin_server main.go` |
| `go test ./...` | 运行所有测试 | `go test ./...` |
| `go fmt ./...` | 格式化代码 | `go fmt ./...` |
| `go vet ./...` | 静态代码分析 | `go vet ./...` |

#### 11.2.2 生产环境命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `systemctl status rbac_admin_server` | 查看服务状态 | `systemctl status rbac_admin_server` |
| `systemctl start rbac_admin_server` | 启动服务 | `systemctl start rbac_admin_server` |
| `systemctl stop rbac_admin_server` | 停止服务 | `systemctl stop rbac_admin_server` |
| `systemctl restart rbac_admin_server` | 重启服务 | `systemctl restart rbac_admin_server` |
| `systemctl enable rbac_admin_server` | 设置开机自启 | `systemctl enable rbac_admin_server` |
| `journalctl -u rbac_admin_server -f` | 查看服务日志（实时） | `journalctl -u rbac_admin_server -f` |
| `tail -f /var/log/rbac_admin_server/app.log` | 查看应用日志（实时） | `tail -f /var/log/rbac_admin_server/app.log` |

#### 11.2.3 Docker相关命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `docker build -t rbac_admin_server .` | 构建Docker镜像 | `docker build -t rbac_admin_server .` |
| `docker run -p 8080:8080 rbac_admin_server` | 运行Docker容器 | `docker run -p 8080:8080 rbac_admin_server` |
| `docker-compose up -d` | 启动所有服务（Docker Compose） | `docker-compose up -d` |
| `docker-compose down` | 停止所有服务（Docker Compose） | `docker-compose down` |
| `docker-compose logs -f` | 查看所有服务日志（实时） | `docker-compose logs -f` |
| `docker-compose ps` | 查看服务状态（Docker Compose） | `docker-compose ps` |