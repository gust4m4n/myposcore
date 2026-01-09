package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"myposcore/dto"
	"myposcore/models"
	"time"

	"gorm.io/gorm"
)

type SyncService struct {
	db *gorm.DB
}

func NewSyncService(db *gorm.DB) *SyncService {
	return &SyncService{db: db}
}

// UploadFromClient - Upload data dari mobile client ke server
func (s *SyncService) UploadFromClient(req *dto.SyncUploadRequest, tenantID, branchID, userID uint) (*dto.SyncUploadResponse, error) {
	// Start sync log
	syncLog := &models.SyncLog{
		TenantID:  tenantID,
		BranchID:  &branchID,
		UserID:    &userID,
		ClientID:  req.ClientID,
		SyncType:  "upload",
		Status:    "started",
		StartedAt: time.Now(),
	}
	if err := s.db.Create(syncLog).Error; err != nil {
		return nil, err
	}

	startTime := time.Now()
	response := &dto.SyncUploadResponse{
		SyncID:          fmt.Sprintf("sync_%d", syncLog.ID),
		Conflicts:       []dto.SyncConflictInfo{},
		TenantMapping:   make(map[string]uint),
		BranchMapping:   make(map[string]uint),
		UserMapping:     make(map[string]uint),
		ProductMapping:  make(map[string]uint),
		CategoryMapping: make(map[string]uint),
		OrderMapping:    make(map[string]uint),
		PaymentMapping:  make(map[string]uint),
		AuditMapping:    make(map[string]uint),
		SyncTimestamp:   time.Now(),
		Errors:          []dto.SyncErrorInfo{},
	}

	// Process all entities in transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Process Tenants
		for _, tenantData := range req.Tenants {
			serverID, err := s.processTenant(tx, &tenantData, userID, req.ClientID)
			if err != nil {
				response.FailedTenants++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "tenant",
					LocalID:    tenantData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedTenants++
			response.TenantMapping[tenantData.LocalID] = serverID
		}

		// 2. Process Branches
		for _, branchData := range req.Branches {
			serverID, err := s.processBranch(tx, &branchData, tenantID, userID, req.ClientID)
			if err != nil {
				response.FailedBranches++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "branch",
					LocalID:    branchData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedBranches++
			response.BranchMapping[branchData.LocalID] = serverID
		}

		// 3. Process Users
		for _, userData := range req.Users {
			serverID, err := s.processUser(tx, &userData, tenantID, userID, req.ClientID)
			if err != nil {
				response.FailedUsers++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "user",
					LocalID:    userData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedUsers++
			response.UserMapping[userData.LocalID] = serverID
		}

		// 4. Process Products
		for _, productData := range req.Products {
			serverID, err := s.processProduct(tx, &productData, tenantID, userID, req.ClientID)
			if err != nil {
				response.FailedProducts++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "product",
					LocalID:    productData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedProducts++
			response.ProductMapping[productData.LocalID] = serverID
		}

		// 5. Process Categories
		for _, categoryData := range req.Categories {
			serverID, err := s.processCategory(tx, &categoryData, tenantID, userID, req.ClientID)
			if err != nil {
				response.FailedCategories++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "category",
					LocalID:    categoryData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedCategories++
			response.CategoryMapping[categoryData.LocalID] = serverID
		}

		// 6. Process Orders
		for _, orderData := range req.Orders {
			serverOrderID, err := s.processOrder(tx, &orderData, tenantID, branchID, userID, req.ClientID)
			if err != nil {
				response.FailedOrders++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "order",
					LocalID:    orderData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedOrders++
			response.OrderMapping[orderData.LocalID] = serverOrderID
		}

		// 7. Process Payments
		for _, paymentData := range req.Payments {
			serverPaymentID, err := s.processPayment(tx, &paymentData, tenantID, branchID, userID, req.ClientID, response.OrderMapping)
			if err != nil {
				response.FailedPayments++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "payment",
					LocalID:    paymentData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedPayments++
			response.PaymentMapping[paymentData.LocalID] = serverPaymentID
		}

		// 8. Process Audit Trails
		for _, auditData := range req.AuditTrails {
			serverID, err := s.processAudit(tx, &auditData, tenantID, userID, req.ClientID)
			if err != nil {
				response.FailedAudits++
				response.Errors = append(response.Errors, dto.SyncErrorInfo{
					EntityType: "audit",
					LocalID:    auditData.LocalID,
					Error:      err.Error(),
				})
				continue
			}
			response.ProcessedAudits++
			response.AuditMapping[auditData.LocalID] = serverID
		}

		return nil
	})

	if err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
	} else {
		syncLog.Status = "completed"
	}

	duration := time.Since(startTime)
	syncLog.CompletedAt = &response.SyncTimestamp
	syncLog.DurationMs = int(duration.Milliseconds())
	syncLog.RecordsUploaded = response.ProcessedTenants + response.ProcessedBranches + response.ProcessedUsers +
		response.ProcessedProducts + response.ProcessedCategories + response.ProcessedOrders +
		response.ProcessedPayments + response.ProcessedAudits
	syncLog.ConflictsDetected = len(response.Conflicts)

	s.db.Save(syncLog)

	return response, err
}

// processOrder - Process single order
func (s *SyncService) processOrder(tx *gorm.DB, orderData *dto.SyncOrderData, tenantID, branchID, userID uint, clientID string) (uint, error) {
	// Check if order already exists (by client_id + local_id)
	var existing models.Order
	err := tx.Where("tenant_id = ? AND branch_id = ? AND client_id = ?", tenantID, branchID, clientID+"_"+orderData.LocalID).First(&existing).Error

	if err == nil {
		// Order exists - check version for conflict
		if existing.Version > orderData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, orderData.Version)
		}
		// Update existing order
		existing.TotalAmount = orderData.TotalAmount
		existing.Status = orderData.Status
		existing.Notes = orderData.Notes
		existing.Version = orderData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	// Create new order
	order := models.Order{
		TenantID:       tenantID,
		BranchID:       branchID,
		UserID:         userID,
		OrderNumber:    orderData.OrderNumber,
		TotalAmount:    orderData.TotalAmount,
		Status:         orderData.Status,
		Notes:          orderData.Notes,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + orderData.LocalID, // Kombinasi untuk uniqueness
		LocalTimestamp: &orderData.LocalTimestamp,
		Version:        orderData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	// Generate order number if not provided
	if order.OrderNumber == "" {
		order.OrderNumber = fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().Unix())
	}

	if err := tx.Create(&order).Error; err != nil {
		return 0, err
	}

	// Create order items
	for _, itemData := range orderData.Items {
		orderItem := models.OrderItem{
			OrderID:    order.ID,
			ProductID:  itemData.ProductID,
			Quantity:   itemData.Quantity,
			Price:      itemData.Price,
			Subtotal:   itemData.Subtotal,
			SyncStatus: "synced",
			ClientID:   clientID + "_" + orderData.LocalID,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			return 0, err
		}

		// Update product stock
		if err := tx.Model(&models.Product{}).Where("id = ?", itemData.ProductID).
			UpdateColumn("stock", gorm.Expr("stock - ?", itemData.Quantity)).Error; err != nil {
			return 0, err
		}
	}

	return order.ID, nil
}

// processPayment - Process single payment
func (s *SyncService) processPayment(tx *gorm.DB, paymentData *dto.SyncPaymentData, tenantID, branchID, userID uint, clientID string, orderMapping map[string]uint) (uint, error) {
	// Get server order ID from mapping
	serverOrderID, ok := orderMapping[paymentData.OrderLocalID]
	if !ok {
		// Try to find order by client_id
		var order models.Order
		if err := tx.Where("tenant_id = ? AND branch_id = ? AND client_id = ?", tenantID, branchID, clientID+"_"+paymentData.OrderLocalID).First(&order).Error; err != nil {
			return 0, fmt.Errorf("order not found for local_id: %s", paymentData.OrderLocalID)
		}
		serverOrderID = order.ID
	}

	// Check if payment already exists
	var existing models.Payment
	err := tx.Where("tenant_id = ? AND branch_id = ? AND client_id = ?", tenantID, branchID, clientID+"_"+paymentData.LocalID).First(&existing).Error

	if err == nil {
		// Payment exists - check version
		if existing.Version > paymentData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, paymentData.Version)
		}
		// Update existing payment
		existing.Amount = paymentData.Amount
		existing.PaymentMethod = paymentData.PaymentMethod
		existing.Status = paymentData.Status
		existing.Notes = paymentData.Notes
		existing.Version = paymentData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	// Create new payment
	payment := models.Payment{
		TenantID:       tenantID,
		BranchID:       branchID,
		OrderID:        serverOrderID,
		Amount:         paymentData.Amount,
		PaymentMethod:  paymentData.PaymentMethod,
		Status:         paymentData.Status,
		Notes:          paymentData.Notes,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + paymentData.LocalID,
		LocalTimestamp: &paymentData.LocalTimestamp,
		Version:        paymentData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := tx.Create(&payment).Error; err != nil {
		return 0, err
	}

	// Update order status to completed
	if err := tx.Model(&models.Order{}).Where("id = ?", serverOrderID).
		Updates(map[string]interface{}{
			"status":     "completed",
			"updated_by": userID,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return 0, err
	}

	return payment.ID, nil
}

// processTenant - Process tenant data from client
func (s *SyncService) processTenant(tx *gorm.DB, tenantData *dto.SyncTenantData, userID uint, clientID string) (uint, error) {
	var existing models.Tenant
	err := tx.Where("client_id = ?", clientID+"_"+tenantData.LocalID).First(&existing).Error

	if err == nil {
		if existing.Version > tenantData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, tenantData.Version)
		}
		existing.Name = tenantData.Name
		existing.Phone = tenantData.Phone
		existing.Address = tenantData.Address
		existing.City = tenantData.City
		existing.Country = tenantData.Country
		existing.PostalCode = tenantData.PostalCode
		existing.Image = tenantData.Image
		existing.IsActive = tenantData.IsActive
		existing.Version = tenantData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	tenant := models.Tenant{
		Name:           tenantData.Name,
		Phone:          tenantData.Phone,
		Address:        tenantData.Address,
		City:           tenantData.City,
		Country:        tenantData.Country,
		PostalCode:     tenantData.PostalCode,
		Image:          tenantData.Image,
		IsActive:       tenantData.IsActive,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + tenantData.LocalID,
		LocalTimestamp: &tenantData.LocalTimestamp,
		Version:        tenantData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := tx.Create(&tenant).Error; err != nil {
		return 0, err
	}
	return tenant.ID, nil
}

// processBranch - Process branch data from client
func (s *SyncService) processBranch(tx *gorm.DB, branchData *dto.SyncBranchData, tenantID, userID uint, clientID string) (uint, error) {
	var existing models.Branch
	err := tx.Where("client_id = ?", clientID+"_"+branchData.LocalID).First(&existing).Error

	if err == nil {
		if existing.Version > branchData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, branchData.Version)
		}
		existing.Name = branchData.Name
		existing.Phone = branchData.Phone
		existing.Address = branchData.Address
		existing.City = branchData.City
		existing.Country = branchData.Country
		existing.PostalCode = branchData.PostalCode
		existing.Image = branchData.Image
		existing.IsActive = branchData.IsActive
		existing.Version = branchData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	branch := models.Branch{
		TenantID:       tenantID,
		Name:           branchData.Name,
		Phone:          branchData.Phone,
		Address:        branchData.Address,
		City:           branchData.City,
		Country:        branchData.Country,
		PostalCode:     branchData.PostalCode,
		Image:          branchData.Image,
		IsActive:       branchData.IsActive,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + branchData.LocalID,
		LocalTimestamp: &branchData.LocalTimestamp,
		Version:        branchData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := tx.Create(&branch).Error; err != nil {
		return 0, err
	}
	return branch.ID, nil
}

// processUser - Process user data from client
func (s *SyncService) processUser(tx *gorm.DB, userData *dto.SyncUserData, tenantID, userID uint, clientID string) (uint, error) {
	var existing models.User
	err := tx.Where("client_id = ?", clientID+"_"+userData.LocalID).First(&existing).Error

	if err == nil {
		if existing.Version > userData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, userData.Version)
		}
		existing.FullName = userData.FullName
		existing.Email = userData.Email
		existing.Phone = userData.Phone
		existing.Role = userData.Role
		existing.Image = userData.Image
		existing.IsActive = userData.IsActive
		if userData.BranchID != nil {
			existing.BranchID = userData.BranchID
		}
		existing.Version = userData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	user := models.User{
		TenantID:       tenantID,
		BranchID:       userData.BranchID,
		FullName:       userData.FullName,
		Email:          userData.Email,
		Password:       userData.Password,
		Phone:          userData.Phone,
		Role:           userData.Role,
		Image:          userData.Image,
		IsActive:       userData.IsActive,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + userData.LocalID,
		LocalTimestamp: &userData.LocalTimestamp,
		Version:        userData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := tx.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// processProduct - Process product data from client
func (s *SyncService) processProduct(tx *gorm.DB, productData *dto.SyncProductData, tenantID, userID uint, clientID string) (uint, error) {
	var existing models.Product
	err := tx.Where("client_id = ?", clientID+"_"+productData.LocalID).First(&existing).Error

	if err == nil {
		if existing.Version > productData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, productData.Version)
		}
		existing.Name = productData.Name
		existing.Description = productData.Description
		existing.SKU = productData.SKU
		existing.Price = productData.Price
		existing.Stock = productData.Stock
		categoryID := productData.CategoryID
		existing.CategoryID = &categoryID
		existing.Image = productData.Image
		existing.IsActive = productData.IsActive
		existing.Version = productData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	categoryID := productData.CategoryID
	product := models.Product{
		TenantID:       tenantID,
		CategoryID:     &categoryID,
		Name:           productData.Name,
		Description:    productData.Description,
		SKU:            productData.SKU,
		Price:          productData.Price,
		Stock:          productData.Stock,
		Image:          productData.Image,
		IsActive:       productData.IsActive,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + productData.LocalID,
		LocalTimestamp: &productData.LocalTimestamp,
		Version:        productData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := tx.Create(&product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

// processCategory - Process category data from client
func (s *SyncService) processCategory(tx *gorm.DB, categoryData *dto.SyncCategoryData, tenantID, userID uint, clientID string) (uint, error) {
	var existing models.Category
	err := tx.Where("client_id = ?", clientID+"_"+categoryData.LocalID).First(&existing).Error

	if err == nil {
		if existing.Version > categoryData.Version {
			return existing.ID, fmt.Errorf("version conflict: server version %d > client version %d", existing.Version, categoryData.Version)
		}
		existing.Name = categoryData.Name
		existing.Description = categoryData.Description
		existing.Image = categoryData.Image
		existing.IsActive = categoryData.IsActive
		existing.Version = categoryData.Version + 1
		existing.SyncStatus = "synced"
		existing.UpdatedBy = &userID
		if err := tx.Save(&existing).Error; err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	category := models.Category{
		TenantID:       tenantID,
		Name:           categoryData.Name,
		Description:    categoryData.Description,
		Image:          categoryData.Image,
		IsActive:       categoryData.IsActive,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + categoryData.LocalID,
		LocalTimestamp: &categoryData.LocalTimestamp,
		Version:        categoryData.Version,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := tx.Create(&category).Error; err != nil {
		return 0, err
	}
	return category.ID, nil
}

// processAudit - Process audit trail data from client
func (s *SyncService) processAudit(tx *gorm.DB, auditData *dto.SyncAuditData, tenantID, userID uint, clientID string) (uint, error) {
	tenantPtr := &tenantID
	audit := models.AuditTrail{
		TenantID:       tenantPtr,
		UserID:         auditData.UserID,
		Action:         auditData.Action,
		EntityType:     auditData.TableName,
		EntityID:       &auditData.RecordID,
		Changes:        auditData.NewValues, // Store new values in changes field
		IPAddress:      auditData.IPAddress,
		UserAgent:      auditData.UserAgent,
		SyncStatus:     "synced",
		ClientID:       clientID + "_" + auditData.LocalID,
		LocalTimestamp: &auditData.LocalTimestamp,
		Version:        auditData.Version,
	}

	if err := tx.Create(&audit).Error; err != nil {
		return 0, err
	}
	return audit.ID, nil
}

// DownloadToClient - Download master data ke mobile client
func (s *SyncService) DownloadToClient(req *dto.SyncDownloadRequest, tenantID uint) (*dto.SyncDownloadResponse, error) {
	// Start sync log
	syncLog := &models.SyncLog{
		TenantID:  tenantID,
		ClientID:  req.ClientID,
		SyncType:  "download",
		Status:    "started",
		StartedAt: time.Now(),
	}
	if err := s.db.Create(syncLog).Error; err != nil {
		return nil, err
	}

	startTime := time.Now()
	response := &dto.SyncDownloadResponse{
		SyncID:        fmt.Sprintf("sync_%d", syncLog.ID),
		SyncTimestamp: time.Now(),
		HasMore:       false,
	}

	// Determine which entities to download
	entitiesToDownload := req.EntityTypes
	if len(entitiesToDownload) == 0 {
		// Default: download all master data
		entitiesToDownload = []string{"tenants", "branches", "users", "products", "categories"}
	}

	recordsCount := 0

	for _, entityType := range entitiesToDownload {
		switch entityType {
		case "tenants":
			tenants, err := s.getTenantsForSync(req.LastSyncAt)
			if err != nil {
				syncLog.Status = "failed"
				syncLog.ErrorMessage = err.Error()
				s.db.Save(syncLog)
				return nil, err
			}
			response.Tenants = tenants
			recordsCount += len(tenants)

		case "branches":
			branches, err := s.getBranchesForSync(tenantID, req.LastSyncAt)
			if err != nil {
				syncLog.Status = "failed"
				syncLog.ErrorMessage = err.Error()
				s.db.Save(syncLog)
				return nil, err
			}
			response.Branches = branches
			recordsCount += len(branches)

		case "users":
			users, err := s.getUsersForSync(tenantID, req.LastSyncAt)
			if err != nil {
				syncLog.Status = "failed"
				syncLog.ErrorMessage = err.Error()
				s.db.Save(syncLog)
				return nil, err
			}
			response.Users = users
			recordsCount += len(users)

		case "products":
			products, err := s.getProductsForSync(tenantID, req.LastSyncAt)
			if err != nil {
				syncLog.Status = "failed"
				syncLog.ErrorMessage = err.Error()
				s.db.Save(syncLog)
				return nil, err
			}
			response.Products = products
			recordsCount += len(products)

		case "categories":
			categories, err := s.getCategoriesForSync(tenantID, req.LastSyncAt)
			if err != nil {
				syncLog.Status = "failed"
				syncLog.ErrorMessage = err.Error()
				s.db.Save(syncLog)
				return nil, err
			}
			response.Categories = categories
			recordsCount += len(categories)
		}
	}

	// Update sync log
	duration := time.Since(startTime)
	syncLog.Status = "completed"
	syncLog.CompletedAt = &response.SyncTimestamp
	syncLog.DurationMs = int(duration.Milliseconds())
	syncLog.RecordsDownloaded = recordsCount
	s.db.Save(syncLog)

	return response, nil
}

// getProductsForSync - Get products for sync (delta or full)
func (s *SyncService) getProductsForSync(tenantID uint, lastSyncAt *time.Time) ([]dto.ProductResponse, error) {
	var products []models.Product
	query := s.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	// Delta sync: only get updated products
	if lastSyncAt != nil {
		query = query.Where("updated_at > ?", lastSyncAt)
	}

	if err := query.Preload("CategoryDetail").Find(&products).Error; err != nil {
		return nil, err
	}

	// Convert to response DTO
	var response []dto.ProductResponse
	for _, p := range products {
		productResp := dto.ProductResponse{
			ID:          p.ID,
			TenantID:    p.TenantID,
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			Price:       p.Price,
			Stock:       p.Stock,
			IsActive:    p.IsActive,
			Image:       p.Image,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if p.CategoryID != nil {
			productResp.CategoryID = p.CategoryID
			if p.CategoryDetail != nil {
				productResp.CategoryDetail = &dto.CategorySummary{
					ID:   p.CategoryDetail.ID,
					Name: p.CategoryDetail.Name,
				}
			}
		}

		response = append(response, productResp)
	}

	return response, nil
}

// getCategoriesForSync - Get categories for sync
func (s *SyncService) getCategoriesForSync(tenantID uint, lastSyncAt *time.Time) ([]dto.CategoryResponse, error) {
	var categories []models.Category
	query := s.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if lastSyncAt != nil {
		query = query.Where("updated_at > ?", lastSyncAt)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	var response []dto.CategoryResponse
	for _, c := range categories {
		response = append(response, dto.CategoryResponse{
			ID:          c.ID,
			TenantID:    c.TenantID,
			Name:        c.Name,
			Description: c.Description,
			Image:       c.Image,
			IsActive:    c.IsActive,
			CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// getTenantsForSync - Get tenants for sync (delta or full)
func (s *SyncService) getTenantsForSync(lastSyncAt *time.Time) ([]dto.TenantResponse, error) {
	var tenants []models.Tenant
	query := s.db.Where("deleted_at IS NULL")

	// Delta sync: only get records updated after last sync
	if lastSyncAt != nil {
		query = query.Where("updated_at > ?", lastSyncAt)
	}

	if err := query.Find(&tenants).Error; err != nil {
		return nil, err
	}

	var response []dto.TenantResponse
	for _, t := range tenants {
		response = append(response, dto.TenantResponse{
			ID:         t.ID,
			Name:       t.Name,
			Phone:      t.Phone,
			Address:    t.Address,
			City:       t.City,
			Country:    t.Country,
			PostalCode: t.PostalCode,
			Image:      t.Image,
			IsActive:   t.IsActive,
			CreatedAt:  t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// getBranchesForSync - Get branches for sync (delta or full)
func (s *SyncService) getBranchesForSync(tenantID uint, lastSyncAt *time.Time) ([]dto.BranchResponse, error) {
	var branches []models.Branch
	query := s.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	// Delta sync: only get records updated after last sync
	if lastSyncAt != nil {
		query = query.Where("updated_at > ?", lastSyncAt)
	}

	if err := query.Find(&branches).Error; err != nil {
		return nil, err
	}

	var response []dto.BranchResponse
	for _, b := range branches {
		response = append(response, dto.BranchResponse{
			ID:         b.ID,
			TenantID:   b.TenantID,
			Name:       b.Name,
			Phone:      b.Phone,
			Address:    b.Address,
			City:       b.City,
			Country:    b.Country,
			PostalCode: b.PostalCode,
			Image:      b.Image,
			IsActive:   b.IsActive,
			CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// getUsersForSync - Get users for sync (delta or full)
func (s *SyncService) getUsersForSync(tenantID uint, lastSyncAt *time.Time) ([]dto.UserResponse, error) {
	var users []models.User
	query := s.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	// Delta sync: only get records updated after last sync
	if lastSyncAt != nil {
		query = query.Where("updated_at > ?", lastSyncAt)
	}

	if err := query.Preload("Branch").Find(&users).Error; err != nil {
		return nil, err
	}

	var response []dto.UserResponse
	for _, u := range users {
		userResp := dto.UserResponse{
			ID:       u.ID,
			TenantID: u.TenantID,
			FullName: u.FullName,
			Email:    u.Email,
			Phone:    u.Phone,
			Role:     u.Role,
			Image:    u.Image,
			IsActive: u.IsActive,
		}
		if u.BranchID != nil {
			userResp.BranchID = u.BranchID
			if u.Branch.ID > 0 {
				userResp.BranchName = u.Branch.Name
			}
		}
		response = append(response, userResp)
	}

	return response, nil
}

// GetSyncStatus - Get sync status for client
func (s *SyncService) GetSyncStatus(clientID string, tenantID uint) (*dto.SyncStatusResponse, error) {
	response := &dto.SyncStatusResponse{
		ClientID:        clientID,
		ServerTimestamp: time.Now(),
	}

	// Get last sync log
	var lastSync models.SyncLog
	if err := s.db.Where("client_id = ? AND tenant_id = ? AND status = ?", clientID, tenantID, "completed").
		Order("created_at DESC").First(&lastSync).Error; err == nil {
		response.LastSyncAt = &lastSync.CreatedAt
		response.LastSyncSuccess = lastSync.Status == "completed"
	}

	// Count total syncs
	var totalSyncs int64
	s.db.Model(&models.SyncLog{}).Where("client_id = ? AND tenant_id = ?", clientID, tenantID).Count(&totalSyncs)
	response.TotalSyncs = int(totalSyncs)

	// Count pending uploads (orders/payments with sync_status = pending)
	var pendingOrders, pendingPayments int64
	s.db.Model(&models.Order{}).Where("tenant_id = ? AND sync_status = ?", tenantID, "pending").Count(&pendingOrders)
	s.db.Model(&models.Payment{}).Where("tenant_id = ? AND sync_status = ?", tenantID, "pending").Count(&pendingPayments)
	response.PendingUploads = int(pendingOrders + pendingPayments)

	// Count unresolved conflicts
	var conflicts int64
	s.db.Model(&models.SyncConflict{}).Where("tenant_id = ? AND resolved = ?", tenantID, false).Count(&conflicts)
	response.PendingConflicts = int(conflicts)

	return response, nil
}

// GetSyncLogs - Get sync history
func (s *SyncService) GetSyncLogs(clientID string, tenantID uint, page, pageSize int) ([]dto.SyncLogResponse, int64, error) {
	var logs []models.SyncLog
	var total int64

	offset := (page - 1) * pageSize
	query := s.db.Model(&models.SyncLog{}).Where("client_id = ? AND tenant_id = ?", clientID, tenantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	var response []dto.SyncLogResponse
	for _, log := range logs {
		response = append(response, dto.SyncLogResponse{
			ID:                log.ID,
			ClientID:          log.ClientID,
			SyncType:          log.SyncType,
			EntityType:        log.EntityType,
			RecordsUploaded:   log.RecordsUploaded,
			RecordsDownloaded: log.RecordsDownloaded,
			ConflictsDetected: log.ConflictsDetected,
			Status:            log.Status,
			ErrorMessage:      log.ErrorMessage,
			DurationMs:        log.DurationMs,
			CreatedAt:         log.StartedAt,
			CompletedAt:       log.CompletedAt,
		})
	}

	return response, total, nil
}

// ResolveConflict - Manually resolve a conflict
func (s *SyncService) ResolveConflict(req *dto.ResolveConflictRequest, userID uint) error {
	var conflict models.SyncConflict
	if err := s.db.First(&conflict, req.ConflictID).Error; err != nil {
		return err
	}

	if conflict.Resolved {
		return fmt.Errorf("conflict already resolved")
	}

	now := time.Now()
	conflict.Resolved = true
	conflict.ResolvedAt = &now
	conflict.ResolvedBy = &userID
	conflict.ResolutionStrategy = req.ResolutionStrategy

	// Apply resolution based on strategy
	switch req.ResolutionStrategy {
	case "server_wins":
		resolvedBytes, _ := json.Marshal(conflict.ServerData)
		resolvedStr := string(resolvedBytes)
		conflict.ResolvedData = &resolvedStr
	case "client_wins":
		resolvedBytes, _ := json.Marshal(conflict.ClientData)
		resolvedStr := string(resolvedBytes)
		conflict.ResolvedData = &resolvedStr
	case "manual":
		conflict.ResolvedData = &req.ResolvedData
	}

	return s.db.Save(&conflict).Error
}
