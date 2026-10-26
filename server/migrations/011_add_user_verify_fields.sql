-- 补齐 users 表的 Phase 2 字段
--
-- 背景：internal/auth/model.go 的 User 结构体自 Phase 2 起就含有
-- EmailVerified 和 LastWorkDate 两个字段，但从未有对应迁移文件，
-- 导致全新部署上邮箱注册直接失败：
--   ERROR: column "email_verified" of relation "users" does not exist
-- （旧服务器上应为手工 ALTER TABLE 添加，未落入版本控制）
--
-- 编号 011 此前为空缺号，此处补上。

ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_work_date TIMESTAMP;

COMMENT ON COLUMN users.email_verified IS '邮箱是否已验证';
COMMENT ON COLUMN users.last_work_date IS '最后工作日，用于计算自由天数';
