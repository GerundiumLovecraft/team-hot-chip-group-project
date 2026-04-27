import {
  Wifi,
  Toilet,
  Plug,
  Volume1,
  Moon,
} from "lucide-react";
import { getAllFeatures } from "../../services/features";
import { useEffect, useState } from "react";

const SpotFilter = ({ onFilterChange }) => {
  const buttonFeatureIcon = {
    "wifi": Wifi,
    "toilets": Toilet,
    "power_sockets": Plug,
    "noise_level": Volume1,
    "open_late": Moon,
  };

  const [features, setFeatures] = useState([]);
  const [selectedFeatures, setSelectedFeatures] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getAllFeatures()
      .then((data) => {
        setFeatures(data.features);
      })
      .catch((errorMessage) => {
        console.error(errorMessage);
      })
      .finally(() => setLoading(false));
  }, []);

  const handleFilterClick = (feature) => {
    let updatedSelectedFeatures = [];

    // when the user selected (to untick) a feature in the selectedFeatures, the feature is removed
    if (selectedFeatures.includes(feature.feat_name)) {
      updatedSelectedFeatures = selectedFeatures.filter(
        (item) => item !== feature.feat_name,
      );
    } else {
      // when the user selected (to tick) a feature not in the selectedFeatures, the feature is added
      updatedSelectedFeatures = [...selectedFeatures, feature.feat_name];
    }
    // Setter to update the list with the updated selected features
    setSelectedFeatures(updatedSelectedFeatures);
    // the new list is send to the parent
    onFilterChange(updatedSelectedFeatures);
  };

  // When clear is clicked:
  const handleClearFeaturesClick = () => {
    setSelectedFeatures([]);
    onFilterChange([]);
  };

  return (
    <>
      <div className="filter-container">
        <div className="filter-header">
          <h3>Filter all spots by feature</h3>
          <button
            onClick={handleClearFeaturesClick}
            className="clear-feature-button"
          >
            Clear all
          </button>
        </div>
        {loading ? (
          <p>Loading features... ⏳</p>
        ) : (
          <div className="filter-grid">
            {features.map((feature) => {
              const FeatureIcon = buttonFeatureIcon[feature.feat_name];
              return (
                <button
                  key={feature.feat_id}
                  onClick={() => handleFilterClick(feature)}
                  className={`filter-button ${
                    selectedFeatures.includes(feature) ? "active" : ""
                  }`}
                >
                  {feature.feat_name === "price" ? (
                    <span className="price-tag">
                      {"£".repeat(feature.value)}
                    </span>
                  ) : (
                    FeatureIcon && <FeatureIcon size={18} />
                  )}
                  <span className="feature-tag">{feature.feat_name}</span>
                </button>
              );
            })}
            ;
          </div>
        )}
      </div>
    </>
  );
};

export default SpotFilter;
