describe("Interview Scheduling", () => {
  beforeEach(() => {
    cy.clearLocalStorage();
    cy.fixture("users").then((users) => {
      cy.login(users.hrUser.email, users.hrUser.password);
    });
  });

  it("should load schedule interviews page", () => {
    cy.visit("/schedule-interviews");
    cy.get("[data-cy=schedule-title]").should("contain", "Schedule Interviews");
  });

  it("should schedule a new interview", () => {
    cy.visit("/schedule-interviews");

    // Select job and candidate
    cy.get("[data-cy=select-job]").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").select("John Doe");

    // Select interviewer and time slot
    cy.get("[data-cy=select-interviewer]").select("Jane Smith");
    cy.get("[data-cy=select-date]").type("2024-04-25");
    cy.get("[data-cy=select-time]").type("14:00");

    // Schedule interview
    cy.get("[data-cy=schedule-button]").click();

    // Verify success
    cy.get("[data-cy=success-message]").should("be.visible");
    cy.get("[data-cy=scheduled-interview]").should("contain", "John Doe");
  });

  it("should show available time slots", () => {
    cy.visit("/schedule-interviews");

    // Select job and candidate
    cy.get("[data-cy=select-job]").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").select("John Doe");

    // Select interviewer
    cy.get("[data-cy=select-interviewer]").select("Jane Smith");

    // Verify available slots are displayed
    cy.get("[data-cy=available-slots]").should("exist");
  });

  it("should prevent double booking", () => {
    cy.visit("/schedule-interviews");

    // Schedule first interview
    cy.get("[data-cy=select-job]").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").select("John Doe");
    cy.get("[data-cy=select-interviewer]").select("Jane Smith");
    cy.get("[data-cy=select-date]").type("2024-04-25");
    cy.get("[data-cy=select-time]").type("14:00");
    cy.get("[data-cy=schedule-button]").click();

    // Try to schedule another interview in same slot
    cy.get("[data-cy=select-candidate]").select("Alice Johnson");
    cy.get("[data-cy=select-date]").type("2024-04-25");
    cy.get("[data-cy=select-time]").type("14:00");
    cy.get("[data-cy=schedule-button]").click();

    // Verify error message
    cy.get("[data-cy=error-message]").should("contain", "Time slot already booked");
  });

  it("should cancel scheduled interview", () => {
    cy.visit("/schedule-interviews");

    // Schedule an interview first
    cy.get("[data-cy=select-job]").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").select("John Doe");
    cy.get("[data-cy=select-interviewer]").select("Jane Smith");
    cy.get("[data-cy=select-date]").type("2024-04-25");
    cy.get("[data-cy=select-time]").type("14:00");
    cy.get("[data-cy=schedule-button]").click();

    // Cancel the interview
    cy.get("[data-cy=cancel-interview]").first().click();
    cy.get("[data-cy=confirm-cancel]").click();

    // Verify cancellation
    cy.get("[data-cy=success-message]").should("contain", "Interview cancelled");
  });
});
