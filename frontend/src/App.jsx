import React from "react";
import { Route, Routes } from "react-router-dom";
import Login from "./components/Login.jsx";
import Dashboard from "./pages/Dashboard.jsx";
import JobPosting from "./components/JobPosting.jsx";
import InterviewerProfile from "./components/InterviewerProfile.jsx";
import InterviewerCalendar from "./components/InterviewerCalendar.jsx";
import Questionnaire from "./components/Questionnaire.jsx";
import ShareJob from "./components/ShareJob.jsx";
import JobApplications from "./components/JobApplications.jsx";
import ViewJobs from "./components/ViewJobs.jsx"; // New import
import Register from "./components/Register.jsx";
import Apply from './components/Apply';  // Adjust the path to where Apply.jsx is located
import InterviewerDashboard from "./pages/InterviewerDashboard.jsx";
import ScheduleInterviews from "./components/ScheduleInterviews.jsx";
import InterviewerInterviews from "./components/InterviewerInterviews.jsx";

console.log("App.jsx is rendering...");

export function App() {
  return (
    <div>
      <Routes>
        <Route exact path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/interviewer-dashboard" element={<InterviewerDashboard />}/>
        <Route path="/job-posting" element={<JobPosting />} />
        <Route path="/questionnaire" element={<Questionnaire />} />
        <Route path="/share-job" element={<ShareJob />} />
        <Route path="/job-applications" element={<JobApplications />} />
        <Route path="/view-jobs" element={<ViewJobs />} /> {/* New Route */}
        <Route path="/apply" element={<Apply />} />
        <Route path="/interviewer-profile" element={<InterviewerProfile />} />
        <Route path="/interviewer-calendar" element={<InterviewerCalendar />} />
        <Route path="/schedule-interviews" element={<ScheduleInterviews />} />
        <Route path="/interviewer-interviews" element={<InterviewerInterviews />} />
      </Routes>
    </div>
  );
}

export default App;