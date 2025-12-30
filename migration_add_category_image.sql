-- Migration: Add image column to categories table
-- This allows categories to have associated images for better visual representation

ALTER TABLE categories ADD COLUMN IF NOT EXISTS image VARCHAR(500);

-- Add comment for documentation
COMMENT ON COLUMN categories.image IS 'Path to category image file (stored in uploads/categories/)';
