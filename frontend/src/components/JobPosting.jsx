import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function JobPosting() {
  const navigate = useNavigate(); // Hook for navigation

  // Initial form state
  const initialState = {
    JobID: "", // Separate Job ID
    Info1: "", // Job Location
    Info2: "", // Job Description
    Info3: "", // Salary
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
    Info2: "Enter Job Description",
    Info3: "Enter Salary Details",
    Info4: "Enter Experience Required",
    Info5: "Skills Required",
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

    // Map frontend fields to backend structure (correct field names as in the backend model)
    const jobDataMapped = {
      job_id: jobData.JobID.trim(), // map to `job_id`
      job_title: jobData.Info2.trim(), // map to `job_title`
      job_description: jobData.Info3.trim(), // map to `job_description`
      job_status: "Open", // map to `job_status`
      skills_required: jobData.Info5.trim() || "EMPTY", // skills_required as a string
      location: jobData.Info1.trim() || "EMPTY", // location as string or "EMPTY"
      experience: jobData.Info4.trim() || "EMPTY", // experience as string or "EMPTY"
      // Pack Info6 to Info10, if empty, send "EMPTY"
      info6: jobData.Info6.trim() || "EMPTY",
      info7: jobData.Info7.trim() || "EMPTY",
      info8: jobData.Info8.trim() || "EMPTY",
      info9: jobData.Info9.trim() || "EMPTY",
      info10: jobData.Info10.trim() || "EMPTY",
    };

    // Log the mapped data to see if everything is correctly formatted
    console.log("Mapped Job Data:", jobDataMapped);

    // Check if all required fields are valid
    if (
      !jobDataMapped.job_id ||
      !jobDataMapped.job_title ||
      !jobDataMapped.job_description ||
      !jobDataMapped.location ||
      !jobDataMapped.experience ||
      !jobDataMapped.skills_required
    ) {
      setError("Please fill in all required fields and provide at least one skill.");
      console.log("Validation failed. Missing required fields.");
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`${API_BASE}/jobs`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`, // Attach JWT token
        },
        body: JSON.stringify(jobDataMapped), // Sending the corrected mapped data
      });

      const data = await response.json();

      console.log("Response Data:", data); // Log the response from the backend (frontend logging)

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
    <div
      style={{ textAlign: "center", marginTop: "20px", position: "relative" }}
    >
      {/* Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        style={{
          position: "absolute",
          top: "10px",
          left: "10px",
          padding: "5px 10px",
          fontSize: "16px",
          cursor: "pointer",
        }}
      >
        ⬅️ Dashboard
      </button>

      <h2>Create Job Posting</h2>

      {error && <p className="error-text" style={{ color: "red" }}>{error}</p>}
      {success && <p className="success-text" style={{ color: "green" }}>{success}</p>}

      <form
        onSubmit={handleSubmit}
        style={{ maxWidth: "600px", margin: "auto" }}
      >
        {/* Separate Job ID Field */}
        <div style={{ marginBottom: "10px" }}>
          <label>Job ID (Required): </label>
          <input
            type="text"
            name="JobID"
            value={jobData.JobID}
            onChange={handleChange}
            required
            style={{ width: "100%", padding: "8px", fontSize: "16px" }}
            placeholder="Enter Job ID"
          />
        </div>

        {/* Job Fields (Info1 - Info5 always visible) */}
        {["Info1", "Info2", "Info3", "Info4", "Info5"].map((key, index) => (
          <div key={index} style={{ marginBottom: "10px" }}>
            <label>{placeholders[key]}:</label>
            <textarea
              name={key}
              value={jobData[key]}
              onChange={handleChange}
              rows="2"
              style={{
                width: "100%",
                padding: "8px",
                fontSize: "16px",
                resize: "none",
              }}
              placeholder={placeholders[key]}
            />
          </div>
        ))}

        {/* Info6 - Info10 with Remove Button */}
        {["Info6", "Info7", "Info8", "Info9", "Info10"].map((key, index) => (
          <div
            key={index}
            style={{
              marginBottom: "10px",
              display: jobData[key] !== "none" ? "block" : "none",
            }}
          >
            <label>{placeholders[key]}:</label>
            <textarea
              name={key}
              value={jobData[key]}
              onChange={handleChange}
              rows="2"
              style={{
                width: "100%",
                padding: "8px",
                fontSize: "16px",
                resize: "none",
              }}
              placeholder={placeholders[key]}
            />
            <button
              type="button"
              onClick={() => handleRemove(key)}
              style={{
                marginLeft: "10px",
                cursor: "pointer",
                background: "red",
                color: "white",
                padding: "5px 10px",
                border: "none",
                borderRadius: "5px",
              }}
            >
              ❌ Remove
            </button>
          </div>
        ))}

        {/* Submit and New Buttons */}
        <button
          type="submit"
          style={{ marginTop: "10px", cursor: "pointer", padding: "10px 15px" }}
          disabled={loading}
        >
          {loading ? "Posting Job..." : "Post Job"}
        </button>

        {formSubmitted && (
          <button
            type="button"
            onClick={handleNewJob}
            style={{
              marginLeft: "10px",
              cursor: "pointer",
              padding: "10px 15px",
            }}
          >
            New
          </button>
        )}
      </form>
    </div>
  );
}

export default JobPosting;