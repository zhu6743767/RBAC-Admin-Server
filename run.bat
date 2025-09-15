@echo off
chcp 65001 > nul
:: =================================================================================
:: 🚀 RBAC管理员服务器 - 通用启动脚本 (Windows)
:: =================================================================================
:: 📋 使用说明：
::   双击运行：默认启动开发环境
::   命令行参数：run.bat [dev|test|prod]
::   示例：run.bat prod
:: =================================================================================

:: 设置窗口标题
title RBAC管理员服务器启动器

:: 清除屏幕
cls

:: 显示启动信息
echo.
echo   ╔═══════════════════════════════════════╗
echo   ║          RBAC管理员服务器启动器          ║
echo   ║    Role-Based Access Control System   ║
echo   ╚═══════════════════════════════════════╝
echo.

:: 设置环境变量，默认为dev
set ENV=dev
if not "%1"=="" set ENV=%1

:: 验证环境参数
if "%ENV%"=="dev" goto VALID_ENV
if "%ENV%"=="test" goto VALID_ENV
if "%ENV%"=="prod" goto VALID_ENV

echo ❌ 错误：不支持的环境参数 "%ENV%"
echo 📋 支持的环境：dev(开发环境) test(测试环境) prod(生产环境)
echo 💡 示例：%0 dev
echo.
pause
exit /b 1

:VALID_ENV
echo   📋 启动环境：%ENV%
echo.

:: 检查Go环境
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误：未检测到Go运行环境
    echo 📥 请访问 https://golang.org/dl/ 下载并安装Go
    echo.
    pause
    exit /b 1
)

:: 显示Go版本
echo ✅ Go环境检测通过
go version
echo.

:: 检查依赖包
echo 📦 检查项目依赖...
go mod tidy >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 依赖包检查失败，尝试修复...
    go mod download
    go mod tidy
)

:: 检查配置文件
set CONFIG_FILE=settings_%ENV%.yaml
if not exist "%CONFIG_FILE%" (
    echo ❌ 错误：未找到配置文件 %CONFIG_FILE%
    echo.
    pause
    exit /b 1
)

:: 创建必要的目录
if not exist "logs" mkdir logs
if not exist "data" mkdir data
if not exist "uploads" mkdir uploads

:: 启动服务器
echo 🚀 正在启动服务器...
echo.
echo ═════════════════════════════════════════════════
echo.

:: 运行服务器
go run main.go -env=%ENV%

:: 检查退出状态
if %errorlevel% neq 0 (
    echo.
    echo ❌ 服务器启动失败，请检查错误信息
    echo.
    pause
    exit /b %errorlevel%
)

echo.
echo ✅ 服务器已关闭
echo.
pause