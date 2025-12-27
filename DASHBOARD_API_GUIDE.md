# Dashboard API Guide

## Overview
Dashboard API for superadmin to view comprehensive system statistics including:
- Total tenants and their list
- Total branches
- Total users
- Total products
- Total orders (all time, today, this week, this month)

**Base URL:** `/api/v1/superadmin`

**Authentication:** Requires Bearer token with superadmin role

---

## Get Dashboard

Get comprehensive dashboard statistics with all system metrics.

### Endpoint
```
GET /api/v1/superadmin/dashboard
```

### Authentication
```
Authorization: Bearer <token>
```

**Required Role:** `superadmin`

### Response

#### Success Response (200 OK)
```json
{
  "data": {
    "total_tenants": 5,
    "total_branches": 12,
    "total_users": 45,
    "total_products": 234,
    "total_orders": 1523,
    "total_orders_today": 23,
    "total_orders_this_week": 156,
    "total_orders_this_month": 678,
    "tenants": [
      {
        "id": 1,
        "name": "Warung Makan Sederhana",
        "code": "WMS001",
        "is_active": true,
        "created_at": "2025-01-15 10:30:00"
      },
      {
        "id": 2,
        "name": "Kopi Kita",
        "code": "KK002",
        "is_active": true,
        "created_at": "2025-01-16 14:20:00"
      },
      {
        "id": 3,
        "name": "Toko Kelontong Jaya",
        "code": "TKJ003",
        "is_active": true,
        "created_at": "2025-01-18 09:15:00"
      },
      {
        "id": 4,
        "name": "Resto Seafood",
        "code": "RS004",
        "is_active": false,
        "created_at": "2025-01-20 16:45:00"
      },
      {
        "id": 5,
        "name": "Bakery Fresh",
        "code": "BF005",
        "is_active": true,
        "created_at": "2025-01-22 11:00:00"
      }
    ]
  }
}
```

#### Error Response (401 Unauthorized)
```json
{
  "error": "Unauthorized: Superadmin access required"
}
```

#### Error Response (500 Internal Server Error)
```json
{
  "error": "Failed to fetch dashboard data"
}
```

---

## Response Fields

### DashboardResponse
| Field | Type | Description |
|-------|------|-------------|
| `total_tenants` | integer | Total number of tenants in the system |
| `total_branches` | integer | Total number of branches across all tenants |
| `total_users` | integer | Total number of users across all tenants |
| `total_products` | integer | Total number of products across all tenants |
| `total_orders` | integer | Total number of orders (all time) |
| `total_orders_today` | integer | Total orders created today (00:00 - 23:59) |
| `total_orders_this_week` | integer | Total orders created this week (Sunday - Saturday) |
| `total_orders_this_month` | integer | Total orders created this month (1st - last day) |
| `tenants` | array | List of all tenants with details |

### TenantResponse (in tenants array)
| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Tenant ID |
| `name` | string | Tenant name |
| `code` | string | Tenant unique code |
| `is_active` | boolean | Tenant active status |
| `created_at` | string | Creation timestamp |

---

## Usage Examples

### Using cURL
```bash
curl -X GET http://localhost:8080/api/v1/superadmin/dashboard \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Using Postman
1. Open Postman
2. Create a new GET request to `{{base_url}}/api/v1/superadmin/dashboard`
3. Go to Authorization tab
4. Select "Bearer Token"
5. Paste your superadmin token in the Token field
6. Click Send

### Using JavaScript (fetch)
```javascript
fetch('http://localhost:8080/api/v1/superadmin/dashboard', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer YOUR_SUPERADMIN_TOKEN_HERE'
  }
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```

---

## Notes

1. **Superadmin Access Only**: This endpoint requires superadmin role. Regular users and branch admins cannot access this endpoint.

2. **Time Zone**: All time-based calculations (today, this week, this month) are based on the server's local timezone.

3. **Week Calculation**: Week starts on Sunday (0) and ends on Saturday (6).

4. **Performance**: This endpoint aggregates data from multiple tables. Consider caching for high-traffic scenarios.

5. **Tenant List**: The tenants array includes both active and inactive tenants. Use the `is_active` field to filter.

---

## Common Use Cases

### 1. System Overview
Use this endpoint on the superadmin dashboard landing page to show overall system health and activity.

### 2. Business Analytics
Monitor business growth by tracking:
- New tenant registrations
- Daily/weekly/monthly order trends
- Product catalog expansion

### 3. Tenant Management
Quickly see all tenants and their status without making separate API calls.

### 4. Performance Monitoring
Track system usage patterns:
- Peak order periods (compare today vs week vs month)
- User adoption rate across tenants
- Product catalog growth

---

## Security Considerations

1. **Role Verification**: Always ensure the requesting user has superadmin role
2. **Token Validation**: Verify JWT token signature and expiration
3. **Rate Limiting**: Consider implementing rate limiting for this endpoint to prevent abuse
4. **Data Privacy**: This endpoint exposes system-wide data. Ensure proper access controls are in place.
