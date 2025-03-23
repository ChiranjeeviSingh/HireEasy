import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/styles.css";

export function Login() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const mockUser = { email: "abcd@gmail.com", password: "abcdef567" };

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

 
    try {
      const response = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      if (!response.ok) {
        throw new Error("Login failed");
      }
      const {token} = await response.json();
      console.log("Token ",token);
      alert("Login Successful!");
      localStorage.setItem("token", token.token);
      navigate("/dashboard");
      setLoading(false);
    } catch (error) {
      setError("Invalid credentials, please try again.");
      setLoading(false);
    }
  };

  return (
    <div
      className="relative min-h-screen flex justify-center items-center bg-cover bg-center"
      style={{
        backgroundImage:
          "url('https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?fm=jpg&q=60&w=3000&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8am9iJTIwcG9ydGFsfGVufDB8fDB8fHww')",
      }}
    >
      {/* ✅ Login Card Centered */}
      <div className="bg-white bg-opacity-90 shadow-lg rounded-xl p-10 w-full max-w-md text-center">
        {/* ✅ HireEasy Title Styled with Tailwind */}
        <h1 className="text-4xl logo font-bold text-gray-800 tracking-wide font-sans mb-4">
          HireEasy
        </h1>

        <p className="text-gray-600">
          Sign in and start hiring the best talent out there.
        </p>

        <form onSubmit={handleLogin} className="mt-6">
          {error && <p className="text-red-500 mb-3">{error}</p>}

          {/* Email Input */}
          <div className="mb-4">
            <input
              id="email"
              className="w-full p-3 border border-gray-400 rounded-lg focus:ring-2 focus:ring-blue-400"
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
              className="w-full p-3 border border-gray-400 rounded-lg focus:ring-2 focus:ring-blue-400"
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {/* Sign-In Button */}
          <button
            type="submit"
            className="w-full py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition"
            disabled={loading}
          >
            {loading ? "Logging in..." : "Sign In"}
          </button>
        </form>

        {/* Register Link */}
        <div className="mt-4 text-gray-600">
          Don't have an account?{" "}
          <a
            href="#"
            onClick={() => navigate("/register")}
            className="text-blue-500"
          >
            Create One Now
          </a>
        </div>
      </div>
    </div>
  );
}

export default Login;
