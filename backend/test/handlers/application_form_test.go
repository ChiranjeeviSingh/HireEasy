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

// Helper function to create a test job for application form tests
func createTestJobForForm(t *testing.T, router *gin.Engine, userID int) string {
	createBody := map[string]interface{}{
		"job_id":          "TEST_JOB_001",
		"job_title":       "Test Job",
		"job_description": "Test Job Description",
		"job_status":      "active",
		"skills_required": []string{"Go", "PostgreSQL"},
		"attributes":      map[string]interface{}{"location": "Remote"},
	}
	jsonBody, _ := json.Marshal(createBody)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Get the job ID from the database
	var jobID int
	err := db.QueryRow("SELECT id FROM jobs WHERE job_id = $1", createBody["job_id"]).Scan(&jobID)
	assert.NoError(t, err)
	return createBody["job_id"].(string)
}

// Helper function to create a test form template for application form tests
func createTestFormTemplateForForm(t *testing.T, router *gin.Engine, userID int) string {
	createBody := map[string]interface{}{
		"form_template_id": "TEST_FORM_001",
		"fields": []map[string]interface{}{
			{
				"id":       "name",
				"type":     "text",
				"label":    "Full Name",
				"required": true,
			},
		},
	}
	jsonBody, _ := json.Marshal(createBody)
	req, _ := http.NewRequest("POST", "/api/forms/templates", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Get the form template ID from the database
	var formTemplateID int
	err := db.QueryRow("SELECT id FROM form_templates WHERE form_template_id = $1", createBody["form_template_id"]).Scan(&formTemplateID)
	assert.NoError(t, err)
	return createBody["form_template_id"].(string)
}

func TestLinkJobToFormTemplateH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and set up router
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	// Set up routes
	router.POST("/api/jobs", handlers.CreateJobH)
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.POST("/api/jobs/:job_id/forms", handlers.LinkJobToFormTemplateH)

	// Create test job and form template
	jobID := createTestJobForForm(t, router, userID)
	formTemplateID := createTestFormTemplateForForm(t, router, userID)

	// Link job to form template
	requestBody := map[string]interface{}{
		"form_template_id": formTemplateID,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/jobs/"+jobID+"/forms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["form_uuid"])
}

func TestUpdateFormStatusH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and set up router
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	// Set up routes
	router.POST("/api/jobs", handlers.CreateJobH)
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.POST("/api/jobs/:job_id/forms", handlers.LinkJobToFormTemplateH)
	router.PATCH("/api/forms/:form_uuid/status", handlers.UpdateFormStatusH)

	// Create test job and form template, then link them
	jobID := createTestJobForForm(t, router, userID)
	formTemplateID := createTestFormTemplateForForm(t, router, userID)

	// Link job to form template
	requestBody := map[string]interface{}{
		"form_template_id": formTemplateID,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/jobs/"+jobID+"/forms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var linkResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &linkResponse)
	assert.NoError(t, err)
	formUUID := linkResponse["form_uuid"].(string)

	// Update form status
	updateRequestBody := map[string]interface{}{
		"status": "inactive",
	}
	jsonBody, _ = json.Marshal(updateRequestBody)
	req, _ = http.NewRequest("PATCH", "/api/forms/"+formUUID+"/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, formUUID, response["form_uuid"])
	assert.Equal(t, "inactive", response["status"])
}

func TestGetFormDetailsH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and set up router
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	// Set up routes
	router.POST("/api/jobs", handlers.CreateJobH)
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.POST("/api/jobs/:job_id/forms", handlers.LinkJobToFormTemplateH)
	router.GET("/api/forms/:form_uuid", handlers.GetFormDetailsH)

	// Create test job and form template, then link them
	jobID := createTestJobForForm(t, router, userID)
	formTemplateID := createTestFormTemplateForForm(t, router, userID)

	// Link job to form template
	requestBody := map[string]interface{}{
		"form_template_id": formTemplateID,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/jobs/"+jobID+"/forms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var linkResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &linkResponse)
	assert.NoError(t, err)
	formUUID := linkResponse["form_uuid"].(string)

	// Get form details
	req, _ = http.NewRequest("GET", "/api/forms/"+formUUID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, formUUID, response["form_uuid"])
	assert.Equal(t, "active", response["status"])

	// Verify job details
	jobDetails := response["job"].(map[string]interface{})
	assert.Equal(t, jobID, jobDetails["job_id"])
	assert.Equal(t, "Test Job", jobDetails["job_title"])
	assert.Equal(t, "Test Job Description", jobDetails["job_description"])

	// Verify form template details
	formTemplateDetails := response["form_template"].(map[string]interface{})
	assert.Equal(t, formTemplateID, formTemplateDetails["form_template_id"])
}

func TestDeleteFormH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and set up router
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	// Set up routes
	router.POST("/api/jobs", handlers.CreateJobH)
	router.POST("/api/forms/templates", handlers.CreateFormTemplateH)
	router.POST("/api/jobs/:job_id/forms", handlers.LinkJobToFormTemplateH)
	router.DELETE("/api/forms/:form_uuid", handlers.DeleteFormH)
	router.GET("/api/forms/:form_uuid", handlers.GetFormDetailsH)

	// Create test job and form template, then link them
	jobID := createTestJobForForm(t, router, userID)
	formTemplateID := createTestFormTemplateForForm(t, router, userID)

	// Link job to form template
	requestBody := map[string]interface{}{
		"form_template_id": formTemplateID,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/jobs/"+jobID+"/forms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var linkResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &linkResponse)
	assert.NoError(t, err)
	formUUID := linkResponse["form_uuid"].(string)

	// Delete the form
	req, _ = http.NewRequest("DELETE", "/api/forms/"+formUUID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify form is deleted by trying to get it
	req, _ = http.NewRequest("GET", "/api/forms/"+formUUID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
