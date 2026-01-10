# Admin Change Password API Guide

API untuk mengubah password user lain (role tinggi ke role rendah).

## Role Hierarchy

```
superadmin > owner > admin > user
```

- **superadmin**: Dapat mengubah password semua user
- **owner**: Dapat mengubah password admin dan user di branch-nya
- **admin**: Dapat mengubah password user di tenant-nya
- **user**: Tidak dapat mengubah password user lain

## Endpoint

### Change Password by Admin

Mengubah password user lain berdasarkan username.

**URL**: `/api/admin/change-password`

**Method**: `PUT`

**Auth Required**: Yes (Bearer Token)

**Permissions**: 
- Superadmin dapat mengubah password semua user
- Owner dapat mengubah password admin dan user di branch yang sama
- Admin dapat mengubah password user di tenant yang sama

### Request Body

```json
{
  "username": "user123",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | Yes | Username dari user yang akan diubah passwordnya |
| password | string | Yes | Password baru (minimal 6 karakter) |
| confirm_password | string | Yes | Konfirmasi password baru (harus sama dengan password) |

### Success Response

**Code**: `200 OK`

**Content Example**:

```json
{
  "message": "Password changed successfully"
}
```

### Error Responses

**Code**: `400 BAD REQUEST`

**Content Examples**:

```json
{
  "error": "password and confirm password do not match"
}
```

```json
{
  "error": "password must be at least 6 characters"
}
```

```json
{
  "error": "target user not found"
}
```

```json
{
  "error": "insufficient permission: can only change password for lower role users"
}
```

```json
{
  "error": "cannot change password for users from different tenant"
}
```

```json
{
  "error": "owner can only change password for users in their branch"
}
```

**Code**: `401 UNAUTHORIZED`

**Content Example**:

```json
{
  "error": "User not authenticated"
}
```

## Usage Examples

### Example 1: Superadmin Mengubah Password Owner

**Request:**

```bash
curl -X PUT http://localhost:8080/api/admin/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_SUPERADMIN_TOKEN" \
  -d '{
    "username": "owner1",
    "password": "newpass123",
    "confirm_password": "newpass123"
  }'
```

**Response:**

```json
{
  "message": "Password changed successfully"
}
```

### Example 2: Owner Mengubah Password User

**Request:**

```bash
curl -X PUT http://localhost:8080/api/admin/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_OWNER_TOKEN" \
  -d '{
    "username": "user123",
    "password": "newuserpass",
    "confirm_password": "newuserpass"
  }'
```

**Response:**

```json
{
  "message": "Password changed successfully"
}
```

### Example 3: Admin Mengubah Password User

**Request:**

```bash
curl -X PUT http://localhost:8080/api/admin/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "username": "cashier1",
    "password": "cashierpass123",
    "confirm_password": "cashierpass123"
  }'
```

**Response:**

```json
{
  "message": "Password changed successfully"
}
```

### Example 4: Error - Password Tidak Match

**Request:**

```bash
curl -X PUT http://localhost:8080/api/admin/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "username": "user123",
    "password": "password1",
    "confirm_password": "password2"
  }'
```

**Response:**

```json
{
  "error": "password and confirm password do not match"
}
```

### Example 5: Error - Insufficient Permission

**Request:**

```bash
curl -X PUT http://localhost:8080/api/admin/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_USER_TOKEN" \
  -d '{
    "username": "admin1",
    "password": "newpass123",
    "confirm_password": "newpass123"
  }'
```

**Response:**

```json
{
  "error": "insufficient permission: can only change password for lower role users"
}
```

## Business Rules

1. **Role Hierarchy**: User yang melakukan perubahan harus memiliki role yang lebih tinggi dari target user
2. **Tenant Isolation**: Non-superadmin hanya dapat mengubah password user di tenant yang sama
3. **Branch Isolation**: Owner hanya dapat mengubah password user di branch yang sama
4. **Password Validation**:
   - Minimal 6 karakter
   - Password dan confirm password harus sama
5. **Username Uniqueness**: Username harus ada di database

## Security Notes

- API ini memerlukan autentikasi dengan Bearer token
- Token didapat dari login endpoint
- Validasi role dilakukan di service layer
- Password di-hash menggunakan bcrypt sebelum disimpan
- Owner tidak dapat mengubah password user di branch lain
- Admin tidak dapat mengubah password user di tenant lain
- User biasa tidak dapat mengubah password user lain

## Integration Testing

Untuk testing API ini di Postman:

1. Login sebagai superadmin/owner/admin untuk mendapatkan token
2. Gunakan token di header Authorization
3. Kirim request dengan username target dan password baru
4. Verifikasi response sukses atau error sesuai dengan permission
