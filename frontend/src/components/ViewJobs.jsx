import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function ViewJobs() {
  const navigate = useNavigate();

  // Mock Job Data
  const mockJobs = [
    { jobId: "JOB123", jobLink: "https://company.com/apply/JOB123" },
    { jobId: "JOB456", jobLink: "https://company.com/apply/JOB456" },
    { jobId: "JOB789", jobLink: "https://company.com/apply/JOB789" },
    { jobId: "JOB101", jobLink: "https://company.com/apply/JOB101" },
  ];

  const [searchJobId, setSearchJobId] = useState("");

  // Filter jobs based on user input
  const filteredJobs = mockJobs.filter((job) =>
    job.jobId.toLowerCase().includes(searchJobId.toLowerCase())
  );

  return (
    <div className="relative min-h-screen bg-gray-100 flex flex-col justify-between">
      {/* ✅ Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 text-lg bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      {/* ✅ Content Section (Proper Layout) */}
      <div className="flex-grow flex flex-col justify-center items-center mt-12 mb-20 px-6">
        <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-8">
          View Jobs
        </h2>

        {/* ✅ Job Filter Input */}
        <div className="w-full max-w-lg mb-6">
          <label className="block font-medium mb-2 text-lg">Filter by Job ID:</label>
          <input
            type="text"
            value={searchJobId}
            onChange={(e) => setSearchJobId(e.target.value)}
            placeholder="Enter Job ID..."
            className="w-full p-3 border border-gray-500 rounded-lg text-lg focus:ring-2 focus:ring-green-400"
          />
        </div>

        {/* ✅ Job Table */}
        <div className="w-full max-w-3xl overflow-x-auto">
          <table className="w-full border-collapse shadow-lg bg-white rounded-lg overflow-hidden">
            <thead>
              <tr className="bg-green-500 text-white">
                <th className="p-4 text-lg border-b">Job ID</th>
                <th className="p-4 text-lg border-b">Job Link</th>
              </tr>
            </thead>
            <tbody>
              {filteredJobs.length > 0 ? (
                filteredJobs.map((job) => (
                  <tr key={job.jobId} className="text-center border-b">
                    <td className="p-4 text-lg">{job.jobId}</td>
                    <td className="p-4 text-lg">
                      <a
                        href={job.jobLink}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 hover:text-blue-800 transition"
                      >
                        {job.jobLink}
                      </a>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="2" className="p-4 text-lg text-gray-500 text-center">
                    No jobs found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* ✅ Footer Section (Properly Separated) */}
      <footer className="bg-gradient-to-r from-gray-900 to-black text-white py-8 px-12">
        <div className="container mx-auto grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* About Section */}
          <div>
            <h2 className="text-2xl font-bold">HireEasy</h2>
            <p className="mt-3 text-gray-400">
              Browse job listings and apply directly. We help you find the best opportunities in the market.
            </p>
          </div>

          {/* Office Information */}
          <div>
            <h3 className="text-xl font-bold mb-3">Office</h3>
            <p>456 Talent Hub,</p>
            <p>San Francisco, USA</p>
            <p>Email: support@hireeasy.com</p>
            <p>Phone: +1 987-654-3210</p>
          </div>

          {/* Useful Links */}
          <div>
            <h3 className="text-xl font-bold mb-3">Quick Links</h3>
            <ul className="space-y-2">
              <li><a href="#" className="hover:text-gray-300 transition">Dashboard</a></li>
              <li><a href="#" className="hover:text-gray-300 transition">Job Applications</a></li>
            </ul>
          </div>
        </div>
        <div className="text-center text-gray-400 mt-8">
          HireEasy © {new Date().getFullYear()} - All Rights Reserved
        </div>
      </footer>
    </div>
  );
}

export default ViewJobs;
