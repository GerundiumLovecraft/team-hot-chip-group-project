import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { login, signup, fetchUserProfile } from "../../src/services/authentication";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

// Mock fetch function
createFetchMock(vi).enableMocks();

describe("authentication service", () => {
  describe("login", () => {
    test("calls the backend url for a token", async () => {
      const testEmail = "test@testEmail.com";
      const testPassword = "12345678";

      fetch.mockResponseOnce(JSON.stringify({ token: "testToken" }), {
        status: 201,
      });

      await login(testEmail, testPassword);

      // This is an array of the arguments that were last passed to fetch
      const fetchArguments = fetch.mock.lastCall;
      const url = fetchArguments[0];
      const options = fetchArguments[1];

      expect(url).toEqual(`${BACKEND_URL}/tokens`);
      expect(options.method).toEqual("POST");
      expect(options.body).toEqual(
        JSON.stringify({ usernameOrEmail: testEmail, password: testPassword })
      );
      expect(options.headers["Content-Type"]).toEqual("application/json");
    });

    test("returns the token if the request was a success", async () => {
      const testEmail = "test@testEmail.com";
      const testPassword = "12345678";

      fetch.mockResponseOnce(JSON.stringify({ token: "testToken" }), {
        status: 201,
      });

      const token = await login(testEmail, testPassword);
      expect(token).toEqual("testToken");
    });

    test("throws an error if the request failed", async () => {
      const testEmail = "test@testEmail.com";
      const testPassword = "12345678";

      fetch.mockResponseOnce(JSON.stringify({ message: "Wrong Password" }), {
        status: 403,
      });

      try {
        await login(testEmail, testPassword);
      } catch (err) {
        expect(err.message).toEqual(
          "Received status 403 when logging in. Expected 201"
        );
      }
    });
  });

  describe("signup", () => {
    test("calls the backend url for a token", async () => {
      const testUsername = "testTest"
      const testEmail = "test@testEmail.com";
      const testPassword = "12345678";

      fetch.mockResponseOnce(JSON.stringify({"message": "OK", "token": "new-token"}), {
        status: 201,
      });

      await signup(testUsername, testEmail, testPassword);

      // This is an array of the arguments that were last passed to fetch
      const fetchArguments = fetch.mock.lastCall;
      const url = fetchArguments[0];
      const options = fetchArguments[1];

      expect(url).toEqual(`${BACKEND_URL}/users`);
      expect(options.method).toEqual("POST");
      expect(options.body).toEqual(
        JSON.stringify({ username: testUsername, email: testEmail, password: testPassword })
      );
      expect(options.headers["Content-Type"]).toEqual("application/json");
    });

    test("returns new token if the signup request was a success", async () => {
      const testUsername = "testTest"
      const testEmail = "test@testEmail.com";
      const testPassword = "12345678";

      fetch.mockResponseOnce(JSON.stringify({"message": "OK", "token": "new-token"}), {
        status: 201,
      });

      const responseData = await signup(testUsername, testEmail, testPassword);

      expect(responseData.token).toEqual("new-token");
    });

    test("throws an error if the request failed", async () => {
      const testEmail = "test@testEmail.com";
      const testPassword = "12345678";

      fetch.mockResponseOnce(
        JSON.stringify({ message: "User already exists" }),
        {
          status: 400,
        }
      );

      try {
        await signup(testEmail, testPassword);
      } catch (err) {
        expect(err.message).toEqual(
          "User already exists"
        );
      }
    });
  });

  describe("fetchUserProfile", () => {
    const testToken = "testToken123";

    const mockUser = {
      id: "1",
      username: "testuser",
      email: "test@example.com",
      createdAt: "2024-01-01T00:00:00Z",
      avatar: "",
    };

    test("calls the correct url with the correct options", async () => {
      fetch.mockResponseOnce(
          JSON.stringify({user: mockUser, token: "newToken"}),
          {status: 200}
      );

      await fetchUserProfile(testToken);

      const [url, options] = fetch.mock.lastCall;

      expect(url).toEqual(`${BACKEND_URL}/profile`);
      expect(options.method).toEqual("GET");
      expect(options.headers["Content-Type"]).toEqual("application/json");
      expect(options.headers["Authorization"]).toEqual(`Bearer ${testToken}`);
    });

    test("returns user and token on success", async () => {
      fetch.mockResponseOnce(
          JSON.stringify({user: mockUser, token: "newToken"}),
          {status: 200}
      );

      const data = await fetchUserProfile(testToken);

      expect(data.user).toEqual(mockUser);
      expect(data.token).toEqual("newToken");
    });

    test("throws an unauthorised error on 401", async () => {
      fetch.mockResponseOnce(
          JSON.stringify({message: "Unauthorised"}),
          {status: 401}
      );

      await expect(fetchUserProfile("badtoken")).rejects.toThrow(
          "Unauthorised: Please log in again."
      );
    });

    test("throws a descriptive error on other failure status", async () => {
        fetch.mockResponseOnce(
            JSON.stringify({message: "Server error"}),
            {status: 500}
        );

        await expect(fetchUserProfile(testToken)).rejects.toThrow(
            "Received status 500 when fetching profile. Expected 200"
        );
    });
  });
});
