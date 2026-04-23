import { useState, useEffect } from "react";
import { getAllSpots } from "../../services/spots";

export const BrowseSpots = () => {
  const [spots, setSpots] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");

    getAllSpots()
      .then((data) => {
        setSpots(data.spots);
      })
      .catch((err) => {
        console.error(err);
      });
  }, []);

  return (
    <>
      <h2>Spots</h2>

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
    </>
  );
};
