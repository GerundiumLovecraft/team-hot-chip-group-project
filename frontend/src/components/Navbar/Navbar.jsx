import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";
import { User } from "lucide-react";
import "./Navbar.css";

const Navbar = () => {
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-logo">
        sipstack
      </Link>

      <div className="navbar-links">
        <Link to="/">Browse Spots</Link>
        <Link to="/leaderboards">Leaderboards</Link>
        {token ? (
          <>
            <Link to="/spots/new">Submit Spot</Link>
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
                  <Link to="/profile">My Profile</Link>
                  <button onClick={handleLogout}>Logout</button>
                </div>
              )}
            </div>
          </>
        ) : (
          <>
            <Link to="/login">Submit Spot</Link>
            <Link to="/login">Login</Link>
            <Link to="/signup">Signup</Link>
          </>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
