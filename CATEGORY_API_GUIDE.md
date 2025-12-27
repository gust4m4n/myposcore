# Category API Guide

API endpoints untuk manajemen kategori produk di MyPOSCore.

## 📋 Endpoints

### 1. Create Category
**POST** `/api/v1/categories`

Membuat kategori produk baru untuk tenant.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Beverages",
  "description": "All kinds of drinks"
}
```

**Success Response (200):**
```json
{
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Beverages",
    "description": "All kinds of drinks",
    "is_active": true,
    "created_at": "2025-12-27 10:00:00",
    "updated_at": "2025-12-27 10:00:00"
  }
}
```

---

### 2. Get All Categories
**GET** `/api/v1/categories`

Mendapatkan daftar semua kategori untuk tenant.

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- `active_only` (boolean, optional) - Filter hanya kategori aktif

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "tenant_id": 1,
      "name": "Beverages",
      "description": "All kinds of drinks",
      "is_active": true,
      "created_at": "2025-12-27 10:00:00",
      "updated_at": "2025-12-27 10:00:00"
    },
    {
      "id": 2,
      "tenant_id": 1,
      "name": "Food",
      "description": "All food items",
      "is_active": true,
      "created_at": "2025-12-27 10:01:00",
      "updated_at": "2025-12-27 10:01:00"
    }
  ]
}
```

---

### 3. Get Category by ID
**GET** `/api/v1/categories/{id}`

Mendapatkan detail kategori berdasarkan ID.

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
    "name": "Beverages",
    "description": "All kinds of drinks",
    "is_active": true,
    "created_at": "2025-12-27 10:00:00",
    "updated_at": "2025-12-27 10:00:00"
  }
}
```

---

### 4. Update Category
**PUT** `/api/v1/categories/{id}`

Update kategori yang sudah ada.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Hot Beverages",
  "description": "All hot drinks including coffee and tea",
  "is_active": true
}
```

**Success Response (200):**
```json
{
  "message": "Category updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Hot Beverages",
    "description": "All hot drinks including coffee and tea",
    "is_active": true,
    "created_at": "2025-12-27 10:00:00",
    "updated_at": "2025-12-27 10:30:00"
  }
}
```

---

### 5. Delete Category
**DELETE** `/api/v1/categories/{id}`

Hapus kategori (hanya jika tidak digunakan oleh produk).

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "message": "Category deleted successfully"
}
```

**Error Response (400):**
```json
{
  "error": "cannot delete category that is used by products"
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
| 404 | Not Found - Kategori tidak ditemukan |
| 500 | Internal Server Error |

## 💡 Notes

- Nama kategori harus unik per tenant
- Kategori yang digunakan oleh produk tidak bisa dihapus
- Field `name` wajib diisi minimal 2 karakter, maksimal 100 karakter
- Kategori otomatis difilter berdasarkan `tenant_id` dari JWT token
