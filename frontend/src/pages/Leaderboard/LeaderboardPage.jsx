import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getLeaderboard } from "../../services/leaderboards";
import "./LeaderboardPage.css";
import { isTokenValid, getTokenUserId } from "../../helpers/authentication.js";

const LeaderboardPage = () => {
    const [leaderboard, setLeaderboard] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();
    const currentUserId = getTokenUserId();

    useEffect(() => {
        if(!isTokenValid()) {
            navigate("/login");
            return;
        }

        const token = localStorage.getItem("token");

        getLeaderboard(token)
            .then((data) => {
                setLeaderboard(data.leaderboard);
            })
            .catch((err) => {
                console.error(err);
                localStorage.removeItem("token");
                navigate("/login");
            })
            .finally(() => setLoading(false));
    }, [navigate]);

    const getAvatar = (avatar, username) => {
        if (avatar && avatar.length > 0) {
            return <img src={avatar} alt={username} className="avatar-img" />;
        }
        return <img src={`https://api.dicebear.com/7.x/notionists-neutral/svg?seed=${username}&backgroundColor=f2cbb6,C1CF9B,F08456,F2E8B6,b8967d`} alt={username} className="avatar-img" />;
    };

    const isCurrentUser = (userId) => {
        return userId === Number(currentUserId);
    };

    const renderAvatar = (entry, sizeClass) => {
        const isOwn = isCurrentUser(entry.user_id);
        return (
            <div className={`avatar ${sizeClass} ${isOwn ? "own-entry" : ""}`}>
                {isOwn
                    ? <a href="/profile">{getAvatar(entry.avatar, entry.username)}</a>
                    : getAvatar(entry.avatar, entry.username)
                }
            </div>
        );
    };

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
                                {renderAvatar(leaderboard[1], "")}
                                <p className="podium-medal">🥈</p>
                                <p className="podium-username">{leaderboard[1].username}</p>
                                <p className="podium-spots">{leaderboard[1].spots_created} {leaderboard[1].spots_created === 1 ? "spot" : "spots"}</p>
                            </div>
                        )}
                        {leaderboard[0] && (
                            <div className="podium-entry first">
                                {renderAvatar(leaderboard[0], "large")}
                                <p className="podium-medal">🏆</p>
                                <p className="podium-username">{leaderboard[0].username}</p>
                                <p className="podium-spots">{leaderboard[0].spots_created} {leaderboard[0].spots_created === 1 ? "spot" : "spots"}</p>
                            </div>
                        )}
                        {leaderboard[2] && (
                            <div className="podium-entry third">
                                {renderAvatar(leaderboard[2], "")}
                                <p className="podium-medal">🥉</p>
                                <p className="podium-username">{leaderboard[2].username}</p>
                                <p className="podium-spots">{leaderboard[2].spots_created} {leaderboard[2].spots_created === 1 ? "spot" : "spots"}</p>
                            </div>
                        )}
                    </div>

                    {/* List — 4th place onwards */}
                    {leaderboard.length > 3 && (
                        <div className="leaderboard-list">
                            {leaderboard.slice(3).map((entry, index) => (
                                <div key={entry.user_id} className={`leaderboard-entry ${isCurrentUser(entry.user_id) ? "own-entry" : ""}`}>
                                    <span className="rank">#{index + 4}</span>
                                    <div className="list-avatar">
                                        {isCurrentUser(entry.user_id)
                                            ? <a href="/profile">{getAvatar(entry.avatar, entry.username)}</a>
                                            : getAvatar(entry.avatar, entry.username)
                                        }
                                    </div>
                                    <span className="username">{entry.username}</span>
                                    <span className="spots-created">{entry.spots_created} {entry.spots_created === 1 ? "spot" : "spots"} submitted</span>
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