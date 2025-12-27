# Product API Guide

API endpoints untuk manajemen produk di MyPOSCore.

## 📋 Endpoints

### 1. Create Product
**POST** `/api/v1/products`

Membuat produk baru untuk tenant.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Nasi Goreng Spesial",
  "description": "Nasi goreng dengan telur, ayam, dan sayuran",
  "category": "Food",
  "sku": "NGS-001",
  "price": 25000,
  "stock": 100,
  "is_active": true
}
```

**Success Response (200):**
```json
{
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran",
    "category": "Food",
    "sku": "NGS-001",
    "price": 25000,
    "stock": 100,
    "is_active": true,
    "created_at": "2025-12-27 10:00:00",
    "updated_at": "2025-12-27 10:00:00"
  }
}
```

---

### 2. Get All Products
**GET** `/api/v1/products`

Mendapatkan daftar semua produk untuk tenant dengan filter opsional.

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- `category` (string, optional) - Filter berdasarkan kategori
- `search` (string, optional) - Pencarian berdasarkan nama, deskripsi, atau SKU

**Examples:**
- `/api/v1/products` - Get all products
- `/api/v1/products?category=Food` - Filter by category
- `/api/v1/products?search=goreng` - Search products

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "tenant_id": 1,
      "name": "Nasi Goreng Spesial",
      "description": "Nasi goreng dengan telur, ayam, dan sayuran",
      "category": "Food",
      "sku": "NGS-001",
      "price": 25000,
      "stock": 100,
      "is_active": true,
      "created_at": "2025-12-27 10:00:00",
      "updated_at": "2025-12-27 10:00:00"
    },
    {
      "id": 2,
      "tenant_id": 1,
      "name": "Es Teh Manis",
      "description": "Teh manis dingin",
      "category": "Beverage",
      "sku": "ETM-001",
      "price": 5000,
      "stock": 200,
      "is_active": true,
      "created_at": "2025-12-27 10:01:00",
      "updated_at": "2025-12-27 10:01:00"
    }
  ]
}
```

---

### 3. Get Product by ID
**GET** `/api/v1/products/{id}`

Mendapatkan detail produk berdasarkan ID.

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran",
    "category": "Food",
    "sku": "NGS-001",
    "price": 25000,
    "stock": 100,
    "is_active": true,
    "created_at": "2025-12-27 10:00:00",
    "updated_at": "2025-12-27 10:00:00"
  }
}
```

---

### 4. Update Product
**PUT** `/api/v1/products/{id}`

Update produk yang sudah ada.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Nasi Goreng Spesial Premium",
  "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran",
  "category": "Food",
  "sku": "NGS-001",
  "price": 30000,
  "stock": 50,
  "is_active": true
}
```

**Notes:**
- Semua field bersifat opsional
- Hanya field yang dikirim yang akan diupdate
- Field yang tidak dikirim akan tetap dengan nilai lama

**Success Response (200):**
```json
{
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Nasi Goreng Spesial Premium",
    "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran",
    "category": "Food",
    "sku": "NGS-001",
    "price": 30000,
    "stock": 50,
    "is_active": true,
    "created_at": "2025-12-27 10:00:00",
    "updated_at": "2025-12-27 10:30:00"
  }
}
```

---

### 5. Delete Product
**DELETE** `/api/v1/products/{id}`

Hapus produk (soft delete).

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "message": "Product deleted successfully"
}
```

---

### 6. Get Product Categories
**GET** `/api/v1/products/categories`

Mendapatkan daftar kategori unik dari produk tenant.

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "data": [
    "Food",
    "Beverage",
    "Snacks",
    "Desserts"
  ]
}
```

---

## 🔒 Authentication

Semua endpoint memerlukan Bearer token yang didapat dari login:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📊 Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request berhasil |
| 400 | Bad Request - Validasi error atau data tidak valid |
| 401 | Unauthorized - Token invalid atau expired |
| 404 | Not Found - Produk tidak ditemukan |
| 500 | Internal Server Error |

## 📝 Field Validations

### Create Product
- `name` (required): Nama produk, min 2 karakter
- `description` (optional): Deskripsi produk
- `category` (optional): Kategori produk
- `sku` (optional): Stock Keeping Unit
- `price` (required): Harga produk, harus > 0
- `stock` (optional): Jumlah stok, default 0
- `is_active` (optional): Status aktif, default true

### Update Product
- Semua field bersifat opsional
- `price` jika dikirim harus > 0
- `stock` jika dikirim harus >= 0

## 💡 Notes

- Produk otomatis difilter berdasarkan `tenant_id` dari JWT token
- Produk yang dihapus menggunakan soft delete (tidak dihapus dari database)
- SKU tidak harus unik, tergantung kebutuhan bisnis
- Field `search` melakukan pencarian case-insensitive di name, description, dan SKU
- Kategori bersifat free text, gunakan API `/api/v1/categories` untuk manajemen kategori terstruktur

## 🔗 Related APIs

- **Category API**: `/api/v1/categories` - Untuk manajemen kategori terstruktur
- **Order API**: `/api/v1/orders` - Order menggunakan product_id
- **Stock Management**: Update stock melalui PUT `/api/v1/products/{id}`

## 📌 Example Workflows

### 1. Create Product with Category
```bash
# Step 1: Create category (optional)
POST /api/v1/categories
{
  "name": "Food",
  "description": "All food items"
}

# Step 2: Create product
POST /api/v1/products
{
  "name": "Nasi Goreng",
  "category": "Food",
  "price": 25000,
  "stock": 100
}
```

### 2. Search and Filter
```bash
# Search by name
GET /api/v1/products?search=nasi

# Filter by category
GET /api/v1/products?category=Food

# Combine search and category
GET /api/v1/products?category=Food&search=goreng
```

### 3. Update Stock
```bash
PUT /api/v1/products/1
{
  "stock": 75
}
```

### 4. Deactivate Product
```bash
PUT /api/v1/products/1
{
  "is_active": false
}
```
