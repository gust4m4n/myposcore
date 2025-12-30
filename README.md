# MyPOSCore - Backend API untuk Aplikasi POS

Backend API untuk aplikasi Point of Sale (POS) yang dibangun dengan Go dan PostgreSQL, mendukung multi-tenancy.

## Fitur

- ✅ Multi-tenancy support
- ✅ User registration dan login
- ✅ JWT authentication
- ✅ Password hashing dengan bcrypt
- ✅ RESTful API dengan Gin framework
- ✅ PostgreSQL database dengan GORM

## Teknologi

- **Golang** 1.21+
- **PostgreSQL** 
- **Gin** - Web framework
- **GORM** - ORM
- **JWT** - Authentication
- **bcrypt** - Password hashing

## Struktur Project

```
myposcore/
├── config/          # Konfigurasi aplikasi
├── database/        # Database connection & migration
├── dto/             # Data Transfer Objects
├── handlers/        # HTTP handlers
├── middleware/      # Middleware (auth, tenant)
├── models/          # Database models
├── routes/          # Route definitions
├── services/        # Business logic
├── utils/           # Utility functions
├── main.go          # Entry point
├── go.mod           # Go dependencies
└── .env.example     # Environment variables template
```

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE myposcore;
```

### 3. Environment Variables

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

Edit file `.env`:

```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=myposcore
JWT_SECRET=your-secret-key-change-this-in-production
```

### 4. Initialize Database

Jalankan script untuk membuat super tenant:

```bash
psql -U postgres -d myposcore -f init_tenant.sql
```

### 5. (Optional) Load Demo Data

Untuk testing, Anda bisa load demo tenant dengan 2 tipe bisnis:

```bash
psql -U postgres -d myposcore -f init_demo_tenants.sql
```

Demo data termasuk:
- 🍽️ **Tenant Restoran** dengan 12 produk makanan & minuman
- 👔 **Tenant Fashion Store** dengan 14 produk pakaian & aksesoris

Lihat [DEMO_TENANTS.md](DEMO_TENANTS.md) untuk detail lengkap credential dan data.

### 6. Jalankan Aplikasi

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /health
```

### Authentication

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "tenant_code": "TENANT001",
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

Response:
```json
{
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "tenant_id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "is_active": true
    }
  }
}
```

#### Login User
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "tenant_code": "TENANT001",
  "username": "john_doe",
  "password": "password123"
}
```

Response:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "tenant_id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "is_active": true
    }
  }
}
```

#### Get Profile (Protected)
```
GET /api/v1/profile
Authorization: Bearer <token>
```

#### Change Password (Protected)
```
PUT /api/v1/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

#### Admin Change Password (Protected)
```
PUT /api/v1/admin/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "user123",
  "password": "newpass123",
  "confirm_password": "newpass123"
}
```

**Role Requirements:**
- Superadmin: Dapat mengubah password semua user
- Owner: Dapat mengubah password admin dan user di branch-nya
- Admin: Dapat mengubah password user di tenant-nya

Lihat [ADMIN_CHANGE_PASSWORD_GUIDE.md](ADMIN_CHANGE_PASSWORD_GUIDE.md) untuk dokumentasi lengkap.

#### Admin Change PIN (Protected)
```
PUT /api/v1/admin/change-pin
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "user123",
  "pin": "123456",
  "confirm_pin": "123456"
}
```

**Role Requirements:**
- Superadmin: Dapat mengubah PIN semua user
- Owner: Dapat mengubah PIN admin dan user di branch-nya
- Admin: Dapat mengubah PIN user di tenant-nya

Lihat [ADMIN_CHANGE_PIN_GUIDE.md](ADMIN_CHANGE_PIN_GUIDE.md) untuk dokumentasi lengkap.

## Database Schema

### Tenants Table
```sql
- id (uint, primary key)
- name (string)
- code (string, unique)
- is_active (boolean)
- created_at (timestamp)
- updated_at (timestamp)
```

### Users Table
```sql
- id (uint, primary key)
- tenant_id (uint, foreign key)
- username (string)
- email (string)
- password (string, hashed)
- full_name (string)
- is_active (boolean)
- created_at (timestamp)
- updated_at (timestamp)
```

## Multi-Tenancy

Aplikasi menggunakan tenant isolation dengan:
- Setiap user terikat ke satu tenant dan branch
- Username unique secara global
- JWT token menyimpan tenant_id dan user_id
- Middleware memvalidasi akses berdasarkan tenant
- Semua resource (products, dll) terisolasi per tenant

## Testing

### Quick Test dengan Demo Data

Setelah menjalankan `init_demo_tenants.sql`, Anda bisa langsung testing dengan credential berikut:

**Restoran:**
```bash
# Login sebagai admin resto
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "resto01",
    "branch_code": "resto01-pusat",
    "username": "admin_resto",
    "password": "demo123"
  }'
```

**Fashion Store:**
```bash
# Login sebagai admin fashion
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "fashion01",
    "branch_code": "fashion01-plaza",
    "username": "admin_fashion",
    "password": "demo123"
  }'
```

Lihat [DEMO_TENANTS.md](DEMO_TENANTS.md) untuk daftar lengkap users, branches, dan products.

### Membuat Tenant Manual (via psql)

Jika ingin membuat tenant sendiri:

```sql
-- Insert tenant
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Tenant Demo', 'TENANT001', true, NOW(), NOW());

-- Insert branch
INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at)
VALUES (1, 'Branch Demo', 'BRANCH001', 'Jl. Demo No. 1', '021-12345678', true, NOW(), NOW());
```

### Test dengan curl

```bash
# Register User
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "TENANT001",
    "branch_code": "BRANCH001",
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "TENANT001",
    "branch_code": "BRANCH001",
    "username": "testuser",
    "password": "password123"
  }'

# Get Profile (ganti <TOKEN> dengan token dari login)
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <TOKEN>"

# Get Products (tenant isolated)
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <TOKEN>"
```

## Development

### Build
```bash
go build -o myposcore
```

### Run Binary
```bash
./myposcore
```

## Security Notes

- Password di-hash dengan bcrypt (cost 14)
- JWT token expire dalam 24 jam
- Ganti `JWT_SECRET` di production dengan nilai yang aman
- Gunakan HTTPS di production
- Implementasikan rate limiting untuk production

## License

MIT License
