# Product Category Feature Guide

## Overview
Product category feature memungkinkan filtering dan pencarian product berdasarkan kategori dan keyword.

## API Endpoints

### 1. Get Categories
**Endpoint:** `GET /api/v1/products/categories`  
**Auth:** Required (Bearer Token)  
**Description:** Mendapatkan list kategori unik untuk tenant

**Response Example:**
```json
{
  "data": [
    "Makanan Utama",
    "Minuman",
    "Snack & Dessert"
  ]
}
```

### 2. List Products with Filters
**Endpoint:** `GET /api/v1/products`  
**Auth:** Required (Bearer Token)  
**Query Parameters:**
- `category` (optional): Filter by exact category name
- `search` (optional): Search keyword (case-insensitive)

**Examples:**
```bash
# Get all products
GET /api/v1/products

# Filter by category
GET /api/v1/products?category=Minuman

# Search by keyword
GET /api/v1/products?search=ayam

# Combined filter + search
GET /api/v1/products?category=Makanan%20Utama&search=goreng
```

**Response Example:**
```json
{
  "data": [
    {
      "id": 42,
      "tenant_id": 17,
      "name": "Es Teh Manis",
      "description": "Teh manis dingin",
      "category": "Minuman",
      "sku": "RESTO-MNM-001",
      "price": 5000,
      "stock": 200,
      "is_active": true,
      "created_at": "2025-12-24 19:44:27",
      "updated_at": "2025-12-24 19:44:27"
    }
  ]
}
```

## Demo Data Categories

### Tenant: resto01 (Restaurant)
- **Makanan Utama** - 15 products (RESTO-MKN-*)
- **Minuman** - 10 products (RESTO-MNM-*)
- **Snack & Dessert** - 10 products (RESTO-SNK-*)

### Tenant: fashion01 (Fashion Store)
- **Pakaian Pria** - 12 products (FSH-MEN-*)
- **Pakaian Wanita** - 13 products (FSH-WMN-*)
- **Aksesoris** - 10 products (FSH-ACC-*)

## Search Behavior
- Search is **case-insensitive** (using PostgreSQL ILIKE)
- Searches across:
  - Product name
  - Product description
  - Product SKU
- Example: search="ayam" will match "Ayam Bakar", "nasi goreng dengan ayam", etc.

## Create/Update Product with Category
When creating or updating products, include the `category` field:

```json
{
  "name": "New Product",
  "description": "Product description",
  "category": "Makanan Utama",
  "sku": "SKU-001",
  "price": 50000,
  "stock": 100,
  "is_active": true
}
```

## Database Schema
```sql
ALTER TABLE products ADD COLUMN category VARCHAR(100);
CREATE INDEX idx_products_category ON products(category);
```

## Testing with cURL

```bash
# Login resto01
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "resto01",
    "branch_code": "resto01-pusat",
    "username": "tenantadmin",
    "password": "123456"
  }' | python3 -m json.tool | grep -o 'eyJ[^"]*' | head -1)

# Get categories
curl -X GET "http://localhost:8080/api/v1/products/categories" \
  -H "Authorization: Bearer $TOKEN"

# Filter by category
curl -X GET "http://localhost:8080/api/v1/products?category=Minuman" \
  -H "Authorization: Bearer $TOKEN"

# Search keyword
curl -X GET "http://localhost:8080/api/v1/products?search=ayam" \
  -H "Authorization: Bearer $TOKEN"

# Combined
curl -X GET "http://localhost:8080/api/v1/products?category=Makanan%20Utama&search=goreng" \
  -H "Authorization: Bearer $TOKEN"
```

## Implementation Files
- **Model:** `models/product.go` - Added Category field
- **DTO:** `dto/product.go` - Updated request/response DTOs
- **Service:** `services/product_service.go` - Filter & search logic
- **Handler:** `handlers/product_handler.go` - Query parameter handling
- **Routes:** `routes/routes.go` - New /categories endpoint

## Notes
- Categories are **tenant-isolated** (resto01 cannot see fashion01 categories)
- Empty category filter returns all products
- Empty search returns all products
- Both filters can be combined for precise filtering
