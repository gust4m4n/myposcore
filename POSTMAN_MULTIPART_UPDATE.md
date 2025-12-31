# Postman Collection - Multipart Support Update

## 📋 Summary

Update Postman collection untuk menambahkan contoh multipart/form-data pada API Create User dan Update User yang mendukung upload profile image.

## ✅ Changes Made

### 1. Create User API
**Before:** Hanya 1 request example (JSON mode)  
**After:** 2 request examples:
- **Create User (JSON)** - Mode JSON standar tanpa image
- **Create User (Multipart with Image)** - Mode form-data dengan optional image upload

### 2. Update User API
**Before:** Hanya 1 request example (JSON mode)  
**After:** 2 request examples:
- **Update User (JSON)** - Mode JSON standar tanpa image  
- **Update User (Multipart with Image)** - Mode form-data dengan optional image replacement

---

## 📊 Postman Collection Structure

```
MyPOSCore.postman_collection.json
└── Users Folder
    ├── List Users (with pagination)
    ├── Get User
    ├── Create User (JSON) ✅ NEW NAME
    │   └── Response: Create User Success
    ├── Create User (Multipart with Image) ✅ NEW
    │   └── Response: Create User with Image Success
    ├── Update User (JSON) ✅ NEW NAME
    │   └── Response: Update User Success
    ├── Update User (Multipart with Image) ✅ NEW
    │   └── Response: Update User with New Image Success
    └── Delete User
```

---

## 🎯 Request Examples

### Create User (JSON) - Lines ~3052
```http
POST /api/v1/users
Content-Type: application/json
Authorization: Bearer {{auth_token}}

{
  "email": "newuser@example.com",
  "password": "password123",
  "full_name": "New User",
  "role": "user",
  "branch_id": 1,
  "is_active": true
}
```

**Response:**
```json
{
  "message": "User created successfully",
  "data": {
    "id": 10,
    "tenant_id": 17,
    "branch_id": 1,
    "email": "newuser@example.com",
    "full_name": "New User",
    "role": "user",
    "is_active": true,
    "image": null,
    "created_at": "2025-12-27 11:00:00"
  }
}
```

---

### Create User (Multipart with Image) - Lines ~3135
```http
POST /api/v1/users
Content-Type: multipart/form-data
Authorization: Bearer {{auth_token}}

Form Fields:
- email: johndoe@example.com
- password: password123
- full_name: John Doe
- role: user
- branch_id: 1
- is_active: true
- image: [FILE] profile.jpg (optional)
```

**Response:**
```json
{
  "message": "User created successfully",
  "data": {
    "id": 11,
    "tenant_id": 17,
    "branch_id": 1,
    "email": "johndoe@example.com",
    "full_name": "John Doe",
    "role": "user",
    "is_active": true,
    "image": "uploads/profiles/user_11_1735291200.jpg",
    "created_at": "2025-12-27 11:00:00"
  }
}
```

---

### Update User (JSON) - Lines ~3283
```http
PUT /api/v1/users/1
Content-Type: application/json
Authorization: Bearer {{auth_token}}

{
  "full_name": "Updated User Name",
  "role": "branchadmin",
  "is_active": true
}
```

**Response:**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 17,
    "branch_id": 1,
    "email": "user1@example.com",
    "full_name": "Updated User Name",
    "role": "branchadmin",
    "is_active": true,
    "image": "uploads/profiles/user_1_1735200000.jpg",
    "created_at": "2025-12-27 10:00:00"
  }
}
```

---

### Update User (Multipart with Image) - Lines ~3368
```http
PUT /api/v1/users/1
Content-Type: multipart/form-data
Authorization: Bearer {{auth_token}}

Form Fields (ALL OPTIONAL):
- full_name: Jane Doe Updated
- role: branchadmin
- is_active: true
- image: [FILE] new_profile.jpg (optional, replaces old image)
```

**Response:**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 17,
    "branch_id": 1,
    "email": "user1@example.com",
    "full_name": "Jane Doe Updated",
    "role": "branchadmin",
    "is_active": true,
    "image": "uploads/profiles/user_1_1735291800.jpg",
    "created_at": "2025-12-27 10:00:00"
  }
}
```

---

## 📝 Key Features

### Multipart Requests Include:
1. **Detailed descriptions** explaining dual content-type support
2. **Form field definitions** with type and requirement status
3. **How-to instructions** for testing in Postman
4. **Response examples** showing image URLs
5. **Partial update support** for Update User (all fields optional)

### Documentation Highlights:
- **Required vs Optional fields** clearly marked
- **Image field behavior** explained (optional, replaces old image)
- **File upload instructions** for Postman users
- **Use cases** for when to use JSON vs Multipart mode

---

## 🔍 Testing Instructions

### In Postman:

#### For JSON Requests:
1. Select request: "Create User (JSON)" or "Update User (JSON)"
2. Body already configured as JSON
3. Click **Send**

#### For Multipart Requests:
1. Select request: "Create User (Multipart with Image)" or "Update User (Multipart with Image)"
2. Go to **Body** tab → **form-data**
3. For image field:
   - Hover over the **image** field
   - Change type dropdown to **File**
   - Click **Select Files**
   - Choose a JPEG/PNG from your computer
4. Click **Send**

---

## 🎨 Response Differences

| Mode | Image Field | Example Value |
|------|-------------|---------------|
| **JSON** (no image) | `null` | `"image": null` |
| **Multipart** (no image) | `null` | `"image": null` |
| **Multipart** (with image) | URL path | `"image": "uploads/profiles/user_11_1735291200.jpg"` |

---

## 📦 Files Modified

1. **MyPOSCore.postman_collection.json**
   - Renamed: "Create User" → "Create User (JSON)"
   - Added: "Create User (Multipart with Image)" with full form-data example
   - Renamed: "Update User" → "Update User (JSON)"
   - Added: "Update User (Multipart with Image)" with partial update example
   - Updated all response examples to include `image` field
   - Added detailed descriptions for both modes

2. **MULTIPART_USER_GUIDE.md** (NEW)
   - Complete guide for multipart support
   - JSON vs Multipart comparison
   - Form field specifications
   - Frontend integration examples (JavaScript, React)
   - Security notes and best practices

3. **README.md**
   - Added link to MULTIPART_USER_GUIDE.md
   - Brief explanation of dual-mode support

---

## ✅ Validation

JSON syntax validation:
```bash
python3 -m json.tool MyPOSCore.postman_collection.json > /dev/null
```
**Result:** ✅ Valid JSON

---

## 🚀 Benefits

1. **Clear Examples** - Users can see both JSON and Multipart modes side by side
2. **Easy Testing** - Pre-configured form-data requests ready to use
3. **Better Documentation** - Detailed descriptions explain when to use each mode
4. **Image Upload Support** - Users can now upload profile images during user creation/update
5. **Flexibility** - Users can choose simple JSON or full multipart based on needs

---

## 📚 Related Documentation

- [MULTIPART_USER_GUIDE.md](MULTIPART_USER_GUIDE.md) - Complete multipart guide with frontend examples
- [USER_API_GUIDE.md](USER_API_GUIDE.md) - General user API documentation
- [PAGINATION_GUIDE.md](PAGINATION_GUIDE.md) - Pagination for list APIs
- [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md) - How to use Postman collection

---

## 🔄 Migration Notes

### For Existing API Consumers:

**No breaking changes!** 
- Existing JSON requests continue to work as before
- Multipart support is ADDITIONAL feature, not replacement
- `image` field added to responses (will be `null` if no image)

### Recommended Actions:
1. ✅ Import updated Postman collection
2. ✅ Test both JSON and Multipart modes
3. ✅ Update frontend to support image uploads if needed
4. ✅ Use multipart mode for user creation forms with photo upload
5. ✅ Continue using JSON mode for simple API integrations

---

**Last Updated:** 2025-12-27  
**Postman Collection Version:** Compatible with Postman v10+  
**Backend Compatibility:** MyPOSCore v1.0+
