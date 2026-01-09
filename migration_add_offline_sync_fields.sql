-- Migration: Add offline sync support fields
-- Purpose: Enable offline-first mobile app with SQLite local database
-- Date: 2026-01-09

-- Add sync fields to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced' CHECK (sync_status IN ('pending', 'synced', 'conflict', 'failed'));
ALTER TABLE orders ADD COLUMN IF NOT EXISTS client_id VARCHAR(100);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMPTZ;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS version INT DEFAULT 1;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to payments table
ALTER TABLE payments ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced' CHECK (sync_status IN ('pending', 'synced', 'conflict', 'failed'));
ALTER TABLE payments ADD COLUMN IF NOT EXISTS client_id VARCHAR(100);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMPTZ;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS version INT DEFAULT 1;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to order_items table
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced' CHECK (sync_status IN ('pending', 'synced', 'conflict', 'failed'));
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS client_id VARCHAR(100);
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMPTZ;
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS version INT DEFAULT 1;

-- Add indexes for sync queries
CREATE INDEX IF NOT EXISTS idx_orders_sync_status ON orders(sync_status);
CREATE INDEX IF NOT EXISTS idx_orders_client_id ON orders(client_id);
CREATE INDEX IF NOT EXISTS idx_orders_updated_at ON orders(updated_at);
CREATE INDEX IF NOT EXISTS idx_orders_sync_lookup ON orders(tenant_id, branch_id, sync_status, updated_at);

CREATE INDEX IF NOT EXISTS idx_payments_sync_status ON payments(sync_status);
CREATE INDEX IF NOT EXISTS idx_payments_client_id ON payments(client_id);
CREATE INDEX IF NOT EXISTS idx_payments_updated_at ON payments(updated_at);

CREATE INDEX IF NOT EXISTS idx_order_items_sync_status ON order_items(sync_status);

-- Add last_synced timestamp to products for delta sync
ALTER TABLE products ADD COLUMN IF NOT EXISTS last_synced_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_products_last_synced ON products(tenant_id, last_synced_at);

-- Add last_synced timestamp to categories for delta sync
ALTER TABLE categories ADD COLUMN IF NOT EXISTS last_synced_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_categories_last_synced ON categories(tenant_id, last_synced_at);

-- Create sync log table to track all sync operations
CREATE TABLE IF NOT EXISTS sync_logs (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id INT REFERENCES branches(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    client_id VARCHAR(100) NOT NULL,
    sync_type VARCHAR(50) NOT NULL, -- 'upload', 'download', 'full', 'delta'
    entity_type VARCHAR(50), -- 'orders', 'payments', 'products', etc.
    records_uploaded INT DEFAULT 0,
    records_downloaded INT DEFAULT 0,
    conflicts_detected INT DEFAULT 0,
    status VARCHAR(20) NOT NULL CHECK (status IN ('started', 'completed', 'failed')),
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    duration_ms INT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sync_logs_tenant ON sync_logs(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_logs_client ON sync_logs(client_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_logs_status ON sync_logs(status);

-- Create conflicts table to store unresolved conflicts
CREATE TABLE IF NOT EXISTS sync_conflicts (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id INT REFERENCES branches(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    client_id VARCHAR(100) NOT NULL,
    client_version INT NOT NULL,
    server_version INT NOT NULL,
    client_data JSONB NOT NULL,
    server_data JSONB NOT NULL,
    conflict_type VARCHAR(50) NOT NULL, -- 'update_conflict', 'delete_conflict'
    resolution_strategy VARCHAR(50), -- 'server_wins', 'client_wins', 'manual', 'merge'
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMPTZ,
    resolved_by INT REFERENCES users(id) ON DELETE SET NULL,
    resolved_data JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_conflicts_unresolved ON sync_conflicts(tenant_id, resolved) WHERE resolved = FALSE;
CREATE INDEX IF NOT EXISTS idx_conflicts_entity ON sync_conflicts(entity_type, entity_id);

-- Update existing records to have synced status
UPDATE orders SET sync_status = 'synced' WHERE sync_status IS NULL;
UPDATE payments SET sync_status = 'synced' WHERE sync_status IS NULL;
UPDATE order_items SET sync_status = 'synced' WHERE sync_status IS NULL;

-- Comments for documentation
COMMENT ON COLUMN orders.sync_status IS 'Sync status: pending (waiting to sync), synced (successfully synced), conflict (merge conflict), failed (sync error)';
COMMENT ON COLUMN orders.client_id IS 'Unique identifier for client device (mobile app instance)';
COMMENT ON COLUMN orders.local_timestamp IS 'Timestamp when record was created on client device';
COMMENT ON COLUMN orders.version IS 'Version number for optimistic locking and conflict detection';
COMMENT ON COLUMN orders.conflict_data IS 'Stores conflicting data for manual resolution';

COMMENT ON TABLE sync_logs IS 'Tracks all sync operations for monitoring and debugging';
COMMENT ON TABLE sync_conflicts IS 'Stores unresolved sync conflicts for manual resolution';
