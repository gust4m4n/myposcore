-- Migration: Add image column to tenants table
-- Description: Adds image field to store tenant logo/image URL

ALTER TABLE tenants 
ADD COLUMN IF NOT EXISTS image VARCHAR(500);

-- Add comment to the column
COMMENT ON COLUMN tenants.image IS 'URL path to tenant logo/image file';
