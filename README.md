# RBAC管理员服务器

一个基于Go语言的现代化RBAC（基于角色的访问控制）管理员服务器，支持多种数据库连接，采用工厂模式设计。

## ✨ 特性

- 🏭 **数据库工厂模式**: 支持MySQL、PostgreSQL、SQLite、SQL Server
- 🔧 **企业级配置**: 11个配置模块，支持环境变量覆盖
- 🔄 **自动迁移**: 数据库表结构自动创建和更新
- 🔐 **JWT认证**: 安全的用户认证和权限管理
- 📊 **结构化日志**: 基于Zap的高性能日志系统
- 🐳 **Docker支持**: 完整的容器化部署方案
- 📈 **监控指标**: Prometheus指标收集和健康检查
- 📚 **API文档**: 自动生成Swagger API文档

## 🚀 快速开始

### 1. 环境要求
- Go 1.21+
- Docker (可选)
- MySQL/PostgreSQL/SQLite (根据配置选择)

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置项目
复制配置文件模板：
```bash
cp settings.yaml.example settings.yaml
```

### 4. 启动项目

#### 使用SQLite（推荐开发环境）
```bash
# 无需安装数据库，直接运行
go run main.go
```

#### 使用MySQL
```bash
# 设置环境变量
export DB_TYPE=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=yourpassword
export DB_NAME=rbac_admin

# 运行项目
go run main.go
```

#### 使用PostgreSQL
```bash
# 设置环境变量
export DB_TYPE=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_NAME=rbac_admin

# 运行项目
go run main.go
```

### 5. 验证启动
访问以下地址验证项目是否正常运行：
- API文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health
- 默认管理员: admin/admin123

## 🏗️ 项目结构

```
rbac_admin_server/
├── 📂 api/                    # RESTful API接口
├── 📂 config/                 # 配置管理系统
│   ├── config.go             # 配置结构体定义
│   └── loader.go             # 配置加载器
├── 📂 core/                   # 核心启动逻辑
│   ├── initializer.go        # 项目初始化器
│   ├── audit/                # 审计日志系统
│   └── errors/               # 错误处理系统
├── 📂 database/               # 数据库工厂系统
│   ├── database_factory.go   # 数据库工厂核心实现
│   ├── migrator.go           # 数据库迁移和初始化
│   └── models/               # 数据模型定义
│       ├── user.go           # 用户、角色、权限模型
├── 📂 examples/               # 使用示例
├── 📂 global/                 # 全局变量管理
├── 📝 main.go                # 程序入口
├── ⚙️ settings.yaml          # 配置文件
├── 🐳 docker-compose.yml     # Docker部署配置
├── 🐳 Dockerfile            # Docker镜像构建
└── 📦 go.mod                 # 依赖管理
```

## 📊 数据库支持

### 配置示例

#### MySQL
```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  name: "rbac_admin"
  ssl_mode: "disable"
  timeout: 30
```

#### PostgreSQL
```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "postgres"
  name: "rbac_admin"
  ssl_mode: "disable"
  timeout: 30
```

#### SQLite（开发推荐）
```yaml
database:
  type: "sqlite"
  path: "./rbac_admin.db"
```

## 🐳 Docker部署

### 使用Docker Compose
```bash
# 启动所有服务（MySQL + Redis + RBAC服务器）
docker-compose up -d

# 查看日志
docker-compose logs -f rbac_server

# 停止服务
docker-compose down
```

### 使用PostgreSQL
```bash
# 修改docker-compose.yml中的数据库配置
# 然后启动
docker-compose up -d
```

### 单独构建镜像
```bash
# 构建镜像
docker build -t rbac-admin-server .

# 运行容器
docker run -p 8080:8080 rbac-admin-server
```

## 🔐 默认账户

首次启动时自动创建：
- **用户名**: admin
- **密码**: admin123
- **角色**: 超级管理员
- **权限**: 全部权限

## 📚 API文档

项目启动后访问：
- Swagger文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health
- 指标监控: http://localhost:8080/metrics

## 🔧 开发指南

### 添加新的API接口
1. 在`api/`目录下创建新的路由文件
2. 实现对应的处理器函数
3. 在main.go中注册路由
4. 更新Swagger文档

### 数据库模型扩展
1. 在`database/models/`目录下添加新的模型
2. 更新`migrator.go`中的迁移逻辑
3. 运行项目自动迁移

### 配置扩展
1. 在`config/config.go`中添加新的配置结构
2. 更新`settings.yaml`模板
3. 在`loader.go`中添加验证逻辑

## 📝 环境变量

### 数据库配置
```bash
# 数据库类型: mysql, postgres, sqlite, sqlserver
DB_TYPE=mysql

# MySQL配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=rbac_admin

# PostgreSQL配置
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=rbac_admin

# SQLite配置
DB_PATH=./rbac_admin.db
```

### 服务器配置
```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=60
SERVER_WRITE_TIMEOUT=60
```

### JWT配置
```bash
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE=7200
JWT_ISSUER=rbac_admin
```

### Redis配置
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 🔍 故障排查

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查数据库服务
   docker-compose ps
   
   # 查看日志
   docker-compose logs mysql
   ```

2. **权限错误**
   ```bash
   # 检查数据库权限
   mysql -u root -p
   GRANT ALL PRIVILEGES ON rbac_admin.* TO 'rbac_user'@'%';
   FLUSH PRIVILEGES;
   ```

3. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tulnp | grep :8080
   
   # 修改配置端口
   export SERVER_PORT=8081
   ```

### 调试模式
```bash
# 启用调试日志
export LOG_LEVEL=debug
./rbac_admin_server
```

## 📈 性能优化

### 数据库优化
- 自动创建索引优化查询性能
- 连接池参数自动配置
- 查询预加载减少N+1问题

### 缓存策略
- Redis缓存热点数据
- JWT令牌缓存验证
- API响应缓存

### 监控指标
- Prometheus指标收集
- 健康检查端点
- 性能监控面板

## 🤝 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送分支
5. 创建Pull Request

## 📄 许可证

MIT License - 详见LICENSE文件

## 🙋‍♂️ 支持

如有问题，请创建Issue或联系维护者。

---

**这是一个企业级的RBAC管理员服务器，具备完整的配置管理和数据库工厂模式，适合生产环境部署。**