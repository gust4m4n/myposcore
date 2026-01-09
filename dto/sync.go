package dto

import "time"

// ============================================
// Sync Request/Response DTOs
// ============================================

// SyncUploadRequest - Request untuk upload data dari client ke server
type SyncUploadRequest struct {
	ClientID        string             `json:"client_id" binding:"required"`
	ClientTimestamp time.Time          `json:"client_timestamp" binding:"required"`
	Tenants         []SyncTenantData   `json:"tenants,omitempty"`
	Branches        []SyncBranchData   `json:"branches,omitempty"`
	Users           []SyncUserData     `json:"users,omitempty"`
	Products        []SyncProductData  `json:"products,omitempty"`
	Categories      []SyncCategoryData `json:"categories,omitempty"`
	Orders          []SyncOrderData    `json:"orders,omitempty"`
	Payments        []SyncPaymentData  `json:"payments,omitempty"`
	AuditTrails     []SyncAuditData    `json:"audit_trails,omitempty"`
	LastSyncAt      *time.Time         `json:"last_sync_at,omitempty"` // Untuk delta sync
}

// SyncDownloadRequest - Request untuk download data dari server ke client
type SyncDownloadRequest struct {
	ClientID    string     `json:"client_id" binding:"required"`
	LastSyncAt  *time.Time `json:"last_sync_at,omitempty"` // Untuk delta sync
	EntityTypes []string   `json:"entity_types,omitempty"` // ["tenants", "branches", "users", "products", "categories"] - optional filter
}

// SyncOrderData - Data order dari client
type SyncOrderData struct {
	LocalID        string              `json:"local_id" binding:"required"` // UUID dari client
	OrderNumber    string              `json:"order_number,omitempty"`
	TotalAmount    float64             `json:"total_amount" binding:"required"`
	Status         string              `json:"status"`
	Notes          string              `json:"notes"`
	Items          []SyncOrderItemData `json:"items" binding:"required,min=1"`
	LocalTimestamp time.Time           `json:"local_timestamp" binding:"required"`
	Version        int                 `json:"version"`
}

// SyncOrderItemData - Data order item dari client
type SyncOrderItemData struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required"`
	Subtotal  float64 `json:"subtotal" binding:"required"`
}

// SyncPaymentData - Data payment dari client
type SyncPaymentData struct {
	LocalID        string    `json:"local_id" binding:"required"`       // UUID dari client
	OrderLocalID   string    `json:"order_local_id" binding:"required"` // Reference ke order local_id
	Amount         float64   `json:"amount" binding:"required"`
	PaymentMethod  string    `json:"payment_method" binding:"required"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncTenantData - Data tenant dari client
type SyncTenantData struct {
	LocalID        string    `json:"local_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
	PostalCode     string    `json:"postal_code"`
	Image          string    `json:"image"`
	IsActive       bool      `json:"is_active"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncBranchData - Data branch dari client
type SyncBranchData struct {
	LocalID        string    `json:"local_id" binding:"required"`
	TenantID       uint      `json:"tenant_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
	PostalCode     string    `json:"postal_code"`
	Image          string    `json:"image"`
	IsActive       bool      `json:"is_active"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncUserData - Data user dari client
type SyncUserData struct {
	LocalID        string    `json:"local_id" binding:"required"`
	TenantID       uint      `json:"tenant_id" binding:"required"`
	BranchID       *uint     `json:"branch_id,omitempty"`
	FullName       string    `json:"full_name" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	Password       string    `json:"password,omitempty"` // Only for new users
	Phone          string    `json:"phone"`
	Role           string    `json:"role" binding:"required"`
	Image          string    `json:"image"`
	IsActive       bool      `json:"is_active"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncProductData - Data product dari client
type SyncProductData struct {
	LocalID        string    `json:"local_id" binding:"required"`
	TenantID       uint      `json:"tenant_id" binding:"required"`
	CategoryID     uint      `json:"category_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Description    string    `json:"description"`
	SKU            string    `json:"sku"`
	Price          float64   `json:"price" binding:"required"`
	Stock          int       `json:"stock"`
	Image          string    `json:"image"`
	IsActive       bool      `json:"is_active"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncCategoryData - Data category dari client
type SyncCategoryData struct {
	LocalID        string    `json:"local_id" binding:"required"`
	TenantID       uint      `json:"tenant_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	IsActive       bool      `json:"is_active"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncAuditData - Data audit trail dari client
type SyncAuditData struct {
	LocalID        string    `json:"local_id" binding:"required"`
	TenantID       uint      `json:"tenant_id" binding:"required"`
	UserID         uint      `json:"user_id" binding:"required"`
	Action         string    `json:"action" binding:"required"`
	TableName      string    `json:"table_name" binding:"required"`
	RecordID       uint      `json:"record_id"`
	OldValues      string    `json:"old_values"`
	NewValues      string    `json:"new_values"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	LocalTimestamp time.Time `json:"local_timestamp" binding:"required"`
	Version        int       `json:"version"`
}

// SyncUploadResponse - Response setelah upload
type SyncUploadResponse struct {
	SyncID              string             `json:"sync_id"`
	ProcessedTenants    int                `json:"processed_tenants"`
	ProcessedBranches   int                `json:"processed_branches"`
	ProcessedUsers      int                `json:"processed_users"`
	ProcessedProducts   int                `json:"processed_products"`
	ProcessedCategories int                `json:"processed_categories"`
	ProcessedOrders     int                `json:"processed_orders"`
	ProcessedPayments   int                `json:"processed_payments"`
	ProcessedAudits     int                `json:"processed_audits"`
	FailedTenants       int                `json:"failed_tenants"`
	FailedBranches      int                `json:"failed_branches"`
	FailedUsers         int                `json:"failed_users"`
	FailedProducts      int                `json:"failed_products"`
	FailedCategories    int                `json:"failed_categories"`
	FailedOrders        int                `json:"failed_orders"`
	FailedPayments      int                `json:"failed_payments"`
	FailedAudits        int                `json:"failed_audits"`
	Conflicts           []SyncConflictInfo `json:"conflicts,omitempty"`
	TenantMapping       map[string]uint    `json:"tenant_mapping,omitempty"`
	BranchMapping       map[string]uint    `json:"branch_mapping,omitempty"`
	UserMapping         map[string]uint    `json:"user_mapping,omitempty"`
	ProductMapping      map[string]uint    `json:"product_mapping,omitempty"`
	CategoryMapping     map[string]uint    `json:"category_mapping,omitempty"`
	OrderMapping        map[string]uint    `json:"order_mapping,omitempty"`
	PaymentMapping      map[string]uint    `json:"payment_mapping,omitempty"`
	AuditMapping        map[string]uint    `json:"audit_mapping,omitempty"`
	SyncTimestamp       time.Time          `json:"sync_timestamp"`
	Errors              []SyncErrorInfo    `json:"errors,omitempty"`
}

// SyncDownloadResponse - Response untuk download data
type SyncDownloadResponse struct {
	SyncID        string             `json:"sync_id"`
	Tenants       []TenantResponse   `json:"tenants,omitempty"`
	Branches      []BranchResponse   `json:"branches,omitempty"`
	Users         []UserResponse     `json:"users,omitempty"`
	Products      []ProductResponse  `json:"products,omitempty"`
	Categories    []CategoryResponse `json:"categories,omitempty"`
	SyncTimestamp time.Time          `json:"sync_timestamp"`
	HasMore       bool               `json:"has_more"` // Untuk pagination jika data besar
}

// SyncConflictInfo - Informasi conflict yang terdeteksi
type SyncConflictInfo struct {
	EntityType    string `json:"entity_type"` // "order", "payment"
	LocalID       string `json:"local_id"`
	ConflictType  string `json:"conflict_type"` // "version_mismatch", "already_exists", "deleted"
	ClientVersion int    `json:"client_version"`
	ServerVersion int    `json:"server_version"`
	Resolution    string `json:"resolution"` // "server_wins", "client_wins", "manual_required"
	Message       string `json:"message"`
}

// SyncErrorInfo - Informasi error saat sync
type SyncErrorInfo struct {
	EntityType string `json:"entity_type"`
	LocalID    string `json:"local_id"`
	Error      string `json:"error"`
}

// SyncStatusRequest - Request untuk cek sync status
type SyncStatusRequest struct {
	ClientID string `json:"client_id" binding:"required"`
}

// SyncStatusResponse - Response status sync
type SyncStatusResponse struct {
	ClientID         string     `json:"client_id"`
	LastSyncAt       *time.Time `json:"last_sync_at"`
	PendingUploads   int        `json:"pending_uploads"`
	PendingConflicts int        `json:"pending_conflicts"`
	TotalSyncs       int        `json:"total_syncs"`
	LastSyncSuccess  bool       `json:"last_sync_success"`
	ServerTimestamp  time.Time  `json:"server_timestamp"`
}

// ResolveConflictRequest - Request untuk manual resolve conflict
type ResolveConflictRequest struct {
	ConflictID         uint   `json:"conflict_id" binding:"required"`
	ResolutionStrategy string `json:"resolution_strategy" binding:"required,oneof=server_wins client_wins manual"`
	ResolvedData       string `json:"resolved_data,omitempty"` // JSON string jika strategy=manual
}

// SyncLogResponse - Response untuk sync history/logs
type SyncLogResponse struct {
	ID                uint       `json:"id"`
	ClientID          string     `json:"client_id"`
	SyncType          string     `json:"sync_type"`
	EntityType        string     `json:"entity_type"`
	RecordsUploaded   int        `json:"records_uploaded"`
	RecordsDownloaded int        `json:"records_downloaded"`
	ConflictsDetected int        `json:"conflicts_detected"`
	Status            string     `json:"status"`
	ErrorMessage      string     `json:"error_message,omitempty"`
	DurationMs        int        `json:"duration_ms"`
	CreatedAt         time.Time  `json:"created_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
}

// ConflictResponse - Response untuk conflict data
type ConflictResponse struct {
	ID                 uint       `json:"id"`
	EntityType         string     `json:"entity_type"`
	EntityID           uint       `json:"entity_id"`
	ClientID           string     `json:"client_id"`
	ConflictType       string     `json:"conflict_type"`
	ResolutionStrategy string     `json:"resolution_strategy,omitempty"`
	Resolved           bool       `json:"resolved"`
	CreatedAt          time.Time  `json:"created_at"`
	ResolvedAt         *time.Time `json:"resolved_at,omitempty"`
}
