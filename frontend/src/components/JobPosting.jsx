import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function JobPosting() {
  const navigate = useNavigate(); // Hook for navigation

  // Initial form state
  const initialState = {
    JobID: "", // Separate Job ID
    Info1: "", // Job Location
    Info2: "", // Job Title
    Info3: "", // Job Description
    Info4: "", // Experience Required
    Info5: "", // Skills Required
    Info6: "",
    Info7: "",
    Info8: "",
    Info9: "",
    Info10: "",
  };

  const placeholders = {
    Info1: "Enter Job Location",
    Info2: "Enter Job Title",
    Info3: "Enter Job Description",
    Info4: "Enter Experience Required",
    Info5: "Skills Required (comma-separated, e.g., Go, JavaScript, Python)",
    Info6: "Additional Info",
    Info7: "Additional Info",
    Info8: "Additional Info",
    Info9: "Additional Info",
    Info10: "Additional Info",
  };

  const [jobData, setJobData] = useState(initialState);
  const [formSubmitted, setFormSubmitted] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Handle input changes
  const handleChange = (e) => {
    setJobData({ ...jobData, [e.target.name]: e.target.value });
  };

  // Remove optional fields (set them to "none")
  const handleRemove = (infoKey) => {
    setJobData({ ...jobData, [infoKey]: "none" });
  };

  // Submit form
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    const token = localStorage.getItem("token"); // Get JWT token

    // Map frontend fields to backend structure
    const jobDataMapped = {
      job_id: jobData.JobID.trim(),
      job_title: jobData.Info2.trim(),
      job_description: jobData.Info3.trim(),
      job_status: "Open",
      skills_required: jobData.Info5.trim()
        ? jobData.Info5.split(",").map(skill => skill.trim())
        : [],
      location: jobData.Info1.trim() || "EMPTY",
      experience: jobData.Info4.trim() || "EMPTY",
      info6: jobData.Info6.trim() || "EMPTY",
      info7: jobData.Info7.trim() || "EMPTY",
      info8: jobData.Info8.trim() || "EMPTY",
      info9: jobData.Info9.trim() || "EMPTY",
      info10: jobData.Info10.trim() || "EMPTY",
    };

    console.log("Mapped Job Data:", jobDataMapped);

    if (
      !jobDataMapped.job_id ||
      !jobDataMapped.job_title ||
      !jobDataMapped.job_description ||
      !jobDataMapped.location ||
      !jobDataMapped.experience ||
      jobDataMapped.skills_required.length === 0
    ) {
      setError("Please fill in all required fields and provide at least one skill.");
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`${API_BASE}/jobs`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify(jobDataMapped),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.msg || "Failed to create job posting.");
      }

      setSuccess(`Job with ID "${data.jobId}" created successfully!`);
      setFormSubmitted(true);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Reset form for a new job
  const handleNewJob = () => {
    setJobData(initialState);
    setFormSubmitted(false);
    setSuccess("");
    setError("");
  };

  return (
    <div className="flex flex-col min-h-screen bg-gray-100">
      {/* ✅ Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 text-lg bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      {/* ✅ Job Posting Form Section */}
      <div className="flex-grow flex justify-center items-center">
        <div className="bg-white shadow-2xl rounded-xl p-8 w-full max-w-3xl">
          <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-8">
            Create Job Posting
          </h2>

          {error && <p className="text-red-500 text-center">{error}</p>}
          {success && <p className="text-green-500 text-center">{success}</p>}

          <form onSubmit={handleSubmit}>
            {/* Job ID Field */}
            <div className="mb-4">
              <label className="block font-medium mb-1 text-lg">Job ID (Required):</label>
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

            {/* Job Fields */}
            {["Info1", "Info2", "Info3", "Info4", "Info5"].map((key, index) => (
              <div key={index} className="mb-4">
                <label className="block font-medium mb-1 text-lg">{placeholders[key]}:</label>
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
              <div key={index} className={`mb-4 ${jobData[key] === "none" ? "hidden" : "block"}`}>
                <label className="block font-medium mb-1 text-lg">{placeholders[key]}:</label>
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
              disabled={loading}
            >
              {loading ? "Posting Job..." : "Post Job"}
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
    </div>
  );
}

export default JobPosting;
