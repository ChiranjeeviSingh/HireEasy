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
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      const data = await response.json();
      console.log("Fetched Jobs:", data);
      if (response.ok) {
        setJobPostings(data);
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
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      const data = await response.json();
      console.log("Fetched Forms:", data);
      if (response.ok) {
        setQuestionnaires(data);
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

    const token = localStorage.getItem("token");

    try {
      const response = await fetch(`${API_BASE}/jobs/${selectedJobId}/forms`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          form_template_id: selectedFormId,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.msg || "Failed to link job to form.");
      }

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
    <div className="relative min-h-screen flex flex-col justify-between bg-gray-100">
      {/* ✅ Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      {/* ✅ Main Card */}
      <div className="flex-grow flex justify-center items-center">
        <div className="bg-white shadow-lg rounded-xl p-10 w-full max-w-3xl">
          <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-8">
            Share Job
          </h2>

          {error && <p className="text-red-500 text-center">{error}</p>}

          {/* Dropdowns for Job ID and Form ID */}
          <div className="mb-6">
            <label className="block font-medium mb-2 text-lg">Select Job ID:</label>
            <select
              value={selectedJobId}
              onChange={(e) => setSelectedJobId(e.target.value)}
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
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

          <div className="mb-6">
            <label className="block font-medium mb-2 text-lg">Select Form ID:</label>
            <select
              value={selectedFormId}
              onChange={(e) => setSelectedFormId(e.target.value)}
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
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
            className="w-full py-3 bg-green-500 text-white text-lg rounded-lg hover:bg-green-600 transition"
          >
            Generate Job Application
          </button>

          {/* Display Job Details and Questionnaire */}
          {jobDetails && questions && (
            <div className="mt-6">
              <h3 className="text-2xl font-semibold">Job Details</h3>
              <p className="text-lg">
                <strong>Job ID:</strong> {jobDetails.job_id}
              </p>

              {Object.keys(jobDetails)
                .filter((key) => key !== "job_id")
                .map((key, index) => (
                  <p key={index} className="text-lg">
                    <strong>{key.replace("Info", "Detail")}:</strong> {jobDetails[key]}
                  </p>
                ))}

              <h3 className="text-2xl font-semibold mt-6">Job Questionnaire</h3>
              {questions.fields.map((question) => (
                <div key={question.question_id} className="mt-3">
                  <p className="text-lg font-bold">{question.question_text}</p>

                  {question.question_type === "text" && (
                    <input type="text" className="w-full p-3 border border-gray-300 rounded-lg" />
                  )}
                  {question.question_type === "radio" &&
                    question.options.map((opt) => (
                      <label key={opt} className="block">
                        <input type="radio" name={question.question_id} value={opt} className="mr-2" /> {opt}
                      </label>
                    ))}
                  {question.question_type === "checkbox" &&
                    question.options.map((opt) => (
                      <label key={opt} className="block">
                        <input type="checkbox" value={opt} className="mr-2" /> {opt}
                      </label>
                    ))}
                  {question.question_type === "file" && <input type="file" className="w-full p-2 border border-gray-300 rounded-lg" />}
                </div>
              ))}

              <button
                type="button"
                onClick={handleShare}
                className="w-full py-3 mt-6 bg-blue-500 text-white text-lg rounded-lg hover:bg-blue-600 transition"
              >
                Share
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ShareJob;