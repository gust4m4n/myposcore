-- Migration: Add indexes for tenant code and branch code
-- Date: 2025-12-27
-- Description: Add unique indexes for tenant.code and branch.code for better query performance

-- Create unique index for tenant code if not exists
CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_code ON tenants(code) WHERE deleted_at IS NULL;

-- Create unique index for branch code if not exists
CREATE UNIQUE INDEX IF NOT EXISTS idx_branch_code ON branches(code) WHERE deleted_at IS NULL;

-- Create index for branch tenant_id if not exists (for faster joins)
CREATE INDEX IF NOT EXISTS idx_branch_tenant ON branches(tenant_id) WHERE deleted_at IS NULL;

-- Verify indexes
SELECT 
    tablename,
    indexname,
    indexdef
FROM 
    pg_indexes
WHERE 
    tablename IN ('tenants', 'branches')
    AND indexname IN ('idx_tenant_code', 'idx_branch_code', 'idx_branch_tenant')
ORDER BY 
    tablename, indexname;
