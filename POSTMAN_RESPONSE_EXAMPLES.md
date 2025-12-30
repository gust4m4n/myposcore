# Postman Response Examples - Success Cases

## Overview
Postman collection sudah dilengkapi dengan contoh response sukses untuk semua endpoint utama.

## Response Structure

### Standard Success Response Format

Semua API menggunakan format response yang konsisten:

```json
{
  "message": "Operation description",
  "data": { /* actual data */ }
}
```

Atau untuk simple operations:

```json
{
  "message": "Operation completed successfully"
}
```

## Authentication & Profile

### Register Success (200 OK)
```json
{
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "tenant_id": 1,
      "branch_id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "user",
      "is_active": true
    }
  }
}
```

### Login Success (200 OK)
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "tenant_id": 1,
      "branch_id": 1,
      "branch_name": "Main Branch",
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "user",
      "is_active": true
    }
  }
}
```

### Get Profile Success (200 OK)
```json
{
  "user": {
    "id": 25,
    "username": "tenantadmin",
    "email": "tenantadmin@resto.com",
    "full_name": "Tenant Admin Resto",
    "image": "/uploads/profiles/user_25_1735296000.jpg",
    "role": "tenantadmin",
    "is_active": true
  },
  "tenant": {
    "id": 17,
    "name": "Food Corner",
    "code": "resto01",
    "is_active": true
  },
  "branch": {
    "id": 25,
    "name": "Cabang Pusat",
    "code": "resto01-pusat",
    "address": "Jl. Sudirman No. 123, Jakarta",
    "phone": "021-12345678",
    "is_active": true
  }
}
```

### Change Password Success (200 OK)
```json
{
  "message": "Password changed successfully"
}
```

## Admin Operations ⭐

### Admin Change Password Success (200 OK)
```json
{
  "message": "Password changed successfully"
}
```

**Request:**
```json
{
  "username": "user123",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

### Admin Change PIN Success (200 OK)
```json
{
  "message": "PIN changed successfully"
}
```

**Request:**
```json
{
  "username": "user123",
  "pin": "123456",
  "confirm_pin": "123456"
}
```

## PIN Management

### Create PIN Success (200 OK)
```json
{
  "message": "PIN created successfully"
}
```

### Change PIN Success (200 OK)
```json
{
  "message": "PIN changed successfully"
}
```

### Check PIN Status (200 OK)
```json
{
  "has_pin": true
}
```

Or:
```json
{
  "has_pin": false
}
```

## FAQ Management

### Get All FAQ Success (200 OK)
```json
{
  "message": "FAQs retrieved successfully",
  "data": [
    {
      "id": 1,
      "question": "Bagaimana cara melakukan pembayaran?",
      "answer": "Anda dapat melakukan pembayaran melalui cash atau e-wallet",
      "category": "Payment",
      "order": 1,
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00"
    }
  ]
}
```

### Create FAQ Success (200 OK)
```json
{
  "message": "FAQ created successfully",
  "data": {
    "id": 25,
    "question": "Apakah MyPOS Core tersedia dalam bahasa Indonesia?",
    "answer": "Ya, MyPOS Core tersedia dalam bahasa Indonesia dan English.",
    "category": "General",
    "order": 3,
    "is_active": true,
    "created_at": "2025-12-26T10:00:00+07:00",
    "updated_at": "2025-12-26T10:00:00+07:00"
  }
}
```

### Update FAQ Success (200 OK)
```json
{
  "message": "FAQ updated successfully",
  "data": {
    "id": 1,
    "question": "Apakah MyPOS Core tersedia dalam bahasa Indonesia?",
    "answer": "Ya, MyPOS Core tersedia dalam bahasa Indonesia dan English. Anda dapat mengubah preferensi bahasa di pengaturan akun Anda.",
    "category": "General",
    "order": 3,
    "is_active": true,
    "created_at": "2025-12-25T10:00:00+07:00",
    "updated_at": "2025-12-26T11:30:00+07:00"
  }
}
```

### Delete FAQ Success (200 OK)
```json
{
  "message": "FAQ deleted successfully"
}
```

## Terms & Conditions Management

### Get Active TnC Success (200 OK)
```json
{
  "message": "Active terms and conditions retrieved successfully",
  "data": {
    "id": 1,
    "title": "Terms and Conditions v1.0",
    "content": "# Terms and Conditions\n\n## 1. Introduction...",
    "version": "1.0",
    "is_active": true,
    "created_at": "2025-12-25T10:00:00+07:00",
    "updated_at": "2025-12-25T10:00:00+07:00"
  }
}
```

### Create TnC Success (200 OK)
```json
{
  "message": "Terms and conditions created successfully",
  "data": {
    "id": 2,
    "title": "Privacy Policy",
    "content": "# Privacy Policy\n\n## Information We Collect...",
    "version": "1.0",
    "is_active": true,
    "created_at": "2025-12-26T10:00:00+07:00",
    "updated_at": "2025-12-26T10:00:00+07:00"
  }
}
```

### Update TnC Success (200 OK)
```json
{
  "message": "Terms and conditions updated successfully",
  "data": {
    "id": 1,
    "title": "Privacy Policy - Updated",
    "content": "# Privacy Policy\n\n## Information We Collect\n- Business information...",
    "version": "1.0.1",
    "is_active": true,
    "created_at": "2025-12-25T10:00:00+07:00",
    "updated_at": "2025-12-26T11:30:00+07:00"
  }
}
```

### Delete TnC Success (200 OK)
```json
{
  "message": "Terms and conditions deleted successfully"
}
```

## Product Management

### List Products Success (200 OK)
```json
{
  "data": [
    {
      "id": 50,
      "tenant_id": 18,
      "name": "Kemeja Batik Pria",
      "description": "Kemeja batik motif modern untuk pria",
      "category": "Pakaian Pria",
      "sku": "FASHION-001",
      "price": 250000,
      "stock": 30,
      "is_active": true,
      "created_at": "2025-12-24T19:48:10+07:00",
      "updated_at": "2025-12-24T19:48:10+07:00"
    }
  ]
}
```

### Create Product Success (200 OK)
```json
{
  "message": "Product created successfully",
  "data": {
    "id": 66,
    "tenant_id": 18,
    "name": "Tas Ransel Premium",
    "description": "Tas ransel dengan banyak kantong",
    "category": "Aksesoris",
    "sku": "FASHION-ACC-002",
    "price": 350000,
    "stock": 20,
    "is_active": true,
    "created_at": "2025-12-26T14:20:00+07:00",
    "updated_at": "2025-12-26T14:20:00+07:00"
  }
}
```

## User Management

### List Users Success (200 OK)
```json
{
  "data": [
    {
      "id": 1,
      "tenant_id": 1,
      "branch_id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-20T10:00:00+07:00",
      "updated_at": "2025-12-20T10:00:00+07:00"
    }
  ]
}
```

## Superadmin Management

### Dashboard Statistics Success (200 OK)
```json
{
  "message": "Dashboard data retrieved successfully",
  "data": {
    "total_tenants": 15,
    "active_tenants": 12,
    "total_branches": 45,
    "total_users": 234
  }
}
```

### List Tenants Success (200 OK)
```json
{
  "data": [
    {
      "id": 1,
      "name": "Restoran Sederhana",
      "code": "resto01",
      "is_active": true,
      "created_at": "2025-12-15T10:00:00+07:00",
      "updated_at": "2025-12-15T10:00:00+07:00"
    }
  ]
}
```

## Error Responses

### Common Error Formats

#### 400 Bad Request - Validation Error
```json
{
  "error": "password and confirm password do not match"
}
```

#### 401 Unauthorized
```json
{
  "error": "User not authenticated"
}
```

```json
{
  "error": "Invalid or expired token"
}
```

#### 403 Forbidden
```json
{
  "error": "insufficient permission: can only change password for lower role users"
}
```

#### 404 Not Found
```json
{
  "error": "target user not found"
}
```

```json
{
  "error": "FAQ not found"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to create FAQ"
}
```

## Response Headers

All responses include standard headers:

```
Content-Type: application/json
```

For authenticated endpoints:
```
Authorization: Bearer {token}
```

## Token Auto-Save

Postman collection sudah dikonfigurasi untuk auto-save token setelah login/register:

```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("auth_token", jsonData.data.token);
    pm.environment.set("user_id", jsonData.data.user.id);
    pm.environment.set("tenant_id", jsonData.data.user.tenant_id);
}
```

## Testing Tips

1. **Check Response Status Code**
   - 200 OK: Success
   - 400: Validation error
   - 401: Not authenticated
   - 403: No permission
   - 404: Not found

2. **Validate Response Structure**
   - Ensure `message` field exists
   - Check `data` structure matches expected format

3. **Test Error Scenarios**
   - Invalid data
   - Missing authentication
   - Insufficient permissions
   - Non-existent resources

4. **Use Response Examples**
   - Reference examples in Postman untuk expected responses
   - Compare actual vs expected results

## Notes

- Semua timestamp dalam format ISO 8601 dengan timezone (+07:00 untuk WIB)
- Boolean fields: `true` atau `false` (lowercase)
- ID fields: integer
- Dates: string dalam format ISO 8601
- Token berlaku 24 jam

---

**Collection Version:** 1.2  
**Last Updated:** December 30, 2025
