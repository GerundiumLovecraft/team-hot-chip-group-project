import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getLeaderboard } from "../../services/leaderboards";
import "./LeaderboardPage.css";

const LeaderboardPage = () => {
    const [leaderboard, setLeaderboard] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (!token) {
            navigate("/login");
            return;
        }
        getLeaderboard(token)
            .then((data) => {
                setLeaderboard(data.leaderboard);
            })
            .catch((err) => {
                console.error(err);
            })
            .finally(() => setLoading(false));
    }, [navigate]);

    return (
        <div className="leaderboard-page">
            <div className="leaderboard-header">
                <h1>Leaderboard</h1>
                <p>Top contributors to our community</p>
            </div>

            {loading ? (
                <p>Loading...</p>
            ) : (
                <>
                    {/* Podium — top 3 */}
                    <div className="podium">
                        {leaderboard[1] && (
                            <div className="podium-entry second">
                                <div className="avatar">{leaderboard[1].username.slice(0, 2).toUpperCase()}</div>
                                <p className="podium-rank">2nd</p>
                                <p className="podium-username">{leaderboard[1].username}</p>
                                <p className="podium-spots">{leaderboard[1].spots_created} spots</p>
                            </div>
                        )}
                        {leaderboard[0] && (
                            <div className="podium-entry first">
                                <div className="avatar large">🏆</div>
                                <p className="podium-rank">1st</p>
                                <p className="podium-username">{leaderboard[0].username}</p>
                                <p className="podium-spots">{leaderboard[0].spots_created} spots</p>
                            </div>
                        )}
                        {leaderboard[2] && (
                            <div className="podium-entry third">
                                <div className="avatar">{leaderboard[2].username.slice(0, 2).toUpperCase()}</div>
                                <p className="podium-rank">3rd</p>
                                <p className="podium-username">{leaderboard[2].username}</p>
                                <p className="podium-spots">{leaderboard[2].spots_created} spots</p>
                            </div>
                        )}
                    </div>

                    {/* List — 4th place onwards */}
                    {leaderboard.length > 3 && (
                        <div className="leaderboard-list">
                            {leaderboard.slice(3).map((entry, index) => (
                                <div key={entry.user_id} className="leaderboard-entry">
                                    <span className="rank">#{index + 4}</span>
                                    <span className="username">{entry.username}</span>
                                    <span className="spots-created">{entry.spots_created} spots submitted</span>
                                </div>
                            ))}
                        </div>
                    )}
                </>
            )}
        </div>
    );
};

export default LeaderboardPage;