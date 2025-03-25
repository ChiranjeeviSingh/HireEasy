describe('HireEasy Login', () => {
  beforeEach(() => {
    // Visit the login page before each test
    cy.visit('http://localhost:5173/');
  });

  it('displays the login form', () => {
    // Check that the login page elements are visible
    cy.contains('h1', 'HireEasy').should('be.visible');
    cy.contains('Sign in and start hiring the best talent out there.').should('be.visible');
    cy.get('input[placeholder="Enter your email"]').should('be.visible');
    cy.get('input[placeholder="Enter your password"]').should('be.visible');
    cy.contains('button', 'Sign In').should('be.visible');
  });

  it('shows validation errors for empty fields', () => {
    // Try to submit the form without entering any data
    cy.contains('button', 'Sign In').click();
    
    // Check for HTML5 validation on required fields
    cy.get('input[placeholder="Enter your email"]')
      .invoke('prop', 'validity')
      .should('have.property', 'valueMissing', true);
  });

  it('shows error message for invalid credentials', () => {
    // Enter invalid credentials
    cy.get('input[placeholder="Enter your email"]').type('wrong@email.com');
    cy.get('input[placeholder="Enter your password"]').type('wrongpassword');
    
    // Submit the form
    cy.contains('button', 'Sign In').click();
    
    // Check for the exact error message from the Login component
    cy.contains('Invalid credentials, please try again.').should('be.visible');
  });

  it('successfully logs in with valid credentials', () => {
    // Enter the mock credentials from the Login component
    cy.get('input[placeholder="Enter your email"]').type('abcd@gmail.com');
    cy.get('input[placeholder="Enter your password"]').type('abcdef567');
    
    // Submit the form
    cy.contains('button', 'Sign In').click();
    
    // Check for the exact alert message from the Login component
    cy.on('window:alert', (text) => {
      expect(text).to.equal('Login Successful!');
    });
    
    // After successful login, should navigate to dashboard
    cy.url().should('include', '/dashboard');
  });

  it('allows navigation to register page', () => {
    // Click on the "Create One Now" link
    cy.contains('a', 'Create One Now').click();
    
    // Verify navigation to register page
    cy.url().should('include', '/register');
    // Check for sign up form elements
    cy.contains('Sign Up').should('be.visible');
  });
});

describe('Dashboard Links Visibility', () => {
  beforeEach(() => {
    // Login first with the correct mock credentials
    cy.visit('http://localhost:5173/');
    cy.get('input[placeholder="Enter your email"]').type('abcd@gmail.com');
    cy.get('input[placeholder="Enter your password"]').type('abcdef567');
    cy.contains('button', 'Sign In').click();
    
    // Ensure we're on the dashboard page
    cy.url().should('include', '/dashboard');
  });

  it('displays all five navigation links', () => {
    // Check that all five links are visible
    cy.contains('Create Job Posting').should('be.visible');
    cy.contains('Create Questionnaire').should('be.visible');
    cy.contains('Share Job').should('be.visible');
    cy.contains('View Job Applications').should('be.visible');
    cy.contains('View Jobs').should('be.visible');
  });

  it('verifies all links are clickable', () => {
    // Check that each link is not just visible but also clickable (has proper HTML element)
    cy.contains('Create Job Posting').should('have.prop', 'tagName').should('eq', 'BUTTON');
    cy.contains('Create Questionnaire').should('have.prop', 'tagName').should('eq', 'BUTTON');
    cy.contains('Share Job').should('have.prop', 'tagName').should('eq', 'BUTTON');
    cy.contains('View Job Applications').should('have.prop', 'tagName').should('eq', 'BUTTON');
    cy.contains('View Jobs').should('have.prop', 'tagName').should('eq', 'BUTTON');
  });


  // Optional: Test each link navigation if needed
  
  // Add similar tests for other links if desired
});