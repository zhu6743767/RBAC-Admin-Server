# RBAC管理服务器 - 简洁版

## 功能说明
这是一个基于Go语言的RBAC（基于角色的访问控制）管理服务器，提供用户管理和权限控制功能。

## 核心功能
- ✅ 管理员用户创建
- ✅ 用户登录验证
- ✅ JWT Token认证
- ✅ 权限管理

## 快速开始

### 1. 创建管理员用户
```bash
go run create_admin_simple.go
```
**默认管理员账号**: `admin/admin123`
**邮箱**: `admin@example.com`
**电话**: `13800138000`

### 2. 启动服务器
```bash
go run main.go -env dev
```
服务器将在 `http://127.0.0.1:8080` 启动

### 3. 测试登录
```bash
# 简单测试
go run simple_login.go

# 详细调试测试
go run debug_login.go
```

## API端点
- **登录**: `POST /api/public/login`
  - 请求体: `{"username": "admin", "password": "admin123"}`
  - 响应: JWT Token + 用户信息

## 文件说明
- `create_admin_simple.go` - 管理员用户创建程序
- `simple_login.go` - 简单登录测试
- `debug_login.go` - 详细登录调试测试
- `main.go` - 服务器主程序

## 环境配置
- 开发环境: `-env dev`
- 生产环境: `-env prod`
- 测试环境: `-env test`

## 数据库
使用MySQL数据库，连接配置在 `settings_*.yaml` 文件中。

## 状态码
- `200` - 成功
- `401` - 认证失败
- `500` - 服务器错误