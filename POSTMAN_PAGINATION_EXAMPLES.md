# Postman Pagination Examples - MyPOSCore API

Berikut adalah contoh lengkap request dan response untuk semua API yang support pagination.

---

## 1. List Users API

### Endpoint
```
GET /api/users
```

### Example 1: Default Pagination (Page 1, 10 items)

**Request:**
```
GET {{base_url}}/api/users?page=1&page_size=10
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 25,
  "total_pages": 3,
  "data": [
    {
      "id": 1,
      "tenant_id": 17,
      "branch_id": 1,
      "email": "user1@example.com",
      "full_name": "User Satu",
      "image": "http://localhost:8080/uploads/profiles/user1.jpg",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-27 10:00:00",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 2,
      "tenant_id": 17,
      "branch_id": 1,
      "email": "admin1@example.com",
      "full_name": "Admin Satu",
      "image": "",
      "role": "branchadmin",
      "is_active": true,
      "created_at": "2025-12-27 10:05:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 3,
      "tenant_id": 17,
      "branch_id": 1,
      "email": "cashier1@example.com",
      "full_name": "Kasir Satu",
      "image": "",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-27 10:10:00",
      "created_by": 2,
      "created_by_name": "Admin Satu",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 2: Page 2 with Custom Size (5 items)

**Request:**
```
GET {{base_url}}/api/users?page=2&page_size=5
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 2,
  "page_size": 5,
  "total_items": 25,
  "total_pages": 5,
  "data": [
    {
      "id": 6,
      "tenant_id": 17,
      "branch_id": 2,
      "email": "user6@example.com",
      "full_name": "User Enam",
      "image": "",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-27 11:00:00",
      "created_by": 2,
      "created_by_name": "Admin Satu",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 7,
      "tenant_id": 17,
      "branch_id": 2,
      "email": "user7@example.com",
      "full_name": "User Tujuh",
      "image": "",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-27 11:05:00",
      "created_by": 2,
      "created_by_name": "Admin Satu",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 3: No Params (Uses Default: page=1, page_size=10)

**Request:**
```
GET {{base_url}}/api/users
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 25,
  "total_pages": 3,
  "data": [
    {
      "id": 1,
      "tenant_id": 17,
      "branch_id": 1,
      "email": "user1@example.com",
      "full_name": "User Satu",
      "image": "",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-27 10:00:00"
    }
  ]
}
```

### Example 4: Large Page Size (50 items)

**Request:**
```
GET {{base_url}}/api/users?page=1&page_size=50
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 50,
  "total_items": 25,
  "total_pages": 1,
  "data": [
    // All 25 users in one page
  ]
}
```

---

## 2. List Categories API

### Endpoint
```
GET /api/categories
```

### Example 1: Page 1 with 10 items

**Request:**
```
GET {{base_url}}/api/categories?page=1&page_size=10
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 8,
  "total_pages": 1,
  "data": [
    {
      "id": 1,
      "tenant_id": 17,
      "name": "Beverages",
      "description": "Coffee, tea, and other drinks",
      "image": "http://localhost:8080/uploads/categories/1_20241230_142208.jpg",
      "is_active": true,
      "created_at": "2025-12-27 10:00:00",
      "updated_at": "2025-12-27 10:00:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 2,
      "tenant_id": 17,
      "name": "Electronics",
      "description": "Electronic devices and accessories",
      "image": "",
      "is_active": true,
      "created_at": "2025-12-27 10:01:00",
      "updated_at": "2025-12-27 10:01:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 3,
      "tenant_id": 17,
      "name": "Food",
      "description": "Food and snacks",
      "image": "",
      "is_active": true,
      "created_at": "2025-12-27 10:02:00",
      "updated_at": "2025-12-27 10:02:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 2: Page 1 with Active Only Filter

**Request:**
```
GET {{base_url}}/api/categories?page=1&page_size=10&active_only=true
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 6,
  "total_pages": 1,
  "data": [
    {
      "id": 1,
      "tenant_id": 17,
      "name": "Beverages",
      "description": "Coffee, tea, and other drinks",
      "image": "http://localhost:8080/uploads/categories/1_20241230_142208.jpg",
      "is_active": true,
      "created_at": "2025-12-27 10:00:00",
      "updated_at": "2025-12-27 10:00:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 3: Page 2 with 3 items per page

**Request:**
```
GET {{base_url}}/api/categories?page=2&page_size=3
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 2,
  "page_size": 3,
  "total_items": 8,
  "total_pages": 3,
  "data": [
    {
      "id": 4,
      "tenant_id": 17,
      "name": "Clothing",
      "description": "Apparel and fashion",
      "image": "",
      "is_active": true,
      "created_at": "2025-12-27 10:03:00",
      "updated_at": "2025-12-27 10:03:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 5,
      "tenant_id": 17,
      "name": "Books",
      "description": "Books and magazines",
      "image": "",
      "is_active": false,
      "created_at": "2025-12-27 10:04:00",
      "updated_at": "2025-12-27 10:04:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

---

## 3. List Products API

### Endpoint
```
GET /api/products
```

### Example 1: Page 1 with 10 items

**Request:**
```
GET {{base_url}}/api/products?page=1&page_size=10
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 45,
  "total_pages": 5,
  "data": [
    {
      "id": 1,
      "tenant_id": 17,
      "name": "Nasi Goreng Spesial",
      "description": "Nasi goreng dengan telur, ayam, dan sayuran",
      "category": "Food",
      "sku": "FOOD-001",
      "price": 25000,
      "stock": 50,
      "image": "http://localhost:8080/uploads/products/product_1_1735296000.jpg",
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 2,
      "tenant_id": 17,
      "name": "Es Teh Manis",
      "description": "Teh manis dingin segar",
      "category": "Beverages",
      "sku": "DRINK-001",
      "price": 5000,
      "stock": 100,
      "image": "http://localhost:8080/uploads/products/product_2_1735296000.jpg",
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 2: Page 1 with Category Filter

**Request:**
```
GET {{base_url}}/api/products?page=1&page_size=10&category=Beverages
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 12,
  "total_pages": 2,
  "data": [
    {
      "id": 2,
      "tenant_id": 17,
      "name": "Es Teh Manis",
      "description": "Teh manis dingin segar",
      "category": "Beverages",
      "sku": "DRINK-001",
      "price": 5000,
      "stock": 100,
      "image": "http://localhost:8080/uploads/products/product_2_1735296000.jpg",
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 5,
      "tenant_id": 17,
      "name": "Kopi Susu",
      "description": "Kopi dengan susu segar",
      "category": "Beverages",
      "sku": "DRINK-002",
      "price": 15000,
      "stock": 80,
      "image": "",
      "is_active": true,
      "created_at": "2025-12-25T10:05:00+07:00",
      "updated_at": "2025-12-25T10:05:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 3: Page 1 with Search Filter

**Request:**
```
GET {{base_url}}/api/products?page=1&page_size=5&search=kopi
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 5,
  "total_items": 3,
  "total_pages": 1,
  "data": [
    {
      "id": 5,
      "tenant_id": 17,
      "name": "Kopi Susu",
      "description": "Kopi dengan susu segar",
      "category": "Beverages",
      "sku": "DRINK-002",
      "price": 15000,
      "stock": 80,
      "image": "",
      "is_active": true,
      "created_at": "2025-12-25T10:05:00+07:00",
      "updated_at": "2025-12-25T10:05:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 8,
      "tenant_id": 17,
      "name": "Kopi Hitam",
      "description": "Kopi tanpa gula",
      "category": "Beverages",
      "sku": "DRINK-005",
      "price": 10000,
      "stock": 90,
      "image": "",
      "is_active": true,
      "created_at": "2025-12-25T10:08:00+07:00",
      "updated_at": "2025-12-25T10:08:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 4: Page 2 with Category and Search Filter

**Request:**
```
GET {{base_url}}/api/products?page=2&page_size=5&category=Food&search=ayam
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 2,
  "page_size": 5,
  "total_items": 8,
  "total_pages": 2,
  "data": [
    {
      "id": 12,
      "tenant_id": 17,
      "name": "Ayam Bakar",
      "description": "Ayam bakar dengan bumbu khas",
      "category": "Food",
      "sku": "FOOD-008",
      "price": 35000,
      "stock": 30,
      "image": "",
      "is_active": true,
      "created_at": "2025-12-25T10:12:00+07:00",
      "updated_at": "2025-12-25T10:12:00+07:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 5: Page 1 with Large Page Size (20 items)

**Request:**
```
GET {{base_url}}/api/products?page=1&page_size=20
Authorization: Bearer {{auth_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 20,
  "total_items": 45,
  "total_pages": 3,
  "data": [
    // 20 products...
  ]
}
```

---

## 4. List Tenants API (Superadmin)

### Endpoint
```
GET /api/superadmin/tenants
```

### Example 1: Page 1 with 10 items

**Request:**
```
GET {{base_url}}/api/superadmin/tenants?page=1&page_size=10
Authorization: Bearer {{superadmin_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 3,
  "total_pages": 1,
  "data": [
    {
      "id": 1,
      "name": "Food Corner 99",
      "description": "Restoran modern dengan berbagai pilihan menu internasional dan lokal",
      "address": "Jl. M.H. Thamrin No. 99, Menteng, Jakarta Pusat, DKI Jakarta 10350",
      "website": "https://www.foodcorner99.com",
      "email": "info@foodcorner99.com",
      "phone": "021-3199-8899",
      "image": "http://localhost:8080/uploads/tenants/tenant_foodcorner99_1704000000.png",
      "is_active": true,
      "created_at": "2025-12-24 10:00:00",
      "updated_at": "2025-12-24 10:00:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": 1,
      "updated_by_name": "Super Admin"
    },
    {
      "id": 2,
      "name": "Fashion Store 101",
      "description": "Pusat fashion terlengkap dengan koleksi pakaian, sepatu, dan aksesoris trendy",
      "address": "Grand Indonesia Shopping Town Tower C, Jl. M.H. Thamrin No. 1, Menteng, Jakarta Pusat",
      "website": "https://www.fashionstore101.com",
      "email": "contact@fashionstore101.com",
      "phone": "021-2358-1010",
      "image": "http://localhost:8080/uploads/tenants/tenant_fashionstore101_1704001000.png",
      "is_active": true,
      "created_at": "2025-12-24 10:00:00",
      "updated_at": "2025-12-25 14:30:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": 2,
      "updated_by_name": "John Admin"
    },
    {
      "id": 3,
      "name": "Tech Shop 2025",
      "description": "Electronics and gadget store",
      "address": "Plaza Senayan, Jakarta",
      "website": "https://www.techshop2025.com",
      "email": "info@techshop2025.com",
      "phone": "021-5555-1234",
      "image": "",
      "is_active": true,
      "created_at": "2025-12-24 11:00:00",
      "updated_at": "2025-12-24 11:00:00",
      "created_by": 1,
      "created_by_name": "Super Admin",
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 2: No Params (Uses Default)

**Request:**
```
GET {{base_url}}/api/superadmin/tenants
Authorization: Bearer {{superadmin_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 3,
  "total_pages": 1,
  "data": [
    // Same as Example 1
  ]
}
```

---

## 5. List Branches API (Superadmin)

### Endpoint
```
GET /api/superadmin/tenants/:tenant_id/branches
```

### Example 1: Page 1 with 10 items

**Request:**
```
GET {{base_url}}/api/superadmin/tenants/17/branches?page=1&page_size=10
Authorization: Bearer {{superadmin_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 5,
  "total_pages": 1,
  "data": [
    {
      "id": 25,
      "tenant_id": 17,
      "name": "Branch Central",
      "description": "Main branch in central business district",
      "address": "Jl. MH Thamrin No. 100",
      "website": "https://foodcorner.com/central",
      "email": "central@foodcorner.com",
      "phone": "+628123456789",
      "image": "http://localhost:8080/uploads/branches/25_20241230_142208.jpg",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 26,
      "tenant_id": 17,
      "name": "Branch North",
      "description": "Branch in north area",
      "address": "Jl. Gatot Subroto No. 200",
      "website": "https://foodcorner.com/north",
      "email": "north@foodcorner.com",
      "phone": "+628123456790",
      "image": "",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 27,
      "tenant_id": 17,
      "name": "Branch South",
      "description": "Branch in south area",
      "address": "Jl. Sudirman No. 300",
      "website": "https://foodcorner.com/south",
      "email": "south@foodcorner.com",
      "phone": "+628123456791",
      "image": "",
      "is_active": true,
      "created_at": "2024-12-30 14:25:08",
      "updated_at": "2024-12-30 14:25:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 2: Page 1 with 2 items per page

**Request:**
```
GET {{base_url}}/api/superadmin/tenants/17/branches?page=1&page_size=2
Authorization: Bearer {{superadmin_token}}
```

**Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 2,
  "total_items": 5,
  "total_pages": 3,
  "data": [
    {
      "id": 25,
      "tenant_id": 17,
      "name": "Branch Central",
      "description": "Main branch in central business district",
      "address": "Jl. MH Thamrin No. 100",
      "website": "https://foodcorner.com/central",
      "email": "central@foodcorner.com",
      "phone": "+628123456789",
      "image": "http://localhost:8080/uploads/branches/25_20241230_142208.jpg",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 26,
      "tenant_id": 17,
      "name": "Branch North",
      "description": "Branch in north area",
      "address": "Jl. Gatot Subroto No. 200",
      "website": "https://foodcorner.com/north",
      "email": "north@foodcorner.com",
      "phone": "+628123456790",
      "image": "",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

### Example 3: Page 2 with 2 items per page

**Request:**
```
GET {{base_url}}/api/superadmin/tenants/17/branches?page=2&page_size=2
Authorization: Bearer {{superadmin_token}}
```

**Response (200 OK):**
```json
{
  "page": 2,
  "page_size": 2,
  "total_items": 5,
  "total_pages": 3,
  "data": [
    {
      "id": 27,
      "tenant_id": 17,
      "name": "Branch South",
      "description": "Branch in south area",
      "address": "Jl. Sudirman No. 300",
      "website": "https://foodcorner.com/south",
      "email": "south@foodcorner.com",
      "phone": "+628123456791",
      "image": "",
      "is_active": true,
      "created_at": "2024-12-30 14:25:08",
      "updated_at": "2024-12-30 14:25:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 28,
      "tenant_id": 17,
      "name": "Branch East",
      "description": "Branch in east area",
      "address": "Jl. Ahmad Yani No. 400",
      "website": "https://foodcorner.com/east",
      "email": "east@foodcorner.com",
      "phone": "+628123456792",
      "image": "",
      "is_active": true,
      "created_at": "2024-12-30 14:28:08",
      "updated_at": "2024-12-30 14:28:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    }
  ]
}
```

---

## 📝 Important Notes

### Query Parameters Summary

| Parameter | Type | Default | Min | Max | Description |
|-----------|------|---------|-----|-----|-------------|
| `page` | integer | 1 | 1 | ∞ | Page number (starts from 1) |
| `page_size` | integer | 10 | 1 | 100 | Items per page |

### Validation Rules

1. **page < 1** → automatically set to 1
2. **page_size < 1** → automatically set to 10
3. **page_size > 100** → automatically capped at 100
4. **No params provided** → defaults to page=1, page_size=10

### Response Fields

All paginated responses include:

- **page**: Current page number
- **page_size**: Items per page
- **total_items**: Total number of items in database
- **total_pages**: Total number of pages (calculated as ceil(total_items / page_size))
- **data**: Array of items for current page

### Postman Variables

Make sure to set these variables in Postman:

- `{{base_url}}`: http://localhost:8080
- `{{auth_token}}`: JWT token for regular user
- `{{superadmin_token}}`: JWT token for superadmin user

### How to Test in Postman

1. **Import Collection**: Import `MyPOSCore.postman_collection.json`
2. **Import Environment**: Import `MyPOSCore.postman_environment.json`
3. **Login First**: Use Login endpoint to get auth token
4. **Set Token**: Token automatically saved to environment variable
5. **Test Pagination**: Try different page and page_size values

### cURL Examples

```bash
# Test List Users
curl -X GET "http://localhost:8080/api/users?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Test List Products with filters
curl -X GET "http://localhost:8080/api/products?page=1&page_size=5&category=Beverages&search=kopi" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Test List Categories with active filter
curl -X GET "http://localhost:8080/api/categories?page=1&page_size=10&active_only=true" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Test List Tenants (Superadmin)
curl -X GET "http://localhost:8080/api/superadmin/tenants?page=1&page_size=10" \
  -H "Authorization: Bearer SUPERADMIN_TOKEN"

# Test List Branches (Superadmin)
curl -X GET "http://localhost:8080/api/superadmin/tenants/17/branches?page=1&page_size=10" \
  -H "Authorization: Bearer SUPERADMIN_TOKEN"
```

---

## 🎯 Testing Scenarios

### Scenario 1: Default Behavior
Test without pagination params to verify default values work.

### Scenario 2: Custom Page Size
Test with different page_size values (5, 20, 50, 100).

### Scenario 3: Navigation
Test page navigation (page 1, 2, 3, etc.) to verify data changes.

### Scenario 4: Edge Cases
- Test page=0 (should default to 1)
- Test page_size=0 (should default to 10)
- Test page_size=200 (should cap at 100)
- Test page beyond total_pages (should return empty data array)

### Scenario 5: Combined Filters
Test pagination with filters (category, search, active_only).

### Scenario 6: Large Datasets
Create 100+ items and test pagination performance.

---

**Last Updated:** December 31, 2024  
**Server:** http://localhost:8080  
**Version:** 1.0
