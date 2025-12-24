package main

import (
	"bytes"
	"encoding/json"
	"myposcore/config"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testRouter *gin.Engine
	testConfig *config.Config
	authToken  string
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	err = db.AutoMigrate(&models.Tenant{}, &models.Branch{}, &models.User{})
	if err != nil {
		return nil, err
	}

	// Seed test data
	tenant := models.Tenant{
		Name:     "Test Company",
		Code:     "TENANT001",
		IsActive: true,
	}
	db.Create(&tenant)

	branch := models.Branch{
		TenantID: tenant.ID,
		Name:     "Main Branch",
		Code:     "BRANCH001",
		Address:  "Test Address",
		Phone:    "021-12345678",
		IsActive: true,
	}
	db.Create(&branch)

	return db, nil
}

func setupTestServer(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Override database connection
	database.DB = db

	// Setup test config
	testConfig = &config.Config{
		ServerPort: "8080",
		JWTSecret:  "test-secret-key",
	}

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Setup routes
	testRouter = gin.Default()
	routes.SetupRoutes(testRouter, testConfig)
}

func TestHealthEndpoint(t *testing.T) {
	setupTestServer(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestRegisterEndpoint(t *testing.T) {
	setupTestServer(t)

	registerReq := dto.RegisterRequest{
		TenantCode: "TENANT001",
		BranchCode: "BRANCH001",
		Username:   "testuser",
		Email:      "test@example.com",
		Password:   "password123",
		FullName:   "Test User",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Extract token for later use
	data := response["data"].(map[string]interface{})
	authToken = data["token"].(string)
	assert.NotEmpty(t, authToken)
}

func TestRegisterEndpoint_InvalidTenant(t *testing.T) {
	setupTestServer(t)

	registerReq := dto.RegisterRequest{
		TenantCode: "INVALID",
		BranchCode: "BRANCH001",
		Username:   "testuser",
		Email:      "test@example.com",
		Password:   "password123",
		FullName:   "Test User",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "tenant not found or inactive", response["error"])
}

func TestLoginEndpoint(t *testing.T) {
	setupTestServer(t)

	// First register a user
	registerReq := dto.RegisterRequest{
		TenantCode: "TENANT001",
		BranchCode: "BRANCH001",
		Username:   "loginuser",
		Email:      "login@example.com",
		Password:   "password123",
		FullName:   "Login User",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// Now login
	loginReq := dto.LoginRequest{
		TenantCode: "TENANT001",
		BranchCode: "BRANCH001",
		Username:   "loginuser",
		Password:   "password123",
	}

	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", response["message"])
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	token := data["token"].(string)
	assert.NotEmpty(t, token)
}

func TestLoginEndpoint_InvalidCredentials(t *testing.T) {
	setupTestServer(t)

	loginReq := dto.LoginRequest{
		TenantCode: "TENANT001",
		BranchCode: "BRANCH001",
		Username:   "nonexistent",
		Password:   "wrongpassword",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid credentials", response["error"])
}

func TestProfileEndpoint(t *testing.T) {
	setupTestServer(t)

	// First register a user
	registerReq := dto.RegisterRequest{
		TenantCode: "TENANT001",
		BranchCode: "BRANCH001",
		Username:   "profileuser",
		Email:      "profile@example.com",
		Password:   "password123",
		FullName:   "Profile User",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var registerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	data := registerResponse["data"].(map[string]interface{})
	token := data["token"].(string)

	// Now get profile
	req = httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "profileuser", response["username"])
	assert.NotNil(t, response["user_id"])
	assert.NotNil(t, response["tenant_id"])
}

func TestProfileEndpoint_NoToken(t *testing.T) {
	setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Authorization header required", response["error"])
}

func TestProfileEndpoint_InvalidToken(t *testing.T) {
	setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"].(string), "Invalid or expired token")
}
