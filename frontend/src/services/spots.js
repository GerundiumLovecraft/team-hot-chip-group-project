// docs: https://vitejs.dev/guide/env-and-mode.html
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const getAllSpots = async () => {
    const requestOptions = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    };

    const response = await fetch(`${BACKEND_URL}/spots`, requestOptions);

    if (response.status !== 200) {
        throw new Error("Unable to fetch all spots!");
    }

    const data = await response.json();
    return data;
};

export const getSpotById = async (id) => {
    const requestOptions = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    };

    const response = await fetch(`${BACKEND_URL}/spots/${id}`, requestOptions);

    if (response.status !== 200) {
        throw new Error("Unable to find that spot!");
    }

    const data = await response.json();
    return data.Spot;
};

export const createSpot = async (spotData, token) => {
    const requestOptions = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify(spotData),
    };

    const response = await fetch(`${BACKEND_URL}/spots`, requestOptions);

    if (response.status !== 201) {
        throw new Error("Unable to create a spot!");
    }

    const data = await response.json();
    return data
};