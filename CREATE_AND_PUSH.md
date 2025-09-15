# 创建GitHub仓库并推送项目

## 🎯 步骤1：手动创建GitHub仓库

### 在GitHub网页创建仓库：

1. **访问GitHub**：打开浏览器访问 https://github.com
2. **登录账号**：使用你的GitHub账号登录
3. **创建新仓库**：
   - 点击右上角的 "+" 图标
   - 选择 "New repository"
   - 或使用直接链接：https://github.com/new

4. **填写仓库信息**：
   ```
   Repository name: RBAC-Admin-Server
   Description: A modern Role-Based Access Control (RBAC) admin server built with Go. Features include user management, role-based permissions, audit logging, and multi-database support.
   
   Public: ✓ (选择公开)
   Initialize this repository with: ❌ 不要勾选任何选项
   ```

5. **点击创建**：点击 "Create repository" 按钮

## 🚀 步骤2：推送本地代码到GitHub

### 打开终端/PowerShell，执行以下命令：

```bash
# 1. 进入项目目录
cd e:\myblog\Go项目学习\rbac_admin_server

# 2. 设置远程仓库地址（替换zhu6743767为你的实际用户名）
git remote set-url origin https://github.com/zhu6743767/RBAC-Admin-Server.git

# 3. 推送代码到GitHub
git push -u origin master
```

### 如果推送遇到认证问题：

#### 方法A：使用Personal Access Token
```bash
# 当提示输入密码时，使用GitHub Personal Access Token
# 获取Token：https://github.com/settings/tokens
```

#### 方法B：使用SSH方式
```bash
# 1. 切换为SSH地址
git remote set-url origin git@github.com:zhu6743767/RBAC-Admin-Server.git

# 2. 推送
git push -u origin master
```

## ✅ 步骤3：验证上传成功

1. **访问仓库**：https://github.com/zhu6743767/RBAC-Admin-Server
2. **检查文件**：确认以下文件存在：
   - `.gitignore`
   - `settings.example.yaml`
   - `README.md`
   - `UPLOAD_GUIDE.md`
   - `PUSH_TO_GITHUB.md`

3. **确认安全**：
   - ✅ `settings.yaml` 不存在
   - ✅ 无敏感信息泄露
   - ✅ 所有源代码已上传

## 🔧 一键创建脚本

### 创建一键推送脚本：

保存为 `create_and_push.bat`：

```batch
@echo off
echo ======================================
echo RBAC Admin Server - GitHub 创建和上传
echo ======================================
echo.

:START
echo 请选择操作：
echo 1. 创建仓库并推送（需要已手动创建仓库）
echo 2. 仅推送（仓库已存在）
echo 3. 查看当前状态
echo 4. 退出
set /p choice=输入选项(1-4): 

if "%choice%"=="1" goto PUSH
if "%choice%"=="2" goto PUSH_ONLY
if "%choice%"=="3" goto STATUS
if "%choice%"=="4" goto EXIT

goto START

:PUSH
echo.
echo 请先确保已在GitHub创建仓库！
echo 地址：https://github.com/new
echo 仓库名：RBAC-Admin-Server
echo.
set /p username=请输入GitHub用户名: 
if "%username%"=="" goto PUSH

git remote set-url origin https://github.com/%username%/RBAC-Admin-Server.git
goto PUSH_CODE

:PUSH_ONLY
set /p username=请输入GitHub用户名: 
if "%username%"=="" goto PUSH_ONLY
git remote set-url origin https://github.com/%username%/RBAC-Admin-Server.git

:PUSH_CODE
echo.
echo 正在推送代码...
git push -u origin master

echo.
echo 操作完成！
echo 访问：https://github.com/%username%/RBAC-Admin-Server
echo.
pause
goto EXIT

:STATUS
echo.
echo 当前Git状态：
git status
echo.
echo 当前远程配置：
git remote -v
echo.
pause
goto START

:EXIT
pause