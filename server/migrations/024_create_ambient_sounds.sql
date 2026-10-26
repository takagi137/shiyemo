-- 023_create_ambient_sounds.sql
-- 静思空间「环境音」改为后端维护列表（对齐 daily_images / readings 框架）：
-- 后端表存元数据 + 外部下载地址，前端 /api/v1/sounds 取列表，不再硬编码、
-- 也不把 mp3 二进制放进仓库/镜像（与 daily_images 用 Unsplash 外链同理）。
-- 音源均为 CC0/公共领域可商用：
--   rain/ocean/forest/fire/whitenoise → BigSoundBank (CC0)
--   piano → OpenGameArt "Calm Piano 1 / Vaporware", The Cynic Project (CC0)

CREATE TABLE IF NOT EXISTS ambient_sounds (
    id         BIGSERIAL PRIMARY KEY,
    key        VARCHAR(50)  NOT NULL,
    name       VARCHAR(50)  NOT NULL DEFAULT '',
    emoji      VARCHAR(16)  NOT NULL DEFAULT '',
    url        VARCHAR(500) NOT NULL,
    credit     VARCHAR(200) NOT NULL DEFAULT '',
    sort_order INT          NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ambient_sounds_is_active ON ambient_sounds(is_active);
CREATE INDEX IF NOT EXISTS idx_ambient_sounds_sort_order ON ambient_sounds(sort_order);

-- 幂等：仅当表为空时插入初始 6 个音效
INSERT INTO ambient_sounds (key, name, emoji, url, credit, sort_order)
SELECT v.key, v.name, v.emoji, v.url, v.credit, v.sort_order
FROM (VALUES
    ('rain',       '雨声',   '🌧️', 'https://bigsoundbank.com/UPLOAD/mp3/1019.mp3',                'BigSoundBank (CC0)',                   1),
    ('piano',      '钢琴',   '🎹', 'https://opengameart.org/sites/default/files/003_Vaporware_2.mp3', 'OpenGameArt / The Cynic Project (CC0)', 2),
    ('forest',     '森林',   '🌲', 'https://bigsoundbank.com/UPLOAD/mp3/0100.mp3',                'BigSoundBank (CC0)',                   3),
    ('ocean',      '海浪',   '🌊', 'https://bigsoundbank.com/UPLOAD/mp3/2565.mp3',                'BigSoundBank (CC0)',                   4),
    ('fire',       '篝火',   '🔥', 'https://bigsoundbank.com/UPLOAD/mp3/0031.mp3',                'BigSoundBank (CC0)',                   5),
    ('whitenoise', '白噪音', '☁️', 'https://bigsoundbank.com/UPLOAD/mp3/1037.mp3',                'BigSoundBank (CC0)',                   6)
) AS v(key, name, emoji, url, credit, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM ambient_sounds);
