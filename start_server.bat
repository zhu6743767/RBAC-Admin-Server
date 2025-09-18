@echo off
echo 🚀 启动RBAC管理服务器...

REM 检查是否有编译好的可执行文件
if exist rbac.admin.exe (
    echo 📦 使用已编译的rbac.admin.exe
    start "" rbac.admin.exe -env dev
) else (
    echo 🔨 使用go run启动服务器
    start "" cmd /c "go run main.go -env dev"
)

echo ⏳ 等待服务器启动...
timeout /t 3 /nobreak > nul
echo ✅ 服务器启动命令已发送
echo.
echo 🌐 服务器地址: http://127.0.0.1:8080
echo 📝 API文档: http://127.0.0.1:8080/swagger/index.html
echo.
echo 💡 提示: 可以运行 simple_login.go 测试登录功能
echo.
pause