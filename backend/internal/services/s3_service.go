package services

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadResumeToS3 uploads the resume to AWS S3 and returns the file URL
func UploadResumeToS3(file *multipart.FileHeader, userID int, username string) (string, error) {
	// Check for test mode first
	testMode := os.Getenv("TEST_MODE") == "true" || os.Getenv("S3_TEST_MODE") == "true"
	
	log.Printf("DEBUG - Environment TEST_MODE: %s, S3_TEST_MODE: %s", os.Getenv("TEST_MODE"), os.Getenv("S3_TEST_MODE"))
	log.Printf("DEBUG - Is test mode: %t", testMode)
	
	// Get AWS configuration
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		// Default region if not set
		awsRegion = "us-east-1"
		log.Printf("Warning: AWS_REGION not set, using default: %s", awsRegion)
	}
	
	bucketName := os.Getenv("S3_BUCKET")
	if bucketName == "" {
		bucketName = "test-bucket"
		log.Printf("Warning: S3_BUCKET not set, using default: %s", bucketName)
	}
	
	log.Printf("Debug - AWS_REGION: %s, S3_BUCKET: %s, TEST_MODE: %t", awsRegion, bucketName, testMode)
	
	// Generate sanitized filename
	sanitizedUsername := sanitizeFilename(username)
	fileName := fmt.Sprintf("%d_%s%s", userID, sanitizedUsername, filepath.Ext(file.Filename))
	s3Path := fmt.Sprintf("resumes/%s", fileName)
	
	// Create the S3 URL that would be used (for both test and real modes)
	resumeURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Path)
	
	// In test mode, just return the URL without uploading
	if testMode {
		log.Println("Test mode: Using mock S3 URL:", resumeURL)
		return resumeURL, nil
	}
	
	// For real S3 uploads
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Println("AWS Session Error:", err)
		return "", fmt.Errorf("AWS session error: %w", err)
	}

	svc := s3.New(sess)

	// Open file
	src, err := file.Open()
	if err != nil {
		log.Println("Error Opening File:", err)
		return "", fmt.Errorf("file open error: %w", err)
	}
	defer src.Close()

	// Read file content
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		log.Println("Error Reading File Content:", err)
		return "", fmt.Errorf("file read error: %w", err)
	}

	// Upload to S3
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Path),
		Body:   bytes.NewReader(buf.Bytes()),
	}
	log.Printf("Debug - PutObjectInput: Bucket=%s, Key=%s", *putObjectInput.Bucket, *putObjectInput.Key)
	
	result, err := svc.PutObject(putObjectInput)
	if err != nil {
		log.Println("S3 Upload Error:", err)
		return "", fmt.Errorf("S3 upload error: %w", err)
	}
	log.Printf("Debug - PutObject successful: %v", result)

	log.Println("Resume Uploaded to S3:", resumeURL)
	return resumeURL, nil
}

// sanitizeFilename ensures filenames are safe for S3 storage
func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return re.ReplaceAllString(name, "_")
}
