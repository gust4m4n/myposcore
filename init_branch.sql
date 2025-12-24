-- Initialize branch data for testing
-- This script creates sample branches for the default tenant

-- Insert sample branches for TENANT001
INSERT INTO branches (tenant_id, name, code, address, phone, is_active, created_at, updated_at) 
VALUES 
  (
    (SELECT id FROM tenants WHERE code = 'TENANT001' LIMIT 1),
    'Main Branch',
    'BRANCH001',
    'Jl. Utama No. 123, Jakarta',
    '021-12345678',
    true,
    NOW(),
    NOW()
  ),
  (
    (SELECT id FROM tenants WHERE code = 'TENANT001' LIMIT 1),
    'Branch Surabaya',
    'BRANCH002',
    'Jl. Pahlawan No. 456, Surabaya',
    '031-87654321',
    true,
    NOW(),
    NOW()
  ),
  (
    (SELECT id FROM tenants WHERE code = 'TENANT001' LIMIT 1),
    'Branch Bandung',
    'BRANCH003',
    'Jl. Asia Afrika No. 789, Bandung',
    '022-11223344',
    true,
    NOW(),
    NOW()
  )
ON CONFLICT (code) DO NOTHING;
