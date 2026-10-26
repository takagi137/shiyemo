#!/bin/bash

# 用户认证功能测试脚本

BASE_URL="http://localhost:8080"

echo "🧪 测试用户认证功能..."
echo ""

# 测试 1: 邮箱注册
echo "📝 测试 1: 邮箱注册"
REGISTER_EMAIL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "a123456",
    "nickname": "测试用户"
  }')
echo "响应: $REGISTER_EMAIL_RESPONSE"
echo ""

# 提取 token
EMAIL_TOKEN=$(echo $REGISTER_EMAIL_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $EMAIL_TOKEN"
echo ""

# 测试 2: 邮箱登录
echo "📝 测试 2: 邮箱登录"
LOGIN_EMAIL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "a123456"
  }')
echo "响应: $LOGIN_EMAIL_RESPONSE"
echo ""

# 测试 3: 手机号注册（使用 mock 验证码 123456）
echo "📝 测试 3: 手机号注册"
REGISTER_PHONE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/phone" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "code": "123456",
    "nickname": "手机用户"
  }')
echo "响应: $REGISTER_PHONE_RESPONSE"
echo ""

# 测试 4: 手机号登录
echo "📝 测试 4: 手机号登录"
LOGIN_PHONE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/phone" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "code": "123456"
  }')
echo "响应: $LOGIN_PHONE_RESPONSE"
echo ""

# 测试 5: 访问受保护接口（获取用户信息）
echo "📝 测试 5: 访问受保护接口"
if [ -n "$EMAIL_TOKEN" ]; then
  PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/profile" \
    -H "Authorization: Bearer $EMAIL_TOKEN")
  echo "响应: $PROFILE_RESPONSE"
else
  echo "⚠️ 无法获取 token，跳过此测试"
fi
echo ""

# 测试 6: 重复注册（应该返回错误）
echo "📝 测试 6: 重复邮箱注册"
DUPLICATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "a123456",
    "nickname": "测试用户2"
  }')
echo "响应: $DUPLICATE_RESPONSE"
echo ""

# 测试 7: 密码错误登录
echo "📝 测试 7: 密码错误登录"
WRONG_PASSWORD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }')
echo "响应: $WRONG_PASSWORD_RESPONSE"
echo ""

# 测试 8: 无 token 访问受保护接口
echo "📝 测试 8: 无 token 访问受保护接口"
NO_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/profile")
echo "响应: $NO_TOKEN_RESPONSE"
echo ""

echo "✅ 测试完成！"
