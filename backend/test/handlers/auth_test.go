package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/api/handlers"
	"backend/test"

	"github.com/stretchr/testify/assert"
)

var db = test.SetupTestDB()

func TestRegisterUserH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	router := test.SetupTestRouter()
	router.POST("/api/register", handlers.RegisterH)

	requestBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
		"role": "HR",
		"company_name": "TCS",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])

	user, ok := response["user"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "testuser", user["username"])
}

func TestLoginUserH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// First register a user
	test.InsertTestUser(db)

	router := test.SetupTestRouter()
	router.POST("/api/login", handlers.LoginH)

	requestBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
}
