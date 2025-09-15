@echo off
color 0C
echo.
echo ╔═══════════════════════════════════════╗
echo ║        🏭 RBAC Admin Server          ║
echo ║     生产环境启动脚本 (Windows)        ║
echo ╚═══════════════════════════════════════╝
echo.
echo 🌍 启动生产环境...
echo 📁 使用配置文件: settings_prod.yaml
echo 🔒 安全模式: 最高级别
echo 📊 日志级别: info

rem 检查go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: Go未安装或未配置环境变量
    echo 📥 请访问 https://golang.org/dl/ 下载安装Go
    pause
    exit /b 1
)

echo ✅ Go环境检查通过

rem 检查必需的环境变量
echo 🔐 检查生产环境配置...

if "%DB_HOST%"=="" (
    echo ❌ 错误: DB_HOST 环境变量未设置
    echo 📝 请设置: set DB_HOST=your-db-host
    pause
    exit /b 1
)

if "%DB_PASSWORD%"=="" (
    echo ❌ 错误: DB_PASSWORD 环境变量未设置
    echo 📝 请设置: set DB_PASSWORD=your-db-password
    pause
    exit /b 1
)

if "%JWT_SECRET%"=="" (
    echo ❌ 错误: JWT_SECRET 环境变量未设置
    echo 📝 请设置: set JWT_SECRET=your-256-bit-secret-key
    pause
    exit /b 1
)

echo ✅ 环境变量检查完成

rem 检查配置文件
if not exist "settings_prod.yaml" (
    echo ❌ 错误: settings_prod.yaml 配置文件不存在
    echo 📝 请确保配置文件在正确位置
    pause
    exit /b 1
)

echo ✅ 配置文件检查完成

rem 检查依赖
echo 📦 检查项目依赖...
go mod tidy >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 依赖下载失败
    pause
    exit /b 1
)

echo ✅ 依赖检查完成
echo.
echo 🚀 正在启动生产服务器...
echo 🌐 生产环境已启动
echo 📊 监控地址: http://localhost:8080/metrics
echo 📚 API文档: http://localhost:8080/swagger/index.html
echo.
echo ⚠️  警告: 生产环境请勿随意停止！
echo ═════════════════════════════════════════
echo.

go run main.go -env=prod

echo.
echo 服务器已停止
pause