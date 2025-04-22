# Sprint 4 Documentation
## Tasks Completed in Sprint 4

### Functionality Added in This Sprint - Frontend 

## Frontend Component Descriptions

#### 1. ScheduleInterviews.jsx
**Purpose:** Facilitates the interview scheduling process.  
**Features:**
- Interface to match candidates with interviewers
- Date and time selection based on availability
- Email notification system
- Interview status tracking

### Interviewer Components

#### 2. InterviewerDashboard.jsx
**Purpose:** Main navigation hub for interviewers.  
**Features:**
- Access to profile management
- Calendar viewing and availability setting
- Interface for accessing scheduled interviews
- Simple, focused design for interviewer needs

#### 3. InterviewerProfile.jsx
**Purpose:** Allows interviewers to manage their professional profiles.  
**Features:**
- Personal information management
- Professional details and expertise areas
- Availability preferences
- Contact information updates

#### 4. InterviewerCalendar.jsx
**Purpose:** Manages interviewer availability scheduling.  
**Features:**
- Calendar interface for setting available time slots
- Recurring availability options
- Integration with scheduling system
- View of upcoming scheduled interviews

#### 5. InterviewerInterviews.jsx
**Purpose:** Provides interviewers with access to their assigned interviews.  
**Features:**
- List of scheduled interviews
- Candidate information and application details
- Interview notes and feedback entry
- Status updates for completed interviews

## Testing
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

### Functionality Added in This Sprint - Backend 

#### Interview Management Routes (Reshma)
The following API endpoints have been implemented for managing interviews:

- **POST /interviews** - Creates a new interview by reserving an availability slot (HR action)
  - Validates availability exists and is not already booked
  - Links interview to job submission and interviewer
  - Returns interview details including scheduled time

- **DELETE /interviews/:id** - Cancels an existing interview booking (HR action)
  - Validates interview exists and belongs to HR's company
  - Frees up the availability slot
  - Returns success confirmation

- **GET /interviews** - Lists all interviews with filtering options
  - For HR users: Shows all interviews they have schedules across interviewers in the company. Also includes feedback and verdict provided by interviewer post interview.
  - For Interviewers: Shows only their assigned interviews
  - Supports filtering by:
    - Candidate's job submission ID
    - Date range
    - Interview status
  - Returns interviews with candidate and scheduling details

- **POST /interviews/:id/feedback** - Submits feedback for completed interview (Interviewer action)  
  - Validates interview exists and interviewer was assigned
  - Stores detailed feedback and updates interview status
  - Returns success confirmation

### Interview Management Tests

<table>
<tr>
<th>Test Case</th>
<th>Description</th>
<th>Scenarios Covered</th>
<th>Expected Outcome</th>
</tr>

<tr>
<td><code>TestCreateInterview</code></td>
<td>Tests interview creation process with various inputs</td>
<td>
- Missing required fields<br>
- Invalid job_id<br>
- Valid request with all fields
</td>
<td>
- HTTP 400 for missing fields<br>
- HTTP 400 for invalid job_id<br>
- HTTP 201 with interview details for valid request<br>
- Validates job_id, hr_user_id, status fields
</td>
</tr>

<tr>
<td><code>TestDeleteInterviewH</code></td>
<td>Tests interview deletion functionality</td>
<td>
- Non-existent interview ID<br>
- Valid interview deletion
</td>
<td>
- HTTP 400 for non-existent ID<br>
- HTTP 200 for successful deletion<br>
- Verifies interview removal from database
</td>
</tr>

<tr>
<td><code>TestSubmitFeedbackH</code></td>
<td>Tests interview feedback submission</td>
<td>
- Non-existent interview ID<br>
- Valid feedback submission with verdict and comments
</td>
<td>
- HTTP 400 for non-existent interview<br>
- HTTP 200 for valid submission<br>
- Verifies status update to "completed"<br>
- Verifies feedback storage
</td>
</tr>

<tr>
<td><code>TestListAllInterviewsH</code></td>
<td>Tests interview listing with different filters and roles</td>
<td>
- List all interviews for HR<br>
- Filter by job_submission_id<br>
- Filter by date range<br>
- List interviews for specific interviewer<br>
- Status checks for past/future interviews
</td>
<td>
- Returns correct number of interviews<br>
- Filters applied correctly<br>
- Past interviews marked as "pending_feedback"<br>
- Future interviews marked as "scheduled"<br>
- Role-based access control
</td>
</tr>
</table>

#### Test Setup Details
- Uses embedded test database
- Creates test users (HR and Interviewers)
- Sets up test jobs, submissions, and availabilities
- Handles proper test data cleanup

Run tests with:
```bash
go test -v HireEasy/backend/test/handlers/interview_test.go
```

### Job Submission Status Updates (Sushmita)
#### Status Management
#### 1. Update Submission Status
```
PUT /api/submissions/:submission_id/status
```
Update a submission's status with:
```json
{
    "status": "under_review"  // Options: applied, under_review, accepted, rejected
}
```
#### 2. Filter Submissions by Status
```
GET /api/jobs/:job_id/submissions?status=under_review
```
Query parameter:
- `status`: Filter submissions by status (applied, under_review, accepted, rejected)

### Test Coverage
We've added comprehensive tests for:
1. Status Updates
   - Successfully updating status
   - Invalid status values
   - Non-existent submissions
2. Status Filtering
   - Filtering by valid status
   - Invalid status parameters
   - Empty results for non-matching status

Run tests with:
```bash
go test -v ./internal/api/handlers/job_submission_test.go
```

## API Documentation
For more details on API endpoints, please refer to the [Backend Documentation](https://documenter.getpostman.com/view/41938964/2sB2cRCQ89).

## Video Demo
[Sprint 4 Demo Video](https://youtu.be/demo-link)