import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchUserProfile } from "../services/authentication";

export default function ProfilePage() {
    const [user, setUser] = useState(null);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem("token");
        
        if (!token) {
            navigate("/login");
            return;
        }

    fetchUserProfile(token)
      .then(({ user, token: newToken }) => {
        setUser(user);
        // Store the rotated token to keep the session alive
        localStorage.setItem("token", newToken);
      })
      .catch((err) => {
        if (err.message.startsWith("Unauthorised")) {
          navigate("/login");
        } else {
          setError(err.message);
        }
      })
      .finally(() => setLoading(false));
  }, [navigate]);

  const memberSince = user?.createdAt
    ? new Date(user.createdAt).toLocaleDateString("en-GB", {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : null;

  return (
    <main className="profile-page">
      <div className="profile-card">
        <div className="profile-header">
          <div className="avatar">
            {user ? user.username.slice(0, 2).toUpperCase() : "??"}
          </div>
          <h1 className="profile-title">My Profile</h1>
        </div>

        {loading && <p className="status-text">Loading profile…</p>}
        {error && <p className="status-text error">{error}</p>}

        {user && (
          <ul className="profile-fields">
            <li className="profile-field">
              <span className="field-label">Username</span>
              <span className="field-value">{user.username}</span>
            </li>
            <li className="profile-field">
              <span className="field-label">Email</span>
              <span className="field-value">{user.email}</span>
            </li>
            <li className="profile-field">
              <span className="field-label">Member since</span>
              <span className="field-value">{memberSince}</span>
            </li>
          </ul>
        )}
      </div>
    </main>
  );
}
