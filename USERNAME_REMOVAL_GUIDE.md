# Migration Guide: Username Removal

## Overview
This document outlines the changes made to remove the `username` field from the User model and use `email` as the primary login identifier.

## Breaking Changes

### 1. Database Schema
- **Removed Field**: `username` column from `users` table
- **New Constraint**: `email` field now has a UNIQUE constraint
- **Migration SQL**: Run `migration_remove_username.sql` to update the database

### 2. API Changes

#### Login API (`POST /api/v1/auth/login`)
**Before:**
```json
{
  "username": "john@foodcorner99.com",
  "password": "password123"
}
```

**After:**
```json
{
  "email": "john@foodcorner99.com",
  "password": "password123"
}
```

#### Register API (`POST /api/v1/auth/register`)
**Before:**
```json
{
  "tenant_code": "resto01",
  "branch_code": "pusat",
  "username": "john",
  "email": "john@foodcorner99.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**After:**
```json
{
  "tenant_code": "resto01",
  "branch_code": "pusat",
  "email": "john@foodcorner99.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

#### Create User API (`POST /api/v1/users`)
**Before:**
```json
{
  "username": "newuser",
  "email": "newuser@foodcorner99.com",
  "password": "password123",
  "full_name": "New User",
  "role": "user",
  "branch_id": 1
}
```

**After:**
```json
{
  "email": "newuser@foodcorner99.com",
  "password": "password123",
  "full_name": "New User",
  "role": "user",
  "branch_id": 1
}
```

#### Update User API (`PUT /api/v1/users/:id`)
**Before:**
```json
{
  "username": "updateduser",
  "email": "updated@foodcorner99.com",
  "full_name": "Updated Name"
}
```

**After:**
```json
{
  "email": "updated@foodcorner99.com",
  "full_name": "Updated Name"
}
```

#### Admin Change Password API (`POST /api/v1/auth/admin/change-password`)
**Before:**
```json
{
  "username": "targetuser",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**After:**
```json
{
  "email": "targetuser@foodcorner99.com",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

### 3. Response Changes

All API responses that previously included `username` field now exclude it:

**Before:**
```json
{
  "id": 1,
  "username": "john@foodcorner99.com",
  "email": "john@foodcorner99.com",
  "full_name": "John Doe",
  "role": "user"
}
```

**After:**
```json
{
  "id": 1,
  "email": "john@foodcorner99.com",
  "full_name": "John Doe",
  "role": "user"
}
```

### 4. JWT Token Changes

**Before:**
```json
{
  "user_id": 1,
  "tenant_id": 1,
  "username": "john@foodcorner99.com",
  "exp": 1735624800
}
```

**After:**
```json
{
  "user_id": 1,
  "tenant_id": 1,
  "email": "john@foodcorner99.com",
  "exp": 1735624800
}
```

## Migration Steps

1. **Backup Database**: Always backup your database before running migrations
   ```bash
   pg_dump myposdb > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

2. **Update Code**: Pull the latest code changes

3. **Run Migration**: Execute the migration SQL file
   ```bash
   psql -U postgres -d myposdb -f migration_remove_username.sql
   ```

4. **Build Application**:
   ```bash
   go build -o myposcore
   ```

5. **Update Postman Collection**: Import the updated `MyPOSCore.postman_collection.json`

6. **Test APIs**: Verify all authentication and user management APIs work correctly

## Files Modified

### Models
- `models/user.go` - Removed `Username` field, added UNIQUE constraint on `Email`

### DTOs
- `dto/auth.go` - Updated LoginRequest, RegisterRequest, UserProfile, UserDetailProfile, AdminChangePasswordRequest
- `dto/user.go` - Updated CreateUserRequest, UpdateUserRequest
- `dto/superadmin.go` - Updated UserResponse

### Services
- `services/login_service.go` - Changed to use email for authentication
- `services/register_service.go` - Removed username validation and assignment
- `services/user_service.go` - Removed username from create/update operations
- `services/auth_service.go` - Updated Register, Login, GetProfile methods
- `services/admin_change_password_service.go` - Changed to lookup user by email

### Handlers
- `handlers/login_handler.go` - Updated to use email in GenerateToken
- `handlers/register_handler.go` - Updated to use email in GenerateToken
- `handlers/user_handler.go` - Removed username from all operations
- `handlers/superadmin_handler.go` - Removed username from responses

### Middleware
- `middleware/auth.go` - Changed context to store email instead of username

### Utils
- `utils/jwt.go` - Updated Claims struct to use email instead of username

## Testing Checklist

- [ ] Login with email works correctly
- [ ] Register new user with email only
- [ ] Create user via API without username
- [ ] Update user without username field
- [ ] Admin change password using email
- [ ] JWT tokens contain email field
- [ ] All user list/get endpoints return data without username
- [ ] Email uniqueness is enforced
- [ ] Migration rollback tested (if needed)

## Rollback Plan

If you need to rollback:

1. Restore database from backup
2. Checkout previous code version
3. Rebuild and restart application

## Notes

- Email is now the ONLY login identifier
- Email must be unique across ALL users (not just per tenant)
- All existing usernames will be removed from the database
- Update any external integrations or scripts that reference username
