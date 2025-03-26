import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function Questionnaire() {
  const navigate = useNavigate(); // Hook for navigation

  // Initial state: Questions & their options
  const initialState = {
    FormID: "",
    Questions: [
      {
        id: "Q_Gender",
        text: "What is your gender?",
        type: "radio",
        options: ["Male", "Female", "Other"],
      },
      { id: "Q_Education", text: "Education Level", type: "text", options: [] },
      {
        id: "Q_Skills",
        text: "Which programming languages do you know?",
        type: "checkbox",
        options: ["JavaScript", "Python", "C++"],
      },
      {
        id: "Q_Experience",
        text: "Work Experience (Years)",
        type: "text",
        options: [],
      },
      { id: "Q_Resume", text: "Upload your resume", type: "file", options: [] },
      {
        id: "Q6",
        text: "What is your expected salary?",
        type: "text",
        options: [],
      },
      { id: "Q7", text: "", type: "text", options: [] },
      { id: "Q8", text: "", type: "text", options: [] },
      { id: "Q9", text: "", type: "text", options: [] },
      { id: "Q10", text: "", type: "text", options: [] },
    ],
  };

  const [formData, setFormData] = useState(initialState);
  const [formSubmitted, setFormSubmitted] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Handle input changes (Question text or options)
  const handleChange = (e, index, isOption = false, optionIndex = null) => {
    const updatedQuestions = [...formData.Questions];

    if (isOption) {
      // Update specific option for a question
      updatedQuestions[index].options[optionIndex] = e.target.value;
    } else {
      updatedQuestions[index].text = e.target.value;
    }

    setFormData({ ...formData, Questions: updatedQuestions });
  };

  // Add new option for multiple-choice questions
  const addOption = (index) => {
    const updatedQuestions = [...formData.Questions];
    updatedQuestions[index].options.push("");
    setFormData({ ...formData, Questions: updatedQuestions });
  };

  // Handle form submission (integration with backend)
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    const token = localStorage.getItem("token"); // Get JWT token

    // Map frontend fields to backend structure (correct field names as in the backend model)
    const questionnaireData = {
      form_template_id: formData.FormID.trim(), // map to form_template_id
      user_id: 1, // This will be the logged-in user's ID; for now, you can fetch it from localStorage or other sources
      fields: formData.Questions.map((question) => ({
        question_id: question.id, // unique question id
        question_text: question.text.trim(), // question text
        question_type: question.type, // question type (radio, checkbox, text, etc.)
        options: question.options.map((option) => option.trim()) // options for radio/checkbox questions
      })),
    };

    // Log the data being sent to the backend
    console.log("Mapped Questionnaire Data:", questionnaireData);

    // Check if all required fields are valid
    if (
      !questionnaireData.form_template_id ||
      questionnaireData.fields.length === 0
    ) {
      setError("Please fill in all required fields.");
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`${API_BASE}/forms/templates`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`, // Attach JWT token for authorization
        },
        body: JSON.stringify(questionnaireData), // Send the questionnaire data
      });

      const data = await response.json();

      console.log("Response Data:", data); // Log the response from the backend

      if (!response.ok) {
        throw new Error(data.msg || "Failed to create questionnaire.");
      }

      setSuccess(`Questionnaire with Form ID "${data.form_template_id}" created successfully!`);
      setFormSubmitted(true);

    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Reset form for a new questionnaire
  const handleNewQuestionnaire = () => {
    setFormData(initialState);
    setFormSubmitted(false);
  };

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

      <h2>Create a Job Questionnaire</h2>

      {error && <p className="error-text" style={{ color: "red" }}>{error}</p>}
      {success && <p className="success-text" style={{ color: "green" }}>{success}</p>}

      <form
        onSubmit={handleSubmit}
        style={{ maxWidth: "600px", margin: "auto" }}
      >
        {/* Form ID Input */}
        <div style={{ marginBottom: "10px" }}>
          <label>Form ID (Required): </label>
          <input
            type="text"
            name="FormID"
            value={formData.FormID}
            onChange={(e) =>
              setFormData({ ...formData, FormID: e.target.value })
            }
            required
            style={{ width: "100%", padding: "8px", fontSize: "16px" }}
          />
        </div>

        {/* Dynamic Question Inputs */}
        {formData.Questions.map((question, index) => (
          <div key={question.id} style={{ marginBottom: "10px" }}>
            <label>{`Question ${index + 1}:`}</label>
            <input
              type="text"
              value={question.text}
              onChange={(e) => handleChange(e, index)}
              style={{ width: "100%", padding: "8px", fontSize: "16px" }}
            />

            {/* Multiple choice (radio/checkbox) options */}
            {(question.type === "radio" || question.type === "checkbox") && (
              <div>
                {question.options.map((option, optionIndex) => (
                  <input
                    key={optionIndex}
                    type="text"
                    value={option}
                    onChange={(e) => handleChange(e, index, true, optionIndex)}
                    placeholder={`Option ${optionIndex + 1}`}
                    style={{ width: "80%", padding: "5px", marginTop: "5px" }}
                  />
                ))}
                <button
                  type="button"
                  onClick={() => addOption(index)}
                  style={{ marginLeft: "5px", padding: "5px 10px" }}
                >
                  ➕ Add Option
                </button>
              </div>
            )}

            {/* File Upload Notice */}
            {question.type === "file" && (
              <p style={{ fontSize: "14px", color: "gray" }}>
                This question requires candidates to upload a file.
              </p>
            )}
          </div>
        ))}

        {/* Submit and New Buttons */}
        <button
          type="submit"
          style={{ marginTop: "10px", cursor: "pointer", padding: "10px 15px" }}
        >
          Submit Questionnaire
        </button>

        {formSubmitted && (
          <button
            type="button"
            onClick={handleNewQuestionnaire}
            style={{
              marginLeft: "10px",
              cursor: "pointer",
              padding: "10px 15px",
            }}
          >
            New
          </button>
        )}
      </form>
    </div>
  );
}

export default Questionnaire;