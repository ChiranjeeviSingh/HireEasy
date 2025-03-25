package services

import (
    "context"
    "database/sql"
    "errors"
    "github.com/lib/pq"
    "time"
    "backend/internal/database"
    "backend/internal/models"
    "github.com/google/uuid"
	"encoding/json"
)

// Error variables for specific cases
var (
    ErrJobNotFound          = errors.New("job not found")
    ErrFormTemplateNotFound = errors.New("form template not found")
	ErrFormNotFound = errors.New("form not found")

)

// LinkJobToFormTemplate links a job to a form template and creates an entry in the application_form table.
func LinkJobToFormTemplate(ctx context.Context, jobID string, formTemplateID string) (*models.ApplicationForm, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var dbJobID int
    var dbFormID int

    // Check if job exists
    err := db.QueryRowContext(ctx, "SELECT id FROM jobs WHERE job_id = $1 and user_id= $2", jobID, userID).Scan(&dbJobID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrJobNotFound
        }
        return nil, err
    }

    // Check if form template exists
    err = db.QueryRowContext(ctx, "SELECT id FROM form_templates WHERE form_template_id = $1 and user_id= $2", formTemplateID, userID).Scan(&dbFormID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrFormTemplateNotFound
        }
        return nil, err
    }

    // Insert into application_form table
    query := `
        INSERT INTO application_form (form_uuid, job_id, form_id, date_created)
        VALUES ($1, $2, $3, $4)
        RETURNING form_uuid, job_id, form_id, date_created
    `

    var applicationForm models.ApplicationForm
    err = db.QueryRowContext(ctx, query, uuid.New().String(), dbJobID, dbFormID, time.Now()).
        Scan(&applicationForm.FormUUID, &applicationForm.JobID, &applicationForm.FormID, &applicationForm.DateCreated)
    if err != nil {
        return nil, err
    }

    return &applicationForm, nil
}

func UpdateFormStatus(ctx context.Context, formUUID string, status string) (*models.ApplicationForm, error) {
    db := database.GetDB()

    // Check if form exists
    var applicationForm models.ApplicationForm
    err := db.QueryRowContext(ctx, "SELECT form_uuid, job_id, form_id, status, date_created FROM application_form WHERE form_uuid = $1", formUUID).
        Scan(&applicationForm.FormUUID, &applicationForm.JobID, &applicationForm.FormID, &applicationForm.Status, &applicationForm.DateCreated)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrFormNotFound
        }
        return nil, err
    }

    // Update form status
    query := "UPDATE application_form SET status = $1 WHERE form_uuid = $2 RETURNING status"
    err = db.QueryRowContext(ctx, query, status, formUUID).Scan(&applicationForm.Status)
    if err != nil {
        return nil, err
    }

    return &applicationForm, nil
}

// GetFormDetails retrieves job and form template details associated with the given form_uuid.
func GetFormDetails(ctx context.Context, formUUID string) (*models.GetFormResponse, error) {
    db := database.GetDB()

    // Query to fetch form details
    var form models.ApplicationForm
    err := db.QueryRowContext(ctx, `SELECT form_uuid, job_id, form_id, status, date_created 
                                    FROM application_form WHERE form_uuid = $1`, formUUID).
        Scan(&form.FormUUID, &form.JobID, &form.FormID, &form.Status, &form.DateCreated)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrFormNotFound
        }
        return nil, err
    }

    // Query to fetch job details
    var job models.JobDetails
	var skillsRequired pq.StringArray
    var attributesJSON []byte
    err = db.QueryRowContext(ctx, `SELECT job_id, job_title, job_description, job_status, skills_required, created_at, updated_at, attributes 
                                   FROM jobs WHERE id = $1`, form.JobID).
        Scan(&job.JobID, &job.JobTitle, &job.JobDescription, &job.JobStatus, &skillsRequired, &job.CreatedAt, &job.UpdatedAt, &attributesJSON)
    if err != nil {
        return nil, err
    }

	job.SkillsRequired = []string(skillsRequired)
	// Unmarshal attributes JSON
    if err := json.Unmarshal(attributesJSON, &job.Attributes); err != nil {
        return nil, err
    }

    // Query to fetch form template details
    var formTemplate models.FormTemplateDetails
	var fieldsJSON []byte
    err = db.QueryRowContext(ctx, `SELECT form_template_id, fields, created_at, updated_at FROM form_templates WHERE id = $1`, form.FormID).
        Scan(&formTemplate.FormTemplateID, &fieldsJSON, &formTemplate.CreatedAt, &formTemplate.UpdatedAt)
    if err != nil {
        return nil, err
    }

	// Unmarshal attributes JSON
	if err := json.Unmarshal(fieldsJSON, &formTemplate.Fields); err != nil {
        return nil, err
    }

    return &models.GetFormResponse{
        FormUUID:     form.FormUUID,
        Status:       form.Status,
        DateCreated:  form.DateCreated,
        JobDetails:   job,
        FormTemplate: formTemplate,
    }, nil
}

// DeleteForm deletes a form by its UUID, unlinking it from the job and form template.
func DeleteForm(ctx context.Context, formUUID string) error {
    db := database.GetDB()

    // Check if form exists before deleting
    var exists bool
    err := db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM application_form WHERE form_uuid = $1)", formUUID).
        Scan(&exists)
    if err != nil {
        return err
    }
    if !exists {
        return ErrFormNotFound
    }

    // Delete the form
    _, err = db.ExecContext(ctx, "DELETE FROM application_form WHERE form_uuid = $1", formUUID)
    if err != nil {
        return err
    }

    return nil
}