import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { User, Lock } from "lucide-react";
import AnimatedButton from "../../components/animation/AnimatedButton";
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
        <form className="login-form" onSubmit={handleSubmit}>
        {errorMessage && <p className="error-message">{errorMessage}</p>}
          <div className="form-group-user">
            <label>Username</label>
            <input
              id="usernameOrEmail"
              type="text"
              value={usernameOrEmail}
              onChange={handleUsernameOrEmailChange}
              placeholder="Your username or email"
              required
            />
            <User className="user-icon" />
          </div>

          <div className="form-group-pass">
            <label>Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={handlePasswordChange}
              placeholder="Your secret password"
              required
            />
            <Lock className="lock-icon" />
          </div>

          <AnimatedButton type="submit" className="login-button">
            Login
          </AnimatedButton>
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
