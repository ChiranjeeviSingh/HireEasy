describe("Authentication Flow", () => {
  beforeEach(() => {
    cy.clearLocalStorage();
  });

  it("should login successfully as HR", () => {
    cy.fixture("users").then((users) => {
      cy.login(users.hrUser.email, users.hrUser.password);
      cy.url().should("include", "/dashboard");
      cy.get("[data-cy=user-role]").should("contain", "HR");
    });
  });

  it("should login successfully as Interviewer", () => {
    cy.fixture("users").then((users) => {
      cy.login(users.interviewerUser.email, users.interviewerUser.password);
      cy.url().should("include", "/dashboard");
      cy.get("[data-cy=user-role]").should("contain", "Interviewer");
    });
  });

  it("should handle invalid credentials", () => {
    cy.login("invalid@example.com", "wrongpassword");
    cy.get("[data-cy=error-message]").should("be.visible");
  });

  it("should logout successfully", () => {
    cy.fixture("users").then((users) => {
      cy.login(users.hrUser.email, users.hrUser.password);
      cy.logout();
      cy.url().should("include", "/login");
    });
  });
});
