// docs: https://vitejs.dev/guide/env-and-mode.html
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const login = async (usernameOrEmail, password) => {
  const payload = {
    usernameOrEmail: usernameOrEmail,
    password: password,
  };

  const requestOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  };

  const response = await fetch(`${BACKEND_URL}/tokens`, requestOptions);

  // docs: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/201
  if (response.status === 201) {
    let data = await response.json();
    return data.token;
  } else {
    throw new Error(
      `Received status ${response.status} when logging in. Expected 201`
    );
  }
};

export const signup = async (username, email, password) => {
  const payload = {
    username: username,
    email: email,
    password: password,
  };

  const requestOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  };

  let response = await fetch(`${BACKEND_URL}/users`, requestOptions);
  const data = await response.json();

  // docs: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/201
  if (response.status === 201) {
    return data;
  } else {
    throw new Error(data.message || `Received status ${response.status} when signing up. Expected 201`);
  }
};

export const fetchUserProfile = async (token) => {
  const requestOptions = {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  };

  const response = await fetch(`${BACKEND_URL}/profile`, requestOptions);

  if (response.status === 200) {
    const data = await response.json();
    return { user: data.user, token: data.token };
  } else if (response.status === 401) {
    throw new Error("Unauthorised: Please log in again.");
  } else {
    throw new Error(
      `Received status ${response.status} when fetching profile. Expected 200`
    );
  }
};

export const updateAvatar = async (token, avatarUrl) => {
    const response = await fetch(`${BACKEND_URL}/users/avatar`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ avatar: avatarUrl }),
    });

    if (response.status !== 200) {
        throw new Error("Unable to update avatar!");
    }

    const data = await response.json();
    return data;
};