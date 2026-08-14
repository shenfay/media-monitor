#!/bin/bash

# 开发环境启动脚本（本地模式，不使用 Docker）
# 
# 用法:
#   ./start.sh [选项]
#
# 选项:
#   -m, --monitor      同时启动 Asynqmon 监控
#   -s, --swagger      同时启动 Swagger UI
#   -a, --all          启动所有服务（包括监控）
#   -c, --clean        清理后台进程
#   -h, --help         显示帮助信息

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

show_help() {
    cat << EOF
用法：$0 [选项]

本地开发环境启动脚本（不使用 Docker，依赖本地 PostgreSQL 和 Redis）

选项:
  -m, --monitor      同时启动 Asynqmon 监控
  -s, --swagger      同时启动 Swagger UI
  -a, --all          启动所有服务（包括监控）
  -c, --clean        清理后台进程
  -h, --help         显示帮助信息

示例:
  ./start.sh                    # 只启动 API 和 Worker
  ./start.sh --monitor          # 启动 + Asynqmon
  ./start.sh --swagger          # 启动 + Swagger UI
  ./start.sh --all              # 启动所有服务
  ./start.sh --clean            # 清理后台进程

EOF
}

# 检查是否在 server 目录
check_server_dir() {
    # 获取脚本所在目录的绝对路径
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"
    SERVER_DIR="$PROJECT_ROOT/server"
    
    if [ ! -f "$SERVER_DIR/cmd/api/main.go" ]; then
        print_error "找不到 server 目录：$SERVER_DIR"
        print_error "请在项目根目录或 server 目录运行此脚本"
        exit 1
    fi
}

# 检查必要工具
check_dependencies() {
    local missing=()
    
    if ! command -v go &> /dev/null; then
        missing+=("go")
    fi
    
    if ! command -v psql &> /dev/null; then
        missing+=("psql")
    fi
    
    if ! command -v redis-cli &> /dev/null; then
        missing+=("redis-cli")
    fi
    
    if ! command -v curl &> /dev/null; then
        missing+=("curl")
    fi
    
    if ! command -v jq &> /dev/null; then
        missing+=("jq")
    fi
    
    if [ ${#missing[@]} -ne 0 ]; then
        print_error "缺少以下工具：${missing[*]}"
        print_info "请使用 brew install 安装"
        exit 1
    fi
}

# 检查 PostgreSQL 是否运行
check_postgresql() {
    print_info "检查 PostgreSQL 连接..."
    
    # 从 .env 文件读取配置
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"
    local env_file="$PROJECT_ROOT/server/configs/.env"
    if [ -f "$env_file" ]; then
        export $(grep -v '^#' "$env_file" | xargs)
    fi
    
    # 使用环境变量或默认值
    local db_host="${APP_DATABASE_HOST:-localhost}"
    local db_port="${APP_DATABASE_PORT:-5432}"
    local db_name="${APP_DATABASE_NAME:-kiqi}"
    local db_user="${APP_DATABASE_USER:-shenfay}"
    local db_password="${APP_DATABASE_PASSWORD:-postgres}"
    
    # 使用 PGPASSWORD 环境变量传递密码（避免交互式提示）
    export PGPASSWORD="$db_password"
    
    if psql -h "$db_host" -p "$db_port" -U "$db_user" -d "$db_name" -c "SELECT 1" &> /dev/null; then
        print_success "PostgreSQL 已就绪 (用户：$db_user, 数据库：$db_name)"
    else
        print_error "无法连接到 PostgreSQL"
        print_info "请确认:"
        echo "  1. PostgreSQL 服务已启动"
        echo "  2. 数据库 $db_name 已创建"
        echo "  3. 用户 $db_user 有访问权限"
        echo "  4. 密码正确 (当前使用：$db_password)"
        exit 1
    fi
    
    # 清理 PGPASSWORD
    unset PGPASSWORD
}

# 检查 Redis 是否运行
check_redis() {
    print_info "检查 Redis 连接..."
    
    if redis-cli ping &> /dev/null; then
        print_success "Redis 已就绪"
    else
        print_error "无法连接到 Redis"
        print_info "请先启动 Redis 服务"
        exit 1
    fi
}

# 停止现有进程
stop_existing() {
    print_info "停止现有进程..."
    
    # 停止 API 相关进程（包括 go run 和编译后的二进制）
    pkill -f "go run ./cmd/api" 2>/dev/null || true
    pkill -f "Caches/go-build.*api" 2>/dev/null || true
    
    # 停止 Worker 相关进程
    pkill -f "go run ./cmd/worker" 2>/dev/null || true
    pkill -f "Caches/go-build.*worker" 2>/dev/null || true
    
    # 停止其他辅助进程
    pkill -f "asynqmon" 2>/dev/null || true
    pkill -f "cmd/docs/main.go" 2>/dev/null || true
    
    sleep 1
    print_success "已清理现有进程"
}

# 启动 API 服务
start_api() {
    print_info "启动 API 服务..."
    
    # 获取 server 目录的绝对路径
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"
    SERVER_DIR="$PROJECT_ROOT/server"
    
    # 保存当前目录
    ORIGINAL_DIR="$(pwd)"
    
    # 切换到 server 目录
    cd "$SERVER_DIR" || exit 1
    
    # 在后台启动
    nohup go run ./cmd/api > /tmp/api.log 2>&1 &
    API_PID=$!
    
    # 等待 API 启动
    sleep 3
    
    # 检查是否成功启动
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "API 服务已启动 (PID: $API_PID)"
    else
        # 检查进程是否还在运行
        if ps -p $API_PID > /dev/null; then
            print_warning "API 正在启动中... (PID: $API_PID)"
            print_info "查看日志：tail -f /tmp/api.log"
        else
            print_error "API 启动失败"
            print_info "错误日志："
            tail -20 /tmp/api.log
            exit 1
        fi
    fi
    
    # 回到原始目录
    cd "$ORIGINAL_DIR"
}

# 启动 Worker 服务
start_worker() {
    print_info "启动 Worker 服务..."
    
    # 获取 server 目录的绝对路径
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"
    SERVER_DIR="$PROJECT_ROOT/server"
    
    # 保存当前目录
    ORIGINAL_DIR="$(pwd)"
    
    # 切换到 server 目录
    cd "$SERVER_DIR" || exit 1
    
    # 在后台启动
    nohup go run ./cmd/worker > /tmp/worker.log 2>&1 &
    WORKER_PID=$!
    
    # 等待启动
    sleep 2
    
    if ps -p $WORKER_PID > /dev/null; then
        print_success "Worker 服务已启动 (PID: $WORKER_PID)"
    else
        print_error "Worker 启动失败"
        tail -20 /tmp/worker.log
        exit 1
    fi
    
    # 回到原始目录
    cd "$ORIGINAL_DIR"
}

# 启动 Asynqmon
start_asynqmon() {
    print_info "启动 Asynqmon 监控..."
    
    # 获取 GOPATH
    GOPATH_BIN="$(go env GOPATH)/bin"
    
    # 检查是否已安装
    if [ ! -f "$GOPATH_BIN/asynqmon" ]; then
        print_warning "asynqmon 未安装，正在安装..."
        go install github.com/hibiken/asynqmon/cmd/asynqmon@latest
    fi
    
    # 确保 GOPATH/bin 在 PATH 中
    export PATH="$GOPATH_BIN:$PATH"
    
    # 再次检查 asynqmon 是否可用
    if ! command -v asynqmon &> /dev/null; then
        print_error "asynqmon 安装失败，请手动安装："
        echo "  go install github.com/hibiken/asynqmon/cmd/asynqmon@latest"
        exit 1
    fi
    
    # 在后台启动（使用 8081 端口，避免与 API 冲突）
    nohup asynqmon --redis-addr=localhost:6379 --port=8081 > /tmp/asynqmon.log 2>&1 &
    ASYNQMON_PID=$!
    
    sleep 2
    
    if ps -p $ASYNQMON_PID > /dev/null; then
        print_success "Asynqmon 已启动 (PID: $ASYNQMON_PID)"
        print_info "访问地址：http://localhost:8081"
    else
        print_error "Asynqmon 启动失败"
        tail -5 /tmp/asynqmon.log
        exit 1
    fi
}

# 启动 Swagger UI
start_swagger() {
    print_info "启动 Swagger UI..."
    
    # Swagger UI 已集成在 API 服务中，无需单独启动
    # 访问 http://localhost:8080/swagger/index.html 即可
    
    print_success "Swagger UI 已在 API 服务中运行"
    print_info "访问地址：http://localhost:8080/swagger/index.html"
}

# 保存 PID 到文件
save_pids() {
    local api_pid=$(pgrep -f "go run ./cmd/api")
    local worker_pid=$(pgrep -f "go run ./cmd/worker")
    local asynqmon_pid=$(pgrep -f "asynqmon")
    local swagger_pid=$(pgrep -f "cmd/docs/main.go")
    
    cat > /tmp/kiqi-pids.txt << EOF
API_PID=$api_pid
WORKER_PID=$worker_pid
ASYNQMON_PID=${asynqmon_pid:-}
SWAGGER_PID=${swagger_pid:-}
TIMESTAMP=$(date)
EOF
    
    print_info "进程 ID 已保存到 /tmp/kiqi-pids.txt"
}

# 主函数
main() {
    START_MONITOR=false
    START_SWAGGER=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--monitor)
                START_MONITOR=true
                shift
                ;;
            -s|--swagger)
                START_SWAGGER=true
                shift
                ;;
            -a|--all)
                START_MONITOR=true
                START_SWAGGER=true
                shift
                ;;
            -c|--clean)
                stop_existing
                print_success "所有进程已停止"
                exit 0
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "未知选项：$1"
                show_help
                exit 1
                ;;
        esac
    done
    
    echo ""
    echo "========================================"
    echo -e "${GREEN}🚀 启动本地开发环境${NC}"
    echo "========================================"
    echo ""
    
    # 检查
    check_server_dir
    check_dependencies
    
    # 停止现有进程
    stop_existing
    
    # 检查依赖服务
    check_postgresql
    check_redis
    
    echo ""
    echo "----------------------------------------"
    echo ""
    
    # 启动服务
    start_api
    start_worker
    
    if [ "$START_MONITOR" = true ]; then
        start_asynqmon
    fi
    
    if [ "$START_SWAGGER" = true ]; then
        start_swagger
    fi
    
    # 保存 PID
    save_pids
    
    echo ""
    echo "========================================"
    echo -e "${GREEN}✅ 所有服务已启动${NC}"
    echo "========================================"
    echo ""
    print_info "服务状态:"
    echo ""
    echo "  📡 API 服务:"
    echo "     - 健康检查：http://localhost:8080/health"
    echo "     - Swagger:   http://localhost:8080/swagger/index.html"
    echo "     - Metrics:   http://localhost:8080/metrics"
    echo ""
    
    if [ "$START_MONITOR" = true ]; then
        echo "  📊 Asynqmon:"
        echo "     - 监控面板：http://localhost:8080"
        echo ""
    fi
    
    if [ "$START_SWAGGER" = true ]; then
        echo "  📖 Swagger Docs:"
        echo "     - API 文档：http://localhost:8080/swagger/index.html"
        echo ""
    fi
    
    echo "  💾 PostgreSQL:"
    echo "     - 连接：localhost:5432 (数据库：kiqi)"
    echo ""
    
    echo "  🗄️  Redis:"
    echo "     - 连接：localhost:6379"
    echo ""
    
    echo "========================================"
    echo ""
    print_info "提示:"
    echo "  - 查看 API 日志：tail -f /tmp/api.log"
    echo "  - 查看 Worker 日志：tail -f /tmp/worker.log"
    echo "  - 执行测试流程：./scripts/dev/core-flow-test.sh"
    echo "  - 停止所有服务：$0 --clean"
    echo ""
    print_success "本地开发环境准备就绪！"
    echo ""
}

# 执行主函数
main "$@"
