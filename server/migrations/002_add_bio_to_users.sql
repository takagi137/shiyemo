-- Add bio column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio VARCHAR(500) DEFAULT '';

-- Comment
COMMENT ON COLUMN users.bio IS '个人简介';
