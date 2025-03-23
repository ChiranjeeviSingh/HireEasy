import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function JobApplications() {
  const navigate = useNavigate();

  // Mock job applications data (JobID -> Candidates -> Responses)
  const jobApplications = {
    JOB123: [
      {
        email: "john@example.com",
        name: "John Doe",
        responses: {
          Q_Gender: "Male",
          Q_Education: "Master's in Computer Science",
          Q_Skills: ["JavaScript", "Python"],
          Q_Experience: "5 years",
          Q_Resume: "john_resume.pdf",
          Q6: "120,000 per year",
        },
      },
      {
        email: "sarah@example.com",
        name: "Sarah Johnson",
        responses: {
          Q_Gender: "Female",
          Q_Education: "Bachelor's in Data Science",
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
        responses: {
          Q_Gender: "Male",
          Q_Education: "PhD in AI",
          Q_Skills: ["Python", "C++"],
          Q_Experience: "7 years",
          Q_Resume: "mike_resume.pdf",
          Q6: "150,000 per year",
        },
      },
    ],
  };

  const [selectedJobId, setSelectedJobId] = useState("");
  const [selectedCandidate, setSelectedCandidate] = useState(null);
  const [minExperience, setMinExperience] = useState(0);
  const [selectedEducation, setSelectedEducation] = useState(""); // Education filter
  const [selectedSkill, setSelectedSkill] = useState(""); // Skill filter

  // Handle job selection
  const handleJobChange = (e) => {
    setSelectedJobId(e.target.value);
    setSelectedCandidate(null);
  };

  // Handle candidate selection
  const handleCandidateChange = (email) => {
    const candidate = jobApplications[selectedJobId].find(
      (c) => c.email === email
    );
    setSelectedCandidate(candidate);
  };

  // Handle experience filter selection
  const handleExperienceChange = (e) => {
    setMinExperience(Number(e.target.value));
    setSelectedCandidate(null);
  };

  // Handle education filter selection
  const handleEducationChange = (e) => {
    setSelectedEducation(e.target.value);
    setSelectedCandidate(null);
  };

  // Handle skills filter selection
  const handleSkillChange = (e) => {
    setSelectedSkill(e.target.value);
    setSelectedCandidate(null);
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
          const matchesEducation =
            selectedEducation === "" ||
            candidate.responses.Q_Education.includes(selectedEducation);
          const matchesSkill =
            selectedSkill === "" ||
            candidate.responses.Q_Skills.includes(selectedSkill);

          return matchesExperience && matchesEducation && matchesSkill;
        })
      : [];

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

      {/* Experience Filter Selection */}
      {selectedJobId && (
        <div style={{ marginBottom: "20px" }}>
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
      )}

      {/* Education Filter Selection */}
      {selectedJobId && (
        <div style={{ marginBottom: "20px" }}>
          <label>Filter by Education: </label>
          <input
            type="text"
            placeholder="Enter education keyword..."
            value={selectedEducation}
            onChange={handleEducationChange}
            style={{ width: "100%", padding: "8px", fontSize: "16px" }}
          />
        </div>
      )}

      {/* Skills Filter Selection */}
      {selectedJobId && (
        <div style={{ marginBottom: "20px" }}>
          <label>Filter by Skill: </label>
          <input
            type="text"
            placeholder="Enter skill..."
            value={selectedSkill}
            onChange={handleSkillChange}
            style={{ width: "100%", padding: "8px", fontSize: "16px" }}
          />
        </div>
      )}

      {/* Show Candidates if a job is selected */}
      {selectedJobId && filteredCandidates.length > 0 ? (
        <div style={{ marginBottom: "20px" }}>
          <h3>Candidates:</h3>
          {filteredCandidates.map((candidate) => (
            <button
              key={candidate.email}
              onClick={() => handleCandidateChange(candidate.email)}
              style={{
                display: "block",
                width: "100%",
                padding: "8px",
                margin: "5px 0",
                backgroundColor: "#007bff",
                color: "white",
                border: "none",
                borderRadius: "5px",
                cursor: "pointer",
                textAlign: "left",
              }}
            >
              {candidate.name} ({candidate.email})
            </button>
          ))}
        </div>
      ) : selectedJobId ? (
        <p>No candidates match the selected criteria.</p>
      ) : null}

      {/* Show Candidate Responses if selected */}
      {selectedCandidate && (
        <div
          style={{
            marginTop: "20px",
            textAlign: "left",
            maxWidth: "600px",
            margin: "auto",
          }}
        >
          <h3>Candidate Details:</h3>
          <p>
            <strong>Name:</strong> {selectedCandidate.name}
          </p>
          <p>
            <strong>Email:</strong> {selectedCandidate.email}
          </p>

          <h3>Submitted Responses:</h3>
          {Object.entries(selectedCandidate.responses).map(([key, value]) => (
            <p key={key}>
              <strong>{key.replace("Q_", "").replace("_", " ")}:</strong>{" "}
              {Array.isArray(value) ? value.join(", ") : value}
            </p>
          ))}
        </div>
      )}
    </div>
  );
}

export default JobApplications;
