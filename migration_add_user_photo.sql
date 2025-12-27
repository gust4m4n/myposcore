-- Migration: Add photo field to users table
-- Date: 2025-12-27
-- Description: Add photo column to store user profile image URL

-- Add photo column if not exists
ALTER TABLE users ADD COLUMN IF NOT EXISTS photo VARCHAR(500);

-- Add index for faster queries (optional)
CREATE INDEX IF NOT EXISTS idx_users_photo ON users(photo) WHERE photo IS NOT NULL AND photo != '';

-- Verify column added
SELECT column_name, data_type, character_maximum_length 
FROM information_schema.columns 
WHERE table_name = 'users' AND column_name = 'photo';
