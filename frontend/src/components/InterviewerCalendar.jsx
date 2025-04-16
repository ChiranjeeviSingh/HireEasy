import React, { useState } from "react";

export function InterviewerCalendar() {
  const [date, setDate] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [slots, setSlots] = useState([]);

  const handleSave = () => {
    if (!date || !startTime || !endTime) return;
    setSlots(prev => [
      ...prev,
      { date, startTime, endTime, id: Date.now() },
    ]);
    setDate("");
    setStartTime("");
    setEndTime("");
  };

  const formatSlot = ({ date, startTime, endTime }) => {
    const d = new Date(date).toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
    const fmt = t =>
      new Date(`1970-01-01T${t}`).toLocaleTimeString(undefined, {
        hour: "numeric",
        minute: "2-digit",
      });
    return `${d}: ${fmt(startTime)} – ${fmt(endTime)}`;
  };

  return (
    <div style={containerStyle}>
      <div style={formContainerStyle}>
        <h2>Set Your Availability</h2>
        <label style={labelStyle}>
          Pick a date
          <input
            type="date"
            value={date}
            onChange={e => setDate(e.target.value)}
            style={inputStyle}
          />
        </label>
        <label style={labelStyle}>
          Start time
          <input
            type="time"
            value={startTime}
            onChange={e => setStartTime(e.target.value)}
            min="13:00"
            max="18:00"
            style={inputStyle}
          />
        </label>
        <label style={labelStyle}>
          End time
          <input
            type="time"
            value={endTime}
            onChange={e => setEndTime(e.target.value)}
            min="13:00"
            max="18:00"
            style={inputStyle}
          />
        </label>
        <button
          onClick={handleSave}
          style={buttonStyle}
          onMouseEnter={e => (e.target.style.backgroundColor = "#0056b3")}
          onMouseLeave={e => (e.target.style.backgroundColor = "#007bff")}
        >
          Save
        </button>
      </div>

      <div style={scheduleContainerStyle}>
        <h3>Your Scheduled Slots</h3>
        {slots.length === 0 && <p style={{ color: "#666" }}>No slots yet.</p>}
        {slots.map(slot => (
          <div key={slot.id} style={slotBoxStyle}>
            {formatSlot(slot)}
          </div>
        ))}
      </div>
    </div>
  );
}

// ----- Styles -----
const containerStyle = {
  display: "flex",
  alignItems: "flex-start",
  gap: "40px",
  padding: "40px",
  fontFamily: "Arial, sans-serif",
};

const formContainerStyle = {
  flex: "0 0 300px",
  display: "flex",
  flexDirection: "column",
  gap: "20px",
};

const scheduleContainerStyle = {
  flex: 1,
  padding: "10px",
  borderLeft: "1px solid #eee",
};

const labelStyle = {
  display: "flex",
  flexDirection: "column",
  fontSize: "14px",
  color: "#333",
};

const inputStyle = {
  marginTop: "8px",
  padding: "10px",
  fontSize: "14px",
  border: "1px solid #ccc",
  borderRadius: "4px",
};

const buttonStyle = {
  padding: "12px",
  fontSize: "16px",
  backgroundColor: "#007bff",
  color: "white",
  border: "none",
  borderRadius: "6px",
  cursor: "pointer",
  transition: "background 0.3s",
  marginTop: "10px",
};

const slotBoxStyle = {
  padding: "12px",
  marginBottom: "10px",
  border: "1px solid #ddd",
  borderRadius: "6px",
  backgroundColor: "#f9f9f9",
  fontSize: "14px",
};

// <-- Add this!
export default InterviewerCalendar;