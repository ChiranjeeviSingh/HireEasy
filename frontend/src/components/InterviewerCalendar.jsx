import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

export default function InterviewerCalendar() {
  const [date, setDate] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [slots, setSlots] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";
  const token = localStorage.getItem("token");
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) return;

    const today = new Date();
    const fromDate = today.toISOString().split("T")[0];
    const toDate = new Date(today.setDate(today.getDate() + 7)).toISOString().split("T")[0];

    fetch(`${API_BASE}/availability/me?from_date=${fromDate}&to_date=${toDate}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error("Failed to fetch availability");
        const data = await res.json();
        setSlots(data);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [token, API_BASE]);

  const handleSave = async () => {
    if (!date || !startTime || !endTime) return;

    const newSlot = {
      date,
      from_time: `${startTime}:00`,
      to_time: `${endTime}:00`,
    };

    try {
      const res = await fetch(`${API_BASE}/availability`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(newSlot),
      });

      if (!res.ok) throw new Error("Could not create availability");

      const created = await res.json();
      setSlots((prev) => [...prev, created]);
      setDate("");
      setStartTime("");
      setEndTime("");
    } catch (err) {
      setError(err.message);
    }
  };

  const handleDelete = async (id) => {
    try {
      const res = await fetch(`${API_BASE}/availability/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) throw new Error("Failed to delete slot");

      setSlots((prev) => prev.filter((s) => s.id !== id));
    } catch (err) {
      setError(err.message);
    }
  };

  const formatSlot = ({ date, from_time, to_time }) => {
    const d = new Date(date).toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });

    const fmt = (t) =>
      new Date(`1970-01-01T${t}`).toLocaleTimeString(undefined, {
        hour: "numeric",
        minute: "2-digit",
      });

    return `${d}: ${fmt(from_time)} – ${fmt(to_time)}`;
  };

  return (
    <div className="relative min-h-screen bg-gray-100">
      {/* Background Image */}
      <div
        className="absolute inset-0 bg-cover bg-center"
        style={{
          backgroundImage:
            "url('https://images.unsplash.com/photo-1521737604893-d14cc237f11d?auto=format&fit=crop&w=1950&q=80')",
          zIndex: -1,
          filter: "brightness(0.6)",
        }}
      ></div>

      <div className="relative z-10 flex flex-col md:flex-row gap-10 p-10 backdrop-blur-sm bg-white/80 rounded-lg max-w-6xl mx-auto mt-10 shadow-lg">
        {/* Form Section */}
        <div className="flex-1">
          <h2 className="text-2xl font-bold mb-4 text-blue-800">Set Your Availability</h2>
          {error && <p className="text-red-600 mb-2">{error}</p>}

          <label className="block mb-4">
            <span className="block text-gray-700 font-medium mb-1">Pick a date</span>
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded-md"
            />
          </label>

          <label className="block mb-4">
            <span className="block text-gray-700 font-medium mb-1">Start time</span>
            <input
              type="time"
              value={startTime}
              onChange={(e) => setStartTime(e.target.value)}
              min="13:00"
              max="18:00"
              className="w-full p-2 border border-gray-300 rounded-md"
            />
          </label>

          <label className="block mb-4">
            <span className="block text-gray-700 font-medium mb-1">End time</span>
            <input
              type="time"
              value={endTime}
              onChange={(e) => setEndTime(e.target.value)}
              min="13:00"
              max="18:00"
              className="w-full p-2 border border-gray-300 rounded-md"
            />
          </label>

          <div className="flex gap-4 mt-4">
            <button
              onClick={handleSave}
              className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
            >
              Save
            </button>
            <button
              onClick={() => navigate("/interviewer-dashboard")}
              className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600 transition"
            >
              Cancel
            </button>
          </div>
        </div>

        {/* Slots Display Section */}
        <div className="flex-1">
          <h3 className="text-xl font-semibold text-gray-800 mb-4">Your Scheduled Slots</h3>
          {loading && <p className="text-gray-600">Loading...</p>}
          {!loading && slots.length === 0 && (
            <p className="text-gray-500">No slots yet. Add your first availability above.</p>
          )}
          <div className="space-y-3">
            {slots.map((slot) => (
              <div
                key={slot.id}
                className="flex justify-between items-center p-3 border border-gray-300 rounded-md bg-white shadow-sm"
              >
                <span className="text-gray-800">{formatSlot(slot)}</span>
                <button
                  onClick={() => handleDelete(slot.id)}
                  className="text-red-600 text-lg hover:text-red-800"
                >
                  ✕
                </button>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}