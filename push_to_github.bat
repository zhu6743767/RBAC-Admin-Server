@echo off
echo.
echo ======================================
echo RBAC Admin Server - GitHub 上传工具
echo ======================================
echo.

:INPUT_USERNAME
echo 请输入你的GitHub用户名：
set /p GITHUB_USERNAME=用户名: 

if "%GITHUB_USERNAME%"=="" (
    echo 用户名不能为空！
    goto INPUT_USERNAME
)

:CONFIRM
echo.
echo 即将上传项目到以下地址：
echo https://github.com/%GITHUB_USERNAME%/RBAC-Admin-Server
echo.
echo 请确认：
echo 1. 已在GitHub创建同名仓库
echo 2. 已检查敏感信息已排除
echo 3. 准备上传
set /p CONFIRM=继续？(y/n): 

if /i "%CONFIRM%"=="y" goto PUSH
if /i "%CONFIRM%"=="yes" goto PUSH
echo 已取消上传。
pause
exit

:PUSH
echo.
echo 正在配置远程仓库...
git remote set-url origin https://github.com/%GITHUB_USERNAME%/RBAC-Admin-Server.git

echo.
echo 正在推送代码...
git push -u origin master

echo.
echo 上传完成！
echo 访问：https://github.com/%GITHUB_USERNAME%/RBAC-Admin-Server
echo.
pause