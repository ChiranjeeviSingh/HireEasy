import React from "react";
import { Route, Routes } from "react-router-dom";
import Login from "./components/Login.jsx";
import Dashboard from "./pages/Dashboard.jsx";
import JobPosting from "./components/JobPosting.jsx";
import Questionnaire from "./components/Questionnaire.jsx";
import ShareJob from "./components/ShareJob.jsx";
import JobApplications from "./components/JobApplications.jsx";
import ViewJobs from "./components/ViewJobs.jsx"; // New import
import Register from "./components/Register.jsx";

console.log("App.jsx is rendering...");

export function App() {
  return (
    <div>
      <Routes>
        <Route exact path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/job-posting" element={<JobPosting />} />
        <Route path="/questionnaire" element={<Questionnaire />} />
        <Route path="/share-job" element={<ShareJob />} />
        <Route path="/job-applications" element={<JobApplications />} />
        <Route path="/view-jobs" element={<ViewJobs />} /> {/* New Route */}
      </Routes>
    </div>
  );
}

export default App;