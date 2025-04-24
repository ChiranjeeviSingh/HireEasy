import React from "react";
import { render, screen, fireEvent, waitFor, act } from "@testing-library/react";
import { JobPosting } from "./components/JobPosting";
import { Questionnaire } from "./components/Questionnaire";
import { Register } from "./components/Register";
import { Login } from "./components/Login";
import { ViewJobs } from "./components/ViewJobs";
import { ShareJob } from "./components/ShareJob";
import InterviewerCalendar from "./components/InterviewerCalendar";
import { BrowserRouter as Router } from "react-router-dom";

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn().mockReturnValue("fake-token"),
  setItem: jest.fn(),
  clear: jest.fn(),
  removeItem: jest.fn(),
};
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
  writable: true
});

// Mock fetch to return empty arrays for data fetching
global.fetch = jest.fn().mockImplementation(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve([]),
  })
);

// Mock alert to prevent Not implemented warning
window.alert = jest.fn();

// Clear mock calls between tests
beforeEach(() => {
  window.alert.mockClear();
  localStorageMock.getItem.mockClear();
  localStorageMock.setItem.mockClear();
  localStorageMock.clear.mockClear();
  localStorageMock.removeItem.mockClear();
  fetch.mockClear();
});

describe("Login Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <Login />
      </Router>
    );
  });

  test("renders login page with correct title and subtitle", () => {
    expect(screen.getByText("HireEasy")).toBeInTheDocument();
    expect(screen.getByText(/welcome to careerbuilder/i)).toBeInTheDocument();
    expect(screen.getByText(/sign in and start hiring/i)).toBeInTheDocument();
  });

  test("renders create account link", () => {
    expect(screen.getByText(/don't have an account/i)).toBeInTheDocument();
    expect(screen.getByText(/create one now/i)).toBeInTheDocument();
  });

  test("validates email format", () => {
    const emailInput = screen.getByPlaceholderText(/enter your email/i);
    fireEvent.change(emailInput, { target: { value: "invalid-email" } });
    fireEvent.blur(emailInput);
    expect(emailInput.validity.valid).toBe(false);
  });
});

describe("JobPosting Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <JobPosting />
      </Router>
    );
  });

  test("renders job posting form with correct title", () => {
    expect(screen.getByText(/create job posting/i)).toBeInTheDocument();
  });

  test("renders dashboard navigation button", () => {
    expect(screen.getByText(/dashboard/i)).toBeInTheDocument();
  });
});

describe("Questionnaire Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <Questionnaire />
      </Router>
    );
  });

  test("renders initial questions with correct types", () => {
    expect(screen.getByDisplayValue("What is your gender?")).toBeInTheDocument();
    expect(screen.getByDisplayValue("Education Level")).toBeInTheDocument();
    expect(screen.getByDisplayValue("Which programming languages do you know?")).toBeInTheDocument();
    expect(screen.getByDisplayValue("Work Experience (Years)")).toBeInTheDocument();
  });

  test("allows editing question text", () => {
    const questionText = "New question text";
    const questionInputs = screen.getAllByRole("textbox");
    fireEvent.change(questionInputs[1], { target: { value: questionText } });
    expect(questionInputs[1].value).toBe(questionText);
  });

  test("handles question options", () => {
    const genderOptions = screen.getAllByDisplayValue(/(Male|Female|Other)/);
    expect(genderOptions).toHaveLength(3);
    
    const programmingOptions = screen.getAllByDisplayValue(/(JavaScript|Python|C\+\+)/);
    expect(programmingOptions).toHaveLength(3);
  });

  test("adds new options to questions", () => {
    const addOptionButtons = screen.getAllByText("➕ Add Option");
    fireEvent.click(addOptionButtons[0]); // Click first Add Option button
    
    const genderOptions = screen.getAllByDisplayValue(/(Male|Female|Other|)/);
    expect(genderOptions.length).toBeGreaterThan(3);
  });

  test("validates form ID input", () => {
    // Modified to use getAllByRole and find the one with proper name attribute
    const formIdInputs = screen.getAllByRole("textbox");
    const formIdInput = formIdInputs.find(input => input.getAttribute('name') === 'FormID');
    
    expect(formIdInput).toBeDefined();
    expect(formIdInput).toHaveAttribute("name", "FormID");
    expect(formIdInput).toBeRequired();
  });
});

describe("Register Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <Register />
      </Router>
    );
  });

  test("renders registration page with correct title", () => {
    expect(screen.getByText("HireEasy")).toBeInTheDocument();
    expect(screen.getByText(/sign up and join/i)).toBeInTheDocument();
  });

  test("shows validation errors for empty required fields", () => {
    fireEvent.click(screen.getByRole("button", { name: /sign up/i }));
    
    const requiredFields = screen.getAllByRole("textbox");
    requiredFields.forEach(field => {
      expect(field.validity.valid).toBe(false);
    });
  });

  test("allows role selection between HR and Interviewer", () => {
    const roleSelect = screen.getByRole("combobox");
    expect(roleSelect).toHaveValue("HR");
    
    fireEvent.change(roleSelect, { target: { value: "Interviewer" } });
    expect(roleSelect).toHaveValue("Interviewer");
  });
});

// New tests for ViewJobs component
describe("ViewJobs Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <ViewJobs />
      </Router>
    );
  });

  test("renders view jobs page with the correct title", () => {
    expect(screen.getByText("View Jobs")).toBeInTheDocument();
  });

  test("renders dashboard button", () => {
    // Use getAllByText and test that at least one is in the document
    const dashboardElements = screen.getAllByText(/dashboard/i);
    expect(dashboardElements.length).toBeGreaterThan(0);
  });

  test("renders job filter input", () => {
    // Test for the label text directly
    expect(screen.getByText("Filter by Job ID:")).toBeInTheDocument();
    
    // And check the input exists
    expect(screen.getByPlaceholderText(/enter job id/i)).toBeInTheDocument();
  });

  test("allows filtering jobs", () => {
    const filterInput = screen.getByPlaceholderText(/enter job id/i);
    fireEvent.change(filterInput, { target: { value: "JOB123" } });
    
    // JOB123 should be visible, JOB456 should be filtered out
    expect(screen.getByText("JOB123")).toBeInTheDocument();
    expect(screen.queryByText("JOB456")).not.toBeInTheDocument();
  });
});

// Tests for ShareJob component
describe("ShareJob Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <ShareJob />
      </Router>
    );
  });

  test("renders share job page with the correct title", () => {
    expect(screen.getByText("Share Job")).toBeInTheDocument();
  });

  test("renders dashboard button", () => {
    expect(screen.getByText(/dashboard/i)).toBeInTheDocument();
  });

  test("renders job ID and form ID selection dropdowns", () => {
    // Use more specific selectors
    expect(screen.getAllByText(/select job id/i)[0]).toBeInTheDocument();
    expect(screen.getAllByText(/select form id/i)[0]).toBeInTheDocument();
  });

  test("renders generate job application button", () => {
    expect(screen.getByText(/generate job application/i)).toBeInTheDocument();
  });
});

// Tests for InterviewerCalendar component
describe("InterviewerCalendar Component", () => {
  beforeEach(() => {
    render(
      <Router>
        <InterviewerCalendar />
      </Router>
    );
  });

  test("renders calendar page with the correct title", () => {
    expect(screen.getByText("Set Your Availability")).toBeInTheDocument();
  });

  test("renders date and time input fields", () => {
    expect(screen.getByText(/pick a date/i)).toBeInTheDocument();
    expect(screen.getByText(/start time/i)).toBeInTheDocument();
    expect(screen.getByText(/end time/i)).toBeInTheDocument();
  });

  test("renders save and cancel buttons", () => {
    expect(screen.getByRole("button", { name: /save/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /cancel/i })).toBeInTheDocument();
  });

  test("renders scheduled slots section", () => {
    expect(screen.getByText("Your Scheduled Slots")).toBeInTheDocument();
  });
});