describe("Interviewer Interviews Management", () => {
  beforeEach(() => {
    cy.clearLocalStorage();
    cy.fixture("users").then((users) => {
      cy.login(users.interviewerUser.email, users.interviewerUser.password);
    });
  });

  it("should load interviewer interviews page", () => {
    cy.visit("/interviewer-interviews");
    cy.get("[data-cy=interviews-title]").should("contain", "My Interviews");
  });

  it("should display upcoming interviews", () => {
    cy.visit("/interviewer-interviews");
    cy.get("[data-cy=upcoming-interviews]").should("exist");
  });

  it("should submit interview feedback", () => {
    cy.visit("/interviewer-interviews");

    // Select a completed interview
    cy.get("[data-cy=completed-interview]").first().click();

    // Fill in feedback
    cy.get("[data-cy=feedback-text]").type("Candidate demonstrated strong technical skills and good communication.");
    cy.get("[data-cy=technical-score]").select("4");
    cy.get("[data-cy=communication-score]").select("5");
    cy.get("[data-cy=verdict]").select("recommended");

    // Submit feedback
    cy.get("[data-cy=submit-feedback]").click();

    // Verify success
    cy.get("[data-cy=success-message]").should("be.visible");
  });

  it("should filter interviews by status", () => {
    cy.visit("/interviewer-interviews");

    // Filter by completed interviews
    cy.get("[data-cy=status-filter]").select("completed");

    // Verify filtered results
    cy.get("[data-cy=interview-status]").each(($el) => {
      cy.wrap($el).should("contain", "Completed");
    });
  });

  it("should view candidate details before interview", () => {
    cy.visit("/interviewer-interviews");

    // Select an upcoming interview
    cy.get("[data-cy=upcoming-interview]").first().click();

    // View candidate details
    cy.get("[data-cy=view-candidate]").click();

    // Verify candidate information
    cy.get("[data-cy=candidate-profile]").should("exist");
    cy.get("[data-cy=candidate-resume]").should("exist");
  });

  it("should handle interview rescheduling", () => {
    cy.visit("/interviewer-interviews");

    // Select an upcoming interview
    cy.get("[data-cy=upcoming-interview]").first().click();

    // Request reschedule
    cy.get("[data-cy=reschedule-button]").click();
    cy.get("[data-cy=reschedule-reason]").type("Unexpected conflict");
    cy.get("[data-cy=submit-reschedule]").click();

    // Verify request
    cy.get("[data-cy=success-message]").should("contain", "Reschedule request submitted");
  });
});
