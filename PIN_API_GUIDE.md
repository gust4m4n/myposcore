# PIN Management API Guide

API endpoints untuk manajemen PIN (Personal Identification Number) 6 digit di MyPOSCore.

## 📋 Endpoints

### 1. Create PIN
**POST** `/api/pin/create`

Membuat PIN 6 digit baru untuk user yang sedang login.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "pin": "123456",
  "confirm_pin": "123456"
}
```

**Field Validations:**
- `pin` (required): Harus 6 digit angka (numeric only)
- `confirm_pin` (required): Harus 6 digit angka dan sama dengan PIN
- User belum boleh memiliki PIN sebelumnya

**Success Response (200):**
```json
{
  "message": "PIN created successfully"
}
```

**Error Responses:**
```json
{
  "error": "PIN and confirm PIN do not match"
}
```
```json
{
  "error": "PIN already exists, use change PIN instead"
}
```

---

### 2. Change PIN
**PUT** `/api/pin/change`

Mengganti PIN yang sudah ada dengan PIN baru.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "old_pin": "123456",
  "new_pin": "654321",
  "confirm_pin": "654321"
}
```

**Field Validations:**
- `old_pin` (required): PIN lama yang valid, harus 6 digit angka
- `new_pin` (required): PIN baru, harus 6 digit angka
- `confirm_pin` (required): Konfirmasi PIN baru, harus 6 digit angka
- `new_pin` harus sama dengan `confirm_pin`
- `new_pin` harus berbeda dari `old_pin`
- User harus sudah memiliki PIN sebelumnya

**Success Response (200):**
```json
{
  "message": "PIN changed successfully"
}
```

**Error Responses:**
```json
{
  "error": "old PIN is incorrect"
}
```
```json
{
  "error": "new PIN and confirm PIN do not match"
}
```
```json
{
  "error": "new PIN must be different from old PIN"
}
```
```json
{
  "error": "PIN not set, use create PIN instead"
}
```

---

### 3. Check PIN Status
**GET** `/api/pin/check`

Mengecek apakah user sudah memiliki PIN atau belum.

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200) - Has PIN:**
```json
{
  "has_pin": true
}
```

**Success Response (200) - No PIN:**
```json
{
  "has_pin": false
}
```

---

## 🔒 Authentication

Semua endpoint memerlukan Bearer token yang didapat dari login:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📊 Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request berhasil |
| 400 | Bad Request - Validasi error atau data tidak valid |
| 401 | Unauthorized - Token invalid atau expired |
| 500 | Internal Server Error |

## 🔐 PIN Format

- **Length**: Harus tepat 6 digit
- **Type**: Numeric only (0-9)
- **Examples**: 
  - ✅ Valid: "123456", "000000", "999999", "654321"
  - ❌ Invalid: "12345" (kurang), "1234567" (lebih), "12a456" (ada huruf), "12-456" (ada karakter)

## 🔑 Security Features

- **PIN Hashing**: PIN otomatis di-hash menggunakan bcrypt (sama seperti password)
- **PIN Storage**: PIN tidak pernah disimpan dalam bentuk plain text
- **PIN Verification**: Old PIN harus valid saat change PIN
- **No Duplication**: PIN baru harus berbeda dari PIN lama
- **User Isolation**: Setiap user punya PIN sendiri

## 💡 Use Cases

### Use Case 1: First Time PIN Setup
```bash
# Step 1: Check if user has PIN
GET /api/pin/check
Response: {"has_pin": false}

# Step 2: Create PIN
POST /api/pin/create
{
  "pin": "123456",
  "confirm_pin": "123456"
}
Response: {"message": "PIN created successfully"}
```

### Use Case 2: Change Existing PIN
```bash
# Step 1: Check if user has PIN
GET /api/pin/check
Response: {"has_pin": true}

# Step 2: Change PIN with old PIN verification
PUT /api/pin/change
{
  "old_pin": "123456",
  "new_pin": "654321",
  "confirm_pin": "654321"
}
Response: {"message": "PIN changed successfully"}
```

### Use Case 3: Forgot PIN (Reset Flow)
```bash
# User harus contact admin atau reset via user management
# Admin dapat update user dan set ulang PIN via user management API
# Atau bisa diimplementasikan flow "Forgot PIN" dengan verifikasi email/SMS
```

## 📝 Common Validation Errors

### PIN Length Error
```json
{
  "error": "Key: 'CreatePINRequest.PIN' Error:Field validation for 'PIN' failed on the 'len' tag"
}
```
Solusi: PIN harus tepat 6 digit

### PIN Not Numeric
```json
{
  "error": "Key: 'CreatePINRequest.PIN' Error:Field validation for 'PIN' failed on the 'numeric' tag"
}
```
Solusi: PIN hanya boleh angka 0-9

### PIN Mismatch
```json
{
  "error": "PIN and confirm PIN do not match"
}
```
Solusi: Pastikan PIN dan confirm_pin sama persis

### Wrong Old PIN
```json
{
  "error": "old PIN is incorrect"
}
```
Solusi: Masukkan old_pin yang benar

## 🎯 Best Practices

1. **Client-Side Validation**: Validate PIN format di client sebelum kirim ke server
2. **PIN Masking**: Tampilkan PIN sebagai dots/asterisks di UI (••••••)
3. **PIN Confirmation**: Selalu minta konfirmasi PIN untuk avoid typo
4. **Check Status First**: Gunakan `/pin/check` untuk determine apakah show create atau change form
5. **Secure Input**: Gunakan numeric keyboard di mobile apps
6. **Rate Limiting**: Implement rate limiting untuk prevent brute force attacks

## 🔗 Related APIs

- **Login API**: `/api/auth/login` - Login dengan username/password
- **Change Password API**: `/api/change-password` - Change password (berbeda dari PIN)
- **Profile API**: `/api/profile` - Get user profile

## ⚠️ Important Notes

- PIN berbeda dengan Password
- PIN untuk transaksi atau operasi sensitive
- Password untuk login ke sistem
- User bisa punya Password tanpa PIN
- User harus create PIN manual setelah register
- PIN tidak auto-created saat register
- PIN stored securely dengan bcrypt hashing
- Tidak ada endpoint untuk "view" atau "get" PIN (security reason)

## 🚀 Integration Example

### Frontend Flow Example (React/Vue/Angular):

```javascript
// 1. Check PIN status on app load
const checkPINStatus = async () => {
  const response = await fetch('/api/pin/check', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  const data = await response.json();
  return data.has_pin;
}

// 2. Create PIN
const createPIN = async (pin, confirmPin) => {
  const response = await fetch('/api/pin/create', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      pin: pin,
      confirm_pin: confirmPin
    })
  });
  return await response.json();
}

// 3. Change PIN
const changePIN = async (oldPin, newPin, confirmPin) => {
  const response = await fetch('/api/pin/change', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      old_pin: oldPin,
      new_pin: newPin,
      confirm_pin: confirmPin
    })
  });
  return await response.json();
}
```

## 📱 Mobile App Considerations

1. Use numeric keyboard for PIN input
2. Show/hide PIN toggle optional
3. Biometric as alternative to PIN (store PIN securely in keychain)
4. Auto-lock after inactivity requires PIN
5. Offline PIN verification dengan encrypted local storage
