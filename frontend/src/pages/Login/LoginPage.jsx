import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { login } from "../../services/authentication";

export const LoginPage = () => {
  const [usernameOrEmail, setUsernameOrEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const token = await login(usernameOrEmail, password);
      localStorage.setItem("token", token);
      navigate("/");
    } catch (err) {
      console.error(err);
      navigate("/login");
    }
  };

  const handleUsernameOrEmailChange = (event) => {
    setUsernameOrEmail(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };

  return (
    <>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <label htmlFor="usernameOrEmail">Username/Email:</label>
        <input
          id="usernameOrEmail"
          type="text"
          value={usernameOrEmail}
          onChange={handleUsernameOrEmailChange}
        />
        <label htmlFor="password">Password:</label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={handlePasswordChange}
        />
        <input role="submit-button" id="submit" type="submit" value="Submit" />
      </form>
    </>
  );
};
