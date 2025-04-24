// Custom command for login
Cypress.Commands.add("login", (email, password) => {
  // Mock login API response
  cy.intercept("POST", "/api/auth/login", {
    statusCode: 200,
    body: {
      token: "fake-token",
      user: {
        email,
        role: email.includes("hr") ? "HR" : "Interviewer",
      },
    },
  }).as("loginRequest");

  cy.visit("/");
  cy.get("[data-cy=email]").should("be.visible").type(email);
  cy.get("[data-cy=password]").should("be.visible").type(password);
  cy.get("[data-cy=login-button]").should("be.visible").click();
  cy.wait("@loginRequest");
  cy.url().should("include", "/dashboard");
});

// Custom command for logout
Cypress.Commands.add("logout", () => {
  cy.intercept("POST", "/api/auth/logout", {
    statusCode: 200,
    body: { success: true },
  }).as("logoutRequest");

  cy.get("[data-cy=logout-button]").should("be.visible").click();
  cy.wait("@logoutRequest");
  cy.url().should("include", "/login");
});

// Custom command for clearing local storage
Cypress.Commands.overwrite("clearLocalStorage", () => {
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
    failOnStatusCode: false,
  });
});

// Custom command for handling file uploads
Cypress.Commands.add("uploadFile", (selector, filePath) => {
  cy.intercept("POST", "/api/upload", {
    statusCode: 200,
    body: { url: "https://example.com/uploaded-file.jpg" },
  }).as("uploadRequest");

  cy.get(selector).should("be.visible").attachFile(filePath);
  cy.wait("@uploadRequest");
});

// Custom command for handling network requests
Cypress.Commands.add("mockApiResponse", (method, url, response) => {
  const alias = `apiRequest-${method}-${url}`;
  cy.intercept(method, url, response).as(alias);
  return cy.wait(`@${alias}`);
});

// Custom command for handling dates
Cypress.Commands.add("getFutureDate", (daysAhead = 1) => {
  const date = new Date();
  date.setDate(date.getDate() + daysAhead);
  return date.toISOString().split('T')[0];
});

// Custom command for handling time slots
Cypress.Commands.add("getTimeSlot", (hours = 14, minutes = 0) => {
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
}); 