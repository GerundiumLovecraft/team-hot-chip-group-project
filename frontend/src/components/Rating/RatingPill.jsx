import { useState } from "react";
import { Star } from "lucide-react";
import "./RatingPill.css";
import { isTokenValid } from "../../helpers/authentication.js";

const RatingPill = ({ rating, spotId, interactive = false, onRatingSubmit }) => {
    const [popupOpen, setPopupOpen] = useState(false);
    const [hoveredIndex, setHoveredIndex] = useState(null);

    const isLoggedIn = isTokenValid();
    const token = isLoggedIn ? localStorage.getItem("token") : null;

    const getColorTier = (r) => {
        if (r === null) return { bg: "#f4f4f5", text: "#71717a", accent: "#a1a1aa" };
        if (r >= 4.0) return { bg: "#dcfce7", text: "#166534", accent: "#22c55e" };
        if (r >= 3.0) return { bg: "#fef3c7", text: "#92400e", accent: "#f59e0b" };
        if (r >= 2.0) return { bg: "#ffedd5", text: "#9a3412", accent: "#f97316" };
        return { bg: "#fee2e2", text: "#991b1b", accent: "#ef4444" };
    };

    const tier = getColorTier(rating ?? null);

    const handleStarClick = (value) => {
        setPopupOpen(false);
        setHoveredIndex(null);
        onRatingSubmit?.(token, spotId, value);
    };

    const handlePillClick = () => {
        if (interactive && isLoggedIn) {
            setPopupOpen((prev) => !prev);
        }
    };

    return (
        <div className="sr-wrapper">
            <div
                className={`sr-pill ${interactive && isLoggedIn ? "sr-pill--clickable" : ""}`}
                style={{ background: tier.bg, color: tier.text }}
                onClick={handlePillClick}
                role={interactive && isLoggedIn ? "button" : undefined}
                tabIndex={interactive && isLoggedIn ? 0 : undefined}
                aria-label={rating != null ? `Rating: ${rating.toFixed(1)}` : "No ratings yet"}
            >
        <span className="sr-value">
          {rating != null ? rating.toFixed(1) : "—"}
        </span>
                <Star size={14} color={tier.accent} fill={tier.accent} />
            </div>

            {popupOpen && (
                <>
                    <div className="sr-backdrop" onClick={() => setPopupOpen(false)} />
                    <div
                        className="sr-popup"
                        onMouseLeave={() => setHoveredIndex(null)}
                    >
                        <p className="sr-popup-label">Rate the spot!</p>
                        <div className="sr-popup-stars">
                            {[1, 2, 3, 4, 5].map((value) => (
                                <button
                                    key={value}
                                    className={`sr-star-btn ${hoveredIndex !== null && value <= hoveredIndex ? "hovered" : ""}`}
                                    onClick={() => handleStarClick(value)}
                                    onMouseEnter={() => setHoveredIndex(value)}
                                    aria-label={`Rate ${value} star${value > 1 ? "s" : ""}`}
                                >
                                    ⭐
                                </button>
                            ))}
                        </div>
                    </div>
                </>
            )}
        </div>
    );
};

export default RatingPill;