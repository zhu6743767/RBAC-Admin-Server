#!/bin/bash

# =================================================================================
# 🚀 RBAC管理员服务器 - 部署脚本
# =================================================================================
# 📋 脚本功能：
#   1. 环境检查
#   2. 敏感信息过滤检查
#   3. Git初始化（可选）
#   4. 构建和运行准备
# =================================================================================

set -e  # 遇到错误立即退出

echo "🚀 RBAC管理员服务器部署脚本"
echo "=========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函数：打印彩色信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 1. 环境检查
check_environment() {
    print_info "检查环境..."
    
    # 检查Go环境
    if ! command -v go &> /dev/null; then
        print_error "未找到Go环境，请先安装Go 1.19+"
        exit 1
    fi
    
    # 检查Go版本
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go版本: $GO_VERSION"
    
    # 检查Git
    if ! command -v git &> /dev/null; then
        print_warning "未找到Git，Git功能将不可用"
    else
        print_success "Git已安装"
    fi
}

# 2. 敏感信息检查
check_sensitive_info() {
    print_info "检查敏感信息..."
    
    # 检查.env文件
    if [ -f ".env" ]; then
        print_warning "发现 .env 文件，请确保不包含真实密码"
        print_info "建议：使用 .env.example 作为模板创建生产环境配置"
    fi
    
    # 检查settings.yaml中的敏感信息
    if [ -f "settings.yaml" ]; then
        if grep -q "192.168.\|Zdj_7819!\|localhost" settings.yaml; then
            print_warning "settings.yaml 中可能包含敏感信息或测试数据"
            print_info "建议：使用环境变量配置敏感信息"
        fi
    fi
    
    # 检查是否已配置.gitignore
    if [ -f ".gitignore" ]; then
        if grep -q "\.env\|\.key\|\.pem" .gitignore; then
            print_success ".gitignore 已配置敏感文件过滤"
        else
            print_warning ".gitignore 可能未正确配置敏感文件过滤"
        fi
    fi
    
    print_success "敏感信息检查完成"
}

# 3. 依赖下载
install_dependencies() {
    print_info "下载依赖..."
    
    go mod download
    go mod tidy
    
    print_success "依赖下载完成"
}

# 4. 配置验证
validate_config() {
    print_info "验证配置..."
    
    # 验证开发环境配置
    if go run main.go -env dev -dry-run 2>/dev/null; then
        print_success "开发环境配置验证通过"
    else
        print_warning "开发环境配置验证失败（可能需要数据库）"
    fi
    
    # 创建简单的配置验证脚本
    cat > validate_config_simple.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
    
    "rbac_admin_server/config"
)

func main() {
    env := "dev"
    if len(os.Args) > 1 {
        env = os.Args[1]
    }
    
    fmt.Printf("验证 %s 环境配置...\n", env)
    
    cfg, err := config.LoadConfig(env, "")
    if err != nil {
        log.Printf("配置加载失败: %v", err)
        os.Exit(1)
    }
    
    fmt.Printf("✅ 配置验证成功！\n")
    fmt.Printf("系统配置: IP=%s, Port=%d\n", cfg.System.IP, cfg.System.Port)
    fmt.Printf("数据库: %s@%s:%d/%s\n", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)
    fmt.Printf("JWT: Issuer=%s, Audience=%s\n", cfg.JWT.Issuer, cfg.JWT.Audience)
}
EOF
    
    if go run validate_config_simple.go dev 2>/dev/null; then
        print_success "配置验证通过"
    else
        print_warning "配置验证失败，请检查配置文件"
    fi
    
    # 清理验证脚本
    rm -f validate_config_simple.go
}

# 5. 构建测试
build_test() {
    print_info "构建测试..."
    
    if go build -o rbac_admin_server_test .; then
        print_success "构建成功"
        rm -f rbac_admin_server_test
    else
        print_error "构建失败，请检查代码"
        exit 1
    fi
}

# 6. Git初始化（可选）
init_git() {
    if [ -d ".git" ]; then
        print_info "Git仓库已存在"
        return
    fi
    
    print_info "初始化Git仓库..."
    
    if command -v git &> /dev/null; then
        git init
        git add .gitignore SECURITY.md README.md
        git add settings_dev.yaml settings_test.yaml settings_prod.yaml
        git add go.mod go.sum
        git add config/ core/ middleware/ models/ routes/ api/
        git commit -m "Initial commit: RBAC admin server"
        
        print_success "Git仓库初始化完成"
        print_info "可以添加远程仓库: git remote add origin <your-repo-url>"
    else
        print_warning "Git未安装，跳过Git初始化"
    fi
}

# 7. 显示后续步骤
show_next_steps() {
    print_success "部署准备完成！"
    echo ""
    echo "📋 后续步骤："
    echo ""
    echo "🏃‍♂️ 开发环境（推荐）："
    echo "  go run main.go -env dev"
    echo ""
    echo "🔧 测试环境："
    echo "  1. 确保MySQL数据库已安装并运行"
    echo "  2. 创建数据库: CREATE DATABASE rbac_admin_test;"
    echo "  3. go run main.go -env test"
    echo ""
    echo "🚀 生产环境："
    echo "  1. 设置环境变量（参考 .env.example）"
    echo "  2. 确保MySQL和Redis已安装"
    echo "  3. 创建数据库: CREATE DATABASE rbac_admin_prod;"
    echo "  4. go run main.go -env prod"
    echo ""
    echo "📖 文档："
    echo "  API文档: http://localhost:8080/swagger/index.html"
    echo "  安全配置: SECURITY.md"
    echo "  配置指南: CONFIG_GUIDE.md"
    echo ""
    print_info "开始您的RBAC管理之旅！🎉"
}

# 主函数
main() {
    check_environment
    check_sensitive_info
    install_dependencies
    validate_config
    build_test
    
    # 询问是否初始化Git
    if command -v git &> /dev/null && [ ! -d ".git" ]; then
        read -p "是否初始化Git仓库？(y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            init_git
        fi
    fi
    
    show_next_steps
}

# 运行主函数
main "$@"