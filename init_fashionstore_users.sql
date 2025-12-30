-- ============================================================================
-- Demo Users for Fashion Store Tenant (Tenant ID: 18)
-- ============================================================================
-- This script creates demo users for Fashion Store tenant
-- Password for all users: 123456
-- 
-- Tenant: Fashion Store (ID: 18)
-- Branches: 
--   - Cabang Mall Plaza (ID: 29)
--   - Cabang Grand Mall (ID: 30)
-- ============================================================================

-- Branch Admins
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at) VALUES
(18, 29, 'admin.plaza@fashionstore.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Admin Mall Plaza', 'branchadmin', true, NOW(), NOW()),
(18, 30, 'admin.grand@fashionstore.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Admin Grand Mall', 'branchadmin', true, NOW(), NOW()),
(18, 30, 'admin.senayan@fashionstore.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Admin Senayan', 'branchadmin', true, NOW(), NOW());

-- ============================================================================
-- User Roles:
-- - superadmin: Full system access across all tenants
-- - tenantadmin: Can manage all branches within their tenant
-- - branchadmin: Can manage their specific branch
-- - user: Regular cashier/employee with basic access
-- ============================================================================
