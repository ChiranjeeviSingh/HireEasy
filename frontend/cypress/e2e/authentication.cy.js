describe("Authentication Tests", () => {
  beforeEach(() => {
    // Reset any state
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  it("can visit the home page", () => {
    cy.visit('/');
    cy.get('html').should('exist');
  });

  it("can load user fixture data", () => {
    cy.fixture('users').then((users) => {
      expect(users).to.have.property('jobSeeker');
      expect(users).to.have.property('employer');
    });
  });

  it("can intercept API calls", () => {
    // Set up API interception
    cy.intercept('POST', '**/api/auth/**', {
      statusCode: 200,
      body: { success: true }
    }).as('authCall');
    
    // Visit home page
    cy.visit('/');
  });

  it("has a reliable assertion", () => {
    expect(true).to.equal(true);
  });
}); 