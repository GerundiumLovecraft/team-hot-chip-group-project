import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { getLeaderboard } from "../../src/services/leaderboards.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

// Create fetch mock
createFetchMock(vi).enableMocks();

describe("leaderboards service", () => {
    describe("getLeaderboard", () => {
        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({
                "leaderboard": [
                    { "user_id": 1, "username": "userA", "spots_created": 4 },
                    { "user_id": 2, "username": "userB", "spots_created": 2 }
                ]
            }));

            const token = "test-token";

            await getLeaderboard(token)

            const fetchArguments = fetch.mock.lastCall;
            const url = fetchArguments[0]
            const options = fetchArguments[1]

            expect(url).toEqual(`${BACKEND_URL}/leaderboard`);
            expect(options.method).toEqual("GET");
            expect(options.headers["Content-Type"]).toEqual("application/json");
            expect(options.headers.Authorization).toEqual(`Bearer ${token}`);
        })
        test("returns array of users for the leaderboard", async () => {
            fetch.mockResponseOnce(JSON.stringify({
                "leaderboard": [
                    { "user_id": 1, "username": "userA", "spots_created": 4 },
                    { "user_id": 2, "username": "userB", "spots_created": 2 }
                ]
            }));

            const token = "test-token";

            const data = await getLeaderboard(token)

            expect(data.leaderboard.length).toEqual(2);
            expect(data.leaderboard[0].username).toEqual("userA");
        })
        test("throws error if the request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({
                message: "Server issue. Try again later."
            }), { status: 500 });

            try{
                await getLeaderboard("test-token");
            } catch (e) {
                expect(e.message).toEqual("Server error. Try again later");
            }
        })
    })
})