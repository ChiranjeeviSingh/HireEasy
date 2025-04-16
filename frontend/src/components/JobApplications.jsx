import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function JobApplications() {
  const navigate = useNavigate();

  // Mock job applications data (JobID -> Candidates -> Responses with ATS and Submission Date)
  const jobApplications = {
    JOB123: [
      {
        email: "john@example.com",
        name: "John Doe",
        atsScore: 85,
        submissionDate: "2025-03-20", // Ensure the date format is consistent
        responses: {
          Q_Gender: "Male",
          Q_Skills: ["JavaScript", "Python"],
          Q_Experience: "5 years",
          Q_Resume: "john_resume.pdf",
          Q6: "120,000 per year",
        },
      },
      {
        email: "sarah@example.com",
        name: "Sarah Johnson",
        atsScore: 92,
        submissionDate: "2025-03-22", // Ensure the date format is consistent
        responses: {
          Q_Gender: "Female",
          Q_Skills: ["Python", "SQL"],
          Q_Experience: "3 years",
          Q_Resume: "sarah_resume.pdf",
          Q6: "110,000 per year",
        },
      },
    ],
    JOB456: [
      {
        email: "mike@example.com",
        name: "Mike Smith",
        atsScore: 78,
        submissionDate: "2025-03-18", // Ensure the date format is consistent
        responses: {
          Q_Gender: "Male",
          Q_Skills: ["Python", "C++"],
          Q_Experience: "7 years",
          Q_Resume: "mike_resume.pdf",
          Q6: "150,000 per year",
        },
      },
    ],
  };

  const [selectedJobId, setSelectedJobId] = useState("");
  const [selectedCandidates, setSelectedCandidates] = useState([]);
  const [minExperience, setMinExperience] = useState(0);
  const [selectedSkill, setSelectedSkill] = useState(""); // Skill filter
  const [sortBy, setSortBy] = useState("atsScore"); // Toggle for sorting by ATS or Submission Date

  // Handle job selection
  const handleJobChange = (e) => {
    setSelectedJobId(e.target.value);
    setSelectedCandidates([]);
  };

  // Handle candidate selection
  const handleCandidateClick = (candidate) => {
    const isAlreadySelected = selectedCandidates.some(
      (c) => c.email === candidate.email
    );

    if (isAlreadySelected) {
      // Deselect candidate
      setSelectedCandidates(
        selectedCandidates.filter((c) => c.email !== candidate.email)
      );
    } else {
      // Select candidate
      setSelectedCandidates([...selectedCandidates, candidate]);
    }
  };

  // Handle experience filter selection
  const handleExperienceChange = (e) => {
    setMinExperience(Number(e.target.value));
    setSelectedCandidates([]);
  };

  // Handle skills filter selection (Partial match and case insensitive)
  const handleSkillChange = (e) => {
    setSelectedSkill(e.target.value.toLowerCase());
    setSelectedCandidates([]);
  };

  // Handle sorting toggle (ATS score or Submission Date)
  const handleSortChange = () => {
    setSortBy((prevSortBy) => (prevSortBy === "atsScore" ? "submissionDate" : "atsScore"));
  };

  // Filter candidates based on selected criteria
  const filteredCandidates =
    selectedJobId && jobApplications[selectedJobId]
      ? jobApplications[selectedJobId].filter((candidate) => {
          const experienceYears = parseInt(
            candidate.responses.Q_Experience.replace(/\D/g, ""),
            10
          );

          const matchesExperience = experienceYears >= minExperience;
          const matchesSkill =
            selectedSkill === "" ||
            candidate.responses.Q_Skills.some((skill) =>
              skill.toLowerCase().includes(selectedSkill)
            );

          return matchesExperience && matchesSkill;
        })
      : [];

  // Sort candidates by ATS Score or Submission Date
  const sortedCandidates = filteredCandidates.sort((a, b) => {
    if (sortBy === "atsScore") {
      return b.atsScore - a.atsScore; // Sort by ATS Score (descending)
    } else {
      // Sort by Submission Date (ascending) --> Reversed (earliest submissions first)
      const dateA = new Date(a.submissionDate); // Ensure it's in Date format
      const dateB = new Date(b.submissionDate); // Ensure it's in Date format
      return dateA - dateB; // Older submissions first
    }
  });

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

      <h2>Job Applications</h2>

      {/* Job ID Selection */}
      <div style={{ marginBottom: "20px" }}>
        <label>Select Job ID: </label>
        <select
          value={selectedJobId}
          onChange={handleJobChange}
          style={{ width: "100%", padding: "8px", fontSize: "16px" }}
        >
          <option value="">-- Select Job ID --</option>
          {Object.keys(jobApplications).map((jobId) => (
            <option key={jobId} value={jobId}>
              {jobId}
            </option>
          ))}
        </select>
      </div>

      {/* Filter and Sort Controls */}
      {selectedJobId && (
        <div style={{ marginBottom: "20px", display: "flex", justifyContent: "space-around" }}>
          <div style={{ flex: 1 }}>
            <label>Filter by Experience (years): </label>
            <select
              value={minExperience}
              onChange={handleExperienceChange}
              style={{ width: "100%", padding: "8px", fontSize: "16px" }}
            >
              <option value="0">All Candidates</option>
              <option value="2">Greater than 2 years</option>
              <option value="3">Greater than 3 years</option>
              <option value="5">Greater than 5 years</option>
            </select>
          </div>

          <div style={{ flex: 1 }}>
            <label>Filter by Skill: </label>
            <input
              type="text"
              placeholder="Enter skill..."
              value={selectedSkill}
              onChange={handleSkillChange}
              style={{ width: "100%", padding: "8px", fontSize: "16px" }}
            />
          </div>
        </div>
      )}

      {/* Sorting Toggle */}
      {selectedJobId && (
        <div style={{ marginBottom: "20px" }}>
          <label>Sort by: </label>
          <button
            onClick={handleSortChange}
            style={{
              padding: "8px",
              backgroundColor: "#007bff",
              color: "white",
              border: "none",
              borderRadius: "5px",
              cursor: "pointer",
            }}
          >
            {sortBy === "atsScore" ? "Sort by Submission Date" : "Sort by ATS Score"}
          </button>
        </div>
      )}

      {/* Show Candidates if a job is selected */}
      {selectedJobId && sortedCandidates.length > 0 ? (
        <div style={{ marginBottom: "20px" }}>
          <h3>Candidates:</h3>
          {sortedCandidates.map((candidate) => (
            <div
              key={candidate.email}
              onClick={() => handleCandidateClick(candidate)}
              style={{
                display: "block",
                width: "90%",
                padding: "15px",
                margin: "10px auto",
                backgroundColor: "#007bff",
                color: "white",
                borderRadius: "5px",
                cursor: "pointer",
                textAlign: "left",
                height: selectedCandidates.some((c) => c.email === candidate.email)
                  ? "auto"
                  : "60px",
                transition: "all 0.3s ease",
              }}
            >
              <h4>{candidate.name}</h4>
              <p>{candidate.email}</p>

              {selectedCandidates.some((c) => c.email === candidate.email) && (
                <>
                  <p><strong>ATS Score:</strong> {candidate.atsScore}</p>
                  <p><strong>Submission Date:</strong> {candidate.submissionDate}</p>
                  {Object.entries(candidate.responses).map(([key, value]) => (
                    <p key={key}>
                      <strong>{key.replace("Q_", "").replace("_", " ")}:</strong>{" "}
                      {Array.isArray(value) ? value.join(", ") : value}
                    </p>
                  ))}
                </>
              )}
            </div>
          ))}
        </div>
      ) : selectedJobId ? (
        <p>No candidates match the selected criteria.</p>
      ) : null}
    </div>
  );
}

export default JobApplications;