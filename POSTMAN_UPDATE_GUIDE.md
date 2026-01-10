# Postman Collection - Update Guide

## Overview
Postman collection untuk MyPOSCore API sudah lengkap dengan semua endpoint termasuk API baru:
- ✅ Admin Change Password
- ✅ Admin Change PIN  
- ✅ FAQ CRUD Operations
- ✅ Terms & Conditions CRUD Operations

## Import ke Postman

### 1. Import Collection
1. Buka Postman
2. Klik **Import** di kiri atas
3. Pilih file `MyPOSCore.postman_collection.json`
4. Collection akan muncul di sidebar

### 2. Import Environment
1. Klik icon ⚙️ (Settings) di kiri atas
2. Pilih **Environments**
3. Klik **Import**
4. Pilih file `MyPOSCore.postman_environment.json`
5. Aktifkan environment "MyPOSCore Local"

## Environment Variables

Collection menggunakan variable berikut (sudah dikonfigurasi):

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `base_url` | `http://localhost:8080` | Base URL server |
| `tenant_code` | `TENANT001` | Kode tenant untuk testing |
| `branch_code` | `BRANCH001` | Kode branch untuk testing |
| `superadmin_tenant` | `supertenant` | Tenant code superadmin |
| `superadmin_branch` | `superbranch` | Branch code superadmin |
| `superadmin_username` | `superadmin` | Username superadmin |
| `superadmin_password` | `123456` | Password superadmin |
| `auth_token` | (auto) | JWT token (otomatis di-set setelah login) |
| `user_id` | (auto) | User ID (otomatis di-set setelah login) |
| `tenant_id` | (auto) | Tenant ID (otomatis di-set setelah login) |

## Struktur Collection

### 1. Health Check
- **GET** `/health` - Cek status server

### 2. Dev Tools
- **GET** `/dev/tenants` - List semua tenant (untuk testing)
- **GET** `/dev/tenants/:tenant_id/branches` - List branches per tenant

### 3. Authentication
- **POST** `/api/auth/register` - Register user baru
- **POST** `/api/auth/login` - Login user
- **POST** `/api/auth/login` (Superadmin) - Login sebagai superadmin

### 4. Profile Management
- **GET** `/api/profile` - Get user profile
- **PUT** `/api/profile` - Update profile
- **PUT** `/api/change-password` - Change own password
- **POST** `/api/profile/photo` - Upload profile photo
- **DELETE** `/api/profile/photo` - Delete profile photo

### 5. Admin Operations ⭐ NEW
- **PUT** `/api/admin/change-password` - Change password user lain (role tinggi)
  - Contoh: Superadmin → Owner → Admin → User
- **PUT** `/api/admin/change-pin` - Change PIN user lain (role tinggi)
  - Contoh: Superadmin → Owner → Admin → User

### 6. PIN Management
- **POST** `/api/pin/create` - Create PIN
- **PUT** `/api/pin/change` - Change own PIN
- **GET** `/api/pin/check` - Check PIN status
- **PUT** `/api/admin/change-pin` ⭐ - Admin change PIN (NEW)

### 7. Category Management
- **GET** `/api/categories` - List categories
- **GET** `/api/categories/:id` - Get category by ID
- **POST** `/api/categories` - Create category
- **PUT** `/api/categories/:id` - Update category
- **DELETE** `/api/categories/:id` - Delete category

### 8. Product Management
- **GET** `/api/products` - List products
- **GET** `/api/products/:id` - Get product by ID
- **POST** `/api/products` - Create product
- **PUT** `/api/products/:id` - Update product
- **DELETE** `/api/products/:id` - Delete product
- **POST** `/api/products/:id/photo` - Upload product photo
- **DELETE** `/api/products/:id/photo` - Delete product photo

### 9. Order Management
- **GET** `/api/orders` - List orders
- **GET** `/api/orders/:id` - Get order by ID
- **POST** `/api/orders` - Create order

### 10. Payment Management
- **GET** `/api/payments` - List payments
- **GET** `/api/payments/:id` - Get payment by ID
- **POST** `/api/payments` - Create payment

### 11. User Management
- **GET** `/api/users` - List users
- **GET** `/api/users/:id` - Get user by ID
- **POST** `/api/users` - Create user
- **PUT** `/api/users/:id` - Update user
- **DELETE** `/api/users/:id` - Delete user

### 12. FAQ Management
- **GET** `/api/faq` - Get all FAQ (Public)
- **GET** `/api/faq/:id` - Get FAQ by ID (Public)
- **POST** `/api/superadmin/faq` - Create FAQ (Superadmin only)
- **PUT** `/api/superadmin/faq/:id` - Update FAQ (Superadmin only)
- **DELETE** `/api/superadmin/faq/:id` - Delete FAQ (Superadmin only)

### 13. Terms & Conditions Management
- **GET** `/api/tnc` - Get all TnC (Public)
- **GET** `/api/tnc/:id` - Get TnC by ID (Public)
- **GET** `/api/tnc/active` - Get active TnC (Public)
- **POST** `/api/superadmin/tnc` - Create TnC (Superadmin only)
- **PUT** `/api/superadmin/tnc/:id` - Update TnC (Superadmin only)
- **DELETE** `/api/superadmin/tnc/:id` - Delete TnC (Superadmin only)

### 14. Superadmin Management
- **GET** `/api/superadmin/dashboard` - Dashboard statistics
- **GET** `/api/superadmin/tenants` - List all tenants
- **POST** `/api/superadmin/tenants` - Create tenant
- **GET** `/api/superadmin/tenants/:tenant_id/branches` - List branches
- **GET** `/api/superadmin/branches/:branch_id/users` - List users

## Cara Penggunaan

### Quick Start - Testing Flow

1. **Start Server**
   ```bash
   cd /path/to/myposcore
   ./myposcore
   ```

2. **Health Check**
   - Jalankan request "Health Check"
   - Pastikan response: `{"status": "OK"}`

3. **Login as Superadmin**
   - Folder: Authentication
   - Request: "Login (Superadmin)"
   - Token akan otomatis tersimpan di environment variable

4. **Test Admin Operations**
   - Login sebagai user dengan role lebih tinggi
   - Gunakan request "Admin Change Password" atau "Admin Change PIN"
   - Masukkan username target user

### Testing Admin Change Password

**Scenario 1: Superadmin changes Owner password**
```json
{
  "username": "owner1",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Scenario 2: Owner changes User password**
```json
{
  "username": "user123",
  "password": "userpass123",
  "confirm_password": "userpass123"
}
```

**Scenario 3: Admin changes User password (same tenant)**
```json
{
  "username": "cashier1",
  "password": "cashierpass",
  "confirm_password": "cashierpass"
}
```

### Testing Admin Change PIN

**Scenario 1: Superadmin changes any user PIN**
```json
{
  "username": "owner1",
  "pin": "123456",
  "confirm_pin": "123456"
}
```

**Scenario 2: Owner changes User PIN (same branch)**
```json
{
  "username": "user123",
  "pin": "999888",
  "confirm_pin": "999888"
}
```

### Testing FAQ Operations

**Create FAQ (Superadmin):**
```json
{
  "question": "Bagaimana cara reset password?",
  "answer": "Hubungi admin untuk reset password",
  "category": "Account",
  "order": 1
}
```

**Update FAQ (Superadmin):**
```json
{
  "answer": "Jawaban yang telah diperbarui",
  "is_active": true
}
```

## Auto-Save Token

Collection sudah dikonfigurasi dengan **Test Scripts** yang otomatis menyimpan token setelah login:

```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("auth_token", jsonData.data.token);
    pm.environment.set("user_id", jsonData.data.user.id);
    pm.environment.set("tenant_id", jsonData.data.user.tenant_id);
}
```

Setiap request yang memerlukan autentikasi akan otomatis menggunakan `{{auth_token}}`.

## Response Examples

Setiap request sudah dilengkapi dengan contoh response untuk:
- ✅ Success response
- ✅ Error response (validation, permission, not found, dll)

## Tips & Best Practices

1. **Gunakan Environment Variables**
   - Mudah switch antara local, staging, production
   - Update `base_url` saja untuk ganti environment

2. **Save Requests**
   - Setelah mengedit request, simpan perubahan (Ctrl/Cmd + S)

3. **Run Collection**
   - Gunakan Collection Runner untuk test semua endpoint sekaligus
   - Pastikan server running sebelum run collection

4. **Test Different Roles**
   - Login sebagai superadmin, owner, admin, user
   - Test permission untuk setiap role

5. **Check Response**
   - Perhatikan status code (200, 400, 401, 404, dll)
   - Validasi structure response sesuai dokumentasi

## Troubleshooting

### Token Expired
- Login ulang untuk mendapatkan token baru
- Token valid 24 jam

### 401 Unauthorized
- Pastikan sudah login dan token tersimpan
- Check environment variable `auth_token`

### 403 Forbidden
- User tidak memiliki permission untuk endpoint tersebut
- Login dengan role yang sesuai (superadmin/owner/admin)

### 404 Not Found
- Endpoint salah atau resource tidak ditemukan
- Check URL dan parameter ID

### Connection Refused
- Server belum running
- Check `base_url` di environment

## Export & Share

### Export Collection
1. Right-click pada collection
2. Pilih **Export**
3. Pilih format Collection v2.1
4. Save file

### Export Environment
1. Klik icon ⚙️ (Settings)
2. Pilih tab **Environments**
3. Klik ... di environment
4. Pilih **Export**

## Updates Log

### Version 1.2 (Latest)
- ✅ Added Admin Change Password API
- ✅ Added Admin Change PIN API
- ✅ Updated FAQ CRUD documentation
- ✅ Updated TnC CRUD documentation
- ✅ Added test scripts for auto-save token
- ✅ Added multiple response examples

### Version 1.1
- Added Superadmin management endpoints
- Added FAQ & TnC public endpoints
- Added PIN management

### Version 1.0
- Initial collection with core API endpoints

---

**Need Help?** Check dokumentasi lengkap di:
- [ADMIN_CHANGE_PASSWORD_GUIDE.md](ADMIN_CHANGE_PASSWORD_GUIDE.md)
- [ADMIN_CHANGE_PIN_GUIDE.md](ADMIN_CHANGE_PIN_GUIDE.md)
- [FAQ_TNC_GUIDE.md](FAQ_TNC_GUIDE.md)
- [README.md](README.md)
