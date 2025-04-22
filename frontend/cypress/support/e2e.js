import "./commands";

// Custom command for login
Cypress.Commands.add("login", (email, password) => {
  cy.visit("/login");
  cy.get("[data-cy=email]").type(email);
  cy.get("[data-cy=password]").type(password);
  cy.get("[data-cy=login-button]").click();
  cy.url().should("include", "/dashboard");
});

// Custom command for logout
Cypress.Commands.add("logout", () => {
  cy.get("[data-cy=logout-button]").click();
  cy.url().should("include", "/login");
});

// Custom command for clearing local storage
Cypress.Commands.add("clearLocalStorage", () => {
  cy.window().then((win) => {
    win.localStorage.clear();
  });
});

// Custom command for API requests
Cypress.Commands.add("apiRequest", (method, url, body) => {
  return cy.request({
    method,
    url,
    body,
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
  });
});

// Custom command for handling file uploads
Cypress.Commands.add("uploadFile", (selector, filePath) => {
  cy.get(selector).attachFile(filePath);
});

// Custom command for handling network requests
Cypress.Commands.add("mockApiResponse", (method, url, response) => {
  cy.intercept(method, url, response);
});

// Global error handling
Cypress.on("uncaught:exception", (err, runnable) => {
  // returning false here prevents Cypress from failing the test
  return false;
});
