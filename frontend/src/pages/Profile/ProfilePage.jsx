import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchUserProfile } from "../../services/authentication";
import { getSpotsByUser } from "../../services/spots";
import SpotModal from "../../components/SpotModal/SpotModal";
import {isTokenValid} from "../../helpers/authentication.js";
// import { Link } from "react-router-dom";
import "./ProfilePage.css";

export default function ProfilePage() {
    const [user, setUser] = useState(null);
    const [mySpots, setMySpots] = useState([]);
    const [error, setError] = useState(null);
    const [spotsError, setSpotsError] = useState(null);
    const [loading, setLoading] = useState(true);
    const [editingAvatar, setEditingAvatar] = useState(false);
    const [avatarInput, setAvatarInput] = useState("");
    const [spotsLoading, setSpotsLoading] = useState(true);
    const [selectedSpot, setSelectedSpot] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        if (!isTokenValid()) {
            navigate("/login");
            return;
        }

        const token = localStorage.getItem("token");

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

        getSpotsByUser(token)
          .then((spots) => setMySpots(spots ?? []))
          .catch((err) => setSpotsError(err.message))
          .finally(() => setSpotsLoading(false));
    }, [navigate]);

  const memberSince = user?.createdAt
    ? new Date(user.createdAt).toLocaleDateString("en-GB", {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : null;

  const handleAvatarSave = () => {
    const token = localStorage.getItem("token");
    fetch(`${import.meta.env.VITE_BACKEND_URL}/users/avatar`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ avatar: avatarInput }),
    })
      .then((res) => res.json())
      .then(({ user, token: newToken }) => {
        setUser(user);
        localStorage.setItem("token", newToken);
        setEditingAvatar(false);
        setAvatarInput("");
      })
      .catch((err) => setError(err.message));
  }

  return (
    <main className="profile-page">
      <div className="profile-card">
        <div className="profile-header">
          <div className="avatar-wrapper">
            <div className="avatar" onClick={() => setEditingAvatar(true)}>
              {user?.avatar
                ? <img src={user.avatar} alt="avatar" className="avatar-image" />
                : user ? user.username.slice(0, 2).toUpperCase() : "??"
              }
              <div className="avatar-overlay">edit</div>
            </div>

            {editingAvatar && (
              <div className="avatar-edit">
                <input
                  type="text"
                  placeholder="Paste image URL…"
                  value={avatarInput}
                  onChange={(e) => setAvatarInput(e.target.value)}
                  className="avatar-input"
                />
                <div className="avatar-edit-actions">
                  <button onClick={handleAvatarSave} className="avatar-save-btn">Save</button>
                  <button onClick={() => setEditingAvatar(false)} className="avatar-cancel-btn">Cancel</button>
                </div>
              </div>
            )}
          </div>
          <h1 className="profile-title">My Profile:</h1>
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
              <span className="field-label">Sipstacker since</span>
              <span className="field-value">{memberSince}</span>
            </li>
          </ul>
        )}
      </div>
      <section className="my-spots-section">
        <h2 className="my-spots-title">✨ Your Hot Spots ✨</h2>

        {spotsLoading && <p className="status-text">Loading your spots…</p>}
        {spotsError && <p className="status-text error">{spotsError}</p>}

        {!spotsLoading && !spotsError && mySpots.length === 0 && (
          <p className="status-text">You have not submitted any spots yet.</p>
        )}

        {mySpots.length > 0 && (
          <div className="my-spots-grid">
            {mySpots.map((spot) => (
              <div key={spot._id} className="my-spot-card" onClick={() => setSelectedSpot(spot)}>
                {spot.image && (
                  <img src={spot.image} alt={spot.name} className="my-spot-image" />
                )}
                <div className="my-spot-info">
                  <p className="my-spot-name">{spot.name}</p>
                  <p className="my-spot-address">{spot.address}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
      <SpotModal spot={selectedSpot} onClose={() => setSelectedSpot(null)} isProfilePage={true} />
    </main>
  );
}
