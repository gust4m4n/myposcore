# Database Index Migration Guide

## Overview
This migration adds database indexes for tenant code and branch code to improve query performance and ensure data integrity.

## Date
December 27, 2025

## Changes Made

### 1. Model Updates

#### Tenant Model (`models/tenant.go`)
```go
Code string `gorm:"size:50;uniqueIndex:idx_tenant_code;not null" json:"code"`
```
- Added named unique index `idx_tenant_code` on tenant code
- Ensures fast lookups by tenant code
- Prevents duplicate tenant codes

#### Branch Model (`models/branch.go`)
```go
TenantID uint   `gorm:"not null;index:idx_branch_tenant" json:"tenant_id"`
Code     string `gorm:"type:varchar(50);uniqueIndex:idx_branch_code;not null" json:"code"`
```
- Added named unique index `idx_branch_code` on branch code
- Added named index `idx_branch_tenant` on tenant_id
- Improves JOIN performance between branches and tenants
- Prevents duplicate branch codes

### 2. Database Migration

Run the migration script to add indexes to existing database:

```bash
psql -h localhost -U your_username -d your_database -f migration_add_indexes.sql
```

Or if using Docker:
```bash
docker exec -i postgres_container psql -U your_username -d your_database < migration_add_indexes.sql
```

## Benefits

### Performance Improvements
1. **Faster Tenant Lookups**: Queries filtering by tenant code will use index scan instead of sequential scan
2. **Faster Branch Lookups**: Queries filtering by branch code will use index scan
3. **Improved JOIN Performance**: Joins between branches and tenants will be faster
4. **Better INSERT Performance**: Duplicate checks are faster with indexes

### Example Query Performance

#### Before Index (Sequential Scan)
```sql
EXPLAIN SELECT * FROM tenants WHERE code = 'WMS001';
-- Seq Scan on tenants (cost=0.00..35.50 rows=1 width=...)
```

#### After Index (Index Scan)
```sql
EXPLAIN SELECT * FROM tenants WHERE code = 'WMS001';
-- Index Scan using idx_tenant_code on tenants (cost=0.15..8.17 rows=1 width=...)
```

## Verification

After running the migration, verify the indexes were created:

```sql
-- Check tenant indexes
\d tenants

-- Check branch indexes
\d branches

-- Or query pg_indexes
SELECT tablename, indexname, indexdef 
FROM pg_indexes 
WHERE tablename IN ('tenants', 'branches')
ORDER BY tablename, indexname;
```

Expected output should include:
- `idx_tenant_code` on `tenants(code)`
- `idx_branch_code` on `branches(code)`
- `idx_branch_tenant` on `branches(tenant_id)`

## Rollback

If you need to remove the indexes:

```sql
-- Remove tenant code index
DROP INDEX IF EXISTS idx_tenant_code;

-- Remove branch code index
DROP INDEX IF EXISTS idx_branch_code;

-- Remove branch tenant index
DROP INDEX IF EXISTS idx_branch_tenant;
```

## Notes

1. **GORM Auto Migration**: GORM will automatically create these indexes when you run the application for the first time or when AutoMigrate is called.

2. **Soft Deletes**: The indexes include `WHERE deleted_at IS NULL` condition to optimize queries that filter out soft-deleted records.

3. **Uniqueness**: Both tenant code and branch code are unique across the entire table (not scoped to tenant).

4. **Index Size**: These indexes are relatively small as they're on VARCHAR(50) fields. They should not significantly impact storage or insert performance.

5. **Production Deployment**: 
   - For large tables, consider running the migration during low-traffic periods
   - Monitor query performance before and after
   - The indexes are created with `IF NOT EXISTS` to be idempotent

## Testing

After migration, test the following scenarios:

### 1. Unique Constraint Test
```sql
-- Should succeed (unique code)
INSERT INTO tenants (name, code, is_active) VALUES ('Test Tenant 1', 'TEST001', true);

-- Should fail (duplicate code)
INSERT INTO tenants (name, code, is_active) VALUES ('Test Tenant 2', 'TEST001', true);
-- ERROR: duplicate key value violates unique constraint "idx_tenant_code"
```

### 2. Performance Test
```sql
-- Test tenant lookup by code
EXPLAIN ANALYZE SELECT * FROM tenants WHERE code = 'WMS001';

-- Test branch lookup by code
EXPLAIN ANALYZE SELECT * FROM branches WHERE code = 'BR001';

-- Test branch-tenant join
EXPLAIN ANALYZE 
SELECT b.*, t.name as tenant_name 
FROM branches b 
JOIN tenants t ON b.tenant_id = t.id 
WHERE b.code = 'BR001';
```

## Impact on Application

### Positive Impacts
- ✅ Faster authentication (tenant lookup by code)
- ✅ Faster API responses for tenant/branch queries
- ✅ Better database performance under load
- ✅ Prevents data integrity issues (duplicate codes)

### Minimal Impacts
- ⚠️ Slightly slower INSERT operations (negligible)
- ⚠️ Additional storage for indexes (minimal - ~1-2KB per 1000 records)

## Related Files

- `models/tenant.go` - Tenant model with index definition
- `models/branch.go` - Branch model with index definitions
- `migration_add_indexes.sql` - SQL migration script
- `INDEX_MIGRATION_GUIDE.md` - This documentation
