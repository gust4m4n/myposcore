# User API Multipart/Form-Data Guide

## Overview

API **Create User** dan **Update User** mendukung **2 mode content-type**:

1. **`application/json`** - Mode JSON standar (simple, tanpa image)
2. **`multipart/form-data`** - Mode form-data dengan support upload image profile

## ✅ Create User - Dual Mode Support

### Mode 1: JSON (Simple, No Image)

**Endpoint:** `POST /api/v1/users`  
**Content-Type:** `application/json`

**Request Body:**
```json
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

### Mode 2: Multipart (With Image Upload)

**Endpoint:** `POST /api/v1/users`  
**Content-Type:** `multipart/form-data`

**Form Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | text | ✅ Yes | Valid email format, globally unique |
| `password` | text | ✅ Yes | Min 6 characters |
| `full_name` | text | ✅ Yes | User's full name |
| `role` | text | ✅ Yes | `user` \| `branchadmin` \| `tenantadmin` |
| `branch_id` | text | ✅ Yes | Must exist and belong to tenant |
| `is_active` | text | ✅ Yes | `true` \| `false` |
| `image` | file | ⭕ Optional | Profile image (JPEG/PNG) |

**Postman Example:**
```
Form Data:
- email: johndoe@example.com
- password: password123
- full_name: John Doe
- role: user
- branch_id: 1
- is_active: true
- image: [Select File] profile.jpg
```

**Response with Image:**
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

## ✅ Update User - Dual Mode Support

### Mode 1: JSON (Simple, No Image)

**Endpoint:** `PUT /api/v1/users/{id}`  
**Content-Type:** `application/json`

**Request Body (All fields optional - partial update):**
```json
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

### Mode 2: Multipart (With Image Upload/Replacement)

**Endpoint:** `PUT /api/v1/users/{id}`  
**Content-Type:** `multipart/form-data`

**Form Fields (ALL OPTIONAL - send only what you want to update):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `full_name` | text | ⭕ Optional | Update full name |
| `role` | text | ⭕ Optional | Update role: `user` \| `branchadmin` \| `tenantadmin` |
| `is_active` | text | ⭕ Optional | Update status: `true` \| `false` |
| `password` | text | ⭕ Optional | New password (will be hashed) |
| `branch_id` | text | ⭕ Optional | New branch assignment |
| `image` | file | ⭕ Optional | New profile image (REPLACES old image) |

**Postman Examples:**

#### Example 1: Update only name
```
Form Data:
- full_name: Jane Doe Updated
```

#### Example 2: Update only image
```
Form Data:
- image: [Select File] new_profile.jpg
```

#### Example 3: Update name + role + image
```
Form Data:
- full_name: Jane Doe Updated
- role: branchadmin
- image: [Select File] new_profile.jpg
```

**Response with New Image:**
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

## 🎯 When to Use Which Mode?

### Use **JSON Mode** when:
- ✅ Creating/updating user data WITHOUT image
- ✅ Quick updates (faster, simpler)
- ✅ API integrations where image upload is not needed
- ✅ Mobile apps without camera/gallery access

### Use **Multipart Mode** when:
- ✅ Need to upload profile image during user creation
- ✅ Need to update/replace existing profile image
- ✅ Mobile apps with camera/gallery integration
- ✅ Web forms with file upload capability

---

## 📝 Important Notes

### Content-Type Detection
Backend automatically detects content-type:
```go
contentType := c.GetHeader("Content-Type")
if strings.Contains(contentType, "multipart/form-data") {
    // Process as multipart
} else {
    // Process as JSON
}
```

### Image Upload Behavior
- **Image field is OPTIONAL** in multipart mode
- If no image provided: `image` field will be `null`
- If image provided: Image will be uploaded to `/uploads/profiles/`
- Image filename format: `user_{id}_{timestamp}.{ext}`
- Supported formats: JPEG, PNG
- On update: New image REPLACES old image

### Error Handling
- User creation/update will **succeed** even if image upload fails
- Error message will be returned if image upload fails
- User record will have `image: null` if upload fails

### Security
- Image uploads are validated for file type
- Only JPEG and PNG allowed
- File size limits may apply (check server config)
- Images stored in `/uploads/profiles/` directory

---

## 🔧 Testing in Postman

### For JSON Mode:
1. Set **Content-Type** header to `application/json`
2. Select **Body** → **raw** → **JSON**
3. Enter JSON body
4. Send request

### For Multipart Mode:
1. **REMOVE** Content-Type header (Postman sets it automatically)
2. Select **Body** → **form-data**
3. Add text fields: email, password, full_name, etc.
4. Add file field: 
   - Key: `image`
   - Type: **File**
   - Click **Select Files** → Choose JPEG/PNG
5. Send request

---

## 💡 Frontend Integration Examples

### JavaScript/Fetch (JSON Mode)
```javascript
const response = await fetch('http://localhost:8080/api/v1/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    email: 'newuser@example.com',
    password: 'password123',
    full_name: 'New User',
    role: 'user',
    branch_id: 1,
    is_active: true
  })
});
```

### JavaScript/Fetch (Multipart Mode)
```javascript
const formData = new FormData();
formData.append('email', 'newuser@example.com');
formData.append('password', 'password123');
formData.append('full_name', 'New User');
formData.append('role', 'user');
formData.append('branch_id', '1');
formData.append('is_active', 'true');
formData.append('image', fileInput.files[0]); // File from <input type="file">

const response = await fetch('http://localhost:8080/api/v1/users', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
    // Don't set Content-Type - browser will set it automatically with boundary
  },
  body: formData
});
```

### React Example (Multipart with Image)
```jsx
const handleSubmit = async (e) => {
  e.preventDefault();
  
  const formData = new FormData();
  formData.append('email', email);
  formData.append('password', password);
  formData.append('full_name', fullName);
  formData.append('role', role);
  formData.append('branch_id', branchId);
  formData.append('is_active', isActive);
  
  // Only append image if file is selected
  if (imageFile) {
    formData.append('image', imageFile);
  }
  
  try {
    const response = await fetch(`${API_URL}/api/v1/users`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    });
    
    const result = await response.json();
    console.log('User created:', result.data);
    console.log('Image URL:', result.data.image); // Use this to display image
  } catch (error) {
    console.error('Error:', error);
  }
};
```

---

## 📊 API Summary

| Endpoint | Method | JSON Mode | Multipart Mode | Image Support |
|----------|--------|-----------|----------------|---------------|
| `/api/v1/users` | POST | ✅ Yes | ✅ Yes | ✅ Optional |
| `/api/v1/users/{id}` | PUT | ✅ Yes | ✅ Yes | ✅ Optional |

**Both modes work independently - choose based on your needs!**
