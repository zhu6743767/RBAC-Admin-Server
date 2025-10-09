@echo off

REM 创建管理员用户
set username=admin
set password=admin123

echo 创建管理员用户: %username%
echo 密码: %password%

go run main.go -m=user -t=create -username=%username% -password=%password%

if %errorlevel% equ 0 (
    echo.
echo ✅ 管理员用户创建成功!
echo 您现在可以尝试使用这些凭据登录系统。
) else (
    echo.
echo ❌ 创建管理员用户失败!
echo 请检查错误信息并重试。
)

pause