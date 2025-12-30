-- Migration: Add audit tracking fields to orders table
-- Description: Adds created_by, updated_by, deleted_by columns for tracking user actions on orders
-- Author: System
-- Date: 2024

-- Step 1: Add audit tracking columns
ALTER TABLE orders ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Step 2: Add foreign key constraints to users table
ALTER TABLE orders 
ADD CONSTRAINT fk_orders_creator 
FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE orders 
ADD CONSTRAINT fk_orders_updater 
FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE orders 
ADD CONSTRAINT fk_orders_deleter 
FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Step 3: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_orders_created_by ON orders(created_by);
CREATE INDEX IF NOT EXISTS idx_orders_updated_by ON orders(updated_by);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_by ON orders(deleted_by);

-- Step 4: Add comments for documentation
COMMENT ON COLUMN orders.created_by IS 'User ID who created this order';
COMMENT ON COLUMN orders.updated_by IS 'User ID who last updated this order';
COMMENT ON COLUMN orders.deleted_by IS 'User ID who soft-deleted this order';

-- Rollback instructions:
-- To rollback this migration, run:
-- DROP INDEX IF EXISTS idx_orders_deleted_by;
-- DROP INDEX IF EXISTS idx_orders_updated_by;
-- DROP INDEX IF EXISTS idx_orders_created_by;
-- ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_deleter;
-- ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_updater;
-- ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_creator;
-- ALTER TABLE orders DROP COLUMN IF EXISTS deleted_by;
-- ALTER TABLE orders DROP COLUMN IF EXISTS updated_by;
-- ALTER TABLE orders DROP COLUMN IF EXISTS created_by;
