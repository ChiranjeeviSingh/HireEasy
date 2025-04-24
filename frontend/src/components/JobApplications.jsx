import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

function JobApplications() {
  const navigate = useNavigate();

  const [jobs, setJobs] = useState([]);
  const [selectedJobId, setSelectedJobId] = useState("");
  const [submissions, setSubmissions] = useState([]);
  const [selectedCandidates, setSelectedCandidates] = useState([]);
  const [minExperience, setMinExperience] = useState(0);
  const [selectedSkill, setSelectedSkill] = useState("");
  const [sortBy, setSortBy] = useState("atsScore");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [statusUpdate, setStatusUpdate] = useState({});

  const [selectedStatus, setSelectedStatus] = useState("all"); // New state for filtering by status

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Fetch all jobs
  const fetchJobs = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/jobs`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch job postings.");
      }

      const data = await response.json();
      setJobs(data);
    } catch (err) {
      setError(err.message || "Error fetching jobs.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Fetch submissions for the selected job based on status filter
  const fetchSubmissions = async (jobId, status = "") => {
    try {
      setLoading(true);
      const token = localStorage.getItem("token");
      if (!token) {
        setError("Unauthorized. Please login again.");
        return;
      }

      let url = `${API_BASE}/jobs/${jobId}/submissions`;
      if (status && status !== "all") {
        url += `?status=${status}`;
      }

      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch job submissions.");
      }
      const data = await response.json();

      if (data && data.data) {
        const transformed = data.data.map((item) => {
          const formResponses = item.form_data || {};
          return {
            email: item.email,
            name: item.username,
            atsScore: item.ats_score,
            submissionDate: item.created_at,
            resumeUrl: item.resume_url,
            id: item.id,
            responses: { ...formResponses, Q_Skills: item.skills || [] },
            status: item.status || "applied", // Default status as 'applied'
          };
        });
        setSubmissions(transformed);
      } else {
        setSubmissions([]);
      }
    } catch (err) {
      setError(err.message || "Error fetching submissions.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchJobs();
  }, []);

  const handleJobChange = (e) => {
    const jobId = e.target.value;
    setSelectedJobId(jobId);
    setSelectedCandidates([]);
    if (jobId) {
      fetchSubmissions(jobId, selectedStatus);
    } else {
      setSubmissions([]);
    }
  };

  const handleCandidateClick = (candidate) => {
    const isAlreadySelected = selectedCandidates.some((c) => c.email === candidate.email);
    if (!isAlreadySelected) {
      setSelectedCandidates([...selectedCandidates, candidate]);
    }
  };

  const handleStatusUpdate = async (sub_id, status) => {
    try {
      const response = await fetch(`${API_BASE}/jobs/submissions/${sub_id}/status`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ status }),
      });

      if (response.ok) {
        setStatusUpdate((prevState) => ({
          ...prevState,
          [sub_id]: status,
        }));
      } else {
        throw new Error("Failed to update status.");
      }
    } catch (err) {
      setError(err.message || "Error updating status.");
      console.error(err);
    }
  };

  // Sorting logic
  const sortedCandidates = submissions.sort((a, b) => {
    if (sortBy === "atsScore") {
      return b.atsScore - a.atsScore;
    } else {
      return new Date(a.submissionDate) - new Date(b.submissionDate);
    }
  });

  return (
    <div className="relative min-h-screen flex flex-col justify-between bg-gray-100 p-6">
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      <div className="flex-grow flex justify-center items-center">
        <div className="bg-white shadow-xl rounded-xl p-8 w-full max-w-3xl">
          <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-6">
            Job Applications
          </h2>

          {error && (
            <p className="text-red-600 font-semibold text-center mb-4">
              {error}
            </p>
          )}
          {loading && (
            <p className="text-blue-600 font-semibold text-center mb-4">
              Loading...
            </p>
          )}

          {/* Job Selection */}
          <div className="mb-6">
            <label className="block font-medium mb-2 text-lg">Select Job ID:</label>
            <select
              value={selectedJobId}
              onChange={handleJobChange}
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
            >
              <option value="">-- Select Job ID --</option>
              {jobs && jobs.length > 0 ? (
                jobs.map((job) => (
                  <option key={job.job_id} value={job.job_id}>
                    {job.job_id}
                  </option>
                ))
              ) : (
                <option value="">No jobs available</option>
              )}
            </select>
          </div>

          {/* Status Filter (Button Group) */}
          {selectedJobId && (
            <div className="mb-6">
              <label className="block font-medium mb-2 text-lg">Filter by Status:</label>
              <div className="flex space-x-4">
                <button
                  onClick={() => {
                    setSelectedStatus("all");
                    fetchSubmissions(selectedJobId, ""); // fetch all
                  }}
                  className={`px-4 py-2 rounded-lg text-lg font-semibold ${
                    selectedStatus === "all" ? "bg-blue-500 text-white" : "bg-gray-200"
                  }`}
                >
                  All
                </button>
                <button
                  onClick={() => {
                    setSelectedStatus("applied");
                    fetchSubmissions(selectedJobId, "applied");
                  }}
                  className={`px-4 py-2 rounded-lg text-lg font-semibold ${
                    selectedStatus === "applied" ? "bg-blue-500 text-white" : "bg-gray-200"
                  }`}
                >
                  Applied
                </button>
                <button
                  onClick={() => {
                    setSelectedStatus("under_review");
                    fetchSubmissions(selectedJobId, "under_review");
                  }}
                  className={`px-4 py-2 rounded-lg text-lg font-semibold ${
                    selectedStatus === "under_review" ? "bg-blue-500 text-white" : "bg-gray-200"
                  }`}
                >
                  Under Review
                </button>
                <button
                  onClick={() => {
                    setSelectedStatus("rejected");
                    fetchSubmissions(selectedJobId, "rejected");
                  }}
                  className={`px-4 py-2 rounded-lg text-lg font-semibold ${
                    selectedStatus === "rejected" ? "bg-blue-500 text-white" : "bg-gray-200"
                  }`}
                >
                  Rejected
                </button>
              </div>
            </div>
          )}

          {/* Candidates List */}
          {selectedJobId && sortedCandidates.length > 0 ? (
            <div className="mb-6">
              <h3 className="text-2xl font-semibold mb-4">Candidates:</h3>
              {sortedCandidates.map((candidate) => {
                const isSelected = selectedCandidates.some((c) => c.email === candidate.email);
                return (
                  <div
                    key={candidate.email}
                    onClick={() => handleCandidateClick(candidate)}
                    className={`p-4 border rounded-lg cursor-pointer transition mb-4 ${
                      isSelected
                        ? "bg-blue-100 border-blue-500"
                        : "bg-gray-200 border-gray-400"
                    }`}
                  >
                    <h4 className="text-xl font-semibold">{candidate.name}</h4>
                    <p className="text-gray-600">{candidate.email}</p>
                    <p className="text-sm text-gray-500">Status: {statusUpdate[candidate.id] || candidate.status}</p>
                    {isSelected && (
                      <>
                        <p>
                          <strong>ATS Score:</strong> {candidate.atsScore}
                        </p>
                        <p>
                          <strong>Submission Date:</strong>{" "}
                          {new Date(candidate.submissionDate).toLocaleDateString()}
                        </p>
                        {candidate.responses &&
                          Object.entries(candidate.responses).map(([key, value]) => {
                            if (Array.isArray(value)) {
                              return (
                                <p key={key}>
                                  <strong>{key}:</strong> {value.join(", ")}
                                </p>
                              );
                            } else {
                              return (
                                <p key={key}>
                                  <strong>{key}:</strong> {value}
                                </p>
                              );
                            }
                          })}
                        <div className="flex mt-4">
                          <button
                            onClick={() => handleStatusUpdate(candidate.id, "under_review")}
                            className="px-4 py-2 bg-green-500 text-white rounded-lg mr-2"
                          >
                            Review
                          </button>
                          <button
                            onClick={() => handleStatusUpdate(candidate.id, "rejected")}
                            className="px-4 py-2 bg-red-500 text-white rounded-lg"
                          >
                            Reject
                          </button>
                        </div>
                      </>
                    )}
                  </div>
                );
              })}
            </div>
          ) : selectedJobId && !loading ? (
            <p>No candidates match the selected criteria.</p>
          ) : null}
        </div>
      </div>
    </div>
  );
}

export default JobApplications;