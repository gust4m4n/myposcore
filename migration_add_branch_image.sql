-- Migration: Add image column to branches table
-- Description: Adds image field to store branch logo/image URL

ALTER TABLE branches 
ADD COLUMN IF NOT EXISTS image VARCHAR(500);

-- Add comment to the column
COMMENT ON COLUMN branches.image IS 'URL path to branch logo/image file';
