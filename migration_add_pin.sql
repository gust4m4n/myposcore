-- Migration untuk menambahkan kolom PIN ke tabel users
-- Execute this SQL if PIN column doesn't exist in your database

-- Add PIN column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS pin VARCHAR(255);

-- PIN column akan di-populate oleh user melalui API /pin/create
-- PIN tidak mandatory, user bisa punya account tanpa PIN
-- PIN akan di-hash sama seperti password untuk security

-- Query untuk cek users yang sudah set PIN
SELECT 
    id, 
    username, 
    email,
    CASE 
        WHEN pin IS NULL OR pin = '' THEN 'No PIN'
        ELSE 'Has PIN'
    END as pin_status
FROM users
ORDER BY id;

-- Note: 
-- - GORM AutoMigrate akan otomatis add column ini saat restart app
-- - Jika manual migration diperlukan, run ALTER TABLE di atas
-- - PIN format: 6 digit numeric, stored as bcrypt hash
