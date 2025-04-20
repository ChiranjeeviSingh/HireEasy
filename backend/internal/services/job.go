package services

import (
    "backend/internal/database"
    "backend/internal/models"
    "context"
    "encoding/json"
    "errors"
    "github.com/lib/pq"
)

var (
    ErrJobExists = errors.New("job already exists for this user")
    ErrJobDoesNotExist = errors.New("job does not exist for this user")
)

func CreateJob(ctx context.Context, req *models.Job) error {

    db := database.GetDB()
    userID := ctx.Value("userID")

    // Check if job already exists for this user
    var count int
    err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM jobs WHERE job_id = $1 AND user_id = $2", req.JobID, userID)
    if err != nil {
        return err
    }
    if count > 0 {
        return ErrJobExists
    }

    // Insert new job

    // Convert map to JSON for attributes
    attributesJSON, err := json.Marshal(req.Attributes)
    if err != nil {
        return err
    }

    query := `INSERT INTO jobs (
        job_id, 
        user_id, 
        job_title, 
        job_description, 
        job_status, 
        skills_required, 
        attributes
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)`
    _, err = db.ExecContext(ctx, query, 
        req.JobID, 
        userID, 
        req.JobTitle, 
        req.JobDescription, 
        req.JobStatus, 
        pq.Array(req.SkillsRequired), 
        attributesJSON)

    return err
}

func UpdateJob(ctx context.Context, req *models.Job) error {
    db := database.GetDB()
    userID := ctx.Value("userID")

    // Check if job exists for this user
    var count int
    err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM jobs WHERE job_id = $1 AND user_id = $2", req.JobID, userID)
    if err != nil {
        return err
    }
    if count == 0 {
        return ErrJobDoesNotExist
    }

    // Convert map to JSON for attributes
    attributesJSON, err := json.Marshal(req.Attributes)
    if err != nil {
        return err
    }

    query := `UPDATE jobs SET 
        job_title = $1,
        job_description = $2,
        job_status = $3,
        skills_required = $4,
        attributes = $5
        WHERE job_id = $6 AND user_id = $7`

    _, err = db.ExecContext(ctx, query,
        req.JobTitle,
        req.JobDescription,
        req.JobStatus,
        pq.Array(req.SkillsRequired),
        attributesJSON,
        req.JobID,
        userID)

    return err
}

func GetJobById(ctx context.Context, jobID string) (*models.Job, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var job models.Job
    var skillsRequired pq.StringArray
    var attributesJSON []byte

    query := `SELECT job_id, job_title, job_description, job_status, skills_required, attributes 
             FROM jobs WHERE job_id = $1 AND user_id = $2`
    
    err := db.QueryRowContext(ctx, query, jobID, userID).Scan(
        &job.JobID,
        &job.JobTitle,
        &job.JobDescription,
        &job.JobStatus,
        &skillsRequired,
        &attributesJSON,
    )
    if err != nil {
        return nil, err
    }

    job.SkillsRequired = []string(skillsRequired)

    // Unmarshal attributes JSON
    if err := json.Unmarshal(attributesJSON, &job.Attributes); err != nil {
        return nil, err
    }

    return &job, nil
}

func GetJobsByTitle(ctx context.Context, jobTitle string) ([]*models.Job, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var jobs []*models.Job
    query := `SELECT job_id, job_title, job_description, job_status, skills_required, attributes 
             FROM jobs WHERE job_title ILIKE $1 AND user_id = $2`
    
    rows, err := db.QueryContext(ctx, query, "%"+jobTitle+"%", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var job models.Job
        var skillsRequired pq.StringArray
        var attributesJSON []byte

        err := rows.Scan(
            &job.JobID,
            &job.JobTitle,
            &job.JobDescription,
            &job.JobStatus,
            &skillsRequired,
            &attributesJSON,
        )
        if err != nil {
            return nil, err
        }

        job.SkillsRequired = []string(skillsRequired)
        if err := json.Unmarshal(attributesJSON, &job.Attributes); err != nil {
            return nil, err
        }

        jobs = append(jobs, &job)
    }

    return jobs, nil
}

func GetJobsByStatus(ctx context.Context, status string) ([]*models.Job, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var jobs []*models.Job
    query := `SELECT job_id, job_title, job_description, job_status, skills_required, attributes 
             FROM jobs WHERE job_status = $1 AND user_id = $2`
    
    rows, err := db.QueryContext(ctx, query, status, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var job models.Job
        var skillsRequired pq.StringArray
        var attributesJSON []byte

        err := rows.Scan(
            &job.JobID,
            &job.JobTitle,
            &job.JobDescription,
            &job.JobStatus,
            &skillsRequired,
            &attributesJSON,
        )
        if err != nil {
            return nil, err
        }

        job.SkillsRequired = []string(skillsRequired)
        if err := json.Unmarshal(attributesJSON, &job.Attributes); err != nil {
            return nil, err
        }

        jobs = append(jobs, &job)
    }

    return jobs, nil
}

func GetJobsByUserId(ctx context.Context) ([]*models.Job, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var jobs []*models.Job
    query := `SELECT job_id, job_title, job_description, job_status, skills_required, attributes 
             FROM jobs WHERE user_id = $1`
    
    rows, err := db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var job models.Job
        var skillsRequired pq.StringArray
        var attributesJSON []byte

        err := rows.Scan(
            &job.JobID,
            &job.JobTitle,
            &job.JobDescription,
            &job.JobStatus,
            &skillsRequired,
            &attributesJSON,
        )
        if err != nil {
            return nil, err
        }

        job.SkillsRequired = []string(skillsRequired)
        if err := json.Unmarshal(attributesJSON, &job.Attributes); err != nil {
            return nil, err
        }

        jobs = append(jobs, &job)
    }

    return jobs, nil
}

func DeleteJob(ctx context.Context, jobID string) error {
    db := database.GetDB()
    userID := ctx.Value("userID")

    // Check if job exists for this user
    var count int
    err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM jobs WHERE job_id = $1 AND user_id = $2", jobID, userID)
    if err != nil {
        return err
    }
    if count == 0 {
        return ErrJobDoesNotExist
    }

    // Delete the job
    query := `DELETE FROM jobs WHERE job_id = $1 AND user_id = $2`
    _, err = db.ExecContext(ctx, query, jobID, userID)
    return err
}