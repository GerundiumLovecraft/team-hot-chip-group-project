const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const getFeatures = async () => {

    const requestOptions = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    };

    const response = await fetch(`${BACKEND_URL}/features`, requestOptions);

    if (!response.ok) {
        throw new Error("Server issue. Try again later!");
    }

    const data = await response.json();

    return data;
};
