# 管理端增强计划

> 创建时间: 2026-05-27
> 状态: 进行中
> 负责: shiyemo-frontend（前端开发）

## 背景

管理端已有基础功能（统计概览、帖子列表、用户列表、敏感词管理），但缺少多个关键管理功能。后端 API 已全部实现并注册路由，前端页面缺失。

## 目标

补齐管理端缺失功能，提升管理员操作效率。

## 功能清单

### Phase 1: 核心管理功能（P0）

| # | 功能 | 页面路径 | 后端 API | 状态 |
|---|------|----------|----------|------|
| 1 | 举报管理页面 | `/admin/reports` | `GET /admin/reports` + `PUT /admin/reports/:id` | ⏳ |
| 2 | 帖子审核页面 | `/admin/review` | `GET /admin/posts/pending` + `POST /admin/posts/:id/approve` + `POST /admin/posts/:id/reject` | ⏳ |
| 3 | 帖子置顶/加精按钮 | `/admin`（增强） | `PUT /admin/posts/:id/top` + `PUT /admin/posts/:id/featured` | ⏳ |

### Phase 2: 管理体验优化（P1）

| # | 功能 | 页面路径 | 后端 API | 状态 |
|---|------|----------|----------|------|
| 4 | 评论管理页面 | `/admin/comments` | `GET /admin/comments` + `PUT /admin/comments/:id/hide` + `PUT /admin/comments/:id/restore` | ⏳ |
| 5 | 用户搜索/筛选 | `/admin`（增强） | `GET /admin/users?search=xxx&status=xxx` | ⏳ |
| 6 | 帖子内容预览 | `/admin`（增强） | — | ⏳ |

### Phase 3: 可选功能（P2）

| # | 功能 | 说明 | 状态 |
|---|------|------|------|
| 7 | 反馈管理 | 需先设计后端 API | ❌ 待讨论 |

## 技术方案

### 1. 举报管理页面 `/admin/reports`

**页面结构：**
- 举报列表（表格）
  - 举报人
  - 被举报内容（帖子/评论/用户）
  - 举报原因
  - 状态（待处理/已处理/已忽略）
  - 举报时间
  - 操作按钮（标记已处理、忽略）
- 筛选：按状态筛选
- 分页

**API 对接：**
- `GET /api/v1/admin/reports?page=1&size=20` - 获取举报列表
- `PUT /api/v1/admin/reports/:id` - 处理举报（status: 2=已处理, 3=已忽略）

### 2. 帖子审核页面 `/admin/review`

**页面结构：**
- 待审核帖子列表
  - 标题、作者、类型、发布时间
  - 内容预览（点击展开）
  - 操作按钮（通过、拒绝需填理由）
- 审核通过/拒绝弹窗

**API 对接：**
- `GET /api/v1/admin/posts/pending?page=1&page_size=20` - 待审核列表
- `POST /api/v1/admin/posts/:id/approve` - 审核通过
- `POST /api/v1/admin/posts/:id/reject` - 审核拒绝（需 reason）

### 3. 帖子置顶/加精按钮

**增强位置：** 现有 `/admin` 页面的帖子列表操作列

**交互：**
- 添加「置顶」按钮（切换状态）
- 添加「加精」按钮（切换状态）
- 标题前显示图标标识（📌 置顶、⭐ 加精）

**API 对接：**
- `PUT /api/v1/admin/posts/:id/top` `{ "is_top": true/false }`
- `PUT /api/v1/admin/posts/:id/featured` `{ "is_featured": true/false }`

### 4. 评论管理页面 `/admin/comments`

**页面结构：**
- 评论列表（表格）
  - 评论内容
  - 评论人
  - 所属帖子
  - 状态（正常/已屏蔽）
  - 评论时间
  - 操作按钮（屏蔽/恢复）
- 按帖子筛选
- 分页

**API 对接：**
- `GET /api/v1/admin/comments?page=1&page_size=20&post_id=xxx` - 评论列表
- `PUT /api/v1/admin/comments/:id/hide` - 屏蔽评论
- `PUT /api/v1/admin/comments/:id/restore` - 恢复评论

### 5. 用户搜索/筛选

**增强位置：** 现有 `/admin` 页面的用户列表区域

**交互：**
- 搜索框（按昵称/邮箱/手机号搜索）
- 筛选下拉（按状态、角色）
- URL 参数同步

**API 对接：**
- `GET /api/v1/admin/users?page=1&page_size=10&search=xxx&status=xxx&role=xxx`

### 6. 帖子内容预览

**增强位置：** 现有 `/admin` 页面的帖子列表

**交互：**
- 鼠标悬停标题显示内容摘要（Tooltip）
- 或点击标题弹窗显示完整内容

## 文件变更计划

### 新增文件
- `web/src/app/admin/reports/page.tsx` - 举报管理页面
- `web/src/app/admin/review/page.tsx` - 帖子审核页面
- `web/src/app/admin/comments/page.tsx` - 评论管理页面

### 修改文件
- `web/src/app/admin/page.tsx` - 增强帖子列表（置顶/加精按钮、内容预览）和用户列表（搜索/筛选）
- `web/src/app/admin/layout.tsx` - 添加导航菜单项（如需要）

## 验证标准

| 检查 | 命令 | 通过条件 |
|------|------|---------|
| 构建检查 | `npm run build` | 0 error |
| 类型检查 | `npx tsc --noEmit` | 0 error |
| 功能验证 | 手动测试 | 所有操作正常 |

## 开发顺序

1. 帖子置顶/加精按钮（增强现有页面）
2. 举报管理页面（新建）
3. 帖子审核页面（新建）
4. 评论管理页面（新建）
5. 用户搜索/筛选（增强现有页面）
6. 帖子内容预览（增强现有页面）

## 里程碑

- [ ] Phase 1 完成：核心管理功能可用
- [ ] Phase 2 完成：管理体验优化
- [ ] 部署验证：生产环境测试通过
