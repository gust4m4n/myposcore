# Postman Collection Code Field Removal Guide

## Issue
The MyPOSCore.postman_collection.json file has pre-existing JSON syntax errors (around line 4081) that prevent automated bulk updates.

## Critical Changes Required

### 1. ❌ REMOVE: Registration Endpoint
**Location**: Line ~195  
**Endpoint**: `POST /api/v1/auth/register`  
**Action**: Delete entire "Register User" request block  
**Reason**: Registration API removed - users created by admins only

### 2. ✏️ UPDATE: Tenant Management

#### Create Tenant (`POST /api/v1/superadmin/tenants`)
**Remove from formdata**:
```json
{
  "key": "code",
  "value": "TENANT002",  // DELETE THIS FIELD
  "type": "text"
}
```

#### Update Tenant (`PUT /api/v1/superadmin/tenants/:id`)
**Remove from formdata**:
```json
{
  "key": "code",
  "value": "TENANT002",  // DELETE THIS FIELD
  "type": "text"
}
```

**Update Response Examples** - Remove code field:
```json
{
  "data": {
    "id": 1,
    "name": "Updated Tenant Name",
    "code": "TENANT002",  // ❌ REMOVE
    "description": "..."
  }
}
```

### 3. ✏️ UPDATE: Branch Management

#### Create Branch (`POST /api/v1/superadmin/branches`)
**Remove from formdata**:
```json
{
  "key": "code",
  "value": "baru",  // DELETE THIS FIELD
  "type": "text"
}
```

#### Update Branch (`PUT /api/v1/superadmin/branches/:id`)
**Remove from formdata**:
```json
{
  "key": "code",
  "value": "updated",  // DELETE THIS FIELD
  "type": "text"
}
```

**Update Response Examples** - Remove code field:
```json
{
  "data": {
    "id": 1,
    "tenant_id": 1,
    "name": "Cabang Pusat",
    "code": "pusat",  // ❌ REMOVE
    "description": "..."
  }
}
```

### 4. ✏️ UPDATE: Login Response Examples
**Endpoint**: `POST /api/v1/auth/login`

Remove `code` from tenant and branch objects in response examples:
```json
{
  "message": "Login successful",
  "data": {
    "tenant": {
      "id": 1,
      "name": "Food Corner 99",
      "code": "resto01",  // ❌ REMOVE
      "description": "..."
    },
    "branch": {
      "id": 1,
      "name": "Cabang Pusat",
      "code": "pusat",  // ❌ REMOVE
      "description": "..."
    }
  }
}
```

### 5. ✏️ UPDATE: List Tenant/Branch Responses
**Endpoints**:
- `GET /api/v1/superadmin/tenants`
- `GET /api/v1/superadmin/branches`
- `GET /dev/tenants`
- `GET /dev/tenants/:id/branches`

Remove `code` field from all response examples.

### 6. ✏️ UPDATE: Collection Description
**Location**: `info.description` (top of file)

Replace:
- "Multi-tenancy support dengan tenant code" → "Multi-tenancy support dengan tenant ID"
- "Multi-branch per tenant dengan branch code" → "Multi-branch per tenant dengan branch ID"

## Quick Reference: Fields to Remove

| Endpoint | Method | Field Location | Field Name |
|----------|--------|----------------|------------|
| Register | POST | request body | N/A - DELETE ENTIRE ENDPOINT |
| Create Tenant | POST | formdata | `code` |
| Update Tenant | PUT | formdata | `code` |
| Create Branch | POST | formdata | `code` |
| Update Branch | PUT | formdata | `code` |
| All responses | - | response body | `"code": "..."` |

## Recommended Update Process

### Method 1: Import to Postman (Easiest)
1. Open Postman application
2. Import `MyPOSCore.postman_collection.json`
3. Postman may auto-fix JSON syntax errors
4. Manually delete/edit requests as described above
5. Export collection as JSON (v2.1)
6. Replace the file

### Method 2: Manual JSON Edit
1. Fix JSON syntax error at line 4081 first:
   - Remove duplicate response block
   - Ensure proper closing braces
2. Use find/replace to remove code fields:
   - Search: `"code":\s*"[^"]*",?\n` (regex)
   - Replace with: (empty)
3. Manually verify critical endpoints
4. Validate JSON syntax: `python3 -m json.tool file.json`

### Method 3: Selective Manual Updates (Current)
Since JSON has errors, update manually in text editor:
1. Search for `"code"` in the file
2. Remove each occurrence that's part of tenant/branch data
3. Be careful not to remove HTTP status codes (200, 404, etc.)
4. Delete Register User endpoint block (lines ~195-290)

## Validation Checklist
After updating:
- [ ] JSON syntax is valid (`python3 -m json.tool`)
- [ ] Register endpoint removed
- [ ] No `code` field in Create Tenant request
- [ ] No `code` field in Update Tenant request
- [ ] No `code` field in Create Branch request
- [ ] No `code` field in Update Branch request
- [ ] No `code` in login response examples
- [ ] No `code` in list tenant/branch responses
- [ ] Import to Postman and test requests

## Testing After Update
Test these requests in Postman:
1. Login - verify response has no tenant/branch code
2. Create Tenant - verify doesn't require code
3. Create Branch - verify doesn't require code
4. List Tenants - verify response has no code
5. List Branches - verify response has no code

## Current Status
⚠️ **JSON file has syntax errors that need to be fixed first**  
✅ **Backend code updated** - no longer uses/returns code fields  
⏳ **Postman collection** - needs manual updates as described above
