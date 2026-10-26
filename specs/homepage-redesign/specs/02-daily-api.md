# spec：每日API - 获取今日内容

## 场景描述
前端请求今日的风景图和温暖文案

## 规格

### Given 系统已有今日内容
- Redis 中缓存了今日的风景图 URL
- tips 表中有多条文案

### When 前端请求 GET /api/v1/daily
- 请求无需认证

### Then 返回今日内容
- 响应状态码 200
- 返回 image.url（有效的 Unsplash 图片 URL）
- 返回 image.description（图片描述）
- 返回 quote.text（温暖文案）
- 响应时间 < 200ms（Redis 缓存命中）

### When Redis 缓存过期
- 系统从 Unsplash API 获取新图片
- 系统从 tips 表随机选取一条文案
- 写入 Redis 缓存（TTL = 24小时）
- 返回新内容
