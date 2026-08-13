#!/bin/bash

# 运行所有测试脚本
# 用法: ./scripts/testing/run-all-tests.sh

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_header() {
    echo ""
    echo "======================================"
    echo -e "${GREEN}$1${NC}"
    echo "======================================"
    echo ""
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

TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

run_test() {
    local test_name=$1
    local test_command=$2
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo ""
    print_header "Running: $test_name"
    
    if eval "$test_command"; then
        print_success "$test_name PASSED"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        print_error "$test_name FAILED"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

echo "🧪 Running all tests..."

# 1. 单元测试
run_test "Unit Tests" "cd server && go test -v ./internal/... ./pkg/... -count=1"

# 2. 集成测试
run_test "Integration Tests" "cd server && go test -v ./test/integration/... -count=1"

# 3. 核心功能测试（需要服务运行）
echo ""
print_warning "Checking if services are running for E2E tests..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    run_test "E2E Core Flow Tests" "./scripts/testing/core-flow-test.sh"
else
    print_warning "Skipping E2E tests (services not running)"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# 总结
echo ""
echo "======================================"
echo "📊 Test Summary"
echo "======================================"
echo "Total:  $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    print_success "All tests passed!"
    exit 0
else
    print_error "Some tests failed"
    exit 1
fi
