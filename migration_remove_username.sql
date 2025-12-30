-- Migration to remove username column and add unique index on email
-- Run this after updating the code

-- Drop any existing unique constraint on username if exists
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS idx_users_username;

-- Drop the username column
ALTER TABLE users DROP COLUMN IF EXISTS username;

-- Add unique index on email if not exists
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Add comment to document the change
COMMENT ON COLUMN users.email IS 'Primary login identifier - must be unique across all users';
