# Code Field Removal Summary

## Overview
Removed `code` fields from `tenants` and `branches` tables. The system now uses only `tenant_id` and `branch_id` for identification.

## Changes Made

### 1. Models (Database Schema)
- **models/tenant.go**: Removed `Code` field from `Tenant` struct
- **models/branch.go**: Removed `Code` field from `Branch` struct

### 2. DTOs (Data Transfer Objects)
- **dto/auth.go**: 
  - Removed `Code` from `TenantInfo`
  - Removed `Code` from `BranchInfo`
  - Removed `Code` from `TenantDetailProfile`
  - Removed `Code` from `BranchDetailProfile`

- **dto/superadmin.go**:
  - Removed `Code` from `CreateTenantRequest`
  - Removed `Code` from `UpdateTenantRequest`
  - Removed `Code` from `CreateBranchRequest`
  - Removed `Code` from `UpdateBranchRequest`
  - Removed `Code` from `TenantResponse`
  - Removed `Code` from `BranchResponse`

### 3. Services (Business Logic)
- **services/superadmin_tenant_service.go**:
  - Removed code uniqueness validation from `CreateTenant()`
  - Removed code uniqueness validation from `UpdateTenant()`
  - Removed code assignment logic

- **services/superadmin_branch_service.go**:
  - Removed code uniqueness validation from `CreateBranch()`
  - Removed code uniqueness validation from `UpdateBranch()`
  - Removed code assignment logic

- **services/auth_service.go**:
  - Removed `Code` from profile responses in `GetProfile()`

- **services/superadmin_dashboard_service.go**:
  - Removed `Code` from `TenantResponse` construction

### 4. Handlers (API Controllers)
- **handlers/login_handler.go**:
  - Removed `Code` from `TenantInfo` in login response
  - Removed `Code` from `BranchInfo` in login response

- **handlers/superadmin_handler.go**:
  - **ListTenants**: Removed `Code` from response
  - **CreateTenant**: 
    - Removed `Code` validation
    - Updated filename generation to use timestamp
    - Removed `Code` from response
  - **UpdateTenant**:
    - Removed `Code` from request parsing
    - Removed `Code` validation
    - Updated filename generation to use tenant ID
    - Removed `Code` from response
  - **ListBranches**: Removed `Code` from response
  - **CreateBranch**:
    - Removed `Code` from request parsing
    - Removed `Code` validation
    - Updated filename generation to use timestamp
    - Removed `Code` from response
  - **UpdateBranch**:
    - Removed `Code` from request parsing
    - Removed `Code` validation
    - Updated filename generation to use branch ID
    - Removed `Code` from response

### 5. Database Migration
- **migration_remove_code_columns.sql**: Created migration to:
  - Drop `idx_tenant_code` index
  - Drop `idx_branch_code` index
  - Drop `code` column from `tenants` table
  - Drop `code` column from `branches` table
  - Includes verification query

## Impact

### API Changes
All tenant and branch APIs no longer accept or return `code` fields:
- `POST /superadmin/tenants` - No longer requires `code`
- `PUT /superadmin/tenants/:id` - No longer accepts `code`
- `POST /superadmin/branches` - No longer requires `code`
- `PUT /superadmin/branches/:id` - No longer accepts `code`
- `GET /superadmin/tenants` - Response no longer includes `code`
- `GET /superadmin/branches` - Response no longer includes `code`
- `POST /auth/login` - Response no longer includes tenant/branch `code`
- `GET /auth/profile` - Response no longer includes tenant/branch `code`

### Validation Changes
- Tenant creation/update now only requires `name` (code no longer required)
- Branch creation/update now only requires `name` (code no longer required)
- No more uniqueness checks on code fields

### File Upload Changes
- Tenant image filenames now use: `tenant_{timestamp}_{nanoseconds}.ext`
- Branch image filenames now use: `branch_{id}_{timestamp}.ext` for updates or `branch_{timestamp}_{nanoseconds}.ext` for creation

## Migration Steps

### Step 1: Backup Database
```bash
pg_dump -U postgres myposcore > myposcore_backup_before_code_removal.sql
```

### Step 2: Run Migration
```bash
psql -U postgres -d myposcore -f migration_remove_code_columns.sql
```

### Step 3: Verify Migration
The migration includes a verification query. Check output shows:
- `tenants` table has no `code` column
- `branches` table has no `code` column
- Both indexes are dropped

### Step 4: Test Application
```bash
# Build and run
go build -o myposcore
./myposcore

# Test key endpoints:
# - Login (should not return code fields)
# - Get profile (should not return code fields)
# - Create/update tenant (should not require code)
# - Create/update branch (should not require code)
```

## Rollback Plan
If rollback is needed:

```sql
-- Add code columns back
ALTER TABLE tenants ADD COLUMN code VARCHAR(50);
ALTER TABLE branches ADD COLUMN code VARCHAR(50);

-- Recreate indexes
CREATE UNIQUE INDEX idx_tenant_code ON tenants(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_branch_code ON branches(tenant_id, code) WHERE deleted_at IS NULL;

-- Note: Old code values will be lost, would need to restore from backup
```

## Verification Checklist
- [x] All Go code compiles successfully
- [x] All models updated (Tenant, Branch)
- [x] All DTOs updated (auth, superadmin)
- [x] All services updated (tenant, branch, auth, dashboard)
- [x] All handlers updated (login, superadmin)
- [x] Migration SQL created
- [ ] Migration SQL executed
- [ ] API tests passed
- [ ] Manual testing completed

## Notes
- The code fields were originally used for human-readable identification
- System now relies entirely on numeric IDs for tenant/branch identification
- This simplifies the data model and removes redundant uniqueness constraints
- Frontend applications should be updated to no longer send/expect code fields
