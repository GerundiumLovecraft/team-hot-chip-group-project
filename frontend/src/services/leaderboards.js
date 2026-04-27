const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export async function getLeaderboard(token) {
    const reqOption = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
    };

    const response = await fetch(`${BACKEND_URL}/leaderboard`, reqOption);

    if (!response.ok) {
        throw new Error("Server error. Try again later");
    }

    return await response.json();
}