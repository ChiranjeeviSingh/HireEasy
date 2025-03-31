package services

import (
    "context"
    "backend/internal/database"
    "backend/internal/models"
	"github.com/lib/pq"
)


func CreateProfile(ctx context.Context, req *models.Profile) error {
	db := database.GetDB()
	
	userID := ctx.Value("userID")

	_, err := db.ExecContext(ctx,
		`INSERT INTO profiles (user_id, job_title, years_of_experience, areas_of_expertise, phone_number)
		 VALUES ($1, $2, $3, $4, $5)`,
		userID, req.JobTitle, req.YearsOfExperience, pq.Array(req.Areas_of_expertise), req.PhoneNumber)
	
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

func UpdateMyProfile(ctx context.Context, req *models.Profile) error {
	db := database.GetDB()

	userID := ctx.Value("userID")
	
	_, err := db.ExecContext(ctx,
		`UPDATE profiles 
		 SET job_title = $2, years_of_experience = $3, areas_of_expertise = $4, phone_number = $5
		 WHERE user_id = $1`,
		userID, req.JobTitle, req.YearsOfExperience, pq.Array(req.Areas_of_expertise), req.PhoneNumber)
	
	return err
}