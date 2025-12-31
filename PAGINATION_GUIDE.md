# Pagination Guide - MyPOSCore API

Semua list API sekarang support pagination untuk meningkatkan performa dan pengalaman pengguna.

## Query Parameters

Semua list API menerima 2 query parameters untuk pagination:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman (mulai dari 1) |
| `page_size` | integer | 10 | Jumlah item per halaman (max: 100) |

## Response Format

Semua paginated response mengikuti format berikut:

```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 45,
  "total_pages": 5,
  "data": [...]
}
```

### Fields Explanation

- **page**: Halaman saat ini
- **page_size**: Jumlah item per halaman
- **total_items**: Total seluruh item dalam database
- **total_pages**: Total jumlah halaman
- **data**: Array berisi data item

## List APIs yang Support Pagination

1. **List Users** - `GET /api/v1/users`
2. **List Tenants** - `GET /api/v1/superadmin/tenants`
3. **List Branches** - `GET /api/v1/superadmin/tenants/:tenant_id/branches`
4. **List Products** - `GET /api/v1/products`
5. **List Categories** - `GET /api/v1/categories`

## Examples

### Example 1: List Users (Page 1, 10 items per page)

**Request:**
```
GET /api/v1/users?page=1&page_size=10
Authorization: Bearer <token>
```

**Response:** (200 OK)
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 25,
  "total_pages": 3,
  "data": [
    {
      "id": 18,
      "tenant_id": 17,
      "branch_id": 25,
      "email": "admin@foodcorner.com",
      "full_name": "Admin Food Corner",
      "image": "",
      "role": "tenantadmin",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 19,
      "tenant_id": 17,
      "branch_id": 25,
      "email": "branchadmin@foodcorner.com",
      "full_name": "Branch Admin Food Corner",
      "image": "",
      "role": "branchadmin",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    }
    // ... 8 more items
  ]
}
```

### Example 2: List Products with Filters (Page 2, 5 items per page)

**Request:**
```
GET /api/v1/products?category=Electronics&search=laptop&page=2&page_size=5
Authorization: Bearer <token>
```

**Response:** (200 OK)
```json
{
  "page": 2,
  "page_size": 5,
  "total_items": 12,
  "total_pages": 3,
  "data": [
    {
      "id": 6,
      "tenant_id": 17,
      "name": "Laptop Asus ROG",
      "description": "Gaming laptop with RTX 3060",
      "category": "Electronics",
      "sku": "LAP-ASUS-ROG-001",
      "price": 18000000,
      "stock": 5,
      "image": "http://localhost:8080/uploads/products/6_20241230_142208.jpg",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": 18,
      "created_by_name": "Admin Food Corner",
      "updated_by": null,
      "updated_by_name": null
    }
    // ... 4 more items
  ]
}
```

### Example 3: List Tenants (First Page with Default Page Size)

**Request:**
```
GET /api/v1/superadmin/tenants
Authorization: Bearer <superadmin_token>
```

**Response:** (200 OK)
```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 3,
  "total_pages": 1,
  "data": [
    {
      "id": 17,
      "name": "Fashion Store",
      "description": "Fashion and apparel retail chain",
      "address": "Jl. Sudirman No. 123",
      "website": "https://fashionstore.com",
      "email": "info@fashionstore.com",
      "phone": "+62812345678",
      "image": "http://localhost:8080/uploads/tenants/17_20241230_142208.jpg",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    },
    {
      "id": 18,
      "name": "Food Corner",
      "description": "Restaurant and cafe chain",
      "address": "Jl. Thamrin No. 456",
      "website": "https://foodcorner.com",
      "email": "info@foodcorner.com",
      "phone": "+62823456789",
      "image": "",
      "is_active": true,
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": null,
      "created_by_name": null,
      "updated_by": null,
      "updated_by_name": null
    }
    // ... more tenants if available
  ]
}
```

### Example 4: List Categories (Page 1, Only Active Categories)

**Request:**
```
GET /api/v1/categories?active_only=true&page=1&page_size=10
Authorization: Bearer <token>
```

**Response:** (200 OK)
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
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": 18,
      "created_by_name": "Admin Food Corner",
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
      "created_at": "2024-12-30 14:22:08",
      "updated_at": "2024-12-30 14:22:08",
      "created_by": 18,
      "created_by_name": "Admin Food Corner",
      "updated_by": null,
      "updated_by_name": null
    }
    // ... 6 more active categories
  ]
}
```

### Example 5: List Branches for Tenant (Page 1)

**Request:**
```
GET /api/v1/superadmin/tenants/17/branches?page=1&page_size=10
Authorization: Bearer <superadmin_token>
```

**Response:** (200 OK)
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

## Tips untuk Frontend Development

### Pagination Controls

```javascript
// Calculate pagination info from response
const response = await api.getUsers({ page: 1, page_size: 10 });

const paginationInfo = {
  currentPage: response.page,
  pageSize: response.page_size,
  totalItems: response.total_items,
  totalPages: response.total_pages,
  hasNextPage: response.page < response.total_pages,
  hasPrevPage: response.page > 1,
};

// Example pagination component
<Pagination
  current={paginationInfo.currentPage}
  total={paginationInfo.totalItems}
  pageSize={paginationInfo.pageSize}
  onChange={(page, pageSize) => fetchData(page, pageSize)}
/>
```

### Infinite Scroll

```javascript
let currentPage = 1;
let allData = [];

async function loadMore() {
  const response = await api.getProducts({ 
    page: currentPage, 
    page_size: 20 
  });
  
  allData = [...allData, ...response.data];
  currentPage++;
  
  // Check if more data available
  if (currentPage > response.total_pages) {
    // No more data
    return false;
  }
  return true;
}
```

### Table with Page Size Selector

```javascript
const [page, setPage] = useState(1);
const [pageSize, setPageSize] = useState(10);

// When page size changes, reset to page 1
const handlePageSizeChange = (newPageSize) => {
  setPageSize(newPageSize);
  setPage(1);
  fetchData(1, newPageSize);
};

// Pagination component
<Select value={pageSize} onChange={handlePageSizeChange}>
  <Option value={10}>10 per page</Option>
  <Option value={20}>20 per page</Option>
  <Option value={50}>50 per page</Option>
  <Option value={100}>100 per page</Option>
</Select>
```

## Error Handling

### Invalid Page Number

Jika request page number yang tidak valid (< 1), akan otomatis di-set ke 1.

### Invalid Page Size

- Page size < 1 akan di-set ke default (10)
- Page size > 100 akan di-set ke 100 (max limit)

### Empty Result

Jika tidak ada data pada halaman tertentu, response tetap valid dengan `data: []`:

```json
{
  "page": 5,
  "page_size": 10,
  "total_items": 12,
  "total_pages": 2,
  "data": []
}
```

## Performance Considerations

- Default page size (10) dioptimalkan untuk mobile devices
- Gunakan page size lebih besar (50-100) untuk desktop/admin panels
- Total count di-cache sehingga performa tetap baik untuk dataset besar
- Index database sudah dioptimalkan untuk queries dengan pagination
