-- Migration: Rename role 'user' to 'staff'
-- Date: 2026-01-02
-- Description: Update all user roles from 'user' to 'staff'

-- Update existing users with role 'user' to 'staff'
UPDATE users SET role = 'staff' WHERE role = 'user';
