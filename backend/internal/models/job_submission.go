package models

import (
	"mime/multipart"
	"time"

	"github.com/lib/pq"
)

// FormSubmissionRequest represents a job submission request from the frontend
// Note: Only HR has accounts in the portal, regular applicants don't have user IDs
type FormSubmissionRequest struct {
	JobID    string                `form:"job_id" binding:"required"`
	Username string                `form:"username" binding:"required"`
	Email    string                `form:"email" binding:"required,email"`
	FormData string                `form:"form_data" binding:"required"`
	Resume   *multipart.FileHeader `form:"resume" binding:"required"`
}

// JobSubmission represents a job application in the database
type JobSubmission struct {
	ID        int            `json:"id" db:"id"`
	JobID     string         `json:"job_id" db:"job_id"`
	Username  string         `json:"username" db:"username"`
	Email     string         `json:"email" db:"email"`
	FormData  []byte         `json:"-" db:"form_data"` // Store raw JSON, process later
	Skills    pq.StringArray `json:"skills" db:"skills"`
	ResumeURL string         `json:"resume_url" db:"resume_url"`
	ATSScore  int            `json:"ats_score" db:"ats_score"`
	Status    string         `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}
