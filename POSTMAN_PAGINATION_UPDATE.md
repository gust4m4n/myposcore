# Postman Collection Update - Pagination Support

## Summary
Updated 5 list endpoints in MyPOSCore.postman_collection.json to support pagination:

## ✅ Updated Endpoints

### 1. List Categories
**Endpoint:** `GET /api/v1/categories`

**Updated URL with pagination:**
```
{{base_url}}/api/v1/categories?page=1&page_size=10
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 10, max: 100)
- `active_only` (boolean, optional): Filter to show only active categories

**Success Response (200 OK):**
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
    }
  ]
}
```

---

### 2. List Products
**Endpoint:** `GET /api/v1/products`

**Updated URL with pagination:**
```
{{base_url}}/api/v1/products?page=1&page_size=10
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 10, max: 100)
- `category` (string, optional): Filter by exact category name
- `search` (string, optional): Search keyword in name, description, or SKU

**Success Response (200 OK):**
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
      "name": "Nasi Goreng Spesial",
      "description": "Nasi goreng dengan telur, ayam, dan sayuran",
      "category": "Makanan Utama",
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
    }
  ]
}
```

---

### 3. List Users
**Endpoint:** `GET /api/v1/users`

**Updated URL with pagination:**
```
{{base_url}}/api/v1/users?page=1&page_size=10
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 10, max: 100)

**Success Response (200 OK):**
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
      "created_at": "2025-12-27 10:00:00"
    }
  ]
}
```

---

### 4. List Tenants (Superadmin)
**Endpoint:** `GET /api/v1/superadmin/tenants`

**Updated URL with pagination:**
```
{{base_url}}/api/v1/superadmin/tenants?page=1&page_size=10
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 10, max: 100)

**Success Response (200 OK):**
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
      "address": "Jl. M.H. Thamrin No. 99, Menteng, Jakarta Pusat",
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
    }
  ]
}
```

---

### 5. List Branches by Tenant (Superadmin)
**Endpoint:** `GET /api/v1/superadmin/tenants/:tenant_id/branches`

**Example URL with pagination:**
```
{{base_url}}/api/v1/superadmin/tenants/17/branches?page=1&page_size=10
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 10, max: 100)

**Success Response (200 OK):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 2,
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
    }
  ]
}
```

---

## Migration Notes

### Response Format Changes

**Before (Old Format):**
```json
{
  "data": [...]
}
```

**After (New Format with Pagination):**
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 25,
  "total_pages": 3,
  "data": [...]
}
```

### Backward Compatibility

All pagination parameters are **optional**:
- If no pagination params provided, defaults to `page=1` and `page_size=10`
- Existing clients will continue to work, but they need to handle the new response structure
- Frontend teams should update their code to read from the new pagination response structure

### Testing Recommendations

1. **Test default pagination:**
   ```
   GET /api/v1/products
   ```
   Should return page 1 with 10 items

2. **Test custom page size:**
   ```
   GET /api/v1/products?page=2&page_size=20
   ```
   Should return page 2 with 20 items (if available)

3. **Test max page size limit:**
   ```
   GET /api/v1/products?page_size=200
   ```
   Should automatically limit to 100 items

4. **Test pagination with filters:**
   ```
   GET /api/v1/products?category=Minuman&search=kopi&page=1&page_size=5
   ```
   Should apply filters AND pagination

5. **Test invalid page numbers:**
   ```
   GET /api/v1/products?page=0
   GET /api/v1/products?page=-1
   ```
   Should automatically default to page 1

---

## Import Instructions

The Postman collection file `MyPOSCore.postman_collection.json` has been updated with:
1. Pagination query parameters added to all 5 list endpoints
2. Updated example responses showing pagination metadata
3. Updated descriptions explaining pagination behavior

To update your Postman:
1. Open Postman
2. Go to Collections
3. Right-click on "MyPOSCore API" collection
4. Select "Replace"
5. Choose the updated `MyPOSCore.postman_collection.json` file
6. All requests will be updated with pagination support

---

## Additional Documentation

For detailed pagination usage examples and frontend integration tips, refer to:
- **[PAGINATION_GUIDE.md](./PAGINATION_GUIDE.md)** - Complete pagination documentation with examples

---

## Breaking Changes

⚠️ **Important:** This is a **breaking change** for response format:
- All list endpoints now return `PaginationResponse` instead of simple `{data: [...]}`
- Frontend applications must be updated to handle the new structure
- The `data` array is now nested inside pagination metadata

### Migration Code Example

**Before:**
```javascript
const response = await fetch('/api/v1/products');
const products = response.data; // Direct access
```

**After:**
```javascript
const response = await fetch('/api/v1/products?page=1&page_size=10');
const products = response.data; // Still works
const totalItems = response.total_items; // Now available
const totalPages = response.total_pages; // Now available
```

---

## Questions?

For questions or issues related to pagination implementation, please refer to:
1. [PAGINATION_GUIDE.md](./PAGINATION_GUIDE.md) - Detailed usage guide
2. Backend API source code in `handlers/` and `services/` directories
3. DTO definitions in `dto/pagination.go`
