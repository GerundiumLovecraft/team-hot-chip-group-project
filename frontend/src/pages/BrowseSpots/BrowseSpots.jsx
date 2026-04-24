import { useState, useEffect } from "react";
import { getAllSpots } from "../../services/spots";
import Spot from "../../components/spotCard/spotCard"

export const BrowseSpots = () => {
  const [spots, setSpots] = useState([]);

  useEffect(() => {
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
        {spots.length !== 0 ?
        spots.map((spot) => (
          <Spot key={spot._id} spot={spot}/>
        )):
        <p>Nothing here yet. Why not add and submit your favourite spot?{"\u2615"} {"\u2728"}</p>
        }
      </div>
    </>
  );
};
