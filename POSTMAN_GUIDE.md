# Panduan Menggunakan Postman Collection untuk MyPOSCore API

## 📥 Import Collection dan Environment

### 1. Import Collection
1. Buka Postman
2. Klik **Import** di pojok kiri atas
3. Pilih file `MyPOSCore.postman_collection.json`
4. Collection akan muncul di sidebar

### 2. Import Environment
1. Klik icon ⚙️ (Settings) di pojok kanan atas
2. Pilih tab **Environments**
3. Klik **Import**
4. Pilih file `MyPOSCore.postman_environment.json`
5. Pilih environment "MyPOSCore Local" dari dropdown di kanan atas

## 🆕 What's New - Latest Updates

### Version 1.2 (December 2025)
- ✅ **Admin Change Password API** - Role tinggi dapat mengubah password role rendah
- ✅ **Admin Change PIN API** - Role tinggi dapat mengubah PIN role rendah
- ✅ **FAQ CRUD Operations** - Create, Read, Update, Delete FAQ (Superadmin)
- ✅ **TnC CRUD Operations** - Create, Read, Update, Delete Terms & Conditions (Superadmin)
- ✅ Auto-save token test scripts
- ✅ Multiple response examples untuk setiap endpoint

## 🚀 Persiapan Testing

### 1. Jalankan Server
```bash
cd /Users/gustaman/Desktop/GUSTAMAN7/myposcore
go run main.go
```

### 2. Pastikan Database Sudah Setup
```sql
-- Buat tenant demo terlebih dahulu
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Demo Tenant', 'TENANT001', true, NOW(), NOW());
```

## 📝 Environment Variables

Collection ini menggunakan environment variables:

| Variable | Deskripsi | Default Value |
|----------|-----------|---------------|
| `base_url` | URL server backend | `http://localhost:8080` |
| `tenant_code` | Kode tenant untuk testing | `TENANT001` |
| `branch_code` | Kode branch untuk testing | `BRANCH001` |
| `superadmin_tenant` | Tenant code superadmin | `supertenant` |
| `superadmin_branch` | Branch code superadmin | `superbranch` |
| `superadmin_username` | Username superadmin | `superadmin` |
| `superadmin_password` | Password superadmin | `123456` |
| `auth_token` | JWT token (auto-saved setelah login) | - |
| `user_id` | ID user (auto-saved setelah login) | - |
| `tenant_id` | ID tenant (auto-saved setelah login) | - |

## 🧪 Testing Flow

### Step 1: Health Check
- Endpoint: `GET /health`
- Untuk memastikan server berjalan dengan baik

### Step 2: Register User
- Endpoint: `POST /api/v1/auth/register`
- Body (JSON):
```json
{
  "tenant_code": "TENANT001",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```
- **Token otomatis tersimpan** ke environment variable `auth_token`

### Step 3: Login User
- Endpoint: `POST /api/v1/auth/login`
- Body (JSON):
```json
{
  "tenant_code": "TENANT001",
  "username": "johndoe",
  "password": "password123"
}
```
- **Token otomatis tersimpan** ke environment variable `auth_token`

### Step 4: Get Profile (Protected Route)
- Endpoint: `GET /api/v1/profile`
- Authorization: Bearer Token (otomatis menggunakan `{{auth_token}}`)
- Mengembalikan informasi user yang sedang login

### Step 5: Testing Admin Operations (NEW) ⭐

#### Admin Change Password
- Endpoint: `PUT /api/v1/admin/change-password`
- Untuk role tinggi mengubah password role rendah
- Request body:
```json
{
  "username": "user123",
  "password": "newpass123",
  "confirm_password": "newpass123"
}
```

#### Admin Change PIN
- Endpoint: `PUT /api/v1/admin/change-pin`
- Untuk role tinggi mengubah PIN role rendah
- Request body:
```json
{
  "username": "user123",
  "pin": "123456",
  "confirm_pin": "123456"
}
```

**Role Hierarchy:**
```
superadmin > owner > admin > user
```

### Step 6: Testing FAQ & TnC Management (Superadmin Only) ⭐

#### Create FAQ
- Endpoint: `POST /api/v1/superadmin/faq`
- Request body:
```json
{
  "question": "Bagaimana cara reset password?",
  "answer": "Hubungi admin",
  "category": "Account",
  "order": 1
}
```

#### Update TnC
- Endpoint: `PUT /api/v1/superadmin/tnc/:id`
- Request body:
```json
{
  "title": "Terms and Conditions v2",
  "content": "Updated content...",
  "version": "2.0",
  "is_active": true
}
```

## 🔐 Authorization

Untuk endpoint yang memerlukan autentikasi:
1. Token JWT otomatis ditambahkan ke header `Authorization: Bearer {{auth_token}}`
2. Token disimpan otomatis setelah Register atau Login
3. Token berlaku selama 24 jam

## 📌 Tips

1. **Auto-save Token**: Setelah Register atau Login, token otomatis tersimpan ke environment
2. **Multiple Tenants**: Ubah nilai `tenant_code` di environment untuk testing tenant berbeda
3. **Response Examples**: Setiap endpoint memiliki contoh response untuk referensi
4. **Pre-request Scripts**: Login dan Register memiliki test scripts untuk auto-save token
5. **Role-Based Testing**: Login dengan user berbeda (superadmin/owner/admin/user) untuk test permission
6. **Collection Runner**: Gunakan untuk run multiple requests sekaligus

## 🔄 Testing Multiple Tenants

1. Buat tenant baru di database:
```sql
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Tenant Dua', 'TENANT002', true, NOW(), NOW());
```

2. Ubah environment variable `tenant_code` menjadi `TENANT002`
3. Register user baru dengan tenant code tersebut
4. User dari TENANT001 dan TENANT002 sepenuhnya terisolasi

## 📊 Status Codes

| Code | Deskripsi |
|------|-----------|
| 200 | OK - Request berhasil |
| 400 | Bad Request - Validasi error atau data tidak valid |
| 401 | Unauthorized - Token invalid atau expired |
| 403 | Forbidden - User tidak memiliki permission |
| 404 | Not Found - Resource tidak ditemukan |
| 500 | Internal Server Error |

## 🐛 Troubleshooting

### Token Invalid/Expired
- Lakukan login ulang untuk mendapatkan token baru
- Token berlaku 24 jam sejak dibuat

### 403 Forbidden
- User tidak memiliki permission untuk endpoint tersebut
- Login dengan role yang sesuai (superadmin untuk FAQ/TnC management)

### Tenant Not Found
- Pastikan tenant dengan kode yang digunakan sudah ada di database
- Cek apakah tenant dalam status `is_active = true`

### Username/Email Already Exists
- Username dan email harus unik per tenant

### Insufficient Permission (Admin Change Password/PIN)
- Role user yang login harus lebih tinggi dari target user
- Check role hierarchy: superadmin > owner > admin > user

## 📚 Dokumentasi Lengkap

Untuk informasi detail, lihat:
- [ADMIN_CHANGE_PASSWORD_GUIDE.md](ADMIN_CHANGE_PASSWORD_GUIDE.md) - Admin change password
- [ADMIN_CHANGE_PIN_GUIDE.md](ADMIN_CHANGE_PIN_GUIDE.md) - Admin change PIN
- [FAQ_TNC_GUIDE.md](FAQ_TNC_GUIDE.md) - FAQ & TnC management
- [POSTMAN_UPDATE_GUIDE.md](POSTMAN_UPDATE_GUIDE.md) - Update guide lengkap
- [README.md](README.md) - Overview project
- Gunakan username/email berbeda atau tenant code berbeda

## 📖 Next Steps

Setelah berhasil testing authentication:
1. Tambahkan endpoint untuk manajemen produk
2. Tambahkan endpoint untuk kategori
3. Tambahkan endpoint untuk transaksi POS
4. Tambahkan endpoint untuk laporan

## 🔗 Links

- [Dokumentasi Postman](https://learning.postman.com/docs/getting-started/introduction/)
- [README.md](README.md) - Dokumentasi lengkap aplikasi
