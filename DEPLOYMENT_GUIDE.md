# RBAC Admin Server 部署指南

## 快速开始

如果你想快速部署和运行RBAC Admin Server，可以按照以下步骤操作：

### 使用Docker（推荐）

```bash
# 克隆项目
 git clone https://github.com/rbacadmin/rbac_admin_server.git
 cd rbac_admin_server

# 创建配置文件
 cp settings.yaml.example settings.yaml
 cp .env.example .env

# 根据实际情况修改配置文件
# vi settings.yaml
# vi .env

# 使用docker-compose启动所有服务
 docker-compose up -d --build

# 验证服务是否正常运行
 curl http://localhost:8080/health
```

### 传统部署

```bash
# 克隆项目
 git clone https://github.com/rbacadmin/rbac_admin_server.git
 cd rbac_admin_server

# 创建配置文件
 cp settings.yaml.example settings.yaml
 cp .env.example .env

# 根据实际情况修改配置文件
# vi settings.yaml
# vi .env

# 安装依赖
 go mod download

# 编译项目
 go build -o rbac_admin_server

# 启动服务
 ./rbac_admin_server
```

## 目录

1. [项目简介](#1-项目简介)
2. [技术栈](#2-技术栈)
3. [环境准备](#3-环境准备)
4. [项目获取](#4-项目获取)
5. [配置文件设置](#5-配置文件设置)
6. [环境变量配置](#6-环境变量配置)
7. [配置文件使用](#7-配置文件使用)
8. [编译和运行](#8-编译和运行)
9. [服务验证](#9-服务验证)
10. [API接口测试](#10-api接口测试)
11. [项目目录结构](#11-项目目录结构)
12. [常见问题与解决方案](#12-常见问题与解决方案)
13. [安全建议](#13-安全建议)
14. [开发指南](#14-开发指南)
15. [部署指南](#15-部署指南)
    - [15.1 Docker部署（推荐）](#151-docker部署推荐)
    - [15.2 传统部署](#152-传统部署)
16. [维护与更新](#16-维护与更新)
    - [16.1 更新流程](#161-更新流程)
    - [16.2 日志管理](#162-日志管理)
    - [16.3 监控与维护](#163-监控与维护)
17. [联系支持](#17-联系支持)

## 1. 项目简介

RBAC Admin Server 是一个基于 Go 语言开发的 RBAC（基于角色的访问控制）系统后端服务，提供用户管理、角色管理、权限管理以及文件管理等功能。它采用前后端分离架构，可与 rbacAdmin 前端项目配合使用，为企业级应用提供完整的用户权限管理解决方案。

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
- **浏览器**：用于访问 Swagger API 文档

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

#### Linux 系统

```bash
# 下载 Go 安装包
wget https://golang.org/dl/go1.24.0.linux-amd64.tar.gz

# 解压安装包
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# 添加到环境变量
cat << EOF >> ~/.profile
export PATH=$PATH:/usr/local/go/bin
export GOPATH=~/go
export GOROOT=/usr/local/go
export PATH="$GOPATH/bin:$PATH"
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
EOF

# 应用环境变量
source ~/.profile

# 验证安装
go version
```

#### macOS 系统

```bash
# 使用 Homebrew 安装 Go
brew install go@1.24

# 添加到环境变量
cat << EOF >> ~/.zshrc
export PATH="/usr/local/opt/go@1.24/bin:$PATH"
export GOPATH=~/go
export GOROOT=/usr/local/opt/go@1.24
export PATH="$GOPATH/bin:$PATH"
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
EOF

# 应用环境变量
source ~/.zshrc

# 验证安装
go version
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

#### Linux 系统

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install redis-server -y

sudo systemctl enable redis-server
sudo systemctl start redis-server

# 验证安装
redis-cli ping
```

#### macOS 系统

```bash
# 使用 Homebrew 安装 Redis
brew install redis

# 启动 Redis 服务
brew services start redis

# 验证安装
redis-cli ping
```

### 3.3 数据库安装

#### MySQL 8.0 安装

**Windows 系统**

1. 访问 [MySQL 官方下载页面](https://dev.mysql.com/downloads/installer/)，下载 MySQL 8.0 安装包
2. 运行安装包并按照提示完成安装
3. 设置 root 密码和其他必要配置

**Linux 系统**

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server -y

sudo systemctl enable mysql
sudo systemctl start mysql

# 安全配置
mysql_secure_installation

# 登录 MySQL
mysql -u root -p
```

**macOS 系统**

```bash
# 使用 Homebrew 安装 MySQL
brew install mysql@8.0

# 启动 MySQL 服务
brew services start mysql@8.0

# 安全配置
mysql_secure_installation

# 登录 MySQL
mysql -u root -p
```

### 3.4 数据库初始化

安装完成后，需要创建一个数据库供 RBAC Admin Server 使用：

```sql
-- 登录 MySQL 后执行以下 SQL 语句
CREATE DATABASE IF NOT EXISTS rbacadmin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON rbacadmin.* TO 'root'@'%' IDENTIFIED BY 'Zdj_7819!';
FLUSH PRIVILEGES;
```

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

settings.yaml 是项目的主要配置文件，包含系统运行所需的各种参数。以下是关键配置项说明：

```yaml
# 🖥️ 服务器配置
system:
  ip: 127.0.0.1           # 服务IP地址
  port: 8090              # 服务端口（注意：需与前端项目配置保持一致）
  mode: "debug"           # 运行模式: debug, release（生产环境建议使用release）

# 🗄️ 数据库配置
db:
  mode: "mysql"            # 数据库类型: mysql, postgres, sqlite
  host: "127.0.0.1"       # 数据库主机
  port: 3306              # 数据库端口
  user: "root"             # 数据库用户名
  password: "admin123"     # 数据库密码（建议在.env中配置）
  dbname: "rbac_admin"     # 数据库名称
  max_open_conns: 100     # 最大连接数
  max_idle_conns: 10      # 空闲连接数
  conn_max_lifetime: 3600 # 连接生命周期(秒)

# 🔄 Redis配置
redis:
  addr: "${REDIS_ADDR:127.0.0.1:6379}"  # Redis地址（支持环境变量替换）
  password: "${REDIS_PASSWORD}"         # Redis密码
  db: ${REDIS_DB:3}                      # Redis数据库编号（建议与前端配置保持一致）
  pool_size: 20                          # 连接池大小
  min_idle_conns: 5                      # 最小空闲连接数

# 🔐 JWT认证配置
jwt:
  secret: "${JWT_SECRET:aB3kL9mN7xY2qR8sT1uV4wE6zC0pF5gH}"  # JWT密钥 (大小写字母+数字组合的强密钥)
  expire_hours: 72                            # Token过期时间(小时)
  refresh_expire_hours: 168                   # 刷新Token过期时间(小时)
  issuer: "rbacAdmin"                         # Token签发者
  audience: "rbac-client"                     # Token受众

# 🧩 验证码配置（重要）
captcha:
  enable: true          # 是否启用验证码（默认启用，登录接口需要验证码）
  width: 120            # 验证码图片宽度
  height: 40            # 验证码图片高度
  length: 4             # 验证码长度
  expire_seconds: 300   # 验证码有效期(秒)

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

# 🔒 安全配置
security:
  xss_protection: "1"    # XSS保护
  content_type_nosniff: "nosniff" # 内容类型嗅探
  x_frame_options: "DENY" # X-Frame-Options
  csrf_protection: true  # CSRF保护
  rate_limit: 100        # 速率限制（请求/秒）
  bcrypt_cost: 12        # BCrypt加密成本

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

# 📊 监控配置
monitoring:
  enabled: true          # 是否启用监控
  prometheus_port: 9090  # Prometheus端口
  health_check_path: "/health" # 健康检查路径
  metrics_path: "/metrics" # 指标路径
  trace_sampling_rate: 0.1 # 跟踪采样率

# 📚 Swagger配置
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

# 🚀 应用配置
app:
  name: "RBAC Admin Server" # 应用名称
  version: "1.0.0"       # 应用版本
  description: "A RBAC Admin Server implemented in Go" # 应用描述
  copyright: "© 2023 RBAC Admin Server" # 版权信息
  timezone: "Asia/Shanghai" # 时区
  language: "zh-CN"      # 语言
  debug: true            # 调试模式

# 📤 上传配置
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

.env文件用于存储敏感信息和环境特定配置，以下是一个生产环境配置示例：

```env
# =================================================================================
# 🏗️ RBAC管理员服务器 - 生产环境配置
# =================================================================================
# 📋 使用说明：
#   此文件包含用户提供的具体环境配置信息
#   基于 .env.example 模板，使用实际的服务器地址和凭据
# =================================================================================

# 🖥️ 系统配置
SYSTEM_PORT=8080

# 🗄️ 数据库配置 - 使用用户提供的具体信息
DB_MODE=mysql
DB_HOST=192.168.10.199
DB_PORT=3306
DB_USER=root
DB_PASSWORD=Zdj_7819!
DB_DBNAME=rbacadmin
DB_PATH=./data/rbac_admin.db

# 🔐 JWT配置 - 大小写字母和数字组合的强密钥
JWT_SECRET=AbCdEfGhIjKlMnOpQrStUvWxYz1234567890
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_HOURS=168
JWT_ISSUER=rbac-admin
JWT_AUDIENCE=rbac-admin

# 🔄 Redis配置 - 使用用户提供的具体信息
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

# 🔒 安全配置
CSRF_SECRET=your_csrf_secret_key

# 🌐 CORS配置
CORS_ORIGINS=https://your-domain.com
```

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

#### 5.3.4 配置验证

系统启动时会自动验证配置的有效性，如遇到配置错误会在日志中显示详细信息。

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

### 8.1 测试说明

系统API接口需要进行身份认证，大部分接口需要在请求头中提供有效的JWT令牌。以下是测试API接口的基本流程：

1. **获取验证码**：调用`/public/captcha/get`接口获取验证码ID和验证码内容
2. **登录**：使用获取的验证码和用户名密码调用`/public/login`接口获取JWT令牌
3. **调用受保护接口**：在请求头中添加`Authorization: Bearer {token}`来调用需要认证的接口

### 8.2 验证码功能说明（重要）

系统默认启用了验证码功能，在调用登录接口前必须先获取验证码：

- 验证码接口：`GET /public/captcha/get`
- 返回内容包含：`captchaId`（验证码ID）和`image`（验证码图片的Base64编码）
- 验证码有效期：默认300秒（可在settings.yaml中配置）
- 登录请求必须包含：`captchaId`和`captchaCode`字段

### 8.3 测试脚本

项目提供了多个PowerShell脚本用于测试API接口：

- `test_login_detailed.ps1` - 详细测试登录功能
- `test_admin_user_list.ps1` - 测试获取用户列表
- `test_admin_crud_operations.ps1` - 测试管理员CRUD操作
- `test_all_admin_apis.ps1` - 测试所有管理员API

运行测试脚本示例：
```powershell
# 运行登录测试
.\test_login_detailed.ps1
```

### 8.4 Go代码测试示例

以下是使用Go代码测试API接口的完整示例，包括获取验证码、登录和调用管理员接口的流程：

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CaptchaResponse 验证码响应结构
type CaptchaResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CaptchaID string `json:"captchaId"`
		Image     string `json:"image"`
		Answer    string `json:"answer"` // 注意：实际环境中不会返回Answer字段
	} `json:"data"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaID   string `json:"captchaId"`
	CaptchaCode string `json:"captchaCode"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// UserListResponse 用户列表响应结构
type UserListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Status    int    `json:"status"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
}

func main() {
	// 1. 获取验证码
	captchaResp, err := getCaptcha()
	if err != nil {
		fmt.Printf("获取验证码失败: %v\n", err)
		return
	}
	fmt.Printf("获取验证码成功: CaptchaID=%s, Answer=%s\n", captchaResp.Data.CaptchaID, captchaResp.Data.Answer)

	// 2. 登录获取Token
	token, err := login(captchaResp.Data.CaptchaID, captchaResp.Data.Answer)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	fmt.Printf("登录成功，Token=%s\n", token)

	// 3. 使用Token调用管理员接口
	users, err := getUserList(token)
	if err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
		return
	}
	fmt.Printf("获取用户列表成功，共%d个用户\n", len(users))
	for _, user := range users {
		fmt.Printf("用户: ID=%d, Username=%s, Nickname=%s\n", user.ID, user.Username, user.Nickname)
	}
}

// 获取验证码
func getCaptcha() (*CaptchaResponse, error) {
	resp, err := http.Get("http://localhost:8080/public/captcha/get")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var captchaResp CaptchaResponse
	if err := json.Unmarshal(body, &captchaResp); err != nil {
		return nil, err
	}

	return &captchaResp, nil
}

// 登录获取Token
func login(captchaID, captchaCode string) (string, error) {
	loginData := LoginRequest{
		Username:    "admin",
		Password:    "admin123",
		CaptchaID:   captchaID,
		CaptchaCode: captchaCode,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:8080/public/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", err
	}

	return loginResp.Data.Token, nil
}

// 获取用户列表
func getUserList(token string) ([]struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/admin/user/list", nil)
	if err != nil {
		return nil, err
	}

	// 添加Authorization头
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userListResp UserListResponse
	if err := json.Unmarshal(body, &userListResp); err != nil {
		return nil, err
	}

	return userListResp.Data, nil
}
```

创建一个名为`test_api.go`的文件，将上述代码保存进去，然后运行：

```bash
go run test_api.go
```

这个示例将演示完整的API调用流程，从获取验证码开始，然后登录获取JWT令牌，最后使用令牌调用需要认证的用户列表接口。

**注意：** 示例代码中存在一个`Answer`字段，仅用于演示目的。在实际环境中，验证码接口不会返回验证码的答案，客户端需要用户手动输入验证码。

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

## 15. 部署指南

### 15.1 Docker部署（推荐）

项目支持Docker部署，以下是完整的Docker部署指南：

#### 15.1.1 Dockerfile

创建一个名为`Dockerfile`的文件，内容如下：

```dockerfile
# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖工具
RUN apk add --no-cache git gcc musl-dev

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译项目
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rbac_admin_server -ldflags="-s -w"

# 运行阶段
FROM alpine:3.18

# 设置工作目录
WORKDIR /app

# 创建日志和上传目录
RUN mkdir -p ./logs ./uploads

# 复制二进制文件和配置文件
COPY --from=builder /app/rbac_admin_server .
COPY settings.yaml .
COPY .env .
COPY config/casbin/ ./config/casbin/

# 设置文件权限
RUN chmod +x ./rbac_admin_server

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV TZ=Asia/Shanghai

# 启动命令
CMD ["./rbac_admin_server"]
```

#### 15.1.2 docker-compose.yml

为了更方便地部署整个应用栈（包括数据库和Redis），可以创建一个`docker-compose.yml`文件：

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
      - redis
    environment:
      - TZ=Asia/Shanghai
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    restart: always

  db:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=Zdj_7819!
      - MYSQL_DATABASE=rbacadmin
      - TZ=Asia/Shanghai
    volumes:
      - mysql_data:/var/lib/mysql
    restart: always

  redis:
    image: redis:7.0
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: always

volumes:
  mysql_data:
  redis_data:
```

#### 15.1.3 构建和运行Docker容器

**使用Docker命令行：**

```bash
# 构建镜像
docker build -t rbac_admin_server .

# 运行容器（单独运行应用，需要外部数据库和Redis）
docker run -p 8080:8080 --env-file .env rbac_admin_server

# 运行容器（包含数据库和Redis）
docker-compose up -d
```

**使用docker-compose：**

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app

# 停止所有服务
docker-compose down
```

#### 15.1.4 生产环境Docker部署建议

1. **使用环境变量**：在生产环境中，避免在Dockerfile中硬编码敏感信息，使用环境变量或Docker secrets

2. **使用标签管理**：为Docker镜像添加版本标签，方便回滚和部署管理
   ```bash
docker build -t rbac_admin_server:v1.0.0 .
   ```

3. **限制资源使用**：在docker-compose.yml中添加资源限制
   ```yaml
   resources:
     limits:
       cpus: "0.5"
       memory: "512M"
   ```

4. **配置健康检查**：添加健康检查确保服务正常运行
   ```yaml
   healthcheck:
     test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
     interval: 30s
     timeout: 5s
     retries: 3
     start_period: 30s
   ```

### 15.2 传统部署

如果不使用Docker，也可以进行传统部署：

#### 15.2.1 环境安装

1. **安装Go环境**
   ```bash
   # 下载Go安装包
   wget https://golang.org/dl/go1.24.0.linux-amd64.tar.gz
   
   # 解压安装包
   tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
   
   # 配置环境变量（添加到/etc/profile或~/.bashrc）
   echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
   echo "export GOPATH=~/go" >> /etc/profile
   source /etc/profile
   
   # 验证安装
   go version
   ```

2. **安装MySQL**
   ```bash
   # Ubuntu/Debian
   apt-get update && apt-get install -y mysql-server-8.0
   
   # CentOS/RHEL
   yum install -y mysql-server-8.0
   
   # 启动MySQL服务
   systemctl start mysql
   systemctl enable mysql
   
   # 初始化MySQL（设置root密码、配置远程访问等）
   mysql_secure_installation
   ```

3. **安装Redis**
   ```bash
   # Ubuntu/Debian
   apt-get install -y redis-server
   
   # CentOS/RHEL
   yum install -y redis
   
   # 启动Redis服务
   systemctl start redis
   systemctl enable redis
   ```

#### 15.2.2 部署步骤

1. **编译项目**（可在本地或构建服务器上进行）
   ```bash
   # 克隆代码
   git clone https://github.com/rbacadmin/rbac_admin_server.git
   cd rbac_admin_server
   
   # 下载依赖
   go mod download
   
   # 编译（Linux环境）
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rbac_admin_server -ldflags="-s -w"
   ```

2. **复制文件到目标服务器**
   ```bash
   # 创建部署目录
   mkdir -p /opt/rbac_admin_server/{bin,config,logs,uploads}
   
   # 复制文件
   scp rbac_admin_server root@your_server:/opt/rbac_admin_server/bin/
   scp settings.yaml root@your_server:/opt/rbac_admin_server/config/
   scp .env root@your_server:/opt/rbac_admin_server/config/
   scp -r config/casbin/ root@your_server:/opt/rbac_admin_server/config/
   
   # 设置权限
   chmod +x /opt/rbac_admin_server/bin/rbac_admin_server
   ```

3. **创建systemd服务文件**
   ```bash
   # 创建服务文件
   vi /etc/systemd/system/rbac_admin_server.service
   
   # 文件内容
   [Unit]
   Description=RBAC Admin Server
   After=network.target mysql.service redis.service
   Requires=mysql.service redis.service
   
   [Service]
   Type=simple
   User=root
   WorkingDirectory=/opt/rbac_admin_server
   EnvironmentFile=/opt/rbac_admin_server/config/.env
   ExecStart=/opt/rbac_admin_server/bin/rbac_admin_server --config /opt/rbac_admin_server/config/settings.yaml
   Restart=on-failure
   RestartSec=5s
   
   [Install]
   WantedBy=multi-user.target
   ```

4. **启动服务**
   ```bash
   # 重新加载systemd配置
   systemctl daemon-reload
   
   # 启动服务
   systemctl start rbac_admin_server
   
   # 设置开机自启
   systemctl enable rbac_admin_server
   
   # 查看服务状态
   systemctl status rbac_admin_server
   ```

#### 15.2.3 Nginx反向代理配置

```bash
# 创建Nginx配置文件
vi /etc/nginx/conf.d/rbac_admin_server.conf

# 文件内容
server {
    listen 80;
    server_name api.yourdomain.com;
    
    # 访问日志
    access_log /var/log/nginx/rbac_admin_server_access.log;
    error_log /var/log/nginx/rbac_admin_server_error.log;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket支持（如果需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # 静态文件（如上传的文件）
    location /uploads {
        alias /opt/rbac_admin_server/uploads;
        expires 30d;
    }
    
    # 限制请求大小
    client_max_body_size 20M;
}
```

重启Nginx服务：
```bash
systemctl restart nginx
```

## 16. 维护与更新

### 16.1 更新流程

当需要更新应用时，请按照以下步骤进行：

1. **备份**
   ```bash
   # 备份数据库
   mysqldump -u root -p rbacadmin > rbacadmin_backup_$(date +%Y%m%d).sql
   
   # 备份配置文件
   cp /opt/rbac_admin_server/config/settings.yaml /opt/rbac_admin_server/config/settings.yaml.bak
   cp /opt/rbac_admin_server/config/.env /opt/rbac_admin_server/config/.env.bak
   
   # 备份日志文件（可选）
   tar -czf logs_backup_$(date +%Y%m%d).tar.gz /opt/rbac_admin_server/logs/
   ```

2. **拉取代码并更新**
   ```bash
   # 克隆最新代码或拉取更新
   git clone https://github.com/rbacadmin/rbac_admin_server.git /tmp/rbac_admin_server_new
   # 或 cd rbac_admin_server && git pull
   
   # 下载依赖
   cd /tmp/rbac_admin_server_new
   go mod download
   
   # 编译
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rbac_admin_server -ldflags="-s -w"
   ```

3. **部署新版本**
   ```bash
   # 停止服务
   systemctl stop rbac_admin_server
   
   # 替换二进制文件
   cp /tmp/rbac_admin_server_new/rbac_admin_server /opt/rbac_admin_server/bin/
   
   # 更新配置文件（如果有更改）
   # cp /tmp/rbac_admin_server_new/settings.yaml /opt/rbac_admin_server/config/
   
   # 更新Casbin配置（如果有更改）
   # cp -r /tmp/rbac_admin_server_new/config/casbin/ /opt/rbac_admin_server/config/
   
   # 重启服务
   systemctl start rbac_admin_server
   
   # 检查服务状态
   systemctl status rbac_admin_server
   ```

### 16.2 日志管理

应用的日志默认保存在`./logs`目录下：

1. **日志路径**：`/opt/rbac_admin_server/logs/app.log`

2. **配置日志**：可以在`settings.yaml`中调整日志级别和格式
   ```yaml
   log:
     level: info  # 可选：debug, info, warn, error
     format: json  # 可选：text, json
     dir: ./logs
     filename: app.log
     max_size: 100  # MB
     max_age: 7  # 天
     max_backups: 5
   ```

3. **日志滚动**：日志会根据配置自动滚动并压缩

4. **日志清理策略**
   ```bash
   # 创建日志清理脚本
   vi /opt/rbac_admin_server/scripts/clean_logs.sh
   
   # 脚本内容
   #!/bin/bash
   find /opt/rbac_admin_server/logs -name "*.log.*" -type f -mtime +30 -delete
   
   # 设置可执行权限
   chmod +x /opt/rbac_admin_server/scripts/clean_logs.sh
   
   # 添加到crontab（每月1号执行）
   echo "0 0 1 * * /opt/rbac_admin_server/scripts/clean_logs.sh" >> /etc/crontab
   ```

5. **实时查看日志**
   ```bash
   # 查看服务日志
   journalctl -u rbac_admin_server -f
   
   # 查看应用日志
   tail -f /opt/rbac_admin_server/logs/app.log
   ```

### 16.3 监控与维护

1. **健康检查**
   ```bash
   # 手动检查服务健康状态
   curl http://localhost:8080/health
   ```

2. **性能监控**
   - 应用内置Prometheus指标，访问：`http://localhost:8080/metrics`
   - 可以配置Grafana仪表板来可视化监控数据

3. **常见问题排查**
   ```bash
   # 检查端口占用
   netstat -tuln | grep 8080
   
   # 检查数据库连接
   mysql -u root -p -h localhost -P 3306 rbacadmin -e "SELECT 1;"
   
   # 检查Redis连接
   redis-cli ping
   ```

## 17. 联系支持

如果在部署过程中遇到任何问题，可以通过以下方式获取支持：

- Email: support@rbacadmin.com
- GitHub: https://github.com/rbacadmin/rbac_admin_server/issues
- 社区论坛: https://forum.rbacadmin.com