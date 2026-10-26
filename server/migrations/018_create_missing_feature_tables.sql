-- 补齐四张缺失的功能表
--
-- 背景：以下四个模块的 Go 模型（含 TableName()）早已存在、路由已注册、
-- 前端也有对应入口，但从未有建表迁移。全新部署上这四个功能全部不可用：
--   daily_checkins      → internal/checkin      → 导航栏「打卡」页
--   email_verifications → internal/verify       → 邮箱验证页
--   sensitive_words     → internal/moderation   → 管理后台「敏感词管理」
--   survival_tips       → internal/tips         → 首页 DailyTip 组件
--
-- 编号 018 此前为空缺号，此处补上。字段定义严格对应各模块 model.go。

-- ── 每日打卡（internal/checkin/model.go: DailyCheckin）──
CREATE TABLE IF NOT EXISTS daily_checkins (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    mood       VARCHAR(20)  NOT NULL,          -- happy/calm/sad/angry/anxious
    content    TEXT,
    images     TEXT,                            -- JSON 数组，图片 URL 列表
    is_public  BOOLEAN      DEFAULT TRUE,
    created_at TIMESTAMPTZ  DEFAULT now(),
    updated_at TIMESTAMPTZ  DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_daily_checkins_user_id ON daily_checkins(user_id);
CREATE INDEX IF NOT EXISTS idx_daily_checkins_created_at ON daily_checkins(created_at DESC);
-- 公开打卡时间线查询：WHERE is_public = true ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_daily_checkins_public_created ON daily_checkins(is_public, created_at DESC);

-- ── 邮箱验证（internal/verify/model.go: EmailVerification）──
CREATE TABLE IF NOT EXISTS email_verifications (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    token       VARCHAR(100) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    verified_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_email_verifications_user_id ON email_verifications(user_id);

-- ── 敏感词（internal/moderation/model.go: SensitiveWord）──
CREATE TABLE IF NOT EXISTS sensitive_words (
    id   BIGSERIAL PRIMARY KEY,
    word VARCHAR(100) NOT NULL UNIQUE
);

-- ── 生存小贴士（internal/tips/model.go: SurvivalTip）──
CREATE TABLE IF NOT EXISTS survival_tips (
    id         BIGSERIAL PRIMARY KEY,
    content    TEXT        NOT NULL,
    category   VARCHAR(50) DEFAULT 'general',
    is_active  BOOLEAN     DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT now()
);
-- 随机取一条启用贴士：WHERE is_active = true
CREATE INDEX IF NOT EXISTS idx_survival_tips_active ON survival_tips(is_active);
