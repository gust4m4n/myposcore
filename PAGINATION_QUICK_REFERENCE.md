# Pagination Quick Reference

## 📋 Quick Start

### Using Pagination in API Calls

```bash
# Default pagination (page 1, 10 items)
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/products

# Custom page and size
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/products?page=2&page_size=20

# With filters
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/products?category=Minuman&page=1&page_size=5
```

### Response Structure

```json
{
  "page": 1,           // Current page number
  "page_size": 10,     // Items per page
  "total_items": 45,   // Total number of items
  "total_pages": 5,    // Total number of pages
  "data": [...]        // Array of items
}
```

---

## 🎯 All Paginated Endpoints

| Endpoint | Method | Auth Required | Default Size | Filters |
|----------|--------|---------------|--------------|---------|
| `/api/v1/users` | GET | ✅ User | 10 | - |
| `/api/v1/categories` | GET | ✅ User | 10 | `active_only` |
| `/api/v1/products` | GET | ✅ User | 10 | `category`, `search` |
| `/api/v1/superadmin/tenants` | GET | ✅ Superadmin | 10 | - |
| `/api/v1/superadmin/tenants/:id/branches` | GET | ✅ Superadmin | 10 | - |

---

## ⚙️ Query Parameters

| Parameter | Type | Default | Min | Max | Description |
|-----------|------|---------|-----|-----|-------------|
| `page` | int | 1 | 1 | ∞ | Page number |
| `page_size` | int | 10 | 1 | 100 | Items per page |

---

## 💻 Frontend Code Examples

### JavaScript/Fetch

```javascript
async function fetchProducts(page = 1, pageSize = 10) {
  const response = await fetch(
    `/api/v1/products?page=${page}&page_size=${pageSize}`,
    {
      headers: { 'Authorization': `Bearer ${token}` }
    }
  );
  const data = await response.json();
  
  return {
    items: data.data,
    pagination: {
      page: data.page,
      pageSize: data.page_size,
      totalItems: data.total_items,
      totalPages: data.total_pages
    }
  };
}
```

### Axios

```javascript
const fetchProducts = async (page = 1, pageSize = 10) => {
  const { data } = await axios.get('/api/v1/products', {
    params: { page, page_size: pageSize },
    headers: { Authorization: `Bearer ${token}` }
  });
  return data;
};
```

### React Hook

```jsx
function useProducts(page = 1, pageSize = 10) {
  const [data, setData] = useState({ items: [], pagination: {} });
  const [loading, setLoading] = useState(false);
  
  useEffect(() => {
    setLoading(true);
    fetchProducts(page, pageSize)
      .then(setData)
      .finally(() => setLoading(false));
  }, [page, pageSize]);
  
  return { ...data, loading };
}

// Usage
function ProductList() {
  const [page, setPage] = useState(1);
  const { items, pagination, loading } = useProducts(page, 10);
  
  return (
    <div>
      {loading ? <Spinner /> : <ProductGrid products={items} />}
      <Pagination
        current={pagination.page}
        total={pagination.totalItems}
        pageSize={pagination.pageSize}
        onChange={setPage}
      />
    </div>
  );
}
```

---

## 🔧 Backend Implementation Pattern

### Service Layer

```go
func (s *Service) ListItems(tenantID, page, pageSize int) ([]Item, int64, error) {
    // Step 1: Count total
    var total int64
    query := s.db.Model(&Item{}).Where("tenant_id = ?", tenantID)
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Step 2: Get page data
    offset := (page - 1) * pageSize
    var items []Item
    if err := s.db.Where("tenant_id = ?", tenantID).
        Limit(pageSize).
        Offset(offset).
        Find(&items).Error; err != nil {
        return nil, 0, err
    }
    
    return items, total, nil
}
```

### Handler Layer

```go
func (h *Handler) ListItems(c *gin.Context) {
    // Parse pagination
    var pagination dto.PaginationRequest
    if err := c.ShouldBindQuery(&pagination); err != nil {
        pagination = *dto.NewPaginationRequest(1, 10)
    } else {
        pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
    }
    
    // Get tenant ID from JWT
    tenantID := c.GetUint("tenant_id")
    
    // Call service
    items, total, err := h.service.ListItems(tenantID, pagination.Page, pagination.PageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Build response
    response := buildItemsResponse(items)
    paginatedResponse := dto.NewPaginationResponse(
        pagination.Page,
        pagination.PageSize,
        total,
        response,
    )
    
    c.JSON(http.StatusOK, paginatedResponse)
}
```

---

## 🧪 Testing Examples

### Test Default Pagination

```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer TOKEN"
# Expected: page=1, page_size=10
```

### Test Custom Page Size

```bash
curl -X GET "http://localhost:8080/api/v1/products?page=2&page_size=50" \
  -H "Authorization: Bearer TOKEN"
# Expected: page=2, page_size=50
```

### Test Max Limit

```bash
curl -X GET "http://localhost:8080/api/v1/products?page_size=200" \
  -H "Authorization: Bearer TOKEN"
# Expected: page_size capped at 100
```

### Test Invalid Page

```bash
curl -X GET "http://localhost:8080/api/v1/products?page=0" \
  -H "Authorization: Bearer TOKEN"
# Expected: page defaults to 1
```

### Test With Filters

```bash
curl -X GET "http://localhost:8080/api/v1/products?category=Minuman&search=kopi&page=1&page_size=5" \
  -H "Authorization: Bearer TOKEN"
# Expected: Filtered results with pagination
```

---

## 📊 Pagination Math

### Calculate Total Pages

```javascript
const totalPages = Math.ceil(totalItems / pageSize);
```

### Calculate Current Offset

```javascript
const offset = (page - 1) * pageSize;
```

### Check if Has Next Page

```javascript
const hasNextPage = page < totalPages;
```

### Check if Has Previous Page

```javascript
const hasPrevPage = page > 1;
```

---

## ⚡ Performance Tips

### Frontend
- Cache results per page
- Prefetch next page
- Use virtual scrolling for large lists
- Show loading state during fetch
- Implement optimistic updates

### Backend
- Index commonly filtered columns
- Use database query caching
- Monitor slow queries
- Consider Redis for COUNT cache
- Optimize Preload queries

---

## 🐛 Common Issues & Solutions

### Issue: Response still shows old format `{data: []}`

**Solution:** Update backend handlers to use `dto.PaginationResponse`

```go
// Wrong
c.JSON(http.StatusOK, gin.H{"data": items})

// Correct
c.JSON(http.StatusOK, dto.NewPaginationResponse(page, pageSize, total, items))
```

### Issue: `page_size` not respected

**Solution:** Ensure validation in DTO

```go
pagination := *dto.NewPaginationRequest(page, pageSize)
// This automatically validates and caps page_size
```

### Issue: Frontend shows wrong total pages

**Solution:** Calculate on backend, not frontend

```go
// Backend does this automatically
totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
```

### Issue: Slow queries with large offsets

**Solution:** Consider cursor-based pagination for very large datasets

```sql
-- Instead of OFFSET 10000 LIMIT 10
-- Use cursor-based:
WHERE id > last_seen_id LIMIT 10
```

---

## 📝 Best Practices

### DO ✅
- Always validate page and page_size
- Return total_items for UI pagination controls
- Cache COUNT queries when possible
- Use consistent default page_size across endpoints
- Document breaking changes clearly

### DON'T ❌
- Don't allow unlimited page_size
- Don't forget to add indexes on filtered columns
- Don't use OFFSET for real-time data (use cursors)
- Don't return all records in one request
- Don't forget to test edge cases (page=0, empty results)

---

## 📖 Documentation Links

- **Full Guide:** [PAGINATION_GUIDE.md](./PAGINATION_GUIDE.md)
- **Postman Update:** [POSTMAN_PAGINATION_UPDATE.md](./POSTMAN_PAGINATION_UPDATE.md)
- **Implementation Summary:** [PAGINATION_IMPLEMENTATION_SUMMARY.md](./PAGINATION_IMPLEMENTATION_SUMMARY.md)

---

## 🚀 Quick Deploy Checklist

- [ ] Backend code updated with pagination
- [ ] Server compiled and tested locally
- [ ] Postman collection updated
- [ ] Frontend API client updated
- [ ] UI pagination controls implemented
- [ ] All endpoints tested with various scenarios
- [ ] Documentation updated
- [ ] Team notified of breaking changes
- [ ] Deployed to staging
- [ ] QA testing completed
- [ ] Production deployment scheduled

---

**Last Updated:** December 31, 2024  
**Version:** 1.0  
**Status:** ✅ Production Ready
