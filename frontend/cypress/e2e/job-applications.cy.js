describe("Job Applications Management", () => {
  beforeEach(() => {
    cy.clearLocalStorage();
    cy.fixture("users").then((users) => {
      cy.login(users.hrUser.email, users.hrUser.password);
    });
  });

  it("should load job applications page", () => {
    cy.visit("/job-applications");
    cy.get("[data-cy=applications-title]").should("contain", "Job Applications");
  });

  it("should filter applications by status", () => {
    cy.visit("/job-applications");

    // Select status filter
    cy.get("[data-cy=status-filter]").select("under_review");

    // Verify filtered results
    cy.get("[data-cy=application-list]").should("exist");
    cy.get("[data-cy=application-status]").each(($el) => {
      cy.wrap($el).should("contain", "Under Review");
    });
  });

  it("should update application status", () => {
    cy.visit("/job-applications");

    // Select an application
    cy.get("[data-cy=application-item]").first().click();

    // Update status
    cy.get("[data-cy=update-status]").select("accepted");
    cy.get("[data-cy=save-status]").click();

    // Verify status update
    cy.get("[data-cy=success-message]").should("be.visible");
    cy.get("[data-cy=current-status]").should("contain", "Accepted");
  });

  it("should sort applications by ATS score", () => {
    cy.visit("/job-applications");

    // Click sort by ATS score
    cy.get("[data-cy=sort-ats]").click();

    // Verify sorting
    cy.get("[data-cy=ats-score]").then(($scores) => {
      const scores = $scores.map((_, el) => parseFloat(el.textContent)).get();
      const sortedScores = [...scores].sort((a, b) => b - a);
      expect(scores).to.deep.equal(sortedScores);
    });
  });

  it("should view candidate details", () => {
    cy.visit("/job-applications");

    // Click on a candidate
    cy.get("[data-cy=candidate-name]").first().click();

    // Verify candidate details
    cy.get("[data-cy=candidate-profile]").should("exist");
    cy.get("[data-cy=candidate-resume]").should("exist");
  });
});
