-- 系统设置表（键值存储）
--
-- 用途：让管理员在后台动态配置 SMTP 等运行时参数，无需改环境变量重启容器。
-- 读取优先级：数据库 > 环境变量 > 默认值（见 internal/setting/service.go）。
--
-- 设计为通用 key-value 而非 SMTP 专表，后续其他可配置项可直接复用。

CREATE TABLE IF NOT EXISTS system_settings (
    key        VARCHAR(100) PRIMARY KEY,
    value      TEXT         NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

COMMENT ON TABLE  system_settings       IS '系统运行时设置，管理员后台可编辑';
COMMENT ON COLUMN system_settings.key   IS '设置键，如 smtp.host / smtp.enabled';
COMMENT ON COLUMN system_settings.value IS '设置值，统一以字符串存储，由 service 层解析';
