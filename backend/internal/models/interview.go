package models

import ( 
    "time"
    )

type Interview struct {
	ID               int       `json:"id" db:"id" binding:"required"`
	JobID            string       `json:"job_id,omitempty" db:"job_id" binding:"required"`
	HRUserID         int       `json:"hr_user_id,omitempty" db:"hr_user_id" binding:"required"`
	JobSubmissionID  int       `json:"job_submission_id,omitempty" db:"job_submission_id" binding:"required"`
	InterviewerID    int       `json:"interviewer_user_id,omitempty" db:"interviewer_user_id" binding:"required"`
	AvailabilityID   int       `json:"availability_id,omitempty" db:"availability_id" binding:"required"`
	Feedback         string   `json:"feedback,omitempty" db:"feedback"`
	Verdict          string    `json:"verdict,omitempty" db:"verdict"`
	Status           string    `json:"status,omitempty" db:"status" binding:"required"`
	CreatedAt        time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type CreateInterviewRequest struct {
	JobID           string `json:"job_id" binding:"required"`
	JobSubmissionID int `json:"job_submission_id" binding:"required"`
	InterviewerID   int `json:"interviewer_user_id" binding:"required"`
	AvailabilityID  int `json:"availability_id" binding:"required"`
}
type FeedbackRequest struct {
	Verdict  string `json:"verdict" binding:"required"`
	Feedback string `json:"feedback" binding:"required"`
}
