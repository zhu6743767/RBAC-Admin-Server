@echo off
title RBAC Admin Server - GitHub 上传工具
color 0A

echo.
echo ======================================
echo RBAC Admin Server - GitHub 创建和上传
echo ======================================
echo.

:MENU
echo 请选择操作：
echo.
echo 1. 手动创建仓库后推送（推荐）
echo 2. 仅查看当前Git状态
echo 3. 查看详细指南
echo 4. 退出
echo.
set /p choice=请输入选项(1-4): 

if "%choice%"=="1" goto PUSH
if "%choice%"=="2" goto STATUS
if "%choice%"=="3" goto GUIDE
if "%choice%"=="4" goto EXIT

goto MENU

:PUSH
echo.
echo --------------------------------------
echo 步骤1：创建GitHub仓库
echo --------------------------------------
echo.
echo 请先访问：https://github.com/new
echo 创建名为 "RBAC-Admin-Server" 的仓库
echo 不要勾选任何初始化选项！
echo.
set /p confirmed=已创建仓库？(按回车继续) 

:GET_USERNAME
echo.
echo --------------------------------------
echo 步骤2：配置推送
echo --------------------------------------
echo.
set /p username=请输入GitHub用户名: 
if "%username%"=="" goto GET_USERNAME

echo.
echo 正在配置远程仓库...
git remote set-url origin https://github.com/%username%/RBAC-Admin-Server.git

echo.
echo 正在推送代码...
git push -u origin master

if %errorlevel% neq 0 (
    echo.
    echo 推送失败！可能的原因：
    echo 1. 仓库未创建
    echo 2. 用户名错误
    echo 3. 需要登录认证
    echo.
    echo 请检查CREATE_AND_PUSH.md获取解决方案
    pause
    goto MENU
)

echo.
echo ======================================
echo 上传成功！
echo ======================================
echo 访问：https://github.com/%username%/RBAC-Admin-Server
echo.
pause
goto EXIT

:STATUS
echo.
echo ======================================
echo 当前Git状态
echo ======================================
echo.
echo 当前目录：
cd
echo.
echo Git状态：
git status
echo.
echo 远程配置：
git remote -v
echo.
echo 提交历史：
git log --oneline -5
echo.
pause
goto MENU

:GUIDE
echo.
echo 正在打开详细指南...
if exist "CREATE_AND_PUSH.md" (
    start notepad "CREATE_AND_PUSH.md"
) else (
    echo 指南文件不存在！
)
pause
goto MENU

:EXIT
echo.
echo 感谢使用！按任意键退出...
pause