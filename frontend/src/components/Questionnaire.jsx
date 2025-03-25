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
      { id: "Q_Resume", text: "Upload your resume", type: "file", options: [] },
    ],
  };

  const [formData, setFormData] = useState(initialState);
  const [formSubmitted, setFormSubmitted] = useState(false);

  // Handle input changes
  const handleChange = (e, index, isOption = false, optionIndex = null) => {
    const updatedQuestions = [...formData.Questions];

    if (isOption) {
      updatedQuestions[index].options[optionIndex] = e.target.value;
    } else {
      updatedQuestions[index].text = e.target.value;
    }

    setFormData({ ...formData, Questions: updatedQuestions });
  };

  // Handle form submission
  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("Questionnaire Created:", formData);
    alert("Questionnaire Created Successfully!");
    setFormSubmitted(true);
  };

  return (
    <div className="relative min-h-screen flex flex-col justify-between bg-gray-100">
      {/* ✅ Dashboard Button */}
      <button
        onClick={() => navigate("/dashboard")}
        className="absolute top-4 left-4 px-4 py-2 text-lg bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition"
      >
        ⬅️ Dashboard
      </button>

      {/* ✅ Styled Form Card (Shorter height for spacing) */}
      <div className="flex-grow flex justify-center items-center mt-12 mb-20">
        <div className="bg-white shadow-2xl rounded-xl p-8 w-full max-w-xl h-[500px] overflow-y-auto relative">
          <h2 className="text-4xl font-bold text-center text-gray-800 tracking-wide mb-8">
            Create a Job Questionnaire
          </h2>

          <form onSubmit={handleSubmit}>
            {/* Form ID Input */}
            <div className="mb-6">
              <label className="block font-medium mb-2 text-lg font-bold">
                Form ID:
              </label>
              <input
                type="text"
                name="FormID"
                value={formData.FormID}
                onChange={(e) =>
                  setFormData({ ...formData, FormID: e.target.value })
                }
                required
                className="w-full p-3 border border-gray-500 rounded-lg text-lg focus:ring-2 focus:ring-green-400"
                placeholder="Enter Form ID"
              />
            </div>

            {/* Dynamic Question Inputs */}
            {formData.Questions.map((question, index) => (
              <div key={question.id} className="mb-6">
                <label className="block font-medium mb-2 text-lg font-bold">
                  {question.text}
                </label>

                {/* Text Input */}
                {question.type === "text" && (
                  <input
                    type="text"
                    className="w-full p-3 border border-gray-500 rounded-lg text-lg focus:ring-2 focus:ring-green-400"
                  />
                )}

                {/* Radio Buttons */}
                {question.type === "radio" &&
                  question.options.map((option, optionIndex) => (
                    <label key={optionIndex} className="flex items-center space-x-3">
                      <input
                        type="radio"
                        name={question.id}
                        value={option}
                        className="w-5 h-5 text-green-500 focus:ring-green-400"
                      />
                      <span className="text-lg">{option}</span>
                    </label>
                  ))}

                {/* Checkbox Options */}
                {question.type === "checkbox" &&
                  question.options.map((option, optionIndex) => (
                    <label key={optionIndex} className="flex items-center space-x-3">
                      <input
                        type="checkbox"
                        value={option}
                        className="w-5 h-5 text-green-500 focus:ring-green-400"
                      />
                      <span className="text-lg">{option}</span>
                    </label>
                  ))}

                {/* File Upload */}
                {question.type === "file" && (
                  <input
                    type="file"
                    className="w-full p-3 border border-gray-500 rounded-lg text-lg"
                  />
                )}
              </div>
            ))}

            {/* Submit Button */}
            <button
              type="submit"
              className="w-full py-3 bg-green-500 text-white text-lg rounded-lg hover:bg-green-600 transition"
            >
              Submit Questionnaire
            </button>

            {formSubmitted && (
              <button
                type="button"
                onClick={() => setFormData(initialState)}
                className="w-full py-3 mt-3 bg-gray-500 text-white text-lg rounded-lg hover:bg-gray-600 transition"
              >
                New Questionnaire
              </button>
            )}
          </form>
        </div>
      </div>

      {/* ✅ Footer Section (Now properly separated) */}
      <footer className="bg-gradient-to-r from-gray-900 to-black text-white py-8 px-12">
        <div className="container mx-auto grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* About Section */}
          <div>
            <h2 className="text-2xl font-bold">HireEasy</h2>
            <p className="mt-3 text-gray-400">
              Helping recruiters streamline hiring with job applications, 
              questionnaires, and data-driven insights.
            </p>
          </div>

          {/* Office Information */}
          <div>
            <h3 className="text-xl font-bold mb-3">Office</h3>
            <p>456 Talent Hub,</p>
            <p>San Francisco, USA</p>
            <p>Email: support@hireeasy.com</p>
            <p>Phone: +1 987-654-3210</p>
          </div>

          {/* Useful Links */}
          <div>
            <h3 className="text-xl font-bold mb-3">Links</h3>
            <ul className="space-y-2">
              <li><a href="#" className="hover:text-gray-300 transition">Dashboard</a></li>
              <li><a href="#" className="hover:text-gray-300 transition">Post a Job</a></li>
            </ul>
          </div>
        </div>
        <div className="text-center text-gray-400 mt-8">
          HireEasy © {new Date().getFullYear()} - All Rights Reserved
        </div>
      </footer>
    </div>
  );
}

export default Questionnaire;
