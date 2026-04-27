import { useState, useEffect } from "react";
import { getAllSpots } from "../../services/spots";
import Spot from "../../components/spotCard/spotCard";
import SpotModal from "../../components/SpotModal/SpotModal";
import SpotFilter from "../../components/filterSpot/FilterBar"
import "./BrowseSpots.css";

export const BrowseSpots = () => {
  const [spots, setSpots] = useState([]);
  const [selectedSpot, setSelectedSpot] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getAllSpots()
      .then((data) => {
        setSpots(data.spots);
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => setLoading(false));
  }, []);

  const handleFilterChange = (selectedFeaturesList) => {
    console.log("Features were changed to: ", selectedFeaturesList)
  }

  return (
    <>
      <div className="browse-header">
        <h1>Discover Your Perfect Spot</h1>
        <p>Browse community submitted coffee shops, filtered to your working needs</p>
      </div>
      <div className="browse-layout">
        <div className="spots-section">
          {loading ? (
            <p>Brewing something good... ☕️</p>
          ) : spots.length !== 0 ? (
            spots.map((spot) => (
              <Spot
                key={spot._id}
                spot={spot}
                onClick={() => setSelectedSpot(spot)}
              />
            ))
          ) : (
            <p>Nothing here yet. Why not submit your favourite spot? ☕ ✨</p>
          )}
        </div>
        <div className="filter-section">
          <h3>Filter</h3>
          <SpotFilter onFilterChange={handleFilterChange}/>
        </div>
      </div>
      <SpotModal
        spot={selectedSpot}
        onClose={() => setSelectedSpot(null)}
      />
    </>
  );
};
