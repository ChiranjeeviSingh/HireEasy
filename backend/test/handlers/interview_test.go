package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"backend/internal/api/handlers"
	"backend/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/lib/pq"
	"time"
)

// The db variable is already defined in another test file in the handlers_test package
// var db = test.SetupTestDB()

func TestCreateInterview(t *testing.T) {
	db := test.SetupTestDB()
	
	// Clear existing data
	test.CleanupTestDB(db)

	// db.Exec("TRUNCATE users, jobs, job_submissions, availabilities, interviews CASCADE")
	
	// Create HR user and interviewer
	hrUserID, _ := test.InsertTestUser(db)
	interviewerID, _ := test.InsertTestInterviewerUser(db)
	
	// Insert test job
	jobID := "J12349"
	var jobIdNum int
	err := db.QueryRow(`INSERT INTO jobs (job_id, user_id, job_title, job_description, skills_required) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		jobID, hrUserID, "Test Job", "Test Description", pq.Array([]string{"Go"})).Scan(&jobIdNum)
	assert.NoError(t, err)
	
	// Insert test form template (needed for job submission)
	var formIdNum int
	err = db.QueryRow(`INSERT INTO form_templates (form_template_id, user_id, fields) 
		VALUES ($1, $2, $3) RETURNING id`,
		"template123", hrUserID, json.RawMessage(`{}`)).Scan(&formIdNum)
	assert.NoError(t, err)

	// Insert test application form
	_, err = db.Exec(`INSERT INTO application_form (form_uuid, job_id, form_id, status) 
		VALUES ($1, $2, $3, $4)`,
		"123e4567-e89b-12d3-a456-426614174000", jobIdNum, formIdNum, "active")
	assert.NoError(t, err)

	// Insert test job submission
	_, err = db.Exec(`INSERT INTO job_submissions (form_uuid, job_id, username, email, form_data, resume_url) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"123e4567-e89b-12d3-a456-426614174000", jobID, "test candidate", "test_candidate@test.com", json.RawMessage(`{}`), "resume.pdf")
	assert.NoError(t, err)
	
	// Insert test availability
	var availID int
	err = db.QueryRow(`INSERT INTO availabilities (user_id, date, from_time, to_time) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		interviewerID, "2024-01-01", "10:00:00", "11:00:00").Scan(&availID)
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", hrUserID)
		c.Next()
	})
	router.POST("/api/interviews", handlers.CreateInterviewH)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedCode   int
	}{
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"job_id": jobID,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid job_id",
			requestBody: map[string]interface{}{
				"job_id": 123,
				"job_submission_id": 1,
				"interviewer_user_id": interviewerID,
				"availability_id": availID,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Valid request",
			requestBody: map[string]interface{}{
				"job_id": jobID,
				"job_submission_id": 1,
				"interviewer_user_id": interviewerID,
				"availability_id": availID,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/interviews", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedCode, resp.Code)

			if tc.expectedCode == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.requestBody["job_id"], response["job_id"])
				assert.Equal(t, float64(hrUserID), response["hr_user_id"])
				assert.Equal(t, "scheduled", response["status"])
				assert.NotEmpty(t, response["created_at"])
				assert.NotEmpty(t, response["updated_at"]) 
			}
		})
	}
}

func TestDeleteInterviewH(t *testing.T) {
	db := test.SetupTestDB()
	
	// Clear existing data
	test.CleanupTestDB(db)
	
	// Create HR user and interviewer
	hrUserID, _ := test.InsertTestUser(db)
	interviewerID, _ := test.InsertTestInterviewerUser(db)
	
	// Insert test job
	jobID := "J12349"
	var jobIdNum int
	err := db.QueryRow(`INSERT INTO jobs (job_id, user_id, job_title, job_description, skills_required) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		jobID, hrUserID, "Test Job", "Test Description", pq.Array([]string{"Go"})).Scan(&jobIdNum)
	assert.NoError(t, err)
	
	// Insert test job submission
	_, err = db.Exec(`INSERT INTO job_submissions (form_uuid, job_id, username, email, form_data, resume_url) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"123e4567-e89b-12d3-a456-426614174000", jobID, "test candidate", "test_candidate@test.com", json.RawMessage(`{}`), "resume.pdf")
	assert.NoError(t, err)
	
	// Insert test availability
	var availID int
	err = db.QueryRow(`INSERT INTO availabilities (user_id, date, from_time, to_time) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		interviewerID, "2024-01-01", "10:00:00", "11:00:00").Scan(&availID)
	assert.NoError(t, err)

	// Insert test interview
	var interviewID int
	err = db.QueryRow(`INSERT INTO interviews (job_id, hr_user_id, job_submission_id, interviewer_user_id, availability_id, status) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		jobID, hrUserID, 1, interviewerID, availID, "scheduled").Scan(&interviewID)
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", hrUserID)
		c.Next()
	})
	router.DELETE("/api/interviews/:id", handlers.DeleteInterviewH)

	testCases := []struct {
		name         string
		interviewID  int
		expectedCode int
	}{
		{
			name:         "Non-existent interview ID",
			interviewID:  99999,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Valid interview deletion",
			interviewID:  interviewID,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/interviews/%d", tc.interviewID), nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedCode, resp.Code)

			if tc.expectedCode == http.StatusOK {
				// Verify interview was actually deleted
				var count int
				err := db.QueryRow("SELECT COUNT(*) FROM interviews WHERE id = $1", tc.interviewID).Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 0, count)
			}
		})
	}
}

func TestSubmitFeedbackH(t *testing.T) {
	db := test.SetupTestDB()
	
	// Clear existing data
	test.CleanupTestDB(db)
	
	// Create HR user and interviewer
	hrUserID, _ := test.InsertTestUser(db)
	interviewerID, _ := test.InsertTestInterviewerUser(db)
	
	// Insert test job
	jobID := "J12349"
	var jobIdNum int
	err := db.QueryRow(`INSERT INTO jobs (job_id, user_id, job_title, job_description, skills_required) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		jobID, hrUserID, "Test Job", "Test Description", pq.Array([]string{"Go"})).Scan(&jobIdNum)
	assert.NoError(t, err)
	
	// Insert test job submission
	_, err = db.Exec(`INSERT INTO job_submissions (form_uuid, job_id, username, email, form_data, resume_url) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"123e4567-e89b-12d3-a456-426614174000", jobID, "test candidate", "test_candidate@test.com", json.RawMessage(`{}`), "resume.pdf")
	assert.NoError(t, err)
	
	// Insert test availability
	var availID int
	err = db.QueryRow(`INSERT INTO availabilities (user_id, date, from_time, to_time) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		interviewerID, "2024-01-01", "10:00:00", "11:00:00").Scan(&availID)
	assert.NoError(t, err)

	// Insert test interview
	var interviewID int
	err = db.QueryRow(`INSERT INTO interviews (job_id, hr_user_id, job_submission_id, interviewer_user_id, availability_id, status) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		jobID, hrUserID, 1, interviewerID, availID, "pending_feedback").Scan(&interviewID)
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", interviewerID)
		c.Next()
	})
	router.POST("/api/interviews/:id/feedback", handlers.SubmitFeedbackH)

	testCases := []struct {
		name         string
		interviewID  int
		requestBody  map[string]interface{}
		expectedCode int
	}{
		{
			name:        "Non-existent interview ID",
			interviewID: 99999,
			requestBody: map[string]interface{}{
				"verdict":  "passed",
				"feedback": "Great candidate",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "Valid feedback submission",
			interviewID: interviewID,
			requestBody: map[string]interface{}{
				"verdict":  "passed",
				"feedback": "Great candidate",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/api/interviews/%d/feedback", tc.interviewID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedCode, resp.Code)

			if tc.expectedCode == http.StatusOK {
				// Verify feedback was actually saved
				var status, verdict string
				err := db.QueryRow("SELECT status, verdict FROM interviews WHERE id = $1", tc.interviewID).Scan(&status, &verdict)
				assert.NoError(t, err)
				assert.Equal(t, "passed", verdict)
				assert.Equal(t, "completed", status)
			}
		})
	}
}

func TestListAllInterviewsH(t *testing.T) {
	db := test.SetupTestDB()
	
	// Clear existing data
	test.CleanupTestDB(db)
	
	// Create HR user and two interviewers
	hrUserID, _ := test.InsertTestUser(db)
	interviewer1ID, _ := test.InsertTestInterviewerUser(db)
	interviewer2ID, _ := test.InsertTestInterviewerUser(db)

	// Insert test job
	jobID := "J12349"
	var jobIdNum int
	err := db.QueryRow(`INSERT INTO jobs (job_id, user_id, job_title, job_description, skills_required) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		jobID, hrUserID, "Test Job", "Test Description", pq.Array([]string{"Go"})).Scan(&jobIdNum)
	assert.NoError(t, err)

	// Insert test job submissions
	_, err = db.Exec(`INSERT INTO job_submissions (form_uuid, job_id, username, email, form_data, resume_url) 
		VALUES ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12)`,
		"uuid1", jobID, "candidate1", "candidate1@test.com", json.RawMessage(`{}`), "resume1.pdf",
		"uuid2", jobID, "candidate2", "candidate2@test.com", json.RawMessage(`{}`), "resume2.pdf")
	assert.NoError(t, err)

	// Insert availabilities - one past, one future
	var pastAvailID, futureAvailID int
	pastDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	
	err = db.QueryRow(`INSERT INTO availabilities (user_id, date, from_time, to_time) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		interviewer1ID, pastDate, "10:00:00", "11:00:00").Scan(&pastAvailID)
	assert.NoError(t, err)

	err = db.QueryRow(`INSERT INTO availabilities (user_id, date, from_time, to_time) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		interviewer2ID, futureDate, "14:00:00", "15:00:00").Scan(&futureAvailID)
	assert.NoError(t, err)

	// Insert interviews
	_, err = db.Exec(`INSERT INTO interviews (job_id, hr_user_id, job_submission_id, interviewer_user_id, availability_id, status) 
		VALUES ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12)`,
		jobID, hrUserID, 1, interviewer1ID, pastAvailID, "scheduled",
		jobID, hrUserID, 2, interviewer2ID, futureAvailID, "scheduled")
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		userID      int
		queryParams string
		expectedLen int
		checkStatus bool
	}{
		{
			name:         "List all interviews for HR",
			userID:      hrUserID,
			queryParams: "",
			expectedLen: 2,
			checkStatus: true,
		},
		{
			name:         "Filter by job_submission_id",
			userID:      hrUserID,
			queryParams: "?job_submission_id=1",
			expectedLen: 1,
			checkStatus: false,
		},
		{
			name:         "Filter by date range",
			userID:      hrUserID,
			queryParams: fmt.Sprintf("?from_date=%s&to_date=%s", pastDate, pastDate),
			expectedLen: 1,
			checkStatus: false,
		},
		{
			name:         "List interviews for interviewer1",
			userID:      interviewer1ID,
			queryParams: "",
			expectedLen: 1,
			checkStatus: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			router := test.SetupTestRouter()
	

			router.Use(func(c *gin.Context) {
				c.Set("userID", tc.userID)
				c.Next()
			})

			router.GET("/api/interviews", handlers.ListAllInterviewsH)

			req, _ := http.NewRequest("GET", "/api/interviews"+tc.queryParams, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)

			var response []map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, tc.expectedLen)

			if tc.checkStatus {
				for _, interview := range response {
					availID := int(interview["availability_id"].(float64))
					if availID == pastAvailID {
						assert.Equal(t, "pending_feedback", interview["status"])
					} else {
						assert.Equal(t, "scheduled", interview["status"])
					}
				}
			}
		})
	}
}