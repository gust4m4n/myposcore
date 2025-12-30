-- Migration: Create audit_trails table
-- Description: Creates table for tracking all user actions across the system
-- Author: System
-- Date: 2024-12-30

-- Step 1: Create audit_trails table
CREATE TABLE IF NOT EXISTS audit_trails (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id) ON DELETE SET NULL,
    branch_id INTEGER REFERENCES branches(id) ON DELETE SET NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    action VARCHAR(20) NOT NULL,
    changes JSONB,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Step 2: Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_audit_trails_tenant_id ON audit_trails(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_branch_id ON audit_trails(branch_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_user_id ON audit_trails(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_entity_type ON audit_trails(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_trails_entity_id ON audit_trails(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_action ON audit_trails(action);
CREATE INDEX IF NOT EXISTS idx_audit_trails_created_at ON audit_trails(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_trails_deleted_at ON audit_trails(deleted_at);

-- Step 3: Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_audit_trails_entity ON audit_trails(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_user_action ON audit_trails(user_id, action);
CREATE INDEX IF NOT EXISTS idx_audit_trails_tenant_created ON audit_trails(tenant_id, created_at DESC);

-- Step 4: Add comments for documentation
COMMENT ON TABLE audit_trails IS 'Tracks all user actions across the system for audit and compliance';
COMMENT ON COLUMN audit_trails.tenant_id IS 'ID of the tenant (null for global entities like FAQ/TnC)';
COMMENT ON COLUMN audit_trails.branch_id IS 'ID of the branch (null for tenant-level or global entities)';
COMMENT ON COLUMN audit_trails.user_id IS 'ID of the user who performed the action';
COMMENT ON COLUMN audit_trails.entity_type IS 'Type of entity (user, product, order, payment, category, faq, tnc)';
COMMENT ON COLUMN audit_trails.entity_id IS 'ID of the entity that was modified';
COMMENT ON COLUMN audit_trails.action IS 'Type of action performed (create, update, delete)';
COMMENT ON COLUMN audit_trails.changes IS 'JSON object containing the changes made';
COMMENT ON COLUMN audit_trails.ip_address IS 'IP address of the user who performed the action';
COMMENT ON COLUMN audit_trails.user_agent IS 'User agent string from the request';
COMMENT ON COLUMN audit_trails.created_at IS 'Timestamp when the action was performed';

-- Rollback instructions:
-- To rollback this migration, run:
-- DROP INDEX IF EXISTS idx_audit_trails_tenant_created;
-- DROP INDEX IF EXISTS idx_audit_trails_user_action;
-- DROP INDEX IF EXISTS idx_audit_trails_entity;
-- DROP INDEX IF EXISTS idx_audit_trails_deleted_at;
-- DROP INDEX IF EXISTS idx_audit_trails_created_at;
-- DROP INDEX IF EXISTS idx_audit_trails_action;
-- DROP INDEX IF EXISTS idx_audit_trails_entity_id;
-- DROP INDEX IF EXISTS idx_audit_trails_entity_type;
-- DROP INDEX IF EXISTS idx_audit_trails_user_id;
-- DROP INDEX IF EXISTS idx_audit_trails_branch_id;
-- DROP INDEX IF EXISTS idx_audit_trails_tenant_id;
-- DROP TABLE IF EXISTS audit_trails;
