-- Demo Tenants Data
-- Run this after init_tenant.sql

-- ============================================
-- TENANT 1: RESTORAN "WARUNG MAKAN SEJAHTERA"
-- ============================================

-- Insert Tenant 1: Restaurant
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Warung Makan Sejahtera', 'resto01', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Get Tenant 1 ID
DO $$
DECLARE
    tenant1_id INT;
    branch1_id INT;
    branch2_id INT;
BEGIN
    SELECT id INTO tenant1_id FROM tenants WHERE code = 'resto01';
    
    -- Insert Branches for Restaurant
    INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
    VALUES 
        (tenant1_id, 'Cabang Pusat', 'resto01-pusat', 'Jl. Sudirman No. 123, Jakarta', '021-12345678', true, NOW(), NOW()),
        (tenant1_id, 'Cabang Menteng', 'resto01-menteng', 'Jl. Menteng Raya No. 45, Jakarta', '021-87654321', true, NOW(), NOW())
    ON CONFLICT (code) DO NOTHING;
    
    -- Get Branch IDs
    SELECT id INTO branch1_id FROM branches WHERE code = 'resto01-pusat';
    SELECT id INTO branch2_id FROM branches WHERE code = 'resto01-menteng';
    
    -- Insert Default Tenant Admin (Username: tenantadmin, Password: 123456)
    IF NOT EXISTS (SELECT 1 FROM users WHERE tenant_id = tenant1_id AND username = 'tenantadmin') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant1_id, branch1_id, 'tenantadmin', 'tenantadmin@resto.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Tenant Admin Resto', 'tenantadmin', true, NOW(), NOW());
    END IF;
    
    -- Insert Default Branch Admins for each branch (Username: branchadmin, Password: 123456)
    IF NOT EXISTS (SELECT 1 FROM users WHERE branch_id = branch1_id AND username = 'branchadmin') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant1_id, branch1_id, 'branchadmin', 'branchadmin.pusat@resto.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Branch Admin Pusat', 'branchadmin', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE branch_id = branch2_id AND username = 'branchadmin') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant1_id, branch2_id, 'branchadmin', 'branchadmin.menteng@resto.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Branch Admin Menteng', 'branchadmin', true, NOW(), NOW());
    END IF;
    
    -- Insert Users for Restaurant (Password: demo123)
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin_resto') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant1_id, branch1_id, 'admin_resto', 'admin@resto.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Admin Resto', 'tenantadmin', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'kasir_pusat') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant1_id, branch1_id, 'kasir_pusat', 'kasir.pusat@resto.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Kasir Pusat', 'user', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'kasir_menteng') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant1_id, branch2_id, 'kasir_menteng', 'kasir.menteng@resto.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Kasir Menteng', 'user', true, NOW(), NOW());
    END IF;
    
    -- Insert Products for Restaurant (35 products in 3 categories)
    IF NOT EXISTS (SELECT 1 FROM products WHERE sku = 'RESTO-MKN-001') THEN
        INSERT INTO products (tenant_id, name, description, sku, price, stock, is_active, created_at, updated_at)
        VALUES 
            -- KATEGORI 1: MAKANAN UTAMA (15 items)
            (tenant1_id, 'Nasi Goreng Spesial', 'Nasi goreng dengan telur, ayam, dan udang', 'RESTO-MKN-001', 35000, 100, true, NOW(), NOW()),
            (tenant1_id, 'Mie Goreng Seafood', 'Mie goreng dengan seafood segar', 'RESTO-MKN-002', 32000, 100, true, NOW(), NOW()),
            (tenant1_id, 'Ayam Bakar', 'Ayam bakar bumbu kecap dengan lalapan', 'RESTO-MKN-003', 28000, 50, true, NOW(), NOW()),
            (tenant1_id, 'Soto Ayam', 'Soto ayam dengan bumbu tradisional', 'RESTO-MKN-004', 25000, 80, true, NOW(), NOW()),
            (tenant1_id, 'Gado-gado', 'Sayuran dengan bumbu kacang', 'RESTO-MKN-005', 20000, 60, true, NOW(), NOW()),
            (tenant1_id, 'Sate Ayam', 'Sate ayam 10 tusuk dengan bumbu kacang', 'RESTO-MKN-006', 30000, 70, true, NOW(), NOW()),
            (tenant1_id, 'Cap Cay', 'Tumis sayuran aneka warna', 'RESTO-MKN-007', 28000, 50, true, NOW(), NOW()),
            (tenant1_id, 'Nasi Uduk Komplit', 'Nasi uduk dengan ayam goreng, telur, tempe', 'RESTO-MKN-008', 27000, 80, true, NOW(), NOW()),
            (tenant1_id, 'Rawon Daging', 'Rawon daging sapi dengan keluak', 'RESTO-MKN-009', 38000, 40, true, NOW(), NOW()),
            (tenant1_id, 'Rendang Sapi', 'Rendang sapi Padang original', 'RESTO-MKN-010', 42000, 35, true, NOW(), NOW()),
            (tenant1_id, 'Nasi Pecel', 'Nasi dengan sayuran dan bumbu pecel', 'RESTO-MKN-011', 18000, 90, true, NOW(), NOW()),
            (tenant1_id, 'Nasi Kuning Komplit', 'Nasi kuning dengan lauk lengkap', 'RESTO-MKN-012', 30000, 60, true, NOW(), NOW()),
            (tenant1_id, 'Ayam Geprek', 'Ayam goreng dengan sambal pedas geprek', 'RESTO-MKN-013', 26000, 75, true, NOW(), NOW()),
            (tenant1_id, 'Ikan Bakar Kecap', 'Ikan bakar dengan bumbu kecap manis', 'RESTO-MKN-014', 35000, 45, true, NOW(), NOW()),
            (tenant1_id, 'Beef Teriyaki', 'Beef teriyaki dengan saus Jepang', 'RESTO-MKN-015', 45000, 30, true, NOW(), NOW()),
            
            -- KATEGORI 2: MINUMAN (10 items)
            (tenant1_id, 'Es Teh Manis', 'Teh manis dingin', 'RESTO-MNM-001', 5000, 200, true, NOW(), NOW()),
            (tenant1_id, 'Es Jeruk', 'Jus jeruk segar', 'RESTO-MNM-002', 8000, 150, true, NOW(), NOW()),
            (tenant1_id, 'Es Kelapa Muda', 'Kelapa muda segar', 'RESTO-MNM-003', 12000, 100, true, NOW(), NOW()),
            (tenant1_id, 'Kopi Susu', 'Kopi susu hangat', 'RESTO-MNM-004', 10000, 120, true, NOW(), NOW()),
            (tenant1_id, 'Jus Alpukat', 'Jus alpukat dengan susu', 'RESTO-MNM-005', 15000, 80, true, NOW(), NOW()),
            (tenant1_id, 'Es Campur', 'Es campur dengan buah dan agar-agar', 'RESTO-MNM-006', 18000, 70, true, NOW(), NOW()),
            (tenant1_id, 'Jus Mangga', 'Jus mangga segar tanpa gula', 'RESTO-MNM-007', 13000, 90, true, NOW(), NOW()),
            (tenant1_id, 'Teh Tarik', 'Teh tarik ala Malaysia', 'RESTO-MNM-008', 9000, 110, true, NOW(), NOW()),
            (tenant1_id, 'Es Cincau', 'Es cincau hitam dengan sirup', 'RESTO-MNM-009', 7000, 130, true, NOW(), NOW()),
            (tenant1_id, 'Cappuccino', 'Cappuccino dengan foam susu', 'RESTO-MNM-010', 16000, 85, true, NOW(), NOW()),
            
            -- KATEGORI 3: SNACK & DESSERT (10 items)
            (tenant1_id, 'Pisang Goreng', 'Pisang goreng krispy', 'RESTO-SNK-001', 12000, 100, true, NOW(), NOW()),
            (tenant1_id, 'Kentang Goreng', 'French fries dengan saus', 'RESTO-SNK-002', 15000, 120, true, NOW(), NOW()),
            (tenant1_id, 'Lumpia Semarang', 'Lumpia rebung isi sayuran', 'RESTO-SNK-003', 18000, 80, true, NOW(), NOW()),
            (tenant1_id, 'Martabak Manis Mini', 'Martabak manis coklat keju', 'RESTO-SNK-004', 22000, 50, true, NOW(), NOW()),
            (tenant1_id, 'Onion Rings', 'Onion rings crispy', 'RESTO-SNK-005', 16000, 90, true, NOW(), NOW()),
            (tenant1_id, 'Batagor', 'Batagor bandung 5 pcs', 'RESTO-SNK-006', 20000, 70, true, NOW(), NOW()),
            (tenant1_id, 'Es Krim Goreng', 'Es krim goreng dengan taburan', 'RESTO-SNK-007', 25000, 40, true, NOW(), NOW()),
            (tenant1_id, 'Risoles Mayo', 'Risoles isi sayur mayo 4 pcs', 'RESTO-SNK-008', 18000, 65, true, NOW(), NOW()),
            (tenant1_id, 'Pudding Buah', 'Pudding dengan potongan buah segar', 'RESTO-SNK-009', 14000, 75, true, NOW(), NOW()),
            (tenant1_id, 'Klepon', 'Klepon gula merah 6 pcs', 'RESTO-SNK-010', 10000, 85, true, NOW(), NOW());
    END IF;
END $$;

-- ============================================
-- TENANT 2: TOKO BAJU "FASHION STORE"
-- ============================================

-- Insert Tenant 2: Clothing Store
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Fashion Store', 'fashion01', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Get Tenant 2 ID
DO $$
DECLARE
    tenant2_id INT;
    branch3_id INT;
    branch4_id INT;
BEGIN
    SELECT id INTO tenant2_id FROM tenants WHERE code = 'fashion01';
    
    -- Insert Branches for Fashion Store
    INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
    VALUES 
        (tenant2_id, 'Cabang Mall Plaza', 'fashion01-plaza', 'Mall Plaza Lt. 2 No. 45, Jakarta', '021-55556666', true, NOW(), NOW()),
        (tenant2_id, 'Cabang Grand Mall', 'fashion01-grand', 'Grand Mall Lt. 1 No. 78, Bandung', '022-77778888', true, NOW(), NOW())
    ON CONFLICT (code) DO NOTHING;
    
    -- Get Branch IDs
    SELECT id INTO branch3_id FROM branches WHERE code = 'fashion01-plaza';
    SELECT id INTO branch4_id FROM branches WHERE code = 'fashion01-grand';
    
    -- Insert Default Tenant Admin (Username: tenantadmin, Password: 123456)
    IF NOT EXISTS (SELECT 1 FROM users WHERE tenant_id = tenant2_id AND username = 'tenantadmin') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch3_id, 'tenantadmin', 'tenantadmin@fashion.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Tenant Admin Fashion', 'tenantadmin', true, NOW(), NOW());
    END IF;
    
    -- Insert Default Branch Admins for each branch (Username: branchadmin, Password: 123456)
    IF NOT EXISTS (SELECT 1 FROM users WHERE branch_id = branch3_id AND username = 'branchadmin') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch3_id, 'branchadmin', 'branchadmin.plaza@fashion.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Branch Admin Plaza', 'branchadmin', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE branch_id = branch4_id AND username = 'branchadmin') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch4_id, 'branchadmin', 'branchadmin.grand@fashion.com', '$2a$10$inqmfpKlWFe/eg2dUwUR1ubLnKtb5oKnNX01JbPhBiAalhh.63Ocq', 'Branch Admin Grand', 'branchadmin', true, NOW(), NOW());
    END IF;
    
    -- Insert Users for Fashion Store (Password: demo123)
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin_fashion') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch3_id, 'admin_fashion', 'admin@fashion.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Admin Fashion', 'tenantadmin', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'kasir_plaza') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch3_id, 'kasir_plaza', 'kasir.plaza@fashion.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Kasir Plaza', 'user', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'sales_plaza') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch3_id, 'sales_plaza', 'sales.plaza@fashion.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Sales Plaza', 'user', true, NOW(), NOW());
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'kasir_grand') THEN
        INSERT INTO users (tenant_id, branch_id, username, email, password, full_name, role, is_active, created_at, updated_at)
        VALUES (tenant2_id, branch4_id, 'kasir_grand', 'kasir.grand@fashion.com', '$2a$10$N8rT5EwF8HKWqZfH5LxPPuKqKqKqKqKqKqKqKqKqKqKqKqKqKqKqK', 'Kasir Grand', 'user', true, NOW(), NOW());
    END IF;
    
    -- Insert Products for Fashion Store (35 products in 3 categories)
    IF NOT EXISTS (SELECT 1 FROM products WHERE sku = 'FSH-MEN-001') THEN
        INSERT INTO products (tenant_id, name, description, sku, price, stock, is_active, created_at, updated_at)
        VALUES 
            -- KATEGORI 1: PAKAIAN PRIA (12 items)
            (tenant2_id, 'Kemeja Pria Formal Putih', 'Kemeja formal lengan panjang', 'FSH-MEN-001', 250000, 30, true, NOW(), NOW()),
            (tenant2_id, 'Kemeja Pria Casual Kotak', 'Kemeja casual motif kotak', 'FSH-MEN-002', 180000, 40, true, NOW(), NOW()),
            (tenant2_id, 'Celana Jeans Pria Slim Fit', 'Celana jeans model slim fit', 'FSH-MEN-003', 320000, 25, true, NOW(), NOW()),
            (tenant2_id, 'Kaos Polo Pria', 'Kaos polo casual berbagai warna', 'FSH-MEN-004', 120000, 50, true, NOW(), NOW()),
            (tenant2_id, 'Jaket Jeans Pria', 'Jaket jeans biru dongker', 'FSH-MEN-005', 450000, 15, true, NOW(), NOW()),
            (tenant2_id, 'Celana Chino Pria', 'Celana chino cotton stretch', 'FSH-MEN-006', 280000, 35, true, NOW(), NOW()),
            (tenant2_id, 'Sweater Rajut Pria', 'Sweater rajut hangat winter', 'FSH-MEN-007', 320000, 20, true, NOW(), NOW()),
            (tenant2_id, 'Kaos Oblong Premium', 'Kaos oblong cotton combed 30s', 'FSH-MEN-008', 85000, 80, true, NOW(), NOW()),
            (tenant2_id, 'Blazer Formal Pria', 'Blazer formal untuk kantor', 'FSH-MEN-009', 650000, 12, true, NOW(), NOW()),
            (tenant2_id, 'Celana Jogger Pria', 'Celana jogger casual sporty', 'FSH-MEN-010', 195000, 45, true, NOW(), NOW()),
            (tenant2_id, 'Hoodie Pria Premium', 'Hoodie fleece dengan zipper', 'FSH-MEN-011', 380000, 25, true, NOW(), NOW()),
            (tenant2_id, 'Kemeja Batik Pria', 'Kemeja batik motif modern', 'FSH-MEN-012', 275000, 30, true, NOW(), NOW()),
            
            -- KATEGORI 2: PAKAIAN WANITA (13 items)
            (tenant2_id, 'Blouse Wanita Floral', 'Blouse motif bunga lengan panjang', 'FSH-WMN-001', 195000, 35, true, NOW(), NOW()),
            (tenant2_id, 'Dress Casual Wanita', 'Dress casual untuk santai', 'FSH-WMN-002', 280000, 20, true, NOW(), NOW()),
            (tenant2_id, 'Rok Span Wanita', 'Rok span A-line berbagai warna', 'FSH-WMN-003', 150000, 40, true, NOW(), NOW()),
            (tenant2_id, 'Celana Kulot Wanita', 'Celana kulot bahan premium', 'FSH-WMN-004', 225000, 30, true, NOW(), NOW()),
            (tenant2_id, 'Cardigan Rajut Wanita', 'Cardigan rajut hangat', 'FSH-WMN-005', 175000, 25, true, NOW(), NOW()),
            (tenant2_id, 'Tunik Wanita Muslim', 'Tunik panjang untuk muslimah', 'FSH-WMN-006', 210000, 35, true, NOW(), NOW()),
            (tenant2_id, 'Jumpsuit Wanita', 'Jumpsuit casual elegan', 'FSH-WMN-007', 320000, 18, true, NOW(), NOW()),
            (tenant2_id, 'Blazer Wanita Formal', 'Blazer wanita untuk kerja', 'FSH-WMN-008', 485000, 15, true, NOW(), NOW()),
            (tenant2_id, 'Celana Jeans Wanita', 'Celana jeans highwaist', 'FSH-WMN-009', 295000, 28, true, NOW(), NOW()),
            (tenant2_id, 'Gamis Syari Modern', 'Gamis syari dengan busui friendly', 'FSH-WMN-010', 350000, 22, true, NOW(), NOW()),
            (tenant2_id, 'Kaos Wanita Basic', 'Kaos wanita cotton premium', 'FSH-WMN-011', 75000, 70, true, NOW(), NOW()),
            (tenant2_id, 'Rok Plisket Panjang', 'Rok plisket panjang casual', 'FSH-WMN-012', 165000, 38, true, NOW(), NOW()),
            (tenant2_id, 'Outer Kimono Wanita', 'Outer kimono cardigan motif', 'FSH-WMN-013', 185000, 32, true, NOW(), NOW()),
            
            -- KATEGORI 3: AKSESORIS & FASHION ITEMS (10 items)
            (tenant2_id, 'Tas Selempang Kulit', 'Tas selempang kulit sintetis', 'FSH-ACC-001', 350000, 20, true, NOW(), NOW()),
            (tenant2_id, 'Ikat Pinggang Kulit', 'Ikat pinggang kulit asli', 'FSH-ACC-002', 125000, 45, true, NOW(), NOW()),
            (tenant2_id, 'Topi Baseball', 'Topi baseball berbagai warna', 'FSH-ACC-003', 85000, 60, true, NOW(), NOW()),
            (tenant2_id, 'Syal Rajut', 'Syal rajut untuk musim dingin', 'FSH-ACC-004', 95000, 35, true, NOW(), NOW()),
            (tenant2_id, 'Dompet Panjang Wanita', 'Dompet panjang kulit sintetis', 'FSH-ACC-005', 165000, 40, true, NOW(), NOW()),
            (tenant2_id, 'Kaos Kaki Premium', 'Kaos kaki cotton import 3 pasang', 'FSH-ACC-006', 45000, 100, true, NOW(), NOW()),
            (tenant2_id, 'Kacamata Fashion', 'Kacamata fashion anti UV', 'FSH-ACC-007', 125000, 50, true, NOW(), NOW()),
            (tenant2_id, 'Jam Tangan Pria', 'Jam tangan analog leather strap', 'FSH-ACC-008', 385000, 25, true, NOW(), NOW()),
            (tenant2_id, 'Jam Tangan Wanita', 'Jam tangan wanita model elegan', 'FSH-ACC-009', 425000, 20, true, NOW(), NOW()),
            (tenant2_id, 'Sepatu Sneakers Unisex', 'Sepatu sneakers canvas casual', 'FSH-ACC-010', 295000, 35, true, NOW(), NOW());
    END IF;
END $$;

-- Display summary
SELECT 'Demo tenants created successfully!' AS status;

-- Show tenant summary
SELECT 
    t.id,
    t.name AS tenant_name,
    t.code AS tenant_code,
    COUNT(DISTINCT b.id) AS total_branches,
    COUNT(DISTINCT u.id) AS total_users,
    COUNT(DISTINCT p.id) AS total_products
FROM tenants t
LEFT JOIN branches b ON t.id = b.tenant_id
LEFT JOIN users u ON t.id = u.tenant_id
LEFT JOIN products p ON t.id = p.tenant_id
WHERE t.code IN ('resto01', 'fashion01')
GROUP BY t.id, t.name, t.code
ORDER BY t.id;
