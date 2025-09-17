# 🔐 RBAC管理员服务器 - 安全配置指南

## 📋 概述
本文档说明如何安全地配置和管理RBAC管理员服务器的敏感信息，确保生产环境的安全性。

## 🛡️ 敏感信息保护策略

### 1. 环境变量优先原则
- **生产环境**：所有敏感信息必须通过环境变量配置
- **开发环境**：可以使用配置文件，但建议使用环境变量
- **禁止**：将真实密码、密钥提交到版本库

### 2. 敏感信息分类

#### 🔑 数据库凭据
- `DB_PASSWORD` - 数据库密码
- `DB_HOST` - 数据库主机地址
- `DB_USER` - 数据库用户名

#### 🔐 JWT安全密钥
- `JWT_SECRET` - JWT签名密钥（建议32位以上）
- `JWT_ISSUER` - JWT发行者
- `JWT_AUDIENCE` - JWT受众

#### 🔄 Redis配置
- `REDIS_PASSWORD` - Redis密码
- `REDIS_ADDR` - Redis地址

#### 🔒 其他安全配置
- `CSRF_SECRET` - CSRF保护密钥
- 各类API密钥和令牌

## 📝 配置方法

### 开发环境配置
1. 复制 `.env.example` 为 `.env`
2. 填写真实的配置值
3. 确保 `.env` 文件在 `.gitignore` 中

### 生产环境配置
1. **不要创建 `.env` 文件**
2. 通过操作系统的环境变量设置
3. 使用密钥管理服务（如Docker Secrets、Kubernetes Secrets等）

### Docker环境示例
```bash
# 使用环境变量运行
docker run -e DB_PASSWORD=your_password \
           -e JWT_SECRET=your_jwt_secret \
           -p 8080:8080 \
           rbac-admin-server
```

### Linux系统环境变量
```bash
# 临时设置
export DB_PASSWORD="your_password"
export JWT_SECRET="your_jwt_secret_key_minimum_32_characters"

# 永久设置（添加到 /etc/environment 或 ~/.bashrc）
echo 'export DB_PASSWORD="your_password"' >> ~/.bashrc
echo 'export JWT_SECRET="your_jwt_secret"' >> ~/.bashrc
source ~/.bashrc
```

### Windows系统环境变量
```powershell
# 临时设置
$env:DB_PASSWORD = "your_password"
$env:JWT_SECRET = "your_jwt_secret_key_minimum_32_characters"

# 永久设置（系统环境变量）
[Environment]::SetEnvironmentVariable("DB_PASSWORD", "your_password", "Machine")
[Environment]::SetEnvironmentVariable("JWT_SECRET", "your_jwt_secret", "Machine")
```

## ⚠️ 安全警告

### ❌ 禁止行为
- 不要将真实密码、密钥提交到Git仓库
- 不要在代码中硬编码敏感信息
- 不要将 `.env` 文件上传到生产服务器
- 不要使用弱密码或短密钥

### ✅ 推荐做法
- 使用强密码（12位以上，包含大小写字母、数字、特殊字符）
- JWT密钥至少32位，使用随机生成的字符串
- 定期更换密钥和密码
- 使用密钥管理服务
- 启用SSL/TLS加密传输
- 限制数据库和Redis的网络访问

## 🔍 配置验证

### 开发环境验证
```bash
# 检查配置是否加载成功
go run validate_config.go -env dev

# 启动服务测试
go run main.go -env dev
```

### 生产环境验证
```bash
# 验证生产环境配置（干运行模式）
go run validate_config.go -env prod

# 检查环境变量是否设置正确
echo $DB_PASSWORD  # Linux/Mac
$env:DB_PASSWORD   # PowerShell
```

## 🚀 部署检查清单

部署到生产环境前，请确认：

- [ ] 所有敏感信息通过环境变量配置
- [ ] `.env` 文件不存在于生产环境
- [ ] 数据库密码强度足够
- [ ] JWT密钥长度至少32位且为随机字符串
- [ ] Redis已设置密码保护（如需要）
- [ ] 数据库和Redis只允许可信IP访问
- [ ] 启用SSL/TLS证书
- [ ] 定期备份和密钥轮换策略

## 📞 问题反馈

如发现安全配置问题，请立即：
1. 停止服务
2. 更换所有受影响的密钥和密码
3. 检查访问日志
4. 联系安全团队

---

**记住：安全是第一优先级！永远不要在代码中暴露敏感信息。**