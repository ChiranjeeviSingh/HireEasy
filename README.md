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



## Detailed Component Descriptions

### Authentication Components

#### 1. Login
**Purpose:** Provides user authentication functionality.  
**Features:**
- Email and password-based login
- Role-based redirection (HR to HR dashboard, Interviewer to Interviewer dashboard)
- Error handling and loading states
- Navigation to registration for new users

#### 2. Register
**Purpose:** Allows new users to create accounts.  
**Features:**
- Registration form with validation
- Role selection (HR or Interviewer)
- Account creation with appropriate permissions
- Redirection to login after successful registration

### HR Components

#### 3. Dashboard
**Purpose:** Main navigation hub for HR professionals.  
**Features:**
- Central menu for accessing all HR functions
- Navigation buttons for job management, questionnaires, applications, and interview scheduling
- Clean, user-friendly interface

#### 4. JobPosting
**Purpose:** Enables HR to create and publish new job openings.  
**Features:**
- Form for creating detailed job descriptions
- Fields for title, description, requirements, and other job details
- Submission to backend API for storage
- Validation of required fields

#### 5. Questionnaire
**Purpose:** Allows HR to create custom application forms and screening questionnaires.  
**Features:**
- Dynamic form builder for creating questions
- Support for multiple question types (multiple choice, text, etc.)
- Ability to edit and reorder questions
- Form ID generation for linking forms to job postings

#### 6. ShareJob
**Purpose:** Facilitates sharing job postings with potential candidates.  
**Features:**
- Interface to select job postings and associate them with questionnaires
- Generation of unique application links for candidates
- Options for sharing methods (link, email, etc.)
- Management of active shared job listings

#### 7. ViewJobs
**Purpose:** Provides an overview of all posted jobs.  
**Features:**
- List view of all created job postings
- Filter functionality for job searches
- Basic job statistics
- Navigation to detailed job views

#### 8. JobApplications
**Purpose:** Manages incoming job applications.  
**Features:**
- List of all received applications
- Filtering and sorting capabilities
- Review interface for applications
- Status management (new, reviewed, interview scheduled, etc.)

#### 9. ScheduleInterviews
**Purpose:** Facilitates the interview scheduling process.  
**Features:**
- Interface to match candidates with interviewers
- Date and time selection based on availability
- Email notification system
- Interview status tracking

### Interviewer Components

#### 10. InterviewerDashboard
**Purpose:** Main navigation hub for interviewers.  
**Features:**
- Access to profile management
- Calendar viewing and availability setting
- Interface for accessing scheduled interviews
- Simple, focused design for interviewer needs

#### 11. InterviewerProfile
**Purpose:** Allows interviewers to manage their professional profiles.  
**Features:**
- Personal information management
- Professional details and expertise areas
- Availability preferences
- Contact information updates

#### 12. InterviewerCalendar
**Purpose:** Manages interviewer availability scheduling.  
**Features:**
- Calendar interface for setting available time slots
- Recurring availability options
- Integration with scheduling system
- View of upcoming scheduled interviews

#### 13. InterviewerInterviews
**Purpose:** Provides interviewers with access to their assigned interviews.  
**Features:**
- List of scheduled interviews
- Candidate information and application details
- Interview notes and feedback entry
- Status updates for completed interviews

### Applicant Components

#### 14. Apply
**Purpose:** Serves as the application interface for job seekers.  
**Features:**
- Job information display
- Dynamic form rendering based on the associated questionnaire
- File upload for resumes and supporting documents
- Submission confirmation and tracking

## Detailed Project Flow

### 1. Authentication Flow

Both portals begin with a common authentication path:

1. **Login Page** (`Login`):
   - Users enter credentials (email and password)
   - The system validates credentials and determines user role
   - Based on role, users are redirected to their respective dashboards
   
2. **Registration Page** (`Register`):
   - New users create accounts by providing personal information
   - Users select their role (HR or Interviewer)
   - Account creation initiates with appropriate permissions
   - Users are redirected to login after successful registration

### 2. HR Portal Flow

After authentication, the HR workflow follows this path:

1. **HR Dashboard** (`Dashboard`):
   - Central navigation hub with access to all HR functions
   - Quick access buttons to key features

2. **Job Posting** (`JobPosting`):
   - HR creates detailed job descriptions
   - Specifies title, requirements, responsibilities, and other key details
   - Publishes listings to the job database

3. **Questionnaire Creation** (`Questionnaire`):
   - HR designs custom application forms
   - Creates screening questions with various formats (multiple choice, text, etc.)
   - Generates a unique form ID for linking to job postings

4. **Job Sharing** (`ShareJob`):
   - HR associates job postings with appropriate questionnaires
   - Generates application links for candidates
   - Manages distribution of job opportunities

5. **Application Review** (`JobApplications`):
   - HR views and evaluates incoming applications
   - Filters applications by various criteria
   - Reviews candidate qualifications and responses
   - Marks promising candidates for interviews

6. **Interview Scheduling** (`ScheduleInterviews`):
   - HR matches candidates with appropriate interviewers
   - Schedules interviews based on interviewer availability
   - Sends notifications to all parties
   - Tracks interview status

7. **Job Monitoring** (`ViewJobs`):
   - HR reviews all active job postings
   - Tracks application statistics
   - Updates or closes positions as needed

### 3. Interviewer Portal Flow

The interviewer workflow includes:

1. **Interviewer Dashboard** (`InterviewerDashboard`):
   - Main hub for interviewer activities
   - Navigation to profile, calendar, and interview sections

2. **Profile Management** (`InterviewerProfile`):
   - Interviewers maintain their professional profiles
   - Specify areas of expertise and qualifications
   - Update contact information and preferences

3. **Availability Setting** (`InterviewerCalendar`):
   - Interviewers define their available time slots
   - Set up recurring availability patterns
   - View their already scheduled interviews
   - Block out unavailable periods

4. **Interview Management** (`InterviewerInterviews`):
   - Access list of assigned interviews
   - Review candidate information before interviews
   - Enter feedback and evaluation after interviews
   - Update interview status (completed, rescheduled, etc.)

### 4. Candidate Flow

While not a separate portal, candidates interact with the system through:

1. **Application Process** (`Apply`):
   - Candidates access job listings through shared links
   - View job details and requirements
   - Complete the questionnaire specific to that position
   - Upload supporting documents
   - Submit applications and receive confirmation

2. **Post-Application**:
   - Candidates may receive notifications about their application status
   - Interview invitations are sent if selected
   - Scheduling confirmations are provided

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