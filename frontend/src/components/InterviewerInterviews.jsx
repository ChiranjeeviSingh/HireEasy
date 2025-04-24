import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

function InterviewerInterviews() {
  const navigate  = useNavigate();
  const API_BASE  = process.env.REACT_APP_API_URL || "http://localhost:8080/api";
  const token     = localStorage.getItem("token");

  /* ───────── state ───────── */
  const [interviews, setInterviews]  = useState([]);
  const [feedbacks,  setFeedbacks]   = useState({});
  const [error,      setError]       = useState("");
  const [loading,    setLoading]     = useState(true);
  const [success,    setSuccess]     = useState("");

  /* ───────── data load ───────── */
  useEffect(() => { loadData(); }, []);

  const loadData = async () => {
    if (!token) {
      setError("Unauthorized. Please login again.");
      setLoading(false);
      return;
    }

    try {
      setLoading(true);

      /* 1️⃣  fetch interviews */
      const itRes = await fetch(`${API_BASE}/interviews`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!itRes.ok) throw new Error("Failed to fetch interviews.");
      const interviewRows = await itRes.json();

      /* DEBUG */
      console.table(interviewRows, ["id", "job_id", "job_submission_id"]);

      /* feedback state from existing rows */
      const fbInit = {};
      interviewRows.forEach((row) => {
        if (row.feedback)
          fbInit[row.id] = { feedback: row.feedback, verdict: row.verdict };
      });
      setFeedbacks(fbInit);

      /* 2️⃣  fetch all availabilities once */
      const avRes = await fetch(`${API_BASE}/availability`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!avRes.ok) throw new Error("Failed to fetch availabilities.");
      const avList = await avRes.json();
      const avMap  = Object.fromEntries(avList.map((a) => [a.id, a]));

      /* 3️⃣  fetch submissions per job & build map */
      const jobIds = [...new Set(interviewRows.map((r) => r.job_id))];
      const submissionMap = {};

      for (const jobId of jobIds) {
        /* convert “123” or “J123” → “job-123” if your API needs that */
        const jobPath = jobId.startsWith("J")
          ? `job-${jobId.replace(/^\D+/, "")}`
          : jobId;

        const url = `${API_BASE}/jobs/${jobPath}/submissions`;
        console.log("↗️  GET", url);

        try {
          const subRes = await fetch(url, {
            headers: { Authorization: `Bearer ${token}` },
          });
          console.log("   ↩️", subRes.status, subRes.statusText);

          if (!subRes.ok) continue;

          /* ─── payload‑unwrapping fix ─── */
          const payload = await subRes.json();
          // console.log("   JSON payload:", payload);   // uncomment if curious

          let list;
          if (Array.isArray(payload)) {
            list = payload;                     // backend returns an array
          } else if (Array.isArray(payload.submissions)) {
            list = payload.submissions;         // { submissions: [ … ] }
          } else if (Array.isArray(payload.data)) {
            list = payload.data;                // { data: [ … ] }
          } else {
            list = [payload];                   // single object fallback
          }

          list.forEach((s) => {
            submissionMap[s.id] = s;
          });
        } catch (e) {
          console.log("   ⚠️  request failed:", e.message);
        }
      }

      /* DEBUG */
      console.log("submissionMap keys:", Object.keys(submissionMap).slice(0, 5));
      const sample = submissionMap[interviewRows[0]?.job_submission_id];
      console.log("sample row for first interview:", sample);

      /* 4️⃣  combine everything */
      const enriched = interviewRows.map((row) => ({
        ...row,
        availability: avMap[row.availability_id] || null,
        candidate   : submissionMap[row.job_submission_id] || null,
      }));

      setInterviews(enriched);
      setLoading(false);
    } catch (e) {
      setError(e.message || "Unknown error.");
      setLoading(false);
    }
  };

  /* ───────── helpers ───────── */
  const formatDateTime = (isoDate, from, to) => {
    if (!isoDate) return "Date not available";
    try {
      const base = new Date(isoDate);
      const dStr = base.toLocaleDateString(undefined, {
        year: "numeric",
        month: "long",
        day: "numeric",
      });

      if (!from) return dStr;

      const mk = (t) => {
        const [h, m] = t.split(":").map(Number);
        const d = new Date(base);
        d.setHours(h, m, 0, 0);
        return d.toLocaleTimeString(undefined, { hour: "numeric", minute: "2-digit" });
      };
      const start = mk(from);
      return to ? `${dStr} ${start} – ${mk(to)}` : `${dStr} ${start}`;
    } catch {
      return "Invalid date";
    }
  };

  const displayName = (cand) => cand?.username || "Candidate";

  /* ───────── feedback & status helpers (unchanged core logic) ───────── */
  const handleFeedbackChange = (id, txt) =>
    setFeedbacks((p) => ({ ...p, [id]: { ...p[id], feedback: txt } }));

  const updateStatus = async (subId, status) => {
    try {
      const r = await fetch(`${API_BASE}/jobs/submissions/${subId}/status`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ status }),
      });
      if (!r.ok) throw new Error((await r.json()).message || "Status update failed.");
      return true;
    } catch (e) {
      alert(e.message);
      return false;
    }
  };

  const saveFeedback = async (intId, verdict) => {
    const row = interviews.find((r) => r.id === intId);
    const fb  = feedbacks[intId]?.feedback?.trim();
    if (!row) return alert("Interview not found.");
    if (!fb)  return alert("Provide feedback before submitting.");

    try {
      const fbRes = await fetch(`${API_BASE}/interviews/${intId}/feedback`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ feedback: fb, verdict }),
      });
      if (!fbRes.ok) throw new Error((await fbRes.json()).message || "Feedback submit failed.");

      const ok = await updateStatus(
        row.job_submission_id,
        verdict === "selected" ? "selected" : "rejected"
      );
      if (!ok) return;

      setInterviews((p) =>
        p.map((r) =>
          r.id === intId
            ? { ...r, feedback: fb, verdict, feedback_submitted: true }
            : r
        )
      );
      setFeedbacks((p) => ({ ...p, [intId]: { feedback: fb, verdict } }));
      setSuccess(
        `Feedback saved • candidate ${verdict === "selected" ? "selected" : "rejected"}`
      );
      setTimeout(() => setSuccess(""), 5000);
    } catch (e) {
      alert(e.message);
    }
  };

  const resetFeedback = (id) => {
    setInterviews((p) =>
      p.map((r) =>
        r.id === id
          ? { ...r, feedback: null, verdict: null, feedback_submitted: false }
          : r
      )
    );
    setFeedbacks((p) => {
      const n = { ...p };
      delete n[id];
      return n;
    });
  };

  /* ───────── UI ───────── */
  return (
    <div className="p-6 min-h-screen bg-gradient-to-br from-blue-50 to-blue-100">
      <div className="max-w-6xl mx-auto">
        <button
          onClick={() => navigate("/interviewer-dashboard")}
          className="mb-6 px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-800"
        >
          ⬅️ Back to Dashboard
        </button>

        <h1 className="text-4xl font-bold text-center text-blue-800 mb-8">
          Scheduled Interviews
        </h1>

        {error && <p className="text-red-600 text-center font-medium mb-4">{error}</p>}
        {success && (
          <p className="text-green-600 text-center font-medium mb-4">{success}</p>
        )}

        {loading ? (
          <p className="text-center">Loading interviews...</p>
        ) : interviews.length === 0 ? (
          <p className="text-center text-gray-600">No interviews scheduled.</p>
        ) : (
          <div className="grid gap-8">
            {interviews.map((itv) => (
              <div key={itv.id} className="p-6 bg-white rounded-lg shadow">
                {/* top row */}
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-xl font-semibold text-gray-800">
                      {displayName(itv.candidate)}
                    </h3>

                    <p className="text-gray-600">
                      🗓️{" "}
                      {formatDateTime(
                        itv.availability?.date || itv.created_at,
                        itv.availability?.from_time,
                        itv.availability?.to_time
                      )}
                    </p>

                    {/* DEBUG: view candidate object
                    <pre className="text-xs text-red-500">
                      {JSON.stringify(itv.candidate, null, 2)}
                    </pre>
                    */}

                    {itv.candidate?.skills && (
                      <p className="text-gray-600">
                        🛠️ <strong>Skills:</strong> {itv.candidate.skills.join(", ")}
                      </p>
                    )}

                    {itv.candidate?.resume_url && (
                      <p className="text-gray-600">
                        📎 <strong>Resume:</strong>{" "}
                        <a
                          href={itv.candidate.resume_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-500 underline"
                        >
                          Download
                        </a>
                      </p>
                    )}

                    <p className="text-gray-600">
                      <strong>Interview ID:</strong> {itv.id}
                    </p>
                    <p className="text-gray-600">
                      <strong>Job ID:</strong> {itv.job_id}
                    </p>
                  </div>

                  <div className="text-right">
                    <p className="text-gray-600">
                      <strong>Status:</strong> {itv.status || "Unknown"}
                    </p>
                    <p
                      className={`font-medium ${
                        itv.verdict === "selected"
                          ? "text-green-600"
                          : itv.verdict === "rejected"
                          ? "text-red-600"
                          : "text-yellow-600"
                      }`}
                    >
                      {itv.verdict ?? "Pending"}
                    </p>
                  </div>
                </div>

                {/* feedback block */}
                {itv.feedback ? (
                  <div className="mt-4">
                    <div className="p-4 bg-gray-50 rounded">
                      <h4 className="font-medium mb-2">Feedback Submitted:</h4>
                      <p>{itv.feedback}</p>
                    </div>

                    <button
                      onClick={() => resetFeedback(itv.id)}
                      className="mt-4 px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
                    >
                      Reset &amp; Update Feedback
                    </button>
                  </div>
                ) : (
                  <div className="mt-4">
                    <div className="mb-4">
                      <label className="block mb-2 font-medium">Feedback:</label>
                      <textarea
                        value={feedbacks[itv.id]?.feedback || ""}
                        onChange={(e) => handleFeedbackChange(itv.id, e.target.value)}
                        className="w-full p-3 border border-gray-300 rounded focus:ring-2 focus:ring-blue-400 focus:border-blue-400"
                        rows="4"
                        placeholder="Enter your feedback about the candidate..."
                      />
                    </div>

                    <div className="flex gap-4">
                      <button
                        onClick={() => saveFeedback(itv.id, "selected")}
                        className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                      >
                        Finalize
                      </button>
                      <button
                        onClick={() => saveFeedback(itv.id, "rejected")}
                        className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                      >
                        Reject
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default InterviewerInterviews;
