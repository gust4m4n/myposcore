# Admin Change PIN API Guide

API untuk mengubah PIN user lain (role tinggi ke role rendah).

## Role Hierarchy

```
superadmin > owner > admin > user
```

- **superadmin**: Dapat mengubah PIN semua user
- **owner**: Dapat mengubah PIN admin dan user di branch-nya
- **admin**: Dapat mengubah PIN user di tenant-nya
- **user**: Tidak dapat mengubah PIN user lain

## Endpoint

### Change PIN by Admin

Mengubah PIN user lain berdasarkan username.

**URL**: `/api/v1/admin/change-pin`

**Method**: `PUT`

**Auth Required**: Yes (Bearer Token)

**Permissions**: 
- Superadmin dapat mengubah PIN semua user
- Owner dapat mengubah PIN admin dan user di branch yang sama
- Admin dapat mengubah PIN user di tenant yang sama

### Request Body

```json
{
  "username": "user123",
  "pin": "123456",
  "confirm_pin": "123456"
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | Yes | Username dari user yang akan diubah PIN-nya |
| pin | string | Yes | PIN baru (6 digit angka) |
| confirm_pin | string | Yes | Konfirmasi PIN baru (harus sama dengan pin) |

### Success Response

**Code**: `200 OK`

**Content Example**:

```json
{
  "message": "PIN changed successfully"
}
```

### Error Responses

**Code**: `400 BAD REQUEST`

**Content Examples**:

```json
{
  "error": "PIN and confirm PIN do not match"
}
```

```json
{
  "error": "PIN must be exactly 6 digits"
}
```

```json
{
  "error": "target user not found"
}
```

```json
{
  "error": "insufficient permission: can only change PIN for lower role users"
}
```

```json
{
  "error": "cannot change PIN for users from different tenant"
}
```

```json
{
  "error": "owner can only change PIN for users in their branch"
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

### Example 1: Superadmin Mengubah PIN Owner

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/admin/change-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_SUPERADMIN_TOKEN" \
  -d '{
    "username": "owner1",
    "pin": "654321",
    "confirm_pin": "654321"
  }'
```

**Response:**

```json
{
  "message": "PIN changed successfully"
}
```

### Example 2: Owner Mengubah PIN User

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/admin/change-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_OWNER_TOKEN" \
  -d '{
    "username": "user123",
    "pin": "111111",
    "confirm_pin": "111111"
  }'
```

**Response:**

```json
{
  "message": "PIN changed successfully"
}
```

### Example 3: Admin Mengubah PIN User

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/admin/change-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "username": "cashier1",
    "pin": "999888",
    "confirm_pin": "999888"
  }'
```

**Response:**

```json
{
  "message": "PIN changed successfully"
}
```

### Example 4: Error - PIN Tidak Match

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/admin/change-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "username": "user123",
    "pin": "123456",
    "confirm_pin": "654321"
  }'
```

**Response:**

```json
{
  "error": "PIN and confirm PIN do not match"
}
```

### Example 5: Error - Insufficient Permission

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/admin/change-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_USER_TOKEN" \
  -d '{
    "username": "admin1",
    "pin": "888999",
    "confirm_pin": "888999"
  }'
```

**Response:**

```json
{
  "error": "insufficient permission: can only change PIN for lower role users"
}
```

### Example 6: Error - Invalid PIN Format

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/admin/change-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "username": "user123",
    "pin": "12345",
    "confirm_pin": "12345"
  }'
```

**Response:**

```json
{
  "error": "Key: 'AdminChangePINRequest.PIN' Error:Field validation for 'PIN' failed on the 'len' tag"
}
```

## Business Rules

1. **Role Hierarchy**: User yang melakukan perubahan harus memiliki role yang lebih tinggi dari target user
2. **Tenant Isolation**: Non-superadmin hanya dapat mengubah PIN user di tenant yang sama
3. **Branch Isolation**: Owner hanya dapat mengubah PIN user di branch yang sama
4. **PIN Validation**:
   - Harus 6 digit angka
   - PIN dan confirm PIN harus sama
5. **Username Uniqueness**: Username harus ada di database
6. **PIN Storage**: PIN di-hash menggunakan bcrypt sebelum disimpan (sama seperti password)

## Security Notes

- API ini memerlukan autentikasi dengan Bearer token
- Token didapat dari login endpoint
- Validasi role dilakukan di service layer
- PIN di-hash menggunakan bcrypt sebelum disimpan (untuk keamanan maksimal)
- Owner tidak dapat mengubah PIN user di branch lain
- Admin tidak dapat mengubah PIN user di tenant lain
- User biasa tidak dapat mengubah PIN user lain
- PIN harus numerik dan tepat 6 digit

## Comparison with Self PIN Change

| Aspek | Admin Change PIN | Self Change PIN |
|-------|------------------|-----------------|
| Endpoint | `/api/v1/admin/change-pin` | `/api/v1/pin/change` |
| Parameter | username, pin, confirm_pin | old_pin, new_pin, confirm_pin |
| Old PIN Required | No | Yes |
| Role Check | Yes (hierarchy) | No |
| Tenant Check | Yes (for non-superadmin) | No |
| Target | Other users | Self only |

## Use Cases

1. **Reset PIN Karyawan**: Owner atau admin dapat mereset PIN karyawan yang lupa
2. **Onboarding Baru**: Admin dapat set PIN awal untuk user baru
3. **Security Issue**: Superadmin dapat mengubah PIN user yang terkompromi
4. **Employee Offboarding**: Manager dapat disable atau change PIN sebelum employee exit

## Integration Testing

Untuk testing API ini di Postman:

1. Login sebagai superadmin/owner/admin untuk mendapatkan token
2. Gunakan token di header Authorization
3. Kirim request dengan username target dan PIN baru
4. Verifikasi response sukses atau error sesuai dengan permission
5. Test dengan berbagai skenario: success, mismatch PIN, insufficient permission, etc.

## Related APIs

- [POST /api/v1/pin/create](PIN_API_GUIDE.md#create-pin) - Create PIN (self)
- [PUT /api/v1/pin/change](PIN_API_GUIDE.md#change-pin) - Change PIN (self)
- [GET /api/v1/pin/check](PIN_API_GUIDE.md#check-pin) - Check PIN (self)
- [PUT /api/v1/admin/change-password](ADMIN_CHANGE_PASSWORD_GUIDE.md) - Admin Change Password
