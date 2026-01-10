# FAQ & Terms & Conditions API Guide

## Overview

MyPOS Core menyediakan endpoints untuk mengelola FAQ (Frequently Asked Questions) dan Terms & Conditions (TnC). API ini dibagi menjadi:
- **Public endpoints**: Akses tanpa authentication untuk membaca FAQ dan TnC
- **Superadmin endpoints**: Memerlukan token superadmin untuk CRUD operations

---

## FAQ Endpoints

### Public Endpoints (No Authentication Required)

#### 1. Get All FAQ
```
GET /api/faq
```

**Query Parameters:**
- `category` (optional): Filter by category
- `active_only` (optional): Show only active FAQs (true/false)

**Example Request:**
```bash
curl http://localhost:8080/api/faq?category=General&active_only=true
```

**Example Response:**
```json
{
  "data": [
    {
      "id": 1,
      "question": "Apa itu MyPOS Core?",
      "answer": "MyPOS Core adalah sistem Point of Sale (POS) berbasis cloud...",
      "category": "General",
      "order": 1,
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00"
    }
  ]
}
```

#### 2. Get FAQ by ID
```
GET /api/faq/:id
```

**Example Request:**
```bash
curl http://localhost:8080/api/faq/1
```

**Example Response:**
```json
{
  "data": {
    "id": 1,
    "question": "Apa itu MyPOS Core?",
    "answer": "MyPOS Core adalah sistem Point of Sale...",
    "category": "General",
    "order": 1,
    "is_active": true,
    "created_at": "2025-12-25T10:00:00+07:00",
    "updated_at": "2025-12-25T10:00:00+07:00"
  }
}
```

### Superadmin Endpoints (Authentication Required)

#### 3. Create FAQ
```
POST /api/superadmin/faq
```

**Headers:**
```
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "question": "Bagaimana cara membuat akun?",
  "answer": "Klik tombol Register di halaman login...",
  "category": "Account",
  "order": 1
}
```

**Validation:**
- `question`: required, min 5 characters, max 500 characters
- `answer`: required
- `category`: optional, max 100 characters
- `order`: optional, default 0

**Example Response:**
```json
{
  "message": "FAQ created successfully",
  "data": {
    "id": 25,
    "question": "Bagaimana cara membuat akun?",
    "answer": "Klik tombol Register di halaman login...",
    "category": "Account",
    "order": 1,
    "is_active": true,
    "created_at": "2025-12-26T10:00:00+07:00",
    "updated_at": "2025-12-26T10:00:00+07:00"
  }
}
```

#### 4. Update FAQ
```
PUT /api/superadmin/faq/:id
```

**Headers:**
```
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Request Body (all fields optional):**
```json
{
  "question": "Bagaimana cara membuat akun di MyPOS Core?",
  "answer": "Updated answer...",
  "category": "Account",
  "order": 2,
  "is_active": false
}
```

**Example Response:**
```json
{
  "message": "FAQ updated successfully",
  "data": {
    "id": 25,
    "question": "Bagaimana cara membuat akun di MyPOS Core?",
    "answer": "Updated answer...",
    "category": "Account",
    "order": 2,
    "is_active": false,
    "created_at": "2025-12-26T10:00:00+07:00",
    "updated_at": "2025-12-26T11:30:00+07:00"
  }
}
```

#### 5. Delete FAQ
```
DELETE /api/superadmin/faq/:id
```

**Headers:**
```
Authorization: Bearer <superadmin_token>
```

**Example Response:**
```json
{
  "message": "FAQ deleted successfully"
}
```

---

## Terms & Conditions Endpoints

### Public Endpoints (No Authentication Required)

#### 1. Get All TnC
```
GET /api/tnc
```

**Example Request:**
```bash
curl http://localhost:8080/api/tnc
```

**Example Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Terms and Conditions - MyPOS Core System",
      "content": "# Terms and Conditions for MyPOS Core\n\n## 1. Introduction\nWelcome to MyPOS Core...",
      "version": "1.0.0",
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00"
    },
    {
      "id": 2,
      "title": "Privacy Policy",
      "content": "# Privacy Policy\n\n## Information We Collect...",
      "version": "1.0.0",
      "is_active": true,
      "created_at": "2025-12-25T10:00:00+07:00",
      "updated_at": "2025-12-25T10:00:00+07:00"
    }
  ]
}
```

#### 2. Get Active TnC
```
GET /api/tnc/active
```

**Example Request:**
```bash
curl http://localhost:8080/api/tnc/active
```

**Example Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Terms and Conditions - MyPOS Core System",
    "content": "# Terms and Conditions for MyPOS Core\n\n## 1. Introduction...",
    "version": "1.0.0",
    "is_active": true,
    "created_at": "2025-12-25T10:00:00+07:00",
    "updated_at": "2025-12-25T10:00:00+07:00"
  }
}
```

#### 3. Get TnC by ID
```
GET /api/tnc/:id
```

**Example Request:**
```bash
curl http://localhost:8080/api/tnc/1
```

**Example Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Terms and Conditions - MyPOS Core System",
    "content": "# Terms and Conditions for MyPOS Core...",
    "version": "1.0.0",
    "is_active": true,
    "created_at": "2025-12-25T10:00:00+07:00",
    "updated_at": "2025-12-25T10:00:00+07:00"
  }
}
```

### Superadmin Endpoints (Authentication Required)

#### 4. Create TnC
```
POST /api/superadmin/tnc
```

**Headers:**
```
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Privacy Policy",
  "content": "# Privacy Policy\n\n## Information We Collect\n- Business information...",
  "version": "1.0.0"
}
```

**Validation:**
- `title`: required, max 255 characters
- `content`: required (supports Markdown)
- `version`: required, max 20 characters

**Example Response:**
```json
{
  "message": "TnC created successfully",
  "data": {
    "id": 2,
    "title": "Privacy Policy",
    "content": "# Privacy Policy\n\n## Information We Collect...",
    "version": "1.0.0",
    "is_active": true,
    "created_at": "2025-12-26T10:00:00+07:00",
    "updated_at": "2025-12-26T10:00:00+07:00"
  }
}
```

#### 5. Update TnC
```
PUT /api/superadmin/tnc/:id
```

**Headers:**
```
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Request Body (all fields optional):**
```json
{
  "title": "Privacy Policy - Updated",
  "content": "# Privacy Policy\n\nUpdated content...",
  "version": "1.0.1",
  "is_active": true
}
```

**Example Response:**
```json
{
  "message": "TnC updated successfully",
  "data": {
    "id": 2,
    "title": "Privacy Policy - Updated",
    "content": "# Privacy Policy\n\nUpdated content...",
    "version": "1.0.1",
    "is_active": true,
    "created_at": "2025-12-26T10:00:00+07:00",
    "updated_at": "2025-12-26T11:30:00+07:00"
  }
}
```

#### 6. Delete TnC
```
DELETE /api/superadmin/tnc/:id
```

**Headers:**
```
Authorization: Bearer <superadmin_token>
```

**Example Response:**
```json
{
  "message": "TnC deleted successfully"
}
```

---

## FAQ Categories

Berikut adalah kategori FAQ yang tersedia dalam dummy data:

- **General**: Informasi umum tentang MyPOS Core
- **Account**: Terkait akun dan login
- **Products**: Manajemen produk dan inventory
- **Orders**: Pembuatan dan pengelolaan order
- **Payments**: Metode pembayaran
- **Reports**: Laporan dan analytics
- **Technical**: Teknis dan sistem
- **Pricing**: Harga dan subscription
- **Support**: Dukungan customer

---

## Common Use Cases

### 1. Display FAQ on Website
```bash
# Get all active FAQs grouped by category
curl "http://localhost:8080/api/faq?active_only=true"
```

### 2. Show FAQ by Category
```bash
# Get only Technical FAQs
curl "http://localhost:8080/api/faq?category=Technical&active_only=true"
```

### 3. Display Terms & Conditions on Registration
```bash
# Get active TnC for user to accept
curl "http://localhost:8080/api/tnc/active"
```

### 4. Admin: Add New FAQ
```bash
curl -X POST http://localhost:8080/api/superadmin/faq \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Apakah ada aplikasi mobile?",
    "answer": "Ya, aplikasi mobile MyPOS tersedia di Play Store dan App Store",
    "category": "General",
    "order": 10
  }'
```

### 5. Admin: Update TnC Version
```bash
curl -X PUT http://localhost:8080/api/superadmin/tnc/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "2.0.0",
    "content": "Updated Terms and Conditions..."
  }'
```

---

## Error Responses

### 404 Not Found
```json
{
  "error": "FAQ not found"
}
```

### 401 Unauthorized (for superadmin endpoints)
```json
{
  "error": "Unauthorized"
}
```

### 400 Bad Request
```json
{
  "error": "Validation error: question must be at least 5 characters"
}
```

---

## Loading Dummy Data

Untuk load dummy content FAQ dan TnC ke database:

```bash
# Load TnC dummy data
mysql -u root -p myposcore < init_tnc_dummy.sql

# Load FAQ dummy data
mysql -u root -p myposcore < init_faq_dummy.sql
```

Atau jika menggunakan PostgreSQL:
```bash
psql -U postgres -d myposcore -f init_tnc_dummy.sql
psql -U postgres -d myposcore -f init_faq_dummy.sql
```

---

## Testing with Postman

1. Import collection: `MyPOSCore.postman_collection.json`
2. Import environment: `MyPOSCore.postman_environment.json`
3. Folder "FAQ" dan "Terms & Conditions" berisi semua endpoint
4. Public endpoints bisa langsung ditest tanpa login
5. Superadmin endpoints butuh login sebagai superadmin terlebih dahulu

### Login as Superadmin
```json
POST /api/auth/login
{
  "tenant_code": "supertenant",
  "branch_code": "superbranch",
  "username": "superadmin",
  "password": "123456"
}
```

Token akan otomatis disimpan di environment variable `auth_token`.

---

## Notes

- **Markdown Support**: Field `content` pada TnC mendukung Markdown formatting
- **Ordering**: FAQ diurutkan berdasarkan field `order` (ASC) kemudian `created_at` (DESC)
- **Soft Delete**: Delete operations menggunakan soft delete (GORM DeletedAt)
- **Active Flag**: Gunakan `is_active` untuk hide/show content tanpa delete
- **Public Access**: FAQ dan TnC public endpoints tidak perlu authentication untuk kemudahan akses user

---

## Related Files

- **Models**: `models/faq.go`, `models/tnc.go`
- **Handlers**: `handlers/faq_handler.go`, `handlers/tnc_handler.go`
- **Services**: `services/faq_service.go`, `services/tnc_service.go`
- **Routes**: `routes/routes.go`
- **Dummy Data**: `init_faq_dummy.sql`, `init_tnc_dummy.sql`
- **Postman**: `MyPOSCore.postman_collection.json`
