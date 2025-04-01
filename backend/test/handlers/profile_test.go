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
	"golang.org/x/crypto/bcrypt"
)

// The db variable is already defined in another test file in the handlers_test package
// var db = test.SetupTestDB()

func TestCreateProfileH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Drop all tables and recreate them to ensure a completely clean state
	_, err := db.Exec("DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Insert a test user with a unique email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	var userID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('createprofile@example.com', $1, 'createprofile', 'HR', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword)).Scan(&userID)
	assert.NoError(t, err)

	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/api/profiles", handlers.CreateProfileH)

	requestBody := map[string]interface{}{
		"job_title":           "Dev",
		"years_of_experience": 5,
		"areas_of_expertise":  []string{"Go", "JS", "SQL"},
		"phone_number":        "1234567890",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/profiles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestCreateProfileH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Don't check for userID in response since it's not included
	assert.Equal(t, "Dev", response["job_title"])
	assert.Equal(t, float64(5), response["years_of_experience"])
	assert.Equal(t, "1234567890", response["phone_number"])

	// Verify areas_of_expertise array
	expertise, ok := response["areas_of_expertise"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, expertise, 3)
	assert.Contains(t, expertise, "Go")
	assert.Contains(t, expertise, "JS")
	assert.Contains(t, expertise, "SQL")
}

func TestUpdateMyProfileH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// No need to clean up since we want the profile created in TestCreateProfileH

	// Insert a test user first and create a profile
	userID, _ := test.InsertTestUser(db)

	// Setup router and middleware for update
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.PUT("/api/profiles", handlers.UpdateMyProfileH)

	// Skip creating a profile since it already exists from TestCreateProfileH
	// Update the existing profile
	updateBody := map[string]interface{}{
		"job_title":           "Dev Sr",
		"years_of_experience": 5,
		"areas_of_expertise":  []string{"Go", "JS", "React", "Docker"},
		"phone_number":        "9876543210",
	}
	jsonBody, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest("PUT", "/api/profiles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestUpdateMyProfileH Update Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Don't check for userID in response since it's not included
	assert.Equal(t, "Dev Sr", response["job_title"])
	assert.Equal(t, float64(5), response["years_of_experience"])
	assert.Equal(t, "9876543210", response["phone_number"])

	// Verify areas_of_expertise array was updated
	expertise, ok := response["areas_of_expertise"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, expertise, 4) // Now has 4 areas of expertise
	assert.Contains(t, expertise, "Go")
	assert.Contains(t, expertise, "JS")
	assert.Contains(t, expertise, "React")  // New area
	assert.Contains(t, expertise, "Docker") // New area
}

func TestGetMyProfileH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// No need to clean up since we want the profile from the previous tests

	// Insert a test user
	userID, _ := test.InsertTestUser(db)

	// Setup router and middleware for get
	router := test.SetupTestRouter()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/api/profiles/me", handlers.GetMyProfileH)

	// Skip creating a profile since it already exists from TestUpdateMyProfileH
	// Now get the profile
	req, _ := http.NewRequest("GET", "/api/profiles/me", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestGetMyProfileH Get Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(userID), response["user_id"])
	assert.Equal(t, "Dev Sr", response["job_title"])             // This is what was set in TestUpdateMyProfileH
	assert.Equal(t, float64(5), response["years_of_experience"]) // This is what was set in TestUpdateMyProfileH
	assert.Equal(t, "9876543210", response["phone_number"])      // This is what was set in TestUpdateMyProfileH

	// Verify areas_of_expertise array - we expect the values from TestUpdateMyProfileH
	expertise, ok := response["areas_of_expertise"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, expertise, 4) // We expect 4 from the update test
	// The areas are lowercase in the response due to the service converting them
	assert.Contains(t, expertise, "go")
	assert.Contains(t, expertise, "js")
	assert.Contains(t, expertise, "react")
	assert.Contains(t, expertise, "docker")
}

func TestGetUserProfileH(t *testing.T) {
	// Initialize db at the start of the test
	db := test.SetupTestDB()

	// Drop all tables and recreate them to ensure a completely clean state
	_, err := db.Exec("DELETE FROM profiles; DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users;")
	assert.NoError(t, err)

	// Create a new user with a unique profile for this test
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// Username to use for retrieval
	username := "userprofiletest"

	// Insert user with the username
	var userID int
	err = db.QueryRow(`INSERT INTO users (email, password_hash, username, role, company_name) 
	                  VALUES ('userprofile@example.com', $1, $2, 'HR', 'Test Company') 
					  RETURNING id`,
		string(hashedPassword), username).Scan(&userID)
	assert.NoError(t, err)

	// Create a profile for this user
	_, err = db.Exec(`INSERT INTO profiles 
	                (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
					VALUES ($1, $2, $3, $4, $5)`,
		userID, "Data Scientist", 3, `{"Python", "R", "ML"}`, "5556667777")
	assert.NoError(t, err)

	// Set up router for getting a user profile by username
	router := test.SetupTestRouter()
	// No need for authentication middleware since we're retrieving by username
	router.GET("/api/profiles/user/:user_name", handlers.GetUserProfileH)

	// Make request to get the profile by username
	req, _ := http.NewRequest("GET", "/api/profiles/user/"+username, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Print the response for debugging
	fmt.Printf("TestGetUserProfileH Response: %s\n", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify response content
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(userID), response["user_id"])
	assert.Equal(t, "Data Scientist", response["job_title"])
	assert.Equal(t, float64(3), response["years_of_experience"])
	assert.Equal(t, "5556667777", response["phone_number"])

	// Verify areas_of_expertise array
	expertise, ok := response["areas_of_expertise"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, expertise, 3)
	// The areas are returned with their original case
	assert.Contains(t, expertise, "Python")
	assert.Contains(t, expertise, "R")
	assert.Contains(t, expertise, "ML")
}
