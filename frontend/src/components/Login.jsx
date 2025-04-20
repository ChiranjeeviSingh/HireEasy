import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/styles.css";

export function Login() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_BASE}/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });

      const data = await response.json();
      console.log("Response Data:", data);

      if (!response.ok) {
        throw new Error(data.msg || "Invalid credentials, please try again.");
      }

      const token = data.token.token;
      const user = data.token.user;

      // Save token and role
      localStorage.setItem("token", token);
      localStorage.setItem("role", user.role);
      localStorage.setItem("username", user.username);

      alert("Login Successful!");

      // Redirect based on role
      if (user.role === "HR") {
        navigate("/dashboard");
      } else if (user.role === "Interviewer") {
        navigate("/interviewer-dashboard");
      } else {
        setError("Unknown user role. Please contact support.");
      }

    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="container main-wrapper"
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
      }}
    >
      <div
        className="login-page"
        style={{
          width: "40%",
          maxWidth: "900px",
          padding: "50px",
          borderRadius: "10px",
          boxShadow: "0 4px 8px rgba(0,0,0,0.1)",
          backgroundColor: "#fff",
          textAlign: "center",
        }}
      >
        <h1 className="logo">HireEasy</h1>
        <div className="my-3">
          <h2>Welcome to Careerbuilder</h2>
          <p>Sign in and start hiring the best talent out there.</p>
        </div>
        <form onSubmit={handleLogin}>
          {error && <p className="error-text">{error}</p>}
          <div className="form-group">
            <input
              id="email"
              className="form-control"
              type="email"
              placeholder="Enter your email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <input
              id="password"
              className="form-control"
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button
            type="submit"
            className="auth-button btn-block my-2"
            style={{ width: "100%" }}
            disabled={loading}
          >
            {loading ? "Logging in..." : "Sign In"}
          </button>
        </form>
        <div className="text-center">
          Don't have an account?{" "}
          <a href="#" onClick={() => navigate("/register")}>
            Create One Now
          </a>
        </div>
      </div>
    </div>
  );
}

export default Login;