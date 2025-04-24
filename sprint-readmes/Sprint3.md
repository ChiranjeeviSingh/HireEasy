# Sprint 3 Documentation
## Tasks Completed in Sprint 3
### Functionality Added in This Sprint - Frontend Development
1. **Authentication System**
   - Implemented JWT-based authentication
   - Created login and registration forms
   - Added token management and storage
   - Implemented protected routes
   - Added form validation and error handling
2. **Dashboard Implementation**
   - Created responsive dashboard layout
   - Implemented navigation system
   - Added quick access buttons for main features
   - Integrated with backend APIs
   - Added loading states and error handling
3. **Job Management System**
   - Implemented job posting creation
   - Added job listing view with filtering
   - Created job details view
   - Integrated with backend APIs
   - Added form validation and error handling
4. **Job Applications System**
   - Created application submission form
   - Implemented application filtering
   - Added sorting functionality
   - Integrated with ATS scoring
   - Added file upload for resumes
5. **Questionnaire System**
   - Implemented dynamic form generation
   - Added multiple question types support
   - Created form template management
   - Integrated with job applications
   - Added validation and error handling

### Testing Implementation
1. **Frontend Testing**
   - Implemented component unit tests
   - Added integration tests
   - Created API integration tests
   - Added form validation tests
   - Implemented error handling tests

## Implementation Details
### Frontend Code Examples
#### Authentication Implementation
```javascript
// Login Component
const handleLogin = async (e) => {
  e.preventDefault();
  setLoading(true);
  try {
    const response = await fetch(`${API_BASE}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    const data = await response.json();
    localStorage.setItem("token", data.token.token);
    navigate("/dashboard");
  } catch (err) {
    setError(err.message);
  }
};
```
#### Job Management Implementation
```javascript
// Job Posting Form
const handleSubmit = async (e) => {
  e.preventDefault();
  const jobDataMapped = {
    job_id: jobData.JobID.trim(),
    job_title: jobData.Info2.trim(),
    job_description: jobData.Info3.trim(),
    skills_required: jobData.Info5.split(",").map(skill => skill.trim()),
  };
  // API call implementation
};
```

## Testing Implementation
### Frontend Testing
```javascript
// Component Test Example
describe('Login Component', () => {
  it('should handle successful login', async () => {
    const mockResponse = { token: { token: 'test-token' } };
    global.fetch = jest.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockResponse)
      })
    );
    
    render(<Login />);
    fireEvent.change(screen.getByPlaceholderText(/email/i), {
      target: { value: 'test@example.com' }
    });
    fireEvent.click(screen.getByText(/sign in/i));
    
    expect(localStorage.getItem('token')).toBe('test-token');
  });
});
```


## Functionality Added in This Sprint - Backend


### Job Submission System Enhancement
The following API endpoints have been implemented for managing job submissions:

- **POST /api/jobs/:job_id/apply** - Handles job application submissions with resume upload.
- **GET /api/jobs/:job_id/submissions** - Retrieves job submissions with advanced filtering and sorting.

### Advanced ATS Scoring System
Implemented a sophisticated ATS scoring algorithm that:

- Compares candidate skills against job requirements.
- Handles exact and partial skill matches.
- Provides normalized scoring based on skill relevance.
- Maintains backward compatibility with existing submissions.

### Database Schema Improvements
- Added `form_templates` and `application_form` tables.
- Implemented proper foreign key constraints.
- Added `form_uuid` column to `job_submissions` table.
- Added UNIQUE constraint to the `username` field.
- Job Submission API Tests

### Profile (Interviewer) Routes for managing interview portal
The following API endpoints have been implemented for managing interviewer profiles:

- **POST /profiles** - Creates a new profile.
- **PUT /profiles** - Updates the authenticated user's profile.
- **GET /profiles/me** - Retrieves the authenticated user's profile using JWT token.
- **GET /profiles/user/:user_name** - Retrieves a specific user's profile using their username.

### Availability Routes
The following API endpoints have been implemented for managing interviewer availability:

- **POST /availability** - Creates a new availability slot (interviewer action).
- **DELETE /availability/:id** - Deletes an existing availability slot (interviewer action).
- **GET /availability/me** - Retrieves all availability slots for the authenticated user (with optional date range filtering).
- **GET /availability/user/:user_name** - Retrieves a specific user's availability slots using their username.
- **GET /availability** - Retrieves all available slots (with optional date range and profile filters).

## API Documentation
For more details on API endpoints, please refer to the [Backend Documentation](https://documenter.getpostman.com/view/41938964/2sB2cRCQ89).

<table>
<tr>
<th>Test Case</th>
<th>Description</th>
<th>Expected Outcome</th>
</tr>
<tr>
<td><code>TestHandleFormSubmission</code></td>
<td>Tests complete job application submission process</td>
<td>- HTTP Status <code>201 Created</code><br>- Form data correctly processed<br>- Resume uploaded to S3<br>- ATS score calculated and stored</td>
</tr>
<tr>
<td><code>TestGetFormSubmissions</code></td>
<td>Tests retrieval of job submissions with various filters</td>
<td>- HTTP Status <code>200 OK</code><br>- Submissions returned with correct sorting<br>- Filter parameters properly applied</td>
</tr>
<tr>
<td><code>TestFormValidation</code></td>
<td>Tests validation of required form fields</td>
<td>- HTTP Status <code>400 Bad Request</code><br>- Appropriate error messages returned<br>- Invalid submissions rejected</td>
</tr>
<tr>
<td><code>TestResumeUpload</code></td>
<td>Tests resume file upload functionality</td>
<td>- File successfully uploaded to S3<br>- Correct URL returned<br>- File type validation works</td>
</tr>
</table>
ATS Scoring Tests
<table>
<tr>
<th>Test Case</th>
<th>Description</th>
<th>Expected Outcome</th>
</tr>
<tr>
<td><code>TestATSScoreCalculation</code></td>
<td>Tests ATS scoring algorithm accuracy</td>
<td>- Perfect skill matches score highest<br>- Partial matches receive proportional scores<br>- Case-insensitive matching works</td>
</tr>
<tr>
<td><code>TestSkillNormalization</code></td>
<td>Tests skill text normalization</td>
<td>- Different skill text formats match correctly<br>- Special characters handled properly<br>- Case variations normalized</td>
</tr>
<tr>
<td><code>TestEmptySkillsHandling</code></td>
<td>Tests handling of missing or empty skills</td>
<td>- System handles empty skill lists gracefully<br>- Appropriate default scores assigned<br>- Error messages clear and helpful</td>
</tr>
</table>
Database Integration Tests
<table>
<tr>
<th>Test Case</th>
<th>Description</th>
<th>Expected Outcome</th>
</tr>
<tr>
<td><code>TestJobSubmissionStorage</code></td>
<td>Tests database storage of submissions</td>
<td>- Submission data correctly stored<br>- Foreign key constraints maintained<br>- Form UUID properly linked</td>
</tr>
<tr>
<td><code>TestSchemaUpdates</code></td>
<td>Tests database schema modifications</td>
<td>- Schema updates applied successfully<br>- Existing data preserved<br>- New constraints enforced</td>
</tr>
<tr>
<td><code>TestQueryPerformance</code></td>
<td>Tests performance of submission queries</td>
<td>- Queries execute within acceptable time<br>- Indexes properly utilized<br>- Large result sets handled efficiently</td>
</tr>
</table>
API Filtering and Sorting Tests
<table>
<tr>
<th>Test Case</th>
<th>Description</th>
<th>Expected Outcome</th>
</tr>
<tr>
<td><code>TestATSScoreSorting</code></td>
<td>Tests sorting by ATS score</td>
<td>- Submissions correctly sorted by score<br>- Higher scores appear first<br>- Ties handled appropriately</td>
</tr>
<tr>
<td><code>TestDateFiltering</code></td>
<td>Tests date-based filtering</td>
<td>- Only submissions within date range returned<br>- Date formats handled correctly<br>- Timezone issues resolved</td>
</tr>
<tr>
<td><code>TestResultLimiting</code></td>
<td>Tests result pagination and limiting</td>
<td>- Correct number of results returned<br>- Pagination parameters respected<br>- Total count included in response</td>
</tr>
</table>
Security Tests
<table>
<tr>
<th>Test Case</th>
<th>Description</th>
<th>Expected Outcome</th>
</tr>
<tr>
<td><code>TestAuthentication</code></td>
<td>Tests authentication requirements</td>
<td>- Unauthenticated requests rejected<br>- JWT tokens properly validated<br>- Token expiration handled</td>
</tr>
<tr>
<td><code>TestAuthorization</code></td>
<td>Tests authorization rules</td>
<td>- Users can only access their own data<br>- Admin access properly enforced<br>- Role-based permissions respected</td>
</tr>
<tr>
<td><code>TestInputSanitization</code></td>
<td>Tests input validation and sanitization</td>
<td>- Malicious input detected and rejected<br>- SQL injection prevented<br>- XSS attacks blocked</td>
</tr>
</table>

## Profile Management Tests

<table>
  <tr>
    <th>Test Case</th>
    <th>Description</th>
    <th>Expected Outcome</th>
  </tr>
  <tr>
    <td><code>TestCreateProfileH</code></td>
    <td>Tests creating a new user profile</td>
    <td>- HTTP Status <code>201 Created</code><br>- Profile details match input data<br>- Profile is stored in database</td>
  </tr>
  <tr>
    <td><code>TestUpdateProfileH</code></td>
    <td>Tests updating an existing profile</td>
    <td>- HTTP Status <code>200 OK</code><br>- Updated profile details match input<br>- Changes are reflected in database</td>
  </tr>
  <tr>
    <td><code>TestGetMyProfileH</code></td>
    <td>Tests retrieving profile for authenticated user</td>
    <td>- HTTP Status <code>200 OK</code><br>- Profile details match user's stored profile</td>
  </tr>
  <tr>
    <td><code>TestGetUserProfileH</code></td>
    <td>Tests retrieving a specific user's profile by username</td>
    <td>- HTTP Status <code>200 OK</code><br>- Profile details match requested user's profile</td>
  </tr>
</table>

## Availability Management Tests

File: <code>availability_test.go</code>

<table>
  <tr>
    <th>Test Case</th>
    <th>Description</th>
    <th>Expected Outcome</th>
  </tr>
  <tr>
    <td><code>TestCreateAvailabilityH</code></td>
    <td>Tests creating a new availability slot as an interviewer</td>
    <td>- HTTP Status <code>201 Created</code><br>- Response contains availability details (ID, date, time range)<br>- Slot is stored in database</td>
  </tr>
  <tr>
    <td><code>TestDeleteAvailabilityH</code></td>
    <td>Tests deleting an existing availability slot</td>
    <td>- HTTP Status <code>200 OK</code><br>- Success message returned<br>- Slot is removed from database</td>
  </tr>
  <tr>
    <td><code>TestGetMyAvailabilityH</code></td>
    <td>Tests retrieving all availability slots for the authenticated user</td>
    <td>- HTTP Status <code>200 OK</code><br>- Response contains array of availability slots<br>- All fields properly formatted</td>
  </tr>
  <tr>
    <td><code>TestGetUserAvailabilityH</code></td>
    <td>Tests retrieving availability slots for a specific user by username</td>
    <td>- HTTP Status <code>200 OK</code><br>- Response contains array of matching user's slots<br>- Fields include date and time range</td>
  </tr>
  <tr>
    <td><code>TestGetAllAvailabilityH</code></td>
    <td>Tests retrieving all available slots (potentially filtered)</td>
    <td>- HTTP Status <code>200 OK</code><br>- Response contains array of available slots<br>- Each entry includes username and correct slot details</td>
  </tr>
</table>


## Video Link:
https://drive.google.com/file/d/1SvmwWFFKe_63g9jpGwtI1uciVjbXSDIl/view?usp=sharing
