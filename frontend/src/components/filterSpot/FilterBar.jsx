import {
  Wifi,
  Toilet,
  Plug,
  Volume,
  Volume1,
  Volume2,
  Moon,
} from "lucide-react";
import { getAllFeatures } from "../../services/features";
import { useEffect, useState } from "react";

const FilterBar = ({ onFilterChange }) => {
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

  // const featureGroups = {
  //   "Amenities": ["wifi", "toilets", "power_sockets"],
  //   "Noise Level": ["noise_level"],
  //   "Price": ["price"],
  //   "Open Late": ["open_late"]
  // }

  const handleFilterClick = (feature) => {
    let updatedSelectedFeatures = [];

    // when the user selected (to untick) a feature in the selectedFeatures, the feature is removed
    if (selectedFeatures.includes(feature.feat_id)) {
      updatedSelectedFeatures = selectedFeatures.filter(
        (item) => item !== feature.feat_id,
      );
    } else {
      // when the user selected (to tick) a feature not in the selectedFeatures, the feature is added
      updatedSelectedFeatures = [...selectedFeatures, feature.feat_id];
    }
    // Setter to update the list with the updated selected features
    setSelectedFeatures(updatedSelectedFeatures);
    // the new list is send to the parent
    console.log("sending to backend:", updatedSelectedFeatures)
    onFilterChange(updatedSelectedFeatures.map(id => ({id: id})));
  };

  //   const handleSingleFilterClick = (featId) => {
  //   let updatedSingleSelectedFeature = [];
      
  //   const isAlreadySelected = selectedFeatures.some(feat => feat.id === featId && feat.value === value);
    
  //   if (isAlreadySelected) {
  //     updatedSingleSelectedFeature = selectedFeatures
  //       .filter(item => item.id !== featId || item.value !== value);
  //   } else {
  //     const newSelectedFeatures = selectedFeatures.filter(item => item.id !== featId);
  //     updatedSingleSelectedFeature = [...newSelectedFeatures, {id: featId, value: value}];
  //   }
  //   // Setter to update the list with the updated selected features
  //   setSelectedFeatures(updatedSingleSelectedFeature);
  //   // the new list is send to the parent
  //   console.log("sending to backend:", updatedSingleSelectedFeature)
  //   onFilterChange(updatedSingleSelectedFeature);
  // };

  // When clear is clicked:
  const handleClearFeaturesClick = () => {
    setSelectedFeatures([]);
    onFilterChange([]);
  };

  const renderAmenityButtons = () => {
    const AmenitiesFeature = features.filter(feat => ["wifi", "toilets", "power_sockets"].includes(feat.feat_name));

    return (
      AmenitiesFeature.map((feature) => {
        const AmentiesIcon = buttonFeatureIcon[feature.feat_name]
        return (
          <button
          key={feature.feat_id} 
          onClick={() => handleFilterClick(feature)}
          className={`amenities-button ${selectedFeatures.includes(feature.feat_id) ? "active" : ""}`}>
            {AmentiesIcon && <AmentiesIcon size={18}/>}
            <span className="amenities-tag">{feature.feat_name.replace(/_/g, ' ')}</span>
          </button>
        )
      })
    )
  }

  // const renderNoiseLevelButtons = () => {
  //   const noiseLevelsFeature = features.filter(feat => ["noise_level"].includes(feat.feat_name))

  //   const NoiseLevelIcon = [
  //     {label: "Quiet", value: 1, icon: Volume},
  //     {label: "Moderate", value: 2, icon: Volume1},
  //     {label: "Loud", value: 3, icon: Volume2}
  //   ]

  //   const noiseFeat = noiseLevelsFeature[0]
  //   return (
  //     NoiseLevelIcon.map((feature) => {
  //       const NoiseIcon = feature.icon
  //       return (
  //         <button
  //         key={feature.value} 
  //         onClick={() => handleSingleFilterClick(noiseFeat.feat_id, feature.value)}
  //         className={`noise-button ${selectedFeatures.some(feat => feat.id === noiseFeat.feat_id && feat.value === feature.value) ? "active" : ""}`}>
  //           {NoiseIcon && <NoiseIcon size={18}/>}
  //           <span className="noise-tag">{feature.label}</span>
  //         </button>
  //       )
  //     })
  //   )
  // }

  // const renderPriceButtons = () => {
  //   const priceFeature = features.filter(feat => ["price"].includes(feat.feat_name))

  //   const priceLabel = [
  //     {label: "£", value: 1},
  //     {label: "££", value: 2},
  //     {label: "£££", value: 3},
  //   ];

  //   return priceFeature.map((feature) => {
  //     priceLabel && <priceLabel size={18}/>
  //     return <button key={feature.feat_id} onClick={() => handleFilterClick(feature)}></button>
  //   })
  // }

  const renderOpeningLateButton = () => {
    const openingLateFeature = features.filter(feat => ["open_late"].includes(feat.feat_name));

    return (
      openingLateFeature.map((feature) => {
        const OpeningLateIcon = buttonFeatureIcon[feature.feat_name]
        return (
          <button
          key={feature.feat_id}
          onClick={() => handleFilterClick(feature)}
          className={`open-late-button ${selectedFeatures.includes(feature.feat_id) ? "active" : ""}`}>
          {OpeningLateIcon && <OpeningLateIcon size={18}/>}
          <span className="open-late-tag">{feature.feat_name.replace(/_/g, ' ')}</span>
          </button>
        )
      })
    ) 
  }

  return (
    <>
      <div className="filter-container">
        <div className="filter-header">
          <h3>Filter your spots</h3>
        </div>
        {loading ? (
          <p>Loading features... ⏳</p>
        ) : (
          <>
          <div className="amenities-section">
            <h3>Amenities</h3>
            {renderAmenityButtons()}
          </div>
          <div className="noise-section">
            <h3>Noise Level</h3>
            {/* {renderNoiseLevelButtons()} */}
          </div>
          <div className="price-section">
            <h3>Price</h3>
            {/* {renderPriceButtons()} */}
          </div>
          <div className="open_late">
            <h3>Open Late</h3>
            {renderOpeningLateButton()}
          </div>
          </>
        )}
        <div className="clear-section">
          <button
            onClick={handleClearFeaturesClick}
            className="clear-feature-button"
          >
            Clear all
          </button>
        </div>
      </div>
    </>
  );
};

export default FilterBar;
