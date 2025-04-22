# HireEasy - Job Application System

HireEasy is a comprehensive job application management system that allows employers to create job listings, customize application forms, and track candidate submissions. The system includes an ATS (Applicant Tracking System) score calculator to help identify the best-fit candidates.

## Features

HireEasy is structured around two main user portals, each tailored to specific roles:

### 1. HR Portal

The HR Portal serves as the central hub for all recruitment activities. HR professionals use this portal to:

- Create and manage job postings
- Design custom application questionnaires
- Share job opportunities with candidates
- Review incoming applications
- Schedule interviews with promising candidates
- Track the entire hiring pipeline

### 2. Interviewer Portal

The Interviewer Portal focuses on the evaluation phase of the hiring process. Interviewers use this portal to:

- Manage their professional profiles and expertise areas
- Set their availability for conducting interviews
- View their assigned interviews and candidate details
- Provide structured feedback on candidates
- Track interview outcomes

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
- React.js
- React Router
- Cypress
- Jest
- React Testing Library 

## Setup Instructions

1. Clone the repository:
   ```
   git clone https://github.com/ChiranjeeviSingh/HireEasy/blob/main/backend/README.md
   ```

2. Install frontend dependencies**:
   ```bash
   cd frontend
   npm install
   ```

3. Start the frontend server**:
   ```bash
   npm start
   ```

4. Navigate to the backend directory:
   ```
   cd HireEasy/backend
   ```

5. Install dependencies:
   ```
   go mod tidy
   ```

6. Set up the PostgreSQL database using the provided schema:
   ```
   psql -U <username> -d <database> -f internal/database/migrations/schema.sql
   ```


7. Setup Environment Variables:

  ```bash
  # PostgreSQL Configuration
  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=your_username
  DB_PASSWORD=your_password
  DB_NAME=app_db
  ```

  Optional Test Mode Configuration: Set to "true" to bypass actual S3 uploads
  ```bash
  S3_TEST_MODE=true
  ```

8. Run the backend application:
  ```bash
  export TEST_MODE=false S3_TEST_MODE=false \
  AWS_REGION=us-east-1 \
  S3_BUCKET=hireeasy-resumes \
  DB_USER=reshma \
  DB_PASSWORD=postgres \
  && go run cmd/server/main.go
  ```

## Testing

1. Run backend unittests:
   ```
   go test ./test/handlers/ -v
   ```

2. Frontend Testing:
  The project includes both unit tests (using Jest and React Testing Library) and end-to-end tests (using Cypress).

  ### Unit Tests
  Unit tests focus on individual components and functionality.
  ```
  npm test
  ```

  ### E2E Tests
  End-to-end tests simulate real user interactions with the application.
  ```
  npm run cypress:run
  ```

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

## API Documentation
For more details on API endpoints, please refer to the [Backend Documentation](https://documenter.getpostman.com/view/41938964/2sB2cRCQ89).