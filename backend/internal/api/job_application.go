package api

import (
	"backend/internal/database"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// JobApplication represents the form data
type JobApplication struct {
	FirstName     string `form:"firstName" binding:"required"`
	LastName      string `form:"lastName" binding:"required"`
	Email         string `form:"email" binding:"required,email"`
	PhoneNumber   string `form:"phoneNumber" binding:"required"`
	Gender        string `form:"gender" binding:"required"`
	VeteranStatus string `form:"veteran" binding:"required"`
}

// HandleJobApplication handles job applications and file uploads
func HandleJobApplication(c *gin.Context) {
	fmt.Println("🔍 HandleJobApplication called!") // Debug log

	// Step 1: Bind form data
	var formData JobApplication
	if err := c.ShouldBind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "details": err.Error()})
		return
	}

	// Step 2: Handle file upload
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	// Step 3: Ensure "uploads" directory exists
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Step 4: Generate unique filename
	extension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("resume_%d%s", time.Now().Unix(), extension)
	filePath := filepath.Join(uploadDir, newFileName)

	// Step 5: Save the uploaded file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"})
		return
	}

	// Step 6: Store job application in PostgreSQL
	db := database.GetDB()
	_, err = db.Exec(`
		INSERT INTO job_applications (first_name, last_name, email, phone_number, resume_path, gender, veteran_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		formData.FirstName, formData.LastName, formData.Email, formData.PhoneNumber, filePath, formData.Gender, formData.VeteranStatus,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save job application"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully", "resume_path": filePath})
}
