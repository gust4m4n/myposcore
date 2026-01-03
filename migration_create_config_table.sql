-- Create config table for key-value storage
CREATE TABLE IF NOT EXISTS configs (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on key (already primary key, so indexed by default)
CREATE INDEX IF NOT EXISTS idx_configs_key ON configs(key);

-- Add trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_configs_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_configs_updated_at
    BEFORE UPDATE ON configs
    FOR EACH ROW
    EXECUTE FUNCTION update_configs_updated_at();
