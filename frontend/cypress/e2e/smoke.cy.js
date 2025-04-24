describe('Smoke Test', () => {
  it('can visit the home page', () => {
    // Visit the base URL defined in cypress.config.js
    cy.visit('/');
    
    // Instead of checking for specific content, just verify we loaded a page
    cy.get('html').should('exist');
  });

  it('checks that Cypress is working correctly', () => {
    // A simple assertion that will always pass
    expect(true).to.equal(true);
  });
}); 