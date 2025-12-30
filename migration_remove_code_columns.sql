-- Migration: Remove code columns from tenants and branches tables
-- Date: 2025-12-30
-- Description: Remove tenant_code and branch_code as system now uses IDs only

-- Drop indexes first
DROP INDEX IF EXISTS idx_tenant_code;
DROP INDEX IF EXISTS idx_branch_code;

-- Drop code column from tenants table
ALTER TABLE tenants DROP COLUMN IF EXISTS code;

-- Drop code column from branches table
ALTER TABLE branches DROP COLUMN IF EXISTS code;

-- Verify columns were dropped
SELECT 
    table_name, 
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_name IN ('tenants', 'branches')
    AND column_name = 'code';
