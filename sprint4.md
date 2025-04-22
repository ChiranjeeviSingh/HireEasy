# Sprint 4 Documentation

## Overview

HireEasy is a comprehensive hiring platform that streamlines the recruitment process through automated job posting, candidate screening, and interview management. This sprint focuses on the interviewer workflow and job application management system.

## Tasks Completed in Sprint 4

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
  "status": "under_review" // Options: applied, under_review, accepted, rejected
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

### Functionality Added in This Sprint - Frontend

#### Core Components

##### 1. Interviewer Workflow Components

###### InterviewerProfile

- **Purpose**: Manages interviewer professional information
- **Features**:
  - Profile creation and updates
  - Job title and experience management
  - Areas of expertise specification
  - Contact information management
- **Key Functions**:
  - Profile data persistence
  - Form validation
  - Real-time updates

###### InterviewerCalendar

- **Purpose**: Manages interviewer availability
- **Features**:
  - Time slot management (13:00-18:00)
  - 7-day availability view
  - Slot creation and deletion
- **Key Functions**:
  - Availability scheduling
  - Conflict prevention
  - Real-time updates

###### InterviewerInterviews

- **Purpose**: Manages scheduled interviews
- **Features**:
  - Interview scheduling
  - Candidate feedback system
  - Status updates
  - Interview details view
- **Key Functions**:
  - Feedback submission
  - Candidate selection/rejection
  - Interview status tracking

##### 2. Job Application Management

###### JobApplications

- **Purpose**: Manages job applications and candidate screening
- **Features**:
  - Application status tracking
  - Candidate filtering
  - ATS score integration
  - Resume management
- **Key Functions**:
  - Status updates
  - Candidate shortlisting
  - Application review

###### ScheduleInterviews

- **Purpose**: Coordinates interview scheduling
- **Features**:
  - Interview slot management
  - Candidate selection
  - Interviewer assignment
  - Status tracking
- **Key Functions**:
  - Interview scheduling
  - Status updates
  - Conflict resolution

#### Added Functionalities

##### InterviewerProfile.jsx

- **New Features**:
  - Professional profile creation and management
  - Expertise area tagging system
  - Experience tracking
  - Contact information management
- **Integration Points**:
  - Connects with backend profile API
  - Updates interviewer availability preferences
  - Links with interview scheduling system

##### InterviewerCalendar.jsx

- **New Features**:
  - Interactive calendar interface
  - Time slot management (13:00-18:00)
  - 7-day availability view
  - Conflict detection
- **Integration Points**:
  - Syncs with backend availability API
  - Integrates with interview scheduling
  - Real-time updates for booked slots

##### InterviewerInterviews.jsx

- **New Features**:
  - Interview dashboard
  - Feedback submission system
  - Candidate evaluation interface
  - Status tracking
- **Integration Points**:
  - Connects with interview management API
  - Links with candidate profiles
  - Updates job submission status

##### JobApplications.jsx

- **New Features**:
  - Application tracking dashboard
  - ATS score visualization
  - Candidate filtering system
  - Status management
- **Integration Points**:
  - Job submission API
  - Candidate profile system
  - Interview scheduling

##### ScheduleInterviews.jsx

- **New Features**:
  - Interview scheduling interface
  - Interviewer availability matching
  - Candidate selection system
  - Conflict resolution
- **Integration Points**:
  - Interview management API
  - Availability system
  - Candidate profiles

### Frontend Testing Framework

#### Component Tests

The application includes comprehensive tests for all major components:

1. **Login Component Tests**

   - Form validation
   - Authentication flow
   - UI rendering

2. **JobPosting Component Tests**

   - Form validation
   - Job creation flow
   - UI elements

3. **Questionnaire Component Tests**

   - Question management
   - Option handling
   - Form validation

4. **Register Component Tests**

   - User registration flow
   - Role selection
   - Form validation

5. **ViewJobs Component Tests**

   - Job listing
   - Filtering functionality
   - UI rendering

6. **ShareJob Component Tests**

   - Job sharing functionality
   - Form generation
   - UI elements

7. **InterviewerCalendar Component Tests**
   - Availability management
   - Slot creation/deletion
   - UI validation

#### Test Infrastructure

- Uses React Testing Library
- Mock implementations for:
  - localStorage
  - fetch API
  - Alert system
- Comprehensive test coverage for:
  - UI rendering
  - User interactions
  - Form validation
  - Data handling

### Cypress Integration Tests

#### Setup and Installation

1. Install Cypress:

```bash
cd frontend
npm install cypress --save-dev
```

2. Add Cypress scripts to package.json:

```json
{
  "scripts": {
    "cypress:open": "cypress open",
    "cypress:run": "cypress run"
  }
}
```

3. Create Cypress configuration file (cypress.config.js):

```javascript
const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    baseUrl: "http://localhost:3000",
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
```

#### Test Structure

Create the following test files in `frontend/cypress/e2e/`:

1. **Authentication Tests** (`auth.cy.js`):

```javascript
describe("Authentication Flow", () => {
  it("should login successfully", () => {
    cy.visit("/login");
    cy.get("[data-cy=email]").type("test@example.com");
    cy.get("[data-cy=password]").type("password123");
    cy.get("[data-cy=login-button]").click();
    cy.url().should("include", "/dashboard");
  });

  it("should handle invalid credentials", () => {
    cy.visit("/login");
    cy.get("[data-cy=email]").type("invalid@example.com");
    cy.get("[data-cy=password]").type("wrongpassword");
    cy.get("[data-cy=login-button]").click();
    cy.get("[data-cy=error-message]").should("be.visible");
  });
});
```

2. **Interviewer Profile Tests** (`interviewer-profile.cy.js`):

```javascript
describe("Interviewer Profile Management", () => {
  beforeEach(() => {
    cy.login("interviewer@example.com", "password123");
  });

  it("should update profile information", () => {
    cy.visit("/interviewer-profile");
    cy.get("[data-cy=job-title]").type("Senior Developer");
    cy.get("[data-cy=experience]").type("5");
    cy.get("[data-cy=expertise]").type("JavaScript, React, Node.js");
    cy.get("[data-cy=save-profile]").click();
    cy.get("[data-cy=success-message]").should("be.visible");
  });
});
```

3. **Interview Scheduling Tests** (`interview-scheduling.cy.js`):

```javascript
describe("Interview Scheduling", () => {
  beforeEach(() => {
    cy.login("hr@example.com", "password123");
  });

  it("should schedule an interview", () => {
    cy.visit("/schedule-interviews");
    cy.get("[data-cy=select-candidate]").select("John Doe");
    cy.get("[data-cy=select-interviewer]").select("Jane Smith");
    cy.get("[data-cy=select-date]").type("2024-04-25");
    cy.get("[data-cy=select-time]").type("14:00");
    cy.get("[data-cy=schedule-button]").click();
    cy.get("[data-cy=success-message]").should("be.visible");
  });
});
```

4. **Interviewer Calendar Tests** (`interviewer-calendar.cy.js`):

```javascript
describe("Interviewer Calendar Management", () => {
  beforeEach(() => {
    cy.login("interviewer@example.com", "password123");
  });

  it("should add availability slot", () => {
    cy.visit("/interviewer-calendar");
    cy.get("[data-cy=add-slot]").click();
    cy.get("[data-cy=date-input]").type("2024-04-25");
    cy.get("[data-cy=start-time]").type("13:00");
    cy.get("[data-cy=end-time]").type("14:00");
    cy.get("[data-cy=save-slot]").click();
    cy.get("[data-cy=slot-list]").should("contain", "13:00 - 14:00");
  });
});
```

#### Running Tests

1. Start the development server:

```bash
cd frontend
npm start
```

2. Open Cypress Test Runner:

```bash
npm run cypress:open
```

3. Run tests in headless mode:

```bash
npm run cypress:run
```

#### Test Commands

```bash
# Install Cypress
npm install cypress --save-dev

# Add Cypress scripts to package.json
npm pkg set scripts.cypress:open="cypress open"
npm pkg set scripts.cypress:run="cypress run"

# Initialize Cypress
npx cypress open

# Run specific test file
npx cypress run --spec "cypress/e2e/auth.cy.js"

# Run tests in specific browser
npx cypress run --browser chrome

# Run tests with video recording
npx cypress run --record --key YOUR_RECORD_KEY
```

#### Best Practices

1. Use data-cy attributes for selectors
2. Implement custom commands for common actions
3. Use fixtures for test data
4. Implement proper cleanup after tests
5. Use environment variables for sensitive data

#### Integration with CI/CD

Add to your CI/CD pipeline:

```yaml
steps:
  - name: Run Cypress Tests
    run: |
      cd frontend
      npm install
      npm run cypress:run
```

## Technical Stack

### Frontend

- React.js
- React Router
- Testing Library
- Tailwind CSS

### Backend

- Go
- Gin Framework
- PostgreSQL
- JWT Authentication

## Getting Started

### Frontend Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```
3. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
4. Start the development server:
   ```bash
   npm start
   ```

### Backend Setup

1. Navigate to backend directory
2. Set up environment variables
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

### Running Tests

```bash
# Frontend Tests
cd frontend
npm test

# Backend Tests
cd backend
go test -v ./...
```

## Security Features

- JWT-based authentication
- Role-based access control
- Secure data transmission
- Input validation
- XSS protection

## Performance Considerations

- Lazy loading of components
- Optimized API calls
- Efficient state management
- Responsive design
- Database indexing
- Caching strategies

## API Documentation

For more details on API endpoints, please refer to the [Backend Documentation](https://documenter.getpostman.com/view/41938964/2sB2cRCQ89).

## Video Demo

[Sprint 4 Demo Video](https://youtu.be/demo-link)

## Future Enhancements

1. Real-time notifications
2. Advanced analytics
3. Integration with external ATS
4. Enhanced reporting features
5. Mobile application support

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
