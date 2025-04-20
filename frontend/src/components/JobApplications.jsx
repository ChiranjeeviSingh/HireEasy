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

  // Adjust to your actual API location or environment variable
  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // 1) Fetch all jobs for the dropdown (requires Auth if your route is protected)
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

  // 2) Fetch submissions for the selected job
  const fetchSubmissions = async (jobId) => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/jobs/${jobId}/submissions`);
      if (!response.ok) {
        throw new Error("Failed to fetch job submissions.");
      }
      const data = await response.json();

      if (data && data.data) {
        // Transform data to the structure we want for filtering, sorting, etc.
        // We'll parse item.form_data for experience, skills, etc.
        const transformed = data.data.map((item) => {
          // If form_data is returned as an object, we can use it directly:
          const formResponses = item.form_data || {};

          // For example, if your form_data might contain { Q_Experience: "3 years", Q_Skills: ["React","Node"] }, etc.
          return {
            email: item.email,
            name: item.username,
            atsScore: item.ats_score,
            submissionDate: item.created_at,
            resumeUrl: item.resume_url,
            // We'll store the entire set of Q_ fields under `responses`
            responses: {
              ...formResponses,
              // Also insert `skills` if you store them separately in the DB
              Q_Skills: item.skills || [],
            },
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

  // On mount, fetch all jobs
  useEffect(() => {
    fetchJobs();
  }, []);

  // Handle job selection from dropdown
  const handleJobChange = (e) => {
    const jobId = e.target.value;
    setSelectedJobId(jobId);
    setSelectedCandidates([]);
    if (jobId) {
      // Fetch submissions for the newly selected job
      fetchSubmissions(jobId);
    } else {
      setSubmissions([]);
    }
  };

  // Toggle candidate expand/collapse
  const handleCandidateClick = (candidate) => {
    const isAlreadySelected = selectedCandidates.some(
      (c) => c.email === candidate.email
    );
    if (isAlreadySelected) {
      setSelectedCandidates(
        selectedCandidates.filter((c) => c.email !== candidate.email)
      );
    } else {
      setSelectedCandidates([...selectedCandidates, candidate]);
    }
  };

  // Filter: minimum years of experience
  const handleExperienceChange = (e) => {
    setMinExperience(Number(e.target.value));
    setSelectedCandidates([]);
  };

  // Filter: skill
  const handleSkillChange = (e) => {
    setSelectedSkill(e.target.value.toLowerCase());
    setSelectedCandidates([]);
  };

  // Sort toggle (ATS score vs. submission date)
  const handleSortChange = () => {
    setSortBy((prevSortBy) =>
      prevSortBy === "atsScore" ? "submissionDate" : "atsScore"
    );
  };

  // Filter logic
  const filteredCandidates = submissions.filter((candidate) => {
    // If Q_Experience is stored as "5 years", parse out the number
    let experienceYears = 0;
    if (candidate.responses?.Q_Experience) {
      experienceYears = parseInt(
        candidate.responses.Q_Experience.replace(/\D/g, ""),
        10
      );
      if (isNaN(experienceYears)) {
        experienceYears = 0;
      }
    }

    const matchesExperience = experienceYears >= minExperience;

    // Skill filter: compare input skill to Q_Skills array (if any)
    const candidateSkills = candidate.responses?.Q_Skills || [];
    const matchesSkill =
      !selectedSkill ||
      candidateSkills.some((skill) =>
        skill.toLowerCase().includes(selectedSkill)
      );

    return matchesExperience && matchesSkill;
  });

  // Sorting logic
  const sortedCandidates = filteredCandidates.sort((a, b) => {
    if (sortBy === "atsScore") {
      return b.atsScore - a.atsScore;
    } else {
      // Ascending by submission date
      return new Date(a.submissionDate) - new Date(b.submissionDate);
    }
  });

  return (
    <div className="relative min-h-screen flex flex-col justify-between bg-gray-100 p-6">
      {/* Dashboard Button */}
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
            <label className="block font-medium mb-2 text-lg">
              Select Job ID:
            </label>
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

          {/* Filter and Sort Controls */}
          {selectedJobId && submissions.length > 0 && (
            <>
              <div className="grid grid-cols-2 gap-4 mb-6">
                <div>
                  <label className="block font-medium mb-2 text-lg">
                    Filter by Experience (years):
                  </label>
                  <select
                    value={minExperience}
                    onChange={handleExperienceChange}
                    className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
                  >
                    <option value="0">All Candidates</option>
                    <option value="2">Greater than 2 years</option>
                    <option value="3">Greater than 3 years</option>
                    <option value="5">Greater than 5 years</option>
                  </select>
                </div>

                <div>
                  <label className="block font-medium mb-2 text-lg">
                    Filter by Skill:
                  </label>
                  <input
                    type="text"
                    placeholder="Enter skill..."
                    value={selectedSkill}
                    onChange={handleSkillChange}
                    className="w-full p-3 border border-gray-300 rounded-lg text-lg"
                  />
                </div>
              </div>

              {/* Sorting Toggle */}
              <div className="mb-6">
                <label className="block font-medium mb-2 text-lg">
                  Sort by:
                </label>
                <button
                  onClick={handleSortChange}
                  className="w-full py-3 bg-blue-500 text-white text-lg rounded-lg hover:bg-blue-600 transition"
                >
                  {sortBy === "atsScore"
                    ? "Sort by Submission Date"
                    : "Sort by ATS Score"}
                </button>
              </div>
            </>
          )}

          {/* Candidates List */}
          {selectedJobId && sortedCandidates.length > 0 ? (
            <div className="mb-6">
              <h3 className="text-2xl font-semibold mb-4">Candidates:</h3>
              {sortedCandidates.map((candidate) => {
                const isSelected = selectedCandidates.some(
                  (c) => c.email === candidate.email
                );
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
                    {isSelected && (
                      <>
                        <p>
                          <strong>ATS Score:</strong> {candidate.atsScore}
                        </p>
                        <p>
                          <strong>Submission Date:</strong>{" "}
                          {new Date(candidate.submissionDate).toLocaleDateString()}
                        </p>
                        {/* Show the form responses (if any) */}
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