-- Create tenant demo untuk testing
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Demo Tenant', 'TENANT001', true, NOW(), NOW());

-- Uncomment untuk create tenant lain jika diperlukan
-- INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
-- VALUES ('Tenant 2', 'TENANT002', true, NOW(), NOW());
