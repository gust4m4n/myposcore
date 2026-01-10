# Postman Collection Update - Category Integration

**Tanggal**: 8 Januari 2026  
**Versi**: 2.0  
**Status**: ✅ UPDATED

## Summary

Postman collection telah di-update untuk mendukung fitur category integration dengan foreign key relationship antara products dan categories table.

## Changes Made

### 1. ✨ New Endpoint Added

**Endpoint**: `GET /api/products/by-category/:category_id`

```
Name: List Products by Category
Method: GET
URL: {{base_url}}/api/products/by-category/13
Query Params: 
  - page (optional)
  - page_size (optional)

Description: Get products filtered by category ID with pagination support
```

**Example Response**:
```json
{
  "code": 0,
  "message": "Operation successful",
  "data": [
    {
      "id": 73,
      "tenant_id": 17,
      "name": "Nasi Goreng Spesial",
      "description": "Nasi goreng dengan telur, ayam, dan sayuran",
      "category": "Makanan Utama",
      "category_id": 13,
      "category_detail": {
        "id": 13,
        "name": "Makanan Utama",
        "description": "Menu makanan utama dan hidangan pokok",
        "image": null
      },
      "sku": "FOOD-001",
      "price": 25000,
      "stock": 50,
      "image": "/uploads/products/product_73.jpg",
      "is_active": true,
      "created_at": "2026-01-08T10:00:00+07:00",
      "updated_at": "2026-01-08T10:00:00+07:00",
      "created_by": 1,
      "created_by_name": "Admin Resto",
      "updated_by": null,
      "updated_by_name": null
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 5,
    "total_pages": 1
  }
}
```

### 2. 📝 Updated Request Bodies

**Create Product** - Added `category_id` field:
```json
{
  "name": "Produk Baru",
  "description": "Deskripsi produk baru",
  "category": "Makanan Utama",
  "category_id": 13,           // ✨ NEW
  "sku": "SKU-NEW-001",
  "price": 50000,
  "stock": 100,
  "is_active": true
}
```

**Update Product** - Added `category_id` field:
```json
{
  "name": "Produk Updated",
  "description": "Deskripsi produk yang sudah diupdate",
  "category": "Minuman",
  "category_id": 13,           // ✨ NEW
  "sku": "SKU-UPD-001",
  "price": 75000,
  "stock": 150,
  "is_active": true
}
```

### 3. 🔄 Updated Response Examples

All product response examples now include:
- `category_id` - Integer foreign key to categories table
- `category_detail` - Nested category object with full details

**Affected Endpoints**:
1. ✅ **List Products** - All products in array include category_detail
2. ✅ **Get Product by ID** - Single product includes category_detail
3. ✅ **Create Product** - Created product returns category_detail
4. ✅ **Update Product** - Updated product returns category_detail
5. ✅ **Upload Product Image** - Product with uploaded image includes category_detail

**Example category_detail structure**:
```json
{
  "category_id": 13,
  "category_detail": {
    "id": 13,
    "name": "Makanan Utama",
    "description": "Menu makanan utama dan hidangan pokok",
    "image": null
  }
}
```

## Migration Notes

### Database Schema
Products table now has:
- `category_id` column (INTEGER, nullable)
- Foreign key constraint to `categories(id)`
- Index on `category_id` for performance

### Backward Compatibility
- Legacy `category` string field is retained for backward compatibility
- Both `category` and `category_id` can be used
- API responses include both fields
- Recommended to use `category_id` for new implementations

### Demo Data
Collection includes demo data with:
- Category ID 13: "Makanan Utama"
- Category ID 14: "Minuman"
- Category ID 15: "Snack & Dessert"
- Products linked to tenant_id 17

## Testing Guide

### 1. Test New Endpoint
```bash
GET {{base_url}}/api/products/by-category/13
Authorization: Bearer {{auth_token}}
```

Expected: List of products with category_id = 13

### 2. Test Create with Category ID
```bash
POST {{base_url}}/api/products
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "name": "Test Product",
  "category_id": 13,
  "price": 10000,
  "stock": 50
}
```

Expected: Product created with category_detail in response

### 3. Test List Products
```bash
GET {{base_url}}/api/products
Authorization: Bearer {{auth_token}}
```

Expected: All products include category_detail (if category_id is set)

### 4. Test Update with Category ID
```bash
PUT {{base_url}}/api/products/73
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "category_id": 14,
  "price": 15000
}
```

Expected: Product updated with new category_detail

## Updated Endpoints Summary

| Endpoint | Method | Changes |
|----------|--------|---------|
| List Products | GET | ✅ Response includes `category_detail` |
| Get Product by ID | GET | ✅ Response includes `category_detail` |
| Create Product | POST | ✅ Request accepts `category_id`, response includes `category_detail` |
| Update Product | PUT | ✅ Request accepts `category_id`, response includes `category_detail` |
| Upload Product Image | POST | ✅ Response includes `category_detail` |
| **List by Category** | GET | ✨ **NEW ENDPOINT** |

## Files Modified

- `MyPOSCore.postman_collection.json` - Main Postman collection (UPDATED)
- `update_postman.py` - Python script for automated updates (NEW)

## Backup

Backup files created with timestamp:
```
MyPOSCore.postman_collection.json.backup_20260108_*
```

## Related Documentation

- [PRODUCT_CATEGORY_IMPLEMENTATION.md](PRODUCT_CATEGORY_IMPLEMENTATION.md) - Backend implementation details
- [API_ROUTES_VERIFICATION_COMPLETE.md](API_ROUTES_VERIFICATION_COMPLETE.md) - Route verification results
- [migration_add_category_id_to_products.sql](migration_add_category_id_to_products.sql) - Database migration
- [init_products_with_categories.sql](init_products_with_categories.sql) - Demo data

## Import to Postman

1. Open Postman
2. Click Import button
3. Select `MyPOSCore.postman_collection.json`
4. Import `MyPOSCore.postman_environment.json` for environment variables
5. Set `auth_token` variable after login

## Environment Variables

Required variables in Postman environment:
```
base_url = http://localhost:8080
auth_token = <JWT token from login response>
```

## Next Steps

1. ✅ Import updated collection to Postman
2. ✅ Test all product endpoints with category_detail
3. ✅ Verify new "List by Category" endpoint
4. ✅ Update mobile/frontend apps to use category_id
5. ✅ Migrate existing products to use category_id

---

**Total Product Endpoints**: 9  
**New Endpoints**: 1  
**Updated Endpoints**: 8  
**Response Examples Updated**: 7  
**Request Bodies Updated**: 2
