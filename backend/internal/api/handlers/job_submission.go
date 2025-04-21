package handlers

import (
	"log"
	"net/http"
	"time"

	"backend/internal/database"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// HandleFormSubmission handles job application form submissions
func HandleFormSubmission(c *gin.Context) {
	// Get job ID from URL
	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job ID is required"})
		return
	}

	// Create form submission service
	db := database.GetDB()
	formService := services.NewFormSubmissionService(db)

	// Process form submission
	submission, err := formService.HandleFormSubmission(c)
	if err != nil {
		log.Printf("Error handling form submission: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job application submitted successfully",
		"data": gin.H{
			"job_id":     submission.JobID,
			"ats_score":  submission.ATSScore,
			"resume_url": submission.ResumeURL,
			"status":     submission.Status,
		},
	})
}

// GetFormSubmissions retrieves job submissions for a specific job
func GetFormSubmissions(c *gin.Context) {
	// Get job ID from URL
	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job ID is required"})
		return
	}

	log.Printf("Getting submissions for job ID: %s", jobID)

	// Get submissions from database
	db := database.GetDB()

	// First check if the job exists
	var jobExists bool
	err := db.Get(&jobExists, "SELECT EXISTS(SELECT 1 FROM jobs WHERE job_id = $1)", jobID)
	if err != nil {
		log.Printf("Error checking if job exists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify job"})
		return
	}

	if !jobExists {
		log.Printf("Job with ID %s not found", jobID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	
	// Get status filter from query parameter
	status := c.Query("status")
	
	// Construct query based on whether status filter is provided
	var query string
	var args []interface{}
	
	if status != "" {
		query = `
			SELECT id, job_id, username, email, skills, resume_url, ats_score, status, created_at
			FROM job_submissions
			WHERE job_id = $1 AND status = $2
			ORDER BY created_at DESC
		`
		args = []interface{}{jobID, status}
	} else {
		query = `
			SELECT id, job_id, username, email, skills, resume_url, ats_score, status, created_at
			FROM job_submissions
			WHERE job_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{jobID}
	}

	// Debug the query
	log.Printf("Executing query: %s with args: %v", query, args)

	var submissions []struct {
		ID        int            `json:"id" db:"id"`
		JobID     string         `json:"job_id" db:"job_id"`
		Username  string         `json:"username" db:"username"`
		Email     string         `json:"email" db:"email"`
		Skills    pq.StringArray `json:"skills" db:"skills"`
		ResumeURL string         `json:"resume_url" db:"resume_url"`
		ATSScore  int            `json:"ats_score" db:"ats_score"`
		Status    string         `json:"status" db:"status"`
		CreatedAt time.Time      `json:"created_at" db:"created_at"`
	}

	err = db.Select(&submissions, query, args...)
	if err != nil {
		log.Printf("Error retrieving job submissions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
		return
	}

	log.Printf("Found %d submissions for job ID: %s with status: %s", len(submissions), jobID, status)

	c.JSON(http.StatusOK, gin.H{
		"message": "Job submissions retrieved successfully",
		"data":    submissions,
	})
}

// UpdateSubmissionStatusH handles updating the status of a job submission
func UpdateSubmissionStatusH(c *gin.Context) {
	// Get submission ID from URL
	submissionID := c.Param("submission_id")
	if submissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "submission ID is required"})
		return
	}

	// Parse request body
	var request struct {
		Status string `json:"status" binding:"required,oneof=applied under_review shortlisted rejected selected"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update submission status in database
	db := database.GetDB()
	query := "UPDATE job_submissions SET status = $1 WHERE id = $2 RETURNING id, status"
	
	var updatedID int
	var updatedStatus string
	err := db.QueryRow(query, request.Status, submissionID).Scan(&updatedID, &updatedStatus)
	if err != nil {
		log.Printf("Error updating submission status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Submission status updated successfully",
		"data": gin.H{
			"id":     updatedID,
			"status": updatedStatus,
		},
	})
} 