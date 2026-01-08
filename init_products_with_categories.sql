-- Demo Products with Category Relationships
-- This script inserts sample products that properly reference categories by ID

-- Products for TENANT001 (FoodCorner)
-- Category IDs: 1=Food, 2=Beverage, 3=Snacks, 4=Desserts, 5=Main Course

-- Main Course (category_id: 5)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(1, 5, 'Nasi Goreng Spesial', 'Fried rice with chicken, vegetables, and fried egg', 'FC-MGC-001', 25000.00, 50, true, NOW(), NOW()),
(1, 5, 'Mie Goreng Seafood', 'Stir-fried noodles with seafood and vegetables', 'FC-MGS-002', 28000.00, 40, true, NOW(), NOW()),
(1, 5, 'Ayam Bakar Taliwang', 'Grilled chicken with special spicy sauce', 'FC-ABT-003', 32000.00, 30, true, NOW(), NOW()),
(1, 5, 'Soto Ayam', 'Traditional chicken soup with rice', 'FC-STA-004', 22000.00, 45, true, NOW(), NOW()),
(1, 5, 'Gado-Gado', 'Indonesian salad with peanut sauce', 'FC-GDG-005', 20000.00, 35, true, NOW(), NOW());

-- Beverage (category_id: 2)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(1, 2, 'Es Teh Manis', 'Sweet iced tea', 'FC-ETM-006', 5000.00, 100, true, NOW(), NOW()),
(1, 2, 'Es Jeruk', 'Fresh orange juice', 'FC-EJR-007', 8000.00, 80, true, NOW(), NOW()),
(1, 2, 'Kopi Hitam', 'Black coffee', 'FC-KPH-008', 7000.00, 90, true, NOW(), NOW()),
(1, 2, 'Cappuccino', 'Coffee with steamed milk', 'FC-CAP-009', 15000.00, 60, true, NOW(), NOW()),
(1, 2, 'Jus Alpukat', 'Avocado smoothie', 'FC-JAV-010', 12000.00, 40, true, NOW(), NOW());

-- Snacks (category_id: 3)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(1, 3, 'Tahu Crispy', 'Crispy fried tofu', 'FC-THC-011', 10000.00, 50, true, NOW(), NOW()),
(1, 3, 'Pisang Goreng', 'Fried banana fritters', 'FC-PNG-012', 8000.00, 60, true, NOW(), NOW()),
(1, 3, 'Lumpia Semarang', 'Spring rolls with vegetables', 'FC-LPS-013', 12000.00, 45, true, NOW(), NOW()),
(1, 3, 'French Fries', 'Crispy french fries', 'FC-FRF-014', 15000.00, 70, true, NOW(), NOW());

-- Desserts (category_id: 4)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(1, 4, 'Es Campur', 'Mixed ice dessert', 'FC-ECM-015', 12000.00, 40, true, NOW(), NOW()),
(1, 4, 'Klepon', 'Sweet rice cake balls', 'FC-KLP-016', 10000.00, 50, true, NOW(), NOW()),
(1, 4, 'Puding Coklat', 'Chocolate pudding', 'FC-PDC-017', 8000.00, 55, true, NOW(), NOW());

-- Products for TENANT002 (RetailStore)
-- Category IDs: 6=Electronics, 7=Furniture, 8=Stationery

-- Electronics (category_id: 6)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(2, 6, 'Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 'RS-WMS-001', 150000.00, 25, true, NOW(), NOW()),
(2, 6, 'USB Flash Drive 32GB', 'High-speed USB 3.0 flash drive', 'RS-UFD-002', 80000.00, 40, true, NOW(), NOW()),
(2, 6, 'Keyboard Mechanical', 'RGB mechanical gaming keyboard', 'RS-KBM-003', 450000.00, 15, true, NOW(), NOW());

-- Furniture (category_id: 7)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(2, 7, 'Office Chair', 'Ergonomic office chair with lumbar support', 'RS-OCH-004', 1200000.00, 10, true, NOW(), NOW()),
(2, 7, 'Desk Lamp LED', 'Adjustable LED desk lamp', 'RS-DLM-005', 250000.00, 20, true, NOW(), NOW()),
(2, 7, 'Bookshelf', 'Wooden bookshelf 5 tiers', 'RS-BSF-006', 800000.00, 8, true, NOW(), NOW());

-- Stationery (category_id: 8)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(2, 8, 'Ballpoint Pen Pack', 'Pack of 12 blue ballpoint pens', 'RS-BPP-007', 25000.00, 100, true, NOW(), NOW()),
(2, 8, 'Notebook A5', 'Ruled notebook 100 pages', 'RS-NBA-008', 15000.00, 80, true, NOW(), NOW()),
(2, 8, 'Stapler', 'Heavy duty stapler', 'RS-STP-009', 35000.00, 50, true, NOW(), NOW());

-- Products for DEMO_TENANT (Café)
-- Category IDs: 9=Coffee, 10=Tea, 11=Pastries, 12=Breakfast

-- Coffee (category_id: 9)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(3, 9, 'Espresso', 'Strong Italian coffee', 'CF-ESP-001', 18000.00, 100, true, NOW(), NOW()),
(3, 9, 'Americano', 'Espresso with hot water', 'CF-AMR-002', 20000.00, 100, true, NOW(), NOW()),
(3, 9, 'Cafe Latte', 'Espresso with steamed milk', 'CF-LAT-003', 25000.00, 80, true, NOW(), NOW()),
(3, 9, 'Cappuccino', 'Espresso with foamed milk', 'CF-CAP-004', 25000.00, 80, true, NOW(), NOW()),
(3, 9, 'Mocha', 'Chocolate flavored coffee', 'CF-MCH-005', 28000.00, 70, true, NOW(), NOW()),
(3, 9, 'Caramel Macchiato', 'Coffee with caramel syrup', 'CF-CAM-006', 30000.00, 60, true, NOW(), NOW());

-- Tea (category_id: 10)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(3, 10, 'Green Tea', 'Japanese green tea', 'CF-GRT-007', 15000.00, 100, true, NOW(), NOW()),
(3, 10, 'Earl Grey', 'Black tea with bergamot', 'CF-EGR-008', 18000.00, 90, true, NOW(), NOW()),
(3, 10, 'Chamomile Tea', 'Relaxing herbal tea', 'CF-CHM-009', 18000.00, 85, true, NOW(), NOW()),
(3, 10, 'Thai Tea', 'Sweet Thai iced tea', 'CF-THT-010', 20000.00, 75, true, NOW(), NOW());

-- Pastries (category_id: 11)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(3, 11, 'Croissant', 'Butter croissant', 'CF-CRS-011', 22000.00, 40, true, NOW(), NOW()),
(3, 11, 'Chocolate Muffin', 'Chocolate chip muffin', 'CF-CMF-012', 18000.00, 50, true, NOW(), NOW()),
(3, 11, 'Blueberry Cheesecake', 'Creamy blueberry cheesecake slice', 'CF-BCC-013', 35000.00, 30, true, NOW(), NOW()),
(3, 11, 'Apple Pie', 'Classic apple pie slice', 'CF-APP-014', 28000.00, 35, true, NOW(), NOW()),
(3, 11, 'Donut', 'Glazed donut', 'CF-DNT-015', 12000.00, 60, true, NOW(), NOW());

-- Breakfast (category_id: 12)
INSERT INTO products (tenant_id, category_id, name, description, sku, price, stock, is_active, created_at, updated_at) VALUES
(3, 12, 'Pancakes', 'Stack of 3 pancakes with maple syrup', 'CF-PNC-016', 35000.00, 40, true, NOW(), NOW()),
(3, 12, 'French Toast', 'French toast with butter and honey', 'CF-FRT-017', 32000.00, 40, true, NOW(), NOW()),
(3, 12, 'Breakfast Sandwich', 'Egg, cheese, and bacon sandwich', 'CF-BFS-018', 38000.00, 35, true, NOW(), NOW()),
(3, 12, 'Oatmeal Bowl', 'Oatmeal with fruits and nuts', 'CF-OTM-019', 25000.00, 50, true, NOW(), NOW());
(3, 12, 'French Toast', 'French toast with butter and honey', 'CF-FRT-017', 32000.00, 40, true, NOW(), NOW()),
(3, 12, 'Breakfast Sandwich', 'Egg, cheese, and bacon sandwich', 'CF-BFS-018', 38000.00, 35, true, NOW(), NOW()),
(3, 12, 'Oatmeal Bowl', 'Oatmeal with fruits and nuts', 'CF-OTM-019', 25000.00, 50, true, NOW(), NOW());
