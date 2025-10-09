@echo off

REM 检查服务器是否已经在运行
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8080 ^| findstr 127.0.0.1') do (
    echo 发现服务器进程正在运行，PID: %%a
    taskkill /PID %%a /F
    echo 已终止进程
)

REM 清理旧的数据库文件（以防SQLite有问题）
if exist rbac_admin.db (
    del rbac_admin.db
    echo 已删除旧的数据库文件
)

REM 设置环境变量
set DB_PASSWORD=123456

echo 正在启动服务器...
start "RBAC Admin Server" /B powershell -Command "go run main.go -m=server > server_log.txt 2>&1"

echo 等待3秒让服务器启动...
ping -n 4 127.0.0.1 > nul

echo 检查服务器进程...
tasklist /FI "IMAGENAME eq go.exe" /FI "WINDOWTITLE eq RBAC Admin Server"

echo 检查端口占用情况...
netstat -ano | findstr :8080

echo 尝试访问API...
curl -v http://localhost:8080/captcha/get || echo API访问失败

echo.
echo 服务器日志预览（最后10行）：
type server_log.txt | tail -n 10

echo.
echo 启动完成。请查看server_log.txt获取完整日志。