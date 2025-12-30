-- Migration: Add audit fields to users table
-- This migration adds created_by, updated_by, and deleted_by fields to track user actions

-- Add created_by column (references users.id)
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_by INTEGER;

-- Add updated_by column (references users.id)
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_by INTEGER;

-- Add deleted_by column (references users.id)
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);
CREATE INDEX IF NOT EXISTS idx_users_updated_by ON users(updated_by);
CREATE INDEX IF NOT EXISTS idx_users_deleted_by ON users(deleted_by);

-- Add foreign key constraints (optional, uncomment if you want strict referential integrity)
-- Note: These constraints might fail if there are existing records with NULL values
-- ALTER TABLE users ADD CONSTRAINT fk_users_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE users ADD CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE users ADD CONSTRAINT fk_users_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Add comments to columns for documentation
COMMENT ON COLUMN users.created_by IS 'User ID who created this user record';
COMMENT ON COLUMN users.updated_by IS 'User ID who last updated this user record';
COMMENT ON COLUMN users.deleted_by IS 'User ID who deleted this user record (soft delete)';
