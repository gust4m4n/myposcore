-- Migration: Rename users.photo to users.image
-- Date: 2025-12-27

-- Rename column from photo to image
ALTER TABLE users RENAME COLUMN photo TO image;

-- Verify the change
SELECT column_name, data_type, character_maximum_length 
FROM information_schema.columns 
WHERE table_name = 'users' AND column_name = 'image';
