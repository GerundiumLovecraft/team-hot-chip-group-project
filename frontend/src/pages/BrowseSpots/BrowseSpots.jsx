import { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";

import { getSpots } from "../../services/spots";

export const BrowseSpots = () => {
  const [spots, setSpots] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");

    getSpots(token)
      .then((data) => {
        setSpots(data.spots);
      })
      .catch((err) => {
        console.error(err);
      });
  }, []);

  const logOutHandler = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <>
      <h2>Spots</h2>

      <nav>
        <Link to="/login">Log In</Link>
        {" | "}
        <Link to="/signup">Sign Up</Link>
      </nav>

      <div className="feed" role="feed">
        {spots.map((spot) => (
          <div key={spot._id}>
            <h3>{spot.name}</h3>
            <p>{spot.address}</p>
            <p>{spot.description}</p>
            <p>
              {spot.open_from} - {spot.open_to}
            </p>
          </div>
        ))}
      </div>

      <button onClick={logOutHandler}>Log Out</button>
    </>
  );
};