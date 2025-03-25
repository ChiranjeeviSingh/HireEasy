# **HireEasy - Job Posting & Application Management Platform**

## **Overview**

HireEasy is a job posting and application management platform designed to streamline the hiring process. It enables HR teams to create job postings, manage candidate applications, and facilitate seamless recruitment workflows. The platform includes features like questionnaire customization, job sharing, application tracking, and comprehensive filtering options.

---

## **Features**

- **User Authentication:** Secure signup, login, and forgot password functionality.

- **Job Posting & Sharing:** Create and manage job postings with customizable questionnaires.

- **Application Management:** View, filter, and manage candidate applications.

- **Customizable Questionnaires:** Add various question types like radio buttons, checkboxes, etc.

- **Analytics & Insights:** Track job postings and candidate responses.

- **Integration with External Platforms:** Generate sharable job URLs for candidate applications.

- **Alerts & Notifications:** Get notified on new applications and job interactions.

---

## **Sprint 2 Report**

### **Completed Work in Sprint 2**

#### **Backend Team Achievements**

##### **Job Posting & Application Management**

- Implemented APIs to create, update, and delete job postings.

- Enabled filtering of job postings based on title, ID, and status.

- Developed candidate application storage and retrieval functionalities.

- Added filtering options for applications based on experience and other criteria.

#### **Frontend Team Achievements**

##### **Chiranjeevi (Structure and Functionality)**

1. **Updated Questionnaire Page**

   - Enabled adding, removing, and modifying questions.

   - Implemented different question types (radio buttons, checkboxes, dropdowns, etc.).

   - Added functionality to add/remove options dynamically.

2. **Updated Share Job Page**

   - Provided a preview feature for job postings before sharing.

   - Integrated support for different question types with live preview.

   - Enabled easy sharing through generated job URLs.

3. **Implemented View Job Applications Page**

   - Created an interface for HR to view candidate applications for a job.

   - Implemented filtering options to refine applications based on criteria such as experience.

   - Designed a detailed application view for individual candidates.

4. **Implemented View Jobs Page**

   - Developed a job listing page displaying all postings with their respective URLs.

   - Added search functionality to filter jobs based on partial job ID.

##### **Rahul Sai (UI Styling & Testing)**

1. **Styled View Jobs Page**

   - Created a structured, clutter-free design for job listings.

   - Implemented popup modals for individual job details.

   - Improved table display for job listings.

2. **Styled View Job Applications Page**

   - Implemented dropdowns for selecting Job IDs.

   - Designed a compact and structured display for candidate applications.

3. **Added Footer to All Pages**

   - Integrated a consistent footer with navigation, company logo, and contact details.

4. **Unit Testing & Integration Testing**

   - Set up unit testing using Vite + JavaScript.

   - Configured integration testing in Vite + React.

   - Partially integrated frontend and backend functionalities.

   - Set up Cypress for end-to-end testing.

---

## **Testing**

### **Frontend Unit Tests**

#### **Component Tests**

- **AuthForm.test.js**

  - Validates form submission and error handling.

  - Ensures correct rendering of input fields.

- **Questionnaire.test.js**

  - Tests adding/removing questions and options.

  - Validates different question types (radio, checkbox, dropdown, etc.).

#### **Page Tests**

- **Login.test.js**

  - Validates form submission and authentication flow.

- **Signup.test.js**

  - Tests user registration and validation.

- **Dashboard.test.js**

  - Ensures job postings and candidate applications render correctly.

- **QuestionnairePage.test.js**

  - Tests the ability to add, edit, and remove questions.

  - Validates different question types (radio, checkbox, text input, dropdown, etc.).

  

- **ShareJobPage.test.js**

  - Tests the preview functionality of a job posting.

  - Validates the URL generation and sharing features.

- **ViewJobApplicationsPage.test.js**

  - Ensures job applications are displayed correctly.

  - Tests filtering and sorting by experience, name, email, and job ID.

  

- **ViewJobsPage.test.js**

  - Validates the display of all job postings.

  - Tests filtering and quick copy functionality for job URLs.

  

- **Footer.test.js**

  - Ensures the footer is displayed on all pages.

  - Tests navigation links and logo visibility.

### **Cypress End-to-End Tests**

- **Login.spec.js**

  - Simulates a user logging in with valid/invalid credentials.

- **Signup.spec.js**

  - Tests the registration process.

- **Dashboard.spec.js**

  - Ensures job postings and applications display correctly.

- **Questionnaire.spec.js**

  - Verifies adding/removing questions and options dynamically.

- **JobApplicationFlow.spec.js**

  - Simulates an HR user reviewing candidate applications.

- **JobPosting.spec.js**

  - Simulates creating and sharing a job posting.

---

## **Pending Tasks and Reasons**

1. **Standalone Job Posting Page**

   - The current job posting uses a mock URL.

   - A standalone page needs to be developed to allow candidates to apply without requiring login (similar to Google Forms).

   - Pending due to time constraints but will be completed ASAP.

2. **Backend Integration Pending**

   - Some API integrations, such as CORS configurations, are awaiting backend updates.

   - Will be completed once the backend team provides the required changes.

3. **More Robust Testing**

   - Additional Cypress end-to-end tests will be added in the next sprint.

---

## **Setup & Installation**

### **Backend Setup**

1. Clone the repository:

   ```bash

   git clone https://github.com/your-repo/HireEasy.git

   cd HireEasy/backend

   ```

2. Install dependencies:

   ```bash

   go mod init github.com/HireEasy/backend

   go mod tidy

   ```

3. Start the server:

   ```bash

   cd backend

   go run cmd/server/main.go

   ```

### **Frontend Setup**

1. Navigate to the frontend directory:

   ```bash

   cd HireEasy/frontend

   ```

2. Install dependencies:

   ```bash

   npm install

   ```

3. Start the frontend application:

   ```bash

   npm start

   ```

### **Running Tests**

#### **Backend Tests**

```bash

cd backend

./test.sh

```

#### **Frontend Tests**

```bash

cd frontend

npm test

```

#### **Cypress End-to-End Tests**

```bash

cd frontend

npx cypress open

```

**<span style="text-decoration:underline;">User Stories assigned for Backend</span>**

**<span style="text-decoration:underline;">( Reshma - 52903493 )</span>**

1. As a Hiring Manager

I want to create form templates for job applications so that I can standardize the application process for different positions.

2. As a Hiring Manager

I want to view my created form templates so that I can manage and reuse them for different job positions.

3. As a Hiring Manager

I want to update my form templates so that I can modify the application requirements as needed.

4. As a Hiring Manager

I want to delete form templates so that I can remove outdated or unused templates.

5. As a Hiring Manager

I want to link form templates to job positions so that I can collect standardized applications for each position via the application form.

6. I want to manage the status of application forms so that I can control when applications are accepted and when they are frozen

7. As a Job Applicant

I want to view the application form.

8. As a Hiring Manager

I want to delete application forms so that I can remove outdated or completed application processes.

**<span style="text-decoration:underline;">Backend API documentation:</span>**

**<span style="text-decoration:underline;">Form Template APIs</span>**

1. Create a form template

POST /api/forms/templates

Request Body:

{

    "form_template_id": "string",

    "user_id": "integer",

    "fields": [

        {

            "field_name": "string",

            "field_type": "string",

            "required": "boolean",

            "label": "string"

        }

    ]

}

Response (201 Created):

{

    "form_template_id": "string",

    "user_id": "integer",

    "fields": [

        {

            "field_name": "string",

            "field_type": "string",

            "required": "boolean",

            "label": "string"

        }

    ]

}

2. Get form template

GET /api/forms/templates/:form_template_id

Response (200 OK):

{

    "form_template_id": "string",

    "user_id": "integer",

    "fields": [

        {

            "field_name": "string",

            "field_type": "string",

            "required": "boolean",

            "label": "string"

        }

    ]

}

3. List form Templates:

GET /api/forms/templates

Response (200 OK):

[

    {

        "form_template_id": "string",

        "user_id": "integer",

        "fields": [

            {

                "field_name": "string",

                "field_type": "string",

                "required": "boolean",

                "label": "string"

            }

        ]

    }

]

4. Delete form template:

DELETE /api/forms/templates/:form_template_id

Response (200 OK):

{

    "message": "Form template deleted successfully"

}

**<span style="text-decoration:underline;">Application Form APIs</span>**

5. Link Job to Form Template

POST /api/jobs/:job_id/forms

Request Body:

{

    "form_template_id": "string"

}

Response (201 Created):

{

    "form_uuid": "string"

}

6. Update form status  as active/inactive

PATCH /api/forms/:form_uuid/status

Request Body:

{

    "status": "string"

}

Response (200 OK):

{

    "form_uuid": "string",

    "status": "string"

}

7. Get Form Details

GET /api/forms/:form_uuid

Response (200 OK):

{

    "form_uuid": "string",

    "status": "string",

    "job": {

        "job_id": "string",

        "job_title": "string",

        "job_description": "string"

    },

    "form_template": {

        "form_template_id": "string",

        "fields": [

            {

                "field_name": "string",

                "field_type": "string",

                "required": "boolean",

                "label": "string"

            }

        ]

    }

}

8. Delete From 

DELETE /api/forms/:form_uuid

Response (200 OK):

{

    "message": "Form deleted successfully"

}

<span style="text-decoration:underline;">Unittests Backend:</span>


### **Authentication Tests (auth_test.go)**


<table>
  <tr>
   <td><strong>Test Case</strong>
   </td>
   <td><strong>Description</strong>
   </td>
   <td><strong>Expected Outcome</strong>
   </td>
  </tr>
  <tr>
   <td>TestRegisterUserH
   </td>
   <td>Tests user registration with valid credentials
   </td>
   <td>- HTTP Status 201 (Created)
<p>
- Response contains valid token
<p>
- User details (email, username) match input
   </td>
  </tr>
  <tr>
   <td>TestLoginUserH
   </td>
   <td>Tests user login with valid credentials
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Response contains valid token
   </td>
  </tr>
</table>



### **Job Management Tests (job_test.go)**


<table>
  <tr>
   <td><strong>Test Case</strong>
   </td>
   <td><strong>Description</strong>
   </td>
   <td><strong>Expected Outcome</strong>
   </td>
  </tr>
  <tr>
   <td>TestCreateJobH
   </td>
   <td>Tests job creation with valid data
   </td>
   <td>- HTTP Status 201 (Created)
<p>
- Job details match input
<p>
- Skills array contains all required skills
   </td>
  </tr>
  <tr>
   <td>TestUpdateJobH
   </td>
   <td>Tests job update with modified data
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Updated job details match input
<p>
- Skills array contains updated skills
   </td>
  </tr>
  <tr>
   <td>TestGetJobByIdH
   </td>
   <td>Tests retrieving a job by ID
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Job details match created job
<p>
- Skills array contains all skills
   </td>
  </tr>
  <tr>
   <td>TestGetJobsByTitleH
   </td>
   <td>Tests retrieving jobs by title
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Returns array of matching jobs
<p>
- Each job contains correct details
   </td>
  </tr>
  <tr>
   <td>TestGetJobsByStatusH
   </td>
   <td>Tests retrieving jobs by status
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Returns array of jobs with matching status
   </td>
  </tr>
  <tr>
   <td>TestListUserJobsH
   </td>
   <td>Tests listing all jobs for a user
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Returns array of user's jobs
   </td>
  </tr>
  <tr>
   <td>TestDeleteJobH
   </td>
   <td>Tests job deletion
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Job is successfully deleted
   </td>
  </tr>
</table>



### **Form Template Tests (form_template_test.go)**


<table>
  <tr>
   <td><strong>Test Case</strong>
   </td>
   <td><strong>Description</strong>
   </td>
   <td><strong>Expected Outcome</strong>
   </td>
  </tr>
  <tr>
   <td>TestCreateFormTemplateH
   </td>
   <td>Tests form template creation
   </td>
   <td>- HTTP Status 201 (Created)
<p>
- Template ID matches input
<p>
- Fields array contains all fields
   </td>
  </tr>
  <tr>
   <td>TestGetFormTemplateH
   </td>
   <td>Tests retrieving a form template
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Template details match created template
<p>
- Fields array contains all fields
   </td>
  </tr>
  <tr>
   <td>TestListFormTemplatesH
   </td>
   <td>Tests listing all form templates
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Returns array of templates
<p>
- Each template contains correct details
   </td>
  </tr>
  <tr>
   <td>TestDeleteFormTemplateH
   </td>
   <td>Tests form template deletion
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Template is successfully deleted
<p>
- GET request returns 404
   </td>
  </tr>
</table>



### **Application Form Tests (application_form_test.go)**


<table>
  <tr>
   <td><strong>Test Case</strong>
   </td>
   <td><strong>Description</strong>
   </td>
   <td><strong>Expected Outcome</strong>
   </td>
  </tr>
  <tr>
   <td>TestLinkJobToFormTemplateH
   </td>
   <td>Tests linking a job to a form template
   </td>
   <td>- HTTP Status 201 (Created)
<p>
- Response contains form UUID
   </td>
  </tr>
  <tr>
   <td>TestUpdateFormStatusH
   </td>
   <td>Tests updating form status
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Status is successfully updated
<p>
- Response contains updated status
   </td>
  </tr>
  <tr>
   <td>TestGetFormDetailsH
   </td>
   <td>Tests retrieving form details
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Response contains form UUID
<p>
- Includes job and form template details
   </td>
  </tr>
  <tr>
   <td>TestDeleteFormH
   </td>
   <td>Tests form deletion
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Form is successfully deleted
<p>
- GET request returns 404
   </td>
  </tr>
</table>


**<span style="text-decoration:underline;">User Stories assigned for Backend</span>**

**<span style="text-decoration:underline;">( Sushmitha - &lt;56737143> )</span>**

1. As a Job Applicant

I want to submit my application through an online form so that I can apply for job positions.

2, As a Hiring Manager

I want to review job applications efficiently so that I can identify the best candidates quickly. I want to view all submissions for a specific job position.

3. As a Hiring Manager

I want to have the option to view the application ranked by ATS score 

4. As a Hiring Manager

I want to have the option to filter the candidates on particular fields

**<span style="text-decoration:underline;">Backend API listing:</span>**

**- Form Submission Tests**

Successful Submissions – Tests if valid applications are accepted.

New User Submission – Ensures new users can apply and are created if needed.

Missing Job ID – Checks if submissions fail when the job ID is missing.

Invalid User ID – Ensures incorrect user IDs trigger an error.

Duplicate Submission – Tests if users are prevented from applying to the same job twice.

Invalid Form Data – Ensures submissions with bad JSON data fail.

Resume Upload Handling – Verifies resume file uploads work correctly.

Failed Resume Upload – Simulates a failed file upload and checks error handling.

**- Form Retrieval Tests**

Fetch All Submissions – Ensures submissions are retrieved correctly.

Invalid Form UUID – Checks that bad UUIDs return an error.

Filter by Date – Ensures submissions can be filtered by a specific date.

Sorting Options – Tests if submissions are sorted correctly by ATS score or creation date.

Result Limiting – Ensures the API returns the correct number of submissions when a limit is set.

Invalid Limit Handling – Checks if negative or invalid limits default to a proper value.

**- ATS Score Calculation Tests**

No Skills – Ensures a minimum base score is given even when no skills are submitted.

Single Skill – Tests if adding one skill increases the score correctly.

Multiple Skills – Ensures multiple skills increase the score up to the maximum limit.

Maximum Score – Checks that scores do not exceed 100.

Empty Resume URL – Ensures the ATS score still calculates correctly if no resume is uploaded.

Duplicate Skills – Verifies that duplicate skills are not counted multiple times.

Long Skill List – Tests if a long list of skills properly maxes out the score.

**<span style="text-decoration:underline;">Unittests Listing:</span>**


<table>
  <tr>
   <td><strong>Test Case</strong>
   </td>
   <td><strong>Description</strong>
   </td>
   <td><strong>Expected Outcome</strong>
   </td>
  </tr>
  <tr>
   <td>TestSubmitFormResponseH
   </td>
   <td>Tests form submission with valid responses
   </td>

   <td>- HTTP Status 201 (Created)
<p>
- Submission is stored with the correct form UUID
<p>
- Submission contains form responses and ATS score
<p>
- Correct score returned from external ATS API
   </td>
  <tr>
   <td>TestSubmitFormInvalidDataH
   </td>
   <td>Tests form submission with invalid data
   </td>
   <td>- HTTP Status 400 (Bad Request)
<p>
- Response contains error message specifying invalid form data
   </td>
  </tr>


  <tr>
   <td>TestSubmitFormATSFailureH
   </td>
   <td>Tests form submission when ATS API fails
   </td>
   <td>- HTTP Status 500 (Internal Server Error)
<p>
- Response contains error message indicating ATS API failure
   </td>
  </tr>


  <tr>
   <td>TestSubmitFormMissingFieldsH
   </td>
   <td>Tests form submission with missing required fields
   </td>
   <td>- HTTP Status 400 (Bad Request)
<p>
- Response contains error message specifying missing fields
   </td>
  </tr>

  <tr>
   <td>TestGetFormSubmissionsH
   </td>
   <td>Tests retrieving form submissions sorted by reception order
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Returns first <code>n</code> submissions (default 20) sorted by order of reception
<p>
- Submissions contain correct responses and scores
   </td>
  </tr>


  <tr>
   <td>TestGetFormSubmissionsByScoreH
   </td>
   <td>Tests retrieving form submissions sorted by ATS score
   </td>
   <td>- HTTP Status 200 (OK)
<p>
- Returns first <code>n</code> submissions (default 20) sorted by ATS score
<p>
- Submissions are correctly ordered according to score
   </td>
  </tr>
</table>