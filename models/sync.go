package models

import (
	"time"

	"gorm.io/gorm"
)

// SyncLog - Track all sync operations
type SyncLog struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	TenantID          uint           `gorm:"not null;index" json:"tenant_id"`
	BranchID          *uint          `gorm:"index" json:"branch_id"`
	UserID            *uint          `gorm:"index" json:"user_id"`
	ClientID          string         `gorm:"size:100;not null;index" json:"client_id"`
	SyncType          string         `gorm:"size:50;not null" json:"sync_type"` // upload, download, full, delta
	EntityType        string         `gorm:"size:50" json:"entity_type"`        // orders, payments, products, categories
	RecordsUploaded   int            `gorm:"default:0" json:"records_uploaded"`
	RecordsDownloaded int            `gorm:"default:0" json:"records_downloaded"`
	ConflictsDetected int            `gorm:"default:0" json:"conflicts_detected"`
	Status            string         `gorm:"size:20;not null;index" json:"status"` // started, completed, failed
	ErrorMessage      string         `gorm:"type:text" json:"error_message"`
	StartedAt         time.Time      `gorm:"not null" json:"started_at"`
	CompletedAt       *time.Time     `json:"completed_at"`
	DurationMs        int            `json:"duration_ms"`
	Metadata          *string        `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant Tenant  `gorm:"foreignKey:TenantID" json:"-"`
	Branch *Branch `gorm:"foreignKey:BranchID;constraint:-" json:"branch,omitempty"`
	User   *User   `gorm:"foreignKey:UserID;constraint:-" json:"user,omitempty"`
}

// SyncConflict - Store unresolved sync conflicts
type SyncConflict struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	TenantID           uint           `gorm:"not null;index" json:"tenant_id"`
	BranchID           *uint          `gorm:"index" json:"branch_id"`
	EntityType         string         `gorm:"size:50;not null;index" json:"entity_type"` // order, payment
	EntityID           uint           `gorm:"not null;index" json:"entity_id"`
	ClientID           string         `gorm:"size:100;not null;index" json:"client_id"`
	ClientVersion      int            `gorm:"not null" json:"client_version"`
	ServerVersion      int            `gorm:"not null" json:"server_version"`
	ClientData         string         `gorm:"type:jsonb;not null" json:"client_data"`
	ServerData         string         `gorm:"type:jsonb;not null" json:"server_data"`
	ConflictType       string         `gorm:"size:50;not null" json:"conflict_type"` // update_conflict, delete_conflict
	ResolutionStrategy string         `gorm:"size:50" json:"resolution_strategy"`    // server_wins, client_wins, manual, merge
	Resolved           bool           `gorm:"default:false;index" json:"resolved"`
	ResolvedAt         *time.Time     `json:"resolved_at"`
	ResolvedBy         *uint          `gorm:"index" json:"resolved_by"`
	ResolvedData       *string        `gorm:"type:jsonb" json:"resolved_data,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant   Tenant  `gorm:"foreignKey:TenantID" json:"-"`
	Branch   *Branch `gorm:"foreignKey:BranchID;constraint:-" json:"branch,omitempty"`
	Resolver *User   `gorm:"foreignKey:ResolvedBy;constraint:-" json:"resolver,omitempty"`
}

// TableName overrides
func (SyncLog) TableName() string {
	return "sync_logs"
}

func (SyncConflict) TableName() string {
	return "sync_conflicts"
}
