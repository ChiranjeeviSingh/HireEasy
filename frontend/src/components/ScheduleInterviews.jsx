import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

function ScheduleInterviews() {
  const navigate = useNavigate();
  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  const [jobs, setJobs] = useState([]);
  const [selectedJobId, setSelectedJobId] = useState("");
  const [candidates, setCandidates] = useState([]);
  const [availability, setAvailability] = useState([]);
  const [showPopupFor, setShowPopupFor] = useState(null);
  const [error, setError] = useState("");
  const [interviews, setInterviews] = useState({});

  const token = localStorage.getItem("token");
  console.log("🪪 Auth Token:", token);

  useEffect(() => {
    const fetchJobs = async () => {
      if (!token) {
        setError("Unauthorized. Please login again.");
        return;
      }

      try {
        const res = await fetch(`${API_BASE}/jobs`, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (res.status === 401) throw new Error("Unauthorized request. Token may be invalid or expired.");

        const data = await res.json();
        setJobs(data);
      } catch (err) {
        setError(err.message || "Failed to load jobs.");
      }
    };

    fetchJobs();
  }, [API_BASE, token]);

  const fetchCandidates = async (jobId) => {
    if (!token) {
      setError("Unauthorized. Please login again.");
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/jobs/${jobId}/submissions`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (res.status === 401) throw new Error("Unauthorized request while fetching candidates.");

      const data = await res.json();
      if (data && data.data) {
        setCandidates(data.data);
      } else {
        setCandidates([]);
      }
    } catch (err) {
      setError(err.message || "Failed to load candidates.");
    }
  };

  const fetchAvailability = async () => {
    if (!token) {
      setError("Unauthorized. Please login again.");
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/availability`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (res.status === 401) throw new Error("Unauthorized request while fetching availability.");

      const data = await res.json();
      if (Array.isArray(data)) {
        setAvailability(data);
      } else {
        setAvailability([]);
        console.error("Availability API returned unexpected format:", data);
      }
    } catch (err) {
      setError(err.message || "Failed to fetch availability.");
      setAvailability([]);
    }
  };

  const handleScheduleInterview = async (candidate, slot) => {
    if (!token) {
      alert("Unauthorized. Please login again.");
      return;
    }

    try {
      const body = {
        job_id: candidate.job_id,
        job_submission_id: candidate.id,
        interviewer_user_id: slot.user_id,
        availability_id: slot.id,
      };

      const res = await fetch(`${API_BASE}/interviews`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (res.status === 401) throw new Error("Unauthorized request while scheduling interview.");
      if (!res.ok) throw new Error("Failed to schedule interview.");

      const result = await res.json();

      alert("Interview scheduled successfully!");
      setShowPopupFor(null);
      setInterviews((prev) => ({
        ...prev,
        [candidate.id]: {
          slot,
          interview_id: result.id || result.data?.id || null,
        },
      }));
    } catch (err) {
      alert(err.message);
    }
  };

  const handleDeleteInterview = async (candidateId) => {
    const interview = interviews[candidateId];
    if (!interview || !interview.interview_id) return;

    try {
      const res = await fetch(`${API_BASE}/interviews/${interview.interview_id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) throw new Error("Failed to delete interview");
      alert("Interview deleted successfully");
      setInterviews((prev) => {
        const updated = { ...prev };
        delete updated[candidateId];
        return updated;
      });
    } catch (err) {
      alert(err.message);
    }
  };

  return (
    <div className="p-6 min-h-screen bg-gradient-to-br from-blue-50 to-blue-100">
      <div className="max-w-6xl mx-auto">
        <button
          onClick={() => navigate("/dashboard")}
          className="mb-6 px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-800"
        >
          ⬅️ Back to Dashboard
        </button>

        <h1 className="text-4xl font-bold text-center text-blue-800 mb-8">
          Schedule Interviews
        </h1>

        {error && <p className="text-red-600 text-center font-medium">{error}</p>}

        <div className="max-w-md mx-auto mb-10">
          <label className="block mb-2 font-medium text-lg">Select Job:</label>
          <select
            value={selectedJobId}
            onChange={(e) => {
              const val = e.target.value;
              setSelectedJobId(val);
              fetchCandidates(val);
              setShowPopupFor(null);
            }}
            className="w-full p-3 border border-gray-300 rounded-lg shadow focus:ring-2 focus:ring-blue-400"
          >
            <option value="">-- Select Job --</option>
            {jobs.map((job) => (
              <option key={job.job_id} value={job.job_id}>
                {job.job_id}
              </option>
            ))}
          </select>
        </div>

        {selectedJobId && candidates.length > 0 && (
          <div className="grid gap-6 md:grid-cols-2">
            {candidates.map((c) => (
              <div
                key={c.id}
                className="p-6 bg-white rounded-lg shadow hover:shadow-lg transition"
              >
                <h3 className="text-xl font-semibold text-gray-800 mb-1">{c.username}</h3>
                <p className="text-gray-600 mb-1">📧 <strong>Email:</strong> {c.email}</p>
                <p className="text-gray-600 mb-1">📄 <strong>Status:</strong> {c.status}</p>
                <p className="text-gray-600 mb-1">📊 <strong>ATS Score:</strong> {c.ats_score}</p>
                <p className="text-gray-600 mb-1">
                  🗓️ <strong>Submitted:</strong> {new Date(c.created_at).toLocaleDateString()}
                </p>
                <p className="text-gray-600 mb-1">
                  📎 <strong>Resume:</strong> <a href={c.resume_url} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline">Download</a>
                </p>
                {c.skills && c.skills.length > 0 && (
                  <p className="text-gray-600 mb-1">🛠️ <strong>Skills:</strong> {c.skills.join(", ")}</p>
                )}

                {interviews[c.id] ? (
                  <div className="mt-4 text-green-700">
                    ✅ <strong>Interview Scheduled:</strong> {interviews[c.id].slot.username} – {new Date(interviews[c.id].slot.date).toLocaleDateString()} | {interviews[c.id].slot.from_time} - {interviews[c.id].slot.to_time}
                    <button
                      onClick={() => handleDeleteInterview(c.id)}
                      className="ml-4 px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
                    >
                      Cancel
                    </button>
                  </div>
                ) : (
                  <button
                    onClick={async () => {
                      await fetchAvailability();
                      setShowPopupFor(c.id);
                    }}
                    className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                  >
                    Schedule Interview
                  </button>
                )}

                {showPopupFor === c.id && !interviews[c.id] && (
                  <div className="mt-4 bg-gray-50 border border-gray-300 p-4 rounded-lg">
                    <h4 className="font-semibold mb-3 text-gray-800">Select a Time Slot:</h4>
                    {availability.length === 0 ? (
                      <p className="text-gray-500">No slots available.</p>
                    ) : (
                      <div className="grid gap-2">
                        {availability.map((slot) => (
                          <button
                            key={slot.id}
                            onClick={() => handleScheduleInterview(c, slot)}
                            className="w-full p-2 bg-green-600 text-white rounded hover:bg-green-700"
                          >
                            {slot.username} – {new Date(slot.date).toLocaleDateString()} | {slot.from_time} - {slot.to_time}
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}

        {selectedJobId && candidates.length === 0 && (
          <p className="text-center text-gray-600 mt-10">No candidates available for this job.</p>
        )}
      </div>
    </div>
  );
}

export default ScheduleInterviews;