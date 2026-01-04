# Branch API - Complete Guide

## Overview
MyPOSCore memiliki beberapa endpoint untuk mengelola branches (cabang):
1. **User Branch API** - untuk user biasa (tenant admin/branch admin)
2. **Superadmin Branch API** - untuk superadmin dengan akses penuh
3. **Dev Branch API** - untuk development/testing (tanpa autentikasi)

---

## 1. USER BRANCH API

### GET /api/v1/branches
**Deskripsi**: Get branches dari tenant user yang sedang login
**Authentication**: Required (Bearer Token)
**Parameter**: Tidak ada (otomatis baca dari JWT token)

**Response Codes:**
- `200` OK - Berhasil
- `401` Unauthorized - Token tidak valid/tidak ada

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/branches \
  -H "Authorization: Bearer <token>"
```

**Example Response (Food Corner):**
```json
{
  "code": 0,
  "message": "Branches retrieved successfully",
  "data": [
    {
      "ID": 26,
      "CreatedAt": "2025-12-24T19:44:27.922193+07:00",
      "UpdatedAt": "2026-01-01T17:34:22.657741+07:00",
      "DeletedAt": null,
      "tenant_id": 17,
      "name": "Cabang Menteng",
      "description": "This is an UPDATED test branch description...",
      "address": "Jl. Kuningan Raya No. 789, Jakarta Selatan",
      "website": "https://www.updated-testbranch.com",
      "email": "updated@testbranch.com",
      "phone": "021-11223344",
      "image": "/uploads/branches/branch_26_1767263662.jpg",
      "is_active": true,
      "updated_by": 18
    },
    {
      "ID": 25,
      "CreatedAt": "2025-12-24T19:44:27.922193+07:00",
      "UpdatedAt": "2026-01-01T17:34:06.163744+07:00",
      "DeletedAt": null,
      "tenant_id": 17,
      "name": "Cabang Pondok Indah",
      "description": "This is an UPDATED test branch description...",
      "address": "Jl. Kuningan Raya No. 789, Jakarta Selatan",
      "website": "https://www.updated-testbranch.com",
      "email": "updated@testbranch.com",
      "phone": "021-11223344",
      "image": "/uploads/branches/branch_25_1767263646.jpg",
      "is_active": true,
      "updated_by": 18
    }
  ]
}
```

---

## 2. SUPERADMIN BRANCH API

### GET /api/v1/superadmin/tenants/:tenant_id/branches
**Deskripsi**: List semua branches dari tenant tertentu
**Authentication**: Required (Superadmin only)
**Parameters:**
- `:tenant_id` (path) - ID tenant

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/superadmin/tenants/17/branches \
  -H "Authorization: Bearer <superadmin_token>"
```

### POST /api/v1/superadmin/branches
**Deskripsi**: Create branch baru untuk tenant
**Authentication**: Required (Superadmin only)
**Content-Type**: multipart/form-data

**Form Fields:**
- `tenant_id` (required) - ID tenant
- `name` (required) - Nama branch
- `description` (optional) - Deskripsi branch
- `address` (optional) - Alamat lengkap
- `website` (optional) - Website URL
- `email` (optional) - Email branch
- `phone` (optional) - Nomor telepon
- `is_active` (optional) - Status aktif (true/false)
- `image` (optional) - File upload gambar branch (jpg, jpeg, png, gif, webp, max 5MB)

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/superadmin/branches \
  -H "Authorization: Bearer <superadmin_token>" \
  -F "tenant_id=17" \
  -F "name=Cabang Baru" \
  -F "description=Cabang yang baru dibuka" \
  -F "address=Jl. Sudirman No. 123" \
  -F "phone=021-12345678" \
  -F "is_active=true" \
  -F "image=@/path/to/branch-logo.png"
```

**Example Response:**
```json
{
  "code": 0,
  "message": "Branch created successfully",
  "data": {
    "id": 31,
    "tenant_id": 17,
    "name": "Cabang Baru",
    "description": "Cabang yang baru dibuka",
    "address": "Jl. Sudirman No. 123",
    "phone": "021-12345678",
    "image": "/uploads/branches/branch_31_1704000000.jpg",
    "is_active": true,
    "created_at": "2026-01-04T21:30:00+07:00",
    "updated_at": "2026-01-04T21:30:00+07:00",
    "created_by": 1
  }
}
```

### PUT /api/v1/superadmin/branches/:branch_id
**Deskripsi**: Update data branch
**Authentication**: Required (Superadmin only)
**Content-Type**: multipart/form-data
**Parameters:**
- `:branch_id` (path) - ID branch yang akan diupdate

**Form Fields:** (sama dengan POST, semua optional)
- `name` - Nama branch
- `description` - Deskripsi branch
- `address` - Alamat lengkap
- `website` - Website URL
- `email` - Email branch
- `phone` - Nomor telepon
- `is_active` - Status aktif (true/false)
- `image` - File upload gambar branch baru

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/superadmin/branches/31 \
  -H "Authorization: Bearer <superadmin_token>" \
  -F "name=Cabang Baru Updated" \
  -F "phone=021-99999999" \
  -F "image=@/path/to/new-logo.png"
```

**Example Response:**
```json
{
  "code": 0,
  "message": "Branch updated successfully",
  "data": {
    "id": 31,
    "tenant_id": 17,
    "name": "Cabang Baru Updated",
    "description": "Cabang yang baru dibuka",
    "address": "Jl. Sudirman No. 123",
    "phone": "021-99999999",
    "image": "/uploads/branches/branch_31_1704000100.jpg",
    "is_active": true,
    "created_at": "2026-01-04T21:30:00+07:00",
    "updated_at": "2026-01-04T21:35:00+07:00",
    "created_by": 1,
    "updated_by": 1
  }
}
```

### DELETE /api/v1/superadmin/branches/:branch_id
**Deskripsi**: Soft delete branch (tetap di database tapi ditandai deleted)
**Authentication**: Required (Superadmin only)
**Parameters:**
- `:branch_id` (path) - ID branch yang akan dihapus

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/superadmin/branches/31 \
  -H "Authorization: Bearer <superadmin_token>"
```

**Example Response:**
```json
{
  "code": 0,
  "message": "Branch deleted successfully"
}
```

### GET /api/v1/superadmin/branches/:branch_id/users
**Deskripsi**: List semua users yang terdaftar di branch tertentu
**Authentication**: Required (Superadmin only)
**Parameters:**
- `:branch_id` (path) - ID branch

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/superadmin/branches/25/users \
  -H "Authorization: Bearer <superadmin_token>"
```

**Example Response:**
```json
{
  "code": 0,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 60,
      "tenant_id": 17,
      "branch_id": 25,
      "email": "admin@foodcorner.com",
      "full_name": "Admin Food Corner 860",
      "role": "tenantadmin",
      "is_active": true,
      "created_at": "2025-12-24T19:44:27+07:00"
    }
  ]
}
```

---

## 3. DEV BRANCH API (Testing Only)

### GET /dev/tenants/:tenant_id/branches
**Deskripsi**: Get active branches dari tenant (tanpa autentikasi, untuk testing)
**Authentication**: None
**Parameters:**
- `:tenant_id` (path) - ID tenant

**Example Request:**
```bash
curl http://localhost:8080/dev/tenants/18/branches
```

**Example Response:**
```json
{
  "code": 0,
  "message": "Branches retrieved successfully",
  "data": [
    {
      "ID": 29,
      "CreatedAt": "2025-12-24T19:48:10.083674+07:00",
      "UpdatedAt": "2026-01-01T17:15:14.856775+07:00",
      "DeletedAt": null,
      "tenant_id": 18,
      "name": "Cabang Bintaro Plaza",
      "description": "This is an UPDATED test branch description...",
      "address": "Jl. Kuningan Raya No. 789, Jakarta Selatan",
      "website": "https://www.updated-testbranch.com",
      "email": "updated@testbranch.com",
      "phone": "021-11223344",
      "image": "/uploads/branches/branch_29_1767262514.jpg",
      "is_active": true,
      "updated_by": 18
    },
    {
      "ID": 30,
      "CreatedAt": "2025-12-24T19:48:10.083674+07:00",
      "UpdatedAt": "2026-01-01T17:24:36.112272+07:00",
      "DeletedAt": null,
      "tenant_id": 18,
      "name": "Cabang Grand Indonesia",
      "description": "This is an UPDATED test branch description...",
      "address": "Jl. Kuningan Raya No. 789, Jakarta Selatan",
      "website": "https://www.updated-testbranch.com",
      "email": "updated@testbranch.com",
      "phone": "021-11223344",
      "image": "/uploads/branches/branch_30_1767263076.jpg",
      "is_active": true,
      "updated_by": 18
    }
  ]
}
```

---

## Response Code Reference

### Success Codes
- `0` - Success (operasi berhasil)

### Error Codes
- `1` - Bad Request (400) - Request tidak valid
- `2` - Unauthorized (401) - Token tidak valid/tidak ada
- `3` - Forbidden (403) - Tidak punya akses
- `4` - Not Found (404) - Resource tidak ditemukan
- `5` - Internal Server Error (500) - Error di server
- `6` - Conflict (409) - Data conflict (misal: duplicate)
- `7` - Unprocessable Entity (422) - Validasi gagal

---

## Testing Workflow

### 1. Test User Branch API
```bash
# Login sebagai admin tenant
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@foodcorner.com","password":"123456"}' \
  2>/dev/null | jq -r '.data.token')

# Get branches (akan tampil branch dari tenant Food Corner saja)
curl -X GET http://localhost:8080/api/v1/branches \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 2. Test dengan Tenant Berbeda
```bash
# Login sebagai admin tenant lain
TOKEN2=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin.plaza@fashionstore.com","password":"123456"}' \
  2>/dev/null | jq -r '.data.token')

# Get branches (akan tampil branch dari tenant Fashion Store saja)
curl -X GET http://localhost:8080/api/v1/branches \
  -H "Authorization: Bearer $TOKEN2" | jq
```

### 3. Test Dev API (Tanpa Token)
```bash
# List tenants
curl http://localhost:8080/dev/tenants | jq

# List branches tenant 18 (Fashion Store)
curl http://localhost:8080/dev/tenants/18/branches | jq
```

---

## Postman Collection Status

✅ **Semua endpoint sudah ada di Postman dengan contoh response:**

1. **User Branch API**
   - ✅ GET /api/v1/branches (3 contoh response)

2. **Superadmin Branch API**
   - ✅ GET /api/v1/superadmin/tenants/:tenant_id/branches
   - ✅ POST /api/v1/superadmin/branches (dengan upload image)
   - ✅ PUT /api/v1/superadmin/branches/:branch_id (dengan upload image)
   - ✅ DELETE /api/v1/superadmin/branches/:branch_id
   - ✅ GET /api/v1/superadmin/branches/:branch_id/users

3. **Dev Branch API**
   - ✅ GET /dev/tenants/:tenant_id/branches

**Total: 6 endpoint branches tersedia**
**Validation: 80/80 responses valid dengan code & message**

---

## Notes

1. **Tenant Isolation**: User branch API (`GET /api/v1/branches`) otomatis filter berdasarkan tenant user yang login
2. **Image Upload**: Create & Update branch support multipart upload dengan validasi ukuran (max 5MB) dan format (jpg, jpeg, png, gif, webp)
3. **Soft Delete**: DELETE endpoint tidak benar-benar menghapus data, hanya menandai sebagai deleted
4. **Response Format**: Semua endpoint menggunakan format standard dengan `code`, `message`, dan `data`
5. **Dev Endpoints**: Hanya untuk development/testing, tidak boleh dipakai di production
