package handlers

import (
	"bytes"
	"encoding/json"
	"myposcore/config"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/services"
	"myposcore/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForPIN() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate
	db.AutoMigrate(&models.Tenant{}, &models.Branch{}, &models.User{})

	// Set global DB
	database.DB = db

	return db
}

func createTestUsersForPIN(db *gorm.DB) (superadmin, owner, admin, user models.User) {
	hashedPassword, _ := utils.HashPassword("password123")

	tenant := models.Tenant{
		Name:     "Test Tenant",
		IsActive: true,
	}
	db.Create(&tenant)

	branch := models.Branch{
		TenantID: tenant.ID,
		Name:     "Test Branch",
		IsActive: true,
	}
	db.Create(&branch)

	superadmin = models.User{
		TenantID: tenant.ID,
		BranchID: branch.ID,
		Email:    "superadmin@test.com",
		Password: hashedPassword,
		FullName: "Super Admin",
		Role:     "superadmin",
		IsActive: true,
	}
	db.Create(&superadmin)

	owner = models.User{
		TenantID: tenant.ID,
		BranchID: branch.ID,
		Email:    "owner@test.com",
		Password: hashedPassword,
		FullName: "Owner",
		Role:     "owner",
		IsActive: true,
	}
	db.Create(&owner)

	admin = models.User{
		TenantID: tenant.ID,
		BranchID: branch.ID,
		Email:    "admin@test.com",
		Password: hashedPassword,
		FullName: "Admin",
		Role:     "admin",
		IsActive: true,
	}
	db.Create(&admin)

	user = models.User{
		TenantID: tenant.ID,
		BranchID: branch.ID,
		Email:    "user@test.com",
		Password: hashedPassword,
		FullName: "User",
		Role:     "user",
		IsActive: true,
	}
	db.Create(&user)

	return
}

func TestAdminChangePIN_Success(t *testing.T) {
	// Setup
	db := setupTestDBForPIN()
	superadmin, _, admin, user := createTestUsersForPIN(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePINHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	// Test Case 1: Superadmin changes admin PIN
	t.Run("Superadmin changes admin PIN", func(t *testing.T) {
		reqBody := dto.AdminChangePINRequest{
			Email:      admin.Email,
			PIN:        "654321",
			ConfirmPIN: "654321",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-pin", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", superadmin.ID)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "PIN changed successfully", response["message"])
	})

	// Test Case 2: Admin changes user PIN
	t.Run("Admin changes user PIN", func(t *testing.T) {
		reqBody := dto.AdminChangePINRequest{
			Email:      user.Email,
			PIN:        "111111",
			ConfirmPIN: "111111",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-pin", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", admin.ID)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminChangePIN_PINMismatch(t *testing.T) {
	// Setup
	db := setupTestDBForPIN()
	superadmin, _, _, user := createTestUsersForPIN(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePINHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	reqBody := dto.AdminChangePINRequest{
		Email:      user.Email,
		PIN:        "123456",
		ConfirmPIN: "654321",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-pin", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", superadmin.ID)

	handler.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "do not match")
}

func TestAdminChangePIN_InsufficientPermission(t *testing.T) {
	// Setup
	db := setupTestDBForPIN()
	_, _, admin, user := createTestUsersForPIN(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePINHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	// User tries to change admin PIN (should fail)
	reqBody := dto.AdminChangePINRequest{
		Email:      admin.Email,
		PIN:        "888999",
		ConfirmPIN: "888999",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-pin", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", user.ID)

	handler.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "insufficient permission")
}

func TestAdminChangePIN_UserNotFound(t *testing.T) {
	// Setup
	db := setupTestDBForPIN()
	superadmin, _, _, _ := createTestUsersForPIN(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePINHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	reqBody := dto.AdminChangePINRequest{
		Email:      "nonexistent@test.com",
		PIN:        "123456",
		ConfirmPIN: "123456",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-pin", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", superadmin.ID)

	handler.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "not found")
}
