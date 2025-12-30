-- Migration: Add audit tracking columns to tenants table
-- This enables tracking of who created, updated, or deleted tenant records

ALTER TABLE tenants ADD COLUMN IF NOT EXISTS created_by INTEGER;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS updated_by INTEGER;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_tenants_created_by ON tenants(created_by);
CREATE INDEX IF NOT EXISTS idx_tenants_updated_by ON tenants(updated_by);
CREATE INDEX IF NOT EXISTS idx_tenants_deleted_by ON tenants(deleted_by);

-- Add foreign key constraints to users table
ALTER TABLE tenants ADD CONSTRAINT fk_tenants_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE tenants ADD CONSTRAINT fk_tenants_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE tenants ADD CONSTRAINT fk_tenants_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

-- Add comments for documentation
COMMENT ON COLUMN tenants.created_by IS 'User ID who created this tenant';
COMMENT ON COLUMN tenants.updated_by IS 'User ID who last updated this tenant';
COMMENT ON COLUMN tenants.deleted_by IS 'User ID who deleted this tenant';
