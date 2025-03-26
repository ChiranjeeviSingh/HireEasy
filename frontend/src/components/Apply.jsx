import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

export function Apply() {
  const location = useLocation();

  const [applicationDetails, setApplicationDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [formData, setFormData] = useState({}); // Store user input

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Extract form_uuid from URL
  const queryParams = new URLSearchParams(location.search);
  const formUUID = queryParams.get("formid"); // Now using form_uuid

  useEffect(() => {
    const fetchApplicationDetails = async () => {
      try {
        setLoading(true);

        const token = localStorage.getItem("token");
        if (!token) {
          throw new Error("Authorization token is missing.");
        }

        // Fetch Job and Form Details using form_uuid
        const response = await fetch(`${API_BASE}/forms/${formUUID}`, {
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch job and form details.");
        }

        const data = await response.json();
        setApplicationDetails(data);
        setLoading(false);
      } catch (error) {
        setError(error.message);
        setLoading(false);
      }
    };

    if (formUUID) {
      fetchApplicationDetails();
    } else {
      setError("Form ID (form_uuid) is missing in the URL.");
    }
  }, [formUUID]);

  // Handle form input changes
  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData((prevData) => {
      if (type === "checkbox") {
        const updatedValues = prevData[name] || [];
        if (checked) {
          updatedValues.push(value);
        } else {
          const index = updatedValues.indexOf(value);
          if (index !== -1) updatedValues.splice(index, 1);
        }
        return { ...prevData, [name]: updatedValues };
      }
      return { ...prevData, [name]: value };
    });
  };

  // Handle form submit (console log for now)
  const handleSubmit = (e) => {
    e.preventDefault();

    const submissionData = {
      application_id: applicationDetails?.form_uuid,
      job_id: applicationDetails?.job?.job_id,
      form_template_id: applicationDetails?.form_template?.form_template_id,
      responses: formData,
    };

    console.log("Form Submission Data:", submissionData);
    alert("Application Submitted! (Check Console for Details)");
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div style={{ textAlign: "center", marginTop: "20px", position: "relative" }}>
      <h2>Job Application</h2>

      {/* Display Job Details */}
      {applicationDetails?.job && (
        <div style={{ maxWidth: "600px", margin: "auto", textAlign: "left" }}>
          <h3>Job Details</h3>
          <p><strong>Job ID:</strong> {applicationDetails.job.job_id}</p>
          <p><strong>Job Title:</strong> {applicationDetails.job.job_title}</p>
          <p><strong>Description:</strong> {applicationDetails.job.job_description}</p>
          <p><strong>Skills Required:</strong> {applicationDetails.job.skills_required.join(", ")}</p>
        </div>
      )}

      {/* Display Form Questions */}
      {applicationDetails?.form_template && (
        <div style={{ maxWidth: "600px", margin: "auto", textAlign: "left", marginTop: "20px" }}>
          <h3>Job Questionnaire</h3>
          <form onSubmit={handleSubmit}>
            {applicationDetails.form_template.fields.map((question) => (
              <div key={question.question_id} style={{ marginBottom: "10px" }}>
                <p><strong>{question.question_text}</strong></p>

                {/* Render inputs based on question type */}
                {question.question_type === "text" && (
                  <input
                    type="text"
                    name={question.question_id}
                    onChange={handleInputChange}
                    value={formData[question.question_id] || ""}
                    required
                  />
                )}
                {question.question_type === "radio" &&
                  question.options.map((opt) => (
                    <label key={opt}>
                      <input
                        type="radio"
                        name={question.question_id}
                        value={opt}
                        onChange={handleInputChange}
                        checked={formData[question.question_id] === opt}
                        required
                      /> {opt}
                    </label>
                  ))}
                {question.question_type === "checkbox" &&
                  question.options.map((opt) => (
                    <label key={opt}>
                      <input
                        type="checkbox"
                        name={question.question_id}
                        value={opt}
                        onChange={handleInputChange}
                        checked={(formData[question.question_id] || []).includes(opt)}
                      /> {opt}
                    </label>
                  ))}
                {question.question_type === "file" && (
                  <input
                    type="file"
                    name={question.question_id}
                    onChange={handleInputChange}
                    required
                  />
                )}
              </div>
            ))}

            <button type="submit" style={{ padding: "10px 15px", cursor: "pointer" }}>
              Submit Application
            </button>
          </form>
        </div>
      )}
    </div>
  );
}

export default Apply;