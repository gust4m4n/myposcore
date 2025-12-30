-- Migration: Add audit fields to products table
-- This migration adds created_by, updated_by, and deleted_by fields to track user actions on products

-- Add created_by column (references users.id)
ALTER TABLE products ADD COLUMN IF NOT EXISTS created_by INTEGER;

-- Add updated_by column (references users.id)
ALTER TABLE products ADD COLUMN IF NOT EXISTS updated_by INTEGER;

-- Add deleted_by column (references users.id)
ALTER TABLE products ADD COLUMN IF NOT EXISTS deleted_by INTEGER;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_products_created_by ON products(created_by);
CREATE INDEX IF NOT EXISTS idx_products_updated_by ON products(updated_by);
CREATE INDEX IF NOT EXISTS idx_products_deleted_by ON products(deleted_by);

-- Add foreign key constraints (optional, uncomment if you want strict referential integrity)
-- Note: These constraints might fail if there are existing records with NULL values
-- ALTER TABLE products ADD CONSTRAINT fk_products_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE products ADD CONSTRAINT fk_products_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE products ADD CONSTRAINT fk_products_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Add comments to columns for documentation
COMMENT ON COLUMN products.created_by IS 'User ID who created this product record';
COMMENT ON COLUMN products.updated_by IS 'User ID who last updated this product record';
COMMENT ON COLUMN products.deleted_by IS 'User ID who deleted this product record (soft delete)';
