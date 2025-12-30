-- Migration: Add index on full_name column in users table
-- Date: 2025-12-30
-- Description: Improves query performance for full name searches (non-unique, allows duplicates)

-- Create index on full_name column
CREATE INDEX IF NOT EXISTS idx_users_full_name ON users(full_name) WHERE deleted_at IS NULL;

-- Verify the index was created
SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'users' 
    AND indexname = 'idx_users_full_name';
