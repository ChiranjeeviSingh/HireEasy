import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function ShareJob() {
  const navigate = useNavigate();

  // Mock job postings data
  const jobPostings = [
    {
      JobID: "JOB123",
      Info1: "New York, USA",
      Info2: "Software Engineer Role",
      Info3: "$120,000 per year",
      Info4: "5+ years experience required",
      Info5: "Remote Work Available",
    },
    {
      JobID: "JOB456",
      Info1: "San Francisco, USA",
      Info2: "Data Scientist Role",
      Info3: "$135,000 per year",
      Info4: "3+ years experience required",
      Info5: "Hybrid Work Model",
    },
  ];

  const [selectedJobId, setSelectedJobId] = useState("");
  const [jobDetails, setJobDetails] = useState(null);

  // Handle dropdown selection
  const handleJobChange = (e) => setSelectedJobId(e.target.value);

  // Fetch job data
  const fetchJobData = () => {
    if (!selectedJobId) {
      alert("Please select a Job ID.");
      return;
    }

    const jobData = jobPostings.find((job) => job.JobID === selectedJobId);

    if (!jobData) {
      alert("Invalid selection. Please try again.");
      return;
    }

    setJobDetails(jobData);
  };

  return (
    <div
      className="relative min-h-screen bg-cover bg-center flex flex-col justify-between"
      // style={{
      //   backgroundImage:
      //     "url('https://www.shutterstock.com/image-vector/vector-business-illustration-small-people-260nw-1022567779.jpg')",
      // }}
    >
      {/* ✅ Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 text-lg bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      {/* ✅ Styled Card */}
      <div className="flex-grow flex justify-center items-center">
        <div className="bg-white bg-opacity-90 shadow-lg rounded-xl p-10 w-full max-w-3xl h-[500px] overflow-y-auto">
          <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-8">
            Share Job
          </h2>

          {/* Dropdown for Job ID */}
          <div className="mb-6">
            <label className="block font-medium mb-2 text-lg">Select Job ID:</label>
            <select
              value={selectedJobId}
              onChange={handleJobChange}
              className="w-full p-3 border border-gray-500 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
            >
              <option value="">-- Select Job ID --</option>
              {jobPostings.map((job) => (
                <option key={job.JobID} value={job.JobID}>
                  {job.JobID}
                </option>
              ))}
            </select>
          </div>

          {/* Generate Button */}
          <button
            onClick={fetchJobData}
            className="w-full py-3 bg-green-500 text-white text-lg rounded-lg hover:bg-green-600 transition"
          >
            Generate Job Application
          </button>

          {/* Display Job Details */}
          {jobDetails && (
            <div className="mt-6">
              <h3 className="text-2xl font-semibold">Job Details</h3>
              {Object.keys(jobDetails).map((key, index) => (
                <p key={index} className="text-lg">
                  <strong>{key.replace("Info", "Detail")}:</strong> {jobDetails[key]}
                </p>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* ✅ Footer Section */}
      <footer className="bg-gradient-to-r from-gray-900 to-black text-white py-8 px-12">
        <div className="container mx-auto grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* About Section */}
          <div>
            <h2 className="text-2xl font-bold">HireEasy</h2>
            <p className="mt-3 text-gray-400">
              HireEasy helps recruiters connect with the best talent by streamlining job postings, 
              applications, and hiring processes efficiently.
            </p>
          </div>

          {/* Office Information */}
          <div>
            <h3 className="text-xl font-bold mb-3">Office</h3>
            <p>123 Recruitment St,</p>
            <p>New York, USA</p>
            <p>Email: contact@hireeasy.com</p>
            <p>Phone: +1 234-567-890</p>
          </div>

          {/* Useful Links */}
          <div>
            <h3 className="text-xl font-bold mb-3">Links</h3>
            <ul className="space-y-2">
              <li>
                <a href="#" className="hover:text-gray-300 transition">
                  Home
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-gray-300 transition">
                  Job Listings
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-gray-300 transition">
                  About Us
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-gray-300 transition">
                  Contact
                </a>
              </li>
            </ul>
          </div>

          {/* Newsletter */}
          <div>
            <h3 className="text-xl font-bold mb-3">Newsletter</h3>
            <p className="text-gray-400 mb-3">Subscribe to stay updated with the latest job postings.</p>
            <div className="flex">
              <input
                type="email"
                placeholder="Enter your email"
                className="w-full p-2 rounded-l-lg text-black"
              />
              <button className="bg-blue-500 px-4 rounded-r-lg hover:bg-blue-600 transition">
                ➝
              </button>
            </div>
          </div>
        </div>

        {/* Bottom Footer */}
        <div className="text-center text-gray-400 mt-8">
          HireEasy © {new Date().getFullYear()} - All Rights Reserved
        </div>
      </footer>
    </div>
  );
}

export default ShareJob;
