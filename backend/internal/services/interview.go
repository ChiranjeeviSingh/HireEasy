package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"errors"
	"fmt"
	"time"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var (
	ErrInterviewNotFound    = errors.New("invalid interview id")
	ErrUnauthorizedFeedback = errors.New("unauthorized: only interviewers can submit feedback")
)

func CreateInterview(ctx context.Context, req *models.CreateInterviewRequest) (*models.Interview, error) {
	db := database.GetDB()
	hrUserID := ctx.Value("userID").(int)

	// Check if availability is already used in an interview
	var count int
	err := db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM interviews 
		WHERE availability_id = $1`,
		req.AvailabilityID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrOverlappingSlot
	}

	// Insert new interview
	_, err = db.ExecContext(ctx, `
		INSERT INTO interviews (
			job_id, 
			hr_user_id, 
			job_submission_id, 
			interviewer_user_id, 
			availability_id,
			status
		) VALUES ($1, $2, $3, $4, $5, 'scheduled')`,
		req.JobID,
		hrUserID,
		req.JobSubmissionID,
		req.InterviewerID,
		req.AvailabilityID)

	interview := &models.Interview{
		JobID:           req.JobID,
		HRUserID:        hrUserID,
		JobSubmissionID: req.JobSubmissionID,
		InterviewerID:   req.InterviewerID,
		AvailabilityID:  req.AvailabilityID,
		Status:          "scheduled",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	return interview, err
}

func DeleteInterview(ctx context.Context, interviewID int) error {
	db := database.GetDB()
	userID := ctx.Value("userID").(int)

	// Check if interview exists for this HR user
	var count int
	err := db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM interviews 
		WHERE id = $1 AND hr_user_id = $2`,
		interviewID, userID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInterviewNotFound
	}

	// Delete interview
	_, err = db.ExecContext(ctx, `
		DELETE FROM interviews 
		WHERE id = $1 AND hr_user_id = $2`,
		interviewID, userID)
	return err
}

func ListAllInterviews(ctx *gin.Context) ([]*models.Interview, error) {
	db := database.GetDB()
	userID := ctx.Value("userID").(int)

	// Get user role
	var userRole string
	err := db.GetContext(ctx.Request.Context(), &userRole, `
		SELECT role FROM users WHERE id = $1`, userID)
	if err != nil {
		return nil, err
	}

	// Handle date parameters
	fromDate := ctx.Query("from_date")
	toDate := ctx.Query("to_date")
	if fromDate == "" || toDate == "" {
		now := time.Now()
		firstDayofCurrMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		lastDayofCurrMonth := firstDayofCurrMonth.AddDate(0, 1, -1)

		fromDate = firstDayofCurrMonth.Format("2006-01-02")
		toDate = lastDayofCurrMonth.Format("2006-01-02")
	}

	// base query
	query := `
		SELECT i.id, i.job_id, i.hr_user_id, i.job_submission_id, 
			   i.interviewer_user_id, i.availability_id, i.status, 
			   i.created_at, i.updated_at, a.date, a.to_time,
			   i.feedback, i.verdict
		FROM interviews i
		JOIN availabilities a ON i.availability_id = a.id
		WHERE a.date BETWEEN $1 AND $2`
		
	args := []interface{}{fromDate, toDate}
	argCount := 3

	// Add role-specific filter
	if userRole == "HR" {
		query += fmt.Sprintf(" AND i.hr_user_id = $%d", argCount)
		args = append(args, userID)
		argCount++
	} else {
		query += fmt.Sprintf(" AND i.interviewer_user_id = $%d", argCount)
		args = append(args, userID)
		argCount++
	}

	// Add job_submission_id filter if provided
	if jobSubmissionID := ctx.Query("job_submission_id"); jobSubmissionID != "" {
		query += fmt.Sprintf(" AND i.job_submission_id = $%d", argCount)
		args = append(args, jobSubmissionID)
		argCount++
	}

	// Add status filter if provided
	if status := ctx.Query("status"); status != "" {
		query += fmt.Sprintf(" AND i.status = $%d", argCount)
		args = append(args, status)
	}

	rows, err := db.QueryContext(ctx.Request.Context(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currentTime := time.Now()
	interviews := []*models.Interview{}

	for rows.Next() {
		interview := &models.Interview{}
		var interviewDate string
		var toTime string

		var feedback sql.NullString
		err := rows.Scan(
			&interview.ID,
			&interview.JobID,
			&interview.HRUserID,
			&interview.JobSubmissionID,
			&interview.InterviewerID,
			&interview.AvailabilityID,
			&interview.Status,
			&interview.CreatedAt,
			&interview.UpdatedAt,
			&interviewDate,
			&toTime,
			&feedback,
			&interview.Verdict,
		)
		if feedback.Valid {
			interview.Feedback = feedback.String
		}

		if err != nil {
			return nil, err
		}

		// Check if scheduled interviews need status update
		if interview.Status == "scheduled" {

			parsedDate, _ := time.Parse(time.RFC3339, interviewDate)
			parsedTime, _ := time.Parse(time.RFC3339, toTime)
			
			interviewDateTime := time.Date(
				parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
				0, time.Local,
			)

			if err != nil {
				return nil, err
			}

			if currentTime.After(interviewDateTime) {
				interview.Status = "pending_feedback"
				_, err = db.ExecContext(ctx.Request.Context(), `
					UPDATE interviews SET status = 'pending_feedback' 
					WHERE id = $1`, interview.ID)
				if err != nil {
					return nil, err
				}
			}
		}

		interviews = append(interviews, interview)
	}

	return interviews, nil
}

func SubmitFeedback(ctx context.Context, req *models.FeedbackRequest, interviewID int) error {

	db := database.GetDB()
	userID := ctx.Value("userID").(int)

	// Check if user is an interviewer
	var role string
	err := db.GetContext(ctx, &role, `SELECT role FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	if role != "INTERVIEWER" {
		return ErrUnauthorizedFeedback
	}

	// Check if interview exists and user is the assigned interviewer
	var count int
	err = db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM interviews 
		WHERE id = $1 AND interviewer_user_id = $2`,
		interviewID, userID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInterviewNotFound
	}

	// Update interview with feedback and status
	_, err = db.ExecContext(ctx, `
		UPDATE interviews 
		SET verdict = $1, feedback = $2, status = 'completed', updated_at = NOW()
		WHERE id = $3`,
		req.Verdict, req.Feedback, interviewID)

	return err
}
