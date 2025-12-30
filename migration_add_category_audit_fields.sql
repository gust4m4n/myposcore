-- Migration: Add audit tracking fields to categories table
-- Description: Adds created_by, updated_by, deleted_by columns for tracking user actions on categories
-- Author: System
-- Date: 2024

-- Step 1: Add audit tracking columns
ALTER TABLE categories ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Step 2: Add foreign key constraints to users table
ALTER TABLE categories 
ADD CONSTRAINT fk_categories_creator 
FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE categories 
ADD CONSTRAINT fk_categories_updater 
FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE categories 
ADD CONSTRAINT fk_categories_deleter 
FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Step 3: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_categories_created_by ON categories(created_by);
CREATE INDEX IF NOT EXISTS idx_categories_updated_by ON categories(updated_by);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_by ON categories(deleted_by);

-- Step 4: Add comments for documentation
COMMENT ON COLUMN categories.created_by IS 'User ID who created this category';
COMMENT ON COLUMN categories.updated_by IS 'User ID who last updated this category';
COMMENT ON COLUMN categories.deleted_by IS 'User ID who soft-deleted this category';

-- Rollback instructions:
-- To rollback this migration, run:
-- DROP INDEX IF EXISTS idx_categories_deleted_by;
-- DROP INDEX IF EXISTS idx_categories_updated_by;
-- DROP INDEX IF EXISTS idx_categories_created_by;
-- ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_categories_deleter;
-- ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_categories_updater;
-- ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_categories_creator;
-- ALTER TABLE categories DROP COLUMN IF EXISTS deleted_by;
-- ALTER TABLE categories DROP COLUMN IF EXISTS updated_by;
-- ALTER TABLE categories DROP COLUMN IF EXISTS created_by;
