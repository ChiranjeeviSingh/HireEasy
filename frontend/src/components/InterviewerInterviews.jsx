import React, { useEffect, useState } from "react";

function InterviewerInterviews() {
  const [interviews, setInterviews] = useState([]);
  const [statusMap, setStatusMap] = useState({});
  const [feedbackMap, setFeedbackMap] = useState({});
  const [submittedFeedbacks, setSubmittedFeedbacks] = useState({});
  const [error, setError] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";
  const token = localStorage.getItem("token");

  useEffect(() => {
    fetch(`${API_BASE}/interviews`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error("Failed to fetch interviews.");
        const data = await res.json();
        setInterviews(data);
      })
      .catch((err) => setError(err.message));
  }, [API_BASE, token]);

  const handleStatusChange = (id, newStatus) => {
    setStatusMap((prev) => ({ ...prev, [id]: newStatus }));
  };

  const handleFeedbackChange = (id, field, value) => {
    setFeedbackMap((prev) => ({
      ...prev,
      [id]: { ...prev[id], [field]: value },
    }));
  };

  const handleFeedbackSubmit = async (interviewId) => {
    const feedbackData = feedbackMap[interviewId];
    if (!feedbackData || !feedbackData.verdict || !feedbackData.feedback) {
      alert("Please fill both verdict and feedback.");
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/interviews/${interviewId}/feedback`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(feedbackData),
      });

      if (!res.ok) throw new Error("Failed to submit feedback");

      setSubmittedFeedbacks((prev) => ({ ...prev, [interviewId]: true }));
      alert("Feedback submitted successfully.");
    } catch (err) {
      alert(err.message);
    }
  };

  return (
    <div className="p-8 min-h-screen bg-gradient-to-br from-gray-100 to-blue-50">
      {/* Back button */}
      <div className="mb-6">
        <a
          href="/interviewer-dashboard"
          className="inline-block bg-gray-700 text-white px-4 py-2 rounded hover:bg-gray-800 transition"
        >
          ⬅ Back to Dashboard
        </a>
      </div>

      <h1 className="text-4xl font-bold text-center text-blue-900 mb-8">
        Scheduled Interviews
      </h1>

      {error && (
        <p className="text-center text-red-600 font-medium mb-6">{error}</p>
      )}

      {interviews.length === 0 && !error ? (
        <p className="text-center text-gray-600 text-lg">No interviews scheduled.</p>
      ) : (
        <div className="grid gap-6 max-w-4xl mx-auto">
          {interviews.map((iv) => (
            <div
              key={iv.id}
              className="bg-white shadow-md rounded-lg p-6 border border-gray-200"
            >
              <h3 className="text-xl font-semibold mb-2 text-gray-800">
                Interview ID: {iv.id}
              </h3>
              <p className="text-gray-700">
                <strong>Job ID:</strong> {iv.job_id}
              </p>
              <p className="text-gray-700">
                <strong>Candidate Submission ID:</strong> {iv.job_submission_id}
              </p>
              <p className="text-gray-700">
                <strong>Status:</strong>{" "}
                {statusMap[iv.id] || iv.status}
              </p>
              <p className="text-gray-700">
                <strong>Scheduled On:</strong>{" "}
                {new Date(iv.created_at).toLocaleString()}
              </p>

              {/* Feedback Section */}
              {!submittedFeedbacks[iv.id] ? (
                <div className="mt-4 bg-gray-50 p-4 rounded-lg border">
                  <label className="block mb-2 font-medium text-gray-800">
                    Feedback
                  </label>
                  <textarea
                    rows="3"
                    value={feedbackMap[iv.id]?.feedback || ""}
                    onChange={(e) =>
                      handleFeedbackChange(iv.id, "feedback", e.target.value)
                    }
                    className="w-full p-2 border border-gray-300 rounded"
                    placeholder="Enter your feedback here"
                  ></textarea>

                  <label className="block mt-4 font-medium text-gray-800">
                    Verdict
                  </label>
                  <select
                    value={feedbackMap[iv.id]?.verdict || ""}
                    onChange={(e) =>
                      handleFeedbackChange(iv.id, "verdict", e.target.value)
                    }
                    className="w-full p-2 border border-gray-300 rounded"
                  >
                    <option value="">-- Select Verdict --</option>
                    <option value="passed">Passed</option>
                    <option value="failed">Failed</option>
                  </select>

                  <button
                    onClick={() => handleFeedbackSubmit(iv.id)}
                    className="mt-4 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                  >
                    Submit Feedback
                  </button>
                </div>
              ) : (
                <p className="mt-4 text-green-600 font-semibold">
                  ✅ Feedback submitted.
                </p>
              )}

              {/* Finalize / Reject */}
              <div className="flex gap-4 mt-4">
                <button
                  onClick={() => handleStatusChange(iv.id, "Candidate Finalized")}
                  className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                >
                  Finalize
                </button>
                <button
                  onClick={() => handleStatusChange(iv.id, "Candidate Rejected")}
                  className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
                >
                  Reject
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default InterviewerInterviews;