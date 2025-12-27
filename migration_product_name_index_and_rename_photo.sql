-- Migration: Add index to product name and rename photo to image
-- Date: 2025-12-27
-- Description: Add index to name field for faster search and rename photo column to image

-- Add index to product name for faster search
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);

-- Rename photo column to image
ALTER TABLE products RENAME COLUMN photo TO image;

-- Verify changes
SELECT 
    tablename,
    indexname,
    indexdef
FROM 
    pg_indexes
WHERE 
    tablename = 'products'
    AND indexname = 'idx_products_name';

SELECT column_name, data_type, character_maximum_length 
FROM information_schema.columns 
WHERE table_name = 'products' AND column_name = 'image';
