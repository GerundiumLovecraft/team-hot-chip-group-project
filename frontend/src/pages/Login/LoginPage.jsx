import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { User, Lock } from "lucide-react";
import "./LoginPage.css";

import { login } from "../../services/authentication";

export const LoginPage = () => {
  const [usernameOrEmail, setUsernameOrEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const token = await login(usernameOrEmail, password);
      localStorage.setItem("token", token);
      navigate("/");
    } catch (err) {
      setErrorMessage("Invalid username/email or password.");
    }
  };

  const handleUsernameOrEmailChange = (event) => {
    setUsernameOrEmail(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };

  return (
    <div className="login-page">
      <div className="login-card">
        <div className="login-header">
          <h2>Login</h2>
        </div>
        {errorMessage && <p className="error-message">{errorMessage}</p>}
        <form className="login-form" onSubmit={handleSubmit}>
          <div className="form-group-user">
            <input
              id="usernameOrEmail"
              type="text"
              value={usernameOrEmail}
              onChange={handleUsernameOrEmailChange}
              placeholder="Username/Email"
              required
            />
            <User className="user-icon" />
          </div>

          <div className="form-group-pass">
            <input
              id="password"
              type="password"
              value={password}
              onChange={handlePasswordChange}
              placeholder="Password"
              required
            />
            <Lock className="lock-icon" />
          </div>

          <button type="submit" className="login-button">
            Login
          </button>
          <p className="signup-option">
            Don&apos;t have an account yet?{" "}
            <Link to="/signup" className="signup-link">
              Sign Up
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
};
