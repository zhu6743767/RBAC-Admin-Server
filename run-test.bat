@echo off
color 0E
echo.
echo ╔═══════════════════════════════════════╗
echo ║        🧪 RBAC Admin Server           ║
echo ║     测试环境启动脚本 (Windows)        ║
echo ╚═══════════════════════════════════════╝
echo.
echo 🌍 启动测试环境...
echo 📁 使用配置文件: settings_test.yaml
echo 🗄️  数据库: MySQL测试数据库
echo 📊 日志级别: info
echo 🔒 安全模式: 启用

rem 检查go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: Go未安装或未配置环境变量
    echo 📥 请访问 https://golang.org/dl/ 下载安装Go
    pause
    exit /b 1
)

echo ✅ Go环境检查通过

rem 检查MySQL连接（可选）
echo 🗄️  检查MySQL测试数据库连接...
mysql -h localhost -P 3306 -u root -p123456 -e "SELECT 1" rbac_admin_test >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  警告: MySQL测试数据库连接失败
    echo 📝 请确保MySQL已启动且测试数据库存在
    echo 📝 数据库名: rbac_admin_test
    echo 📝 用户: root / 密码: 123456
    echo.
    choice /C YN /M "是否继续启动？(Y/N)"
    if errorlevel 2 exit /b 1
)

echo ✅ MySQL连接检查完成

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
if not exist "settings_test.yaml" (
    echo ❌ 错误: settings_test.yaml 配置文件不存在
    echo 📝 请确保配置文件在正确位置
    pause
    exit /b 1
)

echo ✅ 配置文件检查完成
echo.
echo 🚀 正在启动测试服务器...
echo 🌐 访问地址: http://localhost:8081
echo 📚 API文档: http://localhost:8081/swagger/index.html
echo.
echo 按 Ctrl+C 停止服务器
echo ═════════════════════════════════════════
echo.

go run main.go -env=test

echo.
echo 服务器已停止
pause