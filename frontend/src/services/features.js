const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const getAllFeatures = async () => {
    const requestOptions = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    };

    const response = await fetch(`${BACKEND_URL}/spots/features`, requestOptions);

    if (response.status !== 200) {
        throw new Error("Unable to fetch all features!");
    }

    const data = await response.json();
    return data;
};
