#!/bin/bash
# 全量同步帖子和用户到 Meilisearch 索引

MEILI_URL="http://shiyemo-search:7700"
: "${MEILI_KEY:?Set MEILI_KEY before synchronizing the search index}"
DB_CMD="docker exec shiyemo-db psql -U shiyemo -d shiyemo -t -A"

# 导出帖子为 JSON
POSTS=$($DB_CMD -c "SELECT json_agg(json_build_object('id', id, 'title', title, 'content', substring(content from 1 for 500), 'type', type, 'status', status, 'like_count', like_count, 'created_at', extract(epoch from created_at)::int)) FROM posts WHERE deleted_at IS NULL AND status = 1")

# 导出用户为 JSON
USERS=$($DB_CMD -c "SELECT json_agg(json_build_object('id', id, 'nickname', nickname, 'avatar', avatar, 'bio', bio)) FROM users WHERE status = 1")

# 索引帖子
echo "Indexing posts..."
docker exec -i shiyemo-search wget -qO- \
  --post-data="$POSTS" \
  --header="Content-Type: application/json" \
  --header="Authorization: Bearer $MEILI_KEY" \
  "http://127.0.0.1:7700/indexes/posts/documents?primaryKey=id"

echo ""

# 索引用户
echo "Indexing users..."
docker exec -i shiyemo-search wget -qO- \
  --post-data="$USERS" \
  --header="Content-Type: application/json" \
  --header="Authorization: Bearer $MEILI_KEY" \
  "http://127.0.0.1:7700/indexes/users/documents?primaryKey=id"

echo ""
echo "Done."
