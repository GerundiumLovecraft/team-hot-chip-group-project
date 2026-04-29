import { useState, useEffect } from "react";
import { getAllSpots, getSpotByFeature } from "../../services/spots";
import Spot from "../../components/spotCard/spotCard";
import SpotModal from "../../components/SpotModal/SpotModal";
import FilterBar from "../../components/filterSpot/FilterBar"
import "./BrowseSpots.css";

export const BrowseSpots = () => {
  const [spots, setSpots] = useState([]);
  const [selectedSpot, setSelectedSpot] = useState(null);
  const [selectedFeatures, setSelectedFeatures] = useState([])
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);

    const fetchingData = selectedFeatures.length > 0 ? getSpotByFeature(selectedFeatures) : getAllSpots();
    fetchingData
      .then((data) => {
        setSpots(data.spots || []);
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => setLoading(false));
  }, [selectedFeatures]);

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
            <p>No spots were found. Why not submit your favourite spot? ☕ ✨</p>
          )}
        </div>
        <div className="filter-section">
          <FilterBar 
          onFilterChange={setSelectedFeatures}
          selectedFeatures={selectedFeatures}
          />
        </div>
      </div>
      <SpotModal
        spot={selectedSpot}
        onClose={() => setSelectedSpot(null)}
      />
    </>
  );
};
