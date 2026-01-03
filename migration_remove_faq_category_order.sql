-- Migration: Remove category and order columns from faqs table
-- Date: 2026-01-03
-- Description: Remove unused category and order fields from FAQ model

-- Remove category column
ALTER TABLE faqs DROP COLUMN IF EXISTS category;

-- Remove order column
ALTER TABLE faqs DROP COLUMN IF EXISTS "order";
