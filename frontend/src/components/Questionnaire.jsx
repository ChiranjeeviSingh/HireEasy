import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function Questionnaire() {
  const navigate = useNavigate();

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
      {
        id: "Q_Resume",
        text: "Upload your resume",
        type: "file",
        options: [],
      },
      {
        id: "Q6",
        text: "What is your expected salary?",
        type: "text",
        options: [],
      },
    ],
  };

  const [formData, setFormData] = useState(initialState);
  const [additionalQuestions, setAdditionalQuestions] = useState([]);
  const [formSubmitted, setFormSubmitted] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(false);

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Handle input changes for predefined questions
  const handleChange = (e, index, isOption = false, optionIndex = null) => {
    const updatedQuestions = [...formData.Questions];

    if (isOption) {
      updatedQuestions[index].options[optionIndex] = e.target.value;
    } else {
      updatedQuestions[index].text = e.target.value;
    }

    setFormData({ ...formData, Questions: updatedQuestions });
  };

  // Handle input changes for additional questions (Q7-Q10)
  const handleAdditionalChange = (e, index, isOption = false, optionIndex = null) => {
    const updatedQuestions = [...additionalQuestions];

    if (isOption) {
      updatedQuestions[index].options[optionIndex] = e.target.value;
    } else {
      updatedQuestions[index].text = e.target.value;
    }

    setAdditionalQuestions(updatedQuestions);
  };

  // Add new option for multiple-choice questions
  const addOption = (index, isAdditional = false) => {
    if (isAdditional) {
      const updatedQuestions = [...additionalQuestions];
      updatedQuestions[index].options.push("");
      setAdditionalQuestions(updatedQuestions);
    } else {
      const updatedQuestions = [...formData.Questions];
      updatedQuestions[index].options.push("");
      setFormData({ ...formData, Questions: updatedQuestions });
    }
  };

  // Add new question (Q7-Q10)
  const addQuestion = () => {
    if (additionalQuestions.length < 4) {
      const newQuestion = {
        id: `Q${7 + additionalQuestions.length}`,
        text: "",
        type: "",
        options: [],
      };
      setAdditionalQuestions([...additionalQuestions, newQuestion]);
    }
  };

  // Set question type for a given additional question
  const setQuestionType = (index, type) => {
    const updatedQuestions = [...additionalQuestions];
    updatedQuestions[index].type = type;

    if (type === "radio" || type === "checkbox") {
      // Provide at least one default option
      updatedQuestions[index].options = ["Default1"];
    } else {
      updatedQuestions[index].options = [];
    }

    setAdditionalQuestions(updatedQuestions);
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    const token = localStorage.getItem("token");

    // Merge additional questions into main form data
    const finalQuestions = [...formData.Questions, ...additionalQuestions];

    // Map frontend fields to backend structure
    const questionnaireData = {
      form_template_id: formData.FormID.trim(),
      user_id: 1,
      fields: finalQuestions.map((question) => ({
        question_id: question.id,
        question_text: question.text.trim(),
        question_type: question.type,
        options: question.options.map((option) => option.trim()),
      })),
    };

    console.log("Mapped Questionnaire Data:", questionnaireData);

    if (!questionnaireData.form_template_id || questionnaireData.fields.length === 0) {
      setError("Please fill in all required fields.");
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`${API_BASE}/forms/templates`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(questionnaireData),
      });

      const data = await response.json();
      console.log("Response Data:", data);

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

  return (
    <div className="flex flex-col items-center min-h-screen bg-gray-100">
      {/* Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      <h2 className="text-3xl font-bold mt-12">Create a Job Questionnaire</h2>

      {error && <p className="text-red-500 mt-3">{error}</p>}
      {success && <p className="text-green-500 mt-3">{success}</p>}

      <form onSubmit={handleSubmit} className="bg-white p-6 mt-6 shadow-md rounded-lg w-full max-w-xl">
        {/* Form ID Input */}
        <div className="mb-4">
          <label className="block font-bold">Form ID:</label>
          <input
            type="text"
            name="FormID"
            value={formData.FormID}
            onChange={(e) => setFormData({ ...formData, FormID: e.target.value })}
            required
            className="w-full p-2 border border-gray-300 rounded-lg"
          />
        </div>

        {/* Predefined Questions */}
        {formData.Questions.map((question, index) => (
          <div key={question.id} className="mb-4">
            <label className="block font-bold">Question Text:</label>
            <input
              type="text"
              value={question.text}
              onChange={(e) => handleChange(e, index)}
              className="w-full p-2 border border-gray-300 rounded-lg"
            />

            {/* Show file input placeholder if question.type === 'file' */}
            {question.type === "file" && (
              <div className="mt-2">
                <input
                  type="file"
                  disabled
                  className="w-full p-2 border border-gray-300 rounded-lg cursor-not-allowed"
                />
                <p className="text-gray-500 text-sm mt-1">
                  (File upload will be shown to candidates)
                </p>
              </div>
            )}

            {/* If it's radio/checkbox, show the options */}
            {(question.type === "radio" || question.type === "checkbox") && (
              <div className="mt-2">
                {question.options.map((option, optionIndex) => (
                  <input
                    key={optionIndex}
                    type="text"
                    value={option}
                    onChange={(e) => handleChange(e, index, true, optionIndex)}
                    className="w-full p-2 border border-gray-300 rounded-lg mt-1"
                  />
                ))}
                <button
                  type="button"
                  onClick={() => addOption(index)}
                  className="text-blue-500 mt-1"
                >
                  ➕ Add Option
                </button>
              </div>
            )}
          </div>
        ))}

        {/* Additional Questions (Q7-Q10) */}
        {additionalQuestions.map((question, index) => (
          <div key={question.id} className="mb-4">
            {/* Question Text */}
            <label className="block font-bold">New Question Text:</label>
            <input
              type="text"
              value={question.text}
              onChange={(e) => handleAdditionalChange(e, index)}
              className="w-full p-2 border border-gray-300 rounded-lg"
              placeholder="Enter question text"
            />

            {/* Question Type Select */}
            <label className="block font-bold mt-2">Select Question Type:</label>
            <select
              value={question.type}
              onChange={(e) => setQuestionType(index, e.target.value)}
              className="w-full p-2 border border-gray-300 rounded-lg"
            >
              <option value="">Select Question Type</option>
              <option value="text">Text</option>
              <option value="radio">Radio</option>
              <option value="checkbox">Checkbox</option>
              <option value="file">File</option>
            </select>

            {/* If we want to show a disabled file input for newly added “file” questions */}
            {question.type === "file" && (
              <div className="mt-2">
                <input
                  type="file"
                  disabled
                  className="w-full p-2 border border-gray-300 rounded-lg cursor-not-allowed"
                />
                <p className="text-gray-500 text-sm mt-1">
                  (File upload will be shown to candidates)
                </p>
              </div>
            )}

            {/* Options for radio/checkbox */}
            {(question.type === "radio" || question.type === "checkbox") && (
              <div className="mt-2">
                {question.options.map((option, optionIndex) => (
                  <input
                    key={optionIndex}
                    type="text"
                    value={option}
                    onChange={(e) => handleAdditionalChange(e, index, true, optionIndex)}
                    className="w-full p-2 border border-gray-300 rounded-lg mt-1"
                    placeholder={`Option ${optionIndex + 1}`}
                  />
                ))}
                <button
                  type="button"
                  onClick={() => addOption(index, true)}
                  className="text-blue-500 mt-1"
                >
                  ➕ Add Option
                </button>
              </div>
            )}
          </div>
        ))}

        {/* Button to Add Additional Questions */}
        {additionalQuestions.length < 4 && (
          <button
            type="button"
            onClick={addQuestion}
            className="text-green-500 mt-3"
          >
            ➕ Add Question
          </button>
        )}

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full bg-green-500 text-white py-2 mt-4 rounded-lg"
          disabled={loading}
        >
          {loading ? "Submitting..." : "Submit Questionnaire"}
        </button>
      </form>
    </div>
  );
}

export default Questionnaire;