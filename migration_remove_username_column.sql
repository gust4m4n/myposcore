-- ============================================================================
-- Migration: Remove username column from users table
-- ============================================================================
-- Description: 
--   Removes the username column as the system now uses email for authentication
--   All login and authentication logic has been updated to use email instead
-- 
-- Prerequisites:
--   - All application code must be updated to not use username field
--   - All users must have valid email addresses
-- 
-- Impact:
--   - Drops username column from users table
--   - Removes any indexes or constraints related to username
-- ============================================================================

BEGIN;

-- Check if there are any users without email (safety check)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM users WHERE email IS NULL OR email = '') THEN
        RAISE EXCEPTION 'Cannot drop username column: Some users have no email address';
    END IF;
END $$;

-- Drop the username column
ALTER TABLE users DROP COLUMN IF EXISTS username;

-- Log the change
DO $$
BEGIN
    RAISE NOTICE 'Successfully dropped username column from users table';
END $$;

COMMIT;

-- ============================================================================
-- Verification Query (run after migration):
-- SELECT column_name FROM information_schema.columns 
-- WHERE table_name = 'users' AND column_name = 'username';
-- 
-- Should return 0 rows if migration successful
-- ============================================================================
