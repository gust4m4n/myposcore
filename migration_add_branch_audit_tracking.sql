-- Migration: Add audit tracking columns to branches table
-- This enables tracking of who created, updated, or deleted branch records

ALTER TABLE branches ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE branches ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE branches ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_branches_created_by ON branches(created_by);
CREATE INDEX IF NOT EXISTS idx_branches_updated_by ON branches(updated_by);
CREATE INDEX IF NOT EXISTS idx_branches_deleted_by ON branches(deleted_by);

-- Add foreign key constraints to users table
ALTER TABLE branches ADD CONSTRAINT fk_branches_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE branches ADD CONSTRAINT fk_branches_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE branches ADD CONSTRAINT fk_branches_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

-- Add comments for documentation
COMMENT ON COLUMN branches.created_by IS 'User ID who created this branch';
COMMENT ON COLUMN branches.updated_by IS 'User ID who last updated this branch';
COMMENT ON COLUMN branches.deleted_by IS 'User ID who deleted this branch';
