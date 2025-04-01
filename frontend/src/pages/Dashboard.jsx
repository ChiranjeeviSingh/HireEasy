import React from "react";
import { useNavigate } from "react-router-dom";

export function Dashboard() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100">
      {/* Dashboard Container */}
      <div className="bg-white shadow-lg rounded-2xl p-10 w-full max-w-4xl text-center">
        <h2 className="text-4xl font-bold text-gray-800 tracking-wide mb-6">
          Welcome, HR!
        </h2>
        <p className="text-gray-600 text-lg">
          Manage job postings, applications, and recruit talent seamlessly.
        </p>

        {/* Buttons Container */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-8">
          <DashboardButton
            label="Create Job Posting"
            onClick={() => navigate("/job-posting")}
          />
          <DashboardButton
            label="Create Questionnaire"
            onClick={() => navigate("/questionnaire")}
          />
          <DashboardButton
            label="Share Job"
            onClick={() => navigate("/share-job")}
          />
          <DashboardButton
            label="View Job Applications"
            onClick={() => navigate("/job-applications")}
          />
          <DashboardButton
            label="View Jobs"
            onClick={() => navigate("/view-jobs")}
          />
        </div>
      </div>
    </div>
  );
}

// Reusable Button Component
const DashboardButton = ({ label, onClick }) => {
  return (
    <button
      onClick={onClick}
      className="w-full py-4 bg-blue-500 text-white text-lg font-semibold rounded-lg shadow-md hover:bg-blue-600 transition transform hover:scale-105 focus:outline-none"
    >
      {label}
    </button>
  );
};

export default Dashboard;