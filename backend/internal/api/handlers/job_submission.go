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
	
	// Construct query
	query := `
		SELECT id, job_id, username, email, skills, resume_url, ats_score, status, created_at
		FROM job_submissions
		WHERE job_id = $1
		ORDER BY created_at DESC
	`

	// Debug the query
	log.Printf("Executing query: %s with job_id=%s", query, jobID)

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

	err = db.Select(&submissions, query, jobID)
	if err != nil {
		log.Printf("Error retrieving job submissions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
		return
	}

	log.Printf("Found %d submissions for job ID: %s", len(submissions), jobID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Job submissions retrieved successfully",
		"data":    submissions,
	})
} 