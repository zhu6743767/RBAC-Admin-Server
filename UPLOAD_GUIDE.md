# 项目上传指南

## 🔒 安全上传步骤

### 1. 准备工作
确保您已经：
- 移除了所有敏感信息
- 创建了安全的配置文件模板
- 添加了.gitignore文件

### 2. 创建GitHub仓库
在GitHub上创建一个新的仓库：
1. 访问 https://github.com/new
2. 仓库名称：`RBAC-Admin-Server`
3. 设置为公开仓库
4. 不要初始化README（我们已有）

### 3. 本地推送命令

```bash
# 如果远程仓库不存在，先移除旧的远程
# git remote remove origin

# 添加新的远程仓库（替换为您的用户名）
git remote add origin https://github.com/YOUR_USERNAME/RBAC-Admin-Server.git

# 推送代码
git push -u origin master
```

### 4. 已处理的安全事项

✅ **已移除的敏感信息：**
- `settings.yaml` - 包含数据库密码
- 所有可执行文件 (*.exe)
- 日志文件
- 环境变量文件

✅ **已添加的安全措施：**
- `.gitignore` - 排除敏感文件
- `settings.example.yaml` - 安全配置模板
- 更新的README.md - 包含安全提醒

### 5. 配置文件处理

#### 敏感文件列表（已排除）：
```
settings.yaml          # 主配置文件（包含密码）
*.exe                  # 可执行文件
logs/                  # 日志目录
.env*                  # 环境变量文件
*.pem, *.key, *.crt    # 证书文件
```

#### 安全模板文件：
```
settings.example.yaml  # 配置模板
README.md              # 安全的使用文档
.gitignore            # 排除规则
```

### 6. 首次使用者的配置步骤

新用户需要：
1. 克隆仓库：`git clone https://github.com/YOUR_USERNAME/RBAC-Admin-Server.git`
2. 复制配置：`cp settings.example.yaml settings.yaml`
3. 编辑配置：修改数据库密码等信息
4. 运行项目：`go run main.go`

### 7. 验证上传

上传后检查：
- [ ] 仓库中没有settings.yaml文件
- [ ] 包含settings.example.yaml模板
- [ ] .gitignore文件存在并正确配置
- [ ] README.md包含安全提醒
- [ ] 没有可执行文件或日志文件

### 8. 后续维护

每次更新前检查：
```bash
# 检查是否有敏感文件被意外添加
git status

# 如果有敏感文件，移除它们
git rm --cached sensitive_file.yaml
git commit -m "security: remove sensitive file"
git push
```

## 🚨 重要提醒

**永远不要提交到GitHub的文件：**
- 包含真实密码的配置文件
- SSL证书和私钥
- 数据库连接字符串
- API密钥和访问令牌
- 生产环境的具体配置

**安全最佳实践：**
1. 使用环境变量存储敏感信息
2. 定期轮换密码和密钥
3. 为不同环境使用不同配置
4. 启用双因素认证
5. 定期审计仓库内容