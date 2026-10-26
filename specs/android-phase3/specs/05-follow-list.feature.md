# spec：App 关注 — 关注/粉丝列表

## 场景描述
用户查看关注列表和粉丝列表

## 规格

### Given 用户主页已加载
- 用户有 3 个关注、5 个粉丝

### When 用户切换到"关注" Tab
- 系统发送 GET /api/v1/users/{id}/following?page=1&page_size=20
- 显示关注用户列表（头像、昵称、简介）
- 列表支持无限滚动

### When 用户切换到"粉丝" Tab
- 系统发送 GET /api/v1/users/{id}/followers?page=1&page_size=20
- 显示粉丝用户列表
- 列表支持无限滚动

### When 用户点击列表中的某个用户
- 跳转到该用户的主页（user/{userId}）
