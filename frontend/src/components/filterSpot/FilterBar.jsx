import {SlidersHorizontal, Wifi, Toilet, Plug, Volume, Volume1, Volume2, Moon,} from "lucide-react";
import { useEffect, useState } from "react";
import { getFeatures } from "../../services/features";
import "./FilterBar.css";

const amenitiesFeatureIcons = {
    "wifi": Wifi,
    "toilets": Toilet,
    "power_sockets": Plug,
    // "open_late": Moon
};

const amenitiesFeatureLabels = {
    wifi: "Wifi",
    toilets: "Toilets",
    power_sockets: "Power Sockets",
    // open_late: "Open Late"
};

const noiseLevelOptions = [
    {label: "Quiet", value: "1", Icon: Volume},
    {label: "Moderate", value: "2", Icon: Volume1},
    {label: "Loud", value: "3", Icon: Volume2},
];

const priceOptions = [
    {label: "£", value: "1"},
    {label: "££", value: "2"},
    {label: "£££", value: "3"},

];

export default function FilterBar({selectedFeatures, onFilterChange}) {
    const [listedFeatures, setListedFeatures] = useState([]);

    useEffect(() => {
        getFeatures()
        .then((data) => setListedFeatures(data.features))
        .catch((err) => console.error(err));
    }, []);

    function getFeatId(featName) {
        return listedFeatures.find((feat) => feat.feat_name === featName)?.feat_id;
    }

    function toggleFeatureId(featId) {
        const numberId = Number(featId);
        if (selectedFeatures.some((feat) => feat.feat_id === numberId)) {
            onFilterChange(selectedFeatures.filter((feat) => feat.feat_id !== numberId));
        } else {
            onFilterChange([...selectedFeatures, {feat_id: numberId, value: null}]);
        }
    }


    function toggleFeatureIdAndValue(featId, value) {
        const numberId = Number(featId);
        const numberValue = parseInt(value, 10);
        const filteredFeat = selectedFeatures.filter((feat) => feat.feat_id !== numberId);
        if (selectedFeatures.some((feat) => feat.feat_id === numberId && feat.value === numberValue)) {
            onFilterChange(filteredFeat)
        } else {
            onFilterChange([...filteredFeat, {feat_id: numberId, value: numberValue}]);
        }
    }

    function isValueSelected(featId, value) {
        return selectedFeatures.some(
            (feat) => feat.feat_id === Number(featId) && feat.value === parseInt(value, 10)
        )
    }

    function isFeatureIdSelected(featId) {
        return selectedFeatures.some((feat) => feat.feat_id === Number(featId));
    }

    function clearAll(){
        onFilterChange([])
    }

    const noiseFeatId = getFeatId("noise_level");
    const priceFeatId = getFeatId("price");
    const openLateFeatId = getFeatId("open_late")

    return (
        <div className="filter-bar">
      <div className="filter-bar__header">
        <SlidersHorizontal size={19} />
        <span>Filter</span>
      </div>

      <div className="filter-bar__section">
        <h3 className="filter-bar__section-title">Amenities</h3>
        <div className="filter-bar__options">
          {Object.entries(amenitiesFeatureIcons).map(([featName, Icon]) => {
            const featId = getFeatId(featName);
            if (!featId) return null;
            return (
              <button
                key={featName}
                type="button"
                className={`filter-btn ${isFeatureIdSelected(featId) ? "selected" : ""}`}
                onClick={() => toggleFeatureId(featId)}
              >
                <Icon size={16} />
                <span>{amenitiesFeatureLabels[featName]}</span>
              </button>
            );
          })}
        </div>
      </div>

      <div className="filter-bar__section">
        <h3 className="filter-bar__section-title">Noise Level</h3>
        <div className="filter-bar__options">
          {noiseFeatId && noiseLevelOptions.map(({ label, value, Icon }) => (
            <button
              key={value}
              type="button"
              className={`filter-btn ${isValueSelected(noiseFeatId, value) ? "selected" : ""}`}
              onClick={() => toggleFeatureIdAndValue(noiseFeatId, value)}
            >
              <Icon size={16} />
              <span>{label}</span>
            </button>
          ))}
        </div>
      </div>

      <div className="filter-bar__section">
        <h3 className="filter-bar__section-title">Price</h3>
        <div className="filter-bar__options">
          {priceFeatId && priceOptions.map(({ label, value }) => (
            <button
              key={value}
              type="button"
              className={`filter-btn ${isValueSelected(priceFeatId, value) ? "selected" : ""}`}
              onClick={() => toggleFeatureIdAndValue(priceFeatId, value)}
            >
              <span>{label}</span>
            </button>
          ))}
        </div>
      </div>

      <div className="filter-bar__section">
        <h3 className="filter-bar__section-title">Open Late</h3>
        <div className="filter-bar__options">
          {openLateFeatId && (
            <button
              type="button"
              className={`filter-btn ${isFeatureIdSelected(openLateFeatId) ? "selected" : ""}`}
              onClick={() => toggleFeatureId(openLateFeatId)}
            >
              <Moon size={16} />
              <span>Open Late</span>
            </button>
          )}
        </div>
      </div>
      <div className="clear-section">
        <button
        onClick={clearAll}
        >Clear All</button>
      </div>
    </div>
  );
}

