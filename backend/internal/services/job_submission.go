package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type FormSubmissionService struct {
	db *sqlx.DB
}

func NewFormSubmissionService(db *sqlx.DB) *FormSubmissionService {
	return &FormSubmissionService{db: db}
}

func (s *FormSubmissionService) HandleFormSubmission(c *gin.Context) (*models.JobSubmission, error) {
	// Check table structure first (for debugging)
	if err := s.checkJobSubmissionsTable(); err != nil {
		log.Printf("Warning: Failed to check job_submissions table structure: %v", err)
	}

	var submission models.FormSubmissionRequest
	
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		return nil, fmt.Errorf("failed to parse form: %v", err)
	}

	// Bind form data
	if err := c.ShouldBind(&submission); err != nil {
		log.Printf("Failed to bind form data: %v", err)
		return nil, fmt.Errorf("invalid form data: %v", err)
	}

	// Parse form data JSON
	var formDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(submission.FormData), &formDataMap); err != nil {
		log.Printf("Failed to parse form data JSON: %v", err)
		return nil, fmt.Errorf("invalid form data format: %v", err)
	}

	// Extract skills from form data
	var skills []string
	if skillsInterface, ok := formDataMap["skills"].([]interface{}); ok {
		for _, skill := range skillsInterface {
			if skillStr, ok := skill.(string); ok {
				skills = append(skills, skillStr)
			}
		}
	}

	// Validate that skills are not empty
	if len(skills) == 0 {
		return nil, fmt.Errorf("skills are required")
	}

	// Check if we're in test mode
	testMode := c.GetHeader("X-Test-Mode") == "true" || config.GetConfig().TestMode

	// Upload resume to S3
	var resumeURL string
	var err error

	if !testMode {
		if submission.Resume == nil {
			return nil, fmt.Errorf("resume file is required")
		}
		resumeURL, err = UploadResumeToS3(submission.Resume, 0, submission.Username)
		if err != nil {
			log.Printf("Failed to upload resume: %v", err)
			return nil, fmt.Errorf("failed to upload resume: %v", err)
		}
	} else {
		// In test mode, use a mock URL
		resumeURL = fmt.Sprintf("https://test-bucket.s3.amazonaws.com/resumes/%s_test_user.pdf", submission.Username)
		log.Println("Test mode: Using mock resume URL:", resumeURL)
	}

	// Calculate ATS score
	atsScore := calculateATSScore(skills)

	// Create job submission
	jobSubmission := &models.JobSubmission{
		JobID:     submission.JobID,
		Username:  submission.Username,
		Email:     submission.Email,
		FormData:  []byte(submission.FormData),
		Skills:    pq.StringArray(skills),
		ResumeURL: resumeURL,
		ATSScore:  int(atsScore),
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into database
	err = s.insertJobSubmission(jobSubmission)
	if err != nil {
		log.Printf("Failed to insert job submission: %v", err)
		return nil, fmt.Errorf("failed to save submission: %v", err)
	}

	return jobSubmission, nil
}

func (s *FormSubmissionService) insertJobSubmission(submission *models.JobSubmission) error {
	// First let's try to understand the structure of the job_submissions table
	tableInfo, err := s.db.Query(`
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = 'job_submissions'
	`)
	if err != nil {
		log.Printf("Error examining job_submissions table: %v", err)
		// Continue anyway, we'll try the insert
	} else {
		defer tableInfo.Close()
		
		var columns []string
		for tableInfo.Next() {
			var name, dataType, nullable string
			if err := tableInfo.Scan(&name, &dataType, &nullable); err == nil {
				columns = append(columns, fmt.Sprintf("%s(%s, %s)", name, dataType, nullable))
			}
		}
		log.Printf("Table structure: %v", columns)
	}

	// Check if the form_uuid column exists
	var hasFormUUID bool
	err = s.db.Get(&hasFormUUID, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'job_submissions' AND column_name = 'form_uuid'
		)
	`)
	if err != nil {
		log.Printf("Error checking for form_uuid column: %v", err)
		// Continue anyway, we'll try the insert
	}

	var insertQuery string
	var args []interface{}

	if hasFormUUID {
		// Try to get any form_uuid for this job
		var formUUID string
		err = s.db.QueryRow(`
			SELECT id FROM form_templates 
			WHERE job_id = $1 
			LIMIT 1
		`, submission.JobID).Scan(&formUUID)
		
		if err != nil {
			// If no form template exists, use a default UUID for test submissions
			formUUID = "00000000-0000-0000-0000-000000000000"
			log.Printf("Using default UUID as no form template found: %s", formUUID)
		}
		
		// Query with form_uuid
		insertQuery = `
			INSERT INTO job_submissions (
				form_uuid, job_id, username, email, form_data, skills, resume_url, ats_score, status, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id`
		
		args = []interface{}{
			formUUID,
			submission.JobID,
			submission.Username,
			submission.Email,
			submission.FormData,
			submission.Skills,
			submission.ResumeURL,
			submission.ATSScore,
			submission.Status,
			submission.CreatedAt,
			submission.UpdatedAt,
		}
	} else {
		// Query without form_uuid
		insertQuery = `
			INSERT INTO job_submissions (
				job_id, username, email, form_data, skills, resume_url, ats_score, status, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id`
		
		args = []interface{}{
			submission.JobID,
			submission.Username,
			submission.Email,
			submission.FormData,
			submission.Skills,
			submission.ResumeURL,
			submission.ATSScore,
			submission.Status,
			submission.CreatedAt,
			submission.UpdatedAt,
		}
	}

	log.Printf("Using query: %s", insertQuery)
	err = s.db.QueryRow(insertQuery, args...).Scan(&submission.ID)

	if err != nil {
		return fmt.Errorf("failed to insert job submission: %v", err)
	}

	return nil
}

func calculateATSScore(skills []string) int {
	// Simple scoring logic based on number of skills
	baseScore := 70
	skillPoints := len(skills) * 5
	
	score := baseScore + skillPoints
	if score > 100 {
		score = 100
	}
	return score
}

func (s *FormSubmissionService) checkJobSubmissionsTable() error {
	// Check if the job_submissions table exists and get its columns
	query := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'job_submissions'
		ORDER BY ordinal_position;
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query table schema: %v", err)
	}
	defer rows.Close()
	
	var columns []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return fmt.Errorf("failed to scan column name: %v", err)
		}
		columns = append(columns, col)
	}
	
	log.Printf("Job submissions table columns: %v", columns)
	return nil
}
