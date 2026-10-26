#!/bin/bash
# BDD 场景验证脚本
# 用法: ./scripts/verify-bdd.sh [模块名]
# 示例: ./scripts/verify-bdd.sh user-auth

set -e

BASE_URL="${SHIYEMO_BASE_URL:-http://localhost:8080}"
MODULE=${1:-"all"}

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

pass() { echo -e "${GREEN}✅ $1${NC}"; }
fail() { echo -e "${RED}❌ $1${NC}"; exit 1; }

# 健康检查
echo "=== 健康检查 ==="
HEALTH=$(curl -s "$BASE_URL/health")
echo "$HEALTH" | grep -q '"status":"ok"' && pass "服务正常" || fail "服务不可用"

# 用户注册场景
verify_user_auth() {
    echo ""
    echo "=== 用户注册/登录场景 ==="
    
    # 生成随机邮箱
    RANDOM_EMAIL="bdd_test_$(date +%s)@example.com"
    
    # 场景 1: 邮箱注册
    echo "场景 1: 邮箱注册"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/email" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$RANDOM_EMAIL\",\"password\":\"a123456\",\"nickname\":\"BDDTest\"}")
    echo "$RESPONSE" | grep -q '"code":0' && pass "邮箱注册成功" || fail "邮箱注册失败: $RESPONSE"
    
    # 场景 2: 重复注册
    echo "场景 2: 重复邮箱注册"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/email" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$RANDOM_EMAIL\",\"password\":\"a123456\",\"nickname\":\"BDDTest2\"}")
    echo "$RESPONSE" | grep -q '"code":40206' && pass "重复注册被拒绝" || fail "重复注册未被拒绝: $RESPONSE"
    
    # 场景 3: 邮箱登录
    echo "场景 3: 邮箱登录"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/email" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$RANDOM_EMAIL\",\"password\":\"a123456\"}")
    echo "$RESPONSE" | grep -q '"code":0' && pass "邮箱登录成功" || fail "邮箱登录失败: $RESPONSE"
    TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    # 场景 4: 密码错误
    echo "场景 4: 密码错误"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/email" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$RANDOM_EMAIL\",\"password\":\"wrong\"}")
    echo "$RESPONSE" | grep -q '"code":40204' && pass "密码错误被拒绝" || fail "密码错误未被拒绝: $RESPONSE"
    
    # 场景 5: 获取用户信息
    echo "场景 5: 获取用户信息"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/profile" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取用户信息成功" || fail "获取用户信息失败: $RESPONSE"
    
    echo "$TOKEN" > /tmp/shiyemo_test_token.txt
}

# 帖子场景
verify_post() {
    echo ""
    echo "=== 帖子场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 场景 1: 创建帖子
    echo "场景 1: 创建帖子"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"type":"emotion","title":"BDD测试帖子","content":"测试内容","tags":["测试"]}')
    echo "$RESPONSE" | grep -q '"code":0' && pass "创建帖子成功" || fail "创建帖子失败: $RESPONSE"
    POST_ID=$(echo "$RESPONSE" | grep -o '"data":{[^}]*' | grep -o '"id":[0-9]*' | cut -d: -f2)
    
    # 场景 2: 获取帖子详情
    echo "场景 2: 获取帖子详情"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/posts/$POST_ID")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取帖子详情成功" || fail "获取帖子详情失败: $RESPONSE"
    
    # 场景 3: 更新帖子
    echo "场景 3: 更新帖子"
    RESPONSE=$(curl -s -X PUT "$BASE_URL/api/v1/posts/$POST_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"title":"更新后的标题","content":"更新后的内容"}')
    echo "$RESPONSE" | grep -q '"code":0' && pass "更新帖子成功" || fail "更新帖子失败: $RESPONSE"
    
    # 场景 4: 删除帖子
    echo "场景 4: 删除帖子"
    RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/posts/$POST_ID" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "删除帖子成功" || fail "删除帖子失败: $RESPONSE"
    
    echo "$POST_ID" > /tmp/shiyemo_test_post_id.txt
}

# 评论场景
verify_comment() {
    echo ""
    echo "=== 评论场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 创建新帖子用于评论测试
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"type":"emotion","title":"评论测试帖子","content":"测试内容"}')
    POST_ID=$(echo "$RESPONSE" | grep -o '"data":{[^}]*' | grep -o '"id":[0-9]*' | cut -d: -f2)
    
    # 场景 1: 发表评论
    echo "场景 1: 发表评论"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts/$POST_ID/comments" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"content":"BDD测试评论"}')
    echo "$RESPONSE" | grep -q '"code":0' && pass "发表评论成功" || fail "发表评论失败: $RESPONSE"
    # 匹配最后一个 id（避免匹配到 author.id）
    COMMENT_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | tail -1 | cut -d: -f2)
    
    # 场景 2: 回复评论
    echo "场景 2: 回复评论"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts/$POST_ID/comments" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"content\":\"回复评论\",\"parent_id\":$COMMENT_ID}")
    echo "$RESPONSE" | grep -q '"code":0' && pass "回复评论成功" || fail "回复评论失败: $RESPONSE"
    
    # 场景 3: 获取评论列表
    echo "场景 3: 获取评论列表"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/posts/$POST_ID/comments")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取评论列表成功" || fail "获取评论列表失败: $RESPONSE"
}

# 点赞/收藏场景
verify_interaction() {
    echo ""
    echo "=== 点赞/收藏场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 创建新帖子
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"type":"emotion","title":"点赞测试帖子","content":"测试内容"}')
    POST_ID=$(echo "$RESPONSE" | grep -o '"data":{[^}]*' | grep -o '"id":[0-9]*' | cut -d: -f2)
    
    # 场景 1: 点赞
    echo "场景 1: 点赞帖子"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts/$POST_ID/like" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "点赞成功" || fail "点赞失败: $RESPONSE"
    
    # 场景 2: 取消点赞
    echo "场景 2: 取消点赞"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts/$POST_ID/like" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "取消点赞成功" || fail "取消点赞失败: $RESPONSE"
    
    # 场景 3: 收藏
    echo "场景 3: 收藏帖子"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/posts/$POST_ID/bookmark" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "收藏成功" || fail "收藏失败: $RESPONSE"
}

# 关注场景
verify_follow() {
    echo ""
    echo "=== 关注场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 创建第二个用户
    RANDOM_EMAIL2="bdd_test2_$(date +%s)@example.com"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/email" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$RANDOM_EMAIL2\",\"password\":\"a123456\",\"nickname\":\"BDDTest2\"}")
    USER2_ID=$(echo "$RESPONSE" | grep -o '"data":{[^}]*' | grep -o '"id":[0-9]*' | cut -d: -f2)
    
    # 场景 1: 关注用户
    echo "场景 1: 关注用户"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/users/$USER2_ID/follow" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "关注成功" || fail "关注失败: $RESPONSE"
    
    # 场景 2: 取消关注
    echo "场景 2: 取消关注"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/users/$USER2_ID/follow" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "取消关注成功" || fail "取消关注失败: $RESPONSE"
}

# 搜索场景
verify_search() {
    echo ""
    echo "=== 搜索场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 先同步索引
    curl -s -X POST "$BASE_URL/api/v1/search/sync" \
      -H "Authorization: Bearer $TOKEN" > /dev/null
    
    # 场景 1: 搜索帖子
    echo "场景 1: 搜索帖子"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/search/posts?q=BDD")
    echo "$RESPONSE" | grep -q '"code":0' && pass "搜索帖子成功" || fail "搜索帖子失败: $RESPONSE"
    
    # 场景 2: 搜索用户
    echo "场景 2: 搜索用户"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/search/users?q=BDD")
    echo "$RESPONSE" | grep -q '"code":0' && pass "搜索用户成功" || fail "搜索用户失败: $RESPONSE"
    
    # 场景 3: 搜索建议
    echo "场景 3: 搜索建议"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/search/suggestions?q=BDD")
    echo "$RESPONSE" | grep -q '"code":0' && pass "搜索建议成功" || fail "搜索建议失败: $RESPONSE"
}

# 用户统计数据场景
verify_user_stats() {
    echo ""
    echo "=== 用户统计数据场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 场景 1: 获取用户统计数据
    echo "场景 1: 获取用户统计数据"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users/1/stats")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取用户统计数据成功" || fail "获取用户统计数据失败: $RESPONSE"
    echo "$RESPONSE" | grep -q '"post_count"' && pass "包含帖子数" || fail "缺少帖子数字段"
    echo "$RESPONSE" | grep -q '"follower_count"' && pass "包含关注数" || fail "缺少关注数字段"
    echo "$RESPONSE" | grep -q '"following_count"' && pass "包含粉丝数" || fail "缺少粉丝数字段"
    echo "$RESPONSE" | grep -q '"like_received"' && pass "包含获赞数" || fail "缺少获赞数字段"
}

# 通知分类场景
verify_notification() {
    echo ""
    echo "=== 通知分类场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 场景 1: 获取全部通知
    echo "场景 1: 获取全部通知"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/notifications" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取全部通知成功" || fail "获取全部通知失败: $RESPONSE"
    
    # 场景 2: 按类型筛选通知
    echo "场景 2: 按类型筛选通知"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/notifications?type=like" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "按类型筛选成功" || fail "按类型筛选失败: $RESPONSE"
    
    # 场景 3: 筛选未读通知
    echo "场景 3: 筛选未读通知"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/notifications?unread=true" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "筛选未读成功" || fail "筛选未读失败: $RESPONSE"
}

# 审核队列场景（需要管理员权限）
verify_admin() {
    echo ""
    echo "=== 审核队列场景 ==="
    
    # 使用已有的管理员用户（ID=1）
    ADMIN_EMAIL="test@example.com"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/email" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"123456\"}")
    ADMIN_TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    # 场景 1: 获取待审核帖子
    echo "场景 1: 获取待审核帖子"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/posts/pending" \
      -H "Authorization: Bearer $ADMIN_TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取待审核帖子成功" || fail "获取待审核帖子失败: $RESPONSE"
}

# 用户主页数据场景
verify_user_profile() {
    echo ""
    echo "=== 用户主页数据场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 场景 1: 获取用户主页数据
    echo "场景 1: 获取用户主页数据"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users/1" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":0' && pass "获取用户主页数据成功" || fail "获取用户主页数据失败: $RESPONSE"
    echo "$RESPONSE" | grep -q '"post_count"' && pass "包含帖子数" || fail "缺少帖子数字段"
    echo "$RESPONSE" | grep -q '"follower_count"' && pass "包含关注数" || fail "缺少关注数字段"
    echo "$RESPONSE" | grep -q '"following_count"' && pass "包含粉丝数" || fail "缺少粉丝数字段"
    
    # 场景 2: 获取不存在用户
    echo "场景 2: 获取不存在用户"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users/999999" \
      -H "Authorization: Bearer $TOKEN")
    echo "$RESPONSE" | grep -q '"code":40203' && pass "用户不存在错误正确" || fail "用户不存在错误不正确: $RESPONSE"
}

# 通知批量操作场景
verify_notification_batch() {
    echo ""
    echo "=== 通知批量操作场景 ==="
    TOKEN=$(cat /tmp/shiyemo_test_token.txt)
    
    # 场景 1: 批量标记已读
    echo "场景 1: 批量标记已读"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/notifications/batch/read" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"ids":[1,2,3]}')
    echo "$RESPONSE" | grep -q '"code":0' && pass "批量标记已读成功" || fail "批量标记已读失败: $RESPONSE"
    
    # 场景 2: 批量删除通知
    echo "场景 2: 批量删除通知"
    RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/notifications/batch/delete" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"ids":[1,2,3]}')
    echo "$RESPONSE" | grep -q '"code":0' && pass "批量删除成功" || fail "批量删除失败: $RESPONSE"
}

# 主执行逻辑
case $MODULE in
    "user-auth")
        verify_user_auth
        ;;
    "post")
        verify_user_auth
        verify_post
        ;;
    "comment")
        verify_user_auth
        verify_post
        verify_comment
        ;;
    "interaction")
        verify_user_auth
        verify_post
        verify_interaction
        ;;
    "follow")
        verify_user_auth
        verify_follow
        ;;
    "search")
        verify_user_auth
        verify_search
        ;;
    "user-stats")
        verify_user_stats
        ;;
    "notification")
        verify_user_auth
        verify_notification
        ;;
    "admin")
        verify_user_auth
        verify_admin
        ;;
    "all")
        verify_user_auth
        verify_post
        verify_comment
        verify_interaction
        verify_follow
        verify_search
        verify_user_stats
        verify_notification
        verify_admin
        verify_user_profile
        verify_notification_batch
        ;;
    *)
        echo "用法: $0 [user-auth|post|comment|interaction|follow|search|user-stats|notification|admin|user-profile|notification-batch|all]"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=== 所有验证通过！ ===${NC}"
