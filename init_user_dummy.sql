-- Dummy data for users table
-- Make sure you have tenant and branch data first

-- Insert dummy users for tenant_id = 17, branch_id = 1
INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at) VALUES
(17, 1, 'user1', 'user1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User Satu', 'user', true, NOW(), NOW()),
(17, 1, 'user2', 'user2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User Dua', 'user', true, NOW(), NOW()),
(17, 1, 'branchadmin1', 'branchadmin1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Branch Admin Satu', 'branchadmin', true, NOW(), NOW()),
(17, 1, 'tenantadmin1', 'tenantadmin1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Tenant Admin Satu', 'tenantadmin', true, NOW(), NOW()),
(17, 1, 'cashier1', 'cashier1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Satu', 'user', true, NOW(), NOW()),
(17, 1, 'cashier2', 'cashier2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Dua', 'user', true, NOW(), NOW()),
(17, 1, 'inactiveuser', 'inactive@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Inactive User', 'user', false, NOW(), NOW());

-- Note: All passwords above are hashed version of "password123"
-- To use: login with username and password "password123"

-- Example for tenant_id = 18, branch_id = 2 (if exists)
INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at) VALUES
(18, 2, 'user1', 'user1.tenant18@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User Satu Tenant 18', 'user', true, NOW(), NOW()),
(18, 2, 'branchadmin1', 'branchadmin1.tenant18@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Branch Admin Tenant 18', 'branchadmin', true, NOW(), NOW());

-- Query to verify
SELECT id, tenant_id, branch_id, username, email, full_name, role, is_active 
FROM users 
WHERE tenant_id = 17
ORDER BY created_at DESC;
