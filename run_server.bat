@echo off
REM RBAC管理员服务器启动脚本
REM =================================================================================
REM 此脚本用于启动RBAC管理员服务器，支持开发环境和生产环境
REM =================================================================================

REM 设置中文编码
chcp 65001 > nul

REM 定义颜色变量
set "COLOR_SUCCESS=[92m"
set "COLOR_ERROR=[91m"
set "COLOR_INFO=[96m"
set "COLOR_RESET=[0m"

REM 默认环境为开发环境
goto :menu

:menu
cls
echo. 
echo %COLOR_INFO%================================================================================%COLOR_RESET%
echo %COLOR_INFO%                RBAC管理员服务器启动菜单                %COLOR_RESET%
echo %COLOR_INFO%================================================================================%COLOR_RESET%
echo. 
echo %COLOR_INFO%[1] 开发环境启动 (调试模式, 详细日志)%COLOR_RESET%
echo %COLOR_INFO%[2] 生产环境启动 (生产模式, 性能优化)%COLOR_RESET%
echo %COLOR_INFO%[3] 数据库迁移 (更新数据库结构)%COLOR_RESET%
echo %COLOR_INFO%[4] 创建管理员用户%COLOR_RESET%
echo %COLOR_INFO%[0] 退出%COLOR_RESET%
echo. 

set /p choice=请输入选择 (0-4): 

if %choice% equ 1 goto :dev_start
if %choice% equ 2 goto :prod_start
if %choice% equ 3 goto :db_migrate
if %choice% equ 4 goto :create_admin
if %choice% equ 0 goto :exit

echo %COLOR_ERROR%无效的选择，请重试。%COLOR_RESET%
pause
goto :menu

:dev_start
REM 开发环境配置
echo %COLOR_INFO%正在以开发环境模式启动服务器...%COLOR_RESET%
set APP_ENVIRONMENT=development
set APP_DEBUG=true
set SYSTEM_MODE=debug
set LOG_LEVEL=debug

REM 启动服务器
go run main.go -m server -settings settings.yaml
goto :exit

:prod_start
REM 生产环境配置
echo %COLOR_INFO%正在以生产环境模式启动服务器...%COLOR_RESET%
set APP_ENVIRONMENT=production
set APP_DEBUG=false
set SYSTEM_MODE=release
set LOG_LEVEL=info

REM 编译并启动服务器
echo %COLOR_INFO%正在编译服务器...%COLOR_RESET%
go build -o rbac_admin_server.exe
if %errorlevel% neq 0 (
echo %COLOR_ERROR%编译失败，请检查代码错误。%COLOR_RESET%
pause
goto :menu
)
echo %COLOR_SUCCESS%编译成功，正在启动服务器...%COLOR_RESET%
start /b rbac_admin_server.exe -m server -settings settings.yaml

echo %COLOR_SUCCESS%服务器已在后台启动。%COLOR_RESET%
echo %COLOR_INFO%请查看 logs/ 目录下的日志文件了解服务器运行状态。%COLOR_RESET%
echo %COLOR_INFO%按任意键返回菜单。%COLOR_RESET%
pause
goto :menu

:db_migrate
echo %COLOR_INFO%正在执行数据库迁移...%COLOR_RESET%
go run main.go -m db -t migrate -settings settings.yaml
if %errorlevel% neq 0 (
echo %COLOR_ERROR%数据库迁移失败，请检查数据库连接和配置。%COLOR_RESET%
) else (
echo %COLOR_SUCCESS%数据库迁移成功！%COLOR_RESET%
)
pause
goto :menu

:create_admin
echo %COLOR_INFO%正在创建管理员用户...%COLOR_RESET%
set /p username=请输入管理员用户名 (默认: admin): 
if "%username%" == "" set username=admin

setlocal enabledelayedexpansion
set "password="
echo %COLOR_INFO%请输入管理员密码 (默认: admin123): %COLOR_RESET%
echo 注意：密码将不可见
echo. 

REM 隐藏密码输入
for /f "tokens=* delims=" %%p in ('powershell -command "$p=Read-Host -AsSecureString;$b=ConvertFrom-SecureString -SecureString $p -AsPlainText;$b"') do set "password=%%p"

if "!password!" == "" set password=admin123
endlocal & set "password=%password%"

REM 执行创建管理员用户的命令
go run main.go -m user -t create -username "%username%" -password "%password%" -settings settings.yaml
if %errorlevel% neq 0 (
echo %COLOR_ERROR%创建管理员用户失败。%COLOR_RESET%
) else (
echo %COLOR_SUCCESS%管理员用户创建成功！用户名: %username%, 密码: %password%%COLOR_RESET%
)
pause
goto :menu

:exit
echo %COLOR_INFO%感谢使用RBAC管理员服务器！%COLOR_RESET%
pause