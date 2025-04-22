describe("Interviewer Calendar Management", () => {
  beforeEach(() => {
    cy.clearLocalStorage();
    cy.fixture("users").then((users) => {
      cy.login(users.interviewerUser.email, users.interviewerUser.password);
    });
  });

  it("should load calendar page", () => {
    cy.visit("/interviewer-calendar");
    cy.get("[data-cy=calendar-title]").should("contain", "Interview Calendar");
  });

  it("should add availability slot", () => {
    cy.visit("/interviewer-calendar");

    // Mock API response
    cy.mockApiResponse("POST", "/api/availability", {
      statusCode: 200,
      body: { success: true },
    });

    // Open add slot modal
    cy.get("[data-cy=add-slot]").click();

    // Fill in slot details
    cy.get("[data-cy=date-input]").type("2024-04-25");
    cy.get("[data-cy=start-time]").type("13:00");
    cy.get("[data-cy=end-time]").type("14:00");

    // Save slot
    cy.get("[data-cy=save-slot]").click();

    // Verify slot was added
    cy.get("[data-cy=slot-list]").should("contain", "13:00 - 14:00");
  });

  it("should delete availability slot", () => {
    cy.visit("/interviewer-calendar");

    // Mock API response
    cy.mockApiResponse("DELETE", "/api/availability/*", {
      statusCode: 200,
      body: { success: true },
    });

    // Add a slot first
    cy.get("[data-cy=add-slot]").click();
    cy.get("[data-cy=date-input]").type("2024-04-25");
    cy.get("[data-cy=start-time]").type("13:00");
    cy.get("[data-cy=end-time]").type("14:00");
    cy.get("[data-cy=save-slot]").click();

    // Delete the slot
    cy.get("[data-cy=delete-slot]").first().click();
    cy.get("[data-cy=confirm-delete]").click();

    // Verify slot was removed
    cy.get("[data-cy=slot-list]").should("not.contain", "13:00 - 14:00");
  });

  it("should prevent overlapping slots", () => {
    cy.visit("/interviewer-calendar");

    // Add first slot
    cy.get("[data-cy=add-slot]").click();
    cy.get("[data-cy=date-input]").type("2024-04-25");
    cy.get("[data-cy=start-time]").type("13:00");
    cy.get("[data-cy=end-time]").type("14:00");
    cy.get("[data-cy=save-slot]").click();

    // Try to add overlapping slot
    cy.get("[data-cy=add-slot]").click();
    cy.get("[data-cy=date-input]").type("2024-04-25");
    cy.get("[data-cy=start-time]").type("13:30");
    cy.get("[data-cy=end-time]").type("14:30");
    cy.get("[data-cy=save-slot]").click();

    // Verify error message
    cy.get("[data-cy=error-message]").should("contain", "Time slot overlaps with existing slot");
  });

  it("should display scheduled interviews", () => {
    cy.visit("/interviewer-calendar");

    // Mock API response for scheduled interviews
    cy.mockApiResponse("GET", "/api/interviews", {
      statusCode: 200,
      body: {
        interviews: [
          {
            id: 1,
            candidate: "John Doe",
            date: "2024-04-25",
            startTime: "14:00",
            endTime: "15:00",
          },
        ],
      },
    });

    cy.get("[data-cy=scheduled-interviews]").should("exist");
    cy.get("[data-cy=interview-item]").should("contain", "John Doe");
  });

  it("should handle API errors gracefully", () => {
    cy.visit("/interviewer-calendar");

    // Mock API error response
    cy.mockApiResponse("POST", "/api/availability", {
      statusCode: 500,
      body: { error: "Internal Server Error" },
    });

    // Try to add slot
    cy.get("[data-cy=add-slot]").click();
    cy.get("[data-cy=date-input]").type("2024-04-25");
    cy.get("[data-cy=start-time]").type("13:00");
    cy.get("[data-cy=end-time]").type("14:00");
    cy.get("[data-cy=save-slot]").click();

    // Verify error message
    cy.get("[data-cy=error-message]").should("contain", "Failed to add time slot");
  });
});
