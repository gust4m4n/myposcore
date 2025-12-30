-- Migration: Add audit tracking fields to payments table
-- Description: Adds created_by, updated_by, deleted_by columns for tracking user actions on payments
-- Author: System
-- Date: 2024

-- Step 1: Add audit tracking columns
ALTER TABLE payments ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Step 2: Add foreign key constraints to users table
ALTER TABLE payments 
ADD CONSTRAINT fk_payments_creator 
FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE payments 
ADD CONSTRAINT fk_payments_updater 
FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE payments 
ADD CONSTRAINT fk_payments_deleter 
FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Step 3: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_payments_created_by ON payments(created_by);
CREATE INDEX IF NOT EXISTS idx_payments_updated_by ON payments(updated_by);
CREATE INDEX IF NOT EXISTS idx_payments_deleted_by ON payments(deleted_by);

-- Step 4: Add comments for documentation
COMMENT ON COLUMN payments.created_by IS 'User ID who created this payment';
COMMENT ON COLUMN payments.updated_by IS 'User ID who last updated this payment';
COMMENT ON COLUMN payments.deleted_by IS 'User ID who soft-deleted this payment';

-- Rollback instructions:
-- To rollback this migration, run:
-- DROP INDEX IF EXISTS idx_payments_deleted_by;
-- DROP INDEX IF EXISTS idx_payments_updated_by;
-- DROP INDEX IF EXISTS idx_payments_created_by;
-- ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_deleter;
-- ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_updater;
-- ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_creator;
-- ALTER TABLE payments DROP COLUMN IF EXISTS deleted_by;
-- ALTER TABLE payments DROP COLUMN IF EXISTS updated_by;
-- ALTER TABLE payments DROP COLUMN IF EXISTS created_by;
