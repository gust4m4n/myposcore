-- Migration: Add indexes on name columns for better search performance
-- Date: 2026-01-01
-- Description: Add indexes on name columns in users, products, tenants, branches, and categories tables

-- Add index on products.name (if not exists)
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name) WHERE deleted_at IS NULL;

-- Add index on tenants.name (if not exists)
CREATE INDEX IF NOT EXISTS idx_tenants_name ON tenants(name) WHERE deleted_at IS NULL;

-- Add index on branches.name (if not exists)
CREATE INDEX IF NOT EXISTS idx_branches_name ON branches(name) WHERE deleted_at IS NULL;

-- Add index on categories.name (if not exists)
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name) WHERE deleted_at IS NULL;

-- Note: users.full_name index already exists in migration_add_fullname_index.sql

-- Verify the indexes were created
SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename IN ('products', 'tenants', 'branches', 'categories', 'users')
    AND indexname IN ('idx_products_name', 'idx_tenants_name', 'idx_branches_name', 'idx_categories_name', 'idx_users_full_name')
ORDER BY tablename, indexname;
