import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function JobPosting() {
  const navigate = useNavigate();

  const initialState = {
    JobID: "",
    Info1: "",
    Info2: "",
    Info3: "",
    Info4: "",
    Info5: "",
    Info6: "",
    Info7: "",
    Info8: "",
    Info9: "",
    Info10: "",
  };

  const placeholders = {
    Info1: "Enter Job Location",
    Info2: "Enter Job Description",
    Info3: "Enter Salary Details",
    Info4: "Enter Experience Required",
    Info5: "Additional Info",
    Info6: "Additional Info",
    Info7: "Additional Info",
    Info8: "Additional Info",
    Info9: "Additional Info",
    Info10: "Additional Info",
  };

  const [jobData, setJobData] = useState(initialState);
  const [formSubmitted, setFormSubmitted] = useState(false);

  const handleChange = (e) => {
    setJobData({ ...jobData, [e.target.name]: e.target.value });
  };

  const handleRemove = (infoKey) => {
    setJobData({ ...jobData, [infoKey]: "none" });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("Job Created:", jobData);
    alert("Job Posted Successfully!");
    setFormSubmitted(true);
  };

  const handleNewJob = () => {
    setJobData(initialState);
    setFormSubmitted(false);
  };

  return (
    <div className="flex flex-col min-h-screen bg-gray-100">
      {/* ✅ Dashboard Button in Top-Left */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 text-lg bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      {/* ✅ Job Posting Form Section (Card Style) */}
      <div className="flex-grow flex justify-center items-center">
        <div className="bg-white shadow-2xl rounded-xl p-8 w-full max-w-3xl">
          <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-8">
            Create Job Posting
          </h2>

          <form onSubmit={handleSubmit}>
            {/* Job ID Field */}
            <div className="mb-4">
              <label className="block font-medium mb-1 text-lg">
                Job ID (Required):
              </label>
              <input
                type="text"
                name="JobID"
                value={jobData.JobID}
                onChange={handleChange}
                required
                className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
                placeholder="Enter Job ID"
              />
            </div>

            {/* Job Fields (Info1 - Info5 always visible) */}
            {["Info1", "Info2", "Info3", "Info4", "Info5"].map((key, index) => (
              <div key={index} className="mb-4">
                <label className="block font-medium mb-1 text-lg">
                  {placeholders[key]}:
                </label>
                <textarea
                  name={key}
                  value={jobData[key]}
                  onChange={handleChange}
                  rows="3"
                  className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
                  placeholder={placeholders[key]}
                />
              </div>
            ))}

            {/* Info6 - Info10 with Remove Button */}
            {["Info6", "Info7", "Info8", "Info9", "Info10"].map((key, index) => (
              <div
                key={index}
                className={`mb-4 ${jobData[key] === "none" ? "hidden" : "block"}`}
              >
                <label className="block font-medium mb-1 text-lg">
                  {placeholders[key]}:
                </label>
                <textarea
                  name={key}
                  value={jobData[key]}
                  onChange={handleChange}
                  rows="3"
                  className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
                  placeholder={placeholders[key]}
                />
                <button
                  type="button"
                  onClick={() => handleRemove(key)}
                  className="ml-2 mt-2 px-3 py-1 bg-red-500 text-white rounded-lg hover:bg-red-600 transition"
                >
                  ❌ Remove
                </button>
              </div>
            ))}

            {/* ✅ "Post Job" Button */}
            <button
              type="submit"
              className="w-full py-3 mt-6 bg-green-500 text-white rounded-lg text-lg hover:bg-green-600 transition"
            >
              Post Job
            </button>

            {formSubmitted && (
              <button
                type="button"
                onClick={handleNewJob}
                className="w-full py-3 mt-2 bg-gray-500 text-white rounded-lg text-lg hover:bg-gray-600 transition"
              >
                New Job
              </button>
            )}
          </form>
        </div>
      </div>

      {/* ✅ Footer Section */}
      <footer className="bg-gray-900 text-white py-6 px-12">
        <div className="container mx-auto flex flex-col md:flex-row justify-between">
          {/* About Section */}
          <div>
            <h2 className="text-2xl font-bold">HireEasy</h2>
            <p className="mt-3 text-gray-400">
              Your trusted job portal for connecting top talent with top companies. 
              Find and post jobs with ease.
            </p>
          </div>

          {/* Office Information */}
          <div>
            <h3 className="text-xl font-bold mb-3">Office</h3>
            <p>123 Business Park,</p>
            <p>New York, USA</p>
            <p>Email: support@hireeasy.com</p>
            <p>Phone: +1 987-654-3210</p>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-xl font-bold mb-3">Quick Links</h3>
            <ul className="space-y-2">
              <li><a href="#" className="hover:text-gray-300 transition">Dashboard</a></li>
              <li><a href="#" className="hover:text-gray-300 transition">Job Applications</a></li>
              <li><a href="#" className="hover:text-gray-300 transition">View Jobs</a></li>
            </ul>
          </div>
        </div>
        <div className="text-center text-gray-400 mt-6">
          HireEasy © {new Date().getFullYear()} - All Rights Reserved
        </div>
      </footer>
    </div>
  );
}

export default JobPosting;
