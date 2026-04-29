// docs: https://vitejs.dev/guide/env-and-mode.html
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export async function SubmitRating(token, spotId, value){
    try {
        const requestOptions = {
            "method": "POST",
            "headers": {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`,
            },
            body: JSON.stringify({ "rating": value, })
        }

        const response = await fetch(`${BACKEND_URL}/spots/${spotId}/rate`, requestOptions);

        if (!response.ok) {
            throw new Error("Couldn't add the rating. Try again later");
        }

        return;
    } catch (e) {
        console.log(e);
    }
}