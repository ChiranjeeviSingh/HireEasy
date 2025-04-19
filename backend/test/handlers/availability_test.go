package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/api/handlers"
	"backend/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestCreateAvailabilityH tests the POST /api/availability endpoint
func TestCreateAvailabilityH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Drop all tables and recreate them to ensure a completely clean state
	_, err := db.Exec("DELETE FROM availabilities; DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Insert a test user with a unique email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	var userID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('createavail@example.com', $1, 'createavail', 'Interviewer', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&userID)
	assert.NoError(t, err)

	// Create a profile for this user (assuming a profile is required)
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID, "Senior Interviewer", 8, `{"Go", "System Design", "Algorithms"}`, "1234567890")
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/availability", handlers.CreateAvailabilityH)

	// Get tomorrow's date at 10:00 AM
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 10, 0, 0, 0, now.Location())

	requestBody := map[string]interface{}{
		"date":      tomorrow.Format("2006-01-02"),
		"from_time": tomorrow.Format("15:04:05"),
		"to_time":   tomorrow.Add(1 * time.Hour).Format("15:04:05"),
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/availability", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestCreateAvailabilityH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Basic validation of response fields
	assert.NotNil(t, response["id"])
	assert.Equal(t, requestBody["date"], response["date"])
	assert.Equal(t, requestBody["from_time"], response["from_time"])
	assert.Equal(t, requestBody["to_time"], response["to_time"])
}

// TestDeleteAvailabilityH tests the DELETE /api/availability/:id endpoint
func TestDeleteAvailabilityH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Clean up database
	_, err := db.Exec("DELETE FROM availabilities; DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Insert a test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	var userID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('deleteavail@example.com', $1, 'deleteavail', 'Interviewer', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&userID)
	assert.NoError(t, err)

	// Create a profile for this user
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID, "Interviewer", 5, `{"Java", "Spring"}`, "9876543210")
	assert.NoError(t, err)

	// Create an availability slot
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 14, 0, 0, 0, now.Location())

	var availabilityID int
	err = db.QueryRow(`INSERT INTO availabilities (user_id, date, from_time, to_time)
	                 VALUES ($1, $2, $3, $4)
					 RETURNING id`,
		userID, tomorrow.Format("2006-01-02"), tomorrow.Format("15:04:05"), tomorrow.Add(1*time.Hour).Format("15:04:05")).Scan(&availabilityID)
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.DELETE("/api/availability/:id", handlers.DeleteAvailabilityH)

	// Delete the availability slot
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/availability/%d", availabilityID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestDeleteAvailabilityH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify the success message
	assert.Contains(t, response, "message")
	assert.Contains(t, response["message"], "deleted")

	// Verify the availability slot is actually deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM availabilities WHERE id = $1", availabilityID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

// TestGetMyAvailabilityH tests the GET /api/availability/me endpoint
func TestGetMyAvailabilityH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Clean up database
	_, err := db.Exec("DELETE FROM availabilities; DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Insert a test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	var userID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('myavail@example.com', $1, 'myavail', 'Interviewer', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&userID)
	assert.NoError(t, err)

	// Create a profile for this user
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID, "Tech Lead", 7, `{"Architecture", "Leadership"}`, "5551234567")
	assert.NoError(t, err)

	// Create multiple availability slots
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, now.Location())

	// First slot: 9-10 AM tomorrow
	_, err = db.Exec(`INSERT INTO availabilities (user_id, date, from_time, to_time)
	                 VALUES ($1, $2, $3, $4)`,
		userID, tomorrow.Format("2006-01-02"), tomorrow.Format("15:04:05"), tomorrow.Add(1*time.Hour).Format("15:04:05"))
	assert.NoError(t, err)

	// Second slot: 2-3 PM tomorrow
	_, err = db.Exec(`INSERT INTO availabilities (user_id, date, from_time, to_time)
	                 VALUES ($1, $2, $3, $4)`,
		userID, tomorrow.Format("2006-01-02"), tomorrow.Add(5*time.Hour).Format("15:04:05"), tomorrow.Add(6*time.Hour).Format("15:04:05"))
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/api/availability/me", handlers.GetMyAvailabilityH)

	// Get my availability slots
	req, _ := http.NewRequest("GET", "/api/availability/me", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestGetMyAvailabilityH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have two availability slots
	assert.Len(t, response, 2)

	// Verify both slots belong to the user
	for _, slot := range response {
		assert.Equal(t, float64(userID), slot["user_id"])
	}

	// Verify we have all the fields needed
	for _, slot := range response {
		assert.Contains(t, slot, "id")
		assert.Contains(t, slot, "user_id")
		assert.Contains(t, slot, "date")
		assert.Contains(t, slot, "from_time")
		assert.Contains(t, slot, "to_time")
	}
}

// TestGetUserAvailabilityH tests the GET /api/availability/user/:user_name endpoint
func TestGetUserAvailabilityH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Clean up database
	_, err := db.Exec("DELETE FROM availabilities; DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Insert a test user with specific username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	username := "specificuser"
	var userID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('specific@example.com', $1, $2, 'Interviewer', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword), username).Scan(&userID)
	assert.NoError(t, err)

	// Create a profile for this user
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID, "Principal Engineer", 10, `{"System Design", "Algorithms"}`, "1112223333")
	assert.NoError(t, err)

	// Create availability slots
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 11, 0, 0, 0, now.Location())

	// Slot: 11-12 tomorrow
	_, err = db.Exec(`INSERT INTO availabilities (user_id, date, from_time, to_time)
	                 VALUES ($1, $2, $3, $4)`,
		userID, tomorrow.Format("2006-01-02"), tomorrow.Format("15:04:05"), tomorrow.Add(1*time.Hour).Format("15:04:05"))
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	// No authentication needed for this endpoint
	router.GET("/api/availability/user/:user_name", handlers.GetUserAvailabilityH)

	// Get user's availability by username
	req, _ := http.NewRequest("GET", "/api/availability/user/"+username, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestGetUserAvailabilityH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have one availability slot
	assert.Len(t, response, 1)

	if len(response) > 0 {
		// Verify slot details
		slot := response[0]
		assert.Equal(t, float64(userID), slot["user_id"])
		assert.Contains(t, slot, "date")
		assert.Contains(t, slot, "from_time")
		assert.Contains(t, slot, "to_time")
	}
}

// TestGetAllAvailabilityH tests the GET /api/availability endpoint
func TestGetAllAvailabilityH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Clean up database
	_, err := db.Exec("DELETE FROM availabilities; DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Insert multiple test users with availability slots
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// First user
	var userID1 int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('user1@example.com', $1, 'user1', 'Interviewer', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&userID1)
	assert.NoError(t, err)

	// Create a profile for first user
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID1, "Software Engineer", 3, `{"Java", "Spring"}`, "1112223333")
	assert.NoError(t, err)

	// Second user
	var userID2 int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('user2@example.com', $1, 'user2', 'Interviewer', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&userID2)
	assert.NoError(t, err)

	// Create a profile for second user
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID2, "Senior Developer", 7, `{"Python", "Django"}`, "4445556666")
	assert.NoError(t, err)

	// Create availability slots
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 10, 0, 0, 0, now.Location())

	// Slot for first user: 10-11 tomorrow
	_, err = db.Exec(`INSERT INTO availabilities (user_id, date, from_time, to_time)
	                 VALUES ($1, $2, $3, $4)`,
		userID1, tomorrow.Format("2006-01-02"), tomorrow.Format("15:04:05"), tomorrow.Add(1*time.Hour).Format("15:04:05"))
	assert.NoError(t, err)

	// Slot for second user: 2-3 PM tomorrow
	_, err = db.Exec(`INSERT INTO availabilities (user_id, date, from_time, to_time)
	                 VALUES ($1, $2, $3, $4)`,
		userID2, tomorrow.Format("2006-01-02"), tomorrow.Add(4*time.Hour).Format("15:04:05"), tomorrow.Add(5*time.Hour).Format("15:04:05"))
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	// Add authentication for a third user who is viewing all availability
	var viewerID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('viewer@example.com', $1, 'viewer', 'HR', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&viewerID)
	assert.NoError(t, err)

	router.Use(func(c *gin.Context) {
		c.Set("userID", viewerID)
		c.Next()
	})

	router.GET("/api/availability", handlers.GetAllAvailabilityH)

	// Get all available slots
	req, _ := http.NewRequest("GET", "/api/availability", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestGetAllAvailabilityH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have two availability slots
	assert.Len(t, response, 2)

	// Check that we have the expected fields
	for _, slot := range response {
		assert.Contains(t, slot, "user_id")
		assert.Contains(t, slot, "date")
		assert.Contains(t, slot, "from_time")
		assert.Contains(t, slot, "to_time")
	}
}
