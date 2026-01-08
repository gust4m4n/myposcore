# Category & Product Relationship Implementation

## Overview
Implementasi relasi foreign key antara table `products` dan `categories` dengan menambahkan kolom `category_id` di table products.

## Database Changes

### 1. Migration File
**File:** `migration_add_category_id_to_products.sql`

Migration ini melakukan:
- Menambahkan kolom `category_id` (INT UNSIGNED NULL) ke table products
- Membuat index untuk `category_id`
- Menambahkan foreign key constraint ke table categories
- Migrasi data dari kolom `category` (string) ke `category_id` (relasi)

**Cara menjalankan:**
```bash
mysql -u root -p myposcore < migration_add_category_id_to_products.sql
```

**Catatan:** 
- Kolom `category` (string) tetap ada sebagai legacy field untuk backward compatibility
- Setelah yakin semua data termigrasi dengan baik, Anda bisa uncomment baris terakhir di migration untuk drop kolom `category`

### 2. Demo Data
**File:** `init_products_with_categories.sql`

File ini berisi data demo products yang sudah menggunakan `category_id`:

**Tenant 1 (FoodCorner):**
- 5 Main Course items
- 5 Beverage items
- 4 Snacks items
- 3 Desserts items

**Tenant 2 (RetailStore):**
- 3 Electronics items
- 3 Furniture items
- 3 Stationery items

**Tenant 3 (Demo Café):**
- 6 Coffee items
- 4 Tea items
- 5 Pastries items
- 4 Breakfast items

**Cara menjalankan:**
```bash
# Pastikan categories sudah ada (jalankan init_category_dummy.sql dulu jika belum)
mysql -u root -p myposcore < init_category_dummy.sql

# Kemudian insert products
mysql -u root -p myposcore < init_products_with_categories.sql
```

## Code Changes

### 1. Model Updates
**File:** `models/product.go`
- Menambahkan field `CategoryID *uint`
- Menambahkan relasi `CategoryDetail *Category`

### 2. DTO Updates
**File:** `dto/product.go`
- `CreateProductRequest`: ditambah `category_id`
- `UpdateProductRequest`: ditambah `category_id`
- `ProductResponse`: ditambah `category_id` dan `category_detail`
- Menambahkan struct `CategorySummary` untuk menampilkan detail category

### 3. Service Updates
**File:** `services/product_service.go`
- Preload `CategoryDetail` pada `ListProducts()` dan `GetProduct()`
- Support `category_id` pada `CreateProduct()` dan `UpdateProduct()`
- Audit trail mencatat perubahan `category_id`

### 4. Handler Updates
**File:** `handlers/product_handler.go`
- Menambahkan helper function `mapCategoryToDTO()`
- Semua response ProductResponse menyertakan `category_detail`

## API Response Example

### List Products with Category Detail
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "tenant_id": 3,
      "name": "Espresso",
      "description": "Strong Italian coffee",
      "category": "",
      "category_id": 9,
      "category_detail": {
        "id": 9,
        "name": "Coffee",
        "description": "All coffee based beverages",
        "image": "http://localhost:8080/uploads/categories/coffee.jpg"
      },
      "sku": "CF-ESP-001",
      "price": 18000,
      "stock": 100,
      "is_active": true,
      "created_at": "2026-01-08 10:00:00",
      "updated_at": "2026-01-08 10:00:00"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total": 1,
  "total_pages": 1
}
```

## Usage Examples

### 1. Create Product with Category ID
```bash
curl -X POST "http://localhost:8080/api/v1/products" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Iced Latte",
    "description": "Cold coffee with milk",
    "category_id": 9,
    "sku": "CF-ILA-020",
    "price": 28000,
    "stock": 50,
    "is_active": true
  }'
```

### 2. Update Product Category
```bash
curl -X PUT "http://localhost:8080/api/v1/products/1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 10,
    "price": 25000
  }'
```

### 3. List Products (will include category_detail)
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Get Products by Category
```bash
# Menggunakan category string (legacy)
curl -X GET "http://localhost:8080/api/v1/products?category=Coffee&page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 5. List Categories
```bash
curl -X GET "http://localhost:8080/api/v1/categories" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Migration Steps (Complete)

### Step 1: Backup Database
```bash
mysqldump -u root -p myposcore > myposcore_backup_$(date +%Y%m%d).sql
```

### Step 2: Run Migration
```bash
mysql -u root -p myposcore < migration_add_category_id_to_products.sql
```

### Step 3: Verify Migration
```sql
-- Check if category_id column exists
DESCRIBE products;

-- Check if foreign key was created
SHOW CREATE TABLE products;

-- Verify data migration
SELECT id, name, category, category_id FROM products LIMIT 10;
```

### Step 4: Insert Demo Categories (if not exists)
```bash
mysql -u root -p myposcore < init_category_dummy.sql
```

### Step 5: Insert Demo Products
```bash
mysql -u root -p myposcore < init_products_with_categories.sql
```

### Step 6: Test API
```bash
# Start the server
go run main.go

# Test list products
curl -X GET "http://localhost:8080/api/v1/products" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Backward Compatibility

- Kolom `category` (string) tetap ada untuk backward compatibility
- API masih support filter by category string: `GET /api/v1/products?category=Coffee`
- Client bisa menggunakan `category_id` (recommended) atau `category` (legacy)
- Response API mencakup both `category` dan `category_detail`

## Benefits

1. **Data Integrity**: Foreign key constraint memastikan category_id valid
2. **Better Performance**: Index pada category_id mempercepat query
3. **Rich Response**: API response menyertakan full category detail
4. **Cascade Update**: Perubahan category name otomatis terlihat di product
5. **Easy Reporting**: Join query lebih mudah untuk reporting

## Next Steps (Optional)

1. Drop kolom `category` (string) setelah yakin migrasi sukses
2. Update query filter by category untuk menggunakan category_id
3. Add validation di handler untuk memastikan category_id exists
4. Consider adding category filter by category_id in list products API

## Database Schema

### products table (after migration)
```sql
CREATE TABLE products (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  tenant_id INT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100),  -- Legacy field
  category_id INT UNSIGNED,  -- New FK to categories
  sku VARCHAR(100),
  price DECIMAL(10,2) NOT NULL,
  stock INT DEFAULT 0,
  image VARCHAR(500),
  is_active BOOLEAN DEFAULT true,
  created_by INT UNSIGNED,
  updated_by INT UNSIGNED,
  deleted_by INT UNSIGNED,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  INDEX idx_category_id (category_id),
  FOREIGN KEY (category_id) REFERENCES categories(id) 
    ON DELETE SET NULL ON UPDATE CASCADE
);
```
