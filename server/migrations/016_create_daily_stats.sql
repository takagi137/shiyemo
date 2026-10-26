-- 016: 每日统计表（趋势可视化）
CREATE TABLE IF NOT EXISTS daily_stats (
    id BIGSERIAL PRIMARY KEY,
    stat_date VARCHAR(10) NOT NULL,
    dimension VARCHAR(20) NOT NULL,
    label VARCHAR(100) NOT NULL,
    count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_daily_stats_date ON daily_stats(stat_date);
CREATE INDEX IF NOT EXISTS idx_daily_stats_dimension ON daily_stats(dimension);
CREATE INDEX IF NOT EXISTS idx_daily_stats_label ON daily_stats(label);
CREATE INDEX IF NOT EXISTS idx_daily_stats_date_dim ON daily_stats(stat_date, dimension);
