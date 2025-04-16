import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

export function InterviewerProfile() {
  const navigate = useNavigate();
  const [profile, setProfile] = useState({
    name: "",
    email: "",
    phone: "",
  });

  useEffect(() => {
    // TODO: Load existing profile data from your API
    // Example:
    // fetch("/api/interviewer/profile")
    //   .then((res) => res.json())
    //   .then((data) => setProfile(data));
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setProfile((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    // TODO: Send updated profile to your API
    // Example:
    // fetch("/api/interviewer/profile", {
    //   method: "PUT",
    //   headers: { "Content-Type": "application/json" },
    //   body: JSON.stringify(profile),
    // }).then(() => navigate("/dashboard"));
    navigate("/interviewer-dashboard"); // remove when real API call is in place
  };

  return (
    <div style={containerStyle}>
      <h2>Edit Interviewer Profile</h2>
      <form onSubmit={handleSubmit} style={formStyle}>
        <label style={labelStyle}>
          Name
          <input
            type="text"
            name="name"
            value={profile.name}
            onChange={handleChange}
            style={inputStyle}
            required
          />
        </label>

        <label style={labelStyle}>
          Email
          <input
            type="email"
            name="email"
            value={profile.email}
            onChange={handleChange}
            style={inputStyle}
            required
          />
        </label>

        <label style={labelStyle}>
          Phone
          <input
            type="tel"
            name="phone"
            value={profile.phone}
            onChange={handleChange}
            style={inputStyle}
          />
        </label>

        <div style={buttonContainerStyle}>
          <button
            type="button"
            onClick={() => navigate("/interviewer-dashboard")}
            style={cancelButtonStyle}
          >
            Cancel
          </button>
          <button type="submit" style={buttonStyle}>
            Save
          </button>
        </div>
      </form>
    </div>
  );
}

// Styles
const containerStyle = {
  maxWidth: "500px",
  margin: "50px auto",
  padding: "20px",
  textAlign: "left",
};

const formStyle = {
  display: "flex",
  flexDirection: "column",
  gap: "20px",
};

const labelStyle = {
  display: "flex",
  flexDirection: "column",
  fontSize: "16px",
  color: "#333",
};

const inputStyle = {
  marginTop: "8px",
  padding: "12px",
  fontSize: "14px",
  border: "1px solid #ccc",
  borderRadius: "4px",
};

const buttonContainerStyle = {
  display: "flex",
  justifyContent: "flex-end",
  gap: "10px",
  marginTop: "10px",
};

const buttonStyle = {
  padding: "12px 24px",
  fontSize: "16px",
  backgroundColor: "#007bff",
  color: "white",
  border: "none",
  borderRadius: "8px",
  cursor: "pointer",
  transition: "background 0.3s",
};

const cancelButtonStyle = {
  ...buttonStyle,
  backgroundColor: "#6c757d",
};

export default InterviewerProfile;
