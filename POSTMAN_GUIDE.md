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

```bash
# Jalankan init script untuk data demo
psql -U postgres -d myposcore -f init_demo_tenants.sql
```

**Demo Users:**

- Superadmin: `admin@mypos.com` / `123456`
- Restaurant Admin: `admin.pusat@foodcorner99.com` / `123456`
- Fashion Admin: `admin.plaza@fashionhub.com` / `123456`

## 📝 Environment Variables

Collection ini menggunakan environment variables:

| Variable | Deskripsi | Default Value |
|----------|-----------|---------------|
| `base_url` | URL server backend | `http://localhost:8080` |
| `superadmin_email` | Email superadmin | `admin@mypos.com` |
| `superadmin_password` | Password superadmin | `123456` |
| `auth_token` | JWT token (auto-saved setelah login) | - |
| `user_id` | ID user (auto-saved setelah login) | - |
| `tenant_id` | ID tenant (auto-saved setelah login) | - |

## 🧪 Testing Flow

### Step 1: Health Check

- Endpoint: `GET /health`
- Untuk memastikan server berjalan dengan baik

### Step 2: Login User

- Endpoint: `POST /api/v1/auth/login`
- Body (JSON):

```json
{
  "email": "admin.pusat@foodcorner99.com",
  "password": "123456"
}
```

- **Token otomatis tersimpan** ke environment variable `auth_token`
- **Response** mencakup informasi user, tenant, dan branch

**Login dengan Role Berbeda:**

```json
// Superadmin
{
  "email": "admin@mypos.com",
  "password": "123456"
}

// Restaurant Branch Admin
{
  "email": "admin.pusat@foodcorner99.com",
  "password": "123456"
}

// Fashion Branch Admin
{
  "email": "admin.plaza@fashionhub.com",
  "password": "123456"
}
```

### Step 3: Get Profile (Protected Route)

- Endpoint: `GET /api/v1/profile`
- Authorization: Bearer Token (otomatis menggunakan `{{auth_token}}`)
- Mengembalikan informasi user yang sedang login

### Step 4: Testing Admin Operations ⭐

#### Admin Change Password

- Endpoint: `PUT /api/v1/admin/change-password`
- Untuk role tinggi mengubah password role rendah
- Request body:

```json
{
  "email": "user@example.com",
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
  "email": "user@example.com",
  "pin": "123456",
  "confirm_pin": "123456"
}
```

**Role Hierarchy:**

```
superadmin > owner > admin > user
```

### Step 5: Testing FAQ & TnC Management (Superadmin Only) ⭐

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

1. **Auto-save Token**: Setelah Login, token otomatis tersimpan ke environment
2. **Email-Based Login**: Login menggunakan email (unique di seluruh sistem)
3. **Response Examples**: Setiap endpoint memiliki contoh response untuk referensi
4. **Pre-request Scripts**: Login memiliki test scripts untuk auto-save token
5. **Role-Based Testing**: Login dengan user berbeda (superadmin/branchadmin/user) untuk test permission
6. **Collection Runner**: Gunakan untuk run multiple requests sekaligus
7. **No Registration**: User dibuat oleh admin melalui superadmin panel (tidak ada public registration)

## 🔄 Testing Multiple Tenants & Branches

**Demo Data** sudah menyediakan 2 tenant dengan multiple branches:

**Food Corner 99 (Restaurant)**

- Branch Pusat: `admin.pusat@foodcorner99.com` / `123456`
- Branch Menteng: `admin.menteng@foodcorner99.com` / `123456`

**Fashion Hub (Retail)**

- Branch Plaza: `admin.plaza@fashionhub.com` / `123456`
- Branch Grand Mall: `admin.grandmall@fashionhub.com` / `123456`

Login dengan email berbeda untuk test multi-tenant isolation.

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

### User Not Found

- Pastikan user sudah dibuat oleh admin
- Cek apakah user dalam status `is_active = true`
- Email harus valid dan terdaftar di sistem

### Email Already Exists

- Email harus unique di seluruh sistem (global)

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
- [CODE_FIELD_REMOVAL_SUMMARY.md](CODE_FIELD_REMOVAL_SUMMARY.md) - Summary perubahan code fields

## 📖 Next Steps

Setelah berhasil testing authentication:

1. Tambahkan endpoint untuk manajemen produk
2. Tambahkan endpoint untuk kategori
3. Tambahkan endpoint untuk transaksi POS
4. Tambahkan endpoint untuk laporan

## 🔗 Links

- [Dokumentasi Postman](https://learning.postman.com/docs/getting-started/introduction/)
- [README.md](README.md) - Dokumentasi lengkap aplikasi
