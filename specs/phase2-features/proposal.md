# Phase 2 — 失业者专属功能规划

> 核心原则：不贩卖焦虑，用温和的方式传递"你不是一个人" + "你其实没那么糟"

## 背景

Phase 1 的社区功能（帖子、评论、关注、搜索）和普通社区无异，缺少为失业人群量身定制的功能。
用户反馈："除了口号，看不到和其他社区的区别。"

## 设计哲学

| 原则 | 不做 | 要做 |
|------|------|------|
| 不制造焦虑 | ~~存款倒计时~~ ~~你还能撑X天~~ | 正向的生存贴士、"你比想象中安全" |
| 强调陪伴感 | ~~XX人已重建成功~~ | "你所在行业有 X 人跟你一样" |
| 工具化留存 | ~~每日签到~~ | 心情打卡、求职看板（后续） |
| 低门槛表达 | ~~必须实名~~ | 匿名模式、图片打卡 |

---

## 功能列表

### 1. 个人状态仪表盘

#### 1.1 自由天数计时器

**定位**：不叫"失业天数"，叫"自由第 X 天"——强调个人状态而非失业身份。

- 注册时选填"最后工作日"
- 首页/个人中心显示"自由第 X 天"
- 可选择隐藏
- 文案调性："自由第 47 天。无论如何，这段时间属于你。"

#### 1.2 生存小贴士（替代生存倒计时）

**定位**：不告诉你"还能撑多久"，而是告诉你"其实没那么糟"。

每日随机展示一条正向贴士：
- "见过加班猝死的，没见过不上班饿死的。"
- "一个人生活自己做饭，每月 800 就能吃得很好。"
- "假如你有 10 万存款，每月花 3000，你能休息将近 3 年。"
- "社保断缴可以自己续，医保只要连续缴满 X 个月就行。"
- "失业金可以领 X 个月，去社保局申请就行。"
- "今天不想投简历？那就不投。明天再说。"

实现方式：
- 后端维护贴士库（可后台配置）
- 每日/每次打开随机展示一条
- 用户可以收藏、分享

#### 1.3 今日打卡（心情记录 + 图片）

**定位**：不是"打卡学习"，而是"记录今天你的状态"。

- 每天可以发一条打卡（非强制）
- 选择心情：😴 摆烂 / 😐 平淡 / 🙂 还行 / 💪 积极 / 🎉 突破
- 可附带一张图片（今天的天空、做的饭、散步的路）
- 可选择是否公开（默认公开到社区时间线）
- 个人中心展示心情曲线（近 30 天）

---

### 2. 身份与隐私

#### 2.1 匿名模式

- 发帖/评论时可选择"匿名发布"
- 匿名帖显示为"匿名用户"+ 随机头像
- 同一帖子内同一匿名用户显示相同的"匿名 ID"（如"匿名海豚"）
- 作者自己可以看到自己发的匿名帖
- 管理员可查看（用于审核/举报）

#### 2.2 邮箱验证

- 注册时发送验证邮件
- 未验证邮箱限制发帖（只能浏览和评论）
- 验证后解锁全部功能
- 支持重新发送验证邮件

---

### 3. 同温层社交

#### 3.1 行业/城市自动匹配

**定位**：不是"找工作群"，而是"找跟你一样的人"。

用户注册/编辑资料时填写：
- 城市（省/市）
- 行业（互联网/金融/教育/制造/...）
- 离职时间段（可选）

系统自动匹配同温层：
- "深圳 · 互联网 · 近 3 个月" → 自动聚合这些用户
- 社区 Feed 增加"同温层"标签页，只看跟你同城/同行业的人

#### 3.2 失业趋势可视化

**定位**：强调"你不是一个人"，这是结构性问题，不是你的错。

展示数据：
- 平台内统计："你所在城市有 X 人跟你一样"
- 行业趋势："互联网行业本月新增 X 人加入"
- 命运共同体感："全国有 XX 万人正在经历跟你一样的事"

注意事项：
- 不展示具体失业率（避免过于负面）
- 强调"人数"而非"比例"
- 可展示"本周有 X 人上岸了"作为正向平衡

---

## 优先级排序

| 优先级 | 功能 | 工作量估算 | 价值 |
|--------|------|-----------|------|
| P0 | 邮箱验证 | 2-3 天 | 基础安全 |
| P0 | 自由天数计时器 | 1 天 | 低成本高辨识度 |
| P0 | 生存小贴士 | 1 天 | 情感连接 |
| P1 | 今日打卡 + 图片 | 3-4 天 | 日活留存 |
| P1 | 匿名模式 | 2-3 天 | 降低表达门槛 |
| P2 | 行业/城市同温层 | 5-7 天 | 社交差异化 |
| P2 | 失业趋势可视化 | 3-5 天 | 命运共同体感 |

---

## 技术要点

### 数据库扩展

```sql
-- 用户表新增字段
ALTER TABLE users ADD COLUMN last_work_date DATE;          -- 最后工作日
ALTER TABLE users ADD COLUMN city VARCHAR(50);             -- 城市
ALTER TABLE users ADD COLUMN industry VARCHAR(50);         -- 行业
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

-- 心情打卡表
CREATE TABLE daily_checkins (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    mood VARCHAR(20) NOT NULL,        -- sluggish/neutral/okay/energetic/breakthrough
    note TEXT,                         -- 一句话
    image_url VARCHAR(500),           -- 图片
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, DATE(created_at)) -- 每天只能打卡一次
);

-- 生存贴士表
CREATE TABLE survival_tips (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    category VARCHAR(50),             -- finance/life/mindset/policy
    created_at TIMESTAMP DEFAULT NOW()
);

-- 邮箱验证表
CREATE TABLE email_verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token VARCHAR(100) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 新增后端模块

```
server/internal/
├── checkin/          -- 心情打卡
│   ├── handler.go
│   ├── service.go
│   ├── repo.go
│   └── model.go
├── tips/            -- 生存贴士
│   ├── handler.go
│   ├── service.go
│   └── repo.go
└── verify/          -- 邮箱验证
    ├── handler.go
    ├── service.go
    └── repo.go
```

### 新增前端页面

```
web/app/
├── checkin/          -- 打卡页面
│   └── page.tsx
├── dashboard/        -- 个人仪表盘（自由天数 + 心情曲线）
│   └── page.tsx
└── community/        -- 同温层页面（后续）
    └── page.tsx
```

---

## 对现有架构的影响

- 用户模型扩展（新增字段，不破坏现有接口）
- 帖子模型扩展（新增 `is_anonymous` 字段）
- 新增 3 个独立模块（checkin / tips / verify），不侵入现有代码
- 前端新增 2-3 个页面 + 首页仪表盘组件

---

*讨论日期：2026-05-26*
*状态：讨论中，待确认后进入 design 阶段*
