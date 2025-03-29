import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

export function Apply() {
  const location = useLocation();
  
  // Holds the job + form template details once fetched
  const [applicationDetails, setApplicationDetails] = useState(null);
  
  // Basic loading/error states
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Candidate name & email – required by backend
  const [candidateName, setCandidateName] = useState("");
  const [candidateEmail, setCandidateEmail] = useState("");

  // We'll store all question answers here (except the resume file)
  // Keyed by question_id -> value(s)
  const [questionAnswers, setQuestionAnswers] = useState({});

  // We'll store the actual resume file here for the "file" question
  const [resumeFile, setResumeFile] = useState(null);

  // Base URL for your API
  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";

  // Extract form_uuid from the URL query param: ?formid=xxxx
  const queryParams = new URLSearchParams(location.search);
  const formUUID = queryParams.get("formid");

  // On component mount, fetch the form details (which includes the job info)
  useEffect(() => {
    const fetchApplicationDetails = async () => {
      try {
        setLoading(true);
        const token = localStorage.getItem("token");
        if (!token) {
          throw new Error("Authorization token is missing.");
        }

        // GET /api/forms/:form_uuid
        const response = await fetch(`${API_BASE}/forms/${formUUID}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch job and form details.");
        }

        const data = await response.json();
        setApplicationDetails(data);
        setLoading(false);
      } catch (err) {
        setError(err.message);
        setLoading(false);
      }
    };

    if (formUUID) {
      fetchApplicationDetails();
    } else {
      setError("Form ID (form_uuid) is missing in the URL.");
      setLoading(false);
    }
  }, [formUUID, API_BASE]);

  // Handle input changes for each question
  const handleInputChange = (e) => {
    const { name, value, type, checked, files } = e.target;

    // If the question is a file input (Resume)
    if (type === "file") {
      // The backend only wants one 'resume' file
      if (files && files[0]) {
        setResumeFile(files[0]); // store the File in state
      }
      // Optionally record the file name in questionAnswers
      setQuestionAnswers((prev) => ({
        ...prev,
        [name]: files[0]?.name || "",
      }));
      return;
    }

    // If it's a checkbox (could be multiple selections)
    if (type === "checkbox") {
      setQuestionAnswers((prev) => {
        const existing = prev[name] || [];
        if (checked) {
          return { ...prev, [name]: [...existing, value] };
        } else {
          // remove from array
          return { ...prev, [name]: existing.filter((v) => v !== value) };
        }
      });
      return;
    }

    // If it's radio or text
    setQuestionAnswers((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // Submitting the application
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!applicationDetails) return;

    // The job ID from the fetched data
    const jobID = applicationDetails.job?.job_id;
    if (!jobID) {
      alert("No valid job ID found. Cannot submit application.");
      return;
    }

    // Build the multipart FormData
    const submissionFormData = new FormData();

    // 1) job_id
    submissionFormData.append("job_id", jobID);

    // 2) username
    submissionFormData.append("username", candidateName);

    // 3) email
    submissionFormData.append("email", candidateEmail);

    // 4) form_data as JSON (all the question answers except the file itself)
    submissionFormData.append("form_data", JSON.stringify(questionAnswers));

    // 5) resume file (required by the backend)
    if (!resumeFile) {
      alert("Please attach a resume file before submitting.");
      return;
    }
    submissionFormData.append("resume", resumeFile);

    // Console logs for the data we're sending
    console.log("jobID:", jobID);
    console.log("candidateName:", candidateName);
    console.log("candidateEmail:", candidateEmail);
    console.log("questionAnswers:", questionAnswers);
    console.log("resumeFile:", resumeFile);
    console.log("FormData (raw object):", submissionFormData);

    try {
      const response = await fetch(`${API_BASE}/jobs/${jobID}/apply`, {
        method: "POST",
        body: submissionFormData,
      });

      // Attempt to parse JSON response
      let data;
      try {
        data = await response.json();
      } catch (parseError) {
        throw new Error("Invalid JSON response from server.");
      }

      // Console log what we expect or receive from the server
      console.log("Response from server:", data);

      if (!response.ok) {
        throw new Error(data.error || "Failed to submit application.");
      }

      alert(`Application submitted successfully! Submission ID: ${data.id}`);
    } catch (err) {
      alert(`Error: ${err.message}`);
    }
  };

  if (loading) {
    return <div className="text-center text-xl mt-10">Loading...</div>;
  }
  if (error) {
    return <div className="text-center text-red-500 mt-10">Error: {error}</div>;
  }

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-100">
      <div className="bg-white shadow-lg rounded-lg p-8 w-full max-w-3xl">
        <h2 className="text-3xl font-bold text-center mb-6 text-blue-700">
          Job Application
        </h2>

        {/* Show job details */}
        {applicationDetails?.job && (
          <div className="bg-gray-50 p-5 rounded-lg shadow-md mb-6">
            <h3 className="text-xl font-semibold mb-2 text-gray-800">
              Job Details
            </h3>
            <p className="text-gray-700">
              <strong>Job ID:</strong> {applicationDetails.job.job_id}
            </p>
            <p className="text-gray-700">
              <strong>Job Title:</strong> {applicationDetails.job.job_title}
            </p>
            <p className="text-gray-700">
              <strong>Description:</strong>{" "}
              {applicationDetails.job.job_description}
            </p>
            <p className="text-gray-700">
              <strong>Skills Required:</strong>{" "}
              {applicationDetails.job.skills_required.join(", ")}
            </p>
          </div>
        )}

        {/* Display the form questions */}
        {applicationDetails?.form_template && (
          <div className="bg-gray-50 p-5 rounded-lg shadow-md">
            <h3 className="text-xl font-semibold mb-4 text-gray-800">
              Job Questionnaire
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              {/* Candidate Name & Email (required by backend) */}
              <div>
                <label className="font-medium text-gray-900 block mb-1">
                  Full Name <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={candidateName}
                  onChange={(e) => setCandidateName(e.target.value)}
                  required
                  className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400"
                  placeholder="Enter your full name"
                />
              </div>

              <div>
                <label className="font-medium text-gray-900 block mb-1">
                  Email <span className="text-red-500">*</span>
                </label>
                <input
                  type="email"
                  value={candidateEmail}
                  onChange={(e) => setCandidateEmail(e.target.value)}
                  required
                  className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400"
                  placeholder="Enter your email"
                />
              </div>

              {applicationDetails.form_template.fields.map((question) => (
                <div key={question.question_id}>
                  <label className="font-medium text-gray-900 block mb-1">
                    {question.question_text}
                  </label>

                  {/* Text Input */}
                  {question.question_type === "text" && (
                    <input
                      type="text"
                      name={question.question_id}
                      onChange={handleInputChange}
                      value={questionAnswers[question.question_id] || ""}
                      required
                      className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400"
                    />
                  )}

                  {/* Radio Options */}
                  {question.question_type === "radio" &&
                    question.options.map((opt) => (
                      <label key={opt} className="flex items-center space-x-2 mt-2">
                        <input
                          type="radio"
                          name={question.question_id}
                          value={opt}
                          onChange={handleInputChange}
                          checked={questionAnswers[question.question_id] === opt}
                          required
                          className="w-5 h-5"
                        />
                        <span>{opt}</span>
                      </label>
                    ))}

                  {/* Checkbox Options */}
                  {question.question_type === "checkbox" &&
                    question.options.map((opt) => (
                      <label key={opt} className="flex items-center space-x-2 mt-2">
                        <input
                          type="checkbox"
                          name={question.question_id}
                          value={opt}
                          onChange={handleInputChange}
                          checked={
                            (questionAnswers[question.question_id] || []).includes(opt)
                          }
                          className="w-5 h-5"
                        />
                        <span>{opt}</span>
                      </label>
                    ))}

                  {/* File Upload (Resume) */}
                  {question.question_type === "file" && (
                    <input
                      type="file"
                      name={question.question_id}
                      onChange={handleInputChange}
                      required
                      className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-400"
                    />
                  )}
                </div>
              ))}

              {/* Submit Button */}
              <button
                type="submit"
                className="w-full py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition"
              >
                Submit Application
              </button>
            </form>
          </div>
        )}
      </div>
    </div>
  );
}

export default Apply;