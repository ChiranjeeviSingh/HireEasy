package models

import ( 
    "time"
    )
	
type Profile struct {
	ID                 int       `json:"id,omitempty" db:"id"`
	UserID             int       `json:"user_id,omitempty" db:"user_id"`
	JobTitle           string    `json:"job_title,omitempty" binding:"required" db:"job_title"`
	YearsOfExperience  int       `json:"years_of_experience,omitempty" binding:"required" db:"years_of_experience"`
	Areas_of_expertise []string  `json:"areas_of_expertise,omitempty" binding:"required,min=1" db:"areas_of_expertise"`
	PhoneNumber        string    `json:"phone_number,omitempty" db:"phone_number"`
	CreatedAt          time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at,omitempty" db:"updated_at"`
}