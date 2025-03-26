import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

export function Apply() {
  const location = useLocation();

  const [jobDetails, setJobDetails] = useState(null);
  const [formDetails, setFormDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [formData, setFormData] = useState({});

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  // Extract jobid and formid from the query parameters
  const queryParams = new URLSearchParams(location.search);
  const jobId = queryParams.get("jobid");
  const formId = queryParams.get("formid");

  // Fetch job details and form details using jobId and formId
  useEffect(() => {
    const fetchJobAndFormDetails = async () => {
      try {
        setLoading(true);

        const token = localStorage.getItem("token");
        if (!token) {
          throw new Error("Authorization token is missing.");
        }

        // Fetch Job Details
        const jobResponse = await fetch(`${API_BASE}/jobs/${jobId}`, {
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        });
        if (!jobResponse.ok) {
          throw new Error("Failed to fetch job details");
        }
        const jobData = await jobResponse.json();
        setJobDetails(jobData);

        // Fetch Form Details
        const formResponse = await fetch(`${API_BASE}/forms/templates/${formId}`, {
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        });
        if (!formResponse.ok) {
          throw new Error("Failed to fetch form details");
        }
        const formData = await formResponse.json();
        setFormDetails(formData);

        setLoading(false);
      } catch (error) {
        setError(error.message);
        setLoading(false);
      }
    };

    if (jobId && formId) {
      fetchJobAndFormDetails();
    } else {
      setError("Job ID or Form ID is missing in the URL.");
    }
  }, [jobId, formId]);

  // Handle form input change
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

  // Handle form submit
  const handleSubmit = (e) => {
    e.preventDefault();

    // Combine the job details and the form data to send to the backend
    const submissionData = {
      jobDetails,
      formDetails,
      formData,
    };

    // Log the submission data to the console
    console.log("Form Submission Data:", submissionData);

    // For now, we are just logging the data in the frontend.
    // In the future, you can send this to the backend using a POST request.
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
      {jobDetails && (
        <div style={{ maxWidth: "600px", margin: "auto", textAlign: "left" }}>
          <h3>Job Details</h3>
          <p>
            <strong>Job ID:</strong> {jobDetails.job_id}
          </p>
          <p>
            <strong>Job Title:</strong> {jobDetails.job_title}
          </p>
          <p>
            <strong>Description:</strong> {jobDetails.job_description}
          </p>
          <p>
            <strong>Skills Required:</strong> {jobDetails.skills_required.join(", ")}
          </p>
        </div>
      )}

      {/* Display Form Questions */}
      {formDetails && (
        <div style={{ maxWidth: "600px", margin: "auto", textAlign: "left", marginTop: "20px" }}>
          <h3>Job Questionnaire</h3>
          <form onSubmit={handleSubmit}>
            {formDetails.fields.map((question) => (
              <div key={question.question_id} style={{ marginBottom: "10px" }}>
                <p>
                  <strong>{question.question_text}</strong>
                </p>

                {question.question_type === "text" && (
                  <input
                    type="text"
                    name={question.question_id}
                    onChange={handleInputChange}
                    value={formData[question.question_id] || ""}
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
                      />{" "}
                      {opt}
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
                      />{" "}
                      {opt}
                    </label>
                  ))}
                {question.question_type === "file" && (
                  <input
                    type="file"
                    name={question.question_id}
                    onChange={handleInputChange}
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


// http://localhost:3000/apply?jobid=ceqc&formid=456 example URL which works.