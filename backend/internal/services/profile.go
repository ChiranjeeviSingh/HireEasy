package services

import (
    "context"
    "backend/internal/database"
    "backend/internal/models"
	"github.com/lib/pq"
	"errors"
	"strings"
	"database/sql"
	"fmt"
)


var (
    ErrProfileExists = errors.New("profile already exists for this user")
)

func CreateProfile(ctx context.Context, req *models.Profile) error {
	db := database.GetDB()
	
	userID := ctx.Value("userID")

	// First check if a profile already exists for this user
	var exists bool
	err := db.QueryRowContext(ctx, 
		"SELECT EXISTS(SELECT 1 FROM profiles WHERE user_id = $1)", 
		userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrProfileExists
	}

	// Convert all strings in Areas_of_expertise to lowercase
	lowerCaseAOE := make([]string, len(req.Areas_of_expertise))
	for i, area := range req.Areas_of_expertise {
		lowerCaseAOE[i] = strings.ToLower(area)
	}

	// If no profile exists, create new one
	_, err = db.ExecContext(ctx,
		`INSERT INTO profiles (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
		 VALUES ($1, $2, $3, $4, $5)`,
		userID, req.JobTitle, req.YearsOfExperience, pq.Array(lowerCaseAOE), req.PhoneNumber)
	
	return err
}

func GetMyProfile(ctx context.Context) (*models.Profile, error) {
	db := database.GetDB()
	userID := ctx.Value("userID")

	var profile models.Profile
	var areasOfExpertise pq.StringArray
	
	err := db.QueryRowContext(ctx,
		`SELECT user_id, job_title, years_of_experience, areas_of_expertise, phone_number 
		 FROM profiles WHERE user_id = $1`, userID).Scan(
		&profile.UserID,
		&profile.JobTitle,
		&profile.YearsOfExperience,
		&areasOfExpertise,
		&profile.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	profile.Areas_of_expertise = []string(areasOfExpertise)

	return &profile, nil
}

func GetUserProfile(ctx context.Context, userName string) (*models.Profile, error) {
	
	db := database.GetDB()

	var userID int
	err := db.GetContext(ctx, &userID, `
		SELECT id FROM users WHERE userName = $1`, userName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUsernameNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	var profile models.Profile
	var areasOfExpertise pq.StringArray
	
	err = db.QueryRowContext(ctx,
		`SELECT user_id, job_title, years_of_experience, areas_of_expertise, phone_number 
		 FROM profiles WHERE user_id = $1`, userID).Scan(
		&profile.UserID,
		&profile.JobTitle,
		&profile.YearsOfExperience,
		&areasOfExpertise,
		&profile.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	profile.Areas_of_expertise = []string(areasOfExpertise)

	return &profile, nil
}


func UpdateMyProfile(ctx context.Context, req *models.Profile) error {
	db := database.GetDB()

	userID := ctx.Value("userID")
	
	// Convert all strings in Areas_of_expertise to lowercase
	lowerCaseAOE := make([]string, len(req.Areas_of_expertise))
	for i, area := range req.Areas_of_expertise {
		lowerCaseAOE[i] = strings.ToLower(area)
	}

	_, err := db.ExecContext(ctx,
		`UPDATE profiles 
		 SET job_title = $2, years_of_experience = $3, areas_of_expertise = $4, phone_number = $5
		 WHERE user_id = $1`,
		userID, req.JobTitle, req.YearsOfExperience, pq.Array(lowerCaseAOE), req.PhoneNumber)
	
	return err
}