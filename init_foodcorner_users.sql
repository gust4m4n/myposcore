-- ============================================================================
-- Demo Users for Food Corner Tenant (Tenant ID: 17)
-- ============================================================================
-- This script creates demo users for Food Corner tenant
-- Password for all users: 123456
-- 
-- Tenant: Food Corner (ID: 17)
-- Branches: 
--   - Cabang Pusat (ID: 25)
--   - Cabang Menteng (ID: 26)
-- ============================================================================

-- Tenant Admin (can manage all branches)
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at) VALUES
(17, 25, 'admin@foodcorner.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Admin Food Corner', 'tenantadmin', true, NOW(), NOW());

-- Branch Admins
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at) VALUES
(17, 25, 'admin.pusat@foodcorner.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Admin Cabang Pusat', 'branchadmin', true, NOW(), NOW()),
(17, 26, 'admin.menteng@foodcorner.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Admin Cabang Menteng', 'branchadmin', true, NOW(), NOW());

-- Cashiers (regular users)
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at) VALUES
(17, 25, 'cashier.pusat@foodcorner.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Cashier Pusat', 'user', true, NOW(), NOW()),
(17, 26, 'cashier.menteng@foodcorner.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Cashier Menteng', 'user', true, NOW(), NOW());

-- ============================================================================
-- User Roles:
-- - superadmin: Full system access across all tenants
-- - tenantadmin: Can manage all branches within their tenant
-- - branchadmin: Can manage their specific branch
-- - user: Regular cashier/employee with basic access
-- ============================================================================
