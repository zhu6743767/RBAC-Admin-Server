@echo off
REM =================================================================================
REM 🚀 RBAC管理员服务器 - Windows部署脚本
REM =================================================================================
REM 📋 脚本功能：
REM   1. 环境检查
REM   2. 敏感信息过滤检查
REM   3. 依赖下载
REM   4. 配置验证
REM   5. 构建测试
REM =================================================================================

echo 🚀 RBAC管理员服务器部署脚本
echo ==========================================

REM 设置变量
setlocal enabledelayedexpansion

REM 颜色定义（Windows 10+支持）
set "COLOR_INFO=[94m"
set "COLOR_SUCCESS=[92m"
set "COLOR_WARNING=[93m"
set "COLOR_ERROR=[91m"
set "COLOR_RESET=[0m"

REM 函数：打印信息
goto :main

:print_info
echo %COLOR_INFO%[INFO]%COLOR_RESET% %~1%
goto :eof

:print_success
echo %COLOR_SUCCESS%[SUCCESS]%COLOR_RESET% %~1%
goto :eof

:print_warning
echo %COLOR_WARNING%[WARNING]%COLOR_RESET% %~1%
goto :eof

:print_error
echo %COLOR_ERROR%[ERROR]%COLOR_RESET% %~1%
goto :eof

REM 1. 环境检查
:check_environment
call :print_info "检查环境..."

REM 检查Go环境
where go >nul 2>nul
if %errorlevel% neq 0 (
    call :print_error "未找到Go环境，请先安装Go 1.19+"
    exit /b 1
)

REM 检查Go版本
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
call :print_success "Go版本: %GO_VERSION%"

REM 检查Git
where git >nul 2>nul
if %errorlevel% equ 0 (
    call :print_success "Git已安装"
) else (
    call :print_warning "未找到Git，Git功能将不可用"
)
goto :eof

REM 2. 敏感信息检查
:check_sensitive_info
call :print_info "检查敏感信息..."

REM 检查.env文件
if exist ".env" (
    call :print_warning "发现 .env 文件，请确保不包含真实密码"
    call :print_info "建议：使用 .env.example 作为模板创建生产环境配置"
)

REM 检查settings.yaml中的敏感信息
if exist "settings.yaml" (
    findstr /c:"192.168." /c:"Zdj_7819!" /c:"localhost" settings.yaml >nul 2>nul
    if %errorlevel% equ 0 (
        call :print_warning "settings.yaml 中可能包含敏感信息或测试数据"
        call :print_info "建议：使用环境变量配置敏感信息"
    )
)

REM 检查是否已配置.gitignore
if exist ".gitignore" (
    findstr /c:".env" /c:".key" /c:".pem" .gitignore >nul 2>nul
    if %errorlevel% equ 0 (
        call :print_success ".gitignore 已配置敏感文件过滤"
    ) else (
        call :print_warning ".gitignore 可能未正确配置敏感文件过滤"
    )
)

call :print_success "敏感信息检查完成"
goto :eof

REM 3. 依赖下载
:install_dependencies
call :print_info "下载依赖..."

go mod download
if %errorlevel% neq 0 (
    call :print_error "依赖下载失败"
    exit /b 1
)

go mod tidy
if %errorlevel% neq 0 (
    call :print_error "依赖整理失败"
    exit /b 1
)

call :print_success "依赖下载完成"
goto :eof

REM 4. 配置验证
:validate_config
call :print_info "验证配置..."

REM 创建简单的配置验证脚本
echo package main > validate_config_simple.go
echo. >> validate_config_simple.go
echo import ^( >> validate_config_simple.go
echo     "fmt" >> validate_config_simple.go
echo     "log" >> validate_config_simple.go
echo     "os" >> validate_config_simple.go     >> validate_config_simple.go
echo. >> validate_config_simple.go
echo     "rbac_admin_server/config" >> validate_config_simple.go
echo ^) >> validate_config_simple.go
echo. >> validate_config_simple.go
echo func main^(^) ^{ >> validate_config_simple.go
echo     env := "dev" >> validate_config_simple.go
echo     if len^(os.Args^) ^> 1 { >> validate_config_simple.go
echo         env = os.Args[1] >> validate_config_simple.go
echo     } >> validate_config_simple.go
echo. >> validate_config_simple.go
echo     fmt.Printf^("验证 %%s 环境配置...\n", env^) >> validate_config_simple.go
echo. >> validate_config_simple.go
echo     cfg, err := config.LoadConfig^(env, ""^) >> validate_config_simple.go
echo     if err != nil { >> validate_config_simple.go
echo         log.Printf^("配置加载失败: %%v", err^) >> validate_config_simple.go
echo         os.Exit^(1^) >> validate_config_simple.go
echo     } >> validate_config_simple.go
echo. >> validate_config_simple.go
echo     fmt.Printf^("✅ 配置验证成功！\n"^) >> validate_config_simple.go
echo     fmt.Printf^("系统配置: IP=%%s, Port=%%d\n", cfg.System.IP, cfg.System.Port^) >> validate_config_simple.go
echo     fmt.Printf^("数据库: %%s@%%s:%%d/%%s\n", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName^) >> validate_config_simple.go
echo     fmt.Printf^("JWT: Issuer=%%s, Audience=%%s\n", cfg.JWT.Issuer, cfg.JWT.Audience^) >> validate_config_simple.go
echo ^} >> validate_config_simple.go

if exist validate_config_simple.go (
    go run validate_config_simple.go dev >nul 2>nul
    if %errorlevel% equ 0 (
        call :print_success "配置验证通过"
    ) else (
        call :print_warning "配置验证失败，请检查配置文件"
    )
    del validate_config_simple.go >nul 2>nul
)
goto :eof

REM 5. 构建测试
:build_test
call :print_info "构建测试..."

if exist rbac_admin_server_test.exe del rbac_admin_server_test.exe >nul 2>nul

go build -o rbac_admin_server_test.exe .
if %errorlevel% equ 0 (
    call :print_success "构建成功"
    del rbac_admin_server_test.exe >nul 2>nul
) else (
    call :print_error "构建失败，请检查代码"
    exit /b 1
)
goto :eof

REM 6. 显示后续步骤
:show_next_steps
call :print_success "部署准备完成！"
echo.
echo 📋 后续步骤：
echo.
echo 🏃‍♂️ 开发环境（推荐）：
echo   go run main.go -env dev
echo.
echo 🔧 测试环境：
echo   1. 确保MySQL数据库已安装并运行
echo   2. 创建数据库: CREATE DATABASE rbac_admin_test;
echo   3. go run main.go -env test
echo.
echo 🚀 生产环境：
echo   1. 设置环境变量（参考 .env.example）
echo   2. 确保MySQL和Redis已安装
echo   3. 创建数据库: CREATE DATABASE rbac_admin_prod;
echo   4. go run main.go -env prod
echo.
echo 📖 文档：
echo   API文档: http://localhost:8080/swagger/index.html
echo   安全配置: SECURITY.md
echo   配置指南: CONFIG_GUIDE.md
echo.
call :print_info "开始您的RBAC管理之旅！🎉"
goto :eof

REM 主函数
:main
call :check_environment
call :check_sensitive_info
call :install_dependencies
call :validate_config
call :build_test
call :show_next_steps

endlocal
echo.
echo 按任意键退出...
pause >nul