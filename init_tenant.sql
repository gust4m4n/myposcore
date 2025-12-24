-- Create super tenant untuk admin dan dashboard
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Super Admin', 'supertenant', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create tenant demo untuk testing
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Demo Tenant', 'TENANT001', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create super branch
INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
SELECT id, 'Super Branch', 'superbranch', 'Admin HQ', '-', true, NOW(), NOW()
FROM tenants WHERE code = 'supertenant'
ON CONFLICT DO NOTHING;

-- Create superuser (username: superadmin, password: 123456, role: superadmin)
INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'superadmin', 'admin@myposcore.com', 
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Super Admin', 'superadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'supertenant' AND b.code = 'superbranch'
ON CONFLICT DO NOTHING;

-- Uncomment untuk create tenant lain jika diperlukan
-- INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
-- VALUES ('Tenant 2', 'TENANT002', true, NOW(), NOW());
