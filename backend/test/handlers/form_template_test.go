package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/api/handlers"
	"backend/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createTestFormTemplate(t *testing.T, router *gin.Engine, userID int) string {
	createBody := map[string]interface{}{
		"form_template_id": "TEMPLATE_001",
		"user_id":          userID,
		"fields": []map[string]interface{}{
			{
				"field_name": "full_name",
				"field_type": "text",
				"required":   true,
				"label":      "Full Name",
			},
			{
				"field_name": "experience",
				"field_type": "number",
				"required":   true,
				"label":      "Years of Experience",
			},
			{
				"field_name": "skills",
				"field_type": "textarea",
				"required":   true,
				"label":      "Technical Skills",
			},
		},
	}
	jsonBody, _ := json.Marshal(createBody)
	req, _ := http.NewRequest("POST", "/api/forms/templates", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["form_template_id"])
	return response["form_template_id"].(string)
}

func TestCreateFormTemplateH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user first
	userID, _ := test.InsertTestUser(db)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)

	requestBody := map[string]interface{}{
		"form_template_id": "TEMPLATE_001",
		"user_id":          userID,
		"fields": []map[string]interface{}{
			{
				"field_name": "full_name",
				"field_type": "text",
				"required":   true,
				"label":      "Full Name",
			},
			{
				"field_name": "experience",
				"field_type": "number",
				"required":   true,
				"label":      "Years of Experience",
			},
			{
				"field_name": "skills",
				"field_type": "textarea",
				"required":   true,
				"label":      "Technical Skills",
			},
		},
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/forms/templates", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "TEMPLATE_001", response["form_template_id"])
	assert.Equal(t, float64(userID), response["user_id"])

	// Verify fields array
	fields, ok := response["fields"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, fields, 3)

	// Verify first field
	firstField := fields[0].(map[string]interface{})
	assert.Equal(t, "full_name", firstField["field_name"])
	assert.Equal(t, "text", firstField["field_type"])
	assert.Equal(t, "Full Name", firstField["label"])
}

func TestGetFormTemplateH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and form template first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.GET("/api/forms/templates/:form_template_id", handlers.GetFormTemplateH)

	// First create a form template
	templateID := createTestFormTemplate(t, router, userID)

	// Now get the form template
	req, _ := http.NewRequest("GET", "/api/forms/templates/"+templateID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, templateID, response["form_template_id"])
	assert.Equal(t, float64(userID), response["user_id"])

	// Verify fields array
	fields, ok := response["fields"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, fields, 3)

	// Verify first field
	firstField := fields[0].(map[string]interface{})
	assert.Equal(t, "full_name", firstField["field_name"])
	assert.Equal(t, "text", firstField["field_type"])
	assert.Equal(t, "Full Name", firstField["label"])
}

func TestListFormTemplatesH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and form template first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.GET("/api/forms/templates", handlers.ListFormTemplatesH)

	// First create a form template
	createTestFormTemplate(t, router, userID)

	// Now list all form templates
	req, _ := http.NewRequest("GET", "/api/forms/templates", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)

	template := response[0]
	assert.Equal(t, "TEMPLATE_001", template["form_template_id"])
	assert.Equal(t, float64(userID), template["user_id"])

	// Verify fields array
	fields, ok := template["fields"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, fields, 3)

	// Verify first field
	firstField := fields[0].(map[string]interface{})
	assert.Equal(t, "full_name", firstField["field_name"])
	assert.Equal(t, "text", firstField["field_type"])
	assert.Equal(t, "Full Name", firstField["label"])
}

func TestDeleteFormTemplateH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and form template first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.DELETE("/api/forms/templates/:form_template_id", handlers.DeleteFormTemplateH)

	// First create a form template
	templateID := createTestFormTemplate(t, router, userID)

	// Now delete the form template
	req, _ := http.NewRequest("DELETE", "/api/forms/templates/"+templateID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify template is deleted by trying to get it
	req, _ = http.NewRequest("GET", "/api/forms/templates/"+templateID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
