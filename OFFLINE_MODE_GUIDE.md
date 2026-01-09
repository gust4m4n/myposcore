# Offline Mode Implementation Guide

## Overview
Implementasi offline mode untuk MyPOSCore mobile app menggunakan SQLite local database dengan sync mechanism ke PostgreSQL server.

## Architecture

### Client Side (Mobile App)
- **Local Database**: SQLite
- **Sync Strategy**: Offline-first approach
- **Data Flow**: 
  1. Semua transaksi disimpan ke SQLite lokal terlebih dahulu
  2. Mark as `pending_sync`
  3. Ketika online, upload ke server
  4. Server response dengan mapping `local_id → server_id`
  5. Update status menjadi `synced`

### Server Side (Backend)
- **Main Database**: PostgreSQL
- **Sync Fields**: Added to orders, payments, order_items tables
- **Conflict Resolution**: Version-based optimistic locking
- **Sync Logs**: Complete audit trail untuk monitoring

---

## Database Schema Changes

### Migration File
File: `migration_add_offline_sync_fields.sql`

**New Fields Added:**
- `sync_status` - Status sync: pending, synced, conflict, failed
- `client_id` - Unique identifier untuk device
- `local_timestamp` - Timestamp dari client device
- `version` - Version number untuk conflict detection
- `conflict_data` - JSON data untuk manual resolution

**New Tables:**
- `sync_logs` - Track semua sync operations
- `sync_conflicts` - Store unresolved conflicts

**Indexes Added:**
- Performance indexes untuk sync queries
- Composite indexes untuk tenant+branch filtering

### Running Migration
```bash
PGPASSWORD=postgres psql -h localhost -U postgres -d myposcore -f migration_add_offline_sync_fields.sql
```

---

## API Endpoints

### 1. Upload Data dari Client
**POST** `/api/v1/sync/upload`

Upload orders dan payments dari mobile app ke server.

**Request Body:**
```json
{
  "client_id": "device_uuid_12345",
  "client_timestamp": "2026-01-09T10:30:00Z",
  "orders": [
    {
      "local_id": "order_local_uuid_1",
      "order_number": "ORD-20260109-001",
      "total_amount": 150000,
      "status": "completed",
      "notes": "Cash payment",
      "items": [
        {
          "product_id": 27,
          "quantity": 2,
          "price": 50000,
          "subtotal": 100000
        },
        {
          "product_id": 42,
          "quantity": 1,
          "price": 50000,
          "subtotal": 50000
        }
      ],
      "local_timestamp": "2026-01-09T10:25:00Z",
      "version": 1
    }
  ],
  "payments": [
    {
      "local_id": "payment_local_uuid_1",
      "order_local_id": "order_local_uuid_1",
      "amount": 150000,
      "payment_method": "cash",
      "status": "completed",
      "notes": "Cash received",
      "local_timestamp": "2026-01-09T10:26:00Z",
      "version": 1
    }
  ],
  "last_sync_at": "2026-01-09T09:00:00Z"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "Data synced successfully",
  "data": {
    "sync_id": "sync_123",
    "processed_orders": 1,
    "processed_payments": 1,
    "failed_orders": 0,
    "failed_payments": 0,
    "conflicts": [],
    "order_mapping": {
      "order_local_uuid_1": 456
    },
    "payment_mapping": {
      "payment_local_uuid_1": 789
    },
    "sync_timestamp": "2026-01-09T10:30:15Z",
    "errors": []
  }
}
```

---

### 2. Download Master Data
**POST** `/api/v1/sync/download`

Download products dan categories ke mobile app.

**Request Body:**
```json
{
  "client_id": "device_uuid_12345",
  "last_sync_at": "2026-01-09T09:00:00Z",
  "entity_types": ["products", "categories"]
}
```

**Response:**
```json
{
  "code": 0,
  "message": "Data downloaded successfully",
  "data": {
    "sync_id": "sync_124",
    "products": [
      {
        "id": 27,
        "name": "Nasi Goreng",
        "price": 50000,
        "stock": 100,
        "category_id": 13,
        "category_detail": {
          "id": 13,
          "name": "Makanan Utama"
        },
        "updated_at": "2026-01-09T10:00:00Z"
      }
    ],
    "categories": [
      {
        "id": 13,
        "name": "Makanan Utama",
        "description": "Kategori untuk makanan utama",
        "updated_at": "2026-01-08T15:00:00Z"
      }
    ],
    "sync_timestamp": "2026-01-09T10:30:20Z",
    "has_more": false
  }
}
```

---

### 3. Check Sync Status
**GET** `/api/v1/sync/status?client_id=device_uuid_12345`

Get current sync status untuk device.

**Response:**
```json
{
  "code": 0,
  "message": "Sync status retrieved successfully",
  "data": {
    "client_id": "device_uuid_12345",
    "last_sync_at": "2026-01-09T10:30:15Z",
    "pending_uploads": 0,
    "pending_conflicts": 0,
    "total_syncs": 25,
    "last_sync_success": true,
    "server_timestamp": "2026-01-09T10:35:00Z"
  }
}
```

---

### 4. Get Sync History
**GET** `/api/v1/sync/logs?client_id=device_uuid_12345&page=1&page_size=20`

Get paginated sync history logs.

**Response:**
```json
{
  "code": 0,
  "message": "Sync logs retrieved successfully",
  "data": {
    "page": 1,
    "page_size": 20,
    "total_items": 25,
    "total_pages": 2,
    "data": [
      {
        "id": 123,
        "client_id": "device_uuid_12345",
        "sync_type": "upload",
        "entity_type": "orders",
        "records_uploaded": 5,
        "records_downloaded": 0,
        "conflicts_detected": 0,
        "status": "completed",
        "duration_ms": 1250,
        "created_at": "2026-01-09T10:30:15Z",
        "completed_at": "2026-01-09T10:30:16Z"
      }
    ]
  }
}
```

---

### 5. Resolve Conflict
**POST** `/api/v1/sync/conflicts/resolve`

Manually resolve sync conflict.

**Request Body:**
```json
{
  "conflict_id": 45,
  "resolution_strategy": "server_wins",
  "resolved_data": null
}
```

Strategies:
- `server_wins` - Use server data
- `client_wins` - Use client data
- `manual` - Use custom merged data (provide in `resolved_data`)

---

### 6. Get Server Time
**GET** `/api/v1/sync/time`

Get server timestamp untuk time synchronization.

**Response:**
```json
{
  "code": 0,
  "message": "Server time retrieved successfully",
  "data": {
    "server_time": "2026-01-09T10:35:45.123Z",
    "unix_time": 1736420145
  }
}
```

---

## Mobile App Implementation Guide

### 1. SQLite Schema (Client Side)

```sql
-- Orders table
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    local_id TEXT UNIQUE NOT NULL,  -- UUID generated by app
    server_id INTEGER,               -- NULL until synced
    order_number TEXT,
    total_amount REAL NOT NULL,
    status TEXT DEFAULT 'pending',
    notes TEXT,
    sync_status TEXT DEFAULT 'pending', -- pending, synced, failed
    local_timestamp TEXT NOT NULL,
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- Payments table
CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    local_id TEXT UNIQUE NOT NULL,
    server_id INTEGER,
    order_local_id TEXT NOT NULL,   -- Reference to orders.local_id
    amount REAL NOT NULL,
    payment_method TEXT NOT NULL,
    status TEXT DEFAULT 'completed',
    notes TEXT,
    sync_status TEXT DEFAULT 'pending',
    local_timestamp TEXT NOT NULL,
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (order_local_id) REFERENCES orders(local_id)
);

-- Order items table
CREATE TABLE order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_local_id TEXT NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price REAL NOT NULL,
    subtotal REAL NOT NULL,
    sync_status TEXT DEFAULT 'pending',
    FOREIGN KEY (order_local_id) REFERENCES orders(local_id)
);

-- Products (master data - synced from server)
CREATE TABLE products (
    id INTEGER PRIMARY KEY,  -- Server ID
    name TEXT NOT NULL,
    description TEXT,
    sku TEXT,
    price REAL NOT NULL,
    stock INTEGER NOT NULL,
    category_id INTEGER,
    category_name TEXT,
    image TEXT,
    is_active INTEGER DEFAULT 1,
    last_synced_at TEXT NOT NULL
);

-- Categories (master data - synced from server)
CREATE TABLE categories (
    id INTEGER PRIMARY KEY,  -- Server ID
    name TEXT NOT NULL,
    description TEXT,
    image TEXT,
    is_active INTEGER DEFAULT 1,
    last_synced_at TEXT NOT NULL
);

-- Indexes
CREATE INDEX idx_orders_sync_status ON orders(sync_status);
CREATE INDEX idx_payments_sync_status ON payments(sync_status);
CREATE INDEX idx_products_last_synced ON products(last_synced_at);
```

---

### 2. Sync Logic (Pseudo Code)

#### Create Order Offline
```javascript
async function createOrder(orderData) {
  const db = await getDatabase();
  
  // 1. Generate local UUID
  const localId = generateUUID();
  const localTimestamp = new Date().toISOString();
  
  // 2. Save to local SQLite
  await db.execute(`
    INSERT INTO orders (local_id, total_amount, status, notes, 
                        sync_status, local_timestamp, version, 
                        created_at, updated_at)
    VALUES (?, ?, ?, ?, 'pending', ?, 1, ?, ?)
  `, [localId, orderData.total_amount, 'completed', 
      orderData.notes, localTimestamp, localTimestamp, localTimestamp]);
  
  // 3. Save order items
  for (const item of orderData.items) {
    await db.execute(`
      INSERT INTO order_items (order_local_id, product_id, quantity, price, subtotal)
      VALUES (?, ?, ?, ?, ?)
    `, [localId, item.product_id, item.quantity, item.price, item.subtotal]);
    
    // 4. Update local product stock
    await db.execute(`
      UPDATE products SET stock = stock - ? WHERE id = ?
    `, [item.quantity, item.product_id]);
  }
  
  // 5. Try to sync immediately if online
  if (await isOnline()) {
    await syncPendingOrders();
  }
  
  return { localId, status: 'saved' };
}
```

#### Sync Pending Orders
```javascript
async function syncPendingOrders() {
  const db = await getDatabase();
  
  // 1. Get pending orders
  const orders = await db.query(`
    SELECT * FROM orders WHERE sync_status = 'pending'
  `);
  
  if (orders.length === 0) return;
  
  // 2. Get order items for each order
  const ordersToSync = [];
  for (const order of orders) {
    const items = await db.query(`
      SELECT * FROM order_items WHERE order_local_id = ?
    `, [order.local_id]);
    
    ordersToSync.push({
      local_id: order.local_id,
      order_number: order.order_number,
      total_amount: order.total_amount,
      status: order.status,
      notes: order.notes,
      items: items.map(item => ({
        product_id: item.product_id,
        quantity: item.quantity,
        price: item.price,
        subtotal: item.subtotal
      })),
      local_timestamp: order.local_timestamp,
      version: order.version
    });
  }
  
  // 3. Get pending payments
  const payments = await db.query(`
    SELECT * FROM payments WHERE sync_status = 'pending'
  `);
  
  const paymentsToSync = payments.map(p => ({
    local_id: p.local_id,
    order_local_id: p.order_local_id,
    amount: p.amount,
    payment_method: p.payment_method,
    status: p.status,
    notes: p.notes,
    local_timestamp: p.local_timestamp,
    version: p.version
  }));
  
  // 4. Upload to server
  try {
    const response = await fetch('http://server/api/v1/sync/upload', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + getAuthToken(),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        client_id: getDeviceId(),
        client_timestamp: new Date().toISOString(),
        orders: ordersToSync,
        payments: paymentsToSync,
        last_sync_at: await getLastSyncTime()
      })
    });
    
    const result = await response.json();
    
    // 5. Update local database with server IDs
    for (const [localId, serverId] of Object.entries(result.data.order_mapping)) {
      await db.execute(`
        UPDATE orders 
        SET server_id = ?, sync_status = 'synced' 
        WHERE local_id = ?
      `, [serverId, localId]);
    }
    
    for (const [localId, serverId] of Object.entries(result.data.payment_mapping)) {
      await db.execute(`
        UPDATE payments 
        SET server_id = ?, sync_status = 'synced' 
        WHERE local_id = ?
      `, [serverId, localId]);
    }
    
    // 6. Update last sync time
    await setLastSyncTime(result.data.sync_timestamp);
    
    console.log('Sync completed:', result);
    
  } catch (error) {
    console.error('Sync failed:', error);
    // Mark as failed
    for (const order of orders) {
      await db.execute(`
        UPDATE orders SET sync_status = 'failed' WHERE local_id = ?
      `, [order.local_id]);
    }
  }
}
```

#### Download Master Data
```javascript
async function downloadMasterData() {
  const db = await getDatabase();
  
  try {
    const response = await fetch('http://server/api/v1/sync/download', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + getAuthToken(),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        client_id: getDeviceId(),
        last_sync_at: await getLastMasterDataSync(),
        entity_types: ['products', 'categories']
      })
    });
    
    const result = await response.json();
    
    // Update products
    for (const product of result.data.products) {
      await db.execute(`
        INSERT OR REPLACE INTO products 
        (id, name, description, sku, price, stock, category_id, 
         category_name, image, is_active, last_synced_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
      `, [
        product.id,
        product.name,
        product.description,
        product.sku,
        product.price,
        product.stock,
        product.category_id,
        product.category_detail?.name,
        product.image,
        product.is_active ? 1 : 0,
        result.data.sync_timestamp
      ]);
    }
    
    // Update categories
    for (const category of result.data.categories) {
      await db.execute(`
        INSERT OR REPLACE INTO categories 
        (id, name, description, image, is_active, last_synced_at)
        VALUES (?, ?, ?, ?, ?, ?)
      `, [
        category.id,
        category.name,
        category.description,
        category.image,
        category.is_active ? 1 : 0,
        result.data.sync_timestamp
      ]);
    }
    
    await setLastMasterDataSync(result.data.sync_timestamp);
    
    console.log('Master data downloaded:', result);
    
  } catch (error) {
    console.error('Download failed:', error);
  }
}
```

#### Background Sync
```javascript
// Setup periodic sync
setInterval(async () => {
  if (await isOnline()) {
    await syncPendingOrders();
  }
}, 5 * 60 * 1000); // Every 5 minutes

// Listen to online event
window.addEventListener('online', async () => {
  console.log('Back online - starting sync');
  await syncPendingOrders();
  await downloadMasterData();
});

// Initial download when app starts
async function initializeApp() {
  const db = await getDatabase();
  
  if (await isOnline()) {
    // Download latest master data
    await downloadMasterData();
  }
  
  // Check for pending syncs
  const pendingCount = await db.query(`
    SELECT COUNT(*) as count FROM orders WHERE sync_status = 'pending'
  `);
  
  if (pendingCount[0].count > 0 && await isOnline()) {
    await syncPendingOrders();
  }
}
```

---

## Testing Scenarios

### 1. Create Order Offline
```bash
# Mobile app creates order while offline
# Order saved to SQLite with sync_status='pending'
# When online, sync triggers automatically
```

### 2. Sync Multiple Orders
```bash
# Create 5 orders offline
# Go online
# All 5 orders upload in single batch
# Receive mapping of local_id → server_id
```

### 3. Conflict Resolution
```bash
# Device A creates order, sync to server (version=1)
# Device B fetches same order, modifies offline
# Device B tries to sync (version conflict detected)
# Server returns conflict info
# Manual resolution required
```

### 4. Delta Sync
```bash
# First sync: Download all products (1000 items)
# Store last_sync_at timestamp
# Next sync: Only download updated products since last_sync_at
# Much faster, less data transfer
```

---

## Monitoring & Debugging

### Check Sync Status
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/sync/status?client_id=device_123"
```

### View Sync Logs
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/sync/logs?client_id=device_123&page=1"
```

### Check Database
```sql
-- Check pending orders
SELECT * FROM orders WHERE sync_status = 'pending';

-- Check sync logs
SELECT * FROM sync_logs ORDER BY created_at DESC LIMIT 10;

-- Check unresolved conflicts
SELECT * FROM sync_conflicts WHERE resolved = false;
```

---

## Performance Optimization

### 1. Batch Sync
- Upload multiple orders in single request
- Reduces HTTP overhead
- Faster sync completion

### 2. Delta Sync
- Only sync changed data
- Use `last_sync_at` timestamp
- Significantly reduces data transfer

### 3. Compression
- Consider gzip compression for large payloads
- Especially useful for master data download

### 4. Background Sync
- Use background services/workers
- Don't block UI
- Retry failed syncs automatically

---

## Security Considerations

1. **Authentication**: All sync endpoints require JWT token
2. **Tenant Isolation**: Automatic filtering by tenant_id from token
3. **Data Validation**: Server validates all uploaded data
4. **Conflict Detection**: Version-based optimistic locking
5. **Audit Trail**: All sync operations logged

---

## Future Enhancements

1. **Real-time Sync**: WebSocket untuk instant sync
2. **Partial Sync**: Sync specific orders by ID
3. **Image Sync**: Upload product/profile images
4. **Sync Queue**: Priority queue untuk critical data
5. **Conflict Auto-Resolution**: Intelligent merge strategies
6. **Encryption**: Encrypt sensitive data in local SQLite

---

## Support

Untuk pertanyaan atau issues terkait offline mode:
- Check sync logs: `/api/v1/sync/logs`
- Check sync status: `/api/v1/sync/status`
- Review migration file untuk schema details
- Test dengan Postman collection (sync endpoints)
