describe('Interviewer Profile Tests', () => {
  beforeEach(() => {
    // Clear any existing state
    cy.clearCookies();
    cy.clearLocalStorage();
    
    // Mock the profile API response
    cy.intercept('GET', '**/api/profiles/me', {
      statusCode: 200,
      body: {
        job_title: 'Senior Software Engineer',
        years_of_experience: 5,
        areas_of_expertise: ['JavaScript', 'React', 'Node.js'],
        phone_number: '1234567890'
      }
    }).as('getProfile');

    // Mock the profile update API response
    cy.intercept('PUT', '**/api/profiles', {
      statusCode: 200,
      body: {
        job_title: 'Senior Software Engineer',
        years_of_experience: 5,
        areas_of_expertise: ['JavaScript', 'React', 'Node.js'],
        phone_number: '1234567890'
      }
    }).as('updateProfile');

    // Login as interviewer
    cy.login('interviewer@example.com', 'password');
    
    // Visit the profile page
    cy.visit('/interviewer-profile');
  });

  it('should load the profile page successfully', () => {
    cy.get('h2').should('contain', 'Interviewer Profile');
    cy.wait('@getProfile');
  });

  it('should display existing profile data', () => {
    cy.wait('@getProfile');
    cy.get('[data-cy="job-title-input"]').should('have.value', 'Senior Software Engineer');
    cy.get('[data-cy="years-of-experience-input"]').should('have.value', '5');
    cy.get('[data-cy="areas-of-expertise-input"]').should('have.value', 'JavaScript, React, Node.js');
    cy.get('[data-cy="phone-number-input"]').should('have.value', '1234567890');
  });

  it('should update profile successfully', () => {
    cy.wait('@getProfile');
    
    // Update form fields
    cy.get('[data-cy="job-title-input"]').clear().type('Lead Developer');
    cy.get('[data-cy="years-of-experience-input"]').clear().type('7');
    cy.get('[data-cy="areas-of-expertise-input"]').clear().type('JavaScript, React, Node.js, TypeScript');
    cy.get('[data-cy="phone-number-input"]').clear().type('9876543210');

    // Submit form
    cy.get('[data-cy="save-profile-button"]').click();
    cy.wait('@updateProfile');

    // Check success message
    cy.get('[data-cy="success-message"]').should('be.visible')
      .and('contain', 'Profile saved successfully');
  });

  it('should show error message on failed update', () => {
    // Mock failed API response
    cy.intercept('PUT', '**/api/profiles', {
      statusCode: 500,
      body: { message: 'Server error' }
    }).as('failedUpdate');

    cy.wait('@getProfile');
    
    // Update a field
    cy.get('[data-cy="job-title-input"]').clear().type('Lead Developer');
    
    // Submit form
    cy.get('[data-cy="save-profile-button"]').click();
    cy.wait('@failedUpdate');

    // Check error message
    cy.get('[data-cy="error-message"]').should('be.visible')
      .and('contain', 'Something went wrong');
  });

  it('should navigate to dashboard on cancel', () => {
    cy.get('[data-cy="cancel-button"]').click();
    cy.url().should('include', '/interviewer-dashboard');
  });

  it('should validate required fields', () => {
    cy.wait('@getProfile');
    
    // Clear required fields
    cy.get('[data-cy="job-title-input"]').clear();
    cy.get('[data-cy="years-of-experience-input"]').clear();
    cy.get('[data-cy="areas-of-expertise-input"]').clear();

    // Try to submit
    cy.get('[data-cy="save-profile-button"]').click();

    // Check that form submission was prevented
    cy.get('[data-cy="success-message"]').should('not.exist');
    cy.get('[data-cy="error-message"]').should('not.exist');
  });
}); 