-- ========================================
-- TENANT 1: FOOD CORNER
-- ========================================

-- Create tenant Food Corner
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Food Corner', 'resto01', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create branches untuk Warung Makan
INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
SELECT id, 'Cabang Pusat', 'resto01-pusat', 'Jl. Sudirman No. 123, Jakarta', '021-12345678', true, NOW(), NOW()
FROM tenants WHERE code = 'resto01'
ON CONFLICT (tenant_id, code) DO NOTHING;

INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
SELECT id, 'Cabang Menteng', 'resto01-menteng', 'Jl. Menteng Raya No. 45, Jakarta', '021-87654321', true, NOW(), NOW()
FROM tenants WHERE code = 'resto01'
ON CONFLICT (tenant_id, code) DO NOTHING;

-- Create users untuk Warung Makan (password semua: 123456)
-- Tenant Admin
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'tenantadmin@resto.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Tenant Admin Resto', 'tenantadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'resto01' AND b.code = 'resto01-pusat'
ON CONFLICT (email) DO NOTHING;

-- Branch Admin Pusat
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'branchadmin.pusat@resto.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Branch Admin Pusat', 'branchadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'resto01' AND b.code = 'resto01-pusat'
ON CONFLICT (email) DO NOTHING;

-- Branch Admin Menteng
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'branchadmin.menteng@resto.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Branch Admin Menteng', 'branchadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'resto01' AND b.code = 'resto01-menteng'
ON CONFLICT (email) DO NOTHING;

-- Kasir Pusat
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'kasir.pusat@resto.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Kasir Pusat', 'user', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'resto01' AND b.code = 'resto01-pusat'
ON CONFLICT (email) DO NOTHING;

-- Kasir Menteng
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'kasir.menteng@resto.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Kasir Menteng', 'user', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'resto01' AND b.code = 'resto01-menteng'
ON CONFLICT (email) DO NOTHING;

-- ========================================
-- TENANT 2: FASHION STORE
-- ========================================

-- Create tenant Fashion Store
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Fashion Store', 'fashion01', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create branches untuk Fashion Store
INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
SELECT id, 'Cabang Mall Plaza', 'fashion01-plaza', 'Mall Plaza Lt. 2 No. 45, Jakarta', '021-55556666', true, NOW(), NOW()
FROM tenants WHERE code = 'fashion01'
ON CONFLICT (tenant_id, code) DO NOTHING;

INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
SELECT id, 'Cabang Grand Mall', 'fashion01-grand', 'Grand Mall Lt. 1 No. 78, Bandung', '022-77778888', true, NOW(), NOW()
FROM tenants WHERE code = 'fashion01'
ON CONFLICT (tenant_id, code) DO NOTHING;

-- Create users untuk Fashion Store (password semua: 123456)
-- Tenant Admin
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'tenantadmin@fashion.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Tenant Admin Fashion', 'tenantadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'fashion01' AND b.code = 'fashion01-plaza'
ON CONFLICT (email) DO NOTHING;

-- Branch Admin Plaza
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'branchadmin.plaza@fashion.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Branch Admin Plaza', 'branchadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'fashion01' AND b.code = 'fashion01-plaza'
ON CONFLICT (email) DO NOTHING;

-- Branch Admin Grand
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'branchadmin.grand@fashion.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Branch Admin Grand', 'branchadmin', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'fashion01' AND b.code = 'fashion01-grand'
ON CONFLICT (email) DO NOTHING;

-- Kasir Plaza
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'kasir.plaza@fashion.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Kasir Plaza', 'user', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'fashion01' AND b.code = 'fashion01-plaza'
ON CONFLICT (email) DO NOTHING;

-- Sales Plaza
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'sales.plaza@fashion.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Sales Plaza', 'user', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'fashion01' AND b.code = 'fashion01-plaza'
ON CONFLICT (email) DO NOTHING;

-- Kasir Grand
INSERT INTO users (tenant_id, branch_id, email, password, full_name, role, is_active, created_at, updated_at)
SELECT t.id, b.id, 'kasir.grand@fashion.com',
       '$2a$14$7dg5D./t2Un8.SFREKpxsu/nDt8v8oLWb.BKFQXUD0r2bdknIllF6',
       'Kasir Grand', 'user', true, NOW(), NOW()
FROM tenants t
JOIN branches b ON b.tenant_id = t.id
WHERE t.code = 'fashion01' AND b.code = 'fashion01-grand'
ON CONFLICT (tenant_id, username) DO NOTHING;
