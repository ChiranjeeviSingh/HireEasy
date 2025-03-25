import React from "react";
import { render, screen } from "@testing-library/react";
import { BrowserRouter as Router } from "react-router-dom";
import App from "./App";
import Questionnaire from "./components/Questionnaire";
import ShareJob from "./components/ShareJob";
import ViewJobs from "./components/ViewJobs";
import JobPosting from "./components/JobPosting";
import JobApplications from "./components/JobApplications"


test("Render App", () => {
  render(
    <Router>
      <App />
    </Router>
  );

  const hireEasyText = screen.getByText(/HireEasy/i);
  expect(hireEasyText).toBeInTheDocument();
});


test("Share Jobs", () => {
    render(
      <Router>
        <ShareJob />
      </Router>
    );
  
    const hireEasyText = screen.getByText(/Share Job/i);
    expect(hireEasyText).toBeInTheDocument();
  });


test("View Jobs", () => {
    render(
      <Router>
        <ViewJobs />
      </Router>
    );
  
    const hireEasyText = screen.getByText(/View Jobs/i);
    expect(hireEasyText).toBeInTheDocument();
  });

  test("Questionnaire", () => {
    render(
      <Router>
        <Questionnaire />
      </Router>
    );
  
    const hireEasyText = screen.getByText(/Create a Job Questionnaire/i);
    expect(hireEasyText).toBeInTheDocument();
  });


  test("JobPosting", () => {
    render(
      <Router>
        <JobPosting />
      </Router>
    );
  
    const hireEasyText = screen.getByText(/Create Job Posting/i);
    expect(hireEasyText).toBeInTheDocument();
  });

  test("JobApplications", () => {
    render(
      <Router>
        <JobApplications />
      </Router>
    );
  
    // Check if the heading exists
    expect(screen.getByRole("heading", { name: /Job Applications/i })).toBeInTheDocument();
  });
  


  



  