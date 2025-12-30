-- Migration: Add new fields to tenants and branches tables
-- Created: 2025-12-30

-- Add new columns to tenants table
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS website VARCHAR(255);
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS phone VARCHAR(50);

-- Add new columns to branches table
ALTER TABLE branches ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE branches ADD COLUMN IF NOT EXISTS website VARCHAR(255);
ALTER TABLE branches ADD COLUMN IF NOT EXISTS email VARCHAR(255);

-- Note: address and phone columns already exist in branches table
