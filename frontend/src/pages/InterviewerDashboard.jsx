import React from "react";
import { useNavigate } from "react-router-dom";

export function InterviewerDashboard() {
  const navigate = useNavigate();

  return (
    <div style={{ textAlign: "center", marginTop: "50px" }}>
      <h2>Welcome, Interviewer!</h2>

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
        <button onClick={() => navigate("/interviewer-profile")} style={buttonStyle}>
          Edit Profile
        </button>

        <button onClick={() => navigate("/interviewer-calendar")} style={buttonStyle}>
          Calendar
        </button>

        <button onClick={() => navigate("/interviewer-interviews")} style={buttonStyle}>
          Interview
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

// Hover effect (using a simple inline workaround)
const originalButtonStyle = { ...buttonStyle };
const hoverStyle = { backgroundColor: "#0056b3" };

// Add event listeners for hover effect
Object.assign(buttonStyle, {
  onMouseEnter: (e) => {
    Object.assign(e.target.style, hoverStyle);
  },
  onMouseLeave: (e) => {
    Object.assign(e.target.style, {
      backgroundColor: originalButtonStyle.backgroundColor,
    });
  },
});

export default InterviewerDashboard;
