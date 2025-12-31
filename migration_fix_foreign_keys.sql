-- ============================================================================
-- Migration: Fix incorrect foreign key constraints on users table
-- ============================================================================
-- Description: 
--   Remove incorrect foreign key constraints where users.created_by/updated_by/deleted_by
--   incorrectly reference other tables instead of users table
-- 
-- Issue: GORM auto-migration created wrong foreign keys
-- ============================================================================

BEGIN;

-- Drop all incorrect foreign key constraints from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_branches_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_branches_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_branches_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_categories_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_categories_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_categories_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_faqs_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_faqs_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_faqs_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_orders_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_orders_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_orders_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_payments_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_payments_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_payments_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_products_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_products_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_products_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_tenants_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_tenants_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_tenants_deleter;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_terms_and_conditions_creator;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_terms_and_conditions_updater;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_terms_and_conditions_deleter;

-- Keep only the correct self-referencing foreign keys
-- (fk_users_creator, fk_users_updater, fk_users_deleter should already exist correctly)

COMMIT;

-- ============================================================================
-- Verification Query (run after migration):
-- SELECT conname, conrelid::regclass, confrelid::regclass 
-- FROM pg_constraint 
-- WHERE conrelid = 'users'::regclass AND contype = 'f';
-- ============================================================================
