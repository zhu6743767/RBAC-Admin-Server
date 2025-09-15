# 🔄 RBAC管理员服务器 - 环境切换指南

## 📋 项目概述

RBAC管理员服务器是一个基于角色的访问控制系统，支持多环境配置切换，包括开发环境(dev)、测试环境(test)和生产环境(prod)。

## 🚀 快速启动

### 1. 一键启动

#### Windows系统
```bash
# 启动开发环境（默认）
double-click run.bat
# 或命令行
run.bat dev

# 启动测试环境
run.bat test

# 启动生产环境
run.bat prod
```

#### Linux/Mac系统
```bash
# 启动开发环境（默认）
./run-dev.sh

# 启动测试环境
./run-test.sh

# 启动生产环境
./run-prod.sh
```

### 2. 手动启动

```bash
# 使用特定环境配置
go run main.go -env=dev

# 使用自定义配置文件
go run main.go -config=myconfig.yaml

# 查看帮助信息
go run main.go -h
```

## 🌍 环境对比

| 特性 | 开发环境(dev) | 测试环境(test) | 生产环境(prod) |
|---|---|---|---|
| **端口** | 8080 | 8081 | 80/443 |
| **日志级别** | debug | info | warn/error |
| **数据库** | SQLite | MySQL测试库 | MySQL生产库 |
| **CORS** | 允许所有来源 | 限制来源 | 严格限制 |
| **Swagger** | ✅ 启用 | ✅ 启用 | ❌ 禁用 |
| **调试模式** | ✅ 启用 | ❌ 禁用 | ❌ 禁用 |
| **性能分析** | ✅ 启用 | ❌ 禁用 | ❌ 禁用 |

## 📁 配置文件

### 环境配置文件

每个环境都有独立的配置文件：

- **开发环境**: `settings_dev.yaml`
- **测试环境**: `settings_test.yaml`
- **生产环境**: `settings_prod.yaml`

### 配置模板

#### 开发环境配置 (settings_dev.yaml)
```yaml
server:
  port: 8080
  mode: "debug"

database:
  type: "sqlite"
  path: "./data/rbac_admin_dev.db"

log:
  level: "debug"
  output: "both"

cors:
  enable: true
  allow_origins: ["*"]

swagger:
  enable: true
  enable_ui: true
```

#### 测试环境配置 (settings_test.yaml)
```yaml
server:
  port: 8081
  mode: "test"

database:
  type: "mysql"
  host: "localhost"
  port: 3306
  username: "test_user"
  password: "test_password"
  database: "rbac_admin_test"

log:
  level: "info"
  output: "file"

cors:
  enable: true
  allow_origins: ["http://test.example.com"]

swagger:
  enable: true
  enable_ui: false
```

#### 生产环境配置 (settings_prod.yaml)
```yaml
server:
  port: 443
  mode: "release"
  read_timeout: 30s
  write_timeout: 30s

database:
  type: "mysql"
  host: "prod-db.example.com"
  port: 3306
  username: "prod_user"
  password: "secure_password"
  database: "rbac_admin_prod"

log:
  level: "warn"
  output: "file"
  log_dir: "/var/log/rbac-admin"

cors:
  enable: true
  allow_origins: ["https://admin.example.com"]

swagger:
  enable: false
  enable_ui: false

security:
  bcrypt_cost: 12
  max_login_attempts: 5
  lock_duration_minutes: 30
```

## 🔧 环境变量

### 常用环境变量

| 变量名 | 说明 | 示例 |
|---|---|---|
| `ENV` | 运行环境 | `dev`, `test`, `prod` |
| `CONFIG_PATH` | 配置文件路径 | `./config/settings.yaml` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_USER` | 数据库用户 | `root` |
| `DB_PASSWORD` | 数据库密码 | `password` |

### 环境变量优先级

1. 命令行参数 (`-config`)
2. 环境变量
3. 默认配置文件

## 🧪 测试验证

### 1. 环境验证

```bash
# 测试开发环境
go run main.go -env=dev
# 预期：SQLite数据库，端口8080，调试模式

# 测试测试环境
go run main.go -env=test
# 预期：MySQL测试库，端口8081，日志级别info

# 测试生产环境
go run main.go -env=prod
# 预期：MySQL生产库，端口443，日志级别warn
```

### 2. 功能验证

#### 开发环境验证
- [ ] 访问 http://localhost:8080/health
- [ ] 访问 http://localhost:8080/swagger/index.html
- [ ] 检查日志输出是否详细
- [ ] 验证SQLite数据库文件创建

#### 测试环境验证
- [ ] 访问 http://localhost:8081/health
- [ ] 验证MySQL连接
- [ ] 检查日志级别是否为info
- [ ] 验证CORS设置

#### 生产环境验证
- [ ] 验证SSL证书配置
- [ ] 检查安全设置
- [ ] 验证日志文件轮转
- [ ] 检查性能监控

## 🔍 故障排除

### 常见问题

#### 1. 端口占用
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8080
kill -9 <PID>
```

#### 2. 数据库连接失败
```bash
# 检查MySQL服务
systemctl status mysql

# 检查连接配置
mysql -h localhost -u root -p
```

#### 3. 权限问题
```bash
# Linux/Mac
chmod +x run-dev.sh

# Windows
# 确保有管理员权限
```

### 日志分析

#### 开发环境日志位置
- 控制台输出
- `./logs/dev/rbac-admin.log`

#### 测试环境日志位置
- `./logs/test/rbac-admin.log`

#### 生产环境日志位置
- `/var/log/rbac-admin/rbac-admin.log`

## 📊 性能优化

### 开发环境优化
- 启用热重载
- 减少日志输出
- 使用内存数据库

### 测试环境优化
- 启用连接池
- 优化查询性能
- 启用缓存

### 生产环境优化
- 启用CDN
- 数据库读写分离
- 启用集群部署

## 🚀 部署流程

### 1. 开发环境部署
```bash
git clone <repository>
cd rbac-admin-server
go mod tidy
./run-dev.sh
```

### 2. 测试环境部署
```bash
# 使用Docker
docker-compose -f docker-compose.test.yml up -d

# 或直接运行
./run-test.sh
```

### 3. 生产环境部署
```bash
# 使用Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# 或使用Systemd
systemctl start rbac-admin
```

## 📚 相关文档

- [API文档](docs/API.md)
- [部署指南](docs/DEPLOYMENT.md)
- [配置说明](docs/CONFIG.md)
- [开发指南](docs/DEVELOPMENT.md)

## 🆘 技术支持

如有问题，请联系：
- 📧 Email: support@rbac-admin.com
- 💬 微信群: RBAC管理员技术群
- 📱 技术热线: 400-123-4567

---

**最后更新**: 2025年1月
**版本**: v1.0.0