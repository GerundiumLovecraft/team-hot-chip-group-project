import { NavLink, useNavigate } from "react-router-dom";
import { useState } from "react";
import { User } from "lucide-react";
import "./Navbar.css";
import { isTokenValid } from "../../helpers/authentication.js";

const Navbar = () => {
  const navigate = useNavigate();
  const isAuth = isTokenValid();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <nav className="navbar">
      <NavLink to="/" className="navbar-logo">
        sipstack
      </NavLink>

      <div className="navbar-links">
        <NavLink to="/" end className="nav-btn">Browse Spots</NavLink>
        <NavLink to="/leaderboard" className="nav-btn">Leaderboard</NavLink>
        {isAuth ? (
          <>
            <NavLink to="/spots/new" className="nav-btn">Submit Spot</NavLink>
            <div
              className="profile-menu"
              onMouseEnter={() => setMenuOpen(true)}
              onMouseLeave={() => setMenuOpen(false)}
            >
              <button className="profile-button">
                <User size={24} />
              </button>
              {menuOpen && (
                <div className="dropdown-menu">
                  <NavLink to="/profile">My Profile</NavLink>
                  <button onClick={handleLogout}>Logout</button>
                </div>
              )}
            </div>
          </>
        ) : (
          <>
            <NavLink to="/login" className={() => "nav-btn"}>Submit Spot</NavLink>
            <NavLink to="/login" className="nav-btn">Login</NavLink>
            <NavLink to="/signup" className="nav-btn">Signup</NavLink>
          </>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
