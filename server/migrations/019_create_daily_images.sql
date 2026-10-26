-- 创建每日风景图表
CREATE TABLE IF NOT EXISTS daily_images (
    id BIGSERIAL PRIMARY KEY,
    url VARCHAR(500) NOT NULL,
    description VARCHAR(200) NOT NULL DEFAULT '',
    credit VARCHAR(200) NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_daily_images_is_active ON daily_images(is_active);
CREATE INDEX idx_daily_images_sort_order ON daily_images(sort_order);

-- 插入初始风景图（从硬编码迁移过来）
INSERT INTO daily_images (url, description, credit, sort_order) VALUES
('https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=1920&q=80', '山谷湖泊', 'Photo by Baechtold on Unsplash', 1),
('https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=1920&q=80', '晨雾森林', 'Photo by Baechtold on Unsplash', 2),
('https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=1920&q=80', '林间小径', 'Photo by Lucian on Unsplash', 3),
('https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1920&q=80', '海滩日落', 'Photo by Sean O. on Unsplash', 4),
('https://images.unsplash.com/photo-1519681393784-d120267933ba?w=1920&q=80', '星空山脉', 'Photo by Benjamin Voros on Unsplash', 5),
('https://images.unsplash.com/photo-1472214103451-9374bd1c798e?w=1920&q=80', '绿色山谷', 'Photo by Robert Lukeman on Unsplash', 6),
('https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=1920&q=80', '山间湖泊', 'Photo by Ray Hennessy on Unsplash', 7),
('https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=1920&q=80', '雪山风光', 'Photo by Luca Bravo on Unsplash', 8),
('https://images.unsplash.com/photo-1500534314209-a25ddb2bd429?w=1920&q=80', '田野日出', 'Photo by Lukasz Szmigiel on Unsplash', 9),
('https://images.unsplash.com/photo-1495616811223-4d98c6e9c869?w=1920&q=80', '海边岩石', 'Photo by sampling via unsplash', 10);
