-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    company VARCHAR(100) DEFAULT '',
    occupation VARCHAR(100) DEFAULT '',
    avatar VARCHAR(500) DEFAULT '',
    status SMALLINT DEFAULT 1,  -- 1: active, 0: disabled
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_nickname ON users(nickname);
CREATE INDEX idx_users_status ON users(status);

-- Comment
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.phone IS '手机号';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.password IS '密码哈希';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.company IS '前公司';
COMMENT ON COLUMN users.occupation IS '职业';
COMMENT ON COLUMN users.avatar IS '头像URL';
COMMENT ON COLUMN users.status IS '状态: 1-正常, 0-禁用';
