# Superadmin CRUD API Documentation

## Overview
Complete CRUD operations for Tenant and Branch management. All endpoints require superadmin authentication.

## Authentication
All endpoints require:
- **Header**: `Authorization: Bearer <token>`
- **Role**: `superadmin`

---

## Tenant Management

### 1. List All Tenants
**GET** `/api/v1/superadmin/tenants`

**Response Success (200):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Warteg 123",
      "code": "warteg123",
      "description": "Warung Tegal dengan menu tradisional Indonesia",
      "address": "Jl. Raya Sudirman No. 123, Jakarta Pusat",
      "website": "https://warteg123.com",
      "email": "info@warteg123.com",
      "phone": "021-12345678",
      "is_active": true,
      "created_at": "2025-12-24T10:00:00+07:00"
    }
  ]
}
```

---

### 2. Create Tenant
**POST** `/api/v1/superadmin/tenants`

**Request Body:**
```json
{
  "name": "New Tenant",
  "code": "TENANT002",
  "description": "A new business tenant",
  "address": "Jl. Kebon Jeruk No. 789, Jakarta",
  "website": "https://newtenant.com",
  "email": "contact@newtenant.com",
  "phone": "021-99887766",
  "is_active": true
}
```

**Response Success (200):**
```json
{
  "message": "Tenant created successfully",
  "data": {
    "id": 19,
    "name": "New Tenant",
    "code": "TENANT002",
    "description": "A new business tenant",
    "address": "Jl. Kebon Jeruk No. 789, Jakarta",
    "website": "https://newtenant.com",
    "email": "contact@newtenant.com",
    "phone": "021-99887766",
    "is_active": true,
    "created_at": "2025-12-30T16:00:00+07:00"
  }
}
```

---

### 3. Update Tenant
**PUT** `/api/v1/superadmin/tenants/:tenant_id`

**Request Body:**
```json
{
  "name": "Updated Tenant Name",
  "code": "TENANT002",
  "description": "Updated business description",
  "address": "New Address, Jakarta",
  "website": "https://updatedtenant.com",
  "email": "newemail@tenant.com",
  "phone": "021-11112222",
  "is_active": true
}
```

**Response Success (200):**
```json
{
  "message": "Tenant updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Tenant Name",
    "code": "TENANT002",
    "description": "Updated business description",
    "address": "New Address, Jakarta",
    "website": "https://updatedtenant.com",
    "email": "newemail@tenant.com",
    "phone": "021-11112222",
    "is_active": true,
    "created_at": "2025-12-30T10:00:00+07:00"
  }
}
```

---

### 4. Delete Tenant
**DELETE** `/api/v1/superadmin/tenants/:tenant_id`

**Response Success (200):**
```json
{
  "message": "Tenant deleted successfully"
}
```

---

## Branch Management

### 1. List Branches by Tenant
**GET** `/api/v1/superadmin/tenants/:tenant_id/branches`

**Response Success (200):**
```json
{
  "data": [
    {
      "id": 1,
      "tenant_id": 1,
      "name": "Cabang Pusat",
      "code": "pusat",
      "description": "Kantor pusat dan cabang utama",
      "address": "Jl. Sudirman No. 123, Jakarta Pusat",
      "website": "https://warteg123.com/pusat",
      "email": "pusat@warteg123.com",
      "phone": "021-12345678",
      "is_active": true,
      "created_at": "2025-12-24T10:00:00+07:00"
    }
  ]
}
```

---

### 2. Create Branch
**POST** `/api/v1/superadmin/branches`

**Request Body:**
```json
{
  "tenant_id": 1,
  "name": "Cabang Baru",
  "code": "baru",
  "description": "Cabang baru yang baru dibuka",
  "address": "Jl. Gatot Subroto No. 99, Jakarta",
  "website": "https://warteg123.com/baru",
  "email": "baru@warteg123.com",
  "phone": "021-55556666",
  "is_active": true
}
```

**Response Success (200):**
```json
{
  "message": "Branch created successfully",
  "data": {
    "id": 5,
    "tenant_id": 1,
    "name": "Cabang Baru",
    "code": "baru",
    "description": "Cabang baru yang baru dibuka",
    "address": "Jl. Gatot Subroto No. 99, Jakarta",
    "website": "https://warteg123.com/baru",
    "email": "baru@warteg123.com",
    "phone": "021-55556666",
    "is_active": true,
    "created_at": "2025-12-30T16:30:00+07:00"
  }
}
```

---

### 3. Update Branch
**PUT** `/api/v1/superadmin/branches/:branch_id`

**Request Body:**
```json
{
  "name": "Updated Branch Name",
  "code": "baru",
  "description": "Updated branch description",
  "address": "Updated Address, Jakarta",
  "website": "https://warteg123.com/updated",
  "email": "updated@warteg123.com",
  "phone": "021-99998888",
  "is_active": true
}
```

**Response Success (200):**
```json
{
  "message": "Branch updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Updated Branch Name",
    "code": "baru",
    "description": "Updated branch description",
    "address": "Updated Address, Jakarta",
    "website": "https://warteg123.com/updated",
    "email": "updated@warteg123.com",
    "phone": "021-99998888",
    "is_active": true,
    "created_at": "2025-12-24T10:00:00+07:00"
  }
}
```

---

### 4. Delete Branch
**DELETE** `/api/v1/superadmin/branches/:branch_id`

**Response Success (200):**
```json
{
  "message": "Branch deleted successfully"
}
```

---

## Field Descriptions

### Tenant Fields
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Tenant business name |
| code | string | Yes | Unique tenant code (alphanumeric, no spaces) |
| description | string | No | Business description |
| address | string | No | Business address |
| website | string | No | Business website URL |
| email | string | No | Contact email |
| phone | string | No | Contact phone number |
| is_active | boolean | Yes | Active status |

### Branch Fields
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| tenant_id | uint | Yes | Parent tenant ID (only for create) |
| name | string | Yes | Branch name |
| code | string | Yes | Unique branch code per tenant |
| description | string | No | Branch description |
| address | string | No | Branch address |
| website | string | No | Branch website URL |
| email | string | No | Branch contact email |
| phone | string | No | Branch contact phone |
| is_active | boolean | Yes | Active status |

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "tenant code already exists"
}
```

### 404 Not Found
```json
{
  "error": "tenant not found"
}
```

### 401 Unauthorized
```json
{
  "error": "unauthorized"
}
```

### 403 Forbidden
```json
{
  "error": "superadmin access required"
}
```

---

## Notes

1. **Soft Delete**: Delete operations are soft deletes - records are marked as deleted but remain in database
2. **Code Uniqueness**: 
   - Tenant codes must be unique across all tenants
   - Branch codes must be unique within the same tenant
3. **Field Updates**: When updating, all fields must be provided (full object replacement)
4. **Login Response**: Login endpoint now returns tenant and branch with all new fields (description, address, website, email, phone)

---

## Testing with Postman

Import `MyPOSCore.postman_collection.json` which includes:
- All CRUD operations for Tenant and Branch
- Example requests with complete field data
- Success response examples
- Pre-configured authentication

**Demo Superadmin Credentials:**
- Username: `admin@mypos.com`
- Password: `123456`
