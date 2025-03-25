package models

import ( 
    "time"
    )

type FormTemplate struct {
    ID             int       `json:"id" db:"id"`
    FormTemplateID string    `json:"form_template_id" binding:"required" db:"form_template_id"`
    UserID         int       `json:"user_id" db:"user_id"`
    Fields         [] map[string]interface{}  `json:"fields" db:"fields"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}