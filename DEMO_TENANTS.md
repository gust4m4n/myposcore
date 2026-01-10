# Demo Tenants Documentation

File ini berisi dokumentasi untuk 2 tenant demo yang sudah dibuat di sistem MyPOSCore.

## Cara Menjalankan Demo Data

```bash
psql -U postgres -d myposcore -f init_demo_tenants.sql
```

---

## 🍽️ TENANT 1: WARUNG MAKAN SEJAHTERA (Restoran)

### Informasi Tenant
- **Tenant ID**: 9
- **Nama**: Warung Makan Sejahtera
- **Kode**: `resto01`
- **Tipe**: Restoran

### Cabang (Branches)

| ID | Nama | Kode | Alamat |
|----|------|------|--------|
| 17 | Cabang Pusat | `resto01-pusat` | Jl. Sudirman No. 123, Jakarta |
| 18 | Cabang Menteng | `resto01-menteng` | Jl. Menteng Raya No. 45, Jakarta |

### Pengguna (Users)

| Username | Password | Nama Lengkap | Role | Cabang |
|----------|----------|--------------|------|--------|
| `admin_resto` | `demo123` | Admin Resto | tenantadmin | Cabang Pusat |
| `kasir_pusat` | `demo123` | Kasir Pusat | user | Cabang Pusat |
| `kasir_menteng` | `demo123` | Kasir Menteng | user | Cabang Menteng |

### Produk (12 items total)

#### Makanan (7 items)
| SKU | Nama | Harga | Stok |
|-----|------|-------|------|
| RESTO-MKN-001 | Nasi Goreng Spesial | Rp 35,000 | 100 |
| RESTO-MKN-002 | Mie Goreng Seafood | Rp 32,000 | 100 |
| RESTO-MKN-003 | Ayam Bakar | Rp 28,000 | 50 |
| RESTO-MKN-004 | Soto Ayam | Rp 25,000 | 80 |
| RESTO-MKN-005 | Gado-gado | Rp 20,000 | 60 |
| RESTO-MKN-006 | Sate Ayam | Rp 30,000 | 70 |
| RESTO-MKN-007 | Cap Cay | Rp 28,000 | 50 |

#### Minuman (5 items)
| SKU | Nama | Harga | Stok |
|-----|------|-------|------|
| RESTO-MNM-001 | Es Teh Manis | Rp 5,000 | 200 |
| RESTO-MNM-002 | Es Jeruk | Rp 8,000 | 150 |
| RESTO-MNM-003 | Es Kelapa Muda | Rp 12,000 | 100 |
| RESTO-MNM-004 | Kopi Susu | Rp 10,000 | 120 |
| RESTO-MNM-005 | Jus Alpukat | Rp 15,000 | 80 |

### Contoh Login API

```bash
# Login sebagai admin resto
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "resto01",
    "branch_code": "resto01-pusat",
    "username": "admin_resto",
    "password": "demo123"
  }'

# Login sebagai kasir
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "resto01",
    "branch_code": "resto01-pusat",
    "username": "kasir_pusat",
    "password": "demo123"
  }'
```

---

## 👔 TENANT 2: FASHION STORE (Toko Baju)

### Informasi Tenant
- **Tenant ID**: 10
- **Nama**: Fashion Store
- **Kode**: `fashion01`
- **Tipe**: Toko Baju

### Cabang (Branches)

| ID | Nama | Kode | Alamat |
|----|------|------|--------|
| 19 | Cabang Mall Plaza | `fashion01-plaza` | Mall Plaza Lt. 2 No. 45, Jakarta |
| 20 | Cabang Grand Mall | `fashion01-grand` | Grand Mall Lt. 1 No. 78, Bandung |

### Pengguna (Users)

| Username | Password | Nama Lengkap | Role | Cabang |
|----------|----------|--------------|------|--------|
| `admin_fashion` | `demo123` | Admin Fashion | tenantadmin | Cabang Mall Plaza |
| `kasir_plaza` | `demo123` | Kasir Plaza | user | Cabang Mall Plaza |
| `sales_plaza` | `demo123` | Sales Plaza | user | Cabang Mall Plaza |
| `kasir_grand` | `demo123` | Kasir Grand | user | Cabang Grand Mall |

### Produk (14 items total)

#### Pakaian Pria (5 items)
| SKU | Nama | Harga | Stok |
|-----|------|-------|------|
| FSH-MEN-001 | Kemeja Pria Formal Putih | Rp 250,000 | 30 |
| FSH-MEN-002 | Kemeja Pria Casual Kotak | Rp 180,000 | 40 |
| FSH-MEN-003 | Celana Jeans Pria Slim Fit | Rp 320,000 | 25 |
| FSH-MEN-004 | Kaos Polo Pria | Rp 120,000 | 50 |
| FSH-MEN-005 | Jaket Jeans Pria | Rp 450,000 | 15 |

#### Pakaian Wanita (5 items)
| SKU | Nama | Harga | Stok |
|-----|------|-------|------|
| FSH-WMN-001 | Blouse Wanita Floral | Rp 195,000 | 35 |
| FSH-WMN-002 | Dress Casual Wanita | Rp 280,000 | 20 |
| FSH-WMN-003 | Rok Span Wanita | Rp 150,000 | 40 |
| FSH-WMN-004 | Celana Kulot Wanita | Rp 225,000 | 30 |
| FSH-WMN-005 | Cardigan Rajut Wanita | Rp 175,000 | 25 |

#### Aksesoris (4 items)
| SKU | Nama | Harga | Stok |
|-----|------|-------|------|
| FSH-ACC-001 | Tas Selempang Kulit | Rp 350,000 | 20 |
| FSH-ACC-002 | Ikat Pinggang Kulit | Rp 125,000 | 45 |
| FSH-ACC-003 | Topi Baseball | Rp 85,000 | 60 |
| FSH-ACC-004 | Syal Rajut | Rp 95,000 | 35 |

### Contoh Login API

```bash
# Login sebagai admin fashion
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "fashion01",
    "branch_code": "fashion01-plaza",
    "username": "admin_fashion",
    "password": "demo123"
  }'

# Login sebagai kasir
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "fashion01",
    "branch_code": "fashion01-plaza",
    "username": "kasir_plaza",
    "password": "demo123"
  }'
```

---

## 🔑 Credential Summary

Semua user menggunakan password yang sama untuk kemudahan demo:

**Password untuk semua user demo**: `demo123`

### Tenant Codes
- Restoran: `resto01`
- Fashion Store: `fashion01`

### Branch Codes
- `resto01-pusat` - Cabang Pusat Restoran
- `resto01-menteng` - Cabang Menteng Restoran
- `fashion01-plaza` - Cabang Mall Plaza Fashion
- `fashion01-grand` - Cabang Grand Mall Fashion

---

## 📊 Database Statistics

```sql
-- Lihat ringkasan semua tenant
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
```

**Hasil:**
| Tenant ID | Tenant Name | Tenant Code | Branches | Users | Products |
|-----------|-------------|-------------|----------|-------|----------|
| 9 | Warung Makan Sejahtera | resto01 | 2 | 3 | 12 |
| 10 | Fashion Store | fashion01 | 2 | 4 | 14 |

---

## 🧪 Testing dengan Postman

1. Import collection: `MyPOSCore.postman_collection.json`
2. Import environment: `MyPOSCore.postman_environment.json`
3. Tambahkan variabel untuk demo tenants di environment:

```json
{
  "resto_tenant": "resto01",
  "resto_branch": "resto01-pusat",
  "resto_username": "admin_resto",
  "resto_password": "demo123",
  
  "fashion_tenant": "fashion01",
  "fashion_branch": "fashion01-plaza",
  "fashion_username": "admin_fashion",
  "fashion_password": "demo123"
}
```

---

## 🗑️ Menghapus Demo Data

Jika ingin menghapus data demo:

```sql
-- Hapus produk
DELETE FROM products WHERE tenant_id IN (
  SELECT id FROM tenants WHERE code IN ('resto01', 'fashion01')
);

-- Hapus users (kecuali superadmin)
DELETE FROM users WHERE tenant_id IN (
  SELECT id FROM tenants WHERE code IN ('resto01', 'fashion01')
);

-- Hapus branches
DELETE FROM branches WHERE tenant_id IN (
  SELECT id FROM tenants WHERE code IN ('resto01', 'fashion01')
);

-- Hapus tenants
DELETE FROM tenants WHERE code IN ('resto01', 'fashion01');
```

---

## 📝 Notes

1. Data ini untuk **DEMO/DEVELOPMENT** saja, jangan digunakan di production
2. Password `demo123` sangat lemah dan hanya untuk testing
3. Semua data dapat di-reset dan di-recreate dengan menjalankan `init_demo_tenants.sql` lagi
4. Tenant isolation sudah diterapkan - setiap tenant hanya bisa akses data mereka sendiri
5. Role-based access control sudah aktif - tenantadmin memiliki akses penuh ke tenant mereka
