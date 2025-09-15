@echo off
REM RBAC管理员服务器 - Windows日志系统部署脚本

echo 🚀 开始部署RBAC日志系统...

REM 检查Docker
where docker >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Docker未安装，请先安装Docker
    pause
    exit /b 1
)

REM 检查Docker Compose
where docker-compose >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Docker Compose未安装，请先安装Docker Compose
    pause
    exit /b 1
)

REM 创建目录结构
echo 📁 创建目录结构...
if not exist "logs\error" mkdir logs\error
if not exist "logs\operation" mkdir logs\operation
if not exist "logs\access" mkdir logs\access
if not exist "logstash\pipeline" mkdir logstash\pipeline
if not exist "logstash\config" mkdir logstash\config
if not exist "scripts" mkdir scripts

REM 启动ELK栈
echo 🐳 启动ELK栈...
docker-compose -f docker-compose.logging.yml up -d elasticsearch kibana

REM 等待Elasticsearch启动
echo ⏳ 等待Elasticsearch启动...
:wait_es
curl -s http://localhost:9200/_cluster/health >nul 2>nul
if %errorlevel% neq 0 (
    timeout /t 5 /nobreak >nul
    goto wait_es
)

echo ✅ Elasticsearch已启动

REM 启动Filebeat
echo 📊 启动Filebeat...
docker-compose -f docker-compose.logging.yml up -d filebeat

REM 验证服务状态
echo 🔍 验证服务状态...
docker-compose -f docker-compose.logging.yml ps

REM 输出访问信息
echo 🎉 日志系统部署完成！
echo.
echo 📊 Kibana访问地址: http://localhost:5601
echo 📈 Elasticsearch API: http://localhost:9200
echo 📁 日志目录: .\logs\
echo.
echo 🔧 常用命令:
echo   查看日志: docker-compose -f docker-compose.logging.yml logs -f
echo   停止服务: docker-compose -f docker-compose.logging.yml down
echo   重启服务: docker-compose -f docker-compose.logging.yml restart
echo.
echo 📖 下一步:
echo   1. 启动RBAC应用: go run main.go
echo   2. 访问Kibana创建索引模式: rbac-logs-*
echo   3. 在Discover中查看日志
echo.
pause