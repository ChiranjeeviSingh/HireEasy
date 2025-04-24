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

    // Mock API responses
    cy.mockApiResponse("GET", "/api/jobs", {
      statusCode: 200,
      body: [{ id: 1, title: "Software Engineer" }],
    });

    cy.mockApiResponse("GET", "/api/candidates", {
      statusCode: 200,
      body: [{ id: 1, name: "John Doe" }],
    });

    cy.mockApiResponse("GET", "/api/interviewers", {
      statusCode: 200,
      body: [{ id: 1, name: "Jane Smith" }],
    });

    cy.mockApiResponse("POST", "/api/interviews", {
      statusCode: 200,
      body: { success: true },
    });

    // Select job and candidate
    cy.get("[data-cy=select-job]").should("be.visible").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").should("be.visible").select("John Doe");

    // Select interviewer and time slot
    cy.get("[data-cy=select-interviewer]").should("be.visible").select("Jane Smith");
    cy.get("[data-cy=select-date]").should("be.visible").type(cy.getFutureDate());
    cy.get("[data-cy=select-time]").should("be.visible").type(cy.getTimeSlot());

    // Schedule interview
    cy.get("[data-cy=schedule-button]").should("be.visible").click();

    // Verify success
    cy.get("[data-cy=success-message]").should("be.visible");
    cy.get("[data-cy=scheduled-interview]").should("contain", "John Doe");
  });

  it("should show available time slots", () => {
    cy.visit("/schedule-interviews");

    // Mock API responses
    cy.mockApiResponse("GET", "/api/jobs", {
      statusCode: 200,
      body: [{ id: 1, title: "Software Engineer" }],
    });

    cy.mockApiResponse("GET", "/api/candidates", {
      statusCode: 200,
      body: [{ id: 1, name: "John Doe" }],
    });

    cy.mockApiResponse("GET", "/api/interviewers", {
      statusCode: 200,
      body: [{ id: 1, name: "Jane Smith" }],
    });

    cy.mockApiResponse("GET", "/api/availability/*", {
      statusCode: 200,
      body: {
        slots: [
          { date: cy.getFutureDate(), time: "14:00" },
          { date: cy.getFutureDate(), time: "15:00" },
        ],
      },
    });

    // Select job and candidate
    cy.get("[data-cy=select-job]").should("be.visible").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").should("be.visible").select("John Doe");

    // Select interviewer
    cy.get("[data-cy=select-interviewer]").should("be.visible").select("Jane Smith");

    // Verify available slots are displayed
    cy.get("[data-cy=available-slots]").should("exist");
  });

  it("should prevent double booking", () => {
    cy.visit("/schedule-interviews");

    // Mock API responses
    cy.mockApiResponse("GET", "/api/jobs", {
      statusCode: 200,
      body: [{ id: 1, title: "Software Engineer" }],
    });

    cy.mockApiResponse("GET", "/api/candidates", {
      statusCode: 200,
      body: [
        { id: 1, name: "John Doe" },
        { id: 2, name: "Alice Johnson" },
      ],
    });

    cy.mockApiResponse("GET", "/api/interviewers", {
      statusCode: 200,
      body: [{ id: 1, name: "Jane Smith" }],
    });

    const futureDate = cy.getFutureDate();
    const timeSlot = cy.getTimeSlot();

    // Schedule first interview
    cy.get("[data-cy=select-job]").should("be.visible").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").should("be.visible").select("John Doe");
    cy.get("[data-cy=select-interviewer]").should("be.visible").select("Jane Smith");
    cy.get("[data-cy=select-date]").should("be.visible").type(futureDate);
    cy.get("[data-cy=select-time]").should("be.visible").type(timeSlot);
    cy.get("[data-cy=schedule-button]").should("be.visible").click();

    // Try to schedule another interview in same slot
    cy.mockApiResponse("POST", "/api/interviews", {
      statusCode: 409,
      body: { error: "Time slot already booked" },
    });

    cy.get("[data-cy=select-candidate]").should("be.visible").select("Alice Johnson");
    cy.get("[data-cy=select-date]").should("be.visible").type(futureDate);
    cy.get("[data-cy=select-time]").should("be.visible").type(timeSlot);
    cy.get("[data-cy=schedule-button]").should("be.visible").click();

    // Verify error message
    cy.get("[data-cy=error-message]").should("contain", "Time slot already booked");
  });

  it("should cancel scheduled interview", () => {
    cy.visit("/schedule-interviews");

    // Mock API responses
    cy.mockApiResponse("GET", "/api/jobs", {
      statusCode: 200,
      body: [{ id: 1, title: "Software Engineer" }],
    });

    cy.mockApiResponse("GET", "/api/candidates", {
      statusCode: 200,
      body: [{ id: 1, name: "John Doe" }],
    });

    cy.mockApiResponse("GET", "/api/interviewers", {
      statusCode: 200,
      body: [{ id: 1, name: "Jane Smith" }],
    });

    cy.mockApiResponse("POST", "/api/interviews", {
      statusCode: 200,
      body: { success: true },
    });

    cy.mockApiResponse("DELETE", "/api/interviews/*", {
      statusCode: 200,
      body: { success: true },
    });

    // Schedule an interview first
    cy.get("[data-cy=select-job]").should("be.visible").select("Software Engineer");
    cy.get("[data-cy=select-candidate]").should("be.visible").select("John Doe");
    cy.get("[data-cy=select-interviewer]").should("be.visible").select("Jane Smith");
    cy.get("[data-cy=select-date]").should("be.visible").type(cy.getFutureDate());
    cy.get("[data-cy=select-time]").should("be.visible").type(cy.getTimeSlot());
    cy.get("[data-cy=schedule-button]").should("be.visible").click();

    // Cancel the interview
    cy.get("[data-cy=cancel-interview]").first().should("be.visible").click();
    cy.get("[data-cy=confirm-cancel]").should("be.visible").click();

    // Verify cancellation
    cy.get("[data-cy=success-message]").should("contain", "Interview cancelled");
  });
});
