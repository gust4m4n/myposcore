-- Insert dummy categories for TENANT001
INSERT INTO categories (tenant_id, name, description, is_active, created_at, updated_at) VALUES
(1, 'Food', 'All food items including main dishes, appetizers, and desserts', true, NOW(), NOW()),
(1, 'Beverage', 'All kinds of drinks including hot and cold beverages', true, NOW(), NOW()),
(1, 'Snacks', 'Light snacks and finger foods', true, NOW(), NOW()),
(1, 'Desserts', 'Sweet desserts and pastries', true, NOW(), NOW()),
(1, 'Main Course', 'Main dishes and entrees', true, NOW(), NOW());

-- Insert dummy categories for TENANT002
INSERT INTO categories (tenant_id, name, description, is_active, created_at, updated_at) VALUES
(2, 'Electronics', 'Electronic devices and accessories', true, NOW(), NOW()),
(2, 'Furniture', 'Office and home furniture', true, NOW(), NOW()),
(2, 'Stationery', 'Office supplies and stationery items', true, NOW(), NOW());

-- Insert dummy categories for DEMO_TENANT
INSERT INTO categories (tenant_id, name, description, is_active, created_at, updated_at) VALUES
(3, 'Coffee', 'All coffee based beverages', true, NOW(), NOW()),
(3, 'Tea', 'All tea varieties', true, NOW(), NOW()),
(3, 'Pastries', 'Bread, cakes, and pastries', true, NOW(), NOW()),
(3, 'Breakfast', 'Breakfast menu items', true, NOW(), NOW());
