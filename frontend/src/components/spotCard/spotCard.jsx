import {
  Wifi,
  Toilet,
  Plug,
  Volume,
  Volume1,
  Volume2,
  Moon,
} from "lucide-react";
import "./spotCard.css";
import RatingPill from "../Rating/RatingPill.jsx";

const Spot = ({ spot, onClick }) => {
  const featureIcons = {
    wifi: Wifi,
    toilets: Toilet,
    power_sockets: Plug,
    open_late: Moon,
  };

  const getNoiseIcon = (value) => {
    if (value === 1) return Volume;
    if (value === 2) return Volume1;
    return Volume2;
  };

  const renderFeature = (feature) => {
    const Icon = featureIcons[feature.feat_name];
    if (feature.feat_name === "price") {
      return <span>{"£".repeat(feature.value)}</span>;
    }
    if (feature.feat_name === "noise_level") {
      const NoiseIcon = getNoiseIcon(feature.value);
      return <NoiseIcon />;
    }
    if (Icon) {
      const isEnabled = Boolean(feature.value);
      return <Icon style={{ opacity: 1 }} />;
    }
  };

  return (
    <div className="spot-container" onClick={onClick}>
      <div className="spot-image-wrap">
        <img src={spot.image} alt={spot.name} className="spot-image" />
        <div className="spot-rating-overlay">
          <RatingPill rating={spot.average_rating} />
        </div>
      </div>

      <h3 className="spot-name">{spot.name}</h3>
      <p className="spot-address">{spot.address}</p>

      <div className="features">
        <ul>
          {spot.features.map((feature) => (
            //in css use icon name '.wifi, .toilets' to anchor each icon individually
            <li
              className={`features-value ${feature.feat_name} value-${feature.value}`}
              key={feature.feat_name}
            >
              {renderFeature(feature)}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default Spot;
