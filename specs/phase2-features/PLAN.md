# Phase 2 实施计划

> 目标：为"失业么"增加失业人群专属差异化功能，按 P0 → P1 → P2 分批实施。

---

## 总览

```
Sprint 1 (P0)  ─── 邮箱验证 + 自由天数 + 生存贴士 ─── 预计 4-5 天
Sprint 2 (P1)  ─── 今日打卡+图片 + 匿名模式 ─── 预计 5-6 天
Sprint 3 (P2)  ─── 同温层匹配 + 趋势可视化 ─── 预计 7-10 天
```

---

## Sprint 1：P0 功能

### 1.1 邮箱验证

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | User 模型新增 `EmailVerified bool` 字段 | `auth/model.go` |
| 2 | 创建 `verify` 模块（handler/service/repo/model） | `internal/verify/` |
| 3 | 发送验证邮件 Service（SMTP/第三方邮件） | `common/email/sender.go` |
| 4 | 验证 Token 生成与校验 | `verify/service.go` |
| 5 | 注册成功后触发发送验证邮件 | `auth/service.go` 修改 |
| 6 | 新增路由：`POST /auth/verify/send`、`GET /auth/verify/:token` | `cmd/main.go` |
| 7 | 中间件：未验证邮箱的用户发帖时返回提示 | `common/middleware/` |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 注册成功后跳转"验证邮箱"提示页 | `web/src/app/verify-email/page.tsx` |
| 2 | 邮箱验证落地页（点击邮件链接后） | `web/src/app/verify-email/confirm/page.tsx` |
| 3 | 发帖时若未验证弹出提示 | `web/src/app/posts/create/page.tsx` 修改 |

**配置**

| 步骤 | 内容 |
|------|------|
| 1 | `config.yaml` 新增 SMTP 配置段 |
| 2 | `docker-compose.yml` 新增 SMTP 环境变量 |

---

### 1.2 自由天数计时器

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | User 模型新增 `LastWorkDate *time.Time` 字段 | `auth/model.go` |
| 2 | 用户资料更新接口支持 `last_work_date` 字段 | `user/handler.go` + `user/service.go` |
| 3 | 用户详情接口返回计算后的 `freedom_days` | `user/service.go` |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 个人资料编辑：新增"最后工作日"日期选择器 | `profile/page.tsx` 修改 |
| 2 | 首页/个人中心展示"自由第 X 天"卡片 | `components/FreedomDayCard.tsx` |
| 3 | 用户可选择隐藏/显示 | 设置项 |

---

### 1.3 生存小贴士

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 创建 `tips` 模块 | `internal/tips/` |
| 2 | 贴士表 migration | 数据库迁移（GORM AutoMigrate） |
| 3 | 种子数据：初始 20-30 条贴士 | `tips/seed.go` |
| 4 | API：`GET /tips/random` — 随机返回一条 | `tips/handler.go` |
| 5 | API：`GET /tips/list` — 管理员列表（CRUD） | `tips/handler.go` |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 首页展示"今日贴士"卡片组件 | `components/DailyTip.tsx` |
| 2 | 点击可刷新/收藏 | 集成到首页 |

---

## Sprint 2：P1 功能

### 2.1 今日打卡（心情记录 + 图片）

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 创建 `checkin` 模块 | `internal/checkin/` |
| 2 | Model：`DailyCheckin` 表 | `checkin/model.go` |
| 3 | Service：创建打卡、查询历史、心情统计 | `checkin/service.go` |
| 4 | API 路由 | |
|   | `POST /checkin` — 今日打卡 | |
|   | `GET /checkin/today` — 查询今天是否已打卡 | |
|   | `GET /checkin/history` — 我的打卡历史 | |
|   | `GET /checkin/mood-curve` — 近 30 天心情数据 | |
|   | `GET /checkin/timeline` — 社区公开打卡时间线 | |
| 5 | 图片上传复用现有 `upload` 模块 | |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 打卡页面：心情选择 + 文字 + 图片 | `web/src/app/checkin/page.tsx` |
| 2 | 个人中心：心情曲线图（近 30 天） | `components/MoodCurve.tsx` |
| 3 | 社区时间线中穿插打卡卡片 | Feed 组件修改 |
| 4 | 首页入口："记录今天" 按钮 | 首页组件修改 |

---

### 2.2 匿名模式

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | Post 模型新增 `IsAnonymous bool`、`AnonymousName string` | `post/model.go` 修改 |
| 2 | Comment 模型新增同样字段 | `comment/model.go` 修改 |
| 3 | 发帖/评论接口支持 `is_anonymous` 参数 | handler 修改 |
| 4 | 列表/详情接口：匿名帖隐藏真实作者信息 | service 修改 |
| 5 | 匿名名称生成器（动物名随机组合） | `common/anonymous/generator.go` |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 发帖/评论表单增加"匿名发布"开关 | 表单组件修改 |
| 2 | 帖子/评论卡片：匿名显示特殊样式 | 卡片组件修改 |
| 3 | 作者本人可看到自己的匿名帖标记 | 列表逻辑修改 |

---

## Sprint 3：P2 功能

### 3.1 行业/城市同温层

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | User 模型新增 `City`、`Industry`、`LeftJobAt` | `auth/model.go` |
| 2 | 用户资料接口支持新字段 | `user/` 修改 |
| 3 | 同温层匹配 Service：按 city+industry+时间段聚合 | `internal/cohort/service.go` |
| 4 | Feed 新增"同温层"过滤维度 | `post/` 修改 |
| 5 | API：`GET /cohort/stats` — 我的同温层统计 | `cohort/handler.go` |
| 6 | API：`GET /cohort/peers` — 同温层用户列表 | `cohort/handler.go` |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 注册/资料页新增城市、行业、离职时间选择 | 表单修改 |
| 2 | Feed 页新增"同温层"Tab | `feed/page.tsx` 修改 |
| 3 | 同温层主页：展示匹配的人和帖子 | `web/src/app/cohort/page.tsx` |

---

### 3.2 趋势可视化

**后端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 定时任务：每日统计各城市/行业用户数 | `internal/stats/` |
| 2 | API：`GET /stats/trends` — 趋势数据 | `stats/handler.go` |
| 3 | 缓存：Redis 缓存统计结果 | |

**前端**

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 数据可视化组件（简单柱/线图） | `components/TrendChart.tsx` |
| 2 | 集成到首页或同温层页面 | |
| 3 | 文案：如"你所在的互联网行业，本月有 X 人加入" | |

---

## 执行顺序与依赖关系

```
邮箱验证 ──┐
自由天数 ──┼── Sprint 1（并行，无依赖）
生存贴士 ──┘
            │
            ▼
今日打卡 ──┐
匿名模式 ──┘── Sprint 2（并行，无依赖）
            │
            ▼
同温层匹配 ──┐
趋势可视化 ──┘── Sprint 3（趋势依赖同温层的字段）
```

---

## 数据库变更汇总

```sql
-- Sprint 1
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN last_work_date DATE;

-- Sprint 2
ALTER TABLE posts ADD COLUMN is_anonymous BOOLEAN DEFAULT FALSE;
ALTER TABLE posts ADD COLUMN anonymous_name VARCHAR(50) DEFAULT '';
ALTER TABLE comments ADD COLUMN is_anonymous BOOLEAN DEFAULT FALSE;
ALTER TABLE comments ADD COLUMN anonymous_name VARCHAR(50) DEFAULT '';

CREATE TABLE daily_checkins (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    mood VARCHAR(20) NOT NULL,
    note TEXT DEFAULT '',
    image_url VARCHAR(500) DEFAULT '',
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_daily_checkin UNIQUE (user_id, (created_at::date))
);

CREATE TABLE survival_tips (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    category VARCHAR(50) DEFAULT 'general',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE email_verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token VARCHAR(100) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Sprint 3
ALTER TABLE users ADD COLUMN city VARCHAR(50) DEFAULT '';
ALTER TABLE users ADD COLUMN industry VARCHAR(50) DEFAULT '';
ALTER TABLE users ADD COLUMN left_job_at DATE;
```

---

## 新增后端模块清单

| 模块 | 路径 | Sprint |
|------|------|--------|
| verify | `server/internal/verify/` | 1 |
| tips | `server/internal/tips/` | 1 |
| checkin | `server/internal/checkin/` | 2 |
| anonymous | `server/common/anonymous/` | 2 |
| cohort | `server/internal/cohort/` | 3 |
| stats | `server/internal/stats/` | 3 |
| email | `server/common/email/` | 1 |

---

## 新增前端页面清单

| 页面 | 路径 | Sprint |
|------|------|--------|
| 验证邮箱提示 | `web/src/app/verify-email/page.tsx` | 1 |
| 验证确认 | `web/src/app/verify-email/confirm/page.tsx` | 1 |
| 打卡 | `web/src/app/checkin/page.tsx` | 2 |
| 同温层 | `web/src/app/cohort/page.tsx` | 3 |

---

## 新增前端组件清单

| 组件 | 路径 | Sprint |
|------|------|--------|
| 自由天数卡片 | `components/FreedomDayCard.tsx` | 1 |
| 每日贴士 | `components/DailyTip.tsx` | 1 |
| 心情曲线 | `components/MoodCurve.tsx` | 2 |
| 打卡卡片 | `components/CheckinCard.tsx` | 2 |
| 趋势图 | `components/TrendChart.tsx` | 3 |

---

## 部署注意事项

- SMTP 配置：需要在 `docker-compose.yml` 新增邮件服务环境变量
- 数据库迁移：GORM AutoMigrate 自动处理字段新增，新表需注册
- 前端构建：新页面/组件无特殊配置
- Meilisearch：用户索引需新增 `city`、`industry` 为可过滤字段（Sprint 3）

---

## 验收标准

### Sprint 1 完成标志
- [ ] 注册后收到验证邮件，点击链接完成验证
- [ ] 未验证用户发帖时弹出提示
- [ ] 个人资料可设置"最后工作日"
- [ ] 首页/个人中心显示"自由第 X 天"
- [ ] 首页展示随机生存贴士
- [ ] 全部功能部署到服务器可用

### Sprint 2 完成标志
- [ ] 可选择心情 + 文字 + 图片完成打卡
- [ ] 个人中心展示心情曲线
- [ ] 发帖/评论可选匿名
- [ ] 匿名帖在社区正确显示
- [ ] 全部功能部署到服务器可用

### Sprint 3 完成标志
- [ ] 用户可填写城市/行业/离职时间
- [ ] Feed 有"同温层"Tab 并正确过滤
- [ ] 首页展示同温层统计数据
- [ ] 全部功能部署到服务器可用

---

*制定日期：2026-05-26*
*状态：待执行*
