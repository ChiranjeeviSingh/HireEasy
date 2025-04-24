package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend/internal/api/handlers"
	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"backend/test"
)

var router *gin.Engine
var originalDB *sqlx.DB

// Define a variable for the upload function to make it easier to mock
var uploadResumeFunc = services.UploadResumeToS3

// MockDB for database testing
type MockDB struct {
	*sqlx.DB
	mock.Mock
}

// MockUploadResumeToS3 mocks the S3 upload functionality
func mockUploadResumeToS3(file *multipart.FileHeader, userID int, username string) (string, error) {
	return fmt.Sprintf("https://example.com/resumes/%d_%s.pdf", userID, username), nil
}

// MockFailedUploadResumeToS3 mocks a failed S3 upload
func mockFailedUploadResumeToS3(file *multipart.FileHeader, userID int, username string) (string, error) {
	return "", fmt.Errorf("S3 upload failed")
}

// CreateHandleFormSubmissionWithCustomUploader creates a handler with a custom uploader
func CreateHandleFormSubmissionWithCustomUploader(uploader func(*multipart.FileHeader, int, string) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Save the original function
		originalUploader := uploadResumeFunc
		// Replace with our custom function
		uploadResumeFunc = uploader
		// Restore after the handler completes
		defer func() { uploadResumeFunc = originalUploader }()

		// Call the actual handler
		handlers.HandleFormSubmission(c)
	}
}

func TestMain(m *testing.M) {
	log.Println("Initializing test environment...")

	router = gin.Default()

	// Register only required routes to avoid import cycles
	router.GET("/api/forms/:form_uuid/submissions", handlers.GetFormSubmissions)
	router.POST("/api/forms/:form_uuid/submit", handlers.HandleFormSubmission)

	setupTestDB() // Set up test db and data
	log.Println("Test database setup complete")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupTestDB() {

	db := test.SetupTestDB()

	// Step 1: Delete job submissions first (depends on jobs)
	_, err := db.Exec("DELETE FROM job_submissions")
	if err != nil {
		log.Fatal("Failed to clear job_submissions:", err)
	}

	// Step 2: Delete form template with specific form UUID
	_, err = db.Exec("DELETE FROM form_templates WHERE form_template_id = '4d9a4320-f1d1-43f2-8477-edd07f557442'")
	if err != nil {
		log.Fatal("Failed to clear form_templates:", err)
	}

	// Step 3: Delete jobs (depends on users)
	_, err = db.Exec("DELETE FROM jobs WHERE user_id IN (1001, 1002)")
	if err != nil {
		log.Fatal("Failed to clear jobs:", err)
	}

	// Step 4: Delete users
	_, err = db.Exec("DELETE FROM users WHERE id IN (1001, 1002)")
	if err != nil {
		log.Fatal("Failed to clear test users:", err)
	}

	// 🔹 Insert test users (Ensure unique ID)
	_, err = db.Exec(`
		INSERT INTO users (id, username, email, password_hash, role, company_name) 
		VALUES (1001, 'John Doe', 'john@example.com', 'hashedpassword', 'HR', 'Test Company'),
			   (1002, 'Jane Doe', 'jane@example.com', 'hashedpassword', 'HR', 'Test Company')`)
	if err != nil {
		log.Fatal("Failed to insert test users:", err)
	}

	// 🔹 Insert test jobs
	_, err = db.Exec(`
		INSERT INTO jobs (job_id, user_id, job_title, job_description, job_status, skills_required, attributes)
		VALUES ('5678', 1001, 'Software Engineer', 'Develop applications in Go', 'active', '{"Go", "AWS"}', '{}'::jsonb)`)
	if err != nil {
		log.Fatal("Failed to insert test jobs:", err)
	}

	// 🔹 Insert test form templates
	_, err = db.Exec(`
		INSERT INTO form_templates (form_template_id, user_id, fields)
		VALUES (
			'4d9a4320-f1d1-43f2-8477-edd07f557442',
			1001,
			'{"skills": "array", "location": "string", "experience": "string"}'::jsonb
		)`)
	if err != nil {
		log.Fatal("Failed to insert test form templates:", err)
	}

	// 🔹 Insert test job submissions
	_, err = db.Exec(`
		INSERT INTO job_submissions (form_uuid, job_id, username, email, form_data, skills, resume_url, ats_score, status)
		VALUES
			('4d9a4320-f1d1-43f2-8477-edd07f557442', '5678', 'John Doe', 'john@example.com', '{"experience":"5 years","location":"Remote","skills":["Go","AWS"]}', ARRAY['Go', 'AWS'], 'https://s3.aws.com/resume1.pdf', 95, 'applied'),
			('4d9a4320-f1d1-43f2-8477-edd07f557442', '5678', 'Jane Doe', 'jane@example.com', '{"experience":"3 years","location":"Onsite","skills":["Java","Docker"]}', ARRAY['Java', 'Docker'], 'https://s3.aws.com/resume2.pdf', 88, 'applied');
	`)

	if err != nil {
		log.Fatal("Failed to insert test job submissions:", err)
	}

	log.Println("Test database setup complete")
}

func TestGetFormSubmissions(t *testing.T) {
	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/forms/:form_uuid/submissions", handlers.GetFormSubmissions)

	// Setup mock database
	setupTestDB()

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/forms/4d9a4320-f1d1-43f2-8477-edd07f557442/submissions?sort_by=ats_score&limit=2", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Parse response
	var response map[string][]models.JobSubmission
	err := json.Unmarshal(recorder.Body.Bytes(), &response)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Len(t, response["submissions"], 2)
	assert.GreaterOrEqual(t, response["submissions"][0].ATSScore, response["submissions"][1].ATSScore) // Ensure sorting by ATS score
}

func TestGetFormSubmissions_InvalidUUID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/forms/123e4567-e89b-12d3-a456-426614174000/submissions", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	//  Expecting `400 Bad Request`
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request for invalid UUID")
}

// createTestMultipartRequest creates a multipart request for testing file uploads
func createTestMultipartRequest(t *testing.T, url string, formValues map[string]string, fileFieldName, fileName string) (*http.Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add form fields
	for key, value := range formValues {
		if err := w.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	// Add the file
	fileWriter, err := w.CreateFormFile(fileFieldName, fileName)
	if err != nil {
		return nil, err
	}

	// Write dummy content to the file
	fileContent := []byte("This is a test resume content")
	if _, err := fileWriter.Write(fileContent); err != nil {
		return nil, err
	}

	// Close the writer
	if err := w.Close(); err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}

func TestHandleFormSubmission_Success(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with a unique email to avoid constraint violation
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1003",
		"username":  "Test User",
		"email":     "test_unique_email@example.com", // Unique email
		"form_data": `{"experience":"2 years", "location":"Remote", "skills":["Go", "AWS"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with custom handler
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response
	if recorder.Code != http.StatusOK {
		t.Logf("Response Body: %s", recorder.Body.String())
	}
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for success message
	assert.Contains(t, response, "message")
	assert.Equal(t, "Submission received successfully", response["message"])

	// Check ATS score was calculated
	assert.Contains(t, response, "ats_score")
}

func TestHandleFormSubmission_DuplicateApplication(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values for a user who already applied (user_id 1001, job_id 5678)
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1001", // Existing user who already applied
		"username":  "John Doe",
		"email":     "john@example.com",
		"form_data": `{"experience":"6 years", "location":"Remote", "skills":["Go", "AWS", "Docker"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with custom handler
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should fail with duplicate application error
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "User has already applied for this job", response["error"])
}

func TestHandleFormSubmission_InvalidFormData(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with invalid form_data JSON
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1004",
		"username":  "Invalid User",
		"email":     "invalid@example.com",
		"form_data": `{"experience":"3 years", "location":, "skills":["Go"]}`, // Invalid JSON
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with route
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should fail with invalid form data
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid form data format", response["error"])
}

func TestHandleFormSubmission_InvalidUUID(t *testing.T) {
	// Set up an invalid form UUID
	formUUID := "invalid-uuid"

	// Create form values
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1005",
		"username":  "UUID Test User",
		"email":     "uuid@example.com",
		"form_data": `{"experience":"1 year", "location":"Remote", "skills":["Go"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with route
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should fail with invalid UUID
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid form_uuid", response["error"])
}

func TestHandleFormSubmission_S3UploadFailure(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1006",
		"username":  "S3 Test User",
		"email":     "s3test@example.com",
		"form_data": `{"experience":"4 years", "location":"Hybrid", "skills":["Go", "AWS"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with route
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockFailedUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should fail with S3 upload error
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Failed to upload resume", response["error"])
}

func TestHandleFormSubmission_MissingRequiredFields(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with missing required fields
	formValues := map[string]string{
		// Missing job_id
		"user_id":   "1007",
		"username":  "Missing Fields User",
		"email":     "missing@example.com",
		"form_data": `{"experience":"2 years", "location":"Remote", "skills":["Go"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with route
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should fail with binding error
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid submission data", response["error"])
}

func TestHandleFormSubmission_InvalidUserID(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with invalid user_id format (not a number)
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "not-a-number", // Invalid user_id format
		"username":  "Invalid ID User",
		"email":     "invalid_id@example.com",
		"form_data": `{"experience":"3 years", "location":"Remote", "skills":["Go"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with route
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should fail with invalid user ID format
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message about invalid format
	assert.Contains(t, response, "error")
	errorMsg, ok := response["error"].(string)
	assert.True(t, ok, "Error message should be a string")
	assert.Equal(t, "Invalid User ID format", errorMsg)
}

// Add a new test for the case where user_id is not provided (should now be valid)
func TestHandleFormSubmission_MissingUserID(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values without user_id
	formValues := map[string]string{
		"job_id": "5678",
		// No user_id provided
		"username":  "No User ID",
		"email":     "no_userid@example.com",
		"form_data": `{"experience":"3 years", "location":"Remote", "skills":["Go"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Add test header to avoid actual S3 uploads
	req.Header.Set("X-Test-Mode", "true")

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with route
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Assert response - should succeed
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify success message and that an ID was generated
	assert.Contains(t, response, "message")
	assert.Equal(t, "Application submitted successfully", response["message"])
	assert.Contains(t, response, "id")
	assert.Contains(t, response, "ats_score")
}

// New test cases to improve coverage

// TestGetFormSubmissions_DifferentSortOptions tests different sorting options
func TestGetFormSubmissions_DifferentSortOptions(t *testing.T) {
	// Set up test cases for different sort options
	sortOptions := []string{"created_at", "ats_score"}

	for _, sortOption := range sortOptions {
		// Create a test request with the sort option
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/forms/4d9a4320-f1d1-43f2-8477-edd07f557442/submissions?sort_by=%s&limit=2", sortOption), nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Assertions
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Parse response
		var response map[string][]models.JobSubmission
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "submissions")
	}
}

// TestGetFormSubmissions_DifferentLimits tests different limit values
func TestGetFormSubmissions_DifferentLimits(t *testing.T) {
	// Test with various limit values
	limits := []string{"1", "5", "invalid"}

	for _, limit := range limits {
		// Create a test request with the limit option
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/forms/4d9a4320-f1d1-43f2-8477-edd07f557442/submissions?limit=%s", limit), nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Assertions - should succeed even with invalid limit (will use default)
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Parse response
		var response map[string][]models.JobSubmission
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "submissions")
	}
}

// TestHandleFormSubmission_JobNotFound tests the case where the job doesn't exist
func TestHandleFormSubmission_JobNotFound(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with non-existent job ID
	formValues := map[string]string{
		"job_id":    "9999", // Non-existent job ID
		"user_id":   "1008",
		"username":  "Job Not Found User",
		"email":     "jobnotfound@example.com",
		"form_data": `{"experience":"2 years", "location":"Remote", "skills":["Go", "AWS"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with custom handler
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Check if we get an appropriate response for job not found
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// The exact error might vary but should be database-related
	assert.Contains(t, response, "error")
}

// TestHandleFormSubmission_InvalidEmail tests submission with invalid email format
func TestHandleFormSubmission_InvalidEmail(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with invalid email
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1009",
		"username":  "Invalid Email User",
		"email":     "invalid-email", // Invalid email format
		"form_data": `{"experience":"2 years", "location":"Remote", "skills":["Go", "AWS"]}`,
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with custom handler
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Should fail validation
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should be a validation error
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid submission data", response["error"])
}

// TestHandleFormSubmission_EmptySkills tests submission with empty skills array but still valid JSON
func TestHandleFormSubmission_EmptySkills(t *testing.T) {
	// Set up a valid form UUID
	formUUID := "4d9a4320-f1d1-43f2-8477-edd07f557442"

	// Create form values with empty skills array
	formValues := map[string]string{
		"job_id":    "5678",
		"user_id":   "1010",
		"username":  "Empty Skills User",
		"email":     "emptyskills@example.com",
		"form_data": `{"experience":"2 years", "location":"Remote", "skills":[]}`, // Empty skills array
	}

	// Create test request
	req, err := createTestMultipartRequest(t, "/api/forms/"+formUUID+"/submit", formValues, "resume", "test_resume.pdf")
	assert.NoError(t, err)

	// Set URL parameter
	req.URL.Path = "/api/forms/" + formUUID + "/submit"

	// Create recorder and serve request
	recorder := httptest.NewRecorder()

	// Create router with custom handler
	testRouter := gin.Default()
	testRouter.POST("/api/forms/:form_uuid/submit", CreateHandleFormSubmissionWithCustomUploader(mockUploadResumeToS3))
	testRouter.ServeHTTP(recorder, req)

	// Since the skills column has a NOT NULL constraint, we expect an error
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should be a database error related to NOT NULL constraint
	assert.Contains(t, response, "error")
	assert.Equal(t, "Failed to store form submission", response["error"])
	assert.Contains(t, response["details"], "null value in column \"skills\"")
}

// TestGetFormSubmissions_FormDoesNotExist tests retrieving submissions for a non-existent form (but valid UUID format)
func TestGetFormSubmissions_FormDoesNotExist(t *testing.T) {
	// Create a valid UUID that doesn't exist in the database
	req, _ := http.NewRequest("GET", "/api/forms/12345678-1234-1234-1234-123456789abc/submissions", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Should return 400 Bad Request because form doesn't exist
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid form_uuid, form does not exist", response["error"])
}

// TestGetFormSubmissions_DateFilter tests the date filter functionality
func TestGetFormSubmissions_DateFilter(t *testing.T) {
	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/forms/:form_uuid/submissions", handlers.GetFormSubmissions)

	// Setup mock database
	setupTestDB()

	// Test cases for different date filters
	dateFilters := []string{"all", "today", "2025-03-03", "invalid-date"}

	for _, dateFilter := range dateFilters {
		// Create a test request with the date filter
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/forms/4d9a4320-f1d1-43f2-8477-edd07f557442/submissions?date=%s", dateFilter), nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Assertions
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Parse response
		var response map[string][]models.JobSubmission
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "submissions")

		// For "all" and valid date, we should get submissions
		if dateFilter == "all" || dateFilter == "today" || dateFilter == "2025-03-03" {
			assert.GreaterOrEqual(t, len(response["submissions"]), 1, "Expected at least one submission for date filter: %s", dateFilter)
		}
	}
}
