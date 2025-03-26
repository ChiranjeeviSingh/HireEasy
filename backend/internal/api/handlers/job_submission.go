package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/services"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// HandleFormSubmission processes a form submission
func HandleFormSubmission(c *gin.Context) {
	db := database.GetDB()
	service := services.NewFormSubmissionService(db)

	submission, err := service.HandleFormSubmission(c)
	if err != nil {
		log.Printf("Error handling form submission: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        submission.ID,
		"ats_score": submission.ATSScore,
		"message":   "Application submitted successfully",
	})
}

// GetFormSubmissions retrieves all submissions for a specific form
func GetFormSubmissions(c *gin.Context) {
	// Get optional sorting parameter (default to ats_score)
	sortBy := c.DefaultQuery("sort_by", "ats_score")
	// Validate sort parameter
	validSortFields := map[string]bool{
		"ats_score":  true,
		"created_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "ats_score" // Default to ats_score if invalid
	}

	// Get optional limit parameter (default to 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default to 10 if invalid
	}

	// Get date filter (all, today, or specific date)
	dateFilter := c.DefaultQuery("date", "all")
	
	// Build the query based on parameters
	baseQuery := `
		SELECT id, job_id, username, email, form_data, skills, resume_url, ats_score, status, created_at, updated_at
		FROM job_submissions 
		WHERE job_id = $1`
	
	var query string
	var queryParams []interface{}
	queryParams = append(queryParams, c.Param("job_id"))
	
	if dateFilter == "today" {
		query = baseQuery + " AND created_at::date = CURRENT_DATE"
	} else if dateFilter != "all" {
		// Try to parse the date
		_, err := time.Parse("2006-01-02", dateFilter)
		if err == nil {
			query = baseQuery + " AND created_at::date = $2"
			queryParams = append(queryParams, dateFilter)
		} else {
			log.Println("Invalid date format, ignoring date filter:", dateFilter)
			query = baseQuery
		}
	} else {
		query = baseQuery
	}
	
	// Add ordering and limit
	query += fmt.Sprintf(" ORDER BY %s DESC LIMIT $%d", sortBy, len(queryParams)+1)
	queryParams = append(queryParams, limit)
	
	// Log the query for debugging
	log.Println("Executing query:", query, "with params:", queryParams)
	
	// Execute the query
	db := database.GetDB()
	rows, err := db.Query(query, queryParams...)
	if err != nil {
		log.Println("Database query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
		return
	}
	defer rows.Close()
	
	// Process the results
	var submissions []models.JobSubmission
	for rows.Next() {
		var sub models.JobSubmission
		var skills pq.StringArray
		
		if err := rows.Scan(
			&sub.ID,
			&sub.JobID,
			&sub.Username,
			&sub.Email,
			&sub.FormData,
			&skills,
			&sub.ResumeURL,
			&sub.ATSScore,
			&sub.Status,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		
		sub.Skills = []string(skills)
		submissions = append(submissions, sub)
	}
	
	if err = rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing submissions"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"submissions": submissions,
	})
}
