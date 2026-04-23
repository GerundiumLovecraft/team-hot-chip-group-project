// docs: https://vitejs.dev/guide/env-and-mode.html
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const getSpots = async (token) => {
  const headers = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const requestOptions = {
    method: "GET",
    headers,
  };

  const response = await fetch(`${BACKEND_URL}/spots`, requestOptions);

  if (response.status !== 200) {
    throw new Error("Unable to fetch spots");
  }

  const data = await response.json();
  return data;
};

export const createSpot = async (spotData, token) => {
  const headers = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const requestOptions = {
    method: "POST",
    headers,
    body: JSON.stringify(spotData),
  };

  const response = await fetch(`${BACKEND_URL}/spots`, requestOptions);
  const data = await response.json();

  if (response.status !== 201) {
    throw new Error(data.message || "Unable to create spot");
  }

  return data;
};

export const createUser = async (email, password) => {
  const requestOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: email,
      hashed_password: password,
    }),
  };

  const response = await fetch(`${BACKEND_URL}/users`, requestOptions);

  if (response.status !== 201) {
    const data = await response.json();
    throw new Error(data.message);
  }

  return true;
};