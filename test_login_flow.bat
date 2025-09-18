@echo off
echo 🧪 RBAC登录测试流程
echo ======================
echo.

REM 步骤1: 创建管理员用户
echo 📋 步骤1: 创建管理员用户...
if exist create_admin_simple.exe (
    create_admin_simple.exe
) else (
    echo 🔨 编译管理员创建程序...
    go build -o create_admin_simple.exe create_admin_simple.go
    if %errorlevel% equ 0 (
        create_admin_simple.exe
    ) else (
        echo ❌ 编译失败，使用go run...
        go run create_admin_simple.go
    )
)
echo.

REM 步骤2: 检查数据库
echo 📋 步骤2: 检查数据库中的用户...
if exist check_user_existing.exe (
    check_user_existing.exe
) else (
    echo 🔨 编译用户检查程序...
    go build -o check_user_existing.exe check_user_existing.go
    if %errorlevel% equ 0 (
        check_user_existing.exe
    ) else (
        echo ❌ 编译失败，使用go run...
        go run check_user_existing.go
    )
)
echo.

REM 步骤3: 启动服务器
echo 📋 步骤3: 启动服务器...
echo 🚀 启动RBAC服务器...
if exist rbac.admin.exe (
    start "" rbac.admin.exe -env dev
) else (
    start "" cmd /c "go run main.go -env dev"
)
echo ⏳ 等待服务器启动...
timeout /t 5 /nobreak > nul
echo.

REM 步骤4: 测试登录
echo 📋 步骤4: 测试登录...
echo 🧪 运行登录测试...
go run simple_login.go
echo.

REM 步骤5: 额外测试
echo 📋 步骤5: 运行详细登录调试...
go run debug_login.go
echo.

echo ✅ 测试流程完成！
echo 💡 提示: 服务器仍在后台运行，可以按需要测试其他功能
echo 🛑 如需停止服务器，请关闭相应的命令窗口
echo.
pause