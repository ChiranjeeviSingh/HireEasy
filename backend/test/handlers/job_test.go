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

// Custom logger to reduce verbosity
type testLogger struct{}

func (l *testLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.LoggerWithWriter(&testLogger{}))
	return router
}

func createTestJob(t *testing.T, router *gin.Engine, userID int) string {
	createBody := map[string]interface{}{
		"job_id":          "J12345",
		"job_title":       "Software Engineer",
		"job_description": "Looking for a skilled software engineer",
		"job_status":      "active",
		"skills_required": []string{"Go", "PostgreSQL", "Docker"},
		"user_id":         userID,
	}
	jsonBody, _ := json.Marshal(createBody)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["job_id"])
	return response["job_id"].(string)
}

func TestCreateJobH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user first
	userID, _ := test.InsertTestUser(db)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)

	requestBody := map[string]interface{}{
		"job_id":          "J12345",
		"job_title":       "Software Engineer",
		"job_description": "Looking for a skilled software engineer",
		"job_status":      "active",
		"skills_required": []string{"Go", "PostgreSQL", "Docker"},
		"user_id":         userID,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "J12345", response["job_id"])
	assert.Equal(t, "Software Engineer", response["job_title"])
	assert.Equal(t, "Looking for a skilled software engineer", response["job_description"])
	assert.Equal(t, "active", response["job_status"])

	// Verify skills array
	skills, ok := response["skills_required"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "Go")
	assert.Contains(t, skills, "PostgreSQL")
	assert.Contains(t, skills, "Docker")
}

func TestUpdateJobH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and job first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)
	router.PUT("/api/jobs/:jobId", handlers.UpdateJobH)

	// First create a job
	jobID := createTestJob(t, router, userID)

	// Now update the job
	updateBody := map[string]interface{}{
		"job_title":       "Senior Software Engineer",
		"job_description": "Looking for a senior software engineer",
		"job_status":      "active",
		"skills_required": []string{"Go", "PostgreSQL", "Docker", "Kubernetes"},
	}
	jsonBody, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest("PUT", "/api/jobs/"+jobID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, jobID, response["job_id"])
	assert.Equal(t, "Senior Software Engineer", response["job_title"])
	assert.Equal(t, "Looking for a senior software engineer", response["job_description"])
	assert.Equal(t, "active", response["job_status"])

	// Verify skills array
	skills, ok := response["skills_required"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, skills, 4)
	assert.Contains(t, skills, "Go")
	assert.Contains(t, skills, "PostgreSQL")
	assert.Contains(t, skills, "Docker")
	assert.Contains(t, skills, "Kubernetes")
}

func TestGetJobByIdH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and job first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)
	router.GET("/api/jobs/:jobId", handlers.GetJobByIdH)

	// First create a job
	jobID := createTestJob(t, router, userID)

	// Now get the job
	req, _ := http.NewRequest("GET", "/api/jobs/"+jobID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, jobID, response["job_id"])
	assert.Equal(t, "Software Engineer", response["job_title"])
	assert.Equal(t, "Looking for a skilled software engineer", response["job_description"])
	assert.Equal(t, "active", response["job_status"])

	// Verify skills array
	skills, ok := response["skills_required"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "Go")
	assert.Contains(t, skills, "PostgreSQL")
	assert.Contains(t, skills, "Docker")
}

func TestGetJobsByTitleH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and job first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)
	router.GET("/api/jobs/title/:jobtitle", handlers.GetJobsByTitleH)

	// First create a job
	createTestJob(t, router, userID)

	// Now get jobs by title
	req, _ := http.NewRequest("GET", "/api/jobs/title/Software Engineer", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)

	job := response[0]
	assert.Equal(t, "J12345", job["job_id"])
	assert.Equal(t, "Software Engineer", job["job_title"])
	assert.Equal(t, "Looking for a skilled software engineer", job["job_description"])
	assert.Equal(t, "active", job["job_status"])

	// Verify skills array
	skills, ok := job["skills_required"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "Go")
	assert.Contains(t, skills, "PostgreSQL")
	assert.Contains(t, skills, "Docker")
}

func TestGetJobsByStatusH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and job first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)
	router.GET("/api/jobs/status/:status", handlers.GetJobsByStatusH)

	// First create a job
	createTestJob(t, router, userID)

	// Now get jobs by status
	req, _ := http.NewRequest("GET", "/api/jobs/status/active", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)

	job := response[0]
	assert.Equal(t, "J12345", job["job_id"])
	assert.Equal(t, "Software Engineer", job["job_title"])
	assert.Equal(t, "Looking for a skilled software engineer", job["job_description"])
	assert.Equal(t, "active", job["job_status"])

	// Verify skills array
	skills, ok := job["skills_required"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "Go")
	assert.Contains(t, skills, "PostgreSQL")
	assert.Contains(t, skills, "Docker")
}

func TestListUserJobsH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and job first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)
	router.GET("/api/jobs/user", handlers.ListUserJobsH)

	// First create a job
	createTestJob(t, router, userID)

	// Now get user's jobs
	req, _ := http.NewRequest("GET", "/api/jobs/user", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)

	job := response[0]
	assert.Equal(t, "J12345", job["job_id"])
	assert.Equal(t, "Software Engineer", job["job_title"])
	assert.Equal(t, "Looking for a skilled software engineer", job["job_description"])
	assert.Equal(t, "active", job["job_status"])

	// Verify skills array
	skills, ok := job["skills_required"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "Go")
	assert.Contains(t, skills, "PostgreSQL")
	assert.Contains(t, skills, "Docker")
}

func TestDeleteJobH(t *testing.T) {
	// Clean up before test
	test.CleanupTestDB(db)

	// Insert a test user and job first
	userID, _ := test.InsertTestUser(db)
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/jobs", handlers.CreateJobH)
	router.DELETE("/api/jobs/:jobId", handlers.DeleteJobH)

	// First create a job
	jobID := createTestJob(t, router, userID)

	// Now delete the job
	req, _ := http.NewRequest("DELETE", "/api/jobs/"+jobID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Job deleted successfully", response["message"])
}
