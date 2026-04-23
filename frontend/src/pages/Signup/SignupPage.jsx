import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";

import { signup } from "../../services/authentication";

export const SignupPage = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const handleSubmit = async (event) => {
   event.preventDefault();
    setError("");

    if (!username || !email || !password || !confirmPassword) {
        setError("All fields must be filled in");
        return;
    }
    
    const passwordRegex = /^.{8,}$/
    if (!passwordRegex.test(password)) {
      return setError("Password must be at least 8 characters long.")
    }
    if (password !== confirmPassword) {
      return setError("Passwords do not match.")
    }
  
    try {
      await signup(username, email, password);
      navigate("/");
    } catch (err) {
      console.error(err);
      setError(err.message)
      setUsername("")
      setEmail("")
      setPassword("")
      setConfirmPassword("")
    }
  }


  const handleUsernameChange = (event) => {
    setUsername(event.target.value);
  };

  const handleEmailChange = (event) => {
    setEmail(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };

  const handleConfirmPasswordChange = (event) => {
    setConfirmPassword(event.target.value);
  }


return (
  <div className="signup-page">
    <div className="signup-card">

      <div className="signup-header">
        <h2>Create an account</h2>
        <p>Enter your details below to get started</p>
        <Link to="/login" className="login-link">Already have an account? Log in here</Link>
      </div>

      {error && <p className="errorMessage">{error}</p>}

      <form className="signup-form" onSubmit={handleSubmit}>

        <div className="form-group">
          <label>Username</label>
          <input
            type="text"
            value={username}
            onChange={handleUsernameChange}
          />
        </div>

        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            value={email}
            onChange={handleEmailChange}
          />
        </div>

        <div className="form-group">
          <label>Password</label>
          <input
            type="password"
            value={password}
            onChange={handlePasswordChange}
          />
        </div>

        <div className="form-group">
          <label>Confirm Password</label>
          <input
            type="password"
            value={confirmPassword}
            onChange={handleConfirmPasswordChange}
          />
        </div>

        <button type="submit" className="signup-btn">
          Sign Up
        </button>

      </form>
    </div>
  </div>
);
};

