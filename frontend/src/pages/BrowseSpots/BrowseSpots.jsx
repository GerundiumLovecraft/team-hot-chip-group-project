import { useState, useEffect } from "react";
import { getAllSpots } from "../../services/spots";
import Spot from "../../components/spotCard/spotCard"

export const BrowseSpots = () => {
  const [spots, setSpots] = useState([]);

  useEffect(() => {
    // const token = localStorage.getItem("token");

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
          <Spot key={spot._id} spot={spot} />
        ))}
      </div>
    </>
  );
};
