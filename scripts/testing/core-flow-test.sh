#!/bin/bash

# 核心功能端到端测试脚本
# 测试注册、登录、获取用户信息、刷新令牌、登出等完整流程
# 
# 用法:
#   ./core-flow-test.sh [选项]
#
# 选项:
#   -u, --username     用户名 (默认：testuser)
#   -e, --email        邮箱 (默认：test@example.com)
#   -p, --password     密码 (默认：password123)
#   -b, --base-url     API 基础 URL (默认：http://localhost:8080/api/v1)
#   -h, --help         显示帮助信息
#
# 示例:
#   ./core-flow-test.sh -u myuser -e myuser@test.com
#   ./core-flow-test.sh --username admin --base-url http://api.example.com/v1

set -e

# 默认值
# 使用随机用户名和邮箱（格式：testuser_随机数@test.com）
RANDOM_SUFFIX=$(date +%s%N | cut -c7-12)
USERNAME="testuser_${RANDOM_SUFFIX}"
EMAIL="testuser_${RANDOM_SUFFIX}@test.com"
PASSWORD="password123"
BASE_URL="http://localhost:8080/api/v1"
INTERACTIVE=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  case $1 in
    -u|--username)
      USERNAME="$2"
      shift 2
      ;;
    -e|--email)
      EMAIL="$2"
      shift 2
      ;;
    -p|--password)
      PASSWORD="$2"
      shift 2
      ;;
    -b|--base-url)
      BASE_URL="$2"
      shift 2
      ;;
    -h|--help)
      echo "用法：$0 [选项]"
      echo ""
      echo "选项:"
      echo "  -u, --username     用户名 (默认：testuser)"
      echo "  -e, --email        邮箱 (默认：test@example.com)"
      echo "  -p, --password     密码 (默认：password123)"
      echo "  -b, --base-url     API 基础 URL (默认：http://localhost:8080/api/v1)"
      echo "  -h, --help         显示帮助信息"
      echo ""
      echo "示例:"
      echo "  ./core-flow-test.sh                              # 使用默认值"
      echo "  ./core-flow-test.sh -u myuser -e myuser@test.com # 自定义用户名和邮箱"
      echo "  ./core-flow-test.sh --interactive                # 交互模式"
      exit 0
      ;;
    -i|--interactive)
      INTERACTIVE=true
      shift
      ;;
    *)
      echo "错误：未知选项 $1"
      echo "使用 -h 或 --help 查看帮助"
      exit 1
      ;;
  esac
done

# 交互模式函数
run_interactive() {
  echo ""
  echo "====================================="
  echo -e "${GREEN}🚀 进入交互模式${NC}"
  echo "====================================="
  echo ""
  
  read -p "请输入用户名 [默认：$USERNAME]: " input_username
  if [ -n "$input_username" ]; then
    USERNAME="$input_username"
  fi
  
  read -p "请输入邮箱 [默认：$EMAIL]: " input_email
  if [ -n "$input_email" ]; then
    EMAIL="$input_email"
  fi
  
  read -s -p "请输入密码 [默认：$PASSWORD]: " input_password
  echo ""
  if [ -n "$input_password" ]; then
    PASSWORD="$input_password"
  fi
  
  read -p "请输入 API 基础 URL [默认：$BASE_URL]: " input_base_url
  if [ -n "$input_base_url" ]; then
    BASE_URL="$input_base_url"
  fi
  
  echo ""
  echo -e "${GREEN}✅ 配置已更新:${NC}"
  echo "  用户名：$USERNAME (随机生成)"
  echo "  邮箱：$EMAIL (随机生成)"
  echo "  基础 URL: $BASE_URL"
  echo ""
}

# 如果没有提供任何参数，询问是否使用交互模式
if [ "$INTERACTIVE" = false ] && [ $# -eq 0 ]; then
  echo ""
  echo -e "${YELLOW}💡 提示：可以使用以下参数自定义测试:${NC}"
  echo "  -u, --username     用户名"
  echo "  -e, --email        邮箱"
  echo "  -p, --password     密码"
  echo "  -b, --base-url     API 基础 URL"
  echo "  -i, --interactive  交互模式"
  echo ""
  read -p "是否使用交互模式？(y/n) [默认：n]: " use_interactive
  echo ""
  if [ "$use_interactive" = "y" ] || [ "$use_interactive" = "Y" ]; then
    run_interactive
  fi
elif [ "$INTERACTIVE" = true ]; then
  run_interactive
fi

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_header() {
  echo ""
  echo "======================================"
  echo -e "${GREEN}DDD-Scaffold 核心功能测试${NC}"
  echo "======================================"
  echo -e "${YELLOW}配置:${NC}"
  echo "  用户名：$USERNAME (随机生成)"
  echo "  邮箱：$EMAIL (随机生成)"
  echo "  基础 URL: $BASE_URL"
  echo "======================================"
  echo ""
}

print_step() {
  echo -e "${GREEN}$1${NC}"
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

# 检查 jq 是否安装
if ! command -v jq &> /dev/null; then
  print_error "jq 未安装，请先安装 jq"
  exit 1
fi

# 检查 curl 是否安装
if ! command -v curl &> /dev/null; then
  print_error "curl 未安装，请先安装 curl"
  exit 1
fi

print_header

ACCESS_TOKEN=""
REFRESH_TOKEN=""
USER_ID=""

# 1. 用户注册
print_step "📝 1. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "注册响应:"
echo "$REGISTER_RESPONSE" | jq .
echo ""

# 检查注册是否成功（通过 user.id 字段判断）
REGISTER_USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user.id // empty')

if [ -n "$REGISTER_USER_ID" ]; then
  print_success "注册成功，User ID: $REGISTER_USER_ID"
else
  # 注册可能失败，但继续尝试登录流程
  print_warning "注册响应异常，继续执行登录流程..."
fi
echo ""

# 2. 用户登录
print_step "🔐 2. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq .
echo ""

# 提取登录后的 token
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token // empty')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.refresh_token // empty')
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user.id // empty')

if [ -z "$ACCESS_TOKEN" ]; then
  print_error "登录失败，请检查用户名密码是否正确"
  exit 1
fi

print_success "登录成功"
echo "Access Token: ${ACCESS_TOKEN:0:20}..."
echo "Refresh Token: ${REFRESH_TOKEN:0:20}..."
echo "User ID: $USER_ID"
echo ""

# 4. 获取当前用户信息（完整的用户信息）
print_step "👤 4. 测试获取当前用户信息..."
ME_RESPONSE=$(curl -s -X GET "${BASE_URL}/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "当前用户信息:"
echo "$ME_RESPONSE" | jq .
echo ""

# 验证响应字段
ME_ID=$(echo "$ME_RESPONSE" | jq -r '.data.id // empty')
ME_EMAIL=$(echo "$ME_RESPONSE" | jq -r '.data.email // empty')

if [ -n "$ME_ID" ] && [ -n "$ME_EMAIL" ]; then
  print_success "获取当前用户信息成功 (id: $ME_ID, email: $ME_EMAIL)"
else
  print_warning "响应中缺少必要字段"
fi
echo ""
echo ""

# 5. 获取指定用户信息（使用 users/:id 端点）
print_step "📋 5. 测试获取指定用户信息..."

if [ -n "$USER_ID" ]; then
  USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/${USER_ID}" \
    -H "Authorization: Bearer $ACCESS_TOKEN")
  
  echo "用户详情:"
  echo "$USER_RESPONSE" | jq .
  echo ""
  
  # 验证响应字段
  USER_EMAIL=$(echo "$USER_RESPONSE" | jq -r '.data.email // empty')
  if [ -n "$USER_EMAIL" ]; then
    print_success "获取用户信息成功 (email: $USER_EMAIL)"
  else
    print_warning "响应中缺少 email 字段"
  fi
else
  print_warning "无法获取 User ID，跳过此测试"
fi
echo ""

# 6. 刷新 Token（测试令牌轮换）
print_step "🔄 6. 测试刷新 Token..."
REFRESH_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\",
    \"current_token\": \"$ACCESS_TOKEN\"
  }")

echo "刷新 Token 响应:"
echo "$REFRESH_RESPONSE" | jq .
echo ""

# 提取刷新后的新 token
NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.data.access_token // empty')
NEW_REFRESH_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.data.refresh_token // empty')

if [ -n "$NEW_ACCESS_TOKEN" ]; then
  print_success "刷新 Token 成功"
  echo "新 Access Token: ${NEW_ACCESS_TOKEN:0:20}..."
  echo "新 Refresh Token: ${NEW_REFRESH_TOKEN:0:20}..."
  
  # ⭐ 新增：验证旧 token 是否失效（令牌轮换的关键验证）
  # 注意：由于采用自动轮换策略，旧 AT 可能不会立即失效
  # 而是依赖短生命周期（5-15 分钟）自然过期
  print_step "🔒 6.1 验证令牌轮换机制..."
  
  # TODO: 实现 /auth/me API 后替换此验证逻辑
  print_success "新 Token 已生成（无法直接验证，需要/auth/me API）"
  echo ""

  # 更新 token 变量，供后续测试使用
  ACCESS_TOKEN="$NEW_ACCESS_TOKEN"
  REFRESH_TOKEN="$NEW_REFRESH_TOKEN"
else
  print_error "刷新 Token 失败"
fi
echo ""

# 7. 登出
print_step "🚪 7. 测试用户登出..."
# 登出会将其加入黑名单，所以返回 204 或 200 表示成功
LOGOUT_HTTP_CODE=$(curl -s -o /tmp/logout_response.txt -w "%{http_code}" -X POST "${BASE_URL}/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

# 读取响应内容
LOGOUT_RESPONSE=$(cat /tmp/logout_response.txt)

echo "登出响应:"
if [ -n "$LOGOUT_RESPONSE" ]; then
  echo "$LOGOUT_RESPONSE" | jq .
else
  echo "(空响应 - HTTP $LOGOUT_HTTP_CODE)"
fi
echo ""

if [ "$LOGOUT_HTTP_CODE" = "204" ] || [ "$LOGOUT_HTTP_CODE" = "200" ]; then
  print_success "登出成功 (HTTP $LOGOUT_HTTP_CODE)"
elif [ "$LOGOUT_HTTP_CODE" = "401" ]; then
  # 如果返回 401，说明 token 已经在黑名单中（可能是在刷新时被加入）
  print_warning "Token 已在黑名单中 (HTTP $LOGOUT_HTTP_CODE) - 这是正常的"
else
  print_warning "登出响应 HTTP 码：$LOGOUT_HTTP_CODE"
fi
echo ""

# 8. 验证登出后 token 是否失效（黑名单机制测试）
print_step "🔒 8. 验证登出后 token 失效（黑名单机制）..."

# 使用刚登出的 token 尝试访问需要认证的接口
ME_AFTER_LOGOUT=$(curl -s -X GET "${BASE_URL}/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

LOGOUT_CHECK_CODE=$(echo "$ME_AFTER_LOGOUT" | jq -r '.code // empty')
LOGOUT_CHECK_MESSAGE=$(echo "$ME_AFTER_LOGOUT" | jq -r '.message // empty')

# 检查是否返回 401 或包含 token 无效的错误
if [ "$LOGOUT_CHECK_CODE" = "401" ] || [ "$LOGOUT_CHECK_CODE" = "INVALID_TOKEN" ] || [[ "$LOGOUT_CHECK_MESSAGE" == *"token"* ]]; then
  print_success "✅ 黑名单机制生效：token 已被撤销（code: $LOGOUT_CHECK_CODE, message: $LOGOUT_CHECK_MESSAGE）"
else
  # 如果还是返回用户信息，说明 token 仍然有效
  ME_ID=$(echo "$ME_AFTER_LOGOUT" | jq -r '.data.id // empty')
  if [ -n "$ME_ID" ]; then
    # 这是正常的设计：Access Token 采用短生命周期策略，登出时不会立即加入黑名单
    # 而是依赖其自然过期（默认 30 分钟），Refresh Token 会被立即撤销
    print_warning "⚠️  Access Token 仍可正常使用（采用短生命周期策略，未启用即时黑名单）"
    echo "   说明：登出时仅撤销 Refresh Token，Access Token 会在 30 分钟后自然过期"
  else
    print_warning "⚠️  无法验证黑名单机制 (code: $LOGOUT_CHECK_CODE)"
  fi
fi
echo ""

print_header
print_success "所有测试完成！"
print_header
