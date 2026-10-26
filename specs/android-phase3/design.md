# design：Android App Phase 3 — 社交功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 搜索防抖 | 500ms delay + cancel | 避免每次输入都请求，减少服务端压力 |
| 搜索实现 | ViewModel 层 Job 管理 | 新输入取消旧 Job，避免结果乱序 |
| 通知已读 | 乐观更新（先改 UI 再请求） | 提升用户体验，失败静默处理 |
| 关注状态 | 本地 State 即时更新 + API 请求 | 避免重新加载整个页面 |
| 用户主页路由 | `user/{userId}` 参数化导航 | 同一 ProfileScreen 支持自己和他人主页 |
| Tab 设计 | 按 isOwnProfile 动态显示 | 自己主页 4 Tab（帖子/收藏/关注/粉丝），他人 3 Tab（无收藏） |

## 模块影响

| 层级 | 文件 | 改动类型 | 说明 |
|------|------|---------|------|
| data/remote | SearchApi.kt | 新增 | 搜索 API 接口 |
| data/remote | NotificationApi.kt | 新增 | 通知 API 接口 |
| data/remote | UserApi.kt | 修改 | 新增 follow/unfollow/following/followers |
| data/model | PostDto.kt | 修改 | 新增 SearchResultDto、NotificationDto |
| data/model | UserDto.kt | 修改 | 新增 isFollowing 字段 |
| data/repository | SearchRepository.kt | 新增 | 搜索数据仓库 |
| data/repository | NotificationRepository.kt | 新增 | 通知数据仓库 |
| data/repository | ProfileRepository.kt | 修改 | 新增关注/粉丝相关方法 |
| di | NetworkModule.kt | 修改 | 注册 SearchApi、NotificationApi |
| ui/search | SearchViewModel.kt | 新增 | 搜索状态管理 + 防抖 |
| ui/search | SearchScreen.kt | 重写 | 完整搜索 UI |
| ui/notification | NotificationViewModel.kt | 新增 | 通知状态管理 |
| ui/notification | NotificationScreen.kt | 重写 | 完整通知 UI |
| ui/profile | ProfileViewModel.kt | 重写 | 支持他人主页 + 关注/粉丝列表 |
| ui/profile | ProfileScreen.kt | 重写 | 关注按钮 + 用户列表 Tab |
| ui/navigation | AppNavigation.kt | 修改 | 新增 UserProfile 路由 |

## 接口设计（复用现有后端 API）

```
# 搜索
GET /api/v1/search?q={query}&type={type}&page={page}&page_size={size}

# 通知
GET /api/v1/notifications?page={page}&page_size={size}
PUT  /api/v1/notifications/{id}/read
PUT  /api/v1/notifications/read-all

# 关注
POST   /api/v1/users/{id}/follow
DELETE /api/v1/users/{id}/follow
GET    /api/v1/users/{id}/following?page={page}&page_size={size}
GET    /api/v1/users/{id}/followers?page={page}&page_size={size}

# 用户详情（已有）
GET /api/v1/users/{id}  → 返回 is_following 字段
```

## 数据变更
- 无数据库变更（全部复用现有后端 API）
- UserDto 新增 `isFollowing` 字段映射（后端已返回）

## 依赖
- 后端搜索模块（Meilisearch）：已就绪
- 后端通知模块：已就绪
- 后端关注模块：已就绪

## 风险点
- 搜索 API 返回格式可能与现有 `SearchResultDto` 不完全匹配，需联调确认
- 通知类型枚举值需与后端一致（like/comment/follow/system）
