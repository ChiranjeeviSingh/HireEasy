package models

type Job struct {
	ID             int                    `json:"id,omitempty" db:"id"`
	JobID          string                 `json:"job_id,omitempty" binding:"required" db:"job_id"`
	UserID         int                    `json:"user_id,omitempty" db:"user_id"`
	JobTitle       string                 `json:"job_title,omitempty" binding:"required" db:"job_title"`
	JobDescription string                 `json:"job_description,omitempty" binding:"required" db:"job_description"`
	JobStatus      string                 `json:"job_status,omitempty" binding:"required" db:"job_status"`
	SkillsRequired []string               `json:"skills_required,omitempty" binding:"required" db:"skills_required"`
	CreatedAt      string                 `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      string                 `json:"updated_at,omitempty" db:"updated_at"`
	Attributes     map[string]interface{} `json:"attributes,omitempty" db:"attributes"`
}
