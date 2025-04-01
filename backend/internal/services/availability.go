package services

import (
    "backend/internal/database"
    "backend/internal/models"
    "context"
    "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"database/sql"
	"strings"
	"github.com/lib/pq"
	"log"
)

var (
	ErrMaxDurationExceeded = errors.New("time slot cannot be more than one hour")
	ErrInvalidTime = errors.New("invalid start/end time")
	ErrOverlappingSlot = errors.New("overlapping time slot exists")
	ErrSlotScheduled = errors.New("cannot delete: slot is scheduled for interview. Conact Recruiter")
	ErrNoAvailability = errors.New("no availability found with this id")
	ErrUsernameNotFound = errors.New("no user exists with given username")
)


func CreateAvailability(ctx context.Context, req *models.Availability) error {
	db := database.GetDB()
	userID := ctx.Value("userID")

	// Parse time strings into time.Time
	fromTime, err := time.Parse("15:04", req.FromTime)
	if err != nil {
		return ErrInvalidTime
	}
	toTime, err := time.Parse("15:04", req.ToTime)
	if err != nil {
		return ErrInvalidTime
	}

	// Validate time difference is not more than one hour
	if toTime.Sub(fromTime).Hours() > 1 {
		return ErrMaxDurationExceeded
	}

	// Validate hours, minutes and seconds
	if fromTime.Hour() < 0 || fromTime.Hour() > 23 ||
	   toTime.Hour() < 0 || toTime.Hour() > 23 ||
	   fromTime.Minute() < 0 || fromTime.Minute() > 59 ||
	   toTime.Minute() < 0 || toTime.Minute() > 59 {
		return ErrInvalidTime
	}

	// Check for overlapping slots
	var count int

	err = db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM availabilities 
		WHERE user_id = $1 
		AND date = $4
		AND ((from_time <= $2 AND to_time > $2)
		OR (from_time < $3 AND to_time >= $3)
		OR (from_time >= $2 AND to_time <= $3))`,
		userID, req.FromTime, req.ToTime, req.Date)

	if err != nil {
		return err
	}
	if count > 0 {
		return ErrOverlappingSlot
	}

	// Insert availability
	var id int
	err = db.QueryRowContext(ctx, `
		INSERT INTO availabilities (user_id, from_time, to_time, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		userID, req.FromTime, req.ToTime, req.Date).Scan(&id)
	if err != nil {
		return err
	}
	req.ID = id
	return err
}

func DeleteAvailability(ctx context.Context, availabilityID int) error {
	db := database.GetDB()
	userID := ctx.Value("userID")

	// Check if availability exists
	var existingAvailabilityID int
	err := db.GetContext(ctx, &existingAvailabilityID, `
		SELECT id FROM availabilities 
		WHERE id = $1 AND user_id = $2`,
		availabilityID, userID)

	if err != nil || existingAvailabilityID == 0 {
		return ErrNoAvailability
	}

	// Check if slot is used in interviews
	var count int
	err = db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM interviews 
		WHERE availability_id = $1`,
		availabilityID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSlotScheduled
	}

	// Delete availability
	_, err = db.ExecContext(ctx, `
		DELETE FROM availabilities 
		WHERE id = $1 AND user_id = $2`,
		availabilityID, userID)
	return err
}

func GetMyAvailability(ctx *gin.Context) ([]*models.Availability, error) {
	userID := ctx.Value("userID").(int)
	return getUserAvailabilityHelper(ctx, userID)
}

func GetUserAvailability(ctx *gin.Context, userName string) ([]*models.Availability, error) {
	db := database.GetDB()
	
	var userID int
	err := db.GetContext(ctx.Request.Context(), &userID, `
		SELECT id FROM users WHERE userName = $1`, userName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUsernameNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return getUserAvailabilityHelper(ctx, userID)
}

func getUserAvailabilityHelper(ctx *gin.Context, userID int) ([]*models.Availability, error) {
	db := database.GetDB()
	availabilities := []*models.Availability{}
	
	fromDate := ctx.Query("from_date")
    toDate := ctx.Query("to_date")
	
	fmt.Println("fromDate:", fromDate)
    fmt.Println("toDate:", toDate)
	
	query := `SELECT id, user_id, date, from_time, to_time, created_at, updated_at
		FROM availabilities 
		WHERE user_id = $1`
	
	args := []interface{}{userID}
	argCount := 2

	if fromDate != "" {
		query += fmt.Sprintf(` AND date >= $%d`, argCount)
		args = append(args, fromDate)
		argCount++
	}

	if toDate != "" {
		query += fmt.Sprintf(` AND date <= $%d`, argCount)
		args = append(args, toDate)
		argCount++
	}

	query += ` ORDER BY from_time`

	rows, err := db.QueryContext(ctx.Request.Context(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		availability := &models.Availability{}
		err := rows.Scan(
			&availability.ID,
			&availability.UserID,
			&availability.Date,
			&availability.FromTime,
			&availability.ToTime,
			&availability.CreatedAt,
			&availability.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}


func GetAllAvailability(ctx *gin.Context) ([]*models.Availability, error) {
	db := database.GetDB()
	availabilities := []*models.Availability{}
	
	fromDate := ctx.Query("from_date")
    toDate := ctx.Query("to_date")
	yearsExperience := ctx.Query("years_experience")
	jobTitle:= ctx.Query("job_title")

	var areasOfExpertise []string
	if expertise := ctx.Query("areas_of_expertise"); expertise != "" {
    	areasOfExpertise = strings.Split(strings.ToLower(expertise), ",")
	}

	query := `SELECT DISTINCT a.id, a.user_id, a.from_time, a.to_time, a.date, a.created_at, a.updated_at
		FROM availabilities a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE 1=1`
	
	args := []interface{}{}
	argCount := 1

	if fromDate != "" {
		query += fmt.Sprintf(` AND a.date >= $%d`, argCount)
		args = append(args, fromDate)
		argCount++
	}

	if toDate != "" {
		query += fmt.Sprintf(` AND a.date <= $%d`, argCount)
		args = append(args, toDate)
		argCount++
	}

	if yearsExperience != "" {
		query += fmt.Sprintf(` AND p.years_of_experience = $%d`, argCount)
		args = append(args, yearsExperience)
		argCount++
	}

	if jobTitle != "" {
		query += fmt.Sprintf(` AND p.job_title ILIKE $%d`, argCount)
		args = append(args, "%"+jobTitle+"%")
		argCount++
	}

	if len(areasOfExpertise) > 0 {
		log.Println("here again", areasOfExpertise)
		query += fmt.Sprintf(` AND p.areas_of_expertise && $%d`, argCount)
		args = append(args, pq.Array(areasOfExpertise))
		argCount++
	}
	

	query += ` ORDER BY a.from_time`

	rows, err := db.QueryContext(ctx.Request.Context(), query, args...)
	if err != nil {
		log.Println("here error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		availability := &models.Availability{}
		err := rows.Scan(
			&availability.ID,
			&availability.UserID,
			&availability.FromTime,
			&availability.ToTime,
			&availability.Date,
			&availability.CreatedAt,
			&availability.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}