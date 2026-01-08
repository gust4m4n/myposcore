-- Migration: Add category_id to products table and create foreign key relationship
-- This migration adds a proper foreign key relationship between products and categories
-- PostgreSQL syntax

-- Step 1: Add category_id column to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER NULL;

-- Step 2: Add index for category_id
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);

-- Step 3: Add foreign key constraint
ALTER TABLE products ADD CONSTRAINT fk_products_category_id 
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL ON UPDATE CASCADE;

-- Step 4: Migrate existing data from category string to category_id
-- This will match product category names with existing category names in the categories table
UPDATE products p
SET category_id = c.id
FROM categories c
WHERE LOWER(TRIM(p.category)) = LOWER(TRIM(c.name)) 
  AND p.tenant_id = c.tenant_id
  AND p.category IS NOT NULL 
  AND p.category != '';

-- Step 5 (Optional): After verifying data migration, you can drop the old category column
-- Uncomment the line below only after verifying all data is properly migrated
-- ALTER TABLE products DROP COLUMN category;
