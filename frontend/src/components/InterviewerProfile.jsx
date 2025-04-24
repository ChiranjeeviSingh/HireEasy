import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

function InterviewerProfile() {
  const navigate = useNavigate();

  const [profile, setProfile] = useState(null);
  const [form, setForm] = useState({
    job_title: "",
    years_of_experience: "",
    areas_of_expertise: "",
    phone_number: "",
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";
  const token = localStorage.getItem("token");

  useEffect(() => {
    fetch(`${API_BASE}/profiles/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error("No profile found");
        const data = await res.json();
        setProfile(data);
        setForm({
          job_title: data.job_title || "",
          years_of_experience: data.years_of_experience || "",
          areas_of_expertise: data.areas_of_expertise?.join(", ") || "",
          phone_number: data.phone_number || "",
        });
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [API_BASE, token]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError("");
    setSuccess("");

    const body = {
      job_title: form.job_title,
      years_of_experience: parseInt(form.years_of_experience),
      areas_of_expertise: form.areas_of_expertise.split(",").map((s) => s.trim()),
      phone_number: form.phone_number,
    };

    const method = profile ? "PUT" : "POST";
    const endpoint = profile ? `${API_BASE}/profiles/` : `${API_BASE}/profiles`;

    try {
      const res = await fetch(endpoint, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!res.ok) throw new Error("Something went wrong");
      const updated = await res.json();
      setProfile(updated);
      setSuccess("Profile saved successfully!");
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <div className="text-center p-8 text-lg">Loading profile...</div>;

  return (
    <div className="relative min-h-screen bg-gray-100">
      {/* Background Image */}
      <div
        className="absolute inset-0 bg-cover bg-center"
        style={{
          backgroundImage:
            "url('https://images.unsplash.com/photo-1503676260728-1c00da094a0b?auto=format&fit=crop&w=1950&q=80')",
          zIndex: -1,
          filter: "brightness(0.5)",
        }}
      ></div>

      {/* Form Container */}
      <div className="min-h-screen flex items-center justify-center px-4">
        <div className="max-w-2xl w-full bg-white bg-opacity-95 backdrop-blur-sm shadow-xl p-10 rounded-2xl">
          <h2 className="text-3xl font-bold mb-6 text-center text-blue-800">
            Interviewer Profile
          </h2>

          {error && <div data-cy="error-message" className="text-red-600 mb-3 font-medium">{error}</div>}
          {success && <div data-cy="success-message" className="text-green-600 mb-3 font-medium">{success}</div>}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block mb-1 font-medium text-gray-700">Job Title</label>
              <input
                type="text"
                name="job_title"
                value={form.job_title}
                onChange={handleChange}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200"
                required
                data-cy="job-title-input"
              />
            </div>

            <div>
              <label className="block mb-1 font-medium text-gray-700">
                Years of Experience
              </label>
              <input
                type="number"
                name="years_of_experience"
                value={form.years_of_experience}
                onChange={handleChange}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200"
                required
                min={0}
                data-cy="years-of-experience-input"
              />
            </div>

            <div>
              <label className="block mb-1 font-medium text-gray-700">
                Areas of Expertise (comma separated)
              </label>
              <input
                type="text"
                name="areas_of_expertise"
                value={form.areas_of_expertise}
                onChange={handleChange}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200"
                required
                data-cy="areas-of-expertise-input"
              />
            </div>

            <div>
              <label className="block mb-1 font-medium text-gray-700">Phone Number</label>
              <input
                type="text"
                name="phone_number"
                value={form.phone_number}
                onChange={handleChange}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring focus:ring-blue-200"
                data-cy="phone-number-input"
              />
            </div>

            <div className="flex gap-4 mt-6">
              <button
                type="submit"
                disabled={saving}
                className="bg-blue-600 text-white py-2 px-6 rounded-lg hover:bg-blue-700 transition disabled:opacity-60"
                data-cy="save-profile-button"
              >
                {saving ? "Saving..." : "Save Profile"}
              </button>

              <button
                type="button"
                onClick={() => navigate("/interviewer-dashboard")}
                className="bg-gray-500 text-white py-2 px-6 rounded-lg hover:bg-gray-600 transition"
                data-cy="cancel-button"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

export default InterviewerProfile;