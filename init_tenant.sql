-- Create super tenant untuk admin dan dashboard
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Super Admin', 'SUPERTENANT', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create tenant demo untuk testing
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Demo Tenant', 'TENANT001', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create super branch
INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
SELECT id, 'Super Branch', 'SUPERBRANCH', 'Admin HQ', '-', true, NOW(), NOW()
FROM tenants WHERE code = 'SUPERTENANT'
ON CONFLICT DO NOTHING;

-- Create superuser (password: superuser, role: superadmin)
INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'superuser', 'admin@myposcore.com', 
       '$2a$10$R4Je8Eodch/6IQUTsQG6wupEcX0qFlJZdzQS5Sp9vZ90Kvu/aTb8a',
       'Super Admin', 'superadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'SUPERTENANT' AND b.code = 'SUPERBRANCH'
ON CONFLICT DO NOTHING;

-- Uncomment untuk create tenant lain jika diperlukan
-- INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
-- VALUES ('Tenant 2', 'TENANT002', true, NOW(), NOW());
