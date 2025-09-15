# RBAC Admin Server - GitHub 上传指南

## 🔧 步骤1：创建GitHub仓库

1. 访问 [GitHub.com](https://github.com)
2. 点击右上角的 "+" → "New repository"
3. 设置仓库信息：
   - **Repository name**: `RBAC-Admin-Server`
   - **Description**: A modern Role-Based Access Control (RBAC) admin server built with Go
   - **Public**: 选择公开仓库
   - **Initialize repository**: ❌ 不要勾选任何初始化选项
4. 点击 "Create repository"

## 🚀 步骤2：推送代码到GitHub

### 方法一：使用HTTPS（推荐）
```bash
# 设置远程仓库地址（替换为你的实际用户名）
git remote set-url origin https://github.com/YOUR_USERNAME/RBAC-Admin-Server.git

# 推送代码
git push -u origin master
```

### 方法二：使用SSH（需要配置SSH密钥）
```bash
# 设置SSH远程地址
git remote set-url origin git@github.com:YOUR_USERNAME/RBAC-Admin-Server.git

# 推送代码
git push -u origin master
```

## ✅ 验证上传成功

1. 访问 `https://github.com/YOUR_USERNAME/RBAC-Admin-Server`
2. 确认以下内容：
   - ✅ 仓库可见
   - ✅ 包含 `.gitignore` 文件
   - ✅ 包含 `settings.example.yaml`（不是 `settings.yaml`）
   - ✅ 包含 `README.md`
   - ✅ 包含 `UPLOAD_GUIDE.md`

## 🔍 安全检查清单

上传后请确认以下敏感信息已排除：

- [ ] `settings.yaml` 文件不存在
- [ ] `*.exe` 可执行文件不存在
- [ ] `logs/` 目录不存在
- [ ] `.env` 文件不存在
- [ ] 数据库密码未泄露
- [ ] JWT密钥未泄露

## 📋 用户使用指南

新用户克隆项目后：

```bash
git clone https://github.com/YOUR_USERNAME/RBAC-Admin-Server.git
cd RBAC-Admin-Server

# 复制配置文件
cp settings.example.yaml settings.yaml

# 编辑配置文件
# 修改 settings.yaml 中的实际配置

# 运行项目
go mod tidy
go run main.go
```

## 🛠️ 常见问题解决

### 如果推送失败
```bash
# 强制推送（谨慎使用）
git push -f origin master

# 或者先拉取远程更改
git pull origin master --rebase
git push origin master
```

### 如果认证失败
1. 检查GitHub用户名和密码
2. 或者使用Personal Access Token代替密码
3. 或者配置SSH密钥

## 📞 支持

如有问题，请查看：
- [UPLOAD_GUIDE.md](./UPLOAD_GUIDE.md) - 详细上传指南
- [README.md](./README.md) - 项目使用说明
- [docs/](./docs/) - 项目文档目录