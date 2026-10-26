-- 008_add_notification_type.sql
-- 添加通知类型字段

ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(20) DEFAULT 'system';

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- 更新现有通知的类型
UPDATE notifications SET type = 'system' WHERE type IS NULL;
