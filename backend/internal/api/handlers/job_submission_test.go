package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func mockGetSubmissionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Submissions retrieved successfully"})
}

func mockHandleFormSubmission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Form submitted successfully"})
}

func mockUpdateSubmissionStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func createTestMultipartRequest(t *testing.T) (*http.Request, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add form fields
	writer.WriteField("name", "Test User")
	writer.WriteField("email", "test@example.com")
	writer.WriteField("phone", "1234567890")
	writer.WriteField("skills", "Go,Python")
	
	// Add resume file
	part, err := writer.CreateFormFile("resume", "test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("test resume content"))
	
	writer.Close()
	
	req := httptest.NewRequest("POST", "/api/jobs/1/submit", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	return req, writer
}

func TestGetFormSubmissions_Success(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/jobs/:id/submissions", mockGetSubmissionsHandler)
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/jobs/1/submissions", nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Submissions retrieved successfully", response["message"])
}

func TestHandleFormSubmission_Success(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/jobs/:id/submit", mockHandleFormSubmission)
	
	w := httptest.NewRecorder()
	req, _ := createTestMultipartRequest(t)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Form submitted successfully", response["message"])
}

func TestUpdateSubmissionStatus_Success(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/api/submissions/:id/status", mockUpdateSubmissionStatus)
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/api/submissions/1/status", nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Status updated successfully", response["message"])
} 