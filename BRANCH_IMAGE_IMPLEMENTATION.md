# Branch Image Upload Implementation - Quick Reference

## ✅ Completed Implementation

### Files Modified:
1. **models/branch.go** - Added `Image string` field (varchar 500)
2. **dto/superadmin.go** - Added `Image` to BranchResponse
3. **dto/auth.go** - Added `Image` to BranchInfo (login response)
4. **handlers/superadmin_handler.go** - Updated branch handlers:
   - CreateBranch → multipart/form-data with image upload
   - UpdateBranch → multipart/form-data with old image deletion
   - DeleteBranch → auto delete image file
   - ListBranches → include image in response
5. **handlers/login_handler.go** - Include branch image in response
6. **services/superadmin_branch_service.go** - Updated service methods:
   - CreateBranch(req, imageURL) - accepts imageURL parameter
   - UpdateBranch(id, req, imageURL) - conditional image update
   - GetBranchByID(id) - new method to retrieve branch
7. **MyPOSCore.postman_collection.json** - Updated Create/Update Branch requests to multipart
8. **migration_add_branch_image.sql** - Database migration

### Features:
✅ Upload image saat create branch (optional)
✅ Update/replace image saat update branch (optional)
✅ Auto delete old image when updating with new one
✅ Auto delete image when deleting branch
✅ File validation: jpg, jpeg, png, gif, webp, max 5MB
✅ Unique filename: branch_{code}_{timestamp}.{ext}
✅ Rollback image if database operation fails
✅ Image URL in all responses (create, update, list, login)

### Storage Structure:
```
uploads/
├── tenants/
│   └── tenant_CODE_TIMESTAMP.{ext}
└── branches/
    └── branch_CODE_TIMESTAMP.{ext}
```

### Database Migration:
```bash
psql -U <username> -d <database> -f migration_add_branch_image.sql
```

Or manually:
```sql
ALTER TABLE branches ADD COLUMN IF NOT EXISTS image VARCHAR(500);
```

### API Endpoints Changed:

#### Create Branch
**Before:** `Content-Type: application/json`
**After:** `Content-Type: multipart/form-data`

```
POST /api/v1/superadmin/branches
Form fields:
- tenant_id (text)
- name, code, description, address, website, email, phone (text)
- is_active (text: "true"/"false")
- image (file, optional)
```

**Response includes:**
```json
{
  "data": {
    "id": 5,
    "tenant_id": 1,
    "name": "Cabang Baru",
    "code": "baru",
    "image": "/uploads/branches/branch_baru_1704010200.png",
    ...
  }
}
```

#### Update Branch
Same as Create - multipart/form-data with optional image upload.
Old image automatically deleted when new image provided.

#### Delete Branch
Enhanced to delete image file before soft deleting record.

#### List Branches & Login
Response now includes `image` field for each branch.

### Postman Testing:
1. Import updated collection
2. Use Create Branch request with multipart form-data
3. Attach image file in `image` field
4. Verify file saved in `uploads/branches/`
5. Test Update to replace image
6. Test Delete to remove image file

### Build Status:
✅ **Compilation successful** - No errors

### Summary:
Branch image upload implementation complete, mirroring tenant image upload feature. All CRUD operations support image upload with proper file validation, cleanup, and rollback mechanisms.
