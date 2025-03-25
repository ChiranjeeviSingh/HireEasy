package models

import "time"

type UpdateFormStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

type LinkJobToFormRequest struct {
    FormTemplateID string `json:"form_template_id" binding:"required"`
}

type ApplicationForm struct {
    FormUUID    string    `json:"form_uuid" db:"form_uuid"`
    JobID       int       `json:"job_id" db:"job_id"`       
    FormID      int       `json:"form_id" db:"form_id"`      
	Status      string    `json:"job_status" db:"status"`
    DateCreated time.Time `json:"date_created" db:"date_created"`
}

type GetFormResponse struct {
    FormUUID       string              `json:"form_uuid"`
    Status         string              `json:"status"`
    DateCreated    time.Time           `json:"date_created"`
    JobDetails     JobDetails          `json:"job"`
    FormTemplate   FormTemplateDetails `json:"form_template"`
}


type JobDetails struct {
    JobID          string                 `json:"job_id"`
    JobTitle       string                 `json:"job_title"`
    JobDescription string                 `json:"job_description"`
    JobStatus      string                 `json:"job_status"`
    SkillsRequired []string               `json:"skills_required"`
    CreatedAt      time.Time              `json:"created_at"`
    UpdatedAt      time.Time              `json:"updated_at"`
    Attributes     map[string]interface{} `json:"attributes"`
}

type FormTemplateDetails struct {
    FormTemplateID string                   `json:"form_template_id"`
    Fields         []map[string]interface{} `json:"fields"`
    CreatedAt      time.Time                `json:"created_at"`
    UpdatedAt      time.Time                `json:"updated_at"`
}

