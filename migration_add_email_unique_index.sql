-- Migration: Add unique index on email column in users table
-- Date: 2025-12-30
-- Description: Ensures email uniqueness across all users globally

-- Drop existing index if it exists (in case it was created without UNIQUE)
DROP INDEX IF EXISTS idx_users_email;

-- Create unique index on email column
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(email) WHERE deleted_at IS NULL;

-- Verify the index was created
SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'users' 
    AND indexname = 'idx_users_email_unique';
