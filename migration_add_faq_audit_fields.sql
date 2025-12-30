-- Migration: Add audit tracking fields to faqs table
-- Description: Adds created_by, updated_by, deleted_by columns for tracking user actions on FAQs
-- Author: System
-- Date: 2024

-- Step 1: Add audit tracking columns
ALTER TABLE faqs ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE faqs ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE faqs ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Step 2: Add foreign key constraints to users table
ALTER TABLE faqs 
ADD CONSTRAINT fk_faqs_creator 
FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE faqs 
ADD CONSTRAINT fk_faqs_updater 
FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE faqs 
ADD CONSTRAINT fk_faqs_deleter 
FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Step 3: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_faqs_created_by ON faqs(created_by);
CREATE INDEX IF NOT EXISTS idx_faqs_updated_by ON faqs(updated_by);
CREATE INDEX IF NOT EXISTS idx_faqs_deleted_by ON faqs(deleted_by);

-- Step 4: Add comments for documentation
COMMENT ON COLUMN faqs.created_by IS 'User ID who created this FAQ';
COMMENT ON COLUMN faqs.updated_by IS 'User ID who last updated this FAQ';
COMMENT ON COLUMN faqs.deleted_by IS 'User ID who soft-deleted this FAQ';

-- Rollback instructions:
-- To rollback this migration, run:
-- DROP INDEX IF EXISTS idx_faqs_deleted_by;
-- DROP INDEX IF EXISTS idx_faqs_updated_by;
-- DROP INDEX IF EXISTS idx_faqs_created_by;
-- ALTER TABLE faqs DROP CONSTRAINT IF EXISTS fk_faqs_deleter;
-- ALTER TABLE faqs DROP CONSTRAINT IF EXISTS fk_faqs_updater;
-- ALTER TABLE faqs DROP CONSTRAINT IF EXISTS fk_faqs_creator;
-- ALTER TABLE faqs DROP COLUMN IF EXISTS deleted_by;
-- ALTER TABLE faqs DROP COLUMN IF EXISTS updated_by;
-- ALTER TABLE faqs DROP COLUMN IF EXISTS created_by;
