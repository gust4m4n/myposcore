# User Management API Guide

API endpoints untuk manajemen user di MyPOSCore.

## 📋 Endpoints

### 1. Create User
**POST** `/api/v1/users`

Membuat user baru untuk tenant.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "full_name": "New User",
  "role": "user",
  "branch_id": 1,
  "is_active": true
}
```

**Field Validations:**
- `username` (required): Min 3 karakter, unique per tenant
- `email` (required): Valid email format, globally unique
- `password` (required): Min 6 karakter
- `full_name` (required): Nama lengkap user
- `role` (required): Harus salah satu dari: `user`, `branchadmin`, `tenantadmin`
- `branch_id` (required): ID branch yang valid, harus milik tenant
- `is_active` (optional): Default `true`

**Success Response (200):**
```json
{
  "message": "User created successfully",
  "data": {
    "id": 10,
    "tenant_id": 17,
    "branch_id": 1,
    "username": "newuser",
    "email": "newuser@example.com",
    "full_name": "New User",
    "role": "user",
    "is_active": true,
    "created_at": "2025-12-27 11:00:00"
  }
}
```

---

### 2. Get All Users
**GET** `/api/v1/users`

Mendapatkan daftar semua user untuk tenant.

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "tenant_id": 17,
      "branch_id": 1,
      "username": "user1",
      "email": "user1@example.com",
      "full_name": "User Satu",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-27 10:00:00"
    },
    {
      "id": 2,
      "tenant_id": 17,
      "branch_id": 1,
      "username": "admin1",
      "email": "admin1@example.com",
      "full_name": "Admin Satu",
      "role": "branchadmin",
      "is_active": true,
      "created_at": "2025-12-27 10:05:00"
    }
  ]
}
```

---

### 3. Get User by ID
**GET** `/api/v1/users/{id}`

Mendapatkan detail user berdasarkan ID.

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "data": {
    "id": 1,
    "tenant_id": 17,
    "branch_id": 1,
    "username": "user1",
    "email": "user1@example.com",
    "full_name": "User Satu",
    "role": "user",
    "is_active": true,
    "created_at": "2025-12-27 10:00:00"
  }
}
```

---

### 4. Update User
**PUT** `/api/v1/users/{id}`

Update data user yang sudah ada.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body (semua field optional):**
```json
{
  "username": "updateduser",
  "email": "updated@example.com",
  "password": "newpassword123",
  "full_name": "Updated User Name",
  "role": "branchadmin",
  "branch_id": 2,
  "is_active": true
}
```

**Notes:**
- Semua field bersifat **opsional**
- Hanya field yang dikirim yang akan diupdate
- Password akan otomatis di-hash
- Username tetap harus unique per tenant jika diubah
- Email tetap harus unique globally jika diubah
- Branch ID harus valid dan milik tenant

**Partial Update Examples:**

Hanya update role:
```json
{
  "role": "branchadmin"
}
```

Hanya update password:
```json
{
  "password": "newpassword456"
}
```

Update full name dan status:
```json
{
  "full_name": "John Doe Updated",
  "is_active": false
}
```

**Success Response (200):**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 17,
    "branch_id": 1,
    "username": "user1",
    "email": "user1@example.com",
    "full_name": "Updated User Name",
    "role": "branchadmin",
    "is_active": true,
    "created_at": "2025-12-27 10:00:00"
  }
}
```

---

### 5. Delete User
**DELETE** `/api/v1/users/{id}`

Hapus user (soft delete).

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "message": "User deleted successfully"
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
| 404 | Not Found - User tidak ditemukan |
| 500 | Internal Server Error |

## 👤 User Roles

| Role | Description |
|------|-------------|
| `user` | User biasa dengan akses terbatas |
| `branchadmin` | Admin level branch, bisa manage resource di branchnya |
| `tenantadmin` | Admin level tenant, bisa manage semua branch |
| `superadmin` | Super admin sistem (tidak bisa dibuat via API ini) |

## 🔐 Security Features

- **Password Hashing**: Password otomatis di-hash menggunakan bcrypt
- **Tenant Isolation**: User otomatis difilter berdasarkan `tenant_id` dari JWT token
- **Unique Constraints**: 
  - Username unique per tenant
  - Email unique globally
- **Branch Validation**: Branch ID harus valid dan milik tenant
- **Soft Delete**: User yang dihapus tidak benar-benar terhapus dari database

## 💡 Notes

- User otomatis difilter berdasarkan `tenant_id` dari JWT token
- User yang dihapus menggunakan soft delete (tidak dihapus dari database)
- Password tidak pernah di-return dalam response (hidden di model)
- Username harus unique dalam satu tenant (bisa sama di tenant berbeda)
- Email harus unique global (tidak boleh duplicate di seluruh sistem)

## 📝 Common Error Messages

### Username Already Exists
```json
{
  "error": "username already exists for this tenant"
}
```

### Email Already Exists
```json
{
  "error": "email already exists"
}
```

### Invalid Branch
```json
{
  "error": "branch not found or doesn't belong to this tenant"
}
```

### User Not Found
```json
{
  "error": "user not found"
}
```

### Validation Error
```json
{
  "error": "Key: 'CreateUserRequest.Username' Error:Field validation for 'Username' failed on the 'min' tag"
}
```

## 🔗 Related APIs

- **Auth API**: `/api/v1/auth/register` - Register user baru
- **Profile API**: `/api/v1/profile` - Get user profile dari token
- **Change Password API**: `/api/v1/change-password` - Change password user sendiri

## 📌 Example Workflows

### 1. Create New Branch Admin
```bash
# Step 1: Login as tenantadmin
POST /api/v1/auth/login
{
  "username": "tenantadmin",
  "password": "password123"
}

# Step 2: Create new user with branchadmin role
POST /api/v1/users
{
  "username": "branchadmin1",
  "email": "branchadmin1@example.com",
  "password": "password123",
  "full_name": "Branch Admin 1",
  "role": "branchadmin",
  "branch_id": 1
}
```

### 2. Deactivate User
```bash
PUT /api/v1/users/5
{
  "is_active": false
}
```

### 3. Change User Role
```bash
PUT /api/v1/users/5
{
  "role": "tenantadmin"
}
```

### 4. Reset User Password
```bash
PUT /api/v1/users/5
{
  "password": "newpassword123"
}
```

### 5. Move User to Different Branch
```bash
PUT /api/v1/users/5
{
  "branch_id": 2
}
```
