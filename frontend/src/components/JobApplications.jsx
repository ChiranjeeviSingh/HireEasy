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
        submissionDate: "2025-03-20",
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
        submissionDate: "2025-03-22",
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
        submissionDate: "2025-03-18",
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
  const [selectedSkill, setSelectedSkill] = useState("");
  const [sortBy, setSortBy] = useState("atsScore");

  // Handle job selection
  const handleJobChange = (e) => {
    setSelectedJobId(e.target.value);
    setSelectedCandidates([]);
  };

  // Handle candidate selection
  const handleCandidateClick = (candidate) => {
    const isAlreadySelected = selectedCandidates.some((c) => c.email === candidate.email);

    if (isAlreadySelected) {
      setSelectedCandidates(selectedCandidates.filter((c) => c.email !== candidate.email));
    } else {
      setSelectedCandidates([...selectedCandidates, candidate]);
    }
  };

  // Handle experience filter selection
  const handleExperienceChange = (e) => {
    setMinExperience(Number(e.target.value));
    setSelectedCandidates([]);
  };

  // Handle skills filter selection
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
          const experienceYears = parseInt(candidate.responses.Q_Experience.replace(/\D/g, ""), 10);

          const matchesExperience = experienceYears >= minExperience;
          const matchesSkill =
            selectedSkill === "" ||
            candidate.responses.Q_Skills.some((skill) => skill.toLowerCase().includes(selectedSkill));

          return matchesExperience && matchesSkill;
        })
      : [];

  // Sort candidates by ATS Score or Submission Date
  const sortedCandidates = filteredCandidates.sort((a, b) => {
    if (sortBy === "atsScore") {
      return b.atsScore - a.atsScore;
    } else {
      return new Date(a.submissionDate) - new Date(b.submissionDate);
    }
  });

  return (
    <div className="relative min-h-screen flex flex-col justify-between bg-gray-100 p-6">
      {/* ✅ Dashboard Button */}
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

          {/* Job ID Selection */}
          <div className="mb-6">
            <label className="block font-medium mb-2 text-lg">Select Job ID:</label>
            <select
              value={selectedJobId}
              onChange={handleJobChange}
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
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
            <div className="grid grid-cols-2 gap-4 mb-6">
              <div>
                <label className="block font-medium mb-2 text-lg">Filter by Experience (years):</label>
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
                <label className="block font-medium mb-2 text-lg">Filter by Skill:</label>
                <input
                  type="text"
                  placeholder="Enter skill..."
                  value={selectedSkill}
                  onChange={handleSkillChange}
                  className="w-full p-3 border border-gray-300 rounded-lg text-lg"
                />
              </div>
            </div>
          )}

          {/* Sorting Toggle */}
          {selectedJobId && (
            <div className="mb-6">
              <label className="block font-medium mb-2 text-lg">Sort by:</label>
              <button
                onClick={handleSortChange}
                className="w-full py-3 bg-blue-500 text-white text-lg rounded-lg hover:bg-blue-600 transition"
              >
                {sortBy === "atsScore" ? "Sort by Submission Date" : "Sort by ATS Score"}
              </button>
            </div>
          )}

          {/* Show Candidates */}
          {selectedJobId && sortedCandidates.length > 0 ? (
            <div className="mb-6">
              <h3 className="text-2xl font-semibold mb-4">Candidates:</h3>
              {sortedCandidates.map((candidate) => (
                <div
                  key={candidate.email}
                  onClick={() => handleCandidateClick(candidate)}
                  className={`p-4 border rounded-lg cursor-pointer transition ${
                    selectedCandidates.some((c) => c.email === candidate.email)
                      ? "bg-blue-100 border-blue-500"
                      : "bg-gray-200 border-gray-400"
                  }`}
                >
                  <h4 className="text-xl font-semibold">{candidate.name}</h4>
                  <p className="text-gray-600">{candidate.email}</p>
                  {selectedCandidates.some((c) => c.email === candidate.email) && (
                    <>
                      <p><strong>ATS Score:</strong> {candidate.atsScore}</p>
                      <p><strong>Submission Date:</strong> {candidate.submissionDate}</p>
                      {Object.entries(candidate.responses).map(([key, value]) => (
                        <p key={key}>
                          <strong>{key.replace("Q_", "").replace("_", " ")}:</strong> {Array.isArray(value) ? value.join(", ") : value}
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
      </div>
    </div>
  );
}

export default JobApplications;