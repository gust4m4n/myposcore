-- Migration: Add Offline Sync Fields to All Entities
-- Date: 2026-01-09
-- Description: Add sync_status, client_id, local_timestamp, version, and conflict_data 
--              to tenants, branches, users, and audit_trails tables

-- Add sync fields to tenants table
ALTER TABLE tenants
ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced',
ADD COLUMN IF NOT EXISTS client_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to branches table
ALTER TABLE branches
ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced',
ADD COLUMN IF NOT EXISTS client_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced',
ADD COLUMN IF NOT EXISTS client_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to products table (if not already added)
ALTER TABLE products
ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced',
ADD COLUMN IF NOT EXISTS client_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to categories table (if not already added)
ALTER TABLE categories
ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced',
ADD COLUMN IF NOT EXISTS client_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Add sync fields to audit_trails table
ALTER TABLE audit_trails
ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) DEFAULT 'synced',
ADD COLUMN IF NOT EXISTS client_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS local_timestamp TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conflict_data JSONB;

-- Create indexes for sync operations on tenants
CREATE INDEX IF NOT EXISTS idx_tenants_sync_status ON tenants(sync_status);
CREATE INDEX IF NOT EXISTS idx_tenants_client_id ON tenants(client_id);
CREATE INDEX IF NOT EXISTS idx_tenants_local_timestamp ON tenants(local_timestamp);

-- Create indexes for sync operations on branches
CREATE INDEX IF NOT EXISTS idx_branches_sync_status ON branches(sync_status);
CREATE INDEX IF NOT EXISTS idx_branches_client_id ON branches(client_id);
CREATE INDEX IF NOT EXISTS idx_branches_local_timestamp ON branches(local_timestamp);

-- Create indexes for sync operations on users
CREATE INDEX IF NOT EXISTS idx_users_sync_status ON users(sync_status);
CREATE INDEX IF NOT EXISTS idx_users_client_id ON users(client_id);
CREATE INDEX IF NOT EXISTS idx_users_local_timestamp ON users(local_timestamp);

-- Create indexes for sync operations on products
CREATE INDEX IF NOT EXISTS idx_products_sync_status ON products(sync_status);
CREATE INDEX IF NOT EXISTS idx_products_client_id ON products(client_id);
CREATE INDEX IF NOT EXISTS idx_products_local_timestamp ON products(local_timestamp);

-- Create indexes for sync operations on categories
CREATE INDEX IF NOT EXISTS idx_categories_sync_status ON categories(sync_status);
CREATE INDEX IF NOT EXISTS idx_categories_client_id ON categories(client_id);
CREATE INDEX IF NOT EXISTS idx_categories_local_timestamp ON categories(local_timestamp);

-- Create indexes for sync operations on audit_trails
CREATE INDEX IF NOT EXISTS idx_audit_trails_sync_status ON audit_trails(sync_status);
CREATE INDEX IF NOT EXISTS idx_audit_trails_client_id ON audit_trails(client_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_local_timestamp ON audit_trails(local_timestamp);

-- Add comment for documentation
COMMENT ON COLUMN tenants.sync_status IS 'Sync status: pending, synced, conflict';
COMMENT ON COLUMN branches.sync_status IS 'Sync status: pending, synced, conflict';
COMMENT ON COLUMN users.sync_status IS 'Sync status: pending, synced, conflict';
COMMENT ON COLUMN products.sync_status IS 'Sync status: pending, synced, conflict';
COMMENT ON COLUMN categories.sync_status IS 'Sync status: pending, synced, conflict';
COMMENT ON COLUMN audit_trails.sync_status IS 'Sync status: pending, synced, conflict';
