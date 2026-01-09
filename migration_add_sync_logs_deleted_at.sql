-- Add deleted_at column to sync_logs table for soft delete support
ALTER TABLE sync_logs ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

-- Add index on deleted_at for better query performance
CREATE INDEX IF NOT EXISTS idx_sync_logs_deleted_at ON sync_logs(deleted_at);

-- Add deleted_at column to sync_conflicts table if it's missing too
ALTER TABLE sync_conflicts ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

-- Add index on deleted_at for sync_conflicts
CREATE INDEX IF NOT EXISTS idx_sync_conflicts_deleted_at ON sync_conflicts(deleted_at);
