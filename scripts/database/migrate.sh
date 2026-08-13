#!/bin/bash

# 数据库迁移脚本
# 用法: ./scripts/database/migrate.sh [up|down|status|force]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查环境变量
if [ -z "$DATABASE_URL" ]; then
    # 尝试从 .env 文件加载
    if [ -f "server/configs/.env" ]; then
        export $(grep -v '^#' server/configs/.env | xargs)
    fi
fi

if [ -z "$DATABASE_URL" ]; then
    print_error "DATABASE_URL not set"
    echo "Please set DATABASE_URL or create server/configs/.env"
    exit 1
fi

ACTION=${1:-help}

case $ACTION in
    up)
        print_info "🚀 Running migrations up..."
        cd server
        if command -v migrate &> /dev/null; then
            migrate -path migrations -database "$DATABASE_URL" up
        else
            print_warning "migrate CLI not found, using Makefile..."
            make migrate-up
        fi
        print_info "✅ Migrations completed"
        ;;
    
    down)
        print_warning "⬇️  Rolling back last migration..."
        cd server
        if command -v migrate &> /dev/null; then
            migrate -path migrations -database "$DATABASE_URL" down 1
        else
            make migrate-down
        fi
        print_info "✅ Rollback completed"
        ;;
    
    status)
        print_info "📊 Migration status..."
        cd server
        if command -v migrate &> /dev/null; then
            migrate -path migrations -database "$DATABASE_URL" version
        else
            make migrate-status
        fi
        ;;
    
    force)
        VERSION=${2:-}
        if [ -z "$VERSION" ]; then
            print_error "Version required for force"
            echo "Usage: $0 force <version>"
            exit 1
        fi
        print_warning "⚠️  Forcing migration to version $VERSION..."
        cd server
        if command -v migrate &> /dev/null; then
            migrate -path migrations -database "$DATABASE_URL" force "$VERSION"
        else
            print_error "migrate CLI required for force operation"
            exit 1
        fi
        print_info "✅ Force completed"
        ;;
    
    help|*)
        echo "Usage: $0 {up|down|status|force}"
        echo ""
        echo "Commands:"
        echo "  up      - Run all pending migrations"
        echo "  down    - Rollback last migration"
        echo "  status  - Show current migration version"
        echo "  force   - Force migration to specific version"
        echo ""
        echo "Examples:"
        echo "  $0 up              # Run all migrations"
        echo "  $0 down            # Rollback last migration"
        echo "  $0 status          # Check migration status"
        echo "  $0 force 3         # Force to version 3"
        exit 0
        ;;
esac
