-- Migration: Add tenant_id and branch_id to payments table
-- This migration handles existing payment records by deriving tenant_id and branch_id from their orders

-- Step 1: Add columns as nullable first
ALTER TABLE payments 
ADD COLUMN IF NOT EXISTS tenant_id bigint,
ADD COLUMN IF NOT EXISTS branch_id bigint;

-- Step 2: Update existing payments with tenant_id and branch_id from their orders
UPDATE payments p
SET 
    tenant_id = o.tenant_id,
    branch_id = o.branch_id
FROM orders o
WHERE p.order_id = o.id
AND p.tenant_id IS NULL;

-- Step 3: For any payments without matching orders, set to default tenant/branch (if any exist)
-- Get the first available tenant and branch
DO $$
DECLARE
    default_tenant_id bigint;
    default_branch_id bigint;
BEGIN
    -- Get first tenant
    SELECT id INTO default_tenant_id FROM tenants WHERE is_active = true LIMIT 1;
    
    -- Get first branch
    SELECT id INTO default_branch_id FROM branches WHERE is_active = true LIMIT 1;
    
    -- Update payments that still don't have tenant_id/branch_id
    IF default_tenant_id IS NOT NULL AND default_branch_id IS NOT NULL THEN
        UPDATE payments 
        SET 
            tenant_id = default_tenant_id,
            branch_id = default_branch_id
        WHERE tenant_id IS NULL OR branch_id IS NULL;
    END IF;
END $$;

-- Step 4: Now make columns NOT NULL and add constraints
ALTER TABLE payments 
ALTER COLUMN tenant_id SET NOT NULL,
ALTER COLUMN branch_id SET NOT NULL;

-- Step 5: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_branch_id ON payments(branch_id);

-- Step 6: Add foreign key constraints
ALTER TABLE payments 
DROP CONSTRAINT IF EXISTS fk_payments_tenant,
ADD CONSTRAINT fk_payments_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE payments 
DROP CONSTRAINT IF EXISTS fk_payments_branch,
ADD CONSTRAINT fk_payments_branch FOREIGN KEY (branch_id) REFERENCES branches(id);

-- Add comments
COMMENT ON COLUMN payments.tenant_id IS 'Tenant ID for multi-tenancy support';
COMMENT ON COLUMN payments.branch_id IS 'Branch ID for branch-level data isolation';
