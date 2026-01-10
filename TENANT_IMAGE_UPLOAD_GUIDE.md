# Tenant Image Upload Feature

## Overview
Fitur upload image untuk tenant logo/brand image dengan support multipart form data.

## Features
- ✅ Upload image saat create tenant
- ✅ Update/replace image saat update tenant
- ✅ Auto delete old image saat update dengan image baru
- ✅ Auto delete image saat delete tenant
- ✅ File validation (type & size)
- ✅ Unique filename dengan timestamp
- ✅ Rollback image jika database operation gagal

## Supported Image Formats
- JPG/JPEG
- PNG
- GIF
- WebP

**Maximum file size**: 5MB

## API Endpoints

### 1. Create Tenant with Image

**Endpoint**: `POST /api/superadmin/tenants`

**Content-Type**: `multipart/form-data`

**Form Fields**:
- `name` (string, required) - Nama tenant
- `code` (string, required) - Kode unik tenant
- `description` (string, required) - Deskripsi tenant
- `address` (string, required) - Alamat lengkap
- `website` (string, required) - URL website
- `email` (string, required) - Email kontak
- `phone` (string, required) - Nomor telepon
- `is_active` (boolean, required) - Status aktif
- `image` (file, optional) - Logo/image tenant

**Success Response** (200 OK):
```json
{
  "message": "Tenant created successfully",
  "data": {
    "id": 19,
    "name": "New Tenant",
    "code": "TENANT002",
    "description": "A new business tenant with logo",
    "address": "Jl. Kebon Jeruk No. 789, Jakarta",
    "website": "https://newtenant.com",
    "email": "contact@newtenant.com",
    "phone": "021-99887766",
    "image": "/uploads/tenants/tenant_TENANT002_1704003600.png",
    "is_active": true,
    "created_at": "2025-12-30T16:00:00+07:00"
  }
}
```

**Error Responses**:

400 Bad Request - Invalid file type:
```json
{
  "error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"
}
```

400 Bad Request - File too large:
```json
{
  "error": "File size too large. Maximum 5MB"
}
```

### 2. Update Tenant with Image

**Endpoint**: `PUT /api/superadmin/tenants/:tenant_id`

**Content-Type**: `multipart/form-data`

**Form Fields**: Same as Create Tenant

**Behavior**:
- Jika field `image` diisi dengan file baru:
  - Old image akan dihapus dari storage
  - New image akan diupload
  - Database akan diupdate dengan path baru
- Jika field `image` tidak diisi:
  - Image existing tetap tidak berubah
  - Hanya text fields yang diupdate

**Success Response** (200 OK):
```json
{
  "message": "Tenant updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Tenant Name",
    "code": "TENANT002",
    "description": "Updated business description with new logo",
    "address": "New Address, Jakarta",
    "website": "https://updatedtenant.com",
    "email": "newemail@tenant.com",
    "phone": "021-11112222",
    "image": "/uploads/tenants/tenant_TENANT002_1704007200.png",
    "is_active": true,
    "created_at": "2025-12-30T10:00:00+07:00"
  }
}
```

### 3. Delete Tenant

**Endpoint**: `DELETE /api/superadmin/tenants/:tenant_id`

**Behavior**:
- Soft delete tenant record dari database
- **Auto delete image file** dari storage jika ada

**Success Response** (200 OK):
```json
{
  "message": "Tenant deleted successfully"
}
```

### 4. List Tenants

**Endpoint**: `GET /api/superadmin/tenants`

Response sekarang include field `image`:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Warteg 123",
      "code": "warteg123",
      "description": "Warung Tegal dengan menu tradisional Indonesia",
      "address": "Jl. Raya Sudirman No. 123, Jakarta Pusat",
      "website": "https://warteg123.com",
      "email": "info@warteg123.com",
      "phone": "021-12345678",
      "image": "/uploads/tenants/tenant_warteg123_1704000000.png",
      "is_active": true,
      "created_at": "2025-12-24T10:00:00+07:00"
    }
  ]
}
```

### 5. Login Response

Login response sekarang include tenant image:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGc...",
    "user": { ... },
    "tenant": {
      "id": 1,
      "name": "Warteg 123",
      "code": "warteg123",
      "description": "Warung Tegal dengan menu tradisional Indonesia",
      "address": "Jl. Raya Sudirman No. 123, Jakarta Pusat",
      "website": "https://warteg123.com",
      "email": "info@warteg123.com",
      "phone": "021-12345678",
      "image": "/uploads/tenants/tenant_warteg123_1704000000.png",
      "is_active": true
    },
    "branch": { ... }
  }
}
```

## Storage Structure

```
uploads/
└── tenants/
    ├── tenant_warteg123_1704000000.png
    ├── tenant_fashionstore99_1704001000.jpg
    └── tenant_TENANT002_1704003600.webp
```

**Filename Pattern**: `tenant_{code}_{timestamp}.{ext}`

## Database Migration

File: `migration_add_tenant_image.sql`

```sql
ALTER TABLE tenants 
ADD COLUMN IF NOT EXISTS image VARCHAR(500);

COMMENT ON COLUMN tenants.image IS 'URL path to tenant logo/image file';
```

## Testing with Postman

1. Import updated `MyPOSCore.postman_collection.json`
2. Set environment dengan superadmin token
3. Test Create Tenant:
   - Pilih "Create Tenant" request
   - Body sudah dalam format form-data
   - Click pada field `image` dan pilih file dari komputer
   - Send request
4. Verify:
   - Check response body untuk `image` field
   - Verify file ada di `uploads/tenants/` directory
5. Test Update Tenant:
   - Pilih "Update Tenant" request
   - Upload image baru
   - Verify old image terhapus, new image tersimpan
6. Test Delete Tenant:
   - Delete tenant
   - Verify image file juga terhapus dari storage

## Frontend Integration Guide

### Using Fetch API

```javascript
async function createTenantWithImage(formData) {
  const response = await fetch('http://localhost:8080/api/superadmin/tenants', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData // FormData object
  });
  return await response.json();
}

// Usage
const formData = new FormData();
formData.append('name', 'New Tenant');
formData.append('code', 'TENANT002');
formData.append('description', 'Description here');
formData.append('address', 'Address here');
formData.append('website', 'https://example.com');
formData.append('email', 'email@example.com');
formData.append('phone', '021-12345678');
formData.append('is_active', 'true');

// Get file from input element
const fileInput = document.getElementById('imageInput');
if (fileInput.files[0]) {
  formData.append('image', fileInput.files[0]);
}

const result = await createTenantWithImage(formData);
console.log('Image URL:', result.data.image);
```

### Using Axios

```javascript
import axios from 'axios';

const createTenantWithImage = async (tenantData, imageFile) => {
  const formData = new FormData();
  
  // Append text fields
  Object.keys(tenantData).forEach(key => {
    formData.append(key, tenantData[key]);
  });
  
  // Append image file if exists
  if (imageFile) {
    formData.append('image', imageFile);
  }
  
  const response = await axios.post(
    'http://localhost:8080/api/superadmin/tenants',
    formData,
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'multipart/form-data'
      }
    }
  );
  
  return response.data;
};
```

### Displaying Image

```javascript
// Image URL from API
const imageUrl = tenant.image; // e.g., "/uploads/tenants/tenant_warteg123_1704000000.png"

// Display in HTML
<img 
  src={`http://localhost:8080${imageUrl}`} 
  alt={tenant.name}
  onError={(e) => {
    e.target.src = '/default-tenant-logo.png'; // Fallback image
  }}
/>
```

## Security Notes

1. **File Validation**: Backend validates file type dan size sebelum upload
2. **Unique Filenames**: Menggunakan timestamp untuk avoid filename collision
3. **Directory Permissions**: Upload directory harus memiliki write permissions (0755)
4. **Rollback Mechanism**: Jika database operation gagal, uploaded file akan dihapus
5. **Old File Cleanup**: Old image otomatis dihapus saat update/delete

## Error Handling

Handler sudah implement error handling untuk:
- Invalid file type
- File size too large
- Failed to create directory
- Failed to save file
- Database errors dengan rollback
- Tenant not found

## Performance Considerations

1. **File Size Limit**: 5MB maximum untuk maintain performance
2. **Storage**: Local filesystem storage di `uploads/tenants/`
3. **Future Enhancement**: Consider moving to cloud storage (S3, GCS) untuk production scale

## Related Files

- `models/tenant.go` - Added `Image` field
- `dto/superadmin.go` - Added `Image` to TenantResponse
- `dto/auth.go` - Added `Image` to TenantInfo
- `handlers/superadmin_handler.go` - Multipart upload handlers
- `services/superadmin_tenant_service.go` - Service layer dengan imageURL parameter
- `handlers/login_handler.go` - Include image in response
- `migration_add_tenant_image.sql` - Database migration
- `MyPOSCore.postman_collection.json` - Updated requests & responses

## Changelog

**Version 1.1.0** - December 2025
- ✅ Added image upload support untuk tenant
- ✅ Multipart form data handling
- ✅ Auto delete old images
- ✅ File validation (type & size)
- ✅ Rollback mechanism
- ✅ Updated all CRUD operations
- ✅ Updated Postman collection
- ✅ Updated login response
