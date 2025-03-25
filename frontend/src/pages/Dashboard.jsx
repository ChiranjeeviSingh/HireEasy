import React from "react";
import { useNavigate } from "react-router-dom";

export function Dashboard() {
  const navigate = useNavigate();

  return (
    <div style={{ textAlign: "center", marginTop: "50px" }}>
      <h2>Welcome, HR!</h2>

      {/* Container for buttons */}
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          gap: "20px",
          marginTop: "30px",
          flexWrap: "wrap", // Ensures responsiveness
        }}
      >
        <button onClick={() => navigate("/job-posting")} style={buttonStyle}>
          Create Job Posting
        </button>

        <button onClick={() => navigate("/questionnaire")} style={buttonStyle}>
          Create Questionnaire
        </button>

        <button onClick={() => navigate("/share-job")} style={buttonStyle}>
          Share Job
        </button>

        <button
          onClick={() => navigate("/job-applications")}
          style={buttonStyle}
        >
          View Job Applications
        </button>

        <button onClick={() => navigate("/view-jobs")} style={buttonStyle}>
          View Jobs
        </button>
      </div>
    </div>
  );
}

// Button styling
const buttonStyle = {
  padding: "15px 25px",
  fontSize: "16px",
  backgroundColor: "#007bff",
  color: "white",
  border: "none",
  borderRadius: "8px",
  cursor: "pointer",
  transition: "background 0.3s",
  minWidth: "200px", // Ensures buttons are uniform in size
};

// Hover effect
buttonStyle[":hover"] = {
  backgroundColor: "#0056b3",
};

export default Dashboard;