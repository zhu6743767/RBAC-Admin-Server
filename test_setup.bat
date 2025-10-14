@echo off
REM RBAC管理员服务器环境测试脚本
REM =================================================================================
REM 此脚本用于测试RBAC管理员服务器的基本环境配置和功能
REM =================================================================================

REM 设置中文编码
chcp 65001 > nul

REM 定义颜色变量
set "COLOR_SUCCESS=[92m"
set "COLOR_ERROR=[91m"
set "COLOR_INFO=[96m"
set "COLOR_RESET=[0m"

cls
echo %COLOR_INFO%================================================================================%COLOR_RESET%
echo %COLOR_INFO%                     RBAC管理员服务器环境测试                     %COLOR_RESET%
echo %COLOR_INFO%================================================================================%COLOR_RESET%

REM 1. 检查Go环境
echo. 
echo %COLOR_INFO%[1/6] 检查Go环境...%COLOR_RESET%
go version
if %errorlevel% neq 0 (
echo %COLOR_ERROR%❌ Go环境未正确安装，请先安装Go并配置环境变量。%COLOR_RESET%
pause
exit /b 1
) else (
echo %COLOR_SUCCESS%✅ Go环境检查通过！%COLOR_RESET%
)

REM 2. 检查项目依赖
echo. 
echo %COLOR_INFO%[2/6] 安装项目依赖...%COLOR_RESET%
go mod tidy
if %errorlevel% neq 0 (
echo %COLOR_ERROR%❌ 依赖安装失败，请检查网络连接和go.mod文件。%COLOR_RESET%
pause
exit /b 1
) else (
echo %COLOR_SUCCESS%✅ 依赖安装成功！%COLOR_RESET%
)

REM 3. 检查数据库连接
echo. 
echo %COLOR_INFO%[3/6] 尝试连接数据库...%COLOR_RESET%
echo 创建临时测试文件...
echo package main ^
import ( ^
"database/sql" ^
"fmt" ^
"time" ^
_ "github.com/go-sql-driver/mysql" ^
) ^
func main() { ^
    dsn := "root:Zdj_7819!@tcp(192.168.10.199:3306)/rbacadmin?parseTime=true" ^
    db, err := sql.Open("mysql", dsn) ^
    if err != nil { ^
        fmt.Println("数据库连接失败:", err) ^
        return ^
    } ^
    defer db.Close() ^
    
    // 设置连接参数
    db.SetConnMaxLifetime(time.Minute * 3) ^
    db.SetMaxOpenConns(10) ^
    db.SetMaxIdleConns(10) ^
    
    // 测试连接
    if err := db.Ping(); err != nil { ^
        fmt.Println("数据库连接测试失败:", err) ^
        return ^
    } ^
    
    fmt.Println("数据库连接成功！") ^
} > test_db_conn.go

go run test_db_conn.go
if %errorlevel% neq 0 (
echo %COLOR_ERROR%❌ 数据库连接失败，请检查数据库配置和网络连接。%COLOR_RESET%
del test_db_conn.go
pause
exit /b 1
) else (
echo %COLOR_SUCCESS%✅ 数据库连接测试通过！%COLOR_RESET%
del test_db_conn.go
)

REM 4. 检查Redis连接
echo. 
echo %COLOR_INFO%[4/6] 尝试连接Redis...%COLOR_RESET%
echo 创建临时测试文件...
echo package main ^
import ( ^
"fmt" ^
"github.com/go-redis/redis/v8" ^
"context" ^
) ^
func main() { ^
    ctx := context.Background() ^
    client := redis.NewClient(&redis.Options{ ^
        Addr:     "192.168.10.199:6379", ^
        Password: "", // 无密码
        DB:       4,  // 使用DB 4
    }) ^
    
    // 测试连接
    _, err := client.Ping(ctx).Result() ^
    if err != nil { ^
        fmt.Println("Redis连接失败:", err) ^
        return ^
    } ^
    
    fmt.Println("Redis连接成功！") ^
} > test_redis_conn.go

go run test_redis_conn.go
if %errorlevel% neq 0 (
echo %COLOR_ERROR%❌ Redis连接失败，请检查Redis配置和网络连接。%COLOR_RESET%
del test_redis_conn.go
pause
exit /b 1
) else (
echo %COLOR_SUCCESS%✅ Redis连接测试通过！%COLOR_RESET%
del test_redis_conn.go
)

REM 5. 执行数据库迁移
echo. 
echo %COLOR_INFO%[5/6] 执行数据库迁移...%COLOR_RESET%
go run main.go -m db -t migrate -settings settings.yaml
if %errorlevel% neq 0 (
echo %COLOR_ERROR%❌ 数据库迁移失败，请检查错误信息。%COLOR_RESET%
pause
exit /b 1
) else (
echo %COLOR_SUCCESS%✅ 数据库迁移成功！%COLOR_RESET%
)

REM 6. 创建管理员用户
echo. 
echo %COLOR_INFO%[6/6] 创建管理员用户...%COLOR_RESET%
go run main.go -m user -t create -username admin -password admin123 -settings settings.yaml
if %errorlevel% neq 0 (
echo %COLOR_ERROR%❌ 管理员用户创建失败，请检查错误信息。%COLOR_RESET%
pause
exit /b 1
) else (
echo %COLOR_SUCCESS%✅ 管理员用户创建成功！用户名: admin, 密码: admin123%COLOR_RESET%
)

REM 测试完成
echo. 
echo %COLOR_INFO%================================================================================%COLOR_RESET%
echo %COLOR_SUCCESS%🎉 RBAC管理员服务器环境测试全部通过！🎉%COLOR_RESET%
echo %COLOR_INFO%您现在可以使用 run_server.bat 启动服务器了。%COLOR_RESET%
echo %COLOR_INFO%================================================================================%COLOR_RESET%

pause