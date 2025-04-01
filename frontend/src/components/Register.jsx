import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export function Register() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  const handleRegister = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await fetch(`${API_BASE}/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          email,
          password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.msg || "Registration failed");
      }

      setSuccess("Registration successful! You can now log in.");
      setTimeout(() => navigate("/"), 1500); // Redirect to login after success
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative min-h-screen flex items-center justify-center bg-gray-100">
      {/* Background Image */}
      <div
        className="absolute inset-0 bg-cover bg-center"
        style={{
          backgroundImage:
            "url('https://images.unsplash.com/photo-1560264280-88b68371db39?ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2000&q=80')",
          filter: "brightness(0.4)",
        }}
      ></div>

      {/* Register Card */}
      <div className="relative bg-white bg-opacity-95 shadow-lg rounded-2xl p-10 w-full max-w-md text-center">
        {/* Logo & Welcome Text */}
        <h1 className="text-4xl font-bold text-gray-800 tracking-wide font-sans mb-3">
          HireEasy
        </h1>
        <p className="text-gray-600">Sign up and join the best hiring network.</p>

        {/* Registration Form */}
        <form onSubmit={handleRegister} className="mt-6">
          {error && <p className="text-red-500 mb-3">{error}</p>}
          {success && <p className="text-green-500 mb-3">{success}</p>}

          {/* Username Input */}
          <div className="mb-4">
            <input
              id="username"
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
              type="text"
              placeholder="Enter your username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>

          {/* Email Input */}
          <div className="mb-4">
            <input
              id="email"
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
              type="email"
              placeholder="Enter your email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          {/* Password Input */}
          <div className="mb-4">
            <input
              id="password"
              className="w-full p-3 border border-gray-300 rounded-lg text-lg focus:ring-2 focus:ring-blue-400"
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {/* Register Button */}
          <button
            type="submit"
            className="w-full py-3 bg-green-500 text-white text-lg rounded-lg hover:bg-green-600 transition"
            disabled={loading}
          >
            {loading ? "Registering..." : "Sign Up"}
          </button>
        </form>

        {/* Login Redirect Link */}
        <div className="mt-4 text-gray-600">
          Already have an account?{" "}
          <a
            href="#"
            onClick={() => navigate("/")}
            className="text-blue-500 hover:underline"
          >
            Login
          </a>
        </div>
      </div>
    </div>
  );
}

export default Register;