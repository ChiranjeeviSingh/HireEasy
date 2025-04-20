# HireEasy - Job Application System

HireEasy is a comprehensive job application management system that allows employers to create job listings, customize application forms, and track candidate submissions. The system includes an ATS (Applicant Tracking System) score calculator to help identify the best-fit candidates.

## Features

- Job listing management
- Customizable application forms
- Resume upload and storage
- Automatic ATS scoring based on skills
- Filtering and sorting of job applications
- User authentication and management

## System Architecture

The application consists of:

- **Backend**: Go server with Gin framework
- **Database**: PostgreSQL
- **Storage**: AWS S3 for resume storage
- **Frontend**: (Not included in this repository)

## Database Structure

- **users**: Store user information (employers and applicants)
- **jobs**: Job listings with details
- **form_templates**: Customizable application forms for each job
- **job_submissions**: Applications submitted by candidates
- **job_applications**: Detailed application information

## Running the Application

### Prerequisites

- Go 1.16 or newer
- PostgreSQL database
- AWS account (for S3 storage)

### Environment Variables

The application requires the following environment variables:

```bash
# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=app_db

# AWS Configuration
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
S3_BUCKET=your_bucket_name

# Optional Test Mode Configuration
# Set to "true" to bypass actual S3 uploads
S3_TEST_MODE=true
```

### Running the Server

Navigate to the backend directory and run the server:

```bash
cd backend
go run cmd/server/main.go
```

Alternatively, set the required environment variables inline:

```bash
cd backend
AWS_ACCESS_KEY_ID=your_access_key \
AWS_SECRET_ACCESS_KEY=your_secret_key \
AWS_REGION=us-east-1 \
S3_BUCKET=your_bucket_name \
go run cmd/server/main.go
```

For testing without S3 uploads:

```bash
cd backend
S3_TEST_MODE=true go run cmd/server/main.go
```

### API Endpoints

#### Job Submissions

- **POST** `/api/forms/:form_uuid/submissions`: Submit a job application
  - Required fields: job_id, user_id, username, email, form_data, resume
  - Example:
    ```
    curl -X POST http://localhost:8080/api/forms/4d9a4320-f1d1-43f2-8477-edd07f557442/submissions \
      -F "job_id=5678" \
      -F "user_id=2003" \
      -F "username=Test User" \
      -F "email=testuser@example.com" \
      -F "form_data={\"experience\":\"5 years\", \"location\":\"Remote\", \"skills\": [\"Go\", \"AWS\"]}" \
      -F "resume=@resume.pdf"
    ```

- **GET** `/api/forms/:form_uuid/submissions`: Get submissions for a form
  - Optional query parameters:
    - `sort_by`: Field to sort by (ats_score or created_at)
    - `limit`: Maximum number of results to return
    - `date`: Filter by date (all, today, or YYYY-MM-DD format)
  - Example:
    ```
    curl http://localhost:8080/api/forms/4d9a4320-f1d1-43f2-8477-edd07f557442/submissions?sort_by=ats_score&limit=10&date=all
    ```

## ATS Scoring System

The ATS scoring system calculates a score based on:
- Base score of 70 points
- Additional points for each skill (5 points per skill)
- Maximum score is capped at 100

## Troubleshooting

### S3 Upload Issues

If you encounter "MissingRegion" errors:
1. Ensure AWS_REGION is properly set
2. Verify AWS credentials are valid
3. Consider using S3_TEST_MODE=true for testing without S3 uploads

### Database Connection Issues

If you encounter database connection errors:
1. Verify PostgreSQL is running
2. Check database credentials
3. Ensure database and required tables exist
