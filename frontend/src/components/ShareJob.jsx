import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

export function ShareJob() {
  const navigate = useNavigate(); // Hook for navigation

  const [jobPostings, setJobPostings] = useState([]);
  const [questionnaires, setQuestionnaires] = useState([]);
  const [selectedJobId, setSelectedJobId] = useState("");
  const [selectedFormId, setSelectedFormId] = useState("");
  const [jobDetails, setJobDetails] = useState(null);
  const [questions, setQuestions] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Fetch job postings data from backend
  const fetchJobPostings = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/jobs`, {
        headers: {
          "Authorization": `Bearer ${localStorage.getItem("token")}`,
        },
      });
      const data = await response.json();
      console.log("Fetched Jobs:", data);
      if (response.ok) {
        setJobPostings(data); // Assuming the data is an array of job postings
      } else {
        setError("Failed to fetch job postings.");
      }
    } catch (err) {
      setError("Error fetching job postings.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Fetch form templates data from backend
  const fetchFormTemplates = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/forms/templates`, {
        headers: {
          "Authorization": `Bearer ${localStorage.getItem("token")}`,
        },
      });
      const data = await response.json();
      console.log("Fetched Forms:", data);
      if (response.ok) {
        setQuestionnaires(data); // Assuming the data is an array of form templates
      } else {
        setError("Failed to fetch form templates.");
      }
    } catch (err) {
      setError("Error fetching form templates.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Fetch job and form data after selecting both IDs
  const fetchJobAndFormData = () => {
    if (!selectedJobId || !selectedFormId) {
      alert("Please select both a Job ID and a Form ID.");
      return;
    }

    const jobData = jobPostings.find((job) => job.job_id === selectedJobId);
    const questionData = questionnaires.find(
      (form) => form.form_template_id === selectedFormId
    );

    if (!jobData || !questionData) {
      alert("Invalid selection. Please try again.");
      return;
    }

    setJobDetails(jobData);
    setQuestions(questionData);
  };

  // Handle share button click (show dummy share link)
  const handleShare = async () => {
    if (!selectedJobId || !selectedFormId) {
      alert("Please select both a Job ID and a Form ID.");
      return;
    }

    // Assuming that the job has been linked to the form template
    const token = localStorage.getItem("token");

    try {
      const response = await fetch(`${API_BASE}/jobs/${selectedJobId}/forms`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({
          form_template_id: selectedFormId, // Sending the selected form template
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.msg || "Failed to link job to form.");
      }

      // Assuming we get a form UUID from the backend
      const shareableLink = `https://hireeasy.com/apply?job=${selectedJobId}&form=${data.form_uuid}`;
      alert(`Share using this link:\n${shareableLink}`);

    } catch (err) {
      alert(err.message || "Failed to share the job.");
    }
  };

  // Load job postings and form templates when the component mounts
  useEffect(() => {
    fetchJobPostings();
    fetchFormTemplates();
  }, []);

  return (
    <div style={{ textAlign: "center", marginTop: "20px", position: "relative" }}>
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

      <h2>Share Job</h2>

      {/* Dropdowns for Job ID and Form ID */}
      <div style={{ maxWidth: "600px", margin: "auto" }}>
        <div style={{ marginBottom: "10px" }}>
          <label>Select Job ID: </label>
          <select
            value={selectedJobId}
            onChange={(e) => setSelectedJobId(e.target.value)}
            style={{ width: "100%", padding: "8px", fontSize: "16px" }}
          >
            <option value="">-- Select Job ID --</option>
            {jobPostings.length > 0 ? (
              jobPostings.map((job) => (
                <option key={job.job_id} value={job.job_id}>
                  {job.job_id}
                </option>
              ))
            ) : (
              <option value="">No jobs available</option>
            )}
          </select>
        </div>

        <div style={{ marginBottom: "10px" }}>
          <label>Select Form ID: </label>
          <select
            value={selectedFormId}
            onChange={(e) => setSelectedFormId(e.target.value)}
            style={{ width: "100%", padding: "8px", fontSize: "16px" }}
          >
            <option value="">-- Select Form ID --</option>
            {questionnaires.length > 0 ? (
              questionnaires.map((form) => (
                <option key={form.form_template_id} value={form.form_template_id}>
                  {form.form_template_id}
                </option>
              ))
            ) : (
              <option value="">No forms available</option>
            )}
          </select>
        </div>

        <button
          onClick={fetchJobAndFormData}
          style={{ cursor: "pointer", padding: "10px 15px", marginTop: "10px" }}
        >
          Generate Job Application
        </button>
      </div>

      {/* Display Job Details and Questionnaire */}
      {jobDetails && questions && (
        <div
          style={{
            maxWidth: "600px",
            margin: "auto",
            marginTop: "20px",
            textAlign: "left",
          }}
        >
          <h3>Job Details</h3>
          <p>
            <strong>Job ID:</strong> {jobDetails.job_id}
          </p>

          {/* Display all job details dynamically */}
          {Object.keys(jobDetails)
            .filter((key) => key !== "job_id") // Exclude JobID as it's already displayed
            .map((key, index) => (
              <p key={index}>
                <strong>{key.replace("Info", "Detail ")}:</strong>{" "}
                {jobDetails[key]}
              </p>
            ))}

          <h3>Job Questionnaire</h3>
          {questions.fields.map((question) => (
            <div key={question.question_id} style={{ marginBottom: "10px" }}>
              <p>
                <strong>{question.question_text}</strong>
              </p>

              {question.question_type === "text" && <input type="text" />}
              {question.question_type === "radio" &&
                question.options.map((opt) => (
                  <label key={opt}>
                    <input type="radio" name={question.question_id} value={opt} /> {opt}
                  </label>
                ))}
              {question.question_type === "checkbox" &&
                question.options.map((opt) => (
                  <label key={opt}>
                    <input type="checkbox" value={opt} /> {opt}
                  </label>
                ))}
              {question.question_type === "file" && <input type="file" />}
            </div>
          ))}

          {/* Share Button */}
          <button
            type="button"
            onClick={handleShare}
            style={{
              marginTop: "10px",
              cursor: "pointer",
              padding: "10px 15px",
            }}
          >
            Share
          </button>
        </div>
      )}
    </div>
  );
}

export default ShareJob;