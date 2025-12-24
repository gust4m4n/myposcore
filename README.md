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

### 4. Jalankan Aplikasi

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
- Setiap user terikat ke satu tenant
- Username dan email unique per tenant
- JWT token menyimpan tenant_id
- Middleware memvalidasi akses berdasarkan tenant

## Testing

### Membuat Tenant (Manual - via psql)

Sebelum register user, buat tenant terlebih dahulu:

```sql
INSERT INTO tenants (name, code, is_active, created_at, updated_at) 
VALUES ('Tenant Demo', 'TENANT001', true, NOW(), NOW());
```

### Test dengan curl

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_code": "TENANT001",
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
    "username": "testuser",
    "password": "password123"
  }'

# Get Profile (ganti <TOKEN> dengan token dari login)
curl -X GET http://localhost:8080/api/v1/profile \
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
