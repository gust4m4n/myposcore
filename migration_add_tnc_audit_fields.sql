-- Migration: Add audit tracking fields to terms_and_conditions table
-- Description: Adds created_by, updated_by, deleted_by columns for tracking user actions on TnC
-- Author: System
-- Date: 2024

-- Step 1: Add audit tracking columns
ALTER TABLE terms_and_conditions ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE terms_and_conditions ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE terms_and_conditions ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Step 2: Add foreign key constraints to users table
ALTER TABLE terms_and_conditions 
ADD CONSTRAINT fk_tnc_creator 
FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE terms_and_conditions 
ADD CONSTRAINT fk_tnc_updater 
FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE terms_and_conditions 
ADD CONSTRAINT fk_tnc_deleter 
FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Step 3: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_tnc_created_by ON terms_and_conditions(created_by);
CREATE INDEX IF NOT EXISTS idx_tnc_updated_by ON terms_and_conditions(updated_by);
CREATE INDEX IF NOT EXISTS idx_tnc_deleted_by ON terms_and_conditions(deleted_by);

-- Step 4: Add comments for documentation
COMMENT ON COLUMN terms_and_conditions.created_by IS 'User ID who created this TnC';
COMMENT ON COLUMN terms_and_conditions.updated_by IS 'User ID who last updated this TnC';
COMMENT ON COLUMN terms_and_conditions.deleted_by IS 'User ID who soft-deleted this TnC';

-- Rollback instructions:
-- To rollback this migration, run:
-- DROP INDEX IF EXISTS idx_tnc_deleted_by;
-- DROP INDEX IF EXISTS idx_tnc_updated_by;
-- DROP INDEX IF EXISTS idx_tnc_created_by;
-- ALTER TABLE terms_and_conditions DROP CONSTRAINT IF EXISTS fk_tnc_deleter;
-- ALTER TABLE terms_and_conditions DROP CONSTRAINT IF EXISTS fk_tnc_updater;
-- ALTER TABLE terms_and_conditions DROP CONSTRAINT IF EXISTS fk_tnc_creator;
-- ALTER TABLE terms_and_conditions DROP COLUMN IF EXISTS deleted_by;
-- ALTER TABLE terms_and_conditions DROP COLUMN IF EXISTS updated_by;
-- ALTER TABLE terms_and_conditions DROP COLUMN IF EXISTS created_by;
