-- 为用户表添加同温层匹配所需字段
-- City: 城市
-- Industry: 行业
-- LeftJobAt: 离职时间

ALTER TABLE users ADD COLUMN IF NOT EXISTS city VARCHAR(50) DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS industry VARCHAR(100) DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS left_job_at TIMESTAMP;

-- 为城市和行业字段创建索引，用于同温层聚合查询
CREATE INDEX IF NOT EXISTS idx_users_city ON users(city);
CREATE INDEX IF NOT EXISTS idx_users_industry ON users(industry);
