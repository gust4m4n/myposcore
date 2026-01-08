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

func setupTestDB() *gorm.DB {
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

func createTestUsers(db *gorm.DB) (superadmin, owner, admin, user models.User) {
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

func TestAdminChangePassword_Success(t *testing.T) {
	// Setup
	db := setupTestDB()
	superadmin, _, admin, user := createTestUsers(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePasswordHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	// Test Case 1: Superadmin changes owner password
	t.Run("Superadmin changes admin password", func(t *testing.T) {
		reqBody := dto.AdminChangePasswordRequest{
			Email:           admin.Email,
			Password:        "newpassword123",
			ConfirmPassword: "newpassword123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-password", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", superadmin.ID)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Password changed successfully", response["message"])
	})

	// Test Case 2: Admin changes user password
	t.Run("Admin changes user password", func(t *testing.T) {
		reqBody := dto.AdminChangePasswordRequest{
			Email:           user.Email,
			Password:        "newuserpass123",
			ConfirmPassword: "newuserpass123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-password", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", admin.ID)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminChangePassword_PasswordMismatch(t *testing.T) {
	// Setup
	db := setupTestDB()
	superadmin, _, _, user := createTestUsers(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePasswordHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	reqBody := dto.AdminChangePasswordRequest{
		Email:           user.Email,
		Password:        "password1",
		ConfirmPassword: "password2",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-password", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", superadmin.ID)

	handler.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "do not match")
}

func TestAdminChangePassword_InsufficientPermission(t *testing.T) {
	// Setup
	db := setupTestDB()
	_, _, admin, user := createTestUsers(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePasswordHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	// User tries to change admin password (should fail)
	reqBody := dto.AdminChangePasswordRequest{
		Email:           admin.Email,
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-password", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", user.ID)

	handler.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "insufficient permission")
}

func TestAdminChangePassword_UserNotFound(t *testing.T) {
	// Setup
	db := setupTestDB()
	superadmin, _, _, _ := createTestUsers(db)

	cfg := &config.Config{JWTSecret: "test-secret"}
	auditTrailService := services.NewAuditTrailService(db)
	handler := NewAdminChangePasswordHandler(cfg, auditTrailService)

	gin.SetMode(gin.TestMode)

	reqBody := dto.AdminChangePasswordRequest{
		Email:           "nonexistent@test.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/v1/admin/change-password", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", superadmin.ID)

	handler.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "not found")
}
