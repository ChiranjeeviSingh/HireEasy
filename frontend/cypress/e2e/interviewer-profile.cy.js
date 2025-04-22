describe("Interviewer Profile Management", () => {
  beforeEach(() => {
    cy.clearLocalStorage();
    cy.fixture("users").then((users) => {
      cy.login(users.interviewerUser.email, users.interviewerUser.password);
    });
  });

  it("should load interviewer profile page", () => {
    cy.visit("/interviewer-profile");
    cy.get("[data-cy=profile-title]").should("contain", "Interviewer Profile");
  });

  it("should update profile information", () => {
    cy.visit("/interviewer-profile");

    // Mock API response
    cy.mockApiResponse("PUT", "/api/interviewer/profile", {
      statusCode: 200,
      body: { success: true },
    });

    // Fill in profile information
    cy.get("[data-cy=job-title]").clear().type("Senior Software Engineer");
    cy.get("[data-cy=experience]").clear().type("5");
    cy.get("[data-cy=expertise]").clear().type("JavaScript, React, Node.js");
    cy.get("[data-cy=company]").clear().type("Tech Corp");
    cy.get("[data-cy=linkedin]").clear().type("https://linkedin.com/in/interviewer");

    // Save profile
    cy.get("[data-cy=save-profile]").click();

    // Verify success message
    cy.get("[data-cy=success-message]").should("be.visible");

    // Verify updated information
    cy.get("[data-cy=job-title]").should("have.value", "Senior Software Engineer");
    cy.get("[data-cy=experience]").should("have.value", "5");
  });

  it("should handle validation errors", () => {
    cy.visit("/interviewer-profile");

    // Try to save with empty required fields
    cy.get("[data-cy=save-profile]").click();

    // Verify error messages
    cy.get("[data-cy=error-message]").should("be.visible");
    cy.get("[data-cy=job-title-error]").should("contain", "Job title is required");
  });

  it("should display current availability status", () => {
    cy.visit("/interviewer-profile");
    cy.get("[data-cy=availability-status]").should("exist");
  });

  it("should handle API errors gracefully", () => {
    cy.visit("/interviewer-profile");

    // Mock API error response
    cy.mockApiResponse("PUT", "/api/interviewer/profile", {
      statusCode: 500,
      body: { error: "Internal Server Error" },
    });

    // Try to save profile
    cy.get("[data-cy=job-title]").clear().type("Senior Developer");
    cy.get("[data-cy=save-profile]").click();

    // Verify error message
    cy.get("[data-cy=error-message]").should("contain", "Failed to update profile");
  });

  it("should handle file uploads", () => {
    cy.visit("/interviewer-profile");

    // Upload profile picture
    cy.uploadFile("[data-cy=profile-picture]", "profile.jpg");

    // Verify upload success
    cy.get("[data-cy=upload-success]").should("be.visible");
  });
});
