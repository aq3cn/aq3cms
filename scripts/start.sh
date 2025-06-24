#!/bin/bash

# aq3cms 启动脚本
# 用于快速启动 aq3cms 系统

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 命令未找到，请先安装 $1"
        exit 1
    fi
}

# 检查文件是否存在
check_file() {
    if [ ! -f "$1" ]; then
        log_error "文件 $1 不存在"
        exit 1
    fi
}

# 检查目录是否存在
check_directory() {
    if [ ! -d "$1" ]; then
        log_warn "目录 $1 不存在，正在创建..."
        mkdir -p "$1"
    fi
}

# 显示帮助信息
show_help() {
    echo "aq3cms 启动脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -d, --dev           启动开发环境"
    echo "  -p, --prod          启动生产环境"
    echo "  -b, --build         构建应用"
    echo "  -t, --test          运行测试"
    echo "  -c, --clean         清理构建文件"
    echo "  --docker            使用 Docker 启动"
    echo "  --docker-dev        使用 Docker 启动开发环境"
    echo "  --stop              停止服务"
    echo "  --restart           重启服务"
    echo "  --logs              查看日志"
    echo "  --health            检查健康状态"
    echo ""
    echo "示例:"
    echo "  $0 --dev            # 启动开发环境"
    echo "  $0 --docker         # 使用 Docker 启动生产环境"
    echo "  $0 --docker-dev     # 使用 Docker 启动开发环境"
    echo "  $0 --build          # 构建应用"
    echo "  $0 --test           # 运行测试"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    if [ "$USE_DOCKER" = "true" ]; then
        check_command "docker"
        check_command "docker-compose"
    else
        check_command "go"
        check_command "make"
    fi
    
    log_info "依赖检查完成"
}

# 初始化环境
init_environment() {
    log_info "初始化环境..."
    
    # 创建必要的目录
    check_directory "logs"
    check_directory "uploads"
    check_directory "data"
    
    # 检查配置文件
    if [ ! -f "config.yaml" ]; then
        if [ -f "config.yaml.example" ]; then
            log_warn "config.yaml 不存在，正在复制示例配置文件..."
            cp config.yaml.example config.yaml
        else
            log_error "配置文件不存在，请先创建 config.yaml"
            exit 1
        fi
    fi
    
    # 检查环境变量文件
    if [ ! -f ".env" ] && [ -f ".env.example" ]; then
        log_warn ".env 文件不存在，正在复制示例文件..."
        cp .env.example .env
    fi
    
    log_info "环境初始化完成"
}

# 构建应用
build_app() {
    log_info "构建应用..."
    make build
    log_info "构建完成"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    make test
    log_info "测试完成"
}

# 清理构建文件
clean_build() {
    log_info "清理构建文件..."
    make clean
    log_info "清理完成"
}

# 启动开发环境
start_dev() {
    log_info "启动开发环境..."
    
    if [ "$USE_DOCKER" = "true" ]; then
        make dev
    else
        # 检查 Air 是否安装
        if ! command -v air &> /dev/null; then
            log_warn "Air 未安装，正在安装..."
            go install github.com/cosmtrek/air@latest
        fi
        
        # 启动热重载
        air -c .air.toml
    fi
}

# 启动生产环境
start_prod() {
    log_info "启动生产环境..."
    
    if [ "$USE_DOCKER" = "true" ]; then
        make deploy
    else
        # 构建应用
        build_app
        
        # 启动应用
        ./bin/aq3cms
    fi
}

# 停止服务
stop_services() {
    log_info "停止服务..."
    
    if [ "$USE_DOCKER" = "true" ]; then
        if [ "$DEV_MODE" = "true" ]; then
            make dev-down
        else
            make deploy-down
        fi
    else
        # 查找并停止 aq3cms 进程
        pkill -f "aq3cms" || true
        pkill -f "air" || true
    fi
    
    log_info "服务已停止"
}

# 重启服务
restart_services() {
    log_info "重启服务..."
    stop_services
    sleep 2
    
    if [ "$DEV_MODE" = "true" ]; then
        start_dev
    else
        start_prod
    fi
}

# 查看日志
show_logs() {
    log_info "查看日志..."
    
    if [ "$USE_DOCKER" = "true" ]; then
        if [ "$DEV_MODE" = "true" ]; then
            make dev-logs
        else
            make deploy-logs
        fi
    else
        if [ -f "logs/app.log" ]; then
            tail -f logs/app.log
        else
            log_warn "日志文件不存在"
        fi
    fi
}

# 健康检查
health_check() {
    log_info "检查健康状态..."
    
    # 等待服务启动
    sleep 5
    
    # 检查健康端点
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_info "服务健康状态正常"
    else
        log_error "服务健康检查失败"
        exit 1
    fi
}

# 主函数
main() {
    # 默认值
    DEV_MODE=false
    USE_DOCKER=false
    ACTION=""
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -d|--dev)
                DEV_MODE=true
                ACTION="start"
                shift
                ;;
            -p|--prod)
                DEV_MODE=false
                ACTION="start"
                shift
                ;;
            -b|--build)
                ACTION="build"
                shift
                ;;
            -t|--test)
                ACTION="test"
                shift
                ;;
            -c|--clean)
                ACTION="clean"
                shift
                ;;
            --docker)
                USE_DOCKER=true
                DEV_MODE=false
                ACTION="start"
                shift
                ;;
            --docker-dev)
                USE_DOCKER=true
                DEV_MODE=true
                ACTION="start"
                shift
                ;;
            --stop)
                ACTION="stop"
                shift
                ;;
            --restart)
                ACTION="restart"
                shift
                ;;
            --logs)
                ACTION="logs"
                shift
                ;;
            --health)
                ACTION="health"
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定动作，显示帮助
    if [ -z "$ACTION" ]; then
        show_help
        exit 0
    fi
    
    # 检查依赖
    check_dependencies
    
    # 初始化环境
    if [ "$ACTION" != "stop" ] && [ "$ACTION" != "logs" ] && [ "$ACTION" != "health" ]; then
        init_environment
    fi
    
    # 执行动作
    case $ACTION in
        start)
            if [ "$DEV_MODE" = "true" ]; then
                start_dev
            else
                start_prod
            fi
            ;;
        build)
            build_app
            ;;
        test)
            run_tests
            ;;
        clean)
            clean_build
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            show_logs
            ;;
        health)
            health_check
            ;;
    esac
}

# 运行主函数
main "$@"
