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
		IsActive: true,
	}
	db.Create(&tenant)

	branch := models.Branch{
		TenantID: tenant.ID,
		Name:     "Main Branch",
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

func TestLoginEndpoint_InvalidCredentials(t *testing.T) {
	setupTestServer(t)

	loginReq := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
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
