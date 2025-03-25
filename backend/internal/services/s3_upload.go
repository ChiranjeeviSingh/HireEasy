package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"backend/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadResumeToS3 uploads a resume file to S3 and returns the URL
func UploadResumeToS3(file *multipart.FileHeader, username string) (string, error) {
	cfg := config.GetConfig()

	// Check for test mode first
	if cfg.TestMode {
		return fmt.Sprintf("https://test-bucket.s3.amazonaws.com/resumes/%s_test_user.pdf", username), nil
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Create a unique filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	filename := fmt.Sprintf("resumes/%s_%s%s", username, file.Filename, ext)

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Upload the file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(filename),
		Body:   src,
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Return the public URL
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", cfg.S3Bucket, filename), nil
} 