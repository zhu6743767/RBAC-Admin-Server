@echo off
color 0A
echo.
echo ╔═══════════════════════════════════════╗
echo ║        🚀 RBAC Admin Server           ║
echo ║     开发环境启动脚本 (Windows)        ║
echo ╚═══════════════════════════════════════╝
echo.
echo 🌍 启动开发环境...
echo 📁 使用配置文件: settings_dev.yaml
echo 🗄️  数据库: SQLite (无需MySQL)
echo 📊 日志级别: debug
echo.

rem 检查go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: Go未安装或未配置环境变量
    echo 📥 请访问 https://golang.org/dl/ 下载安装Go
    pause
    exit /b 1
)

echo ✅ Go环境检查通过

rem 检查依赖
echo 📦 检查项目依赖...
go mod tidy >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 依赖下载失败
    pause
    exit /b 1
)

echo ✅ 依赖检查完成

rem 检查配置文件
if not exist "settings_dev.yaml" (
    echo ❌ 错误: settings_dev.yaml 配置文件不存在
    echo 📝 请确保配置文件在正确位置
    pause
    exit /b 1
)

echo ✅ 配置文件检查完成
echo.
echo 🚀 正在启动服务器...
echo 🌐 访问地址: http://localhost:8080
echo 📚 API文档: http://localhost:8080/swagger/index.html
echo.
echo 按 Ctrl+C 停止服务器
echo ═════════════════════════════════════════
echo.

go run main.go -env=dev

echo.
echo 服务器已停止
pause