@echo off
echo === 测试所有示例程序 ===
echo.

echo 1. 编译运行程序...
go build -o run_examples.exe run_examples.go
if %errorlevel% neq 0 (
    echo 编译失败！
    pause
    exit /b 1
)
echo 编译成功！
echo.

echo 2. 编译所有示例文件...
go build -v ./examples/...
if %errorlevel% neq 0 (
    echo 示例文件编译失败！
    pause
    exit /b 1
)
echo 所有示例文件编译成功！
echo.

echo 3. 运行数据库示例...
go run run_examples.go
if %errorlevel% neq 0 (
    echo 运行失败！
    pause
    exit /b 1
)

echo.
echo === 所有测试通过！ ===
pause