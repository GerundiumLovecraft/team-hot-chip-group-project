import "./SpotModal.css";

const SpotModal = ({ spot, onClose }) => {
    if (!spot) return null;

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                <button className="modal-close" onClick={onClose}>✕</button>
                <img src={spot.image} alt={spot.name} className="modal-image" />
                <div className="modal-details">
                    <h2>{spot.name}</h2>
                    <p className="modal-address">{spot.address}</p>
                    <p className="modal-description">{spot.description}</p>
                    <p className="modal-hours">Opening Hours: {spot.open_from} - {spot.open_to}</p>
                    <div className="modal-features">
                        <h3>Amenities</h3>
                        <ul>
                            {spot.features.map((feature) => (
                                <li key={feature.feat_name}>
                                    <span className="feature-name">{feature.feat_name.replace(/_/g, ' ')}</span>
                                    <span className="feature-value">
                                        {feature.feat_name === 'price' ? '£'.repeat(feature.value) :
                                         feature.feat_name === 'noise_level' ? 
                                            feature.value === 1 ? 'Quiet' : 
                                            feature.value === 2 ? 'Moderate' : 'Loud'
                                         : 'Available'}
                                    </span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SpotModal;