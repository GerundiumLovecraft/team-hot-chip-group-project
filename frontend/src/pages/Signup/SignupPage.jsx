import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import "./SignupPage.css";
import { signup } from "../../services/authentication";
import AnimatedButton from "../../components/animation/AnimatedButton";

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
      const res = await signup(username, email, password);
      localStorage.setItem("token", res.token)
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
        <h2>Sign Up Here</h2>
        <p>Enter your details below to get started</p>
      </div>
      <form className="signup-form" onSubmit={handleSubmit}>
      {error && <p className="errorMessage">{error}</p>}

        <div className="form-group">
          <label>Username</label>
          <input
            type="text"
            value={username}
            onChange={handleUsernameChange}
            placeholder="Username"
          />
        </div>

        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            value={email}
            onChange={handleEmailChange}
            placeholder="Email address"
          />
        </div>

        <div className="form-group">
          <label>Password</label>
          <input
            type="password"
            value={password}
            onChange={handlePasswordChange}
            placeholder="Create a password"
          />
        </div>

        <div className="form-group">
          <label>Confirm Password</label>
          <input
            type="password"
            value={confirmPassword}
            onChange={handleConfirmPasswordChange}
            placeholder="Confirm password"
          />
        </div>

        <AnimatedButton type="submit" className="signup-btn">
          Sign Up
        </AnimatedButton>
        <p className="login-option">
          Already have an account?{" "}
          <Link to="/login" className="login-link">Log in here</Link>
        </p>
      </form>
    </div>
  </div>
);
};

